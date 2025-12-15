package server

import (
	"sync"
	"testing"

	"dorcs-v2/internal/site"
)

func TestBuildNavItemsWithSite(t *testing.T) {
	handler := New(Config{
		BasePath:  "",
		HideDraft: false,
	})

	tests := []struct {
		name       string
		targetSite *site.Site
		wantNil    bool
	}{
		{"nil site", nil, true},
		// Note: We can't easily test with a real site without setting up files,
		// but we can test the nil case and basic structure
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler.buildNavItemsWithSite(tt.targetSite)
			if tt.wantNil && got != nil {
				t.Errorf("buildNavItemsWithSite() = %v, want nil", got)
			}
		})
	}
}

func TestGetRootTitleWithSite(t *testing.T) {
	handler := New(Config{
		HideDraft: false,
	})

	tests := []struct {
		name       string
		targetSite *site.Site
		want       string
	}{
		{"nil site", nil, "Home"},
		// Note: Testing with real site would require file setup
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler.getRootTitleWithSite(tt.targetSite)
			if got != tt.want {
				t.Errorf("getRootTitleWithSite() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertNavNodesWithLangDetailed(t *testing.T) {
	tests := []struct {
		name        string
		nodes       []*site.NavNode
		basePath    string
		currentLang string
		wantCount   int
		wantPath    string
	}{
		{
			name:        "empty nodes",
			nodes:       nil,
			basePath:    "",
			currentLang: "",
			wantCount:   0,
		},
		{
			name: "single node - default lang",
			nodes: []*site.NavNode{
				{
					Name:  "Getting Started",
					Key:   "getting-started",
					IsDir: false,
					Page:  &site.Doc{Key: "getting-started", Title: "Getting Started"},
				},
			},
			basePath:    "",
			currentLang: "",
			wantCount:   1,
			wantPath:    "/getting-started",
		},
		{
			name: "single node - with base path",
			nodes: []*site.NavNode{
				{
					Name:  "Getting Started",
					Key:   "getting-started",
					IsDir: false,
					Page:  &site.Doc{Key: "getting-started", Title: "Getting Started"},
				},
			},
			basePath:    "/docs",
			currentLang: "",
			wantCount:   1,
			wantPath:    "/docs/getting-started",
		},
		{
			name: "single node - with language",
			nodes: []*site.NavNode{
				{
					Name:  "Erste Schritte",
					Key:   "getting-started",
					IsDir: false,
					Page:  &site.Doc{Key: "getting-started", Title: "Erste Schritte"},
				},
			},
			basePath:    "",
			currentLang: "de",
			wantCount:   1,
			wantPath:    "/de/getting-started",
		},
		{
			name: "single node - with base path and language",
			nodes: []*site.NavNode{
				{
					Name:  "Erste Schritte",
					Key:   "getting-started",
					IsDir: false,
					Page:  &site.Doc{Key: "getting-started", Title: "Erste Schritte"},
				},
			},
			basePath:    "/docs",
			currentLang: "de",
			wantCount:   1,
			wantPath:    "/docs/de/getting-started",
		},
		{
			name: "directory without page",
			nodes: []*site.NavNode{
				{
					Name:  "Guide",
					Key:   "guide",
					IsDir: true,
					Page:  nil,
				},
			},
			basePath:    "",
			currentLang: "",
			wantCount:   1,
			wantPath:    "", // Directories without index.md have empty path
		},
		{
			name: "directory with page",
			nodes: []*site.NavNode{
				{
					Name:  "Guide",
					Key:   "guide",
					IsDir: true,
					Page:  &site.Doc{Key: "guide", Title: "Guide"},
				},
			},
			basePath:    "",
			currentLang: "",
			wantCount:   1,
			wantPath:    "/guide",
		},
		{
			name: "nested nodes",
			nodes: []*site.NavNode{
				{
					Name:  "Guide",
					Key:   "guide",
					IsDir: true,
					Page:  &site.Doc{Key: "guide", Title: "Guide"},
					Children: []*site.NavNode{
						{
							Name:  "Advanced",
							Key:   "guide/advanced",
							IsDir: false,
							Page:  &site.Doc{Key: "guide/advanced", Title: "Advanced"},
						},
					},
				},
			},
			basePath:    "",
			currentLang: "",
			wantCount:   1,
		},
		{
			name: "prefer Page.Title over Name",
			nodes: []*site.NavNode{
				{
					Name:  "Old Name",
					Key:   "page",
					IsDir: false,
					Page:  &site.Doc{Key: "page", Title: "New Title"},
				},
			},
			basePath:    "",
			currentLang: "",
			wantCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertNavNodesWithLang(tt.nodes, tt.basePath, tt.currentLang)

			if len(got) != tt.wantCount {
				t.Errorf("convertNavNodesWithLang() returned %d items, want %d", len(got), tt.wantCount)
				return
			}

			if tt.wantCount > 0 && tt.wantPath != "" {
				if got[0].Path != tt.wantPath {
					t.Errorf("convertNavNodesWithLang() path = %q, want %q", got[0].Path, tt.wantPath)
				}
			}

			if tt.wantCount > 0 {
				// Verify structure
				if got[0].Title == "" {
					t.Error("convertNavNodesWithLang() item missing title")
				}

				// For nested nodes, check children
				if len(tt.nodes) > 0 && len(tt.nodes[0].Children) > 0 {
					if len(got[0].Children) == 0 {
						t.Error("convertNavNodesWithLang() missing children")
					}
				}

				// Check Page.Title preference
				if len(tt.nodes) > 0 && tt.nodes[0].Page != nil && tt.nodes[0].Page.Title != "" {
					if got[0].Title != tt.nodes[0].Page.Title {
						t.Errorf("convertNavNodesWithLang() title = %q, want %q (should prefer Page.Title)", got[0].Title, tt.nodes[0].Page.Title)
					}
				}
			}
		})
	}
}

func TestConvertNavNodesWithLangConcurrency(t *testing.T) {
	// Test that convertNavNodesWithLang is safe for concurrent use
	nodes := []*site.NavNode{
		{
			Name:  "Test",
			Key:   "test",
			IsDir: false,
			Page:  &site.Doc{Key: "test", Title: "Test"},
		},
	}

	var wg sync.WaitGroup
	concurrency := 10
	iterations := 100

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				result := convertNavNodesWithLang(nodes, "", "")
				if len(result) != 1 {
					t.Errorf("convertNavNodesWithLang() returned %d items, want 1", len(result))
				}
			}
		}()
	}

	wg.Wait()
}
