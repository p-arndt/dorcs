package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// SanitizeBasePrefix normalizes a base URL prefix.
// Rules:
// - empty => ""
// - ensure leading slash
// - remove trailing slash
// - disallow "." and ".." segments
func SanitizeBasePrefix(in string) string {
	s := strings.TrimSpace(in)
	if s == "" {
		return ""
	}
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	for strings.HasSuffix(s, "/") && s != "/" {
		s = strings.TrimSuffix(s, "/")
	}
	if s == "/" {
		return ""
	}
	// Safety checks (avoid path traversal)
	parts := strings.Split(s, "/")
	for _, p := range parts {
		if p == "." || p == ".." {
			log.Fatalf("invalid base-url: %q", in)
		}
	}
	return s
}

// CachingFileServer creates a file server with proper ETag and cache support for embedded files.
func CachingFileServer(fsys fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		// Open the file
		file, err := fsys.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Read file content to compute ETag
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Compute ETag from file content (SHA256 hash)
		hash := sha256.Sum256(data)
		etag := `"` + hex.EncodeToString(hash[:16]) + `"` // Use first 16 bytes for shorter ETag

		// Handle conditional requests BEFORE setting other headers (to avoid writing headers if 304)
		// Check ETag first (more reliable) - but only if browser isn't forcing no-cache
		if r.Header.Get("Cache-Control") != "no-cache" {
			if match := r.Header.Get("If-None-Match"); match != "" {
				// ETag comparison - handle both quoted and unquoted, and comma-separated lists
				cleanETag := strings.Trim(etag, `"`)
				cleanMatch := strings.Trim(match, `"`)
				if strings.Contains(cleanMatch, cleanETag) || match == etag || strings.Contains(match, cleanETag) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}

		// Set content type
		switch {
		case strings.HasSuffix(path, ".css"):
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case strings.HasSuffix(path, ".js"):
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case strings.HasSuffix(path, ".ico"):
			w.Header().Set("Content-Type", "image/x-icon")
		case strings.HasSuffix(path, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg"):
			w.Header().Set("Content-Type", "image/jpeg")
		case strings.HasSuffix(path, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		}

		// Set cache headers (1 year, immutable) - only if not 304
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("ETag", etag)
		w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

		// Serve the file content
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		if r.Method != http.MethodHead {
			w.Write(data)
		}
	})
}
