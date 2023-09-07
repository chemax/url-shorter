package storage

import (
	"fmt"
	. "github.com/chemax/url-shorter/util"
	"net/url"
)

type UrlManger struct {
	urls map[string]*url.URL
}

var manager = &UrlManger{urls: make(map[string]*url.URL)}

func Get() *UrlManger {
	return manager
}

func (u *UrlManger) GetUrl(code string) (parsedURL *url.URL, err error) {
	parsedURL, ok := u.urls[code]
	if !ok {
		return nil, fmt.Errorf("404")
	}
	return parsedURL, nil
}

func (u *UrlManger) AddNewURL(parsedURL *url.URL) (code string, err error) {
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

//func (u *UrlManger) ServeCreate(res http.ResponseWriter, req *http.Request) {
//	if !CheckHeader(req.Header.Get("Content-Type")) {
//		fmt.Println("not plain text", req.Header.Get("Content-Type"))
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	body, err := io.ReadAll(req.Body)
//	if err != nil {
//		fmt.Println(err.Error())
//		res.WriteHeader(http.StatusBadRequest)
//	}
//	fmt.Println(string(body))
//	parsedURL, err := url.ParseRequestURI(string(body))
//	if err != nil {
//		fmt.Println(err.Error())
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	code, err := u.AddNewURL(parsedURL)
//	if err != nil {
//		fmt.Println(err.Error())
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	if code == "" {
//		fmt.Println("code is empty", code)
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	res.Header().Set("content-type", "text/plain")
//	res.WriteHeader(http.StatusCreated)
//	_, err = res.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", code)))
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
//func (u *UrlManger) ServeGET(res http.ResponseWriter, req *http.Request) {
//	fmt.Println("URL requested: ", req.URL)
//	if req.URL.Path == "/" {
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	shortCode := strings.TrimPrefix(req.URL.Path, "/")
//	if len(shortCode) != CodeLength {
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	parsedURL, err := u.GetUrl(shortCode)
//	if err != nil {
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	res.Header().Set("Location", parsedURL.String())
//	res.WriteHeader(http.StatusTemporaryRedirect)
//}
//
//func (u *UrlManger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
//	if req.Method == http.MethodPost {
//		u.ServeCreate(res, req)
//	} else if req.Method == http.MethodGet {
//		u.ServeGET(res, req)
//	} else {
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//}
