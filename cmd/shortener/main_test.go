package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0xdreamerr/url-shortener/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetShortUrl(t *testing.T) {
	config.SetConfig()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Regular test #1",
			url:  "https://practicum.yandex.ru/",
			want: want{
				code:        201,
				response:    "http://localhost:8080/42b3e75",
				contentType: "text/plain",
			},
		},
		{
			name: "Regular test #2",
			url:  "https://youtube.com/",
			want: want{
				code:        201,
				response:    "http://localhost:8080/e1d5e2c",
				contentType: "text/plain",
			},
		},
		{
			name: "Regular test #3",
			url:  "https://vk.com/",
			want: want{
				code:        201,
				response:    "http://localhost:8080/b2f21cf",
				contentType: "text/plain",
			},
		},
		{
			name: "Wrong request method",
			url:  "https://practicum.yandex.ru/",
			want: want{
				code:        405,
				response:    "Only POST requests are allowed!\n",
				contentType: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(test.url)

			if test.name == "Wrong request method" {
				request := httptest.NewRequest(http.MethodGet, "/", bodyReader)

				w := httptest.NewRecorder()
				getShortURL(w, request)

				res := w.Result()

				assert.Equal(t, res.StatusCode, test.want.code)

				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Equal(t, string(resBody), test.want.response)

			} else {
				request := httptest.NewRequest(http.MethodPost, "/", bodyReader)

				w := httptest.NewRecorder()
				getShortURL(w, request)

				res := w.Result()

				assert.Equal(t, res.StatusCode, test.want.code)

				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Equal(t, string(resBody), test.want.response)
				assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
			}

		})
	}
}

func TestRedirectTo(t *testing.T) {
	type want struct {
		code     int
		location string
		response string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Regular test #1",
			url:  "/42b3e75",
			want: want{
				code:     307,
				location: "https://practicum.yandex.ru/",
			},
		},
		{
			name: "Regular test #2",
			url:  "/e1d5e2c",
			want: want{
				code:     307,
				location: "https://youtube.com/",
			},
		},
		{
			name: "Wrong request method",
			url:  "/e1d5e2c",
			want: want{
				code:     405,
				response: "Only GET requests are allowed!\n",
			},
		},
		{
			name: "ShortURL not exist",
			url:  "/e1d5e23",
			want: want{
				code:     404,
				response: "ShortURL not found\n",
			},
		},
	}

	// creating short links
	for _, test := range tests {
		bodyReader := strings.NewReader(test.want.location)

		request := httptest.NewRequest(http.MethodPost, "/", bodyReader)

		w := httptest.NewRecorder()
		getShortURL(w, request)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "Wrong request method" {
				request := httptest.NewRequest(http.MethodPost, "/", nil)

				w := httptest.NewRecorder()
				redirectTo(w, request)

				res := w.Result()

				assert.Equal(t, res.StatusCode, test.want.code)

				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Equal(t, string(resBody), test.want.response)

			} else {
				request := httptest.NewRequest(http.MethodGet, test.url, nil)

				w := httptest.NewRecorder()
				redirectTo(w, request)

				res := w.Result()

				assert.Equal(t, res.StatusCode, test.want.code)

				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Equal(t, string(resBody), test.want.response)
				assert.Equal(t, res.Header.Get("Location"), test.want.location)
			}
		})
	}
}
