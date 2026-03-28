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
			Key:      n.Key,
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

// buildSectionTabs builds SectionTab data from config sections using the given site.
// It returns the tabs and the index of the active section based on currentPath.
func (h *Handler) buildSectionTabs(targetSite *site.Site, currentPath string) ([]SectionTab, int) {
	h.mu.RLock()
	siteConfig := h.cfg.SiteConfig
	basePath := h.cfg.BasePath
	hideDraft := h.cfg.HideDraft
	h.mu.RUnlock()

	if siteConfig == nil || len(siteConfig.Nav.Sections) == 0 {
		return nil, -1
	}

	currentLang := ""
	if targetSite != nil && targetSite.Language != "" {
		currentLang = targetSite.Language
	}
	currentVersion := ""
	if targetSite != nil && targetSite.Version != "" {
		currentVersion = targetSite.Version
	} else if siteConfig.IsMultiVersion() {
		currentVersion = siteConfig.GetDefaultVersion()
	}

	tabs := make([]SectionTab, 0, len(siteConfig.Nav.Sections))
	activeIdx := -1

	for i, section := range siteConfig.Nav.Sections {
		// Build a nav tree from this section's items using the site's index
		navNodes := targetSite.NavTreeFromItems(section.Items, !hideDraft)
		items := convertNavNodesWithVersionAndLang(navNodes, basePath, currentVersion, currentLang, siteConfig)

		// Determine the first linkable path for this section tab
		tabPath := firstNavItemPath(items)

		// Check if current path belongs to this section
		isActive := navItemsContainPath(items, currentPath)
		if isActive && activeIdx == -1 {
			activeIdx = i
		}

		tabs = append(tabs, SectionTab{
			Title:    section.Title,
			Path:     tabPath,
			IsActive: isActive,
			Items:    items,
		})
	}

	// If no section matched, default to first section
	if activeIdx == -1 && len(tabs) > 0 {
		activeIdx = 0
		tabs[0].IsActive = true
	}

	return tabs, activeIdx
}

// firstNavItemPath returns the Path of the first linkable item in the nav tree.
func firstNavItemPath(items []NavItem) string {
	for _, item := range items {
		if item.Path != "" {
			return item.Path
		}
		if p := firstNavItemPath(item.Children); p != "" {
			return p
		}
	}
	return ""
}

// navItemsContainPath checks if any nav item (recursively) matches the given path.
func navItemsContainPath(items []NavItem, path string) bool {
	for _, item := range items {
		if item.Path != "" && item.Path == path {
			return true
		}
		if navItemsContainPath(item.Children, path) {
			return true
		}
	}
	return false
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
