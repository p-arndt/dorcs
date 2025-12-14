// Package server provides HTTP middleware and utilities.
package server

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// LoggingMiddleware wraps an http.Handler to log request method, path, status, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, status: 200}

		next.ServeHTTP(lrw, r)

		dur := time.Since(start)

		// Skip logging for SSE reload endpoint to avoid log spam
		if strings.HasSuffix(r.URL.Path, "/__reload") {
			// Only log initial connection and errors
			if dur < 100*time.Millisecond || lrw.status >= 400 {
				log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, lrw.status, dur.Truncate(time.Millisecond))
			}
			return
		}

		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, lrw.status, dur.Truncate(time.Millisecond))
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code before delegating to the underlying ResponseWriter.
func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Flush implements http.Flusher for SSE support.
func (w *loggingResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
