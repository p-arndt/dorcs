package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/p-arndt/dorcs/internal/site"
)

func setupTestSite(t *testing.T) (*site.Site, *Handler) {
	tmpDir := t.TempDir()

	// Create test markdown files
	files := map[string]string{
		"index.md":           `# Home Page`,
		"getting-started.md": `# Getting Started`,
		"guide/index.md":     `# Guide`,
		"guide/advanced.md":  `# Advanced Topics`,
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

	testSite, err := site.New(tmpDir, "github", "", "")
	if err != nil {
		t.Fatalf("failed to create site: %v", err)
	}
	if err := testSite.BuildIndex(); err != nil {
		t.Fatalf("failed to build index: %v", err)
	}

	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))
	handler := New(Config{
		DocsDir:      tmpDir,
		RootDir:      tmpDir,
		SiteTitle:    "Test Site",
		BasePath:     "",
		Site:         testSite,
		DocumentTmpl: tmpl,
		HideDraft:    false,
	})

	return testSite, handler
}

func TestServeSearch(t *testing.T) {
	_, handler := setupTestSite(t)

	tests := []struct {
		name        string
		method      string
		query       string
		wantStatus  int
		wantResults bool
		checkJSON   bool
	}{
		{"GET with query", "GET", "getting", http.StatusOK, true, true},
		{"GET with empty query", "GET", "", http.StatusOK, false, true},
		{"GET with no query param", "GET", "", http.StatusOK, false, true},
		{"POST method", "POST", "test", http.StatusMethodNotAllowed, false, false},
		{"PUT method", "PUT", "test", http.StatusMethodNotAllowed, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/search"
			if tt.query != "" {
				url += "?q=" + tt.query
			}

			req := httptest.NewRequest(tt.method, url, nil)
			w := httptest.NewRecorder()

			handler.ServeSearch(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ServeSearch() status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.checkJSON {
				contentType := w.Header().Get("Content-Type")
				if !strings.Contains(contentType, "application/json") {
					t.Errorf("ServeSearch() Content-Type = %q, want application/json", contentType)
				}

				if tt.wantStatus == http.StatusOK {
					var response SearchResponse
					if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
						t.Errorf("ServeSearch() failed to decode JSON: %v", err)
					}

					if tt.wantResults && len(response.Results) == 0 {
						t.Error("ServeSearch() expected results but got none")
					}

					if !tt.wantResults && len(response.Results) > 0 {
						t.Errorf("ServeSearch() expected no results but got %d", len(response.Results))
					}
				}
			}
		})
	}
}

func TestHandleSearchWithBasePath(t *testing.T) {
	_, handler := setupTestSite(t)
	handler.cfg.BasePath = "/docs"

	req := httptest.NewRequest("GET", "/api/search?q=getting", nil)
	w := httptest.NewRecorder()

	handler.handleSearch(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleSearch() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("handleSearch() failed to decode JSON: %v", err)
	}

	if len(response.Results) > 0 {
		// Check that paths include base path
		for _, result := range response.Results {
			if !strings.HasPrefix(result.Path, "/docs") {
				t.Errorf("handleSearch() result path = %q, want prefix /docs", result.Path)
			}
		}
	}
}

func TestServeSitemap(t *testing.T) {
	_, handler := setupTestSite(t)

	tests := []struct {
		name       string
		method     string
		wantStatus int
		checkXML   bool
	}{
		{"GET method", "GET", http.StatusOK, true},
		{"HEAD method", "HEAD", http.StatusOK, false},
		{"POST method", "POST", http.StatusMethodNotAllowed, false},
		{"PUT method", "PUT", http.StatusMethodNotAllowed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/sitemap.xml", nil)
			w := httptest.NewRecorder()

			handler.ServeSitemap(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ServeSitemap() status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.checkXML {
				contentType := w.Header().Get("Content-Type")
				if !strings.Contains(contentType, "application/xml") {
					t.Errorf("ServeSitemap() Content-Type = %q, want application/xml", contentType)
				}

				body := w.Body.String()
				if !strings.Contains(body, "<?xml") {
					t.Error("ServeSitemap() response missing XML declaration")
				}
				if !strings.Contains(body, "<urlset") {
					t.Error("ServeSitemap() response missing urlset element")
				}
				if !strings.Contains(body, "</urlset>") {
					t.Error("ServeSitemap() response missing closing urlset tag")
				}
			}
		})
	}
}

