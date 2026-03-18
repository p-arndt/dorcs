package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/p-arndt/dorcs/internal/config"
)

func TestBuildNavTree(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.md": `---
title: Home
---
# Welcome`,
		"getting-started.md": `---
title: Getting Started
---
# Getting Started`,
		"guide/index.md": `---
title: Guide
---
# Guide`,
		"guide/intro.md": `---
title: Introduction
---
# Intro`,
		"api/v1/index.md": `---
title: API v1
---
# API`,
		"api/v1/endpoints.md": `---
title: Endpoints
---
# Endpoints`,
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %q: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %q: %v", fullPath, err)
		}
	}

	s, err := New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if err := s.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex() failed: %v", err)
	}

	t.Run("nav tree structure", func(t *testing.T) {
		nav := s.NavTree(false)

		if nav == nil {
			t.Fatal("NavTree should not be nil")
		}

		// Root should have children
		if len(nav.Children) == 0 {
			t.Error("Root nav should have children")
		}

		// Check that guide folder exists
		foundGuide := false
		for _, child := range nav.Children {
			if child.Key == "guide" && child.IsDir {
				foundGuide = true
				if child.Page == nil {
					t.Error("Guide folder should have a page (index.md)")
				}
				if len(child.Children) == 0 {
					t.Error("Guide folder should have children")
				}
				break
			}
		}
		if !foundGuide {
			t.Error("Guide folder not found in nav tree")
		}
	})

	t.Run("nav tree with drafts", func(t *testing.T) {
		// Add a draft file
		draftPath := filepath.Join(tmpDir, "draft.md")
		if err := os.WriteFile(draftPath, []byte("---\ndraft: true\n---\n# Draft"), 0644); err != nil {
			t.Fatalf("failed to write draft: %v", err)
		}

		if err := s.BuildIndex(); err != nil {
			t.Fatalf("BuildIndex() failed: %v", err)
		}

		navWithoutDrafts := s.NavTree(false)
		navWithDrafts := s.NavTree(true)

		// Count total items (recursively)
		var countItems func(*NavNode) int
		countItems = func(n *NavNode) int {
			if n == nil {
				return 0
			}
			count := 0
			if n.Page != nil {
				count = 1
			}
			for _, child := range n.Children {
				count += countItems(child)
			}
			return count
		}

		countWithout := countItems(navWithoutDrafts)
		countWith := countItems(navWithDrafts)

		if countWith <= countWithout {
			t.Errorf("Nav with drafts should have more items, got %d vs %d", countWith, countWithout)
		}
	})
}

