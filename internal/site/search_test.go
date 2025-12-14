package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchDocs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.md": `---
title: Home
---
# Welcome to Documentation`,
		"getting-started.md": `---
title: Getting Started
description: Learn how to get started
tags: [tutorial, beginner]
---
# Getting Started

This guide will help you get started.`,
		"advanced.md": `---
title: Advanced Topics
description: Advanced features and techniques
---
# Advanced Topics

Here are some advanced concepts.`,
		"guide/index.md": `---
title: Guide
---
# Guide

Welcome to the guide section.`,
		"guide/installation.md": `---
title: Installation
---
# Installation

Install the software using these steps.`,
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

	t.Run("search by title", func(t *testing.T) {
		results := s.SearchDocs("Getting Started", false, 10)

		if len(results) == 0 {
			t.Error("Should find 'Getting Started' document")
			return
		}

		// First result should be the matching document
		if results[0].Title != "Getting Started" {
			t.Errorf("First result title = %q; want %q", results[0].Title, "Getting Started")
		}

		// Should have high score for title match
		if results[0].Score < 100 {
			t.Errorf("Title match should have high score, got %d", results[0].Score)
		}
	})

	t.Run("search by content", func(t *testing.T) {
		results := s.SearchDocs("installation", false, 10)

		if len(results) == 0 {
			t.Error("Should find document containing 'installation'")
			return
		}

		found := false
		for _, r := range results {
			if r.Title == "Installation" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should find Installation document")
		}
	})

	t.Run("search multiple words", func(t *testing.T) {
		results := s.SearchDocs("advanced topics", false, 10)

		if len(results) == 0 {
			t.Error("Should find documents matching multiple words")
		}
	})

	t.Run("empty query returns nil", func(t *testing.T) {
		results := s.SearchDocs("", false, 10)
		if results != nil {
			t.Errorf("Empty query should return nil, got %v", results)
		}
	})

	t.Run("whitespace query returns nil", func(t *testing.T) {
		results := s.SearchDocs("   ", false, 10)
		if results != nil {
			t.Errorf("Whitespace query should return nil, got %v", results)
		}
	})

	t.Run("no matches returns empty", func(t *testing.T) {
		results := s.SearchDocs("nonexistentterm12345", false, 10)
		if len(results) != 0 {
			t.Errorf("No matches should return empty slice, got %d results", len(results))
		}
	})

	t.Run("max results limit", func(t *testing.T) {
		results := s.SearchDocs("guide", false, 2)
		if len(results) > 2 {
			t.Errorf("Results should be limited to 2, got %d", len(results))
		}
	})

	t.Run("excludes drafts", func(t *testing.T) {
		// Add a draft file
		draftPath := filepath.Join(tmpDir, "draft.md")
		if err := os.WriteFile(draftPath, []byte("---\ndraft: true\n---\n# Draft Content\nSearch term here."), 0644); err != nil {
			t.Fatalf("failed to write draft: %v", err)
		}

		if err := s.BuildIndex(); err != nil {
			t.Fatalf("BuildIndex() failed: %v", err)
		}

		results := s.SearchDocs("Search term", false, 10)
		for _, r := range results {
			if r.Title == "Draft" {
				t.Error("Draft documents should be excluded from search")
			}
		}

		// With drafts included
		resultsWithDrafts := s.SearchDocs("Search term", true, 10)
		foundDraft := false
		for _, r := range resultsWithDrafts {
			if r.Title == "Draft" || strings.Contains(r.Key, "draft") {
				foundDraft = true
				break
			}
		}
		if !foundDraft {
			t.Error("Draft documents should be included when includeDraft is true")
		}
	})
}
