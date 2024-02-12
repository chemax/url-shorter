// Package db синглтон для общения с БД.
package db

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/chemax/url-shorter/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Loggerer интерфейс логера
type Loggerer interface {
	Warnln(args ...interface{})
	Error(args ...interface{})
}

// DeleteTask задача на удаление
type DeleteTask struct {
	Codes  []string
	UserID string
}

type managerDB struct {
	conn       *pgxpool.Pool
	url        string
	pingSync   sync.Mutex
	configured bool
	delete     chan DeleteTask
	log        Loggerer
}

var database *managerDB

func (db *managerDB) createURLsTable() error {
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

func (db *managerDB) backgroundDeleteHandler() {
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
			db.log.Error("batch delete managerDB.conn.Acquire error %w", err)
			continue
		}

		_, err = conn.Query(context.Background(), buf.String(), task.UserID)
		conn.Release()
		if err != nil {
			db.log.Error("batch delete error %w", err)
			continue
		}

	}
	db.log.Warnln("managerDB.delete channel was closed")
}

// BatchDelete принимает пакет ид для пакетного удаления и юзерИД
func (db *managerDB) BatchDelete(forDelete []string, userID string) {
	db.delete <- DeleteTask{
		Codes:  forDelete,
		UserID: userID,
	}
}

// Use используется ли БД постгре
func (db *managerDB) Use() bool {
	return db.configured
}

// GetAllURLs возвращает все сокращенные URL пользователя
func (db *managerDB) GetAllURLs(userID string) ([]models.URLWithShort, error) {
	var URLs []models.URLWithShort
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("managerDB.conn.Acquire error: %w", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), `SELECT url, shortcode FROM urls WHERE userid = $1 AND deleted = false`, userID)
	if err != nil {
		return nil, fmt.Errorf("query URLs get error: %w", err)
	}
	for rows.Next() {
		url := models.URLWithShort{}
		err := rows.Scan(&url.URL, &url.Shortcode)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		URLs = append(URLs, url)
	}

	return URLs, err
}

// Get возвращает URL по коду
func (db *managerDB) Get(shortcode string) (string, error) {
	var URL string
	var deleted bool
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return "", fmt.Errorf("managerDB.conn.Acquire error: %w", err)
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), `SELECT url, deleted FROM urls WHERE shortcode = $1`, shortcode).Scan(&URL, &deleted)
	if err != nil {
		return "", fmt.Errorf("query shortcode error: %w", err)
	}
	if deleted {
		return "", models.ErrMissingContent
	}
	return URL, err
}

// SaveURL сохраняет сокращенный урл
func (db *managerDB) SaveURL(shortcode string, URL string, userID string) (string, error) {
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
		return "", fmt.Errorf("SaveURL managerDB.conn.Acquire error %w", err)
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
	return rowString, &models.AlreadyHaveThisURLError{}

}

// CreateUser создает в бд нового пользователя
func (db *managerDB) CreateUser() (string, error) {
	sqlString := `INSERT INTO users values(default) RETURNING id;`
	conn, err := db.conn.Acquire(context.Background())
	if err != nil {
		return "", fmt.Errorf("CreateUser managerDB.conn.Acquire error %w", err)
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
func (db *managerDB) Ping() error {
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
		return fmt.Errorf("ping managerDB.conn.Acquire error %w", err)
	}
	defer conn.Release()
	return conn.Ping(context.Background())
}

func (db *managerDB) pingAllTime() {
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

func (db *managerDB) connect() error {
	if db.url == "" {
		return nil
	}
	config, err := pgxpool.ParseConfig(db.url)
	config.MaxConns = 10
	if err != nil {
		return fmt.Errorf("connect ParseConfig error: %w", err)
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), config) //pgx.Connect(context.Background(), managerDB.url)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	db.conn = conn
	go db.pingAllTime()
	err = db.createURLsTable()
	return err
}

// NewDB Синглтон, возвращает ссылку на структуру для работы с постгре
func NewDB(url string, log Loggerer) (*managerDB, error) {
	if database == nil {
		database = &managerDB{
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
			database.delete = make(chan DeleteTask)
			go database.backgroundDeleteHandler()
		}
	}

	return database, nil
}