func TestIsFolderLandingDoc(t *testing.T) {
	tests := []struct {
		name     string
		doc      *Doc
		expected bool
	}{
		{
			name: "root index",
			doc: &Doc{
				RelPath: "index.md",
				Key:     "",
			},
			expected: true,
		},
		{
			name: "folder index",
			doc: &Doc{
				RelPath: "guide/index.md",
				Key:     "guide",
			},
			expected: true,
		},
		{
			name: "regular file",
			doc: &Doc{
				RelPath: "guide/intro.md",
				Key:     "guide/intro",
			},
			expected: false,
		},
		{
			name: "nested folder index",
			doc: &Doc{
				RelPath: "api/v1/index.md",
				Key:     "api/v1",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFolderLandingDoc(tt.doc)
			if result != tt.expected {
				t.Errorf("isFolderLandingDoc() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestFilterNavDrafts(t *testing.T) {
	// Create a nav tree with drafts
	root := &NavNode{
		Name:  "Root",
		Key:   "",
		IsDir: true,
		Page: &Doc{
			Title: "Home",
			Draft: false,
		},
		Children: []*NavNode{
			{
				Name:  "Draft Page",
				Key:   "draft",
				IsDir: false,
				Page: &Doc{
					Title: "Draft",
					Draft: true,
				},
			},
			{
				Name:  "Published Page",
				Key:   "published",
				IsDir: false,
				Page: &Doc{
					Title: "Published",
					Draft: false,
				},
			},
			{
				Name:  "Draft Folder",
				Key:   "draft-folder",
				IsDir: true,
				Page: &Doc{
					Title: "Draft Folder",
					Draft: true,
				},
				Children: []*NavNode{
					{
						Name:  "Child",
						Key:   "draft-folder/child",
						IsDir: false,
						Page: &Doc{
							Title: "Child",
							Draft: false,
						},
					},
				},
			},
		},
	}

	filtered := filterNavDrafts(root)

	// Should have published page
	foundPublished := false
	for _, child := range filtered.Children {
		if child.Key == "published" {
			foundPublished = true
		}
		if child.Key == "draft" {
			t.Error("Draft page should be filtered out")
		}
	}
	if !foundPublished {
		t.Error("Published page should be present")
	}

	// Draft folder with children should be kept (but page removed)
	foundDraftFolder := false
	for _, child := range filtered.Children {
		if child.Key == "draft-folder" {
			foundDraftFolder = true
			if child.Page != nil {
				t.Error("Draft folder page should be removed")
			}
			if len(child.Children) == 0 {
				t.Error("Draft folder should keep its children")
			}
		}
	}
	if !foundDraftFolder {
		t.Error("Draft folder with children should be kept")
	}
}

func TestBuildIndexWithExplicitNav(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"index.md": `---
title: Home
---
# Home`,
		"01_getting-started.md": `---
title: Getting Started
---
# Getting Started`,
		"usage/index.md": `---
title: Usage Index
---
# Usage`,
		"usage/writing-your-docs.md": `---
title: Writing
---
# Writing`,
		"usage/metadata.md": `---
title: Metadata
---
# Metadata`,
		"external-content/github.md": `---
title: GitHub
---
# GitHub`,
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %q: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %q: %v", fullPath, err)
		}
	}

	s, err := New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	s.SetExplicitNav(config.NavItems{
		{Label: "Home", Page: "index.md"},
		{Label: "Getting Started", Page: "01_getting-started.md"},
		{
			Label: "Usage",
			Page:  "usage/index.md",
			Items: []config.NavItemConfig{
				{Label: "Writing Your Docs", Page: "usage/writing-your-docs.md"},
				{Label: "Metadata", Page: "usage/metadata.md"},
			},
		},
		{
			Label: "External Content",
			Items: []config.NavItemConfig{
				{Label: "GitHub", Page: "external-content/github.md"},
			},
		},
	})

	if err := s.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex() failed: %v", err)
	}

	nav := s.NavTree(false)
	if nav == nil {
		t.Fatal("NavTree should not be nil")
	}
	if len(nav.Children) != 3 {
		t.Fatalf("expected 3 root nav children (root index should not duplicate), got %d", len(nav.Children))
	}

	if nav.Children[0].Name != "Getting Started" || nav.Children[0].Key != "01_getting-started" {
		t.Fatalf("unexpected first nav child: %+v", nav.Children[0])
	}

	usage := nav.Children[1]
	if !usage.IsDir || usage.Page == nil || usage.Key != "usage" {
		t.Fatalf("expected clickable Usage section, got %+v", usage)
	}
	if len(usage.Children) != 2 || usage.Children[0].Name != "Writing Your Docs" {
		t.Fatalf("unexpected Usage children: %+v", usage.Children)
	}

	external := nav.Children[2]
	if !external.IsDir || external.Page != nil {
		t.Fatalf("expected non-clickable External Content group, got %+v", external)
	}
	if !strings.HasPrefix(external.Key, "__group__/") {
		t.Fatalf("expected synthetic group key, got %q", external.Key)
	}

	indexHTML := string(BuildIndex(PlaceholderContext{Site: s}))
	if !strings.Contains(indexHTML, "Getting Started") || !strings.Contains(indexHTML, "External Content") {
		t.Fatalf("expected BuildIndex to use explicit nav tree, got %q", indexHTML)
	}
}

func TestBuildIndexWithExplicitNavValidationErrors(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatalf("failed to write index: %v", err)
	}

	s, err := New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	s.SetExplicitNav(config.NavItems{
		{Label: "Missing", Page: "missing.md"},
	})

	if err := s.BuildIndex(); err == nil {
		t.Fatal("expected BuildIndex to fail for missing explicit nav doc")
	}
}

func TestBuildIndexKeepsPreviousStateWhenExplicitNavFails(t *testing.T) {
	tmpDir := t.TempDir()
	files := map[string]string{
		"index.md": "# Home",
		"guide.md": "# Guide",
	}
	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %q: %v", fullPath, err)
		}
	}

	s, err := New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	s.SetExplicitNav(config.NavItems{{Label: "Guide", Page: "guide.md"}})
	if err := s.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex() failed: %v", err)
	}

	if err := os.Remove(filepath.Join(tmpDir, "guide.md")); err != nil {
		t.Fatalf("failed to remove guide doc: %v", err)
	}

	if err := s.BuildIndex(); err == nil {
		t.Fatal("expected BuildIndex to fail after removing an explicitly referenced doc")
	}

	if _, ok := s.GetDoc("guide"); !ok {
		t.Fatal("expected previous index to remain available after failed rebuild")
	}

	nav := s.NavTree(false)
	if nav == nil || len(nav.Children) != 1 || nav.Children[0].Key != "guide" {
		t.Fatalf("expected previous nav tree to remain intact after failed rebuild, got %+v", nav)
	}
}

