package pprofserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/chemax/url-shorter/interfaces"
	"github.com/go-chi/chi/v5"
)

func Init(ctx context.Context, log interfaces.LoggerInterface) {
	r := chi.NewRouter()
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/*", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	//todo конфигурировать это
	s := &http.Server{
		Addr:    "0.0.0.0:7777",
		Handler: r,
	}
	notifyErr := make(chan error, 1)
	go func(s *http.Server, notErr chan error) {
		notErr <- s.ListenAndServe()
	}(s, notifyErr)

	go func(ctx context.Context, notErr chan error) {
		select {
		case err := <-notErr:
			log.Error(fmt.Errorf("pprof server error: %w", err))
		case <-ctx.Done():
			s.Shutdown(ctx)
		}
		log.Info("pprof server stopped")
	}(ctx, notifyErr)
}
