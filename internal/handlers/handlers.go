package handlers

import (
	"fmt"
	util "github.com/chemax/url-shorter/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
)

type Handlers struct {
	storage util.StorageInterface
	Router  *chi.Mux
	Cfg     util.ConfigInterface
}

func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			// We set a default error status response code if one hasn't been set.
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			// We log the error
			fmt.Printf("Logging err: %s\n", err.Error())

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

func New(s util.StorageInterface, cfg util.ConfigInterface) *Handlers {

	r := chi.NewRouter()
	h := &Handlers{
		storage: s,
		Router:  r,
		Cfg:     cfg,
	}
	r.MethodNotAllowed(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/", h.ServeCreate)
	r.Get("/{id}", h.serveGET)
	return h
}

//
