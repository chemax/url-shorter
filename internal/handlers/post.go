package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/chemax/url-shorter/util"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getBody(req *http.Request) ([]byte, error) {
	var body []byte
	var err error
	var reader *gzip.Reader
	if strings.Contains(req.Header.Get("Content-Type"), "gzip") {
		reader, err = gzip.NewReader(req.Body)
		if err != nil {
			return body, err
		}
		body, err = io.ReadAll(reader)
	} else {
		body, err = io.ReadAll(req.Body)
	}
	return body, err
}

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
	body, err := getBody(req)
	if err != nil {
		return
	}
	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		return
	}
	code, err := h.store(parsedURL)
	if err != nil {
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(fmt.Sprintf("%s/%s", h.Cfg.GetHTTPAddr(), code)))
	if err != nil {
		h.Log.Warn(err.Error())
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
	if !util.CheckHeaderJSONType(req.Header.Get("Content-Type")) {
		err = fmt.Errorf("not application/json: %s", req.Header.Get("Content-Type"))
		return
	}
	body, err := getBody(req)
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
	code, err := h.storage.AddNewURL(parsedURL.String())
	if err != nil {
		return "", err
	}
	if code == "" {
		return "", fmt.Errorf("cannot generate short url")
	}
	return code, nil
}
