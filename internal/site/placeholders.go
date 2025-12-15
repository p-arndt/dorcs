package site

import (
	"fmt"
	"html"
	"html/template"
	"net/url"
	"sort"
	"strings"
	"time"
)

// PlaceholderContext holds all the data needed to generate placeholders
type PlaceholderContext struct {
	CurrentDoc *Doc
	CurrentKey string
	BasePath   string
	Site       *Site
}

// BuildBreadcrumbs generates breadcrumb navigation HTML with horizontal layout and separators
func BuildBreadcrumbs(ctx PlaceholderContext) template.HTML {
	if ctx.CurrentKey == "" {
		// Root page - no breadcrumbs needed
		return ""
	}

	parts := strings.Split(ctx.CurrentKey, "/")
	var b strings.Builder
	b.WriteString(`<nav class="breadcrumbs" aria-label="Breadcrumb">`)
	b.WriteString(`<ol class="breadcrumb-list" itemscope itemtype="https://schema.org/BreadcrumbList">`)

	// Home link
	homeURL := "/"
	if ctx.BasePath != "" {
		homeURL = ctx.BasePath + "/"
	}
	b.WriteString(`<li class="breadcrumb-item" itemprop="itemListElement" itemscope itemtype="https://schema.org/ListItem">`)
	b.WriteString(`<a href="`)
	b.WriteString(html.EscapeString(homeURL))
	b.WriteString(`" itemprop="item">`)
	b.WriteString(`<span itemprop="name">Home</span>`)
	b.WriteString(`</a>`)
	b.WriteString(`<meta itemprop="position" content="1" />`)
	b.WriteString(`</li>`)

	// Build path segments
	currentPath := ""
	position := 2
	for i, part := range parts {
		currentPath += part
		if i < len(parts)-1 {
			currentPath += "/"
		}

		// Try to get the doc to get its title
		doc, ok := ctx.Site.GetDoc(currentPath)
		displayName := part
		if ok && doc != nil && doc.Title != "" {
			displayName = doc.Title
		}

		// Build URL
		urlPath := "/" + escapePathForTOC(currentPath)
		if ctx.BasePath != "" {
			urlPath = ctx.BasePath + urlPath
		}

		// Add separator before each item (except first)
		b.WriteString(`<li class="breadcrumb-separator" aria-hidden="true">`)
		b.WriteString(`<svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">`)
		b.WriteString(`<path d="M6 12l4-4-4-4"/>`)
		b.WriteString(`</svg>`)
		b.WriteString(`</li>`)

		b.WriteString(`<li class="breadcrumb-item" itemprop="itemListElement" itemscope itemtype="https://schema.org/ListItem">`)
		if i < len(parts)-1 {
			b.WriteString(`<a href="`)
			b.WriteString(html.EscapeString(urlPath))
			b.WriteString(`" itemprop="item">`)
			b.WriteString(`<span itemprop="name">`)
			b.WriteString(html.EscapeString(displayName))
			b.WriteString(`</span>`)
			b.WriteString(`</a>`)
		} else {
			b.WriteString(`<span class="breadcrumb-current" aria-current="page" itemprop="name">`)
			b.WriteString(html.EscapeString(displayName))
			b.WriteString(`</span>`)
		}
		b.WriteString(fmt.Sprintf(`<meta itemprop="position" content="%d" />`, position))
		b.WriteString(`</li>`)
		position++
	}

	b.WriteString(`</ol>`)
	b.WriteString(`</nav>`)
	return template.HTML(b.String())
}

