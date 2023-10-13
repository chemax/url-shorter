package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/chemax/url-shorter/internal/db"
	"github.com/chemax/url-shorter/internal/logger"
	"github.com/chemax/url-shorter/internal/storage"
	mock_util "github.com/chemax/url-shorter/mocks/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TODO mockgen
type cfgMock struct {
}

func gzipString(src string) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(src))
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func (c *cfgMock) GetHTTPAddr() string {
	return "http://localhost:8080"
}
func (c *cfgMock) GetNetAddr() string {
	return "localhost:8080"
}

func Test_urlManger_ServeHTTP(t *testing.T) {
	lg, _ := logger.Init()
	var tmpCode string
	const urlURL = "http://q7mtomi69.yandex/ahqas693eln9/sl3q8kiiwh4/mdcwekmdbq"
	type fields struct {
		urls map[string]*url.URL
	}
	urls := map[string]*url.URL{
		"xxxxxxxx": &url.URL{Path: "http://yandex.ru"},
		"yyyyyyyy": &url.URL{Path: "http://ya.ru"},
		"zzzzzzzz": &url.URL{Path: "http://google.com"},
		"vvvvvvvv": &url.URL{Path: "http://habr.com"},
		"qqqqqqqq": &url.URL{Path: "https://pikabu.ru"},
		"wwwwwwww": &url.URL{Path: "http://ixbt.games"},
		"rrrrrrrr": &url.URL{Path: "https://habr.com"},
	}
	type want struct {
		httpCode int
		Location string
	}
	type args struct {
		target string
		body   io.Reader
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{name: "1",
			fields: fields{urls: urls},
			args: args{
				target: "/rrrrrrrr",
				body:   nil,
				method: http.MethodPut,
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{name: "2",
			fields: fields{urls: urls},
			args: args{
				target: "/",
				body:   nil,
				method: http.MethodGet,
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		}, {name: "3",
			fields: fields{urls: urls},
			args: args{
				target: "/12345",
				body:   nil,
				method: http.MethodGet,
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{name: "4",
			fields: fields{urls: urls},
			args: args{
				target: "/12345678",
				body:   nil,
				method: http.MethodGet,
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{name: "5",
			fields: fields{urls: urls},
			args: args{
				target: "/",
				body:   nil,
				method: http.MethodPost,
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{name: "6",
			fields: fields{urls: urls},
			args: args{
				target: "/",
				body:   bytes.NewBuffer([]byte(urlURL)),
				method: http.MethodPost,
			},
			want: want{
				httpCode: http.StatusCreated,
			},
		},
		{name: "7",
			fields: fields{urls: urls},
			args: args{
				target: "replaceme",
				body:   nil,
				method: http.MethodGet,
			},
			want: want{
				httpCode: http.StatusTemporaryRedirect,
				Location: urlURL,
			},
		},
		{name: "8",
			fields: fields{urls: urls},
			args: args{
				target: "/",
				body:   bytes.NewBuffer([]byte("urlURL")),
				method: http.MethodPost,
			},
			want: want{
				httpCode: http.StatusBadRequest,
				Location: "",
			},
		},
	}
	for _, tt := range tests {
		bd, _ := db.Init(context.Background(), "")
		t.Run(tt.name, func(t *testing.T) {
			u, _ := storage.Init("", lg, bd)

			h := New(u, &cfgMock{}, lg)
			if tt.args.target == "replaceme" {
				tt.args.target = strings.Replace(tmpCode, "http://localhost:8080", "", 1)
			}
			request := httptest.NewRequest(tt.args.method, tt.args.target, tt.args.body)
			if tt.args.body != nil {
				request.Header.Set("Content-Type", "text/plain")
			}
			w := httptest.NewRecorder()
			h.Router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.httpCode, res.StatusCode)
			if tt.want.httpCode == http.StatusOK && tt.args.method == http.MethodGet {
				assert.Equal(t, tt.want.Location, res.Header.Get("Location"))
			}
			if tt.args.body != nil {
				body, err := io.ReadAll(res.Body)
				require.Equal(t, tt.want.httpCode, res.StatusCode)
				require.Nil(t, err)
				require.NotNil(t, body)
				tmpCode = string(body)

			}
		})
	}
}

func Test_urlManger_ApiServeCreate(t *testing.T) {
	lg, _ := logger.Init()
	path := "/api/shorten"
	const urlURL = "http://q7mtomi69.yandex/ahqas693eln9/sl3q8kiiwh4/mdcwekmdbq"
	type fields struct {
		urls map[string]*url.URL
	}
	urls := map[string]*url.URL{
		"xxxxxxxx": &url.URL{Path: "http://yandex.ru"},
		"yyyyyyyy": &url.URL{Path: "http://ya.ru"},
		"zzzzzzzz": &url.URL{Path: "http://google.com"},
		"vvvvvvvv": &url.URL{Path: "http://habr.com"},
		"qqqqqqqq": &url.URL{Path: "https://pikabu.ru"},
		"wwwwwwww": &url.URL{Path: "http://ixbt.games"},
		"rrrrrrrr": &url.URL{Path: "https://habr.com"},
	}
	type want struct {
		httpCode int
		Location string
	}
	type args struct {
		body   io.Reader
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{name: "5",
			fields: fields{urls: urls},
			args: args{
				body:   nil,
				method: http.MethodPost,
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{name: "6",
			fields: fields{urls: urls},
			args: args{
				body:   bytes.NewBuffer([]byte(fmt.Sprintf("{\"url\": \"%s\"}", urlURL))),
				method: http.MethodPost,
			},
			want: want{
				httpCode: http.StatusCreated,
			},
		},
		{name: "7",
			fields: fields{urls: urls},
			args: args{
				body:   bytes.NewBuffer([]byte("")),
				method: http.MethodPost,
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{name: "8",
			fields: fields{urls: urls},
			args: args{
				body:   bytes.NewBuffer([]byte("{\"url\": \"url\"}")),
				method: http.MethodPost,
			},
			want: want{
				httpCode: http.StatusBadRequest,
				Location: "",
			},
		},
	}
	for _, tt := range tests {
		bd, _ := db.Init(context.Background(), "")
		t.Run(tt.name, func(t *testing.T) {
			u, _ := storage.Init("", lg, bd)

			h := New(u, &cfgMock{}, lg)
			request := httptest.NewRequest(tt.args.method, path, tt.args.body)

			if tt.args.body != nil {
				request.Header.Set("content-type", "application/json")
			}
			w := httptest.NewRecorder()
			h.Router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.httpCode, res.StatusCode)
			if tt.args.body != nil {
				body, err := io.ReadAll(res.Body)
				require.Equal(t, tt.want.httpCode, res.StatusCode)
				require.Nil(t, err)
				require.NotNil(t, body)
			}
		})
	}
}

func TestHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := mock_util.NewMockConfigInterface(ctrl)
	cfg.EXPECT().GetHTTPAddr().Return("http://127.0.0.1:8080").AnyTimes()
	st := mock_util.NewMockStorageInterface(ctrl)
	st.EXPECT().AddNewURL(gomock.Any()).Return("12345678", nil).Times(2)
	st.EXPECT().GetURL(gomock.Any()).AnyTimes()
	//st.EXPECT().Ping().AnyTimes()
	log, _ := logger.Init()
	handlers := New(st, cfg, log)
	assert.NotNil(t, handlers)
	JSONURL := "{\"url\": \"http://ya.ru\"}"
	JSONBadURL := "{\"url\": \".ru\"}"
	t.Run("all ok gzip", func(t *testing.T) {
		gzString, _ := gzipString(JSONURL)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(gzString))
		request.Header.Set("Content-Type", "application/x-gzip")
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
	t.Run("invalid header gzip", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONURL)))
		request.Header.Set("Content-Type", "application/x-gzip")
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("all ok", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONURL)))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
	t.Run("ping bad", func(t *testing.T) {
		st.EXPECT().Ping().Times(1).Return(false)
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("ping good", func(t *testing.T) {
		st.EXPECT().Ping().Times(1).Return(true)
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("bad json", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte("JSONURL")))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("bad URL", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONBadURL)))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("bad content type", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONBadURL)))
		request.Header.Set("Content-Type", "alication/js")
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	st.EXPECT().AddNewURL(gomock.Any()).Times(1).Return("", fmt.Errorf("test error"))
	t.Run("store URL error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(JSONURL)))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handlers.Router.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
