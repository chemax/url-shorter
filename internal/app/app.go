package app

import (
	"github.com/chemax/url-shorter/internal/config"
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/logger"
	"github.com/chemax/url-shorter/internal/storage"
	"net/http"
)

func Run() error {
	cfg := config.Get()
	log := logger.New()
	defer log.Sync()
	handler := handlers.New(storage.Get(cfg.SavePath.String(), log), cfg, log)
	err := http.ListenAndServe(cfg.GetNetAddr(), handler.Router)
	return err
}
