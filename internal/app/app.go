package app

import (
	"github.com/chemax/url-shorter/internal/handlers"
	"github.com/chemax/url-shorter/internal/storage"
	"net/http"
)

func Run() error {
	s := storage.Get()
	h := handlers.New(s)

	err := http.ListenAndServe(":8080", h.Router)
	return err

}
