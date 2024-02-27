package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/chemax/url-shorter/certgen"
	"github.com/chemax/url-shorter/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Loggerer интерфейс логера
type Loggerer interface {
	Infoln(args ...interface{})
	Debugln(args ...interface{})
	Fatalln(args ...interface{})
	Errorln(args ...interface{})
}

func New(ctx context.Context, cfg *config.Config, log Loggerer, r http.Handler) error {
	server := http.Server{
		IdleTimeout: time.Second * 30,
		ReadTimeout: time.Second * 30,
		Addr:        cfg.GetNetAddr(),
		Handler:     r,
	}
	serverCtx, serverStopCtx := context.WithCancel(ctx)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatalln("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalln(err) // если шатдаун не прошёл - кладём сервер принудительно, а то зависнет
		}
		serverStopCtx()
	}()
	var serverErr error
	if cfg.HTTPSEnabled {
		log.Infoln("start https server")
		// с одной стороны, стоит передавать в генератор серта свой хост. из конфига.
		// с другой стороны, в нормальной ситуации такой фигней заниматься не придётся, серт будет рядом лежать.
		// девопсы, let's encrypt или ещё кто
		// да и я бы предпочел, для вебсервера тупого, прикрыть его nginx'ом. нежели поднимать на нём хттпс. просто удобней.
		c1, c2 := certgen.NewCert()
		var pair tls.Certificate
		pair, err := tls.X509KeyPair(c1.Bytes(), c2.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{pair},
		}
		serverErr = server.ListenAndServeTLS("", "")
	} else {
		log.Infoln("start http server")
		serverErr = server.ListenAndServe()
	}
	<-serverCtx.Done()
	return serverErr
}
