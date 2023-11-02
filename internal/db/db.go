package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/chemax/url-shorter/util"
	"github.com/jackc/pgx/v5"
	"sync"
	"time"
)

type DB struct {
	conn       *pgx.Conn
	url        string
	pingSync   sync.Mutex
	configured bool
}

var database *DB

func (db *DB) createURLsTable() error {
	//это должно миграциями делаться, но как вкрутить миграции внутрь сервиса я пока не знаю. Обычно они снаружи.
	// Снаружи я готов уже. Но задача пока такая.
	_, err := db.conn.Exec(context.Background(), `create table if not exists URLs(
  id serial primary key,
  shortCode varchar unique not null,
  URL text unique not null,
  userID varchar
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
func (db *DB) Use() bool {
	return db.configured
}
func (db *DB) GetAllURLs(userID string) ([]util.URLStructUser, error) {
	var URLs []util.URLStructUser
	rows, err := db.conn.Query(context.Background(), `SELECT url, shortcode FROM urls WHERE userid = $1`, userID)
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
func (db *DB) Get(shortcode string) (string, error) {
	var URL string
	err := db.conn.QueryRow(context.Background(), `SELECT url FROM urls WHERE shortcode = $1`, shortcode).Scan(&URL)
	if err != nil {
		return "", fmt.Errorf("query shortcode error: %w", err)
	}

	return URL, err
}
func (db *DB) SaveURL(shortcode string, URL string, userID string) (string, error) {
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
	row := db.conn.QueryRow(context.Background(), sqlString, shortcode, URL, userID)
	var rowString string
	err := row.Scan(&rowString)
	if err != nil {
		//Ошибка это хорошо, конкретно эта. Она означает отсутствие дюпа.
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return rowString, &util.AlreadyHaveThisURLError{}

}
func (db *DB) CreateUser() (string, error) {
	sqlString := `INSERT INTO users values(default) RETURNING id;`
	row := db.conn.QueryRow(context.Background(), sqlString)
	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("error create new user: %w", err)
	}
	return id, nil
}

func (db *DB) Ping() error {
	// мьютекс нужен, потому что я использую пинг для реконнекта
	// если мьютекса нет, возможны коллизии и ошибка conn busy
	// потом, возможно, эту проблему решит пул коннектов
	// пока он один, я сделал так
	db.pingSync.Lock()
	defer db.pingSync.Unlock()
	if db.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	return db.conn.Ping(context.Background())
}

func (db *DB) pingAllTime() {
	defer db.conn.Close(context.Background())
	for {
		select {
		case <-context.Background().Done():
			return
		default:
			var err error
			<-time.After(500 * time.Millisecond)
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

}

func (db *DB) connect() error {
	if db.url == "" {
		return nil
	}
	conn, err := pgx.Connect(context.Background(), db.url)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	db.conn = conn
	go db.pingAllTime()
	err = db.createURLsTable()
	return err
}

func Init(url string) (*DB, error) {
	if database == nil {
		database = &DB{
			url:        url,
			configured: false,
		}
		if url != "" {
			err := database.connect()
			if err != nil {
				return nil, fmt.Errorf("database init error: %w", err)
			}
			database.configured = true
		}
	}

	return database, nil
}
