package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	GzWriter io.Writer
	compress bool
}

func (w *gzipWriter) WriteHeader(code int) {
	if strings.Contains(w.Header().Get("Content-Type"), "application/json") ||
		strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		w.compress = true
		w.Header().Set("Content-Encoding", "gzip")
	}
	if w.Header().Get("Content-Encoding") != "" {
		w.Header().Del("Content-Length")
		return
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	if w.compress {
		fmt.Println("GzWriter")
		gz, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w.ResponseWriter, err.Error())
		}
		defer gz.Close()
		return gz.Write(b)
	} else {
		fmt.Println("ResponseWriter")
		return w.ResponseWriter.Write(b)
	}
}

func Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(&gzipWriter{ResponseWriter: w}, r)
	})

}
