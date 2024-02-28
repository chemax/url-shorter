package handlers

import (
	"bytes"
	"github.com/chemax/url-shorter/config"
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/storage"
	"github.com/chemax/url-shorter/logger"
	"github.com/chemax/url-shorter/models"
	"github.com/chemax/url-shorter/users"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
)

// Ещё одна попытка сделать интеграционный тест
// тут по сути будет клон  app.Run
func TestIntegrated(t *testing.T) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cfg, _ := config.NewConfig()
	cfg.DBConfig = "123"
	log, _ := logger.NewLogger()
	defer log.Shutdown()
	dbObj, _ := db.NewDB(cfg.DBConfig, log)
	mock, err := pgxmock.NewPool()
	if err != nil {
		return
	}
	defer mock.Close()
	dbObj.SetCon(mock)
	st, err := storage.NewStorage(cfg, log, dbObj)
	assert.Nil(t, err)
	usersObj, _ := users.NewUser(cfg, log, dbObj)
	handlersForTest := NewHandlers(st, cfg, log, usersObj)

	// Панеслася
	// делаем первый запрос на сокращение урла
	JSONURL1 := "{\"url\": \"http://ya.ru\"}"
	//shortCode1 := "1234"
	userID1 := "userid1"
	rows := mock.NewRows([]string{"id"}). //создание пользователя
						AddRow(userID1)
	rows2 := mock.NewRows([]string{"shortcode"}) //сохранение URL
	//AddRow(shortCode1)
	mock.ExpectQuery(`INSERT INTO users .+`).WillReturnRows(rows)
	mock.ExpectQuery(`.+`).WithArgs(pgxmock.AnyArg(), "http://ya.ru", userID1).WillReturnRows(rows2)
	rand.New(rand.NewSource(1234))
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONURL1)))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlersForTest.Router.ServeHTTP(w, request)
	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.True(t, strings.Contains(res.Cookies()[0].String(), models.TokenCookieName))
	res.Body.Close()
}
