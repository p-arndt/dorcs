// Package server provides HTTP middleware and utilities.
package server

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// isStaticRequest checks if a request is for a static asset that should not be logged.
func isStaticRequest(path string) bool {
	// Skip logging for static files served from /static/ path
	if strings.Contains(path, "/static/") {
		return true
	}

	// Skip logging for files with static file extensions
	ext := strings.ToLower(filepath.Ext(path))
	staticExts := map[string]bool{
		".css":   true,
		".js":    true,
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".svg":   true,
		".webp":  true,
		".ico":   true,
		".pdf":   true,
		".zip":   true,
		".json":  true,
		".xml":   true,
		".txt":   true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
		".otf":   true,
	}
	return staticExts[ext]
}

// LoggingMiddleware wraps an http.Handler to log request method, path, status, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, status: 200}

		next.ServeHTTP(lrw, r)

		dur := time.Since(start)

		// Skip logging for static assets to avoid log spam
		if isStaticRequest(r.URL.Path) {
			// Only log errors for static files
			if lrw.status >= 400 {
				log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, lrw.status, dur.Truncate(time.Millisecond))
			}
			return
		}

		// Skip logging for SSE reload endpoint to avoid log spam
		if strings.HasSuffix(r.URL.Path, "/__reload") {
			// Only log initial connection and errors
			if dur < 100*time.Millisecond || lrw.status >= 400 {
				log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, lrw.status, dur.Truncate(time.Millisecond))
			}
			return
		}

		// Skip logging for health check endpoint
		if strings.HasSuffix(r.URL.Path, "/api/health") {
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
