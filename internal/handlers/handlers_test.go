package handlers

import (
	"bytes"
	"github.com/chemax/url-shorter/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_urlManger_ServeHTTP(t *testing.T) {
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

		t.Run(tt.name, func(t *testing.T) {
			u := storage.Get()
			h := New(u)
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
