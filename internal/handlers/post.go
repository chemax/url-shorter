package handlers

import (
	"compress/gzip"
	"encoding/json"
	"errors"
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
			return body, fmt.Errorf("gzip newReader error: %w", err)
		}
		body, err = io.ReadAll(reader)
	} else {
		body, err = io.ReadAll(req.Body)
	}
	if err != nil {
		return nil, fmt.Errorf("get body error: %w", err)
	}
	return body, nil
}

func (h *Handlers) PostHandler(res http.ResponseWriter, req *http.Request) {
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
		err = fmt.Errorf("get body error: %w", err)
		return
	}
	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		err = fmt.Errorf("parse URL error: %w", err)
		return
	}
	code, err := h.store(parsedURL, req.Context().Value(util.UserID).(string))
	var statusCreated = http.StatusCreated
	if err != nil {
		if !errors.Is(err, &util.AlreadyHaveThisURLError{}) {
			err = fmt.Errorf("store error: %w", err)
			return
		}
		statusCreated = http.StatusConflict
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(statusCreated)
	_, err = res.Write([]byte(fmt.Sprintf("%s/%s", h.Cfg.GetHTTPAddr(), code)))
	if err != nil {
		h.Log.Warn("response write error: ", err.Error())
		err = nil
	}
}

func (h *Handlers) JSONBatchPostHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	defer func() {
		if err != nil {
			h.Log.Warn(err.Error())
			res.WriteHeader(http.StatusBadRequest)
		}
	}()

	if !util.CheckHeaderIsValidType(req.Header.Get("Content-Type")) {
		err = fmt.Errorf("not application/json: %s", req.Header.Get("Content-Type"))
		return
	}

	var URLBatchArr []*util.URLStructForBatch
	body, err := getBody(req)
	if err != nil {
		err = fmt.Errorf("get body error: %w", err)
		return
	}
	err = json.Unmarshal(body, &URLBatchArr)
	if err != nil {
		err = fmt.Errorf("JSON unmarshal error: %w", err)
		return
	}
	save, err := h.storage.BatchSave(URLBatchArr, h.Cfg.GetHTTPAddr())
	if err != nil {
		return
	}
	resultData, err := json.Marshal(save)
	if err != nil {
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write(resultData)
	if err != nil {
		h.Log.Warn("response write error: ", err.Error())
		err = nil
	}
	//TODO как-то тут многовато общего кода с JSONPostHandler, DRY or not DRY?
}

func (h *Handlers) JSONPostHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	defer func() {
		if err != nil {
			h.Log.Warn(err.Error())
			res.WriteHeader(http.StatusBadRequest)
		}
	}()
	type URLStruct struct {
		URL    string `json:"url"`
		UserID string `json:"userID"`
	}
	type ResultStruct struct {
		Result string `json:"result"`
	}
	if !util.CheckHeaderIsValidType(req.Header.Get("Content-Type")) {
		err = fmt.Errorf("not application/json: %s", req.Header.Get("Content-Type"))
		return
	}
	body, err := getBody(req)
	if err != nil {
		err = fmt.Errorf("get body error: %w", err)
		return
	}
	userID := req.Context().Value(util.UserID).(string)
	URLObj := URLStruct{UserID: userID}
	err = json.Unmarshal(body, &URLObj)
	if err != nil {
		err = fmt.Errorf("JSON unmarshal error: %w", err)
		return
	}
	parsedURL, err := url.ParseRequestURI(URLObj.URL)
	if err != nil {
		err = fmt.Errorf("parse URL error: %w", err)
		return
	}
	code, err := h.store(parsedURL, req.Context().Value(util.UserID).(string))
	var statusCreated = http.StatusCreated
	if err != nil {
		if errors.Is(err, &util.AlreadyHaveThisURLError{}) {
			statusCreated = http.StatusConflict
		} else {
			err = fmt.Errorf("store error: %w", err)
			return
		}
	}

	result := ResultStruct{Result: fmt.Sprintf("%s/%s", h.Cfg.GetHTTPAddr(), code)}
	resultData, err := json.Marshal(result)
	if err != nil {
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(statusCreated)
	_, err = res.Write(resultData)
	if err != nil {
		h.Log.Warn("response write error: ", err.Error())
		err = nil
	}
}

func (h *Handlers) store(parsedURL *url.URL, userID string) (string, error) {
	code, err := h.storage.AddNewURL(parsedURL.String(), userID)
	if err != nil {
		return code, err
	}
	if code == "" {
		return "", fmt.Errorf("cannot generate short url")
	}
	return code, nil
}
