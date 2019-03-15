package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

// A whitelist is a list of strings that can be joined into a single string
type whitelist []string

// String joins the whitelisted strings together with commas
func (w whitelist) String() string {
	return strings.Join(w, ",")
}

// AllowMethods is a whitelist of HTTP methods for CORS
var AllowMethods = whitelist([]string{
	http.MethodPost,
	http.MethodGet,
	http.MethodOptions,
})

// AllowHeaders is a whiteist of HTTP headers for CORS
var AllowHeaders = whitelist([]string{
	"Accept",
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"X-CSRF-Token",
	"Authorization",
})

// main starts a server.
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/get/", get)
	mux.HandleFunc("/convert/", convert)
	mux.Handle("/", http.FileServer(http.Dir(getEnvDirectory("APP", "app"))))

	srv := &http.Server{
		Addr:         getEnvString("ADDR", ":8080"),
		IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", time.Second*60),
		ReadTimeout:  getEnvDuration("READ_TIMEOUT", time.Second*15),
		WriteTimeout: getEnvDuration("WRITE_TIMEOUT", time.Minute*8),
		Handler:      mux,
	}

	go func() {
		log.Printf("poddle server listening on %s...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	waitTimeout := getEnvDuration("WAIT_TIMEOUT", time.Second*15)
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Goodbye...")
	// os.Exit(0)
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
		if result, err := time.ParseDuration(val); err != nil {
			f := "WARN: '%s' ParseDuration(%s) => %e; falling back to %s"
			log.Printf(f, name, val, err, fallback.String())
			return fallback
		} else {
			return result
		}
	}
	return fallback
}

// getEnvDuration uses getEnvString, and ensures that the string refers to a
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
		http.Error(w, "mixsing uri parameter", http.StatusBadRequest)
		return
	}

	uri, err := url.Parse(param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := http.Get(uri.String())
	if err != nil {
		log.Printf("✗ %s %e", uri.String(), err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	io.Copy(w, resp.Body)
	log.Printf("✓ %s %s", uri.String(), resp.Status)
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
		http.Error(w, "mixsing uri parameter", http.StatusBadRequest)
		return
	}

	uri, err := url.Parse(param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := http.Get(uri.String())
	if err != nil {
		log.Printf("✗ %s %e", uri.String(), err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	if err := ffmpeg(w, resp.Body); err != nil {
		log.Print(err)
		http.Error(w, "something went awry...", http.StatusInternalServerError)
	}
	log.Printf("✓ %s %s", uri.String(), resp.Status)
}

// cors sets CORS headers, and then returns true if the request isn't a
// preflight request.
func cors(w http.ResponseWriter, r *http.Request) bool {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", AllowMethods.String())
		w.Header().Set("Access-Control-Allow-Headers", AllowHeaders.String())
	}
	if r.Method == http.MethodOptions {
		return true
	}
	return false
}

// ffmpeg sets up an ffmpeg process to convert stdin to opus on stdout
func ffmpeg(stdout io.Writer, stdin io.ReadCloser) error {
	defer stdin.Close()

	cmd := exec.Command(
		"ffmpeg", "-hide_banner", "-loglevel", "warning", "-i", "-", "-f", "opus",
		"-vn", "-c:a", "libopus", "-b:a", "16k", "-application", "voip", "-",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = stdout
	cmd.Stdin = stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	log.Print("converting...")

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
