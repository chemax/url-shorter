package handlers

import (
	"github.com/chemax/url-shorter/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handlers) serveGET(res http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if len(id) != util.CodeLength {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	parsedURL, err := h.storage.GetURL(id)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", parsedURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
