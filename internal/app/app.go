// Package app главная и единственная точка входа в программу
package app

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/chemax/url-shorter/certgen"

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

	dbObj, err := db.NewDB(cfg.DBConfig.String(), log)
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
	if cfg.HTTPSEnabled {
		log.Infoln("start https server")
		// с одной стороны, стоит передавать в генератор серта свой хост. из конфига.
		// с другой стороны, в нормальной ситуации такой фигней заниматься не придётся, серт будет рядом лежать.
		// девопсы, let's encrypt или ещё кто
		// да и я бы предпочел, для вебсервера тупого, прикрыть его nginx'ом. нежели поднимать на нём хттпс. просто удобней.
		c1, c2 := certgen.NewCert()
		var pair tls.Certificate
		pair, err = tls.X509KeyPair(c1.Bytes(), c2.Bytes())
		if err != nil {
			return err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{pair},
		}
		server := http.Server{
			Addr:      cfg.GetNetAddr(),
			TLSConfig: tlsConfig,
			Handler:   handler.Router,
		}
		err = server.ListenAndServeTLS("", "")
	} else {
		log.Infoln("start http server")
		err = http.ListenAndServe(cfg.GetNetAddr(), handler.Router)
	}
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	log.Debugln(err)
	return err
}
