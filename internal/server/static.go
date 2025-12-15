// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
// First checks the root directory (where dorcs is running), then falls back to docs directory.
// Returns true if the file was found and served, false otherwise.
func (h *Handler) tryServeStaticAsset(w http.ResponseWriter, _ *http.Request, relPath string) bool {
	// Only serve files that look like static assets
	if !isStaticAsset(relPath) {
		return false
	}

	h.mu.RLock()
	rootDir := h.cfg.RootDir
	docsDir := h.cfg.DocsDir
	h.mu.RUnlock()

	// Clean the relative path to prevent directory traversal
	cleanRelPath := filepath.Clean(filepath.FromSlash(relPath))
	if strings.HasPrefix(cleanRelPath, "..") || filepath.IsAbs(cleanRelPath) {
		return false
	}

	// Try root directory first (for logo, favicon, etc.)
	var filePath string
	var absBaseDir string
	var err error

	if rootDir != "" {
		filePath = filepath.Join(rootDir, cleanRelPath)
		absBaseDir, err = filepath.Abs(rootDir)
		if err == nil {
			absFilePath, err := filepath.Abs(filePath)
			if err == nil {
				// Security: ensure the file is within the root directory
				if strings.HasPrefix(absFilePath, absBaseDir+string(filepath.Separator)) || absFilePath == absBaseDir {
					if file, err := os.Open(filePath); err == nil {
						defer file.Close()
						if stat, err := file.Stat(); err == nil && !stat.IsDir() {
							// Found in root directory, serve it
							return h.serveStaticFile(w, file, stat, cleanRelPath)
						}
					}
				}
			}
		}
	}

	// Fall back to docs directory
	filePath = filepath.Join(docsDir, cleanRelPath)
	absDocsDir, err := filepath.Abs(docsDir)
	if err != nil {
		return false
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return false
	}
	// Security: ensure the file is within the docs directory
	if !strings.HasPrefix(absFilePath, absDocsDir+string(filepath.Separator)) && absFilePath != absDocsDir {
		return false
	}

	// Check if file exists
	file, err := os.Open(filePath)
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
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	// Reset file pointer to beginning (in case it was already read)
	file.Seek(0, 0)

	// Copy file content to response
	_, err := io.Copy(w, file)
	return err == nil
}
