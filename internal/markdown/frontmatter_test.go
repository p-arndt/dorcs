package markdown

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadFrontMatterAndHash(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		content   string
		wantErr   bool
		checkFM   func(*testing.T, FrontMatter)
		checkHash func(*testing.T, string)
	}{
		{
			name: "valid front matter",
			content: `---
title: Test Page
description: A test description
date: 2024-01-15
tags:
  - tutorial
  - test
draft: false
order: 5
---
# Content`,
			wantErr: false,
			checkFM: func(t *testing.T, fm FrontMatter) {
				if fm.Title != "Test Page" {
					t.Errorf("Title = %q; want %q", fm.Title, "Test Page")
				}
				if fm.Description != "A test description" {
					t.Errorf("Description = %q; want %q", fm.Description, "A test description")
				}
				if len(fm.Tags) != 2 {
					t.Errorf("Tags length = %d; want 2", len(fm.Tags))
				}
				if fm.Draft {
					t.Error("Draft should be false")
				}
				if fm.Order != 5 {
					t.Errorf("Order = %d; want 5", fm.Order)
				}
			},
			checkHash: func(t *testing.T, hash string) {
				if hash == "" {
					t.Error("Hash should not be empty")
				}
			},
		},
		{
			name: "no front matter",
			content: `# Just Content
No front matter here.`,
			wantErr: false,
			checkFM: func(t *testing.T, fm FrontMatter) {
				if fm.Title != "" {
					t.Errorf("Title should be empty, got %q", fm.Title)
				}
			},
			checkHash: func(t *testing.T, hash string) {
				if hash == "" {
					t.Error("Hash should not be empty")
				}
			},
		},
		{
			name: "draft true",
			content: `---
title: Draft Page
draft: true
---
# Content`,
			wantErr: false,
			checkFM: func(t *testing.T, fm FrontMatter) {
				if !fm.Draft {
					t.Error("Draft should be true")
				}
			},
			checkHash: func(t *testing.T, hash string) {},
		},
		{
			name: "tags as array",
			content: `---
tags: [tutorial, test, example]
---
# Content`,
			wantErr: false,
			checkFM: func(t *testing.T, fm FrontMatter) {
				if len(fm.Tags) != 3 {
					t.Errorf("Tags length = %d; want 3", len(fm.Tags))
				}
			},
			checkHash: func(t *testing.T, hash string) {},
		},
		{
			name: "order as float",
			content: `---
order: 10.0
---
# Content`,
			wantErr: false,
			checkFM: func(t *testing.T, fm FrontMatter) {
				if fm.Order != 10 {
					t.Errorf("Order = %d; want 10", fm.Order)
				}
			},
			checkHash: func(t *testing.T, hash string) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, "test.md")
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			fm, hash, err := ReadFrontMatterAndHash(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrontMatterAndHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.checkFM(t, fm)
			tt.checkHash(t, hash)
		})
	}
}

func TestReadMarkdownStripFrontMatter(t *testing.T) {
	tmpDir := t.TempDir()

	content := `---
title: Test
date: 2024-01-15
---
# Content Here
Some markdown content.`
	filePath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	raw, fm, hash, modTime, err := ReadMarkdownStripFrontMatter(filePath)
	if err != nil {
		t.Fatalf("ReadMarkdownStripFrontMatter() error = %v", err)
	}

	if raw == "" {
		t.Error("Raw content should not be empty")
	}

	if fm.Title != "Test" {
		t.Errorf("FrontMatter.Title = %q; want %q", fm.Title, "Test")
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if modTime.IsZero() {
		t.Error("ModTime should not be zero")
	}

	// Verify modTime is recent (within last minute)
	if time.Since(modTime) > time.Minute {
		t.Error("ModTime should be recent")
	}
}