func TestServeSitemapContent(t *testing.T) {
	_, handler := setupTestSite(t)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	req.Host = "example.com"
	w := httptest.NewRecorder()

	handler.ServeSitemap(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ServeSitemap() status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()

	// Check for expected URLs
	if !strings.Contains(body, "http://example.com/") {
		t.Error("ServeSitemap() missing root URL")
	}
	if !strings.Contains(body, "http://example.com/getting-started") {
		t.Error("ServeSitemap() missing getting-started URL")
	}

	// Check for required sitemap elements
	if !strings.Contains(body, "<loc>") {
		t.Error("ServeSitemap() missing <loc> elements")
	}
	if !strings.Contains(body, "<lastmod>") {
		t.Error("ServeSitemap() missing <lastmod> elements")
	}
	if !strings.Contains(body, "<priority>") {
		t.Error("ServeSitemap() missing <priority> elements")
	}
	if !strings.Contains(body, "<changefreq>") {
		t.Error("ServeSitemap() missing <changefreq> elements")
	}
}

func TestServeSitemapWithBasePath(t *testing.T) {
	_, handler := setupTestSite(t)
	handler.cfg.BasePath = "/docs"

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	req.Host = "example.com"
	w := httptest.NewRecorder()

	handler.ServeSitemap(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ServeSitemap() status = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()

	// Check that URLs include base path
	if !strings.Contains(body, "http://example.com/docs/") {
		t.Error("ServeSitemap() URLs missing base path")
	}
}

func TestServeSitemapHTTPS(t *testing.T) {
	_, handler := setupTestSite(t)

	// Create a request with TLS (simulated)
	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	req.Host = "example.com"
	// Note: httptest doesn't support TLS directly, so we test the HTTP case
	// In production, req.TLS would be set by the server
	w := httptest.NewRecorder()

	handler.ServeSitemap(w, req)

	body := w.Body.String()

	// Should use http when TLS is not present (httptest default)
	if !strings.Contains(body, "http://example.com") {
		t.Error("ServeSitemap() should use http for non-TLS requests")
	}
}

func TestSearchResponseStructure(t *testing.T) {
	_, handler := setupTestSite(t)

	req := httptest.NewRequest("GET", "/api/search?q=getting", nil)
	w := httptest.NewRecorder()

	handler.ServeSearch(w, req)

	var response SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response structure
	if response.Query != "getting" {
		t.Errorf("SearchResponse.Query = %q, want %q", response.Query, "getting")
	}

	if len(response.Results) > 0 {
		result := response.Results[0]
		if result.Key == "" {
			t.Error("SearchResultItem.Key is empty")
		}
		if result.Title == "" {
			t.Error("SearchResultItem.Title is empty")
		}
		if result.Path == "" {
			t.Error("SearchResultItem.Path is empty")
		}
	}
}

func TestSitemapPriorityAndChangefreq(t *testing.T) {
	_, handler := setupTestSite(t)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	req.Host = "example.com"
	w := httptest.NewRecorder()

	handler.ServeSitemap(w, req)

	body := w.Body.String()

	// Root page should have highest priority
	if !strings.Contains(body, `<loc>http://example.com/</loc>`) {
		t.Error("ServeSitemap() missing root URL")
	}
	// Check that root has priority 1.0
	rootIndex := strings.Index(body, `<loc>http://example.com/</loc>`)
	if rootIndex != -1 {
		// Find the priority for this URL (should be in the same <url> block)
		urlBlock := body[max(0, rootIndex-200):min(len(body), rootIndex+200)]
		if !strings.Contains(urlBlock, `<priority>1.0</priority>`) {
			t.Error("ServeSitemap() root page should have priority 1.0")
		}
		if !strings.Contains(urlBlock, `<changefreq>weekly</changefreq>`) {
			t.Error("ServeSitemap() root page should have changefreq weekly")
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestSearchWithNoSite(t *testing.T) {
	handler := New(Config{
		DocsDir: t.TempDir(),
		RootDir: t.TempDir(),
	})

	req := httptest.NewRequest("GET", "/api/search?q=test", nil)
	w := httptest.NewRecorder()

	handler.handleSearch(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("handleSearch() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestSitemapWithNoSite(t *testing.T) {
	handler := New(Config{
		DocsDir: t.TempDir(),
		RootDir: t.TempDir(),
	})

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	handler.ServeSitemap(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ServeSitemap() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
