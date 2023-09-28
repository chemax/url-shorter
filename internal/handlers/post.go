package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/chemax/url-shorter/util"
	"io"
	"net/http"
	"net/url"
)

func (h *Handlers) ServeCreate(res http.ResponseWriter, req *http.Request) {
	var err error
	defer func() {
		if err != nil {
			h.Log.Warn(err.Error())
			res.WriteHeader(http.StatusBadRequest)
		}
	}()
	if !util.CheckHeader(req.Header.Get("Content-Type")) {
		err = fmt.Errorf("not plain text: %s", req.Header.Get("Content-Type"))
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		return
	}
	code, err := h.store(parsedURL)
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(fmt.Sprintf("%s/%s", h.Cfg.GetHTTPAddr(), code)))
	if err != nil {
		fmt.Println(err.Error())
		err = nil
	}
}

func (h *Handlers) APIServeCreate(res http.ResponseWriter, req *http.Request) {
	var err error
	defer func() {
		if err != nil {
			h.Log.Warn(err.Error())
			res.WriteHeader(http.StatusBadRequest)
		}
	}()
	type URLStruct struct {
		URL string `json:"url"`
	}
	type ResultStruct struct {
		Result string `json:"result"`
	}
	if !util.CheckHeaderJsonType(req.Header.Get("Content-Type")) {
		err = fmt.Errorf("not application/json: %s", req.Header.Get("Content-Type"))
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	URLObj := URLStruct{}
	err = json.Unmarshal(body, &URLObj)
	if err != nil {
		return
	}
	parsedURL, err := url.ParseRequestURI(URLObj.URL)
	if err != nil {
		return
	}
	code, err := h.store(parsedURL)
	if err != nil {
		return
	}

	result := ResultStruct{Result: fmt.Sprintf("%s/%s", h.Cfg.GetHTTPAddr(), code)}
	resultData, err := json.Marshal(result)
	if err != nil {
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(resultData)
	if err != nil {
		h.Log.Warn(err.Error())
		err = nil
	}
}

func (h *Handlers) store(parsedURL *url.URL) (string, error) {
	code, err := h.storage.AddNewURL(parsedURL)
	if err != nil {
		return "", err
	}
	if code == "" {
		return "", fmt.Errorf("cannot generate short url")
	}
	return code, nil
}
