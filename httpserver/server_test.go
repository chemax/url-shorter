package httpserver

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/chemax/url-shorter/config"
	"github.com/chemax/url-shorter/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		log, err := logger.NewLogger()
		assert.Nil(t, err)
		sig := make(chan os.Signal, 1)
		go func() {
			t := time.NewTicker(time.Second * 1)
			<-t.C
			sig <- syscall.SIGINT
		}()
		err = New(context.Background(), &config.Config{NetAddr: "localhost:5555"}, log, service(), sig)
		assert.ErrorIs(t, err, http.ErrServerClosed)
	})
	t.Run("httpS", func(t *testing.T) {
		log, err := logger.NewLogger()
		assert.Nil(t, err)
		sig := make(chan os.Signal, 1)
		go func() {
			t := time.NewTicker(time.Second * 1)
			<-t.C
			sig <- syscall.SIGINT
		}()
		err = New(context.Background(), &config.Config{NetAddr: "localhost:5555", HTTPSEnabled: true}, log, service(), sig)
		assert.ErrorIs(t, err, http.ErrServerClosed)
	})
}

func service() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sup"))
	})
	return r
}
