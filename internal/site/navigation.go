package site

import (
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// NavTree returns the cached navigation tree for sidebar rendering.
// Draft documents are filtered if includeDraft is false.
func (s *Site) NavTree(includeDraft bool) *NavNode {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// If including drafts, return as-is.
	if includeDraft {
		return s.nav
	}

	// Otherwise filter drafts out (copy-on-read).
	return filterNavDrafts(s.nav)
}

// buildNavTree creates a directory tree suitable for a sidebar from the indexed docs.
// Folder landing pages (index.md) are attached to the folder node via NavNode.Page.
func buildNavTree(index map[string]*Doc) *NavNode {
	root := &NavNode{Name: "", Key: "", IsDir: true}

	// Helper to find or create a folder node
	findOrCreateFolder := func(parent *NavNode, name, key string) *NavNode {
		for _, c := range parent.Children {
			if c.IsDir && c.Key == key {
				return c
			}
		}
		child := &NavNode{Name: name, Key: key, IsDir: true}
		parent.Children = append(parent.Children, child)
		return child
	}

	// Insert all docs into folder structure.
	for key, d := range index {
		// key can be "" (root index)
		if key == "" {
			root.Page = d
			root.Name = d.Title // Use title from index.md
			continue
		}

		parts := strings.Split(key, "/")
		cur := root

		// Check if this is a folder landing page (index.md)
		if isFolderLandingDoc(d) {
			// For folder index: create/find folder nodes for ALL parts
			// e.g., "architecture" -> create "architecture" folder
			// e.g., "guide/advanced" -> create "guide" then "guide/advanced"
			for i := 0; i < len(parts); i++ {
				dirName := parts[i]
				dirKey := strings.Join(parts[:i+1], "/")
				cur = findOrCreateFolder(cur, dirName, dirKey)
			}
			// Attach the index.md to this folder and use its title
			cur.Page = d
			if d.Title != "" {
				cur.Name = d.Title
			}
			continue
		}

		// Regular page: walk/create parent directories, then add leaf
		for i := 0; i < len(parts)-1; i++ {
			dirName := parts[i]
			dirKey := strings.Join(parts[:i+1], "/")
			cur = findOrCreateFolder(cur, dirName, dirKey)
		}

		leaf := &NavNode{
			Name:  path.Base(key),
			Key:   key,
			IsDir: false,
			Page:  d,
		}
		cur.Children = append(cur.Children, leaf)
	}

	// Sort children consistently: dirs first, then pages; by title/name.
	sortNav(root)
	return root
}

// isFolderLandingDoc checks if a document is a folder landing page (index.md).
func isFolderLandingDoc(d *Doc) bool {
	// A folder landing doc is represented by a RelPath ending in "/index.md" (or "index.md" at root),
	// and its Key is the folder key ("" for root, "guide", ...).
	rel := filepath.ToSlash(d.RelPath)
	return rel == "index.md" || strings.HasSuffix(rel, "/index.md")
}

// sortNav recursively sorts navigation nodes.
func sortNav(n *NavNode) {
	if n == nil {
		return
	}

	sort.SliceStable(n.Children, func(i, j int) bool {
		a, b := n.Children[i], n.Children[j]
		if a.IsDir != b.IsDir {
			return a.IsDir // dirs first
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

	for _, c := range n.Children {
		if c.IsDir {
			sortNav(c)
		}
	}
}

// filterNavDrafts creates a copy of the navigation tree with draft documents filtered out.
func filterNavDrafts(n *NavNode) *NavNode {
	if n == nil {
		return nil
	}

	cp := &NavNode{
		Name:     n.Name,
		Key:      n.Key,
		IsDir:    n.IsDir,
		Page:     n.Page,
		Children: nil,
	}

	// Drop draft landing pages
	if cp.Page != nil && cp.Page.Draft {
		cp.Page = nil
	}

	for _, c := range n.Children {
		if c == nil {
			continue
		}
		if c.IsDir {
			cc := filterNavDrafts(c)
			// Keep dir if it has a non-draft landing page or any children.
			if cc.Page != nil || len(cc.Children) > 0 {
				cp.Children = append(cp.Children, cc)
			}
			continue
		}
		// Leaf page
		if c.Page != nil && c.Page.Draft {
			continue
		}
		cp.Children = append(cp.Children, c)
	}

	return cp
}
