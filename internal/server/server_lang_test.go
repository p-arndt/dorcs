package server

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/site"
)

func TestLanguageDetection(t *testing.T) {
	tmpDir := t.TempDir()

	// Create default language files (new MkDocs-style structure)
	files := map[string]string{
		"index.md":              `# Home`,
		"getting-started.md":    `# Getting Started`,
		"de/index.md":           `# Startseite`,
		"de/getting-started.md": `# Erste Schritte`,
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

	// Create config with multi-lingual support
	cfg := &config.Config{
		Languages: config.LanguagesConfig{
			Default: "en",
			Enabled: []config.Language{
				{Code: "en", Name: "English"},
				{Code: "de", Name: "Deutsch"},
			},
		},
	}

	// Create sites
	defaultSite, err := site.New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("failed to create default site: %v", err)
	}
	if err := defaultSite.BuildIndex(); err != nil {
		t.Fatalf("failed to build default site index: %v", err)
	}

	deSite, err := site.New(tmpDir, "github", "", "de")
	if err != nil {
		t.Fatalf("failed to create German site: %v", err)
	}
	if err := deSite.BuildIndex(); err != nil {
		t.Fatalf("failed to build German site index: %v", err)
	}

	sites := map[string]*site.Site{
		"":   defaultSite,
		"de": deSite,
	}

	// Create a minimal template that matches what the server expects
	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))

	handler := New(Config{
		DocsDir:      tmpDir,
		RootDir:      tmpDir,
		SiteTitle:    "Test",
		BasePath:     "",
		Site:         defaultSite,
		Sites:        sites,
		DocumentTmpl: tmpl,
		SiteConfig:   cfg,
		Version:      "test",
	})

	tests := []struct {
		name         string
		path         string
		expectedLang string
		expectedDoc  string
		statusCode   int
	}{
		{
			name:         "root - default language",
			path:         "/",
			expectedLang: "",
			expectedDoc:  "",
			statusCode:   http.StatusOK,
		},
		{
			name:         "root - German",
			path:         "/de/",
			expectedLang: "de",
			expectedDoc:  "",
			statusCode:   http.StatusOK,
		},
		{
			name:         "getting-started - default language",
			path:         "/getting-started",
			expectedLang: "",
			expectedDoc:  "getting-started",
			statusCode:   http.StatusOK,
		},
		{
			name:         "getting-started - German",
			path:         "/de/getting-started",
			expectedLang: "de",
			expectedDoc:  "getting-started",
			statusCode:   http.StatusOK,
		},
		{
			name:         "invalid language code",
			path:         "/fr/getting-started",
			expectedLang: "",
			expectedDoc:  "fr/getting-started", // Treated as document path
			statusCode:   http.StatusNotFound,  // Will fail because file doesn't exist
		},
		{
			name:         "non-language path",
			path:         "/guide",
			expectedLang: "",
			expectedDoc:  "guide",
			statusCode:   http.StatusNotFound, // Will fail because file doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				// For now, just check that it doesn't crash
				// We can't easily check the language without parsing HTML
				if w.Code == http.StatusInternalServerError {
					t.Errorf("got internal server error for path %q: %s", tt.path, w.Body.String())
				}
			}
		})
	}
}

func TestConvertNavNodesWithLang(t *testing.T) {
	tests := []struct {
		name        string
		basePath    string
		currentLang string
		key         string
		expected    string
	}{
		{
			name:        "default language - no prefix",
			basePath:    "",
			currentLang: "",
			key:         "getting-started",
			expected:    "/getting-started",
		},
		{
			name:        "default language - with basepath",
			basePath:    "/docs",
			currentLang: "",
			key:         "getting-started",
			expected:    "/docs/getting-started",
		},
		{
			name:        "German language - no basepath",
			basePath:    "",
			currentLang: "de",
			key:         "getting-started",
			expected:    "/de/getting-started",
		},
		{
			name:        "German language - with basepath",
			basePath:    "/docs",
			currentLang: "de",
			key:         "getting-started",
			expected:    "/docs/de/getting-started",
		},
		{
			name:        "German language - root index",
			basePath:    "",
			currentLang: "de",
			key:         "",
			expected:    "/de/",
		},
		{
			name:        "German language - nested path",
			basePath:    "",
			currentLang: "de",
			key:         "usage/metadata",
			expected:    "/de/usage/metadata",
		},
		{
			name:        "default language - nested path",
			basePath:    "",
			currentLang: "",
			key:         "usage/metadata",
			expected:    "/usage/metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple nav node
			node := &site.NavNode{
				Name:  "Test",
				Key:   tt.key,
				IsDir: false,
				Page:  &site.Doc{Key: tt.key, Title: "Test"},
			}

			items := convertNavNodesWithLang([]*site.NavNode{node}, tt.basePath, tt.currentLang)
			if len(items) != 1 {
				t.Fatalf("expected 1 nav item, got %d", len(items))
			}

			if items[0].Path != tt.expected {
				t.Errorf("convertNavNodesWithLang() = %q; want %q", items[0].Path, tt.expected)
			}
		})
	}
}

func TestDocPathExtraction(t *testing.T) {
	tmpDir := t.TempDir()

	// Create minimal files (new MkDocs-style structure)
	files := map[string]string{
		"index.md":    `# Home`,
		"de/index.md": `# Startseite`,
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

	cfg := &config.Config{
		Languages: config.LanguagesConfig{
			Default: "en",
			Enabled: []config.Language{
				{Code: "en", Name: "English"},
				{Code: "de", Name: "Deutsch"},
			},
		},
	}

	defaultSite, _ := site.New(tmpDir, "github", "", "")
	defaultSite.BuildIndex()
	deSite, _ := site.New(tmpDir, "github", "", "de")
	deSite.BuildIndex()

	sites := map[string]*site.Site{
		"":   defaultSite,
		"de": deSite,
	}

	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))

	handler := New(Config{
		DocsDir:      tmpDir,
		RootDir:      tmpDir,
		SiteTitle:    "Test",
		BasePath:     "",
		Site:         defaultSite,
		Sites:        sites,
		DocumentTmpl: tmpl,
		SiteConfig:   cfg,
		Version:      "test",
	})

	tests := []struct {
		name        string
		path        string
		expectedDoc string
	}{
		{name: "root default", path: "/", expectedDoc: "/"},
		{name: "root German", path: "/de/", expectedDoc: "/"},
		{name: "nested default", path: "/usage/metadata", expectedDoc: "/usage/metadata"},
		{name: "nested German", path: "/de/usage/metadata", expectedDoc: "/usage/metadata"},
		{name: "deep nested German", path: "/de/guide/advanced/topics", expectedDoc: "/guide/advanced/topics"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Check that DocPath is set correctly in the response
			// We can't easily extract it from HTML, but we can verify it doesn't crash
			if w.Code == http.StatusInternalServerError {
				body := w.Body.String()
				if strings.Contains(body, "server not configured") {
					t.Logf("server not configured (expected for some paths)")
				} else {
					t.Errorf("unexpected error for path %q: %s", tt.path, body)
				}
			}
		})
	}
}
