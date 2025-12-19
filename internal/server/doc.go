// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/p-arndt/dorcs/internal/site"
)

func (h *Handler) tryServeDocWithSite(w http.ResponseWriter, r *http.Request, key string, targetSite *site.Site, currentLang string, currentVersion string) bool {
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
		h.handleDocByKeyWithSite(w, r, key, targetSite, currentLang, currentVersion, docPathForCall)
		return true
	}
	return false
}

func (h *Handler) handleDocByKeyWithSite(w http.ResponseWriter, r *http.Request, key string, targetSite *site.Site, currentLang string, currentVersion string, docPath string) {
	h.mu.RLock()
	documentTmpl := h.cfg.DocumentTmpl
	basePath := h.cfg.BasePath
	siteConfig := h.cfg.SiteConfig
	siteTitle := h.cfg.SiteTitle
	reloadBroadcaster := h.cfg.ReloadBroadcaster
	version := h.cfg.Version
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
	// Set current path - use the actual request path for consistency
	// This ensures CurrentPath matches the URL structure (language-first)
	m.CurrentPath = r.URL.Path
	// DocPath is the document path without version/language prefix (used for building version/language links)
	m.DocPath = docPath
	m.BasePath = basePath
	m.CurrentLanguage = currentLang
	m.CurrentVersion = currentVersion

	// Security headers for public-internet deployments.
	setCommonSecurityHeaders(w)
	setCSPHeaders(w, r)

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

func setCommonSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
	// If you need framing, remove this and adjust CSP frame-ancestors.
	w.Header().Set("X-Frame-Options", "DENY")
}

func setCSPHeaders(w http.ResponseWriter, r *http.Request) {
	// Determine if the request is effectively HTTPS (direct or behind a proxy).
	secure := r.TLS != nil
	if !secure {
		if xfproto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); strings.EqualFold(xfproto, "https") {
			secure = true
		}
	}

	// Allow our own scripts, required CDN scripts (KaTeX / mermaid), and inline scripts.
	csp := "default-src 'self'; " +
		"base-uri 'self'; " +
		"object-src 'none'; " +
		"frame-ancestors 'none'; " +
		"img-src 'self' data: https:; " +
		"font-src 'self' data:; " +
		"connect-src 'self'; " +
		"style-src 'self' 'unsafe-inline'; " +
		"script-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'"
	if secure {
		csp += "; upgrade-insecure-requests"
	}
	w.Header().Set("Content-Security-Policy", csp)
}
