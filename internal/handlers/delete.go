package handlers

import (
	"encoding/json"
	"github.com/chemax/url-shorter/util"
	"net/http"
)

func (h *Handlers) DeleteUserURLsHandler(res http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, err := res.Write(nil)
		if err != nil {
			h.Log.Error(err)
		}
		return
	}
	var forDelete []string
	err = json.Unmarshal(body, &forDelete)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, err := res.Write(nil)
		if err != nil {
			h.Log.Error(err)
		}
		return
	}
	res.WriteHeader(http.StatusAccepted)
	_, err = res.Write(nil)
	if err != nil {
		h.Log.Error(err)
	}
	h.storage.DeleteListFor(forDelete, r.Context().Value(util.UserID).(string))
}
