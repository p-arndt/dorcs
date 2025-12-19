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

	// VersionSites is a map of version+language key to Site instance (for versioning support)
	// Key format: "{version}:{language}" or "{version}:" for default language, or ":{language}" for default version
	// If nil or empty, only the default Sites are used
	VersionSites map[string]*site.Site

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

	// CurrentVersion is the version identifier for the current page
	CurrentVersion string
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

	// Detect version and language from URL path (MkDocs-style: language-first)
	// Priority: language first, then version
	// Examples: /en/..., /en/v1/..., /v1/... (version-only, no language), /... (default)
	var currentVersion string
	var currentLang string
	var docPath string
	if reqPath == "/" {
		docPath = "/"
	} else {
		// Extract path segments
		parts := strings.Split(strings.TrimPrefix(reqPath, "/"), "/")
		if len(parts) > 0 && parts[0] != "" {
			// Check if first segment is a language code (language-first)
			if siteConfig != nil && siteConfig.IsLanguageEnabled(parts[0]) {
				currentLang = parts[0]
				// Check if second segment is a version ID
				if len(parts) > 1 && parts[1] != "" {
					if siteConfig.IsVersionEnabled(parts[1]) {
						currentVersion = parts[1]
						// Remove language and version prefix from path
						if len(parts) > 2 {
							docPath = "/" + strings.Join(parts[2:], "/")
						} else {
							docPath = "/"
						}
					} else {
						// Language but no version, remove only language prefix
						if len(parts) > 1 {
							docPath = "/" + strings.Join(parts[1:], "/")
						} else {
							docPath = "/"
						}
					}
				} else {
					// Language only, no version
					docPath = "/"
				}
			} else if siteConfig != nil && siteConfig.IsVersionEnabled(parts[0]) {
				// Not a language, check if it's a version ID (version-only mode, no languages)
				currentVersion = parts[0]
				// Remove version prefix from path
				if len(parts) > 1 {
					docPath = "/" + strings.Join(parts[1:], "/")
				} else {
					docPath = "/"
				}
			} else {
				// Not a language or version code, treat as document path
				docPath = reqPath
			}
		} else {
			docPath = reqPath
		}
	}

	// Get the appropriate Site instance for this version and language
	var targetSite *site.Site
	h.mu.RLock()
	// If versioning is enabled, we need to handle default version
	if siteConfig != nil && siteConfig.IsMultiVersion() {
		// If no version specified in URL, use default version
		if currentVersion == "" {
			currentVersion = siteConfig.GetDefaultVersion()
		}
		// First, try version-specific sites
		if len(h.cfg.VersionSites) > 0 {
			// Build key: "{version}:{language}" or "{version}:" for default language
			siteKey := currentVersion + ":"
			if currentLang != "" {
				siteKey = currentVersion + ":" + currentLang
			}
			if verSite, ok := h.cfg.VersionSites[siteKey]; ok {
				targetSite = verSite
			}
		}
	}
	// If no version-specific site found, try language-only sites (for backward compatibility)
	if targetSite == nil && len(h.cfg.Sites) > 0 {
		// Default language uses empty key in Sites map
		siteKey := currentLang
		if currentLang == "" {
			siteKey = "" // Default language
		}
		if langSite, ok := h.cfg.Sites[siteKey]; ok {
			targetSite = langSite
		}
	}
	// Fall back to default site if no version/language-specific site found
	if targetSite == nil {
		targetSite = h.cfg.Site
		// Reset to defaults if we couldn't find a specific site
		if siteConfig != nil && siteConfig.IsMultiVersion() {
			currentVersion = siteConfig.GetDefaultVersion()
		} else {
			currentVersion = ""
		}
		if siteConfig != nil && siteConfig.IsMultiLingual() {
			currentLang = siteConfig.GetDefaultLanguage()
		} else {
			currentLang = ""
		}
	}
	h.mu.RUnlock()

	// Root -> docs/index.md
	if docPath == "/" {
		h.handleDocByKeyWithSite(w, r, "", targetSite, currentLang, currentVersion, docPath)
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
		h.handleDocByKeyWithSite(w, r, path.Clean(rel), targetSite, currentLang, currentVersion, docPath)
		return
	}

	// Check if this is a static asset (image, etc.) before trying markdown
	if h.tryServeStaticAsset(w, r, rel, targetSite) {
		return
	}

	// First: try docs/<rel>.md
	if h.tryServeDocWithSite(w, r, rel, targetSite, currentLang, currentVersion) {
		return
	}
	// Second: try docs/<rel>/index.md (folder index)
	if h.tryServeDocWithSite(w, r, path.Clean(rel), targetSite, currentLang, currentVersion) {
		return
	}

	http.NotFound(w, r)
}
