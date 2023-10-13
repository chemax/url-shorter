package db

import (
	"context"
	"fmt"
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
  shortCode varchar not null,
  URL text not null
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
func (db *DB) SaveURL(shortcode string, URL string) error {
	_, err := db.conn.Exec(db.ctx, `insert into urls(shortcode, url) values ($1, $2)`, shortcode, URL)
	if err != nil {
		return err
	}
	return nil
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
