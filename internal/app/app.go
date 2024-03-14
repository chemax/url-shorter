// Package app главная и единственная точка входа в программу
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chemax/url-shorter/grpcserver"
	"github.com/chemax/url-shorter/internal/gRPCHandlers"

	"github.com/chemax/url-shorter/httpserver"

	"github.com/chemax/url-shorter/config"

	"github.com/chemax/url-shorter/users"

	"github.com/chemax/url-shorter/logger"

	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/storage"
	"github.com/chemax/url-shorter/pprofserver"
)

// Run точка входа в приложение
func Run() (err error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("error init config: %w", err)
	}
	log, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("error setup logger: %w", err)
	}
	defer log.Shutdown()
	pprofserver.NewPprof(ctx, log)

	dbObj, err := db.NewDB(cfg.DBConfig, log)
	if err != nil {
		return fmt.Errorf("db init error: %w", err)
	}
	st, err := storage.NewStorage(cfg, log, dbObj)
	if err != nil {
		return fmt.Errorf("error storage init: %w", err)
	}
	usersObj, err := users.NewUser(cfg, log, dbObj)
	if err != nil {
		return fmt.Errorf("error users init: %w", err)
	}

	handler := handlers.NewHandlers(st, cfg, log, usersObj)
	grpcHandler := gRPCHandlers.New(st, cfg, log, usersObj)
	//тут надо вэйтгруппу пробросить
	go func() {
		err := grpcserver.New(ctx, cfg, log, sig, grpcHandler, usersObj)
		if err != nil {
			log.Error(err)
		}
	}()
	err = httpserver.New(ctx, cfg, log, handler.Router, sig)

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	log.Debug(err)
	return err
}
