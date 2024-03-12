package grpcserver

import (
	"context"
	"fmt"
	"github.com/chemax/url-shorter/config"
	"github.com/chemax/url-shorter/internal/gRPCHandlers"
	pb "github.com/chemax/url-shorter/proto"
	"google.golang.org/grpc"
	"net"
	"os"
)

// Loggerer интерфейс логера
type Loggerer interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

func New(ctx context.Context, cfg *config.Config, log Loggerer, sig chan os.Signal, h gRPCHandlers.URLShortenerServer) error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterURLShortenerV1Server(s, h)

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	err = s.Serve(listen)
	return err
}
