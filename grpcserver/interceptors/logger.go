package interceptors

import (
	"context"
	"time"

	"github.com/chemax/url-shorter/logger"
	"google.golang.org/grpc"
)

// LoggerInterceptor регистрирует запросы и ответы сервера grpc
func LoggerInterceptor(log *logger.Log) grpc.UnaryServerInterceptor {
	log.Debug("logger interceptor enabled")

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		args := make([]any, 0, 3)
		timeStart := time.Now()

		resp, err := handler(ctx, req)

		args = append(args, []any{
			"response_duration", time.Since(timeStart).String(),
		}...)

		log.Info("request is done", args)
		return resp, err
	}
}
