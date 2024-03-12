// Package gRPCHandlers обработчики и перехватчики для gRPC сервера
package gRPCHandlers

import (
	"context"
	"github.com/chemax/url-shorter/models"
	pb "github.com/chemax/url-shorter/proto"
	"google.golang.org/grpc"
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

func (h *URLShortenerServer) Ping(ctx context.Context, in *pb.PingRequest, opts ...grpc.CallOption) (*pb.PingResponse, error) {
	return nil, nil
}
func (h *URLShortenerServer) GetOriginalURL(ctx context.Context, in *pb.UnshortURLRequest, opts ...grpc.CallOption) (*pb.UnshortURLResponse, error) {
	return nil, nil
}
func (h *URLShortenerServer) GetURLsByUserID(ctx context.Context, in *pb.GetUserURLsRequest, opts ...grpc.CallOption) (*pb.GetUserURLsResponse, error) {
	return nil, nil
}
func (h *URLShortenerServer) CreateURL(ctx context.Context, in *pb.ShortURLRequest, opts ...grpc.CallOption) (*pb.ShortURLResponse, error) {
	return nil, nil
}
func (h *URLShortenerServer) CreateURLs(ctx context.Context, in *pb.ShortURLsBatchRequest, opts ...grpc.CallOption) (*pb.ShortURLsBatchResponse, error) {
	return nil, nil
}
func (h *URLShortenerServer) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest, opts ...grpc.CallOption) (*pb.DeleteURLsResponse, error) {
	return nil, nil
}
func (h *URLShortenerServer) Stat(ctx context.Context, in *pb.StatRequest, opts ...grpc.CallOption) (*pb.StatResponse, error) {
	return nil, nil
}
