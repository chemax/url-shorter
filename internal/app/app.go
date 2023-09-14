package app

import (
	"github.com/chemax/url-shorter/internal/config"
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/storage"
	"net/http"
)

func Run() error {
	cfg := config.Get()
	handler := handlers.New(storage.Get(), cfg)
	err := http.ListenAndServe(cfg.GetNetAddr(), handler.Router)
	return err
}
