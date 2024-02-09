package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/chemax/url-shorter/models"
	"github.com/go-chi/chi/v5"
)

func (h *handlers) getUserURLsHandler(res http.ResponseWriter, r *http.Request) {
	URLs, err := h.storage.GetUserURLs(r.Context().Value(models.UserID).(string))
	if err != nil {
		h.Log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(URLs) < 1 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	var updatedURLs []models.URLWithShort
	for _, v := range URLs {
		v.Shortcode = fmt.Sprintf("%s/%s", h.Cfg.GetHTTPAddr(), v.Shortcode)
		updatedURLs = append(updatedURLs, v)
	}
	data, err := json.Marshal(updatedURLs)
	if err != nil {
		h.Log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(data)
	if err != nil {
		h.Log.Warn("response write error: ", err.Error())
		err = nil
	}
}

func (h *handlers) getHandler(res http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedURL, err := h.storage.GetURL(id)
	if err != nil {
		h.Log.Error(fmt.Errorf("getURL error %w", err))
		if errors.Is(err, models.ErrMissingContent) {
			res.WriteHeader(http.StatusGone)
			return
		}
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", parsedURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *handlers) pingHandler(res http.ResponseWriter, r *http.Request) {
	if h.storage.Ping() {
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusInternalServerError)
}
