package server

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dorcs-v2/internal/config"
	"dorcs-v2/internal/site"
)

func TestNew(t *testing.T) {
	cfg := Config{
		DocsDir:   "/tmp/docs",
		RootDir:   "/tmp/root",
		SiteTitle: "Test Site",
		BasePath:  "/docs",
	}

	handler := New(cfg)

	if handler == nil {
		t.Fatal("New() returned nil")
	}

	if handler.cfg.DocsDir != cfg.DocsDir {
		t.Errorf("New() DocsDir = %q, want %q", handler.cfg.DocsDir, cfg.DocsDir)
	}
	if handler.cfg.SiteTitle != cfg.SiteTitle {
		t.Errorf("New() SiteTitle = %q, want %q", handler.cfg.SiteTitle, cfg.SiteTitle)
	}
}

func TestUpdateConfig(t *testing.T) {
	handler := New(Config{
		SiteTitle: "Old Title",
	})

	newConfig := &config.Config{
		Site: config.SiteConfig{
			Title: "New Title",
		},
	}

	handler.UpdateConfig(newConfig)

	if handler.cfg.SiteTitle != "New Title" {
		t.Errorf("UpdateConfig() SiteTitle = %q, want %q", handler.cfg.SiteTitle, "New Title")
	}

	if handler.cfg.SiteConfig != newConfig {
		t.Error("UpdateConfig() SiteConfig not updated")
	}
}

func TestUpdateConfigEmptyTitle(t *testing.T) {
	handler := New(Config{
		SiteTitle: "Original Title",
	})

	newConfig := &config.Config{
		Site: config.SiteConfig{
			Title: "", // Empty title should not override
		},
	}

	handler.UpdateConfig(newConfig)

	// Title should remain unchanged if new title is empty
	if handler.cfg.SiteTitle != "Original Title" {
		t.Errorf("UpdateConfig() should not change title when empty, got %q", handler.cfg.SiteTitle)
	}
}

func TestServeHTTPMethodNotAllowed(t *testing.T) {
	handler := New(Config{})

	tests := []struct {
		name   string
		method string
	}{
		{"POST", "POST"},
		{"PUT", "PUT"},
		{"DELETE", "DELETE"},
		{"PATCH", "PATCH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("ServeHTTP() status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestServeHTTPAllowedMethods(t *testing.T) {
	tmpDir := t.TempDir()

	// Create minimal site
	files := map[string]string{
		"index.md": `# Home`,
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	testSite, _ := site.New(tmpDir, "github", "", "")
	testSite.BuildIndex()

	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))
	handler := New(Config{
		DocsDir:      tmpDir,
		RootDir:      tmpDir,
		Site:         testSite,
		DocumentTmpl: tmpl,
	})

	tests := []struct {
		name   string
		method string
		wantOK bool
	}{
		{"GET", "GET", true},
		{"HEAD", "HEAD", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if tt.wantOK && w.Code == http.StatusMethodNotAllowed {
				t.Errorf("ServeHTTP() method %s should be allowed, got %d", tt.method, w.Code)
			}
		})
	}
}

func TestServeHTTPBasePath(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"index.md": `# Home`,
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		dir := filepath.Dir(fullPath)
		os.MkdirAll(dir, 0755)
		os.WriteFile(fullPath, []byte(content), 0644)
	}

	testSite, _ := site.New(tmpDir, "github", "", "")
	testSite.BuildIndex()

	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))
	handler := New(Config{
		DocsDir:      tmpDir,
		RootDir:      tmpDir,
		BasePath:     "/docs",
		Site:         testSite,
		DocumentTmpl: tmpl,
	})

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"correct base path", "/docs/", http.StatusOK},
		{"wrong base path", "/", http.StatusNotFound},
		{"base path with doc", "/docs/getting-started", http.StatusNotFound}, // File doesn't exist
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ServeHTTP() path %q status = %d, want %d", tt.path, w.Code, tt.wantStatus)
			}
		})
	}
}

func TestServeHTTPLanguageDetection(t *testing.T) {
	tmpDir := t.TempDir()

	files := map[string]string{
		"index.md":                       `# Home`,
		"__lang__/de/index.md":           `# Startseite`,
		"__lang__/de/getting-started.md": `# Erste Schritte`,
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		dir := filepath.Dir(fullPath)
		os.MkdirAll(dir, 0755)
		os.WriteFile(fullPath, []byte(content), 0644)
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
		Site:         defaultSite,
		Sites:        sites,
		DocumentTmpl: tmpl,
		SiteConfig:   cfg,
	})

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"default language root", "/", http.StatusOK},
		{"German language root", "/de/", http.StatusOK},
		{"invalid language", "/fr/", http.StatusNotFound}, // Treated as document path
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ServeHTTP() path %q status = %d, want %d", tt.path, w.Code, tt.wantStatus)
			}
		})
	}
}

func TestServeHTTPStaticAsset(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := tmpDir

	// Create a test image file
	testImage := filepath.Join(docsDir, "logo.png")
	os.WriteFile(testImage, []byte("fake png content"), 0644)

	testSite, _ := site.New(tmpDir, "github", "", "")
	testSite.BuildIndex()

	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))
	handler := New(Config{
		DocsDir:      docsDir,
		RootDir:      tmpDir,
		Site:         testSite,
		DocumentTmpl: tmpl,
	})

	req := httptest.NewRequest("GET", "/logo.png", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ServeHTTP() static asset status = %d, want %d", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "image/png") {
		t.Errorf("ServeHTTP() static asset Content-Type = %q, want image/png", contentType)
	}
}

func TestServeHTTPNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	testSite, _ := site.New(tmpDir, "github", "", "")
	testSite.BuildIndex()

	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))
	handler := New(Config{
		DocsDir:      tmpDir,
		RootDir:      tmpDir,
		Site:         testSite,
		DocumentTmpl: tmpl,
	})

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("ServeHTTP() nonexistent path status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
