package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/chemax/url-shorter/config"
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/storage"
	"github.com/chemax/url-shorter/logger"
	"github.com/chemax/url-shorter/models"
	"github.com/chemax/url-shorter/users"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
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
	URLSForBatch := "[{\"correlation_id\": \"1\", \"original_url\": \"youtube.com\"}, {\"correlation_id\": \"2\", \"original_url\": \"vk.com\"}]"
	//shortCode1 := "1234"
	userID1 := "userid1"
	rows := mock.NewRows([]string{"id"}). //создание пользователя
						AddRow(userID1)
	rows2 := mock.NewRows([]string{"shortcode"}) //сохранение URL
	mock.ExpectQuery(`INSERT INTO users .+`).WillReturnRows(rows)
	mock.ExpectQuery(`.+`).WithArgs(pgxmock.AnyArg(), "http://ya.ru", userID1).WillReturnRows(rows2)

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONURL1)))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlersForTest.Router.ServeHTTP(w, request)
	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.True(t, strings.Contains(res.Cookies()[0].String(), models.TokenCookieName))
	authCookie := res.Cookies()[0]
	res.Body.Close()

	mock.ExpectQuery(`.+`).WithArgs(pgxmock.AnyArg(), "http://habr.com", userID1).WillReturnRows(rows2)
	request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("http://habr.com")))
	request.Header.Set("Content-Type", "text/plain")
	request.AddCookie(authCookie)
	w = httptest.NewRecorder()
	handlersForTest.Router.ServeHTTP(w, request)
	res = w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	res.Body.Close()

	mock.ExpectQuery(`.+`).WithArgs(pgxmock.AnyArg(), "youtube.com", "").WillReturnRows(rows2)
	mock.ExpectQuery(`.+`).WithArgs(pgxmock.AnyArg(), "vk.com", "").WillReturnRows(rows2)
	request = httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer([]byte(URLSForBatch)))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(authCookie)
	w = httptest.NewRecorder()
	handlersForTest.Router.ServeHTTP(w, request)
	res2 := w.Result()
	assert.Equal(t, http.StatusCreated, res2.StatusCode)
	res2.Body.Close()

	rows3 := mock.NewRows([]string{"url", "shortcode"}). //URLs пользователя
								AddRow("1", "vk.com").
								AddRow("2", "youtube.com").
								AddRow("3", "http://ya.ru")
	mock.ExpectQuery("SELECT url, shortcode FROM urls .+").WithArgs(userID1).WillReturnRows(rows3)
	request = httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	request.AddCookie(authCookie)
	handlersForTest.Router.ServeHTTP(w, request)
	w = httptest.NewRecorder()
	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()

	//rows3 := mock.NewRows([]string{"url", "shortcode"}). //URLs пользователя
	//	AddRow("1", "vk.com").
	//	AddRow("2", "youtube.com").
	//	AddRow("3", "http://ya.ru")
	mock.ExpectQuery("SELECT url, shortcode FROM urls .+").WithArgs(userID1).WillReturnRows(rows3)
	request = httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	request.AddCookie(authCookie)
	handlersForTest.Router.ServeHTTP(w, request)
	w = httptest.NewRecorder()
	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res.Body.Close()
}