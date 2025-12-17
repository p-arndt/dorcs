// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"html/template"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/site"
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
	if len(h.cfg.Sites) > 0 {
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
