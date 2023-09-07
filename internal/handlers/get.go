package handlers

import (
	"fmt"
	"github.com/chemax/url-shorter/internal/storage"
	. "github.com/chemax/url-shorter/util"
	"net/http"
	"strings"
)

func (h *Handlers) serveGET(res http.ResponseWriter, req *http.Request, u *storage.UrlManger) {
	fmt.Println("URL requested: ", req.URL)
	if req.URL.Path == "/" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortCode := strings.TrimPrefix(req.URL.Path, "/")
	if len(shortCode) != CodeLength {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	parsedURL, err := u.GetUrl(shortCode)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", parsedURL.String())
	res.WriteHeader(http.StatusTemporaryRedirect)
}
