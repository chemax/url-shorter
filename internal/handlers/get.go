package handlers

import (
	"fmt"
	util "github.com/chemax/url-shorter/util"
	"net/http"
	"strings"
)

func (h *Handlers) serveGET(res http.ResponseWriter, req *http.Request) {
	fmt.Println("URL requested: ", req.URL)
	if req.URL.Path == "/" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortCode := strings.TrimPrefix(req.URL.Path, "/")
	if len(shortCode) != util.CodeLength {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	parsedURL, err := h.storage.GetURL(shortCode)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", parsedURL.String())
	res.WriteHeader(http.StatusTemporaryRedirect)
}
