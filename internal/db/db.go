package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"sync"
	"time"
)

type DB struct {
	conn     *pgx.Conn
	url      string
	ctx      context.Context
	pingSync sync.Mutex
}

var db *DB

func (db *DB) Ping() error {
	db.pingSync.Lock()
	defer db.pingSync.Unlock()
	if db.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	return db.conn.Ping(db.ctx)
}

func (db *DB) pingAllTime() {
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
	db.conn.Close(db.ctx)
}

func (db *DB) connect() error {
	conn, err := pgx.Connect(db.ctx, db.url)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	db.conn = conn
	go db.pingAllTime()
	return nil
}

func Init(ctx context.Context, url string) (*DB, error) {
	if db == nil {
		db = &DB{
			url: url,
			ctx: ctx,
		}
		err := db.connect()
		if err != nil {
			return nil, fmt.Errorf("db init error: %w", err)
		}
	}

	return db, nil
}
