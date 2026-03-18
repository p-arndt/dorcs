// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/p-arndt/dorcs/internal/site"
)

var (
	// staticExts is a set of file extensions that are considered static assets.
	staticExts = map[string]bool{
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
		".css":   true,
		".js":    true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
	}

	// mimeOverrides contains MIME type overrides for cases where we want different behavior than Go's default.
	mimeOverrides = map[string]string{
		".zip": "application/zip",                       // Go returns "application/x-zip-compressed"
		".xml": "application/xml",                       // Go returns "text/xml; charset=utf-8"
		".js":  "application/javascript; charset=utf-8", // Ensure charset for JS
	}

	// mimeFallbacks contains MIME types for extensions not in Go's default mime database.
	mimeFallbacks = map[string]string{
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".eot":   "application/vnd.ms-fontobject",
	}
)

// isStaticAsset checks if a file path looks like a static asset (image, etc.)
func isStaticAsset(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return staticExts[ext]
}

// tryServeStaticAsset attempts to serve a static asset.
// Serves from the site's root directory (which may be a language/version subdirectory).
// Returns true if the file was found and served, false otherwise.
func (h *Handler) tryServeStaticAsset(w http.ResponseWriter, _ *http.Request, relPath string, targetSite *site.Site) bool {
	// Only serve files that look like static assets
	if !isStaticAsset(relPath) {
		return false
	}

	// Use the site's actual root directory (e.g., docs/en/ for default language using its folder)
	siteRootDir := targetSite.RootDir

	// Clean the relative path to prevent directory traversal
	cleanRelPath := filepath.Clean(filepath.FromSlash(relPath))
	if _, err := sanitizeRelPath(cleanRelPath); err != nil {
		return false
	}

	// First try in the site's root directory (e.g., docs/en/logo.png)
	_, resolved, err := resolveExistingPathWithin(siteRootDir, cleanRelPath)
	if err != nil {
		// If not found in site root, fall back to base docs directory (for shared assets)
		h.mu.RLock()
		docsDir := h.cfg.DocsDir
		h.mu.RUnlock()
		_, resolved, err = resolveExistingPathWithin(docsDir, cleanRelPath)
		if err != nil {
			content, fetchErr := targetSite.FetchGitHubAsset(relPath)
			if fetchErr != nil {
				return false
			}
			return h.serveStaticBytes(w, content, cleanRelPath)
		}
	}

	// Check if file exists
	if fi, err := os.Lstat(resolved); err != nil || fi.Mode()&os.ModeSymlink != 0 {
		return false
	}
	file, err := os.Open(resolved)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		return false
	}

	return h.serveStaticFile(w, file, stat, cleanRelPath)
}

// serveStaticFile serves a static file with appropriate headers.
func (h *Handler) serveStaticFile(w http.ResponseWriter, file *os.File, stat os.FileInfo, relPath string) bool {
	h.setStaticHeaders(w, relPath, stat.Size())

	// Reset file pointer to beginning (in case it was already read)
	file.Seek(0, 0)

	// Copy file content to response
	_, err := io.Copy(w, file)
	return err == nil
}

func (h *Handler) serveStaticBytes(w http.ResponseWriter, content []byte, relPath string) bool {
	h.setStaticHeaders(w, relPath, int64(len(content)))

	_, err := io.Copy(w, bytes.NewReader(content))
	return err == nil
}

func (h *Handler) setStaticHeaders(w http.ResponseWriter, relPath string, contentLength int64) {
	// Set content type based on extension
	ext := strings.ToLower(filepath.Ext(relPath))
	contentType := ""

	// Check overrides first
	if override, ok := mimeOverrides[ext]; ok {
		contentType = override
	} else if fallback, ok := mimeFallbacks[ext]; ok {
		// Use fallback for extensions not in Go's mime database
		contentType = fallback
	} else {
		// Use Go's mime package
		contentType = mime.TypeByExtension(ext)
	}

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Set cache headers
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	// Set content length (use strconv for better performance than fmt.Sprintf)
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
}
