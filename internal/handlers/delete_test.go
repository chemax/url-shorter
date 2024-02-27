package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chemax/url-shorter/config"
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/logger"
	mock_util "github.com/chemax/url-shorter/mocks/storage"
	"github.com/chemax/url-shorter/models"
	"github.com/chemax/url-shorter/users"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_DeleteUserURLsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	JSONURL := "{\"url\": \"http://ya.ru\"}"
	URLSForDelete := "[\"1234\"]"
	newUrlCode := "1234"
	userUrls := []models.URLWithShort{{Shortcode: newUrlCode, URL: "http://ya.ru"}}
	t.Run("1", func(t *testing.T) {
		cfg, _ := config.NewConfig()
		st := mock_util.NewMockStorageInterface(ctrl)
		st.EXPECT().GetURL(gomock.Any()).AnyTimes()

		log, _ := logger.NewLogger()
		bd, _ := db.NewDB("", log)
		usersManager, _ := users.NewUser(cfg, log, bd)
		handlersForTest := NewHandlers(st, cfg, log, usersManager)
		assert.NotNil(t, handlersForTest)

		st.EXPECT().AddNewURL(gomock.Any(), gomock.Any()).AnyTimes().Return(newUrlCode, nil)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONURL)))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handlersForTest.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		st.EXPECT().GetUserURLs(gomock.Any()).AnyTimes().Return(userUrls, nil)
		request2 := httptest.NewRequest(http.MethodGet, "/api/user/urls", bytes.NewBuffer([]byte(JSONURL)))
		request2.AddCookie(res.Cookies()[0])
		request2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		handlersForTest.Router.ServeHTTP(w2, request2)
		res2 := w2.Result()
		defer res2.Body.Close()
		assert.Equal(t, http.StatusOK, res2.StatusCode)

		st.EXPECT().DeleteListFor(gomock.Any(), gomock.Any()).Times(1)
		request3 := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer([]byte(URLSForDelete)))
		request3.AddCookie(res.Cookies()[0])
		request3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		handlersForTest.Router.ServeHTTP(w3, request3)
		res3 := w3.Result()
		defer res3.Body.Close()
		assert.Equal(t, http.StatusAccepted, res3.StatusCode)

		request4 := httptest.NewRequest(http.MethodDelete, "/api/user/urls", nil)
		request4.AddCookie(res.Cookies()[0])
		request4.Header.Set("Content-Type", "application/json")
		w4 := httptest.NewRecorder()
		handlersForTest.Router.ServeHTTP(w4, request4)
		res4 := w4.Result()
		defer res4.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res4.StatusCode)

		request5 := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer([]byte("{{")))
		request5.AddCookie(res.Cookies()[0])
		request5.Header.Set("Content-Type", "application/json")
		w5 := httptest.NewRecorder()
		handlersForTest.Router.ServeHTTP(w5, request5)
		res5 := w5.Result()
		defer res5.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res5.StatusCode)
	})
}
