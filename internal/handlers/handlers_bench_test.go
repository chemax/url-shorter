package handlers

import (
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/logger"
	"github.com/chemax/url-shorter/internal/users"
	mock_util "github.com/chemax/url-shorter/mocks/storage"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func BenchmarkHandlers_PostHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	cfg := mock_util.NewMockConfigInterface(ctrl)
	cfg.EXPECT().SecretKey().Return("false").AnyTimes()
	cfg.EXPECT().TokenExp().Return(1 * time.Hour).AnyTimes()
	cfg.EXPECT().GetHTTPAddr().Return("http://127.0.0.1:8080").AnyTimes()
	st := mock_util.NewMockStorageInterface(ctrl)
	st.EXPECT().GetURL(gomock.Any()).AnyTimes()
	log, _ := logger.Init()
	bd, _ := db.Init("", log)
	usersManager, _ := users.Init(cfg, log, bd)
	handlers := New(st, cfg, log, usersManager)
	JSONURL := "{\"url\": \"http://ya.ru\"}"
	JSONURLArray := "[{\"correlation_id\":\"1\",\"original_url\": \"http://ya.ru\"}, {\"correlation_id\":\"2\",\"original_url\": \"http://ya.ru\"}]"
	JSONBadURL := "{\"url\": \".ru\"}"
}
