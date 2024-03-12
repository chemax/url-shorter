// Package gRPCHandlers обработчики и перехватчики для gRPC сервера
package gRPCHandlers

import (
	"github.com/chemax/url-shorter/models"
	pb "github.com/chemax/url-shorter/proto"
	"net"
	"net/http"
)

// Configer интерфейс конфиг-структуры
type Configer interface {
	GetHTTPAddr() string
	GetTrustedSubnet() string
}

// Userser интерфейс юзер-менеджера
type Userser interface {
	Middleware(next http.Handler) http.Handler
}

// Loggerer интерфейс логера
type Loggerer interface {
	Middleware(next http.Handler) http.Handler
	Warn(args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
}

// Storager интерфейс хранилища
type Storager interface {
	GetUserURLs(userID string) ([]models.URLWithShort, error)
	GetURL(code string) (parsedURL string, err error)
	DeleteListFor(forDelete []string, userID string)
	AddNewURL(parsedURL string, userID string) (code string, err error)
	Ping() bool
	BatchSave(arr []*models.URLForBatch, httpPrefix string) (responseArr []models.URLForBatchResponse, err error)
	GetStats() (models.Stats, error)
}

// URLShortenerServer is the server that provides the URLShortener service.
type URLShortenerServer struct {
	pb.UnimplementedURLShortenerV1Server

	log           Loggerer
	storage       Storager
	Cfg           Configer
	TrustedSubnet *net.IPNet
}

// New creates a new URLShortenerServer
func New(s Storager, cfg Configer, log Loggerer, users Userser) URLShortenerServer {
	return URLShortenerServer{
		log:     log,
		storage: s,
		Cfg:     cfg,
	}
}
