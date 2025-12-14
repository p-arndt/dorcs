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

	// Build a map of keys to nodes for resolving "after" references
	// Keys can be folder keys (for directories) or page keys (for leaf pages)
	keyMap := make(map[string]*NavNode)
	var indexNode *NavNode // The index.md node (folder landing page)

	for _, child := range n.Children {
		// Map by the node's key (folder key for dirs, page key for leaf pages)
		keyMap[child.Key] = child
		// For pages (both leaf and folder with Page), also map by the page's Key field if different
		if child.Page != nil && child.Page.Key != child.Key {
			keyMap[child.Page.Key] = child
		}
		// Track index.md node (folder landing page) - this is a directory with a Page
		if child.IsDir && child.Page != nil {
			indexNode = child
			// Also map "index" to this node for easy reference
			keyMap["index"] = child
		}
	}

	sort.SliceStable(n.Children, func(i, j int) bool {
		a, b := n.Children[i], n.Children[j]

		// Get "after" references - simple: just specify what item this should come after
		aAfter := ""
		bAfter := ""
		if a.Page != nil {
			aAfter = a.Page.After
		}
		if b.Page != nil {
			bAfter = b.Page.After
		}

		// "after" relationships take ABSOLUTE precedence over everything else
		// Special handling for root level: root index is n.Page, not in children
		isRootLevel := n.Key == ""

		// Helper to resolve "after" target
		resolveAfterTarget := func(after string) *NavNode {
			if after == "" {
				return nil
			}
			if after == "index" {
				// For root level, index is the parent's Page (not in children)
				// For folder level, index is a child directory with a Page
				if isRootLevel && n.Page != nil {
					// Root index - return a special marker (we'll use n as marker)
					return n // Use parent node as marker for root index
				}
				return indexNode
			}
			// Try exact match first
			if target, ok := keyMap[after]; ok {
				return target
			}
			// Try matching by stripping numeric prefixes (e.g., "getting-started" matches "01_getting-started")
			// This allows referencing pages without their numeric prefix
			for key, target := range keyMap {
				// Check both the node's key and the page's key (if it exists)
				keysToCheck := []string{key}
				if target.Page != nil && target.Page.Key != key {
					keysToCheck = append(keysToCheck, target.Page.Key)
				}

				for _, checkKey := range keysToCheck {
					// Extract the base name from the key (last part after /)
					keyBase := checkKey
					if idx := strings.LastIndex(checkKey, "/"); idx >= 0 {
						keyBase = checkKey[idx+1:]
					}
					// Strip numeric prefix from keyBase (e.g., "01_getting-started" -> "getting-started")
					keyBaseNoPrefix := keyBase
					for i, r := range keyBase {
						if r >= '0' && r <= '9' {
							continue
						}
						if i > 0 && (r == '_' || r == '-' || r == ' ') {
							keyBaseNoPrefix = keyBase[i+1:]
						} else {
							keyBaseNoPrefix = keyBase[i:]
						}
						break
					}
					// If the stripped key matches the after value, use it
					if keyBaseNoPrefix == after {
						return target
					}
				}
			}
			return nil
		}

		aTarget := resolveAfterTarget(aAfter)
		bTarget := resolveAfterTarget(bAfter)

		// CRITICAL: Items with "after: index" must come right after the index
		// At root level, this means they come before all other items
		// At folder level, this means they come right after the folder's index.md
		if aAfter == "index" && bAfter != "index" {
			// a wants to be after index, b doesn't - a comes first (right after index)
			// This works for both root level and folder level
			return true
		}
		if bAfter == "index" && aAfter != "index" {
			// b wants to be after index, a doesn't - b comes first (right after index)
			return false
		}

		// Direct relationship: a wants to be after b (b is a's target)
		if aTarget != nil && aTarget != n { // n is only used as marker for root index
			if aTarget == b {
				return false // a wants to be after b, so b comes before a
			}
		}

		// Direct relationship: b wants to be after a (a is b's target)
		if bTarget != nil && bTarget != n { // n is only used as marker for root index
			if bTarget == a {
				return true // b wants to be after a, so a comes before b
			}
		}

		// Extract numeric prefixes from filenames or folder keys
		aPrefix := 0
		bPrefix := 0
		if a.IsDir {
			aPrefix = extractNumericPrefix(a.Key)
		} else if a.Page != nil {
			aPrefix = extractNumericPrefix(a.Page.RelPath)
		}
		if b.IsDir {
			bPrefix = extractNumericPrefix(b.Key)
		} else if b.Page != nil {
			bPrefix = extractNumericPrefix(b.Page.RelPath)
		}

		// If a has "after: index" and b doesn't, but b would normally come after index (has prefix),
		// then a should come before b (a wants to be right after index)
		if aTarget == indexNode && bTarget == nil && bPrefix > 0 {
			return true // a (after index) comes before b (numbered item after index)
		}

		// If b has "after: index" and a doesn't, but a would normally come after index (has prefix),
		// then b should come before a
		if bTarget == indexNode && aTarget == nil && aPrefix > 0 {
			return false // b (after index) comes before a (numbered item after index)
		}

		// If a has an "after" target (not index) and b doesn't, check if b would come after the target
		// If so, a should come before b (a wants to be right after its target)
		if aTarget != nil && aTarget != n && aTarget != indexNode && bTarget == nil && aTarget != b {
			// Check if b would normally come after aTarget by comparing their positions
			// We can do this by checking if b has a prefix that's greater than aTarget's prefix
			targetPrefix := 0
			if aTarget.IsDir {
				targetPrefix = extractNumericPrefix(aTarget.Key)
			} else if aTarget.Page != nil {
				targetPrefix = extractNumericPrefix(aTarget.Page.RelPath)
			}
			// If b has a prefix and it's greater than the target's prefix, a should come before b
			if bPrefix > 0 && targetPrefix > 0 && bPrefix > targetPrefix {
				return true // a (after target) comes before b (numbered item after target)
			}
			// Also, if target has no prefix but b has a prefix, a should come before b
			// (items with after should come right after their target, before numbered items)
			if targetPrefix == 0 && bPrefix > 0 {
				return true
			}
		}

		// If b has an "after" target (not index) and a doesn't, similar logic
		if bTarget != nil && bTarget != n && bTarget != indexNode && aTarget == nil && bTarget != a {
			targetPrefix := 0
			if bTarget.IsDir {
				targetPrefix = extractNumericPrefix(bTarget.Key)
			} else if bTarget.Page != nil {
				targetPrefix = extractNumericPrefix(bTarget.Page.RelPath)
			}
			if aPrefix > 0 && targetPrefix > 0 && aPrefix > targetPrefix {
				return false // b (after target) comes before a (numbered item after target)
			}
			if targetPrefix == 0 && aPrefix > 0 {
				return false
			}
		}

		// If both have "after" to the same target, sort them normally (they're both after the same thing)
		// If both have "after" to different targets, we'll handle it in normal sorting

		// Get order values (0 means no order specified)
		aOrder := 0
		bOrder := 0
		if a.Page != nil {
			aOrder = a.Page.Order
		}
		if b.Page != nil {
			bOrder = b.Page.Order
		}

		// Sorting logic (intuitive):
		// 1. Numeric prefixes are compared first (00 < 01 < 02 < ...)
		// 2. Order fields are only used when there's no prefix, or to fine-tune within same prefix
		// 3. Items with prefixes always come before items with only order fields
		// 4. This means: 00_ < 01_ < 02_ < ... < items with order fields

		// Both have numeric prefixes - compare them directly
		if aPrefix != 0 && bPrefix != 0 {
			if aPrefix != bPrefix {
				return aPrefix < bPrefix
			}
			// Prefixes are equal, check order fields for fine-tuning within same prefix
			if aOrder != 0 && bOrder != 0 {
				if aOrder != bOrder {
					return aOrder < bOrder
				}
			} else if aOrder != 0 {
				// a has order for fine-tuning, b doesn't - a comes first
				return true
			} else if bOrder != 0 {
				// b has order for fine-tuning, a doesn't - b comes first
				return false
			}
			// Both have same prefix and no order, continue to title/key
		} else if aPrefix != 0 {
			// a has prefix, b doesn't - a always comes first (prefixes take priority)
			return true
		} else if bPrefix != 0 {
			// b has prefix, a doesn't - b always comes first (prefixes take priority)
			return false
		}

		// Neither has prefix - check order fields
		if aOrder != 0 && bOrder != 0 {
			if aOrder != bOrder {
				return aOrder < bOrder
			}
			// Orders are equal, continue to title/key
		} else if aOrder != 0 {
			// a has order, b doesn't - a comes first
			return true
		} else if bOrder != 0 {
			// b has order, a doesn't - b comes first
			return false
		}
		// Both have no sort value, fall through to title/key sorting

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

	// Re-sort the filtered tree to maintain correct order
	sortNav(cp)

	return cp
}
