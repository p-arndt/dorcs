// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"strings"

	"dorcs-v2/internal/site"
)

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
		var pathBuilder strings.Builder
		if basePath != "" {
			pathBuilder.WriteString(basePath)
		}
		if currentLang != "" {
			pathBuilder.WriteByte('/')
			pathBuilder.WriteString(currentLang)
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
