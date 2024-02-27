package db

import (
	"fmt"
	"github.com/chemax/url-shorter/logger"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDB(t *testing.T) {
	t.Run("empty conn string", func(t *testing.T) {
		log, _ := logger.NewLogger()
		db, err := NewDB("", log)
		assert.Nil(t, err)
		assert.NotNil(t, db)
		assert.False(t, db.Use())
	})
	t.Run("create table good", func(t *testing.T) {
		log, _ := logger.NewLogger()
		db, err := NewDB("", log)
		assert.Nil(t, err)
		assert.NotNil(t, db)
		mock, err := pgxmock.NewPool()
		defer mock.Close()
		assert.Nil(t, err)
		db.conn = mock
		mock.ExpectExec(`create table if not exists URLs`).WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec(`create table if not exists users`).WillReturnResult(pgxmock.NewResult("CREATE", 1))
		err = db.createURLsTable()
		assert.Nil(t, err)
	})
	t.Run("create table bad 1", func(t *testing.T) {
		log, _ := logger.NewLogger()
		db, err := NewDB("", log)
		assert.Nil(t, err)
		assert.NotNil(t, db)
		mock, err := pgxmock.NewPool()
		defer mock.Close()
		assert.Nil(t, err)
		db.conn = mock
		mock.ExpectExec(`create table if not exists URLs`).WillReturnError(fmt.Errorf("test error URLs"))
		err = db.createURLsTable()
		assert.NotNil(t, err)
	})
	t.Run("create table bad 2", func(t *testing.T) {
		log, _ := logger.NewLogger()
		db, err := NewDB("", log)
		assert.Nil(t, err)
		assert.NotNil(t, db)
		mock, err := pgxmock.NewPool()
		defer mock.Close()
		assert.Nil(t, err)
		db.conn = mock
		mock.ExpectExec(`create table if not exists URLs`).WillReturnResult(pgxmock.NewResult("CREATE", 1))
		mock.ExpectExec(`create table if not exists users`).WillReturnError(fmt.Errorf("test error users"))
		err = db.createURLsTable()
		assert.NotNil(t, err)
	})
	t.Run("db.GetAllURLs error", func(t *testing.T) {
		log, _ := logger.NewLogger()
		db, err := NewDB("", log)
		assert.Nil(t, err)
		assert.NotNil(t, db)
		mock, err := pgxmock.NewPool()
		defer mock.Close()
		assert.Nil(t, err)
		db.conn = mock
		mock.ExpectQuery(`SELECT url, shortcode FROM urls`).WillReturnError(fmt.Errorf("test error getAllURLs"))
		_, err = db.GetAllURLs("123")
		assert.NotNil(t, err)
	})
	t.Run("db.GetAllURLs", func(t *testing.T) {
		log, _ := logger.NewLogger()
		db, err := NewDB("", log)
		assert.Nil(t, err)
		assert.NotNil(t, db)
		mock, err := pgxmock.NewPool()
		defer mock.Close()
		assert.Nil(t, err)

		db.conn = mock
		rows := mock.NewRows([]string{"url", "shortcode"}).
			AddRow("xxx", "hello").
			AddRow("yyy", "world")
		mock.ExpectQuery(`SELECT url, shortcode FROM urls`).WithArgs("123").WillReturnRows(rows)

		_, err = db.GetAllURLs("123")
		assert.Nil(t, err)
	})
}
