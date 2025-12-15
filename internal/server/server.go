// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"dorcs-v2/internal/config"
	"dorcs-v2/internal/site"
)

// Config contains runtime configuration for the HTTP app.
type Config struct {
	DocsDir   string
	RootDir   string // Root directory where dorcs is running (for static assets like logo/favicon)
	SiteTitle string
	BasePath  string
	Cache     bool
	HideDraft bool

	Site         *site.Site // Default language site (for backward compatibility)
	DocumentTmpl *template.Template

	// SiteConfig holds the loaded dorcs config for theming/branding
	SiteConfig *config.Config

	// Sites is a map of language code to Site instance (for multi-lingual support)
	// If nil or empty, only the default Site is used
	Sites map[string]*site.Site

	// ReloadBroadcaster enables live reload in watch mode
	ReloadBroadcaster *site.ReloadBroadcaster

	// Version is the version identifier for dorcs
	Version string
}

// Handler is an http.Handler that routes requests to doc handlers.
type Handler struct {
	cfg Config
	mu  sync.RWMutex
}

// New constructs the HTTP handler.
func New(cfg Config) *Handler {
	return &Handler{cfg: cfg}
}

// UpdateConfig updates the handler's configuration atomically.
// This is thread-safe and can be called from a watcher goroutine.
func (h *Handler) UpdateConfig(cfg *config.Config) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cfg.SiteConfig = cfg
	// Update site title if config provides one
	if cfg.Site.Title != "" {
		h.cfg.SiteTitle = cfg.Site.Title
	}
}

// NavItem represents a navigation tree item for sidebar rendering.
type NavItem struct {
	Title    string
	Path     string // URL path, begins with "/" (e.g. "/guide/getting-started")
	IsDir    bool
	Children []NavItem
}

// DocPageModel is the view model for document pages.
type DocPageModel struct {
	SiteTitle        string
	DocTitleFallback string
	CanonicalURL     string
	LastModified     string
	TOCHTML          template.HTML
	HTML             template.HTML

	// TOC is optional HTML for a right-side "On this page" table of contents.
	TOC template.HTML

	// CurrentPath is the request path used to highlight the active item in the sidebar.
	CurrentPath string

	// DocPath is the document path without language prefix (e.g., "/getting-started" or "/")
	DocPath string

	// Navigation sidebar tree (layout expects `.Nav.Nodes`)
	Nav struct {
		Nodes []NavItem
	}

	// RootTitle is the title of the root index.md page (used for the "Home" link)
	RootTitle string

	Meta struct {
		Title       string
		Description string
		Date        string
		Tags        []string
		SourcePath  string
	}

	// Config holds the site configuration for theming and branding
	Config *config.Config

	// ThemeCSS contains generated CSS custom properties for theming
	ThemeCSS template.CSS

	// BasePath for building URLs
	BasePath string

	// LiveReload indicates whether live reload is enabled
	LiveReload bool

	// Version is the version identifier for dorcs
	Version string

	// CurrentLanguage is the language code for the current page
	CurrentLanguage string
}

