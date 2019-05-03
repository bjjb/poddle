package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	strTimes := func(s string, n int) string {
		r := s
		for i := 0; i < n; i++ {
			r = r + s
		}
		return r
	}

	t.Run("iTunes", func(t *testing.T) {

		s := &iTunesSearcher{}

		t.Run("newRequest", func(t *testing.T) {
			for _, tc := range []struct{ q, term, err string }{
				{"foo", "foo", ""},
				{"foo bar", "foo+bar", ""},
				{"", "", "q cannot be blank"},
				{"l" + strTimes("o", 256) + "ng", "", "q cannot be longer than 255 characters"},
			} {
				t.Run(tc.q, func(t *testing.T) {
					req, err := s.newRequest(tc.q)
					if tc.err != "" {
						if err == nil || err.Error() != tc.err {
							t.Fatalf("expected error %q, got %q", tc.err, err)
						}
						return
					}
					if err != nil {
						t.Fatal(err)
					}
					if req.Method != http.MethodGet {
						t.Errorf("expected %s, got %s", http.MethodGet, req.Method)
					}
					exp := "https://itunes.apple.com/search?entity=podcast&term=" + tc.term
					if req.URL.String() != exp {
						t.Errorf("expected %s, got %s", exp, req.URL.String())
					}
					exp = "application/json"
					if req.Header.Get("Accept") != exp {
						t.Errorf("expected %s, got %s", exp, req.Header.Get("Accept"))
					}
				})
			}
		})

		t.Run("parseResponse", func(t *testing.T) {
			t.Run("200", func(t *testing.T) {
				t.Run("valid", func(t *testing.T) {
					json := `{
						"results":
							[
								{"collectionName":"T"},
								{"feedUrl":"L"},
								{"artworkUrl100":"1"},
								{"artworkUrl600":"6"}
							]
						}`
					body := ioutil.NopCloser(strings.NewReader(json))
					resp := &http.Response{StatusCode: 200, Body: body}
					results, err := s.parseResponse(resp)
					if err != nil {
						t.Fatal(err)
					}
					if results == nil {
						t.Fatal("expected results")
					}
					if len(results) != 4 {
						t.Fatalf("expected 2 results, got %d", len(results))
					}
					if results[0].Title != "T" {
						t.Fatalf("expected title=%q, got %q", "T", results[0].Title)
					}
					if results[1].URL != "L" {
						t.Fatalf("expected URL=%q, got %q", "L", results[1].URL)
					}
					if results[2].Image.URL != "1" {
						t.Fatalf("expected Image.URL=%q, got %q", "1", results[2].Image.URL)
					}
					if results[3].Image.URL != "6" {
						t.Fatalf("expected Image.URL=%q, got %q", "6", results[3].Image.URL)
					}
				})
			})
			t.Run("400", func(t *testing.T) {
				_, err := s.parseResponse(&http.Response{StatusCode: 414})
				if err == nil {
					t.Fatal("expected an error, not none")
				}
			})
			t.Run("invalid JSON", func(t *testing.T) {
				body := ioutil.NopCloser(strings.NewReader("invalid JSON"))
				_, err := s.parseResponse(&http.Response{StatusCode: 200, Body: body})
				if err == nil {
					t.Fatal("expected an error, not none")
				}
			})
		})
	})
}
