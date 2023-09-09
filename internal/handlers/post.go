package handlers

import (
	"fmt"
	util "github.com/chemax/url-shorter/util"
	"io"
	"net/http"
	"net/url"
)

func (h *Handlers) ServeCreate(res http.ResponseWriter, req *http.Request) {
	if !util.CheckHeader(req.Header.Get("Content-Type")) {
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
	code, err := h.storage.AddNewURL(parsedURL)
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
	_, err = res.Write([]byte(fmt.Sprintf("%s/%s", h.Cfg.GetHTTPAddr(), code)))
	if err != nil {
		fmt.Println(err.Error())
	}
}