// ServeHTTP routes requests:
// - GET {base}/           => render docs/index.md
// - GET {base}/<path>     => render docs/<path>.md
//
// URLs are extensionless. If a folder contains an index.md, access it as "/folder".
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Path
	h.mu.RLock()
	basePath := h.cfg.BasePath
	siteConfig := h.cfg.SiteConfig
	h.mu.RUnlock()
	if basePath != "" {
		if !strings.HasPrefix(reqPath, basePath) {
			http.NotFound(w, r)
			return
		}
		reqPath = strings.TrimPrefix(reqPath, basePath)
		if reqPath == "" {
			reqPath = "/"
		}
	}

	// Detect language from URL path
	var currentLang string
	var docPath string
	if reqPath == "/" {
		docPath = "/"
	} else {
		// Extract first path segment as potential language code
		parts := strings.Split(strings.TrimPrefix(reqPath, "/"), "/")
		if len(parts) > 0 && parts[0] != "" {
			// Check if first segment is a valid language code
			if siteConfig != nil && siteConfig.IsLanguageEnabled(parts[0]) {
				currentLang = parts[0]
				// Remove language prefix from path
				if len(parts) > 1 {
					docPath = "/" + strings.Join(parts[1:], "/")
				} else {
					docPath = "/"
				}
			} else {
				// Not a language code, treat as document path
				docPath = reqPath
			}
		} else {
			docPath = reqPath
		}
	}

	// Get the appropriate Site instance for this language
	var targetSite *site.Site
	h.mu.RLock()
	if h.cfg.Sites != nil && len(h.cfg.Sites) > 0 {
		// Default language uses empty key in Sites map
		siteKey := currentLang
		if currentLang == "" {
			siteKey = "" // Default language
		}
		if langSite, ok := h.cfg.Sites[siteKey]; ok {
			targetSite = langSite
		}
	}
	// Fall back to default site if no language-specific site found
	if targetSite == nil {
		targetSite = h.cfg.Site
		currentLang = "" // Use default language
	}
	h.mu.RUnlock()

	// Root -> docs/index.md
	if docPath == "/" {
		h.handleDocByKeyWithSite(w, r, "", targetSite, currentLang, docPath)
		return
	}

	// Everything else: strip leading "/" and try to resolve to doc.
	rel := strings.TrimPrefix(docPath, "/")
	rel = strings.Trim(rel, "/")
	if rel == "" {
		http.NotFound(w, r)
		return
	}

	// Support folder index: /guide -> docs/guide/index.md
	if strings.HasSuffix(docPath, "/") {
		h.handleDocByKeyWithSite(w, r, path.Clean(rel), targetSite, currentLang, docPath)
		return
	}

	// Check if this is a static asset (image, etc.) before trying markdown
	if h.tryServeStaticAsset(w, r, rel) {
		return
	}

	// First: try docs/<rel>.md
	if h.tryServeDocWithSite(w, r, rel, targetSite, currentLang) {
		return
	}
	// Second: try docs/<rel>/index.md (folder index)
	if h.tryServeDocWithSite(w, r, path.Clean(rel), targetSite, currentLang) {
		return
	}

	http.NotFound(w, r)
}

// isStaticAsset checks if a file path looks like a static asset (image, etc.)
func isStaticAsset(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	staticExts := map[string]bool{
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
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	case ".pdf":
		w.Header().Set("Content-Type", "application/pdf")
	case ".zip":
		w.Header().Set("Content-Type", "application/zip")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".xml":
		w.Header().Set("Content-Type", "application/xml")
	case ".txt":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".woff":
		w.Header().Set("Content-Type", "font/woff")
	case ".woff2":
		w.Header().Set("Content-Type", "font/woff2")
	case ".ttf":
		w.Header().Set("Content-Type", "font/ttf")
	case ".eot":
		w.Header().Set("Content-Type", "application/vnd.ms-fontobject")
	}

	// Set cache headers
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	// Set content length
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	// Reset file pointer to beginning (in case it was already read)
	file.Seek(0, 0)

	// Copy file content to response
	_, err := io.Copy(w, file)
	return err == nil
}

func (h *Handler) tryServeDocWithSite(w http.ResponseWriter, r *http.Request, key string, targetSite *site.Site, currentLang string) bool {
	if targetSite == nil {
		return false
	}
	key = cleanKey(key)
	if _, ok := targetSite.GetDoc(key); ok {
		// Reconstruct docPath from key for tryServeDocWithSite
		var docPathForCall string
		if key == "" {
			docPathForCall = "/"
		} else {
			docPathForCall = "/" + key
		}
		h.handleDocByKeyWithSite(w, r, key, targetSite, currentLang, docPathForCall)
		return true
	}
	return false
}

