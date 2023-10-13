package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handlers) GetHandler(res http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	parsedURL, err := h.storage.GetURL(id)
	if err != nil {
		h.Log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", parsedURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) PingHandler(res http.ResponseWriter, r *http.Request) {
	if h.storage.Ping() {
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusInternalServerError)
}
