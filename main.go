package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gregjones/httpcache"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// A Poddle provides HTTP handlers with settings
type Poddle struct {
	db         *sql.DB
	httpClient *http.Client
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// start starts a server.
func start() {
	srv := &Server{}
	srv.Start()
}

// openDB opens a database, by determining the driver from the string, and
// then calling sql.Open.
func openDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, nil
	}
	if regexp.MustCompile("^sqlite3?:").MatchString(dsn) {
		return sql.Open("sqlite3", strings.SplitN(dsn, ":", 2)[1])
	}
	if regexp.MustCompile("^postgres://").MatchString(dsn) {
		return sql.Open("postgres", dsn)
	}
	return nil, fmt.Errorf("invalid database connection string: %q", dsn)
}

// getEnvString tries to obtain the value of name from the environment, and
// returns fallback if it's empty or not present.
func getEnvString(name, fallback string) string {
	if val, found := os.LookupEnv(name); found {
		return val
	}
	return fallback
}

// getEnvDuration uses getEnvString and parses the values as a duration.
func getEnvDuration(name string, fallback time.Duration) time.Duration {
	if val, found := os.LookupEnv(name); found {
		result, err := time.ParseDuration(val)
		if err == nil {
			return result
		}
		f := "WARN: '%s' ParseDuration(%s) => %e; falling back to %s"
		log.Printf(f, name, val, err, fallback.String())
	}
	return fallback
}

// getEnvDirectory uses getEnvString, and ensures that the string refers to a
// directory which exists.
func getEnvDirectory(name, fallback string) string {
	dirName := getEnvString(name, fallback)
	fileInfo, err := os.Stat(dirName)
	if err != nil {
		log.Fatal(err)
	}
	if !fileInfo.IsDir() {
		log.Fatalf("not a directory: %s", dirName)
	}
	return dirName
}

// parseURL always returns a URL or panics.
func parseURL(input string) *url.URL {
	url, err := url.Parse(input)
	if err != nil {
		panic(err)
	}
	return url
}

// get simply fetches the url in the request's query params, and sends back
// the result.
func get(w http.ResponseWriter, r *http.Request) {
	if cors(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	param := r.URL.Query().Get("uri")
	if param == "" {
		http.Error(w, "missing uri parameter", http.StatusBadRequest)
		return
	}

	uri, err := url.Parse(param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpClient := http.Client{Transport: httpcache.NewMemoryCacheTransport()}
	resp, err := httpClient.Get(uri.String())
	if err != nil {
		log.Printf("✗ %s %e", uri.String(), err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for name, values := range resp.Header {
		log.Printf("⇒ %q: %v", name, values)
	}

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("✓ (%d bytes) %s %s", n, uri.String(), resp.Status)
}

// convert uses ffmpeg to convert the content at the url in the request's
// query params and sends back the result.
func convert(w http.ResponseWriter, r *http.Request) {
	if cors(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	param := r.URL.Query().Get("uri")
	if param == "" {
		http.Error(w, "missing uri parameter", http.StatusBadRequest)
		return
	}

	uri, err := url.Parse(param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpClient := http.Client{Transport: httpcache.NewMemoryCacheTransport()}
	resp, err := httpClient.Get(uri.String())
	if err != nil {
		log.Printf("✗ %s %e", uri.String(), err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := FFMpeg(w, resp.Body); err != nil {
		log.Print(err)
		http.Error(w, "something went awry...", http.StatusInternalServerError)
	}
	log.Printf("✓ %s %s", uri.String(), resp.Status)
}
