package middleware

import (
	"log"
	"net/http"
	"time"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		latency := time.Since(start)
		log.Printf("remote address: %s; method: %s; URL: %s; latency %s\n", r.RemoteAddr, r.Method, r.URL, latency)
	})
}
