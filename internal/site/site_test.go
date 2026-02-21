package site

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/p-arndt/dorcs/internal/markdown"
)

func TestKeyFromRel(t *testing.T) {
	tests := []struct {
		name     string
		rel      string
		expected string
	}{
		{"root index", "index.md", ""},
		{"simple file", "getting-started.md", "getting-started"},
		{"nested file", "guide/intro.md", "guide/intro"},
		{"folder index", "guide/index.md", "guide"},
		{"deep folder index", "api/v2/index.md", "api/v2"},
		{"deep file", "api/v2/endpoints.md", "api/v2/endpoints"},
		{"uppercase extension", "README.MD", "README"},
		{"mixed case", "Guide/Intro.md", "Guide/Intro"},
		{"empty string", "", ""},
		{"no extension", "readme", ""},
		{"wrong extension", "file.txt", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := keyFromRel(tt.rel)
			if result != tt.expected {
				t.Errorf("keyFromRel(%q) = %q; want %q", tt.rel, result, tt.expected)
			}
		})
	}
}

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"empty", "", ""},
		{"simple", "guide", "guide"},
		{"with slashes", "/guide/", "guide"},
		{"double slashes", "guide//intro", "guide/intro"},
		{"backslashes", "guide\\intro", "guide/intro"},
		{"leading slash", "/guide", "guide"},
		{"trailing slash", "guide/", "guide"},
		{"both slashes", "/guide/intro/", "guide/intro"},
		{"dot prefix", ".hidden", ""},
		{"dot dot", "../escape", ""},
		{"whitespace", "  guide  ", "guide"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeKey(tt.key)
			if result != tt.expected {
				t.Errorf("normalizeKey(%q) = %q; want %q", tt.key, result, tt.expected)
			}
		})
	}
}

func TestTitleFromKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"simple", "getting-started", "Getting started"},
		{"underscores", "api_reference", "Api reference"},
		{"nested", "guide/getting-started", "Getting started"},
		{"empty", "", "Untitled"},
		{"single char", "a", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := titleFromKey(tt.key)
			if result != tt.expected {
				t.Errorf("titleFromKey(%q) = %q; want %q", tt.key, result, tt.expected)
			}
		})
	}
}

func TestDirKeyFromKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"simple", "getting-started", ""},
		{"nested", "guide/getting-started", "guide"},
		{"deep nested", "api/v2/endpoints", "api/v2"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dirKeyFromKey(tt.key)
			if result != tt.expected {
				t.Errorf("dirKeyFromKey(%q) = %q; want %q", tt.key, result, tt.expected)
			}
		})
	}
}

func TestRewriteExtensionlessDocLinks(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		currentDirKey string
		expected      string
	}{
		{
			"simple link",
			"[Getting Started](getting-started)",
			"",
			"[Getting Started](/getting-started)",
		},
		{
			"link with .md extension",
			"[Guide](guide.md)",
			"",
			"[Guide](/guide)",
		},
		{
			"nested link",
			"[API](api/endpoints)",
			"",
			"[API](/api/endpoints)",
		},
		{
			"nested with .md",
			"[API](api/endpoints.md)",
			"",
			"[API](/api/endpoints)",
		},
		{
			"folder index link",
			"[Guide](guide/index.md)",
			"",
			"[Guide](/guide)",
		},
		{
			"absolute URL unchanged",
			"[Example](https://example.com)",
			"",
			"[Example](https://example.com)",
		},
		{
			"anchor link unchanged",
			"[Section](#section)",
			"",
			"[Section](#section)",
		},
		{
			"image link unchanged",
			"![Alt](image.png)",
			"",
			"![Alt](image.png)",
		},
		{
			"multiple links",
			"See [one](one.md) and [two](two.md).",
			"",
			"See [one](/one) and [two](/two).",
		},
		{
			"relative link with ./",
			"[Explain](./explain.md)",
			"user",
			"[Explain](/user/explain)",
		},
		{
			"relative link without extension",
			"[Explain](./explain)",
			"user",
			"[Explain](/user/explain)",
		},
		{
			"relative link from root",
			"[Explain](./explain.md)",
			"",
			"[Explain](/explain)",
		},
		{
			"parent directory link",
			"[Home](../index.md)",
			"user/docs",
			"[Home](/user)",
		},
		{
			"parent directory to sibling",
			"[Other](../other/page.md)",
			"user/docs",
			"[Other](/user/other/page)",
		},
		{
			"multiple parent traversals",
			"[Root](../../index.md)",
			"user/docs/api",
			"[Root](/user)",
		},
		{
			"relative link to index",
			"[Section](./section/index.md)",
			"user",
			"[Section](/user/section)",
		},
		{
			"current directory index",
			"[Home](./index.md)",
			"user",
			"[Home](/user)",
		},
		{
			"relative without ./ prefix",
			"[Troubleshooting](user/troubleshooting.md)",
			"guide",
			"[Troubleshooting](/guide/user/troubleshooting)",
		},
		{
			"relative without ./ prefix from root",
			"[Page](subdir/page.md)",
			"",
			"[Page](/subdir/page)",
		},
		{
			"nested relative without ./",
			"[Docker Proxy](docker-socket-proxy.md)",
			"guide",
			"[Docker Proxy](/guide/docker-socket-proxy)",
		},
		{
			"deep nested without ./",
			"[API](api/v2/endpoints.md)",
			"docs",
			"[API](/docs/api/v2/endpoints)",
		},
		{
			"relative from index.md in subdirectory",
			"[User Guide](./user/index.md)",
			"guide",
			"[User Guide](/guide/user)",
		},
		{
			"relative sibling from index.md",
			"[Getting Started](./getting-started.md)",
			"guide",
			"[Getting Started](/guide/getting-started)",
		},
		{
			"relative nested folder from index.md",
			"[Troubleshooting](user/troubleshooting.md)",
			"guide",
			"[Troubleshooting](/guide/user/troubleshooting)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := markdown.RewriteExtensionlessDocLinks(tt.input, tt.currentDirKey, "", "", "", "", "")
			if result != tt.expected {
				t.Errorf("RewriteExtensionlessDocLinks(%q, %q, %q) = %q; want %q", tt.input, tt.currentDirKey, "", result, tt.expected)
			}
		})
	}
}

func TestStripYAMLFrontMatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"with front matter",
			"---\ntitle: Test\n---\n# Content",
			"# Content",
		},
		{
			"no front matter",
			"# Just Content",
			"# Just Content",
		},
		{
			"empty front matter",
			"---\n---\n# Content",
			"# Content",
		},
		{
			"unclosed front matter",
			"---\ntitle: Test\n# Content",
			"---\ntitle: Test\n# Content",
		},
		{
			"front matter with CRLF",
			"---\r\ntitle: Test\r\n---\r\n# Content",
			"# Content",
		},
		{
			"with BOM",
			"\uFEFF---\ntitle: Test\n---\n# Content",
			"# Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := markdown.StripYAMLFrontMatter(tt.input)
			if result != tt.expected {
				t.Errorf("stripYAMLFrontMatter() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestSiteNew(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "site-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("valid directory", func(t *testing.T) {
		s, err := New(tmpDir, "github", "", "")
		if err != nil {
			t.Errorf("New(%q) returned error: %v", tmpDir, err)
		}
		if s == nil {
			t.Error("New() returned nil site")
		}
	})

	t.Run("empty path", func(t *testing.T) {
		_, err := New("", "github", "", "")
		if err == nil {
			t.Error("New(\"\") should return error")
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		_, err := New("/non/existent/path", "github", "", "")
		if err == nil {
			t.Error("New() with non-existent path should return error")
		}
	})

	t.Run("file instead of directory", func(t *testing.T) {
		tmpFile := filepath.Join(tmpDir, "file.txt")
		if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		_, err := New(tmpFile, "github", "", "")
		if err == nil {
			t.Error("New() with file path should return error")
		}
	})
}

func TestSiteBuildIndex(t *testing.T) {
	// Create a temp directory with test markdown files
	tmpDir, err := os.MkdirTemp("", "site-index-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"index.md": `---
title: Home
---
# Welcome`,
		"getting-started.md": `---
title: Getting Started
description: Learn the basics
---
# Getting Started`,
		"guide/index.md": `---
title: Guide
---
# Guide`,
		"guide/intro.md": `---
title: Introduction
draft: true
---
# Intro`,
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

	t.Run("root index exists", func(t *testing.T) {
		doc, ok := s.GetDoc("")
		if !ok {
			t.Error("root index not found")
		}
		if doc.Title != "Home" {
			t.Errorf("root index title = %q; want %q", doc.Title, "Home")
		}
	})

	t.Run("simple doc exists", func(t *testing.T) {
		doc, ok := s.GetDoc("getting-started")
		if !ok {
			t.Error("getting-started not found")
		}
		if doc.Title != "Getting Started" {
			t.Errorf("title = %q; want %q", doc.Title, "Getting Started")
		}
		if doc.Description != "Learn the basics" {
			t.Errorf("description = %q; want %q", doc.Description, "Learn the basics")
		}
	})

	t.Run("folder index exists", func(t *testing.T) {
		doc, ok := s.GetDoc("guide")
		if !ok {
			t.Error("guide index not found")
		}
		if doc.Title != "Guide" {
			t.Errorf("title = %q; want %q", doc.Title, "Guide")
		}
	})

	t.Run("nested doc exists", func(t *testing.T) {
		doc, ok := s.GetDoc("guide/intro")
		if !ok {
			t.Error("guide/intro not found")
		}
		if !doc.Draft {
			t.Error("guide/intro should be marked as draft")
		}
	})

	t.Run("ListDocs excludes drafts", func(t *testing.T) {
		docs := s.ListDocs(false)
		for _, d := range docs {
			if d.Draft {
				t.Errorf("ListDocs(false) returned draft doc: %q", d.Key)
			}
		}
	})

	t.Run("ListDocs includes drafts", func(t *testing.T) {
		docs := s.ListDocs(true)
		hasDraft := false
		for _, d := range docs {
			if d.Draft {
				hasDraft = true
				break
			}
		}
		if !hasDraft {
			t.Error("ListDocs(true) should include draft docs")
		}
	})
}

func TestSiteRenderDoc(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "site-render-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file with headings for TOC
	content := `---
title: Test Page
---
# Main Title

## First Section

Some content here.

### Subsection

More content.

## Second Section

Final content.
`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.md"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	s, err := New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if err := s.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex() failed: %v", err)
	}

	t.Run("renders HTML", func(t *testing.T) {
		rendered, err := s.RenderDoc("test")
		if err != nil {
			t.Fatalf("RenderDoc() failed: %v", err)
		}
		if rendered.HTML == "" {
			t.Error("HTML should not be empty")
		}
		if rendered.Doc.Title != "Test Page" {
			t.Errorf("title = %q; want %q", rendered.Doc.Title, "Test Page")
		}
	})

	t.Run("generates TOC", func(t *testing.T) {
		rendered, err := s.RenderDoc("test")
		if err != nil {
			t.Fatalf("RenderDoc() failed: %v", err)
		}
		if rendered.TocHTML == "" {
			t.Error("TocHTML should not be empty")
		}
	})

	t.Run("non-existent doc", func(t *testing.T) {
		_, err := s.RenderDoc("non-existent")
		if err == nil {
			t.Error("RenderDoc() should return error for non-existent doc")
		}
	})
}

func TestPreprocessMarkdownGitHubRelativePaths(t *testing.T) {
	tests := []struct {
		name          string
		markdown      string
		doc           *Doc
		basePath      string
		language      string
		version       string
		useGitHubAuth bool
		expected      string
		description   string
	}{
		{
			name:          "GitHub root index with relative image path (with GitHub client)",
			markdown:      "![Logo](./logo.png)",
			doc:           &Doc{Key: "", RelPath: "index.md", DirKey: "", IsGitHub: true, GitHubPath: "docs/index.md"},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: true,
			expected:      "![Logo](https://raw.githubusercontent.com/owner/repo/main/docs/logo.png)",
			description:   "With GitHub client, relative image resolves to raw.githubusercontent.com",
		},
		{
			name:          "GitHub root index with relative image path (no GitHub client)",
			markdown:      "![Logo](../logo.png)",
			doc:           &Doc{Key: "", RelPath: "index.md", DirKey: "", IsGitHub: true, GitHubPath: "docs/index.md"},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: false,
			expected:      "![Logo](/logo.png)",
			description:   "Without GitHub client, falls back to local path resolution",
		},
		{
			name:          "GitHub nested index with relative image path",
			markdown:      "![Guide](../images/guide.png)",
			doc:           &Doc{Key: "guide", RelPath: "guide/index.md", DirKey: "guide", IsGitHub: true, GitHubPath: "docs/guide/index.md"},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: false,
			expected:      "![Guide](/images/guide.png)",
			description:   "Image at sibling directory should resolve correctly",
		},
		{
			name:          "GitHub document with local relative image path (with GitHub client)",
			markdown:      "![Screenshot](./screenshot.png)",
			doc:           &Doc{Key: "guide/intro", RelPath: "guide/intro.md", DirKey: "guide", IsGitHub: true, GitHubPath: "docs/guide/intro.md"},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: true,
			expected:      "![Screenshot](https://raw.githubusercontent.com/owner/repo/main/docs/guide/screenshot.png)",
			description:   "With GitHub client, same-dir image resolves to raw URL",
		},
		{
			name:          "GitHub document with parent directory reference (with GitHub client)",
			markdown:      "![Parent](../../assets/logo.png)",
			doc:           &Doc{Key: "api/v2/endpoints", RelPath: "api/v2/endpoints.md", DirKey: "api/v2", IsGitHub: true, GitHubPath: "docs/api/v2/endpoints.md"},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: true,
			expected:      "![Parent](https://raw.githubusercontent.com/owner/repo/main/docs/assets/logo.png)",
			description:   "Parent ref resolves correctly to raw URL",
		},
		{
			name:          "GitHub doc with basePath set (with GitHub client)",
			markdown:      "![Logo](../logo.png)",
			doc:           &Doc{Key: "", RelPath: "index.md", DirKey: "", IsGitHub: true, GitHubPath: "docs/index.md"},
			basePath:      "/myapp",
			language:      "",
			version:       "",
			useGitHubAuth: true,
			expected:      "![Logo](https://raw.githubusercontent.com/owner/repo/main/logo.png)",
			description:   "With GitHub client, basePath does not apply to raw URLs",
		},
		{
			name:          "GitHub doc with basePath set (no GitHub client)",
			markdown:      "![Logo](../logo.png)",
			doc:           &Doc{Key: "", RelPath: "index.md", DirKey: "", IsGitHub: true, GitHubPath: "docs/index.md"},
			basePath:      "/myapp",
			language:      "",
			version:       "",
			useGitHubAuth: false,
			expected:      "![Logo](/myapp/logo.png)",
			description:   "Without GitHub client, image paths should include basePath",
		},
		{
			name:          "GitHub doc with absolute URL unchanged",
			markdown:      "![External](https://example.com/logo.png)",
			doc:           &Doc{Key: "", RelPath: "index.md", DirKey: "", IsGitHub: true, GitHubPath: "docs/index.md"},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: true,
			expected:      "![External](https://example.com/logo.png)",
			description:   "Absolute URLs should not be modified",
		},
		{
			name:          "GitHub doc with already absolute path",
			markdown:      "![Logo](/logo.png)",
			doc:           &Doc{Key: "", RelPath: "index.md", DirKey: "", IsGitHub: true, GitHubPath: "docs/index.md"},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: true,
			expected:      "![Logo](/logo.png)",
			description:   "Already absolute paths should not be modified",
		},
		{
			name:          "Local (non-GitHub) document uses mapped key",
			markdown:      "![Logo](../logo.png)",
			doc:           &Doc{Key: "guide/intro", RelPath: "guide/intro.md", DirKey: "guide", IsGitHub: false, GitHubPath: ""},
			basePath:      "",
			language:      "",
			version:       "",
			useGitHubAuth: false,
			expected:      "![Logo](/logo.png)",
			description:   "Local docs should use their mapped key, not GitHubPath",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Site{
				BasePath:        tt.basePath,
				Language:        tt.language,
				Version:         tt.version,
				DefaultVersion:  tt.version,
				DefaultLanguage: tt.language,
			}
			if tt.useGitHubAuth {
				s.SetGitHubConfig(newMockGitHubClient(), "owner", "repo", "main", "docs")
			}
			result := s.preprocessMarkdown(tt.markdown, tt.doc)
			if result != tt.expected {
				t.Errorf("%s\npreprocessMarkdown(%q, %+v) = %q\nwant %q",
					tt.description, tt.markdown, tt.doc, result, tt.expected)
			}
		})
	}
}

// TestPreprocessMarkdownGitHubRelativeDocLinks ensures that relative doc links (e.g. ./xxx.md)
// in GitHub docs resolve to the correct path based on the doc's Key/DirKey, NOT doc.GitHubPath.
// Bug: GitHubPath includes repo root (e.g. "docs/") which wrongly resolved ./xxx.md to /docs/xxx
// instead of /external-content/xxx when the doc is in docs/external-content/.
func TestPreprocessMarkdownGitHubRelativeDocLinks(t *testing.T) {
	s := &Site{BasePath: "/docs", DefaultVersion: "", DefaultLanguage: ""}
	s.SetGitHubConfig(newMockGitHubClient(), "owner", "repo", "main", "docs")

	// Doc at docs/external-content/index.md: Key="external-content", GitHubPath="docs/external-content/index.md"
	// Relative link ./github.md must resolve to /docs/external-content/github, NOT /docs/github
	md := "[GitHub Guide](./github.md)"
	doc := &Doc{
		Key:        "external-content",
		RelPath:    "external-content/index.md",
		DirKey:     "external-content",
		IsGitHub:   true,
		GitHubPath: "docs/external-content/index.md",
	}
	result := s.preprocessMarkdown(md, doc)
	expected := "[GitHub Guide](/docs/external-content/github)"
	if result != expected {
		t.Errorf("GitHub doc in subfolder: preprocessMarkdown(%q) = %q, want %q", md, result, expected)
	}
}
