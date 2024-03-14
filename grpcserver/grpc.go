package grpcserver

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/chemax/url-shorter/config"
	"github.com/chemax/url-shorter/grpcserver/interceptors"
	"github.com/chemax/url-shorter/internal/gRPCHandlers"
	"github.com/chemax/url-shorter/logger"
	pb "github.com/chemax/url-shorter/proto"
	"google.golang.org/grpc"
)

// Loggerer интерфейс логера
type Loggerer interface {
	Warn(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
}

// Userser интерфейс юзер-менеджера
type Userser interface {
	JWTInterceptor(log *logger.Log) grpc.UnaryServerInterceptor
}

// New возвращает grpc сервер
// контекст должен использоваться для грейсфула, конфиг для порта и всё такое, сиг тоже для грейсфула.
func New(ctx context.Context, cfg *config.Config, log *logger.Log, sig chan os.Signal, h *gRPCHandlers.URLShortenerServer, users Userser) error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors.LoggerInterceptor(log)),
		grpc.ChainUnaryInterceptor(users.JWTInterceptor(log)),
	)
	// регистрируем сервис
	pb.RegisterURLShortenerV1Server(s, h)

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	err = s.Serve(listen)
	return err
}
