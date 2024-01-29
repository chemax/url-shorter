package handlers

import (
	"bytes"
	"fmt"
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/logger"
	"github.com/chemax/url-shorter/internal/users"
	mock_util "github.com/chemax/url-shorter/mocks/storage"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
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
	st.EXPECT().AddNewURL(gomock.Any(), gomock.Any()).Return("12345678", nil).AnyTimes()

	log, _ := logger.Init()
	bd, _ := db.Init("", log)
	usersManager, _ := users.Init(cfg, log, bd)
	handlers := New(st, cfg, log, usersManager)
	JSONURLFmt := "{\"url\": \"http://%d.ya.ru\"}"
	var cmps int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.Run("all ok", func(b *testing.B) {
			jsString := fmt.Sprintf(JSONURLFmt, cmps)
			log.Debug(jsString)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten",
				bytes.NewBuffer([]byte(jsString)))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handlers.Router.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			cmps++
		})
	}
}
