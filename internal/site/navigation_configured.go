package site

import (
	"fmt"
	"strings"

	"github.com/p-arndt/dorcs/internal/config"
)

func buildConfiguredNavTree(index map[string]*Doc, items config.NavItems) (*NavNode, error) {
	root := &NavNode{Name: "", Key: "", IsDir: true}
	if rootDoc, ok := index[""]; ok {
		root.Page = rootDoc
		root.Name = rootDoc.Title
	}

	seenPages := make(map[string]string)
	nodes, err := buildConfiguredNavNodes(index, items, seenPages, "")
	if err != nil {
		return nil, err
	}
	root.Children = nodes
	return root, nil
}

func buildConfiguredNavNodes(index map[string]*Doc, items config.NavItems, seenPages map[string]string, parentPath string) ([]*NavNode, error) {
	nodes := make([]*NavNode, 0, len(items))

	for _, item := range items {
		node, include, err := buildConfiguredNavNode(index, item, seenPages, parentPath)
		if err != nil {
			return nil, err
		}
		if include {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func buildConfiguredNavNode(index map[string]*Doc, item config.NavItemConfig, seenPages map[string]string, parentPath string) (*NavNode, bool, error) {
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
			return nil, false, fmt.Errorf("nav item %q references missing doc %q", entryPath, item.Page)
		}
		page = doc
		key = doc.Key
		if prev, exists := seenPages[doc.Key]; exists {
			return nil, false, fmt.Errorf("nav item %q references %q more than once (already used by %q)", entryPath, item.Page, prev)
		}
		seenPages[doc.Key] = entryPath
	}

	childNodes, err := buildConfiguredNavNodes(index, item.Items, seenPages, entryPath)
	if err != nil {
		return nil, false, err
	}

	if page != nil && page.Key == "" && len(childNodes) == 0 {
		// Root index is already rendered as the dedicated home link in the sidebar.
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

func configuredGroupKey(labelPath string) string {
	slug := strings.ToLower(strings.TrimSpace(labelPath))
	slug = strings.ReplaceAll(slug, " > ", "/")
	slug = strings.ReplaceAll(slug, " ", "-")
	if slug == "" {
		slug = "group"
	}
	return "__group__/" + slug
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
