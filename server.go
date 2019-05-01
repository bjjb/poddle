package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// IdleTimeout is the length of time a server will remain idle before closing
var IdleTimeout = time.Second * 60

// ReadTimeout is the max read-time for a HTTP request
var ReadTimeout = time.Second * 30

// WriteTimeout is the max write-time for a HTTP response
var WriteTimeout = time.Minute * 10

// WaitTimeout is the time a killed server will wait before dying
var WaitTimeout = WriteTimeout * 2

// Handler is the default handler
var Handler *http.ServeMux

// A Server can start a configured http.Server.
type Server struct {
	Handler                                             http.Handler
	ReadTimeout, WriteTimeout, IdleTimeout, WaitTimeout time.Duration
	Addr, TLSCertFile, TLSKeyFile                       string
	Logger                                              *log.Logger
}

// Start starts the server. It will use the configured Addr, and timeouts, or
// or sensible defaults (which might be overridden by environment variables if
// they are zero.
// If TLSCertFile and TLSKeyFile are both set to filenames (which exist), then
// the server will serve TLS (and Addr will default to 443).
// The Mux will default to the package's Mux.
// The server can be killed with an INT signal. Errors are logged to the
// server's logger (which defaults to standard out).
func (s *Server) Start() {
	readTimeout, writeTimeout, idleTimeout, waitTimeout := s.Timeouts()
	tlsOn, tlsCertFile, tlsKeyFile, addr := s.TLSConfig()

	logger := s.Logger
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	handler := s.Handler
	if handler == nil {
		handler = Handler
	}

	srv := &http.Server{
		Addr:         addr,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      handler,
	}

	start := func() error {
		if tlsOn {
			return srv.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
		}
		return srv.ListenAndServe()
	}

	go func() {
		logger.Printf("poddle server listening on %s", addr)
		if err := start(); err != nil {
			logger.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	logger.Println("Goodbye.")
}

// Timeouts gets the server's timeouts, either configured or from the
// environment.
func (s *Server) Timeouts() (idle, read, write, wait time.Duration) {
	d := func(cur, fb time.Duration, env string) time.Duration {
		if cur != 0 {
			return cur
		}
		if v, found := os.LookupEnv(env); found {
			d, err := time.ParseDuration(v)
			if err == nil {
				return d
			}
			f := "WARN: '%s' ParseDuration(%s) => %e; falling back to %s"
			log.Printf(f, env, v, err, fb.String())
		}
		return fb
	}
	idle = d(s.IdleTimeout, IdleTimeout, "IDLE_TIMEOUT")
	read = d(s.ReadTimeout, ReadTimeout, "READ_TIMEOUT")
	write = d(s.WriteTimeout, WriteTimeout, "WRITE_TIMEOUT")
	wait = d(s.WaitTimeout, WaitTimeout, "WAIT_TIMEOUT")
	return
}

// TLSConfig gets the configured or default TLS files from the server, a
// boolean indicating whether TLS is enabled or not, and a default listen
// address.
func (s *Server) TLSConfig() (tlsOn bool, tlsCertFile, tlsKeyFile, addr string) {
	var found bool
	f := func(cur, fb, env string) string {
		if cur != "" {
			return cur
		}
		if v, found := os.LookupEnv(env); found {
			return v
		}
		return fb
	}
	tlsCertFile = f(s.TLSCertFile, "", "TLS_CERT_FILE")
	tlsKeyFile = f(s.TLSKeyFile, "", "TLS_KEY_FILE")

	if addr, found = os.LookupEnv("ADDR"); !found {
		addr = ":80"
		if tlsOn = tlsCertFile != "" && tlsKeyFile != ""; tlsOn {
			addr = ":443"
		}
	}
	return
}

func init() {
	Handler = &http.ServeMux{}
}
