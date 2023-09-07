package handlers

import (
	util "github.com/chemax/url-shorter/util"
	"net/http"
)

type Handlers struct {
	storage util.StorageInterface
}

func New(s util.StorageInterface) *Handlers {
	return &Handlers{storage: s}
}

func (h *Handlers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		h.serveCreate(res, req)
	} else if req.Method == http.MethodGet {
		h.serveGET(res, req)
	} else {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
