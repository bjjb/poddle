package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func Test_cors(t *testing.T) {
	for _, method := range []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPut,
		http.MethodPost,
		http.MethodPatch,
		http.MethodDelete,
	} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodOptions, "http://foo.co", nil)
			req.Header.Set("Origin", "http://foo.co")
			rec := httptest.NewRecorder()
			cors(rec, req)
			resp := rec.Result()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected %d; got %d", http.StatusOK, rec.Code)
			}
			v := resp.Header.Get("Access-Control-Allow-Origin")
			if v != "http://foo.co" {
				t.Fatalf("expected %q, got %q", "http://foo.co", v)
			}
			v = resp.Header.Get("Access-Control-Allow-Methods")
			if v != AllowMethods.String() {
				t.Fatalf("expected %q, got %q", AllowMethods.String(), v)
			}
			v = resp.Header.Get("Access-Control-Allow-Headers")
			if v != AllowHeaders.String() {
				t.Fatalf("expected %q, got %q", AllowHeaders.String(), v)
			}
		})
	}
}

func Test_getEnvString(t *testing.T) {
	expected := os.Getenv("GOPATH")
	t.Run("existing values are returned", func(t *testing.T) {
		if actual := getEnvString("GOPATH", "!wrong!"); actual != expected {
			t.Fatalf("expected %q, got %q", expected, actual)
		}
	})
	t.Run("missing values use the fallback", func(t *testing.T) {
		expected := "fallback value"
		if actual := getEnvString("blob blub", expected); actual != expected {
			t.Fatalf("expected %q, got %q", expected, actual)
		}
	})
}

func Test_getEnvDuration(t *testing.T) {
	expected := time.Second * 7
	t.Run("gets a duration from env or fallback", func(t *testing.T) {
		if actual := getEnvDuration("missing value", expected); actual != expected {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})
}

func Test_get(t *testing.T) {
	t.Run("gives 400 if the URI isn't supplied", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		get(w, r)
		expected := http.StatusBadRequest
		resp := w.Result()
		if resp.StatusCode != expected {
			t.Fatalf("expected %d, got %d", expected, resp.StatusCode)
		}
	})
	t.Run("gives 400 if the URI is malformed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		r.URL.Query().Set("uri", "bad uri")
		get(w, r)
		expected := http.StatusBadRequest
		resp := w.Result()
		if resp.StatusCode != expected {
			t.Fatalf("expected %d, got %d", expected, resp.StatusCode)
		}
	})
	t.Run("proxies to the given URI ", func(t *testing.T) {
		expected := http.StatusOK
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "proxied!", expected)
		}))
		u := fmt.Sprintf("https://example.com/?uri=%s", url.QueryEscape(ts.URL))
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, u, nil)
		get(w, r)
		resp := w.Result()
		if resp.StatusCode != expected {
			t.Errorf("expected %d, got %d", expected, resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != "proxied!\n" {
			t.Errorf("expected %q, got %q", "proxied!\n", string(body))
		}
	})
}

func Test_convert(t *testing.T) {
	t.Run("gives 400 if the URI isn't supplied", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		convert(w, r)
		expected := http.StatusBadRequest
		resp := w.Result()
		if resp.StatusCode != expected {
			t.Fatalf("expected %d, got %d", expected, resp.StatusCode)
		}
	})
	t.Run("gives 400 if the URI is malformed", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		r.URL.Query().Set("uri", "bad uri")
		convert(w, r)
		expected := http.StatusBadRequest
		resp := w.Result()
		if resp.StatusCode != expected {
			t.Fatalf("expected %d, got %d", expected, resp.StatusCode)
		}
	})
}
