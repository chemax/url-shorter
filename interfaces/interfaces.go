package interfaces

import (
	"net/http"
	"time"

	"github.com/chemax/url-shorter/util"
)

type DBInterface interface {
	BatchDelete([]string, string)
	Ping() error
	SaveURL(code string, URL string, userID string) (string, error)
	Get(code string) (string, error)
	GetAllURLs(userID string) ([]util.URLStructUser, error)
	Use() bool
}
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
type StorageInterface interface {
	GetUserURLs(userID string) ([]util.URLStructUser, error)
	GetURL(code string) (parsedURL string, err error)
	DeleteListFor(forDelete []string, userID string)
	AddNewURL(parsedURL string, userID string) (code string, err error)
	Ping() bool
	BatchSave(arr []*util.URLStructForBatch, httpPrefix string) (responseArr []util.URLStructForBatchResponse, err error)
}
type ConfigInterface interface {
	SecretKey() string
	TokenExp() time.Duration
	GetNetAddr() string
	GetHTTPAddr() string
	GetSavePath() string
	GetDBUse() bool
}

type UsersInterface interface {
	Middleware(next http.Handler) http.Handler
}
