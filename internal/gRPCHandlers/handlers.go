// Package gRPCHandlers обработчики и перехватчики для gRPC сервера
package gRPCHandlers

import (
	"context"
	"fmt"
	"github.com/chemax/url-shorter/models"
	pb "github.com/chemax/url-shorter/proto"
	"net"
	"net/http"
	"net/url"
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
func New(s Storager, cfg Configer, log Loggerer, users Userser) *URLShortenerServer {
	return &URLShortenerServer{
		log:     log,
		storage: s,
		Cfg:     cfg,
	}
}

// Ping связь с базой
func (h *URLShortenerServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse
	response.Message = h.storage.Ping()
	return &response, nil
}

// GetOriginalURL получить адрес по его короткому коду
func (h *URLShortenerServer) GetOriginalURL(ctx context.Context, in *pb.UnshortURLRequest) (*pb.UnshortURLResponse, error) {
	var response pb.UnshortURLResponse
	var err error
	response.OriginalUrl, err = h.storage.GetURL(in.ShortUrl)
	return &response, err
}

// GetURLsByUserID получить все адреса пользователя
func (h *URLShortenerServer) GetURLsByUserID(ctx context.Context, in *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	var response pb.GetUserURLsResponse
	userID := ctx.Value(models.UserID).(string)
	userURLs, err := h.storage.GetUserURLs(userID)
	if err != nil {
		return &pb.GetUserURLsResponse{}, err
	}
	for _, u := range userURLs {
		uGRPC := &pb.URLEntity{
			ShortUrl:    u.Shortcode,
			OriginalUrl: u.URL,
		}
		response.Result = append(response.Result, uGRPC)
	}
	return &response, nil
}

// CreateURL зашортить URL
func (h *URLShortenerServer) CreateURL(ctx context.Context, in *pb.ShortURLRequest) (*pb.ShortURLResponse, error) {
	var response pb.ShortURLResponse
	requestURI, err := url.ParseRequestURI(in.Url)
	if err != nil {
		return nil, err
	}
	newURL, err := h.storage.AddNewURL(requestURI.RequestURI(), ctx.Value(models.UserID).(string))
	if err != nil {
		return nil, err
	}
	response.Result = newURL
	return &response, nil
}

// CreateURLs зашортить пачку URL
func (h *URLShortenerServer) CreateURLs(ctx context.Context, in *pb.ShortURLsBatchRequest) (*pb.ShortURLsBatchResponse, error) {
	var response pb.ShortURLsBatchResponse
	var URLsArr []*models.URLForBatch
	for _, v := range in.Dto {
		URLsArr = append(URLsArr, &models.URLForBatch{
			CorrelationID: v.CorrelationId,
			OriginalURL:   v.OriginalUrl,
		})
	}
	saveResponse, err := h.storage.BatchSave(URLsArr, h.Cfg.GetHTTPAddr())
	if err != nil {
		return nil, err
	}
	for _, v := range saveResponse {
		var d pb.ShortURLsBatchResponse_Data
		d.CorrelationId = v.CorrelationID
		d.ShortUrl = v.ShortURL
		response.Result = append(response.Result, &d)
	}
	return &response, nil
}

// DeleteURLs удалить пачку урлов.
func (h *URLShortenerServer) DeleteURLs(ctx context.Context, in *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	h.storage.DeleteListFor(in.Urls, ctx.Value(models.UserID).(string))
	return &pb.DeleteURLsResponse{}, nil
}

// Stat статистика сервиса под защитой могучей проверки адреса
func (h *URLShortenerServer) Stat(ctx context.Context, in *pb.StatRequest) (*pb.StatResponse, error) {
	if h.checkIP(ctx.Value(models.RealIP).(string)) {
		stat, err := h.storage.GetStats()
		if err != nil {
			return &pb.StatResponse{}, fmt.Errorf("400")
		}
		return &pb.StatResponse{Urls: stat.URLs, Users: stat.Users}, nil
	}
	return &pb.StatResponse{}, fmt.Errorf("403")
}

// checkIP дублирование кода, кек, но все говорят что папка Utils это плохо. А куда мне ещё эту фигню выносить, отдельный пакет на две строки делать?
func (h *URLShortenerServer) checkIP(ipAsString string) bool {
	realIP := net.ParseIP(ipAsString)
	return realIP != nil && h.TrustedSubnet.Contains(realIP)
}
