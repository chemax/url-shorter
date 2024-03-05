// Package handlers роутер
package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/chemax/url-shorter/compress"
	"github.com/chemax/url-shorter/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// configer интерфейс конфиг-структуры
type configer interface {
	GetHTTPAddr() string
	GetTrustedSubnet() string
}

// Userser интерфейс юзер-менеджера
type Userser interface {
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
	GetUserURLs(userID string) ([]models.URLWithShort, error)
	GetURL(code string) (parsedURL string, err error)
	DeleteListFor(forDelete []string, userID string)
	AddNewURL(parsedURL string, userID string) (code string, err error)
	Ping() bool
	BatchSave(arr []*models.URLForBatch, httpPrefix string) (responseArr []models.URLForBatchResponse, err error)
	GetStats() (models.Stats, error)
}

type handlers struct {
	storage       storager
	Router        *chi.Mux
	Cfg           configer
	Log           loggerer
	TrustedSubnet *net.IPNet
}

func checkHeaderIsValidType(header string) bool {
	return strings.Contains(header, "application/json") || strings.Contains(header, "application/x-gzip")
}

func checkHeader(header string) bool {
	return strings.Contains(header, "text/plain") || strings.Contains(header, "application/x-gzip")
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

// NewHandlers возвращает хендлер всех ручек
func NewHandlers(s storager, cfg configer, log loggerer, users Userser) *handlers {
	initRender()
	r := chi.NewRouter()
	h := &handlers{
		storage: s,
		Router:  r,
		Cfg:     cfg,
		Log:     log,
	}
	_, network, err := net.ParseCIDR(cfg.GetTrustedSubnet())
	if err != nil {
		h.Log.Warn(fmt.Errorf("WARN err setup trusted subnet: %w", err))
	} else {
		h.TrustedSubnet = network
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
	r.Get("/api/internal/stats", h.statHandler)
	r.Delete("/api/user/urls", h.DeleteUserURLsHandler)
	return h
}

//
