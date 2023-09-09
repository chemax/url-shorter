package app

import (
	"github.com/chemax/url-shorter/internal/config"
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/storage"
	"net/http"
)

func Run() error {
	cfg := config.Get()
	s := storage.Get()
	h := handlers.New(s, cfg)

	err := http.ListenAndServe(cfg.GetNetAddr(), h.Router)
	return err

}
