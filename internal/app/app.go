package app

import (
	"context"
	"fmt"
	"github.com/chemax/url-shorter/internal/config"
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/logger"
	"github.com/chemax/url-shorter/internal/storage"
	"net/http"
)

func Run() error {
	ctx := context.Background()
	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("error init config: %w", err)
	}
	log, err := logger.Init()
	if err != nil {
		return fmt.Errorf("error setup logger: %w", err)
	}
	defer log.Shutdown()
	//TODO возможно стоит использовать только интерфейс хранилища с конфигом внутри и там внутри уже разбираться
	//Кто кого и как инициализирует и использует...

	dbObj, err := db.Init(ctx, cfg.DBConfig.String())
	st, err := storage.Init(cfg.SavePath.String(), log, dbObj)
	if err != nil {
		return fmt.Errorf("error storage init: %w", err)
	}
	handler := handlers.New(st, cfg, log)
	err = http.ListenAndServe(cfg.GetNetAddr(), handler.Router)
	return err
}
