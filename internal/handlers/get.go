package handlers

import (
	util "github.com/chemax/url-shorter/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handlers) serveGET(res http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	//fmt.Println("URL requested: ", req.URL)
	//if req.URL.Path == "/" {
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//shortCode := strings.TrimPrefix(req.URL.Path, "/")
	if len(id) != util.CodeLength {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	parsedURL, err := h.storage.GetURL(id)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", parsedURL.String())
	res.WriteHeader(http.StatusTemporaryRedirect)
}
