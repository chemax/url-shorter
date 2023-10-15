package util

import (
	"math/rand"
	"net/http"
	"strings"
)

type AlreadyHaveThisUrlError struct {
}

func (au *AlreadyHaveThisUrlError) Error() string {
	return "already have this url in db"
}

type DBInterface interface {
	Ping() error
	SaveURL(code string, URL string) (string, error)
	Get(code string) (string, error)
	Use() bool
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
	BatchSave(arr []*URLStructForBatch, httpPrefix string) (responseArr []URLStructForBatchResponse, err error)
}
type URLStructForBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type URLStructForBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
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

func CheckHeaderIsValidType(header string) bool {
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
