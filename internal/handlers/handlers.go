package handlers

import (
	"fmt"
	"github.com/chemax/url-shorter/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
)

type Handlers struct {
	storage util.StorageInterface
	Router  *chi.Mux
	Cfg     util.ConfigInterface
	Log     util.LoggerInterface
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

func New(s util.StorageInterface, cfg util.ConfigInterface, log util.LoggerInterface) *Handlers {
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
	r.Use(middleware.Recoverer)
	r.Post("/api/shorten", h.ApiServeCreate)
	r.Post("/", h.ServeCreate)
	r.Get("/{id}", h.serveGET)
	return h
}

//
