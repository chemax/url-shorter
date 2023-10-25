package handlers

import (
	"fmt"
	"github.com/chemax/url-shorter/interfaces"
	"github.com/chemax/url-shorter/internal/compress"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
)

type Handlers struct {
	storage interfaces.StorageInterface
	Router  *chi.Mux
	Cfg     interfaces.ConfigInterface
	Log     interfaces.LoggerInterface
}

func initRender() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(http.StatusBadRequest)
			}
			fmt.Printf("Logging err: %s\n", err.Error())
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

func New(s interfaces.StorageInterface, cfg interfaces.ConfigInterface, log interfaces.LoggerInterface) *Handlers {
	initRender()
	r := chi.NewRouter()
	h := &Handlers{
		storage: s,
		Router:  r,
		Cfg:     cfg,
		Log:     log,
	}
	r.MethodNotAllowed(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})
	r.Use(log.Middleware)
	r.Use(compress.Middleware)
	r.Use(middleware.Recoverer)
	r.Post("/api/shorten", h.JSONPostHandler)
	r.Post("/api/shorten/batch", h.JSONBatchPostHandler)
	r.Post("/", h.PostHandler)
	r.Get("/ping", h.PingHandler)
	r.Get("/{id}", h.GetHandler)
	return h
}

//
