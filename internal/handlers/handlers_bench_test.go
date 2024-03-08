package handlers

import (
	"bytes"
	"fmt"
	"github.com/chemax/url-shorter/config"
	mock_handlers "github.com/chemax/url-shorter/mocks/handlers"
	"github.com/chemax/url-shorter/users"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chemax/url-shorter/logger"

	"github.com/chemax/url-shorter/internal/db"
	"github.com/golang/mock/gomock"
)

func BenchmarkHandlers_PostHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	cfg, _ := config.NewConfig()
	st := mock_handlers.NewMockStorager(ctrl)
	st.EXPECT().GetURL(gomock.Any()).AnyTimes()
	st.EXPECT().AddNewURL(gomock.Any(), gomock.Any()).Return("12345678", nil).AnyTimes()

	log, _ := logger.NewLogger()
	bd, _ := db.NewDB("", log)
	usersManager, _ := users.NewUser(cfg, log, bd)
	handlers := NewHandlers(st, cfg, log, usersManager)
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
