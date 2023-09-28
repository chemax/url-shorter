package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
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

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

type Logger struct {
	sugar *zap.SugaredLogger
}

func New() *Logger {
	l := &Logger{}
	cfgLogger := zap.NewDevelopmentConfig()
	cfgLogger.DisableStacktrace = true
	lx, err := cfgLogger.Build()
	if err != nil {
		panic(err)
	}
	l.sugar = lx.Sugar()
	return l
}

func (l *Logger) Sync() error {
	return l.sugar.Sync()
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debugln(args)
}
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Infoln(args)
}
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warnln(args)
}
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Errorln(args)
}
func (l *Logger) Middleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

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
		l.sugar.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}

	return http.HandlerFunc(fn)
}