// BuildChildren generates a list of all children pages in the current folder/section
func BuildChildren(ctx PlaceholderContext) template.HTML {
	ctx.Site.mu.RLock()
	nav := ctx.Site.nav
	ctx.Site.mu.RUnlock()

	// Find the node matching the current key
	var targetNode *NavNode
	var findNode func(node *NavNode) bool
	findNode = func(node *NavNode) bool {
		if node.Key == ctx.CurrentKey {
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

	// Get all children (both pages and directories)
	children := targetNode.Children
	if len(children) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<nav class="children-nav">`)
	b.WriteString(`<h3 class="children-title">Pages in this section</h3>`)
	b.WriteString(`<ul class="children-list">`)

	// Sort children using the same logic as sortNav: directories first, then by order, numeric prefix, or title
	sortedChildren := make([]*NavNode, len(children))
	copy(sortedChildren, children)
	sort.Slice(sortedChildren, func(i, j int) bool {
		a, b := sortedChildren[i], sortedChildren[j]
		if a.IsDir != b.IsDir {
			return a.IsDir // directories first
		}

		// Get order values (0 means no order specified)
		aOrder := 0
		bOrder := 0
		if a.Page != nil {
			aOrder = a.Page.Order
		}
		if b.Page != nil {
			bOrder = b.Page.Order
		}

		// If both have order set (non-zero), sort by order
		if aOrder != 0 && bOrder != 0 {
			if aOrder != bOrder {
				return aOrder < bOrder
			}
		} else if aOrder != 0 {
			// a has order, b doesn't - a comes first
			return true
		} else if bOrder != 0 {
			// b has order, a doesn't - b comes first
			return false
		}
		// Both have no order (or both are 0), check filename prefixes

		// Extract numeric prefixes from filenames
		aPrefix := 0
		bPrefix := 0
		if a.Page != nil {
			aPrefix = extractNumericPrefix(a.Page.RelPath)
		}
		if b.Page != nil {
			bPrefix = extractNumericPrefix(b.Page.RelPath)
		}

		// If both have numeric prefixes, sort by prefix
		if aPrefix != 0 && bPrefix != 0 {
			if aPrefix != bPrefix {
				return aPrefix < bPrefix
			}
		} else if aPrefix != 0 {
			// a has prefix, b doesn't - a comes first
			return true
		} else if bPrefix != 0 {
			// b has prefix, a doesn't - b comes first
			return false
		}
		// Both have no prefix, fall back to title/key sorting

		// Prefer page titles for dirs if present.
		aName := a.Name
		if a.IsDir && a.Page != nil && strings.TrimSpace(a.Page.Title) != "" {
			aName = a.Page.Title
		}
		bName := b.Name
		if b.IsDir && b.Page != nil && strings.TrimSpace(b.Page.Title) != "" {
			bName = b.Page.Title
		}

		// For leaf pages prefer doc title.
		if !a.IsDir && a.Page != nil && strings.TrimSpace(a.Page.Title) != "" {
			aName = a.Page.Title
		}
		if !b.IsDir && b.Page != nil && strings.TrimSpace(b.Page.Title) != "" {
			bName = b.Page.Title
		}

		aName = strings.ToLower(strings.TrimSpace(aName))
		bName = strings.ToLower(strings.TrimSpace(bName))
		if aName != bName {
			return aName < bName
		}
		return a.Key < b.Key
	})

	for _, child := range sortedChildren {
		// Skip draft pages
		if child.Page != nil && child.Page.Draft {
			continue
		}

		// Build URL
		childURL := buildDocURL(child.Key, ctx.BasePath)

		// Get display name
		childName := child.Name
		if child.Page != nil && child.Page.Title != "" {
			childName = child.Page.Title
		}
		if childName == "" {
			childName = child.Key
		}

		// Determine if it's a directory or page
		itemClass := "children-item"
		if child.IsDir {
			itemClass += " children-item-dir"
		} else {
			itemClass += " children-item-page"
		}

		b.WriteString(`<li class="`)
		b.WriteString(itemClass)
		b.WriteString(`">`)

		// Add icon or indicator for directories
		if child.IsDir {
			b.WriteString(`<span class="children-icon" aria-hidden="true">📁</span>`)
		}

		b.WriteString(`<a href="`)
		b.WriteString(html.EscapeString(childURL))
		b.WriteString(`" class="children-link">`)
		b.WriteString(html.EscapeString(childName))
		b.WriteString(`</a>`)

		// Add description if available
		if child.Page != nil && child.Page.Description != "" {
			b.WriteString(`<span class="children-description"> - `)
			b.WriteString(html.EscapeString(child.Page.Description))
			b.WriteString(`</span>`)
		}

		// Add child count for directories
		if child.IsDir && len(child.Children) > 0 {
			nonDraftCount := 0
			for _, c := range child.Children {
				if c.Page == nil || !c.Page.Draft {
					nonDraftCount++
				}
			}
			if nonDraftCount > 0 {
				b.WriteString(`<span class="children-count">`)
				b.WriteString(fmt.Sprintf("(%d)", nonDraftCount))
				b.WriteString(`</span>`)
			}
		}

		b.WriteString(`</li>`)
	}

	b.WriteString(`</ul>`)
	b.WriteString(`</nav>`)
	return template.HTML(b.String())
}

// BuildSiblings generates sibling pages (pages at same level)
func BuildSiblings(ctx PlaceholderContext) template.HTML {
	ctx.Site.mu.RLock()
	nav := ctx.Site.nav
	ctx.Site.mu.RUnlock()

	// Find current node
	var currentNode *NavNode
	var findNode func(node *NavNode) bool
	findNode = func(node *NavNode) bool {
		if node.Key == ctx.CurrentKey {
			currentNode = node
			return true
		}
		for _, child := range node.Children {
			if findNode(child) {
				return true
			}
		}
		return false
	}

	if !findNode(nav) || currentNode == nil {
		return ""
	}

	// Find parent
	var parent *NavNode
	var findParent func(node *NavNode) *NavNode
	findParent = func(node *NavNode) *NavNode {
		for _, child := range node.Children {
			if child == currentNode {
				return node
			}
			if p := findParent(child); p != nil {
				return p
			}
		}
		return nil
	}

	parent = findParent(nav)
	if parent == nil || len(parent.Children) <= 1 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<nav class="siblings-nav">`)
	b.WriteString(`<h3 class="siblings-title">Related Pages</h3>`)
	b.WriteString(`<ul class="siblings-list">`)

	for _, sibling := range parent.Children {
		if sibling.Page != nil && sibling.Page.Draft {
			continue
		}
		if sibling.Key == ctx.CurrentKey {
			continue // Skip current page
		}

		siblingURL := buildDocURL(sibling.Key, ctx.BasePath)
		siblingTitle := sibling.Name
		if sibling.Page != nil && sibling.Page.Title != "" {
			siblingTitle = sibling.Page.Title
		}

		b.WriteString(`<li class="sibling-item">`)
		b.WriteString(`<a href="`)
		b.WriteString(html.EscapeString(siblingURL))
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(siblingTitle))
		b.WriteString(`</a>`)
		b.WriteString(`</li>`)
	}

	b.WriteString(`</ul>`)
	b.WriteString(`</nav>`)
	return template.HTML(b.String())
}

