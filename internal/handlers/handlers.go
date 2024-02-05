package handlers

import (
	"fmt"
	"net/http"

	"github.com/chemax/url-shorter/compress"
	"github.com/chemax/url-shorter/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// configer интерфейс конфиг-структуры
type configer interface {
	GetHTTPAddr() string
}

// userser интерфейс юзер-менеджера
type userser interface {
	Middleware(next http.Handler) http.Handler
}

// loggerer интерфейс логера
type loggerer interface {
	Middleware(next http.Handler) http.Handler
	Warn(args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
}

// storager интерфейс хранилища
type storager interface {
	GetUserURLs(userID string) ([]util.URLWithShort, error)
	GetURL(code string) (parsedURL string, err error)
	DeleteListFor(forDelete []string, userID string)
	AddNewURL(parsedURL string, userID string) (code string, err error)
	Ping() bool
	BatchSave(arr []*util.URLForBatch, httpPrefix string) (responseArr []util.URLForBatchResponse, err error)
}

type handlers struct {
	storage storager
	Router  *chi.Mux
	Cfg     configer
	Log     loggerer
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

// New возвращает хендлер всех ручек
func New(s storager, cfg configer, log loggerer, users userser) *handlers {
	initRender()
	r := chi.NewRouter()
	h := &handlers{
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
	r.Use(users.Middleware)
	r.Use(compress.Middleware)
	r.Post("/api/shorten", h.xJSONPostHandler)
	r.Post("/api/shorten/batch", h.xJSONBatchPostHandler)
	r.Post("/", h.postHandler)
	r.Get("/ping", h.pingHandler)
	r.Get("/{id}", h.getHandler)
	r.Get("/api/user/urls", h.getUserURLsHandler)
	r.Delete("/api/user/urls", h.DeleteUserURLsHandler)
	return h
}

//
