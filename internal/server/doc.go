// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/github"
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
	hideDraft := h.cfg.HideDraft
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
	doc, ok := targetSite.GetDoc(key)
	if !ok || (hideDraft && doc.Draft) {
		http.NotFound(w, r)
		return
	}

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
	m.ActiveSection = -1
	m.SiteTitle = siteTitle
	m.DocTitleFallback = key
	if key == "" {
		m.DocTitleFallback = "index"
	}
	m.HTML = rendered.HTML
	m.TOCHTML = rendered.TocHTML
	m.TOC = rendered.TocHTML
	m.Slides = rendered.Slides
	m.PresentationHeader = rendered.Doc.PresentationHeader
	m.PresentationFooter = rendered.Doc.PresentationFooter
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

	// Compute "Edit on GitHub" URL when docs are hosted on GitHub
	if siteConfig != nil {
		m.EditOnGitHubURL = computeEditOnGitHubURL(rendered.Doc, siteConfig, currentLang, currentVersion)
	}

	// Add config and theme CSS if available
	if siteConfig != nil {
		m.Config = siteConfig
		m.ThemeCSS = template.CSS(siteConfig.GenerateThemeCSS())
		// Override site title from config if set
		if siteConfig.Site.Title != "" {
			m.SiteTitle = siteConfig.Site.Title
		}
	}

	// Build section tabs if configured
	if siteConfig != nil && len(siteConfig.Nav.Sections) > 0 {
		sections, activeIdx := h.buildSectionTabs(targetSite, m.CurrentPath)
		m.Sections = sections
		m.ActiveSection = activeIdx
		// Override nav items with the active section's items
		if activeIdx >= 0 && activeIdx < len(sections) {
			m.Nav.Nodes = sections[activeIdx].Items
		}
	}

	// Enable live reload if broadcaster is configured
	m.LiveReload = reloadBroadcaster != nil
	m.Version = version

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Use presentation template for slide decks
	if rendered.Doc.Presentation && len(rendered.Slides) > 0 {
		if err := documentTmpl.ExecuteTemplate(w, "presentation", m); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Try "doc" template first (which calls "layout"), fall back to "layout" for compatibility
	if err := documentTmpl.ExecuteTemplate(w, "doc", m); err != nil {
		if err2 := h.cfg.DocumentTmpl.ExecuteTemplate(w, "layout", m); err2 != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// computeEditOnGitHubURL returns the GitHub "edit this page" URL when applicable.
// Supports: (1) GitHub content source (doc.IsGitHub), (2) edit_on_github for local docs.
func computeEditOnGitHubURL(doc *site.Doc, cfg *config.Config, currentLang, currentVersion string) string {
	if cfg == nil {
		return ""
	}
	gh := cfg.GitHub
	// Case 1: Docs sourced from GitHub - use exact GitHubPath
	if doc.IsGitHub && gh.Enabled && gh.Repository != "" && doc.GitHubPath != "" {
		repoInfo, err := github.ParseRepositoryURL(gh.Repository)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("https://github.com/%s/%s/edit/%s/%s",
			repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, doc.GitHubPath)
	}
	// Case 2: Local docs with edit_on_github - derive path from RelPath and repo structure
	if gh.EditOnGitHub.Repository != "" && doc.RelPath != "" {
		repoInfo, err := github.ParseRepositoryURL(gh.EditOnGitHub.Repository)
		if err != nil {
			return ""
		}
		path := repoInfo.Path
		if path == "" {
			path = "docs"
		}
		if currentLang != "" {
			path = path + "/" + currentLang
		}
		defaultVersion := ""
		if cfg.IsMultiVersion() {
			defaultVersion = cfg.GetDefaultVersion()
		}
		if currentVersion != "" && currentVersion != defaultVersion {
			path = path + "/" + currentVersion
		}
		filePath := strings.TrimPrefix(path+"/"+doc.RelPath, "/")
		return fmt.Sprintf("https://github.com/%s/%s/edit/%s/%s",
			repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, filePath)
	}
	return ""
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
