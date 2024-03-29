// Package httpserver создаёт новый http(s) сервер
package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/chemax/url-shorter/certgen"
	"github.com/chemax/url-shorter/config"
)

// Loggerer интерфейс логера
type Loggerer interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

// New создаёт новый http(s) сервер
func New(ctx context.Context, cfg *config.Config, log Loggerer, r http.Handler, sig chan os.Signal) error {
	server := http.Server{
		IdleTimeout: time.Second * 30,
		ReadTimeout: time.Second * 30,
		Addr:        cfg.GetNetAddr(),
		Handler:     r,
	}
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	go func() {
		// https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
		<-sig

		shutdownCtx, shutDownCancelFunc := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutDownCancelFunc() // подавись, линтер

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err) // если шатдаун не прошёл - кладём сервер принудительно, а то зависнет
		}
		serverStopCtx()
	}()
	var serverErr error
	if cfg.HTTPSEnabled {
		log.Info("start https server")
		// с одной стороны, стоит передавать в генератор серта свой хост. из конфига.
		// с другой стороны, в нормальной ситуации такой фигней заниматься не придётся, серт будет рядом лежать.
		// девопсы, let's encrypt или ещё кто
		// да и я бы предпочел, для вебсервера тупого, прикрыть его nginx'ом. нежели поднимать на нём хттпс. просто удобней.
		c1, c2 := certgen.NewCert()
		var pair tls.Certificate
		pair, err := tls.X509KeyPair(c1.Bytes(), c2.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{pair},
		}
		serverErr = server.ListenAndServeTLS("", "")
	} else {
		log.Info("start http server")
		serverErr = server.ListenAndServe()
	}
	<-serverCtx.Done()
	return serverErr
}
