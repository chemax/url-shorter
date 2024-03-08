package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_handlers "github.com/chemax/url-shorter/mocks/handlers"
	mock_users "github.com/chemax/url-shorter/mocks/users"

	"github.com/chemax/url-shorter/users"

	"github.com/chemax/url-shorter/logger"

	"github.com/chemax/url-shorter/internal/db"
	"github.com/golang/mock/gomock"
)

func BenchmarkHandlers_PostHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	cfg := mock_users.NewMockConfiger(ctrl)
	cfgHandlers := mock_handlers.NewMockConfiger(ctrl)
	cfg.EXPECT().SecretKey().Return("false").AnyTimes()
	cfg.EXPECT().TokenExp().Return(1 * time.Hour).AnyTimes()
	cfgHandlers.EXPECT().GetHTTPAddr().Return("http://127.0.0.1:8080").AnyTimes()
	st := mock_handlers.NewMockStorager(ctrl)
	st.EXPECT().GetURL(gomock.Any()).AnyTimes()
	st.EXPECT().AddNewURL(gomock.Any(), gomock.Any()).Return("12345678", nil).AnyTimes()

	log, _ := logger.NewLogger()
	bd, _ := db.NewDB("", log)
	usersManager, _ := users.NewUser(cfg, log, bd)
	handlers := NewHandlers(st, cfgHandlers, log, usersManager)
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
