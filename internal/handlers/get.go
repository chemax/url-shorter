package handlers

import (
	"encoding/json"
	"github.com/chemax/url-shorter/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handlers) GetUserURLsHandler(res http.ResponseWriter, r *http.Request) {
	if r.Context().Value(util.UserID).(string) == "" {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	URLs, err := h.storage.GetUserURLs(r.Context().Value(util.UserID).(string))
	if err != nil {
		h.Log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(URLs) < 1 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	data, err := json.Marshal(URLs)
	if err != nil {
		h.Log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "application/json")
	_, err = res.Write(data)
	if err != nil {
		h.Log.Warn("response write error: ", err.Error())
		err = nil
	}
}

func (h *Handlers) GetHandler(res http.ResponseWriter, r *http.Request) {
	h.Log.Debug(util.UserID, r.Context().Value(util.UserID))
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
