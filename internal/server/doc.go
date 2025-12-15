// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"html/template"
	"net/http"
	"time"

	"dorcs-v2/internal/site"
)

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
