package storage

import (
	"fmt"
	"github.com/chemax/url-shorter/util"
	"net/url"
)

type URLManager struct {
	urls map[string]*url.URL
}

var manager = &URLManager{urls: make(map[string]*url.URL)}

func Get() *URLManager {
	return manager
}

func (u *URLManager) GetURL(code string) (parsedURL *url.URL, err error) {
	parsedURL, ok := u.urls[code]
	if !ok {
		return nil, fmt.Errorf("requested url not found")
	}
	return parsedURL, nil
}

func (u *URLManager) AddNewURL(parsedURL *url.URL) (code string, err error) {
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
	return code, err
}
