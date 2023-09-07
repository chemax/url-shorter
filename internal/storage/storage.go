package storage

import (
	"fmt"
	"github.com/chemax/url-shorter/util"
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
		code = util.RandStringRunes(util.CodeLength)
		_, ok = u.urls[code]
		loop++
		if loop > util.CodeGenerateAttempts {
			code = ""
			return code, fmt.Errorf("can not found free code for short url")
		}
	}
	u.urls[code] = parsedURL
	return
}
