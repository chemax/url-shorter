package compress

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	middleware := Middleware(handler)

	t.Run("no gzip compression", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello, World!", rec.Body.String())
	})

}

func TestGzipWriter_WriteHeader(t *testing.T) {
	w := gzipWriter{
		ResponseWriter: httptest.NewRecorder(),
	}

	t.Run("no content type", func(t *testing.T) {
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, http.StatusOK, w.ResponseWriter.(*httptest.ResponseRecorder).Code)
	})

	t.Run("application/json content type", func(t *testing.T) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, http.StatusOK, w.ResponseWriter.(*httptest.ResponseRecorder).Code)
		assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	})

	t.Run("text/html content type", func(t *testing.T) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, http.StatusOK, w.ResponseWriter.(*httptest.ResponseRecorder).Code)
		assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	})

	t.Run("content encoding already set", func(t *testing.T) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)

		assert.Equal(t, http.StatusOK, w.ResponseWriter.(*httptest.ResponseRecorder).Code)
		assert.Equal(t, "", w.Header().Get("Content-Length"))
	})
}
