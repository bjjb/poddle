package main

import (
	"net/http"
	"strings"
)

// A whitelist is a list of strings that can be joined into a single string
type whitelist []string

// Join joins the whitelisted strings together with commas
func (w whitelist) Join() string {
	return strings.Join(w, ",")
}

// AllowedCORSMethods is a whitelist of HTTP methods for CORS
var AllowedCORSMethods = whitelist([]string{
	http.MethodPost,
	http.MethodGet,
	http.MethodOptions,
})

// AllowedCORSHeaders is a whiteist of HTTP headers for CORS
var AllowedCORSHeaders = whitelist([]string{
	"Accept",
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"X-CSRF-Token",
	"Authorization",
})

// cors sets CORS headers, and then returns true if the request isn't a
// preflight request.
func cors(w http.ResponseWriter, r *http.Request) bool {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", AllowedCORSMethods.Join())
		w.Header().Set("Access-Control-Allow-Headers", AllowedCORSHeaders.Join())
	}
	return r.Method == http.MethodOptions
}
