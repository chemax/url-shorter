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
	ctx        context.Context
	pingSync   sync.Mutex
	configured bool
}

var database *DB

func (db *DB) createURLsTable() error {
	//это должно миграциями делаться, но как вкрутить миграции внутрь сервиса я пока не знаю. Обычно они снаружи.
	// Снаружи я готов уже. Но задача пока такая.
	_, err := db.conn.Exec(db.ctx, `create table if not exists URLs(
  id serial primary key,
  shortCode varchar unique not null,
  URL text unique not null
);`)
	if err != nil {
		return fmt.Errorf("create table 'URLs' error: %w", err)
	}
	return nil
}
func (db *DB) Use() bool {
	return db.configured
}
func (db *DB) Get(shortcode string) (string, error) {
	var URL string
	err := db.conn.QueryRow(db.ctx, `SELECT url FROM urls WHERE shortcode = $1`, shortcode).Scan(&URL)
	if err != nil {
		return "", fmt.Errorf("query shortcode error: %w", err)
	}

	return URL, err
}
func (db *DB) SaveURL(shortcode string, URL string) (string, error) {
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
	sqlString := `with new(id,shortcode,url) as (
values
(nextval('urls_id_seq'::regclass), $1, $2) 
), dup as (
	select urls.* from urls
	where (url) in ( select url from new)
), ins as (
	insert into urls
	select * from new
	where (url) not in ( select url from dup)
	returning *
) 
select shortcode from dup ;`
	row := db.conn.QueryRow(db.ctx, sqlString, shortcode, URL)
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
	return db.conn.Ping(db.ctx)
}

func (db *DB) pingAllTime() {
	defer db.conn.Close(db.ctx)
	for {
		select {
		case <-db.ctx.Done():
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
	conn, err := pgx.Connect(db.ctx, db.url)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	db.conn = conn
	go db.pingAllTime()
	err = db.createURLsTable()
	return err
}

func Init(ctx context.Context, url string) (*DB, error) {
	if database == nil {
		database = &DB{
			url:        url,
			ctx:        ctx,
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
