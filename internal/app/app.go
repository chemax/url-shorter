package app

import (
	"fmt"
	"github.com/chemax/url-shorter/internal/config"
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/logger"
	"github.com/chemax/url-shorter/internal/storage"
	"net/http"
)

func Run() error {
	cfg,err := config.Init()
	if err != nil {
		return fmt.Errorf("error load config: %w",err)
	}
	log, err := logger.Init()
	if err != nil {
		return fmt.Errorf("error setup logger: %w",err)
	}
	defer log.Shutdown()
	st, err := storage.Init(cfg.SavePath.String(), log)
	if err != nil {
		return fmt.Errorf("error storage init: %w",err)
	}
	handler := handlers.New(st, cfg, log)
	err = http.ListenAndServe(cfg.GetNetAddr(), handler.Router)
	return err
}
