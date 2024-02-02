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

type log struct {
	*zap.SugaredLogger
}

// Init делает новый логгер
func Init() (*log, error) {
	l := &log{}
	cfgLogger := zap.NewDevelopmentConfig()
	cfgLogger.DisableStacktrace = true
	lx, err := cfgLogger.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger error: %w", err)
	}
	l.SugaredLogger = lx.Sugar()
	return l, nil
}

// Shutdown зачищает логгер
func (l *log) Shutdown() error {
	return l.Sync()
}

// Debug debug
func (l *log) Debug(args ...interface{}) {
	l.Debugln(args)
}

// Info info log
func (l *log) Info(args ...interface{}) {
	l.Infoln(args)
}

// Warn warning log
func (l *log) Warn(args ...interface{}) {
	l.Warnln(args)
}

// Error error log
func (l *log) Error(args ...interface{}) {
	l.Errorln(args)
}

// Middleware для логирования хттп запросов
func (l *log) Middleware(next http.Handler) http.Handler {

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
