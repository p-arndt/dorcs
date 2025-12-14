package site

import (
	"os"
	"path/filepath"
	"testing"
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

	s, err := New(tmpDir, "github", "")
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
