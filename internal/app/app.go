package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/chemax/url-shorter/users"
	"net/http"

	"github.com/chemax/url-shorter/logger"

	"github.com/chemax/url-shorter/internal/config"
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/storage"
	"github.com/chemax/url-shorter/pprofserver"
)

// Run точка входа в приложение
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("error init config: %w", err)
	}
	log, err := logger.Init()
	if err != nil {
		return fmt.Errorf("error setup logger: %w", err)
	}
	defer log.Shutdown()
	pprofserver.Init(ctx, log)

	dbObj, err := db.Init(cfg.DBConfig.String(), log)
	if err != nil {
		return fmt.Errorf("db init error: %w", err)
	}
	st, err := storage.Init(cfg, log, dbObj)
	if err != nil {
		return fmt.Errorf("error storage init: %w", err)
	}
	usersObj, err := users.Init(cfg, log, dbObj)
	if err != nil {
		return fmt.Errorf("error users init: %w", err)
	}

	handler := handlers.New(st, cfg, log, usersObj)
	err = http.ListenAndServe(cfg.GetNetAddr(), handler.Router)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
