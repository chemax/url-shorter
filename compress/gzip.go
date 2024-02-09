package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	GzWriter io.Writer
	compress bool
}

// WriteHeader записывает нужный хидер
func (w *gzipWriter) WriteHeader(statusCode int) {
	if strings.Contains(w.Header().Get("Content-Type"), "application/json") ||
		strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		w.compress = true
		w.Header().Set("Content-Encoding", "gzip")
	}
	if w.Header().Get("Content-Encoding") != "" {
		w.Header().Del("Content-Length")
	}
	w.ResponseWriter.WriteHeader(statusCode)

}

// Write обертка над gzip.Writer чтобы определять нужно ли сжимать
func (w *gzipWriter) Write(b []byte) (int, error) {
	if w.compress {
		gz, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w.ResponseWriter, err.Error())
		}
		defer gz.Close()
		return gz.Write(b)
	} else {
		return w.ResponseWriter.Write(b)
	}
}

// Middleware для работы со сжатым контентом
func Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(&gzipWriter{ResponseWriter: w}, r)
	})

}
