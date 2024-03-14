// Package logger пакет-обёртка над логгером zap
package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// Write для мидлварь
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// WriteHeader логирует статус код
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

type Log struct {
	*zap.SugaredLogger
}

// NewLogger делает новый логгер
func NewLogger() (*Log, error) {
	l := &Log{}
	cfgLogger := zap.NewDevelopmentConfig()
	opts := []zap.Option{
		zap.AddCallerSkip(1), // traverse call depth for more useful log lines
		zap.AddCaller(),
	}
	cfgLogger.DisableStacktrace = true
	lx, err := cfgLogger.Build(opts...)
	if err != nil {
		return nil, fmt.Errorf("build logger error: %w", err)
	}
	l.SugaredLogger = lx.Sugar()
	return l, nil
}

// Shutdown зачищает логгер
func (l *Log) Shutdown() error {
	return l.Sync()
}

// Debug debug
func (l *Log) Debug(args ...interface{}) {
	l.Debugln(args)
}

// Info info log
func (l *Log) Info(args ...interface{}) {
	l.Infoln(args)
}

// Warn warning log
func (l *Log) Warn(args ...interface{}) {
	l.Warnln(args)
}

// Error error log
func (l *Log) Error(args ...interface{}) {
	l.Errorln(args)
}

// Fatal debug
func (l *Log) Fatal(args ...interface{}) {
	l.Fatalln(args)
}

// Middleware для логирования хттп запросов
func (l *Log) Middleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("Content-Length")
		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		l.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
			//"userID", r.Context().Value(util.UserID).(string),
		)
	}

	return http.HandlerFunc(fn)
}