func (h *Handler) handleDocByKeyWithSite(w http.ResponseWriter, r *http.Request, key string, targetSite *site.Site, currentLang string, docPath string) {
	h.mu.RLock()
	documentTmpl := h.cfg.DocumentTmpl
	basePath := h.cfg.BasePath
	siteConfig := h.cfg.SiteConfig
	h.mu.RUnlock()

	if targetSite == nil || documentTmpl == nil {
		http.Error(w, "server not configured", http.StatusInternalServerError)
		return
	}

	key = cleanKey(key)

	// Build navigation tree for this site
	nav := h.buildNavItemsWithSite(targetSite)

	// Get root title from navigation tree
	rootTitle := h.getRootTitleWithSite(targetSite)

	// Use site renderer
	rendered, err := targetSite.RenderDoc(key)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var m DocPageModel
	// Read site title with lock
	h.mu.RLock()
	siteTitle := h.cfg.SiteTitle
	h.mu.RUnlock()
	m.SiteTitle = siteTitle
	m.DocTitleFallback = key
	if key == "" {
		m.DocTitleFallback = "index"
	}
	m.HTML = rendered.HTML
	m.TOCHTML = rendered.TocHTML
	m.TOC = rendered.TocHTML
	m.LastModified = rendered.Doc.UpdatedAt.UTC().Format(time.RFC3339)
	m.Nav.Nodes = nav
	m.RootTitle = rootTitle
	// Set current path - for root index, use "/" for default language or "/{lang}/" for others
	if key == "" {
		if currentLang != "" {
			m.CurrentPath = basePath + "/" + currentLang + "/"
		} else {
			m.CurrentPath = basePath + "/"
		}
	} else {
		m.CurrentPath = r.URL.Path
	}
	// DocPath is the document path without language prefix (used for building language links)
	m.DocPath = docPath
	m.BasePath = basePath
	m.CurrentLanguage = currentLang

	m.Meta.Title = rendered.Doc.Title
	m.Meta.Description = rendered.Doc.Description
	m.Meta.Date = formatDate(rendered.Doc.Date)
	m.Meta.Tags = append([]string(nil), rendered.Doc.Tags...)
	m.Meta.SourcePath = rendered.Doc.RelPath

	// Add config and theme CSS if available
	if siteConfig != nil {
		m.Config = siteConfig
		m.ThemeCSS = template.CSS(siteConfig.GenerateThemeCSS())
		// Override site title from config if set
		if siteConfig.Site.Title != "" {
			m.SiteTitle = siteConfig.Site.Title
		}
	}

	// Enable live reload if broadcaster is configured
	h.mu.RLock()
	reloadBroadcaster := h.cfg.ReloadBroadcaster
	version := h.cfg.Version
	h.mu.RUnlock()
	m.LiveReload = reloadBroadcaster != nil
	m.Version = version

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Try "doc" template first (which calls "layout"), fall back to "layout" for compatibility
	if err := documentTmpl.ExecuteTemplate(w, "doc", m); err != nil {
		if err2 := h.cfg.DocumentTmpl.ExecuteTemplate(w, "layout", m); err2 != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// buildNavItemsWithSite converts the given site's NavTree to the template-friendly NavItem slice.
func (h *Handler) buildNavItemsWithSite(targetSite *site.Site) []NavItem {
	if targetSite == nil {
		return nil
	}
	h.mu.RLock()
	hideDraft := h.cfg.HideDraft
	basePath := h.cfg.BasePath
	currentLang := ""
	if targetSite.Language != "" {
		currentLang = targetSite.Language
	}
	h.mu.RUnlock()
	tree := targetSite.NavTree(!hideDraft)
	if tree == nil {
		return nil
	}
	return convertNavNodesWithLang(tree.Children, basePath, currentLang)
}



// getRootTitleWithSite extracts the title from the root index.md page of the given site.
func (h *Handler) getRootTitleWithSite(targetSite *site.Site) string {
	if targetSite == nil {
		return "Home"
	}
	h.mu.RLock()
	hideDraft := h.cfg.HideDraft
	h.mu.RUnlock()
	tree := targetSite.NavTree(!hideDraft)
	if tree == nil {
		return "Home"
	}
	// Root node has the title from index.md
	if tree.Page != nil && tree.Page.Title != "" {
		return tree.Page.Title
	}
	if tree.Name != "" {
		return tree.Name
	}
	return "Home"
}



func convertNavNodesWithLang(nodes []*site.NavNode, basePath string, currentLang string) []NavItem {
	if len(nodes) == 0 {
		return nil
	}
	items := make([]NavItem, 0, len(nodes))
	for _, n := range nodes {
		// n.Name now contains the title from index.md for folders (set in buildNavTree)
		title := n.Name

		// For pages and folders with Page, prefer Page.Title
		if n.Page != nil && n.Page.Title != "" {
			title = n.Page.Title
		}

		item := NavItem{
			Title:    title,
			IsDir:    n.IsDir,
			Children: convertNavNodesWithLang(n.Children, basePath, currentLang),
		}

		// Set path - folders are only clickable if they have a landing page
		pathPrefix := basePath
		if currentLang != "" {
			if pathPrefix != "" {
				pathPrefix += "/" + currentLang
			} else {
				pathPrefix = "/" + currentLang
			}
		}
		if n.IsDir {
			if n.Page != nil {
				item.Path = pathPrefix + "/" + n.Key
			}
			// Path stays empty for folders without index.md
		} else {
			item.Path = pathPrefix + "/" + n.Key
		}

		items = append(items, item)
	}
	return items
}

func cleanKey(in string) string {
	s := strings.TrimSpace(in)
	s = strings.ReplaceAll(s, "\\", "/")
	s = path.Clean("/" + s)
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, "/")
	if s == "." || s == ".." || strings.HasPrefix(s, "../") || strings.Contains(s, "/../") {
		return ""
	}
	return s
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}


// SearchResponse represents the JSON response for search API.
type SearchResponse struct {
	Results []SearchResultItem `json:"results"`
	Query   string             `json:"query"`
}

// SearchResultItem represents a single search result in the API response.
type SearchResultItem struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	Snippet     string `json:"snippet"`
	Score       int    `json:"score"`
	HeadingID   string `json:"heading_id,omitempty"`
	HeadingText string `json:"heading_text,omitempty"`
}

