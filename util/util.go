package util

import (
	"math/rand"
	"net/http"
	"strings"
)

type DBInterface interface {
	Ping() error
	SaveURL(code string, URL string) error
	Get(code string) (string, error)
}
type LoggerInterface interface {
	Middleware(next http.Handler) http.Handler
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type StorageInterface interface {
	GetURL(code string) (parsedURL string, err error)
	AddNewURL(parsedURL string) (code string, err error)
	Ping() bool
}

type ConfigInterface interface {
	GetNetAddr() string
	GetHTTPAddr() string
}

const (
	ServerAddressEnv     = "SERVER_ADDRESS"
	BaseURLEnv           = "BASE_URL"
	SavePath             = "FILE_STORAGE_PATH"
	DBConnectString      = "DATABASE_DSN"
	CodeLength           = 8
	CodeGenerateAttempts = 20
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func CheckHeaderJSONType(header string) bool {
	return strings.Contains(header, "application/json") || strings.Contains(header, "application/x-gzip")
}
func CheckHeader(header string) bool {
	return strings.Contains(header, "text/plain") || strings.Contains(header, "application/x-gzip")
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