// BuildRelated generates related pages based on tags
func BuildRelated(ctx PlaceholderContext) template.HTML {
	if ctx.CurrentDoc == nil || len(ctx.CurrentDoc.Tags) == 0 {
		return ""
	}

	allDocs := ctx.Site.ListDocs(false)
	tagMap := make(map[string]int)
	relatedDocs := make(map[string]*Doc)

	// Count tag matches
	for _, doc := range allDocs {
		if doc.Key == ctx.CurrentKey {
			continue
		}
		if doc.Draft {
			continue
		}

		score := 0
		for _, tag := range ctx.CurrentDoc.Tags {
			for _, docTag := range doc.Tags {
				if tag == docTag {
					score++
					break
				}
			}
		}

		if score > 0 {
			tagMap[doc.Key] = score
			relatedDocs[doc.Key] = doc
		}
	}

	if len(relatedDocs) == 0 {
		return ""
	}

	// Sort by score
	type scoredDoc struct {
		doc   *Doc
		score int
	}
	scored := make([]scoredDoc, 0, len(relatedDocs))
	for key, doc := range relatedDocs {
		scored = append(scored, scoredDoc{doc: doc, score: tagMap[key]})
	}
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		return scored[i].doc.Title < scored[j].doc.Title
	})

	// Limit to top 5
	if len(scored) > 5 {
		scored = scored[:5]
	}

	var b strings.Builder
	b.WriteString(`<nav class="related-nav">`)
	b.WriteString(`<h3 class="related-title">Related Pages</h3>`)
	b.WriteString(`<ul class="related-list">`)

	for _, item := range scored {
		docURL := buildDocURL(item.doc.Key, ctx.BasePath)
		docTitle := item.doc.Title
		if docTitle == "" {
			docTitle = item.doc.Key
		}

		b.WriteString(`<li class="related-item">`)
		b.WriteString(`<a href="`)
		b.WriteString(html.EscapeString(docURL))
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(docTitle))
		b.WriteString(`</a>`)
		if item.doc.Description != "" {
			b.WriteString(`<span class="related-description"> - `)
			b.WriteString(html.EscapeString(item.doc.Description))
			b.WriteString(`</span>`)
		}
		b.WriteString(`</li>`)
	}

	b.WriteString(`</ul>`)
	b.WriteString(`</nav>`)
	return template.HTML(b.String())
}

