package server

import (
	"sync"
	"testing"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/site"
)

type navTestSiteConfig struct {
	defaultLang    string
	defaultVersion string
	multiLang      bool
	multiVersion   bool
}

func (c navTestSiteConfig) IsMultiVersion() bool       { return c.multiVersion }
func (c navTestSiteConfig) GetDefaultVersion() string  { return c.defaultVersion }
func (c navTestSiteConfig) IsMultiLingual() bool       { return c.multiLang }
func (c navTestSiteConfig) GetDefaultLanguage() string { return c.defaultLang }

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

func TestConvertNavNodesWithVersionAndLangOrder(t *testing.T) {
	nodes := []*site.NavNode{
		{
			Name:  "Guide",
			Key:   "guide",
			IsDir: false,
			Page:  &site.Doc{Key: "guide", Title: "Guide"},
		},
	}

	cfg := navTestSiteConfig{
		defaultLang:    "en",
		defaultVersion: "latest",
		multiLang:      true,
		multiVersion:   true,
	}

	items := convertNavNodesWithVersionAndLang(nodes, "/docs", "v1", "de", cfg)
	if len(items) != 1 {
		t.Fatalf("expected 1 nav item, got %d", len(items))
	}
	if items[0].Path != "/docs/de/v1/guide" {
		t.Fatalf("path = %q, want %q", items[0].Path, "/docs/de/v1/guide")
	}
}

func TestComputeEditOnGitHubURLForDefaultVersion(t *testing.T) {
	cfg := &config.Config{
		Versions: config.VersionsConfig{
			Default: "latest",
			Enabled: []config.Version{
				{ID: "latest", Name: "Latest"},
				{ID: "v1", Name: "Version 1"},
			},
		},
		GitHub: config.GitHubConfig{
			EditOnGitHub: config.EditOnGitHubConfig{
				Repository: "https://github.com/example/repo/tree/main/docs",
			},
		},
	}

	doc := &site.Doc{RelPath: "guide.md"}

	got := computeEditOnGitHubURL(doc, cfg, "", "latest")
	want := "https://github.com/example/repo/edit/main/docs/guide.md"
	if got != want {
		t.Fatalf("default version edit URL = %q, want %q", got, want)
	}
}

func TestComputeEditOnGitHubURLForLanguageAndVersion(t *testing.T) {
	cfg := &config.Config{
		Languages: config.LanguagesConfig{
			Default: "en",
			Enabled: []config.Language{
				{Code: "en", Name: "English"},
				{Code: "de", Name: "Deutsch"},
			},
		},
		Versions: config.VersionsConfig{
			Default: "latest",
			Enabled: []config.Version{
				{ID: "latest", Name: "Latest"},
				{ID: "v1", Name: "Version 1"},
			},
		},
		GitHub: config.GitHubConfig{
			EditOnGitHub: config.EditOnGitHubConfig{
				Repository: "https://github.com/example/repo/tree/main/docs",
			},
		},
	}

	doc := &site.Doc{RelPath: "guide.md"}

	got := computeEditOnGitHubURL(doc, cfg, "de", "v1")
	want := "https://github.com/example/repo/edit/main/docs/de/v1/guide.md"
	if got != want {
		t.Fatalf("language/version edit URL = %q, want %q", got, want)
	}
}

func TestFirstNavItemPath(t *testing.T) {
	tests := []struct {
		name  string
		items []NavItem
		want  string
	}{
		{
			name:  "empty items",
			items: nil,
			want:  "",
		},
		{
			name: "direct path",
			items: []NavItem{
				{Title: "Home", Path: "/getting-started"},
			},
			want: "/getting-started",
		},
		{
			name: "first item has no path, child does",
			items: []NavItem{
				{Title: "Folder", IsDir: true, Children: []NavItem{
					{Title: "Child", Path: "/folder/child"},
				}},
			},
			want: "/folder/child",
		},
		{
			name: "skips empty paths",
			items: []NavItem{
				{Title: "No Path"},
				{Title: "Has Path", Path: "/second"},
			},
			want: "/second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstNavItemPath(tt.items)
			if got != tt.want {
				t.Errorf("firstNavItemPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNavItemsContainPath(t *testing.T) {
	items := []NavItem{
		{Title: "Home", Path: "/"},
		{Title: "Guide", Path: "/guide", Children: []NavItem{
			{Title: "Intro", Path: "/guide/intro"},
			{Title: "Advanced", Path: "/guide/advanced"},
		}},
		{Title: "API", Path: "/api"},
	}

	tests := []struct {
		path string
		want bool
	}{
		{"/", true},
		{"/guide", true},
		{"/guide/intro", true},
		{"/guide/advanced", true},
		{"/api", true},
		{"/unknown", false},
		{"/guide/unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := navItemsContainPath(items, tt.path)
			if got != tt.want {
				t.Errorf("navItemsContainPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestBuildSectionTabsNilConfig(t *testing.T) {
	handler := New(Config{
		BasePath:   "",
		SiteConfig: nil,
	})

	tabs, activeIdx := handler.buildSectionTabs(nil, "/")
	if tabs != nil {
		t.Errorf("expected nil tabs for nil config, got %v", tabs)
	}
	if activeIdx != -1 {
		t.Errorf("expected activeIdx -1, got %d", activeIdx)
	}
}

func TestBuildSectionTabsNoSections(t *testing.T) {
	cfg := config.Default()
	handler := New(Config{
		BasePath:   "",
		SiteConfig: cfg,
	})

	tabs, activeIdx := handler.buildSectionTabs(nil, "/")
	if tabs != nil {
		t.Errorf("expected nil tabs when no sections configured, got %v", tabs)
	}
	if activeIdx != -1 {
		t.Errorf("expected activeIdx -1, got %d", activeIdx)
	}
}
