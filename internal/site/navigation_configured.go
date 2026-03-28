package site

import (
	"fmt"
	"strings"

	"github.com/p-arndt/dorcs/internal/config"
)

func buildConfiguredNavTree(site *Site, index map[string]*Doc, items config.NavItems) (*NavNode, error) {
	root := &NavNode{Name: "", Key: "", IsDir: true}
	if rootDoc, ok := index[""]; ok {
		root.Page = rootDoc
		root.Name = rootDoc.Title
	}

	seenPages := make(map[string]string)
	nodes, err := buildConfiguredNavNodes(site, index, items, seenPages, "")
	if err != nil {
		return nil, err
	}
	root.Children = nodes
	return root, nil
}

func buildConfiguredNavNodes(site *Site, index map[string]*Doc, items config.NavItems, seenPages map[string]string, parentPath string) ([]*NavNode, error) {
	nodes := make([]*NavNode, 0, len(items))

	for _, item := range items {
		node, include, err := buildConfiguredNavNode(site, index, item, seenPages, parentPath)
		if err != nil {
			return nil, err
		}
		if include {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func buildConfiguredNavNode(site *Site, index map[string]*Doc, item config.NavItemConfig, seenPages map[string]string, parentPath string) (*NavNode, bool, error) {
	entryPath := item.Label
	if parentPath != "" {
		entryPath = parentPath + " > " + item.Label
	}

	var page *Doc
	var key string
	if item.Page != "" {
		resolvedKey := keyFromRel(strings.TrimSpace(item.Page))
		if strings.TrimSpace(item.Page) == "" || (resolvedKey == "" && !isIndexRel(strings.TrimSpace(item.Page))) {
			return nil, false, fmt.Errorf("nav item %q references invalid markdown path %q", entryPath, item.Page)
		}
		doc, ok := index[resolvedKey]
		if !ok {
			if allowMissingConfiguredNavPage(site) {
				return nil, false, nil
			}
			return nil, false, fmt.Errorf("nav item %q references missing doc %q", entryPath, item.Page)
		}
		page = doc
		key = doc.Key
		if prev, exists := seenPages[doc.Key]; exists {
			return nil, false, fmt.Errorf("nav item %q references %q more than once (already used by %q)", entryPath, item.Page, prev)
		}
		seenPages[doc.Key] = entryPath
	}

	childNodes, err := buildConfiguredNavNodes(site, index, item.Items, seenPages, entryPath)
	if err != nil {
		return nil, false, err
	}

	if page != nil && page.Key == "" && len(childNodes) == 0 && !site.hasSections() {
		// Root index is already rendered as the dedicated home link in the sidebar.
		// When sections are configured, the home link is hidden, so include it.
		return nil, false, nil
	}

	node := &NavNode{
		Name:          item.Label,
		ExplicitTitle: true,
		Key:           key,
		Page:          page,
		Children:      childNodes,
	}
	if len(childNodes) > 0 {
		node.IsDir = true
		if node.Key == "" {
			node.Key = configuredGroupKey(entryPath)
		}
	} else if page == nil {
		node.IsDir = true
		node.Key = configuredGroupKey(entryPath)
	}

	return node, true, nil
}

func allowMissingConfiguredNavPage(site *Site) bool {
	if site == nil {
		return false
	}
	// In multi-version mode, a shared nav config may reference docs that only exist
	// in the default version. Non-default version sites should omit those entries
	// instead of failing the whole index build.
	return strings.TrimSpace(site.DefaultVersion) != "" && strings.TrimSpace(site.Version) != ""
}

func configuredGroupKey(labelPath string) string {
	slug := strings.ToLower(strings.TrimSpace(labelPath))
	slug = strings.ReplaceAll(slug, " > ", "/")
	slug = strings.ReplaceAll(slug, " ", "-")
	if slug == "" {
		slug = "group"
	}
	return "__group__/" + slug
}

// NavTreeFromItems builds navigation nodes from a set of config NavItems using
// the site's document index. Used by section tabs to build per-section navigation.
// Returns the children nodes (not the root).
func (s *Site) NavTreeFromItems(items config.NavItems, includeDraft bool) []*NavNode {
	s.mu.RLock()
	index := s.index
	s.mu.RUnlock()

	if len(items) == 0 || len(index) == 0 {
		return nil
	}

	seenPages := make(map[string]string)
	nodes, err := buildConfiguredNavNodes(s, index, items, seenPages, "")
	if err != nil {
		return nil
	}

	if !includeDraft {
		filtered := make([]*NavNode, 0, len(nodes))
		for _, n := range nodes {
			fn := filterNavDrafts(n)
			if fn != nil && (fn.Page != nil || len(fn.Children) > 0) {
				filtered = append(filtered, fn)
			}
		}
		return filtered
	}

	return nodes
}

// hasSections returns true when nav.sections is configured.
func (s *Site) hasSections() bool {
	if s == nil {
		return false
	}
	return s.sectionsConfigured
}

func navNodeDisplayName(node *NavNode) string {
	if node == nil {
		return ""
	}
	if node.ExplicitTitle && strings.TrimSpace(node.Name) != "" {
		return node.Name
	}
	if node.Page != nil && strings.TrimSpace(node.Page.Title) != "" {
		return node.Page.Title
	}
	if strings.TrimSpace(node.Name) != "" {
		return node.Name
	}
	return node.Key
}

func navNodeURL(node *NavNode, basePath string, site *Site) string {
	if node == nil {
		return ""
	}
	if node.Page == nil && node.IsDir {
		return ""
	}
	return buildDocURL(node.Key, basePath, site)
}
