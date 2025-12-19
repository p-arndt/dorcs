package site

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewWithLanguage(t *testing.T) {
	tmpDir := t.TempDir()

	// Create default language files
	files := map[string]string{
		"index.md": `---
title: Home
---
# Welcome`,
		"getting-started.md": `---
title: Getting Started
---
# Getting Started`,
		"__lang__/de/index.md": `---
title: Startseite
---
# Willkommen`,
		"__lang__/de/getting-started.md": `---
title: Erste Schritte
---
# Erste Schritte`,
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

	t.Run("default language - empty string", func(t *testing.T) {
		s, err := New(tmpDir, "github", "", "")
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}

		if s.Language != "" {
			t.Errorf("expected Language to be empty, got %q", s.Language)
		}

		if err := s.BuildIndex(); err != nil {
			t.Fatalf("BuildIndex() failed: %v", err)
		}

		// Should find default language files
		if _, ok := s.GetDoc(""); !ok {
			t.Error("expected to find root index")
		}
		if _, ok := s.GetDoc("getting-started"); !ok {
			t.Error("expected to find getting-started")
		}
		// Should NOT find German files
		if _, ok := s.GetDoc("__lang__/de/index"); ok {
			t.Error("should NOT find __lang__ folder in default language site")
		}
	})

	t.Run("German language", func(t *testing.T) {
		s, err := New(tmpDir, "github", "", "de")
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}

		if s.Language != "de" {
			t.Errorf("expected Language to be 'de', got %q", s.Language)
		}

		if err := s.BuildIndex(); err != nil {
			t.Fatalf("BuildIndex() failed: %v", err)
		}

		// Should find German language files
		if _, ok := s.GetDoc(""); !ok {
			t.Error("expected to find German root index")
		}
		if _, ok := s.GetDoc("getting-started"); !ok {
			t.Error("expected to find German getting-started")
		}

		// Check that files are from __lang__/de folder
		doc, ok := s.GetDoc("")
		if !ok {
			t.Fatal("expected to find root index")
		}
		if !filepath.IsAbs(doc.FilePath) {
			t.Error("expected FilePath to be absolute")
		}
		expectedPath := filepath.Join(tmpDir, "__lang__", "de", "index.md")
		if doc.FilePath != expectedPath {
			t.Errorf("expected FilePath %q, got %q", expectedPath, doc.FilePath)
		}
	})

	t.Run("non-existent language folder", func(t *testing.T) {
		// When language folder doesn't exist, New() should still succeed
		// This allows GitHub-only mode where language directories might not exist locally
		s, err := New(tmpDir, "github", "", "fr")
		if err != nil {
			t.Errorf("New() should not error for non-existent language folder (GitHub mode support): %v", err)
			return
		}
		// Should use base directory as RootDir when language folder doesn't exist
		if s.RootDir != tmpDir {
			t.Errorf("expected RootDir to be base directory %q when language folder doesn't exist, got %q", tmpDir, s.RootDir)
		}
	})
}

func TestBuildIndexSkipsLangFolder(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files including __lang__ folder
	files := map[string]string{
		"index.md":                   `# Home`,
		"guide/index.md":             `# Guide`,
		"__lang__/de/index.md":       `# Startseite`,
		"__lang__/de/guide/index.md": `# Anleitung`,
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

	// Should find default language files
	if _, ok := s.GetDoc(""); !ok {
		t.Error("expected to find root index")
	}
	if _, ok := s.GetDoc("guide"); !ok {
		t.Error("expected to find guide")
	}

	// Should NOT find __lang__ folder in navigation
	docs := s.ListDocs(false)
	for _, doc := range docs {
		if filepath.Base(filepath.Dir(doc.FilePath)) == "__lang__" {
			t.Errorf("found file from __lang__ folder in default language site: %q", doc.FilePath)
		}
	}
}
