package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_urlManger_ServeHTTP(t *testing.T) {
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
				method: http.MethodGet,
			},
			want: want{
				httpCode: http.StatusTemporaryRedirect,
				Location: urls["rrrrrrrr"].String(),
			},
		}, {name: "2",
			fields: fields{urls: urls},
			args: args{
				target: "/",
				body:   nil,
				method: http.MethodGet,
			},
			want: want{
				httpCode: http.StatusBadRequest,
				Location: "",
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
				Location: "",
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
				Location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &urlManger{
				urls: tt.fields.urls,
			}
			request := httptest.NewRequest(tt.args.method, tt.args.target, tt.args.body)
			w := httptest.NewRecorder()
			u.ServeHTTP(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.httpCode, res.StatusCode)
			if tt.want.httpCode == http.StatusOK {
				assert.Equal(t, tt.want.Location, res.Header.Get("Location"))
			}
		})
	}
}