func TestBuildIndexWithExplicitNavSkipsMissingDocsForNonDefaultVersion(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"index.md":    "# Latest Home",
		"guide.md":    "# Latest Guide",
		"v1/index.md": "# V1 Home",
	}
	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %q: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %q: %v", fullPath, err)
		}
	}

	s, err := NewWithVersion(tmpDir, "github", "", "", "v1")
	if err != nil {
		t.Fatalf("NewWithVersion() failed: %v", err)
	}
	s.SetDefaultVersion("latest")
	s.SetExplicitNav(config.NavItems{
		{Label: "Home", Page: "index.md"},
		{Label: "Guide", Page: "guide.md"},
	})

	if err := s.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex() failed: %v", err)
	}

	nav := s.NavTree(false)
	if nav == nil {
		t.Fatal("NavTree should not be nil")
	}
	if len(nav.Children) != 0 {
		t.Fatalf("expected missing non-default version docs to be skipped, got %+v", nav.Children)
	}
}

func TestBuildIndexWithExplicitNavAndGitHubPreservesConfiguredOrder(t *testing.T) {
	tmpDir := t.TempDir()

	s, err := New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	mockClient := newMockGitHubClient()
	owner := "test-owner"
	repo := "test-repo"
	branch := "main"
	repoPath := "docs"

	mockClient.discoverFiles[owner+"/"+repo+"/"+branch+"/"+repoPath] = []string{
		"index.md",
		"usage/watch-mode.md",
		"usage/index.md",
		"01_getting-started.md",
		"usage/file-structure.md",
	}

	mockClient.fetchContent[owner+"/"+repo+"/"+branch+"/"+repoPath+"/index.md"] = []byte(`---
title: Home
---
# Home

[[INDEX]]`)
	mockClient.fetchContent[owner+"/"+repo+"/"+branch+"/"+repoPath+"/01_getting-started.md"] = []byte(`---
title: Getting Started
---
# Getting Started`)
	mockClient.fetchContent[owner+"/"+repo+"/"+branch+"/"+repoPath+"/usage/index.md"] = []byte(`---
title: Usage
---
# Usage

[[CHILDREN]]`)
	mockClient.fetchContent[owner+"/"+repo+"/"+branch+"/"+repoPath+"/usage/file-structure.md"] = []byte(`---
title: File Structure
---
# File Structure`)
	mockClient.fetchContent[owner+"/"+repo+"/"+branch+"/"+repoPath+"/usage/watch-mode.md"] = []byte(`---
title: Watch Mode
---
# Watch Mode`)

	s.SetGitHubConfig(mockClient, owner, repo, branch, repoPath)
	s.SetExplicitNav(config.NavItems{
		{Label: "Home", Page: "index.md"},
		{
			Label: "Usage",
			Page:  "usage/index.md",
			Items: []config.NavItemConfig{
				{Label: "Watch Mode", Page: "usage/watch-mode.md"},
				{Label: "File Structure", Page: "usage/file-structure.md"},
			},
		},
		{Label: "Getting Started", Page: "01_getting-started.md"},
	})

	if err := s.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex() failed: %v", err)
	}

	nav := s.NavTree(false)
	if nav == nil {
		t.Fatal("NavTree should not be nil")
	}
	if len(nav.Children) != 2 {
		t.Fatalf("expected 2 root nav children (home should not duplicate), got %d", len(nav.Children))
	}
	if nav.Children[0].Key != "usage" || nav.Children[1].Key != "01_getting-started" {
		t.Fatalf("expected configured root order to be preserved, got %+v", nav.Children)
	}

	if len(nav.Children[0].Children) != 2 {
		t.Fatalf("expected Usage to have 2 children, got %d", len(nav.Children[0].Children))
	}
	if nav.Children[0].Children[0].Key != "usage/watch-mode" || nav.Children[0].Children[1].Key != "usage/file-structure" {
		t.Fatalf("expected configured child order to be preserved, got %+v", nav.Children[0].Children)
	}

	rootIndexHTML := string(BuildIndex(PlaceholderContext{Site: s, BasePath: s.BasePath}))
	usagePos := strings.Index(rootIndexHTML, ">Usage<")
	gettingStartedPos := strings.Index(rootIndexHTML, ">Getting Started<")
	if usagePos == -1 || gettingStartedPos == -1 || usagePos >= gettingStartedPos {
		t.Fatalf("expected site index to follow configured root order, got %q", rootIndexHTML)
	}

	usageChildrenHTML := string(BuildChildren(PlaceholderContext{Site: s, CurrentKey: "usage", BasePath: s.BasePath}))
	watchModePos := strings.Index(usageChildrenHTML, ">Watch Mode<")
	fileStructurePos := strings.Index(usageChildrenHTML, ">File Structure<")
	if watchModePos == -1 || fileStructurePos == -1 || watchModePos >= fileStructurePos {
		t.Fatalf("expected children placeholder to follow configured child order, got %q", usageChildrenHTML)
	}
}
