package db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/chemax/url-shorter/interfaces"
	"github.com/chemax/url-shorter/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type structDB struct {
	conn       *pgxpool.Pool
	url        string
	pingSync   sync.Mutex
	configured bool
	delete     chan util.DeleteTask
	log        interfaces.LoggerInterface
}

var database *structDB

func (db *structDB) createURLsTable() error {
	//поздно переезжать на миграции
	_, err := db.conn.Exec(context.Background(), `create table if not exists URLs(
  id serial primary key,
  shortCode varchar unique not null,
  URL text unique not null,
  userID varchar,
  deleted bool default false                             
);`)
	if err != nil {
		return fmt.Errorf("create table 'URLs' error: %w", err)
	}
	_, err = db.conn.Exec(context.Background(), `create table if not exists users(
  id serial primary key
);`)
	if err != nil {
		return fmt.Errorf("create table 'users' error: %w", err)
	}
	return nil
}

func (db *structDB) backgroundDeleteHandler() {
	for task := range db.delete {
		buf := bytes.NewBufferString("UPDATE urls SET deleted = true  WHERE shortcode IN (")
		for i, v := range task.Codes {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(fmt.Sprintf("'%s'", v))
		}
		buf.WriteString(") AND userid = $1;")
		conn, err := db.conn.Acquire(context.Background())
		if err != nil {
			db.log.Error("batch delete structDB.conn.Acquire error %w", err)
			continue
		}

		_, err = conn.Query(context.Background(), buf.String(), task.UserID)
		conn.Release()
		if err != nil {
			db.log.Error("batch delete error %w", err)
			continue
		}

	}
	db.log.Warnln("structDB.delete channel was closed")
}

// BatchDelete принимает пакет ид для пакетного удаления и юзерИД
func (db *structDB) BatchDelete(forDelete []string, userID string) {
	db.delete <- util.DeleteTask{
		Codes:  forDelete,
		UserID: userID,
	}
}

// Use используется ли БД постгре
func (db *structDB) Use() bool {
	return db.configured
}

// GetAllURLs возвращает все сокращенные URL пользователя
func (db *structDB) GetAllURLs(userID string) ([]util.URLStructUser, error) {
	var URLs []util.URLStructUser
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("structDB.conn.Acquire error: %w", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), `SELECT url, shortcode FROM urls WHERE userid = $1 AND deleted = false`, userID)
	if err != nil {
		return nil, fmt.Errorf("query URLs get error: %w", err)
	}
	for rows.Next() {
		url := util.URLStructUser{}
		err := rows.Scan(&url.URL, &url.Shortcode)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		URLs = append(URLs, url)
	}

	return URLs, err
}

// Get возвращает URL по коду
func (db *structDB) Get(shortcode string) (string, error) {
	var URL string
	var deleted bool
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return "", fmt.Errorf("structDB.conn.Acquire error: %w", err)
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), `SELECT url, deleted FROM urls WHERE shortcode = $1`, shortcode).Scan(&URL, &deleted)
	if err != nil {
		return "", fmt.Errorf("query shortcode error: %w", err)
	}
	if deleted {
		return "", util.ErrMissingContent
	}
	return URL, err
}

// SaveURL сохраняет сокращенный урл
func (db *structDB) SaveURL(shortcode string, URL string, userID string) (string, error) {
	// https://stackoverflow.com/questions/34708509/how-to-use-returning-with-on-conflict-in-postgresql
	/*
		with new(id,shortcode,url) as (
		-- the rows you want to insert
		values
		(nextval('urls_id_seq'::regclass), '12345678', 'http://yandex2.ru')
		---
		), dup as (
		-- the one that already exists (conflicting key)
		select urls.* from urls
		where (url)     in ( select url from new)
		---
		), ins as (
		-- the ones to insert
		insert into urls
		select * from new
		where (url) not in ( select url from dup)
		returning *
		---
		)
		--- finally the concatenation of inserted values and old value for skipped ones
		select * from dup union all select * from ins
		---
		;
	*/
	//TODO избавиться от * в запросе //вроде избавился (см returning)
	sqlString := `with new(id,shortcode,url,userid) as (
values
(nextval('urls_id_seq'::regclass), $1, $2, $3) 
), dup as (
	select urls.* from urls
	where (url) in ( select url from new)
), ins as (
	insert into urls
	select * from new
	where (url) not in ( select url from dup)
	returning url, shortcode
) 
select shortcode from dup ;`
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return "", fmt.Errorf("SaveURL structDB.conn.Acquire error %w", err)
	}
	defer conn.Release()
	row := conn.QueryRow(context.Background(), sqlString, shortcode, URL, userID)
	var rowString string
	err = row.Scan(&rowString)
	if err != nil {
		//Ошибка это хорошо, конкретно эта. Она означает отсутствие дюпа.
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return rowString, &util.AlreadyHaveThisURLError{}

}

// CreateUser создает в бд нового пользователя
func (db *structDB) CreateUser() (string, error) {
	sqlString := `INSERT INTO users values(default) RETURNING id;`
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return "", fmt.Errorf("CreateUser structDB.conn.Acquire error %w", err)
	}
	defer conn.Release()
	row := conn.QueryRow(context.Background(), sqlString)
	var id string
	err = row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("error create new user: %w", err)
	}
	return id, nil
}

// Ping базы данных
func (db *structDB) Ping() error {
	// мьютекс нужен, потому что я использую пинг для реконнекта
	// если мьютекса нет, возможны коллизии и ошибка conn busy
	// потом, возможно, эту проблему решит пул коннектов
	// пока он один, я сделал так
	db.pingSync.Lock()
	defer db.pingSync.Unlock()
	if db.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("ping structDB.conn.Acquire error %w", err)
	}
	defer conn.Release()
	return conn.Ping(context.Background())
}

func (db *structDB) pingAllTime() {
	defer db.conn.Close()
	tickTack := time.NewTicker(500 * time.Millisecond)
	for range tickTack.C {
		var err error
		if db.conn != nil {
			err = db.Ping()
		}
		if err != nil || db.conn == nil {
			err := db.connect()
			if err != nil {
				continue
			}
		}
	}
}

func (db *structDB) connect() error {
	if db.url == "" {
		return nil
	}
	config, err := pgxpool.ParseConfig(db.url)
	config.MaxConns = 10
	if err != nil {
		return fmt.Errorf("connect ParseConfig error: %w", err)
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), config) //pgx.Connect(context.Background(), structDB.url)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	db.conn = conn
	go db.pingAllTime()
	err = db.createURLsTable()
	return err
}

// Init Синглтон, возвращает ссылку на структуру для работы с постгре
func Init(url string, log interfaces.LoggerInterface) (*structDB, error) {
	if database == nil {
		database = &structDB{
			url:        url,
			log:        log,
			configured: false,
		}
		if url != "" {
			err := database.connect()
			if err != nil {
				return nil, fmt.Errorf("database init error: %w", err)
			}
			database.configured = true
			database.delete = make(chan util.DeleteTask)
			go database.backgroundDeleteHandler()
		}
	}

	return database, nil
}
