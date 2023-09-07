package handlers

import (
	"fmt"
	"github.com/chemax/url-shorter/internal/storage"
	. "github.com/chemax/url-shorter/util"
	"io"
	"net/http"
	"net/url"
)

func (h *Handlers) serveCreate(res http.ResponseWriter, req *http.Request, u *storage.UrlManger) {
	if !CheckHeader(req.Header.Get("Content-Type")) {
		fmt.Println("not plain text", req.Header.Get("Content-Type"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
	}
	fmt.Println(string(body))
	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	code, err := u.AddNewURL(parsedURL)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if code == "" {
		fmt.Println("code is empty", code)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", code)))
	if err != nil {
		fmt.Println(err.Error())
	}
}
