package site

import (
	"html"
	"html/template"
	"net/url"
	"strings"

	"dorcs-v2/internal/markdown"
)

// generateAndProcessTOC generates a table of contents and processes [[TOC]] placeholders.
// Returns the TOC HTML and the processed markdown (with [[TOC]] replaced if present).
func (s *Site) generateAndProcessTOC(raw string, doc *Doc, key string) (template.HTML, string) {
	hasTOCPlaceholder := strings.Contains(raw, "[[TOC]]")

	var toc template.HTML
	if hasTOCPlaceholder {
		// Generate TOC based on page type
		if isIndexRel(doc.RelPath) {
			// For index pages, generate navigation-based TOC
			toc = s.BuildNavTOC(key, s.BasePath)
		} else {
			// For regular pages, generate heading-based TOC
			toc = markdown.BuildTOC(s.md, raw)
		}
		// Replace [[TOC]] placeholder with the generated TOC HTML
		raw = markdown.ProcessTOCPlaceholder(raw, toc)
	} else {
		// Generate a table of contents from headings (h2/h3 by default) for sidebar
		toc = markdown.BuildTOC(s.md, raw)
	}

	return toc, raw
}

// BuildNavTOC generates a table of contents from navigation children for index pages.
// It finds the navigation node matching the given key and builds a TOC from its children.
func (s *Site) BuildNavTOC(key string, basePath string) template.HTML {
	s.mu.RLock()
	nav := s.nav
	s.mu.RUnlock()

	// Find the node matching the key
	var targetNode *NavNode
	var findNode func(node *NavNode) bool
	findNode = func(node *NavNode) bool {
		if node.Key == key {
			targetNode = node
			return true
		}
		for _, child := range node.Children {
			if findNode(child) {
				return true
			}
		}
		return false
	}

	if !findNode(nav) || targetNode == nil {
		return ""
	}

	// Build TOC from children
	if len(targetNode.Children) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<ul class="toc-list toc-nav">`)

	var buildNavTOCItems func(nodes []*NavNode, level int)
	buildNavTOCItems = func(nodes []*NavNode, level int) {
		for _, node := range nodes {
			// Skip draft pages
			if node.Page != nil && node.Page.Draft {
				continue
			}

			// Build URL path (matching RewriteExtensionlessDocLinks pattern)
			var urlPath string
			if node.Key == "" {
				// Root index
				urlPath = "/"
				if basePath != "" {
					urlPath = basePath + "/"
				}
			} else {
				escapedKey := escapePathForTOC(node.Key)
				urlPath = "/" + escapedKey
				if basePath != "" {
					urlPath = basePath + urlPath
				}
			}

			// Get display name
			name := node.Name
			if node.Page != nil && strings.TrimSpace(node.Page.Title) != "" {
				name = node.Page.Title
			}
			if name == "" {
				name = node.Key
			}

			// Write list item
			b.WriteString(`<li class="toc-item toc-nav-item">`)
			b.WriteString(`<a href="`)
			b.WriteString(html.EscapeString(urlPath))
			b.WriteString(`">`)
			b.WriteString(html.EscapeString(name))
			b.WriteString(`</a>`)

			// If it's a directory with children, nest them
			if node.IsDir && len(node.Children) > 0 {
				b.WriteString(`<ul class="toc-nested">`)
				buildNavTOCItems(node.Children, level+1)
				b.WriteString(`</ul>`)
			}

			b.WriteString(`</li>`)
		}
	}

	buildNavTOCItems(targetNode.Children, 1)
	b.WriteString(`</ul>`)

	return template.HTML(b.String())
}

// escapePathForTOC escapes path segments for use in URLs.
func escapePathForTOC(path string) string {
	// Replace backslashes and ensure forward slashes
	path = strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(path, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}