// ServeSearch handles the search API endpoint.
func (h *Handler) ServeSearch(w http.ResponseWriter, r *http.Request) {
	h.handleSearch(w, r)
}

// handleSearch handles the search API endpoint.
func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		// Return empty results for empty query
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{
			Results: []SearchResultItem{},
			Query:   "",
		})
		return
	}

	h.mu.RLock()
	site := h.cfg.Site
	hideDraft := h.cfg.HideDraft
	basePath := h.cfg.BasePath
	h.mu.RUnlock()

	if site == nil {
		http.Error(w, "server not configured", http.StatusInternalServerError)
		return
	}

	// Perform search
	searchResults := site.SearchDocs(query, !hideDraft, 100) // Max 100 results

	// Convert to API response format
	results := make([]SearchResultItem, 0, len(searchResults))
	for _, sr := range searchResults {
		// Build full path with base path
		path := sr.Path
		if basePath != "" {
			path = basePath + path
		}

		results = append(results, SearchResultItem{
			Key:         sr.Key,
			Title:       sr.Title,
			Path:        path,
			Snippet:     sr.Snippet,
			Score:       sr.Score,
			HeadingID:   sr.HeadingID,
			HeadingText: sr.HeadingText,
		})
	}

	response := SearchResponse{
		Results: results,
		Query:   query,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ServeSitemap handles the sitemap.xml endpoint.
func (h *Handler) ServeSitemap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.mu.RLock()
	site := h.cfg.Site
	hideDraft := h.cfg.HideDraft
	basePath := h.cfg.BasePath
	h.mu.RUnlock()

	if site == nil {
		http.Error(w, "server not configured", http.StatusInternalServerError)
		return
	}

	// Get all non-draft documents
	docs := site.ListDocs(!hideDraft)

	// Build base URL from request
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	host := r.Host
	if host == "" {
		host = "localhost"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	if basePath != "" {
		baseURL += basePath
	}

	// Generate sitemap XML
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Write XML header and root element
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	w.Write([]byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"))

	// Add each document to the sitemap
	for _, doc := range docs {
		// Build URL path
		urlPath := baseURL
		if doc.Key == "" {
			// Root index page
			if !strings.HasSuffix(urlPath, "/") {
				urlPath += "/"
			}
		} else {
			// Ensure single slash between basePath and key
			if !strings.HasSuffix(urlPath, "/") {
				urlPath += "/"
			}
			// URL-encode each path segment
			parts := strings.Split(doc.Key, "/")
			encodedParts := make([]string, len(parts))
			for i, part := range parts {
				encodedParts[i] = url.PathEscape(part)
			}
			urlPath += strings.Join(encodedParts, "/")
		}

		// Format last modified date (W3C datetime format)
		lastmod := doc.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z")

		// Determine priority based on document type
		priority := "0.7" // Default priority
		if doc.Key == "" {
			priority = "1.0" // Homepage gets highest priority
		} else if !strings.Contains(doc.Key, "/") {
			priority = "0.9" // Top-level pages
		}

		// Determine changefreq (how often the page is likely to change)
		changefreq := "monthly" // Default
		if doc.Key == "" {
			changefreq = "weekly" // Homepage changes more frequently
		}

		// Write URL entry
		w.Write([]byte("  <url>\n"))
		fmt.Fprintf(w, "    <loc>%s</loc>\n", escapeXML(urlPath))
		fmt.Fprintf(w, "    <lastmod>%s</lastmod>\n", lastmod)
		fmt.Fprintf(w, "    <changefreq>%s</changefreq>\n", changefreq)
		fmt.Fprintf(w, "    <priority>%s</priority>\n", priority)
		w.Write([]byte("  </url>\n"))
	}

	// Close root element
	w.Write([]byte("</urlset>\n"))
}

// escapeXML escapes special XML characters in a string.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
