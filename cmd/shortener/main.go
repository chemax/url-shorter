package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const (
	CodeLength           = 8
	CodeGenerateAttempts = 20
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type urlManger struct {
	urls map[string]*url.URL
}

func (u *urlManger) addNewURL(parsedURL *url.URL) (code string, err error) {
	var ok = true
	var loop int
	for ok {
		code = RandStringRunes(CodeLength)
		_, ok = u.urls[code]
		loop++
		if loop > CodeGenerateAttempts {
			code = ""
			return code, fmt.Errorf("can not found free code for short url")
		}
	}
	u.urls[code] = parsedURL
	return
}

func (u *urlManger) ServeCreate(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "text/plain" {
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
	code, err := u.addNewURL(parsedURL)
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
	_, err = res.Write([]byte(code))
	if err != nil {
		fmt.Println(err.Error())
	}
}
func (u *urlManger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		u.ServeCreate(res, req)
	} else if req.Method == http.MethodGet {
		if req.URL.Path == "/" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		shortCode := strings.TrimPrefix(req.URL.Path, "/")
		if len(shortCode) != CodeLength {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		parsedURL, ok := u.urls[shortCode]
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		http.Redirect(res, req, parsedURL.String(), http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	u := &urlManger{urls: make(map[string]*url.URL)}

	err := http.ListenAndServe("localhost:8080", u)
	if err != nil {
		panic(err)
	}

}
