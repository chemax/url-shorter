package handlers

import (
	"github.com/chemax/url-shorter/internal/storage"
	. "github.com/chemax/url-shorter/util"
	"net/http"
)

type Handlers struct {
	storage StorageInterface
}

func New(s StorageInterface) *Handlers {
	return &Handlers{storage: s}
}

func (h *Handlers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	u := storage.Get()
	if req.Method == http.MethodPost {
		h.serveCreate(res, req, u)
	} else if req.Method == http.MethodGet {
		h.serveGET(res, req, u)
	} else {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
