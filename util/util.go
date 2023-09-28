package util

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type LoggerInterface interface {
	Middleware(next http.Handler) http.Handler
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
}

type StorageInterface interface {
	GetURL(code string) (parsedURL *url.URL, err error)
	AddNewURL(parsedURL *url.URL) (code string, err error)
}

type ConfigInterface interface {
	GetNetAddr() string
	GetHTTPAddr() string
}

const (
	ServerAddressEnv     = "SERVER_ADDRESS"
	BaseURLEnv           = "BASE_URL"
	CodeLength           = 8
	CodeGenerateAttempts = 20
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func CheckHeaderJSONType(header string) bool {
	return strings.Contains(header, "application/json")
}
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
