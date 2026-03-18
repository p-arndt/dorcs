// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"strings"

	"github.com/p-arndt/dorcs/internal/site"
)

// buildNavItemsWithSite converts the given site's NavTree to the template-friendly NavItem slice.
func (h *Handler) buildNavItemsWithSite(targetSite *site.Site) []NavItem {
	if targetSite == nil {
		return nil
	}
	h.mu.RLock()
	hideDraft := h.cfg.HideDraft
	basePath := h.cfg.BasePath
	siteConfig := h.cfg.SiteConfig
	currentLang := ""
	if targetSite.Language != "" {
		currentLang = targetSite.Language
	}
	currentVersion := ""
	if targetSite.Version != "" {
		currentVersion = targetSite.Version
	} else if siteConfig != nil && siteConfig.IsMultiVersion() {
		// If no version set but versioning is enabled, use default version
		currentVersion = siteConfig.GetDefaultVersion()
	}
	h.mu.RUnlock()
	tree := targetSite.NavTree(!hideDraft)
	if tree == nil {
		return nil
	}
	return convertNavNodesWithVersionAndLang(tree.Children, basePath, currentVersion, currentLang, siteConfig)
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

func convertNavNodesWithVersionAndLang(nodes []*site.NavNode, basePath string, currentVersion string, currentLang string, siteConfig interface {
	IsMultiVersion() bool
	GetDefaultVersion() string
	IsMultiLingual() bool
	GetDefaultLanguage() string
}) []NavItem {
	if len(nodes) == 0 {
		return nil
	}
	items := make([]NavItem, 0, len(nodes))
	for _, n := range nodes {
		title := n.Name
		if !n.ExplicitTitle && n.Page != nil && n.Page.Title != "" {
			title = n.Page.Title
		}
		if strings.TrimSpace(title) == "" && n.Page != nil && n.Page.Title != "" {
			title = n.Page.Title
		}

		item := NavItem{
			Title:    title,
			IsDir:    n.IsDir,
			Children: convertNavNodesWithVersionAndLang(n.Children, basePath, currentVersion, currentLang, siteConfig),
		}

		// Set path - folders are only clickable if they have a landing page
		var pathBuilder strings.Builder
		if basePath != "" {
			pathBuilder.WriteString(basePath)
		}
		// Add language prefix if not default language
		if currentLang != "" && siteConfig != nil && siteConfig.IsMultiLingual() {
			defaultLang := siteConfig.GetDefaultLanguage()
			if currentLang != defaultLang {
				pathBuilder.WriteByte('/')
				pathBuilder.WriteString(currentLang)
			}
		}
		// Add version prefix if not default version.
		// URLs are language-first: /de/v1/page
		if currentVersion != "" && siteConfig != nil && siteConfig.IsMultiVersion() {
			defaultVersion := siteConfig.GetDefaultVersion()
			if currentVersion != defaultVersion {
				pathBuilder.WriteByte('/')
				pathBuilder.WriteString(currentVersion)
			}
		}
		if n.IsDir {
			if n.Page != nil {
				pathBuilder.WriteByte('/')
				pathBuilder.WriteString(n.Key)
				item.Path = pathBuilder.String()
			}
			// Path stays empty for folders without index.md
		} else {
			pathBuilder.WriteByte('/')
			pathBuilder.WriteString(n.Key)
			item.Path = pathBuilder.String()
		}

		items = append(items, item)
	}
	return items
}

// convertNavNodesWithLang is a helper function for backward compatibility in tests.
// It wraps convertNavNodesWithVersionAndLang with empty version.
// If siteConfig is nil, language prefixes won't be added (assumes single-language mode).
func convertNavNodesWithLang(nodes []*site.NavNode, basePath string, currentLang string) []NavItem {
	// Create a minimal siteConfig that treats any non-empty language as multi-lingual
	var siteConfig interface {
		IsMultiVersion() bool
		GetDefaultVersion() string
		IsMultiLingual() bool
		GetDefaultLanguage() string
	}
	if currentLang != "" {
		// If currentLang is set, assume multi-lingual mode for test compatibility
		siteConfig = &mockSiteConfig{isMultiLingual: true, defaultLang: "en"}
	}
	return convertNavNodesWithVersionAndLang(nodes, basePath, "", currentLang, siteConfig)
}

// mockSiteConfig is a minimal implementation for tests
type mockSiteConfig struct {
	isMultiLingual bool
	defaultLang    string
}

func (m *mockSiteConfig) IsMultiVersion() bool       { return false }
func (m *mockSiteConfig) GetDefaultVersion() string  { return "" }
func (m *mockSiteConfig) IsMultiLingual() bool       { return m.isMultiLingual }
func (m *mockSiteConfig) GetDefaultLanguage() string { return m.defaultLang }
