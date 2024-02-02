package interfaces

import (
	"net/http"
	"time"

	"github.com/chemax/url-shorter/util"
)

// DBInterface интерфейс для базы данных
type DBInterface interface {
	BatchDelete([]string, string)
	Ping() error
	SaveURL(code string, URL string, userID string) (string, error)
	Get(code string) (string, error)
	GetAllURLs(userID string) ([]util.URLWithShort, error)
	Use() bool
}

// LoggerInterface интерфейс логера
type LoggerInterface interface {
	Middleware(next http.Handler) http.Handler
	Debug(args ...interface{})
	Debugln(args ...interface{})
	Info(args ...interface{})
	Infoln(args ...interface{})
	Warn(args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
}

// StorageInterface интерфейс хранилища
type StorageInterface interface {
	GetUserURLs(userID string) ([]util.URLWithShort, error)
	GetURL(code string) (parsedURL string, err error)
	DeleteListFor(forDelete []string, userID string)
	AddNewURL(parsedURL string, userID string) (code string, err error)
	Ping() bool
	BatchSave(arr []*util.URLForBatch, httpPrefix string) (responseArr []util.URLForBatchResponse, err error)
}

// ConfigInterface интерфейс конфиг-структуры
type ConfigInterface interface {
	SecretKey() string
	TokenExp() time.Duration
	GetNetAddr() string
	GetHTTPAddr() string
	GetSavePath() string
	GetDBUse() bool
}

// UsersInterface интерфейс юзер-менеджера
type UsersInterface interface {
	Middleware(next http.Handler) http.Handler
}