// BuildRecent generates recently updated pages
func BuildRecent(ctx PlaceholderContext) template.HTML {
	allDocs := ctx.Site.ListDocs(false)
	if len(allDocs) == 0 {
		return ""
	}

	// Sort by UpdatedAt
	sorted := make([]*Doc, len(allDocs))
	copy(sorted, allDocs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].UpdatedAt.After(sorted[j].UpdatedAt)
	})

	// Limit to 10 most recent
	limit := 10
	if len(sorted) < limit {
		limit = len(sorted)
	}
	sorted = sorted[:limit]

	var b strings.Builder
	b.WriteString(`<nav class="recent-nav">`)
	b.WriteString(`<h3 class="recent-title">Recently Updated</h3>`)
	b.WriteString(`<ul class="recent-list">`)

	for _, doc := range sorted {
		docURL := buildDocURL(doc.Key, ctx.BasePath)
		docTitle := doc.Title
		if docTitle == "" {
			docTitle = doc.Key
		}

		b.WriteString(`<li class="recent-item">`)
		b.WriteString(`<a href="`)
		b.WriteString(html.EscapeString(docURL))
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(docTitle))
		b.WriteString(`</a>`)
		if !doc.UpdatedAt.IsZero() {
			b.WriteString(`<time class="recent-date" datetime="`)
			b.WriteString(html.EscapeString(doc.UpdatedAt.Format(time.RFC3339)))
			b.WriteString(`"> - `)
			b.WriteString(html.EscapeString(formatDate(doc.UpdatedAt)))
			b.WriteString(`</time>`)
		}
		b.WriteString(`</li>`)
	}

	b.WriteString(`</ul>`)
	b.WriteString(`</nav>`)
	return template.HTML(b.String())
}

// BuildTags generates tags display for current page
func BuildTags(ctx PlaceholderContext) template.HTML {
	if ctx.CurrentDoc == nil || len(ctx.CurrentDoc.Tags) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<div class="tags-container">`)
	b.WriteString(`<span class="tags-label">Tags:</span>`)
	b.WriteString(`<ul class="tags-list">`)

	for _, tag := range ctx.CurrentDoc.Tags {
		b.WriteString(`<li class="tag-item">`)
		b.WriteString(`<a href="`)
		// Link to tag page (could be enhanced to create tag index pages)
		tagURL := buildDocURL("", ctx.BasePath) + "?tag=" + url.QueryEscape(tag)
		b.WriteString(html.EscapeString(tagURL))
		b.WriteString(`" class="tag-link">`)
		b.WriteString(html.EscapeString(tag))
		b.WriteString(`</a>`)
		b.WriteString(`</li>`)
	}

	b.WriteString(`</ul>`)
	b.WriteString(`</div>`)
	return template.HTML(b.String())
}

// BuildIndex generates full site index
func BuildIndex(ctx PlaceholderContext) template.HTML {
	allDocs := ctx.Site.ListDocs(false)
	if len(allDocs) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<nav class="site-index">`)
	b.WriteString(`<h2 class="index-title">Site Index</h2>`)
	b.WriteString(`<ul class="index-list">`)

	var buildIndexItems func(docs []*Doc, level int)
	buildIndexItems = func(docs []*Doc, level int) {
		// Group by directory
		dirMap := make(map[string][]*Doc)
		for _, doc := range docs {
			dir := doc.DirKey
			if dir == "" {
				dir = "root"
			}
			dirMap[dir] = append(dirMap[dir], doc)
		}

		// Sort directories
		dirs := make([]string, 0, len(dirMap))
		for dir := range dirMap {
			dirs = append(dirs, dir)
		}
		sort.Strings(dirs)

		for _, dir := range dirs {
			dirDocs := dirMap[dir]
			sort.Slice(dirDocs, func(i, j int) bool {
				return dirDocs[i].Title < dirDocs[j].Title
			})

			for _, doc := range dirDocs {
				docURL := buildDocURL(doc.Key, ctx.BasePath)
				docTitle := doc.Title
				if docTitle == "" {
					docTitle = doc.Key
				}

				b.WriteString(`<li class="index-item">`)
				b.WriteString(`<a href="`)
				b.WriteString(html.EscapeString(docURL))
				b.WriteString(`">`)
				b.WriteString(html.EscapeString(docTitle))
				b.WriteString(`</a>`)
				if doc.Description != "" {
					b.WriteString(`<span class="index-description"> - `)
					b.WriteString(html.EscapeString(doc.Description))
					b.WriteString(`</span>`)
				}
				b.WriteString(`</li>`)
			}
		}
	}

	buildIndexItems(allDocs, 0)
	b.WriteString(`</ul>`)
	b.WriteString(`</nav>`)
	return template.HTML(b.String())
}

