package util

import (
	"math/rand"
	"net/url"
	"strings"
)

type StorageInterface interface {
	GetURL(code string) (parsedURL *url.URL, err error)
	AddNewURL(parsedURL *url.URL) (code string, err error)
}

type ConfigInterface interface {
	GetNetAddr() string
	GetHTTPAddr() string
}

const (
	CodeLength           = 8
	CodeGenerateAttempts = 20
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func CheckHeader(header string) bool {
	return strings.Contains(header, "text/plain")
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
