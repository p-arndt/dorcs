package site

import (
	"html"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"dorcs-v2/internal/markdown"
)

// generateAndProcessTOC generates a table of contents and processes all placeholders.
// Returns the TOC HTML (for sidebar), root TOC HTML, and the processed markdown.
func (s *Site) generateAndProcessTOC(raw string, doc *Doc, key string) (template.HTML, template.HTML, string) {
	// Create placeholder context
	ctx := PlaceholderContext{
		CurrentDoc: doc,
		CurrentKey: key,
		BasePath:   s.BasePath,
		Site:       s,
	}

	// Build all placeholders
	placeholders := make(markdown.PlaceholderMap)

	// Check for TOC placeholders
	hasTOCPlaceholder := strings.Contains(raw, "[[TOC]]")
	hasRootTOCPlaceholder := strings.Contains(raw, "[[TOC-ROOT]]")
	hasTOCFiltered := strings.Contains(raw, "[[TOC:")

	var toc template.HTML
	var rootTOC template.HTML

	// Generate root navigation TOC if needed
	if hasRootTOCPlaceholder {
		rootTOC = s.BuildRootNavTOC(s.BasePath)
		placeholders["TOC-ROOT"] = rootTOC
	}

	// Generate page TOC if needed
	if hasTOCPlaceholder {
		// [[TOC]] always generates heading-based TOC for the current page
		toc = markdown.BuildTOC(s.md, raw)
		placeholders["TOC"] = toc
	}

	// Handle filtered TOC (e.g., [[TOC:2]] for depth 2)
	if hasTOCFiltered {
		// Extract depth from placeholders like [[TOC:2]]
		lines := strings.Split(raw, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "[[TOC:") && strings.HasSuffix(trimmed, "]]") {
				depthStr := strings.TrimPrefix(strings.TrimSuffix(trimmed, "]]"), "[[TOC:")
				if depth, err := strconv.Atoi(depthStr); err == nil && depth > 0 {
					// Build filtered TOC with specified depth
					filteredTOC := markdown.BuildTOCWithDepth(s.md, raw, depth)
					placeholders["TOC:"+depthStr] = filteredTOC
				}
			}
		}
	}

	// Generate all other placeholders
	placeholders["BREADCRUMBS"] = BuildBreadcrumbs(ctx)
	placeholders["CHILDREN"] = BuildChildren(ctx)
	placeholders["SIBLINGS"] = BuildSiblings(ctx)
	placeholders["RELATED"] = BuildRelated(ctx)
	placeholders["RECENT"] = BuildRecent(ctx)
	placeholders["TAGS"] = BuildTags(ctx)
	placeholders["INDEX"] = BuildIndex(ctx)
	placeholders["DATE"] = BuildDate(ctx)
	placeholders["PUBLISHED"] = BuildDate(ctx) // Alias
	placeholders["LAST-UPDATED"] = BuildLastUpdated(ctx)
	placeholders["AUTHOR"] = BuildAuthor(ctx)
	placeholders["SUMMARY"] = BuildSummary(ctx)
	placeholders["PAGES-BY-TAG"] = BuildPagesByTag(ctx)

	// Process all placeholders
	raw = markdown.ProcessPlaceholders(raw, placeholders)

	// If no TOC placeholders were found, generate default TOC for sidebar
	if !hasTOCPlaceholder && !hasRootTOCPlaceholder {
		// Generate a table of contents from headings (h2/h3 by default) for sidebar
		toc = markdown.BuildTOC(s.md, raw)
	}

	return toc, rootTOC, raw
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

// BuildRootNavTOC generates a table of contents from the root navigation node's children.
// This shows the main pages/sections at the root level.
func (s *Site) BuildRootNavTOC(basePath string) template.HTML {
	s.mu.RLock()
	nav := s.nav
	s.mu.RUnlock()

	if nav == nil {
		return ""
	}

	// Build TOC from root's children
	if len(nav.Children) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<ul class="toc-list toc-nav toc-root">`)

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

	buildNavTOCItems(nav.Children, 1)
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