// BuildDate generates publication date display
func BuildDate(ctx PlaceholderContext) template.HTML {
	if ctx.CurrentDoc == nil || ctx.CurrentDoc.Date.IsZero() {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<time class="publish-date" datetime="`)
	b.WriteString(html.EscapeString(ctx.CurrentDoc.Date.Format(time.RFC3339)))
	b.WriteString(`">`)
	b.WriteString(`Published: `)
	b.WriteString(html.EscapeString(formatDate(ctx.CurrentDoc.Date)))
	b.WriteString(`</time>`)
	return template.HTML(b.String())
}

// BuildLastUpdated generates last updated date display
func BuildLastUpdated(ctx PlaceholderContext) template.HTML {
	if ctx.CurrentDoc == nil || ctx.CurrentDoc.UpdatedAt.IsZero() {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<time class="last-updated" datetime="`)
	b.WriteString(html.EscapeString(ctx.CurrentDoc.UpdatedAt.Format(time.RFC3339)))
	b.WriteString(`">`)
	b.WriteString(`Last updated: `)
	b.WriteString(html.EscapeString(formatDate(ctx.CurrentDoc.UpdatedAt)))
	b.WriteString(`</time>`)
	return template.HTML(b.String())
}

// BuildAuthor generates author information display
func BuildAuthor(ctx PlaceholderContext) template.HTML {
	if ctx.CurrentDoc == nil || ctx.CurrentDoc.Author == "" {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<div class="author-info">`)
	b.WriteString(`<span class="author-label">Author:</span>`)
	b.WriteString(`<span class="author-name">`)
	b.WriteString(html.EscapeString(ctx.CurrentDoc.Author))
	b.WriteString(`</span>`)
	b.WriteString(`</div>`)
	return template.HTML(b.String())
}

// BuildSummary generates page summary from description or first paragraph
func BuildSummary(ctx PlaceholderContext) template.HTML {
	if ctx.CurrentDoc == nil {
		return ""
	}

	var b strings.Builder
	if ctx.CurrentDoc.Description != "" {
		b.WriteString(`<div class="page-summary">`)
		b.WriteString(html.EscapeString(ctx.CurrentDoc.Description))
		b.WriteString(`</div>`)
		return template.HTML(b.String())
	}

	// Could extract first paragraph from markdown, but for now just use description
	return ""
}

// BuildPagesByTag generates pages grouped by tag
func BuildPagesByTag(ctx PlaceholderContext) template.HTML {
	allDocs := ctx.Site.ListDocs(false)
	tagMap := make(map[string][]*Doc)

	// Group docs by tags
	for _, doc := range allDocs {
		for _, tag := range doc.Tags {
			tagMap[tag] = append(tagMap[tag], doc)
		}
	}

	if len(tagMap) == 0 {
		return ""
	}

	// Sort tags
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	var b strings.Builder
	b.WriteString(`<nav class="pages-by-tag">`)
	b.WriteString(`<h2 class="pages-by-tag-title">Pages by Tag</h2>`)

	for _, tag := range tags {
		docs := tagMap[tag]
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].Title < docs[j].Title
		})

		b.WriteString(`<section class="tag-section">`)
		b.WriteString(`<h3 class="tag-section-title">`)
		b.WriteString(html.EscapeString(tag))
		b.WriteString(`</h3>`)
		b.WriteString(`<ul class="tag-section-list">`)

		for _, doc := range docs {
			docURL := buildDocURL(doc.Key, ctx.BasePath)
			docTitle := doc.Title
			if docTitle == "" {
				docTitle = doc.Key
			}

			b.WriteString(`<li class="tag-section-item">`)
			b.WriteString(`<a href="`)
			b.WriteString(html.EscapeString(docURL))
			b.WriteString(`">`)
			b.WriteString(html.EscapeString(docTitle))
			b.WriteString(`</a>`)
			b.WriteString(`</li>`)
		}

		b.WriteString(`</ul>`)
		b.WriteString(`</section>`)
	}

	b.WriteString(`</nav>`)
	return template.HTML(b.String())
}

// Helper functions

func buildDocURL(key string, basePath string) string {
	if key == "" {
		urlPath := "/"
		if basePath != "" {
			urlPath = basePath + "/"
		}
		return urlPath
	}
	urlPath := "/" + escapePathForTOC(key)
	if basePath != "" {
		urlPath = basePath + urlPath
	}
	return urlPath
}

func formatDate(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 24*time.Hour {
		return "Today"
	}
	if diff < 48*time.Hour {
		return "Yesterday"
	}
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
	if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	// Format as date
	return t.Format("January 2, 2006")
}
