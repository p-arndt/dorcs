// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
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
	SiteTitle string
	BasePath  string
	Cache     bool
	HideDraft bool

	Site         *site.Site
	DocumentTmpl *template.Template

	// SiteConfig holds the loaded dorcs config for theming/branding
	SiteConfig *config.Config

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

	// Root -> docs/index.md
	if reqPath == "/" {
		h.handleDocByKey(w, r, "")
		return
	}

	// Everything else: strip leading "/" and try to resolve to doc.
	rel := strings.TrimPrefix(reqPath, "/")
	rel = strings.Trim(rel, "/")
	if rel == "" {
		http.NotFound(w, r)
		return
	}

	// Support folder index: /guide -> docs/guide/index.md
	if strings.HasSuffix(reqPath, "/") {
		h.handleDocByKey(w, r, path.Clean(rel))
		return
	}

	// Check if this is a static asset (image, etc.) before trying markdown
	if h.tryServeStaticAsset(w, r, rel) {
		return
	}

	// First: try docs/<rel>.md
	if h.tryServeDoc(w, r, rel) {
		return
	}
	// Second: try docs/<rel>/index.md (folder index)
	if h.tryServeDoc(w, r, path.Clean(rel)) {
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

// tryServeStaticAsset attempts to serve a static asset from the docs folder.
// Returns true if the file was found and served, false otherwise.
func (h *Handler) tryServeStaticAsset(w http.ResponseWriter, _ *http.Request, relPath string) bool {
	// Only serve files that look like static assets
	if !isStaticAsset(relPath) {
		return false
	}

	h.mu.RLock()
	docsDir := h.cfg.DocsDir
	h.mu.RUnlock()

	// Build the full file path
	filePath := filepath.Join(docsDir, filepath.FromSlash(relPath))

	// Security: ensure the file is within the docs directory
	absDocsDir, err := filepath.Abs(docsDir)
	if err != nil {
		return false
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return false
	}
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

	// Copy file content to response
	_, err = io.Copy(w, file)
	return err == nil
}

func (h *Handler) tryServeDoc(w http.ResponseWriter, r *http.Request, key string) bool {
	key = cleanKey(key)
	h.mu.RLock()
	site := h.cfg.Site
	h.mu.RUnlock()
	if _, ok := site.GetDoc(key); ok {
		h.handleDocByKey(w, r, key)
		return true
	}
	return false
}

func (h *Handler) handleDocByKey(w http.ResponseWriter, r *http.Request, key string) {
	h.mu.RLock()
	site := h.cfg.Site
	documentTmpl := h.cfg.DocumentTmpl
	basePath := h.cfg.BasePath
	h.mu.RUnlock()

	if site == nil || documentTmpl == nil {
		http.Error(w, "server not configured", http.StatusInternalServerError)
		return
	}

	key = cleanKey(key)

	// Build navigation tree
	nav := h.buildNavItems()

	// Get root title from navigation tree
	rootTitle := h.getRootTitle()

	// Use site renderer
	rendered, err := site.RenderDoc(key)
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
	m.CurrentPath = r.URL.Path
	m.BasePath = basePath

	m.Meta.Title = rendered.Doc.Title
	m.Meta.Description = rendered.Doc.Description
	m.Meta.Date = formatDate(rendered.Doc.Date)
	m.Meta.Tags = append([]string(nil), rendered.Doc.Tags...)
	m.Meta.SourcePath = rendered.Doc.RelPath

	// Add config and theme CSS if available (read with lock)
	h.mu.RLock()
	siteConfig := h.cfg.SiteConfig
	h.mu.RUnlock()

	if siteConfig != nil {
		m.Config = siteConfig
		m.ThemeCSS = template.CSS(siteConfig.GenerateThemeCSS())
		// Override site title from config if set
		if siteConfig.Site.Title != "" {
			m.SiteTitle = siteConfig.Site.Title
		}
		// Otherwise, m.SiteTitle is already set from the earlier read
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

// buildNavItems converts the site's NavTree to the template-friendly NavItem slice.
func (h *Handler) buildNavItems() []NavItem {
	h.mu.RLock()
	hideDraft := h.cfg.HideDraft
	site := h.cfg.Site
	h.mu.RUnlock()
	tree := site.NavTree(!hideDraft)
	if tree == nil {
		return nil
	}
	return convertNavNodes(tree.Children)
}

// getRootTitle extracts the title from the root index.md page.
func (h *Handler) getRootTitle() string {
	h.mu.RLock()
	hideDraft := h.cfg.HideDraft
	site := h.cfg.Site
	h.mu.RUnlock()
	tree := site.NavTree(!hideDraft)
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

func convertNavNodes(nodes []*site.NavNode) []NavItem {
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
			Children: convertNavNodes(n.Children),
		}

		// Set path - folders are only clickable if they have a landing page
		if n.IsDir {
			if n.Page != nil {
				item.Path = "/" + n.Key
			}
			// Path stays empty for folders without index.md
		} else {
			item.Path = "/" + n.Key
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

func firstNonEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
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
