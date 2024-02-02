package handlers

import (
	"bytes"
	"fmt"
	"github.com/chemax/url-shorter/users"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/logger"
	mock_util "github.com/chemax/url-shorter/mocks/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Example() {
	t := &testing.T{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := mock_util.NewMockConfigInterface(ctrl)
	cfg.EXPECT().SecretKey().Return("false").AnyTimes()
	cfg.EXPECT().TokenExp().Return(1 * time.Hour).AnyTimes()
	cfg.EXPECT().GetHTTPAddr().Return("http://127.0.0.1:8080").AnyTimes()
	st := mock_util.NewMockStorageInterface(ctrl)
	st.EXPECT().GetURL(gomock.Any()).AnyTimes()
	st.EXPECT().AddNewURL(gomock.Any(), gomock.Any()).AnyTimes().Return("123445", nil)
	log, _ := logger.Init()
	bd, _ := db.Init("", log)
	usersManager, _ := users.Init(cfg, log, bd)
	handlers := New(st, cfg, log, usersManager)
	assert.NotNil(t, handlers)
	JSONURL := "{\"url\": \"http://ya.ru\"}"

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONURL)))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.Router.ServeHTTP(w, request)
	res := w.Result()
	defer res.Body.Close()

	// Print the result.
	fmt.Printf("Status code: %v\n", res.StatusCode)

	// Output:
	// Status code: 201
}
