package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/p-arndt/dorcs/internal/site"
)

type mockGitHubAssetClient struct {
	fetchContent map[string][]byte
}

func newMockGitHubAssetClient() *mockGitHubAssetClient {
	return &mockGitHubAssetClient{
		fetchContent: make(map[string][]byte),
	}
}

func (m *mockGitHubAssetClient) DiscoverMarkdownFiles(owner, repo, branch, rootPath string) ([]string, error) {
	return nil, nil
}

func (m *mockGitHubAssetClient) FetchMarkdown(owner, repo, branch, filePath string) ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s", owner, repo, branch, filePath)
	content, ok := m.fetchContent[key]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}
	return content, nil
}

func TestIsStaticAsset(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{"PNG image", "image.png", true},
		{"JPG image", "photo.jpg", true},
		{"JPEG image", "photo.jpeg", true},
		{"GIF image", "animation.gif", true},
		{"SVG image", "icon.svg", true},
		{"WebP image", "image.webp", true},
		{"ICO favicon", "favicon.ico", true},
		{"PDF document", "document.pdf", true},
		{"ZIP archive", "archive.zip", true},
		{"JSON file", "data.json", true},
		{"XML file", "data.xml", true},
		{"TXT file", "readme.txt", true},
		{"CSS file", "style.css", true},
		{"JS file", "script.js", true},
		{"WOFF font", "font.woff", true},
		{"WOFF2 font", "font.woff2", true},
		{"TTF font", "font.ttf", true},
		{"EOT font", "font.eot", true},
		{"Markdown file", "readme.md", false},
		{"HTML file", "index.html", false},
		{"No extension", "file", false},
		{"Uppercase extension", "IMAGE.PNG", true}, // Should be case-insensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStaticAsset(tt.filePath); got != tt.want {
				t.Errorf("isStaticAsset(%q) = %v, want %v", tt.filePath, got, tt.want)
			}
		})
	}
}

func TestTryServeStaticAsset(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")

	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}

	// Create test files
	testContent := []byte("test content")
	if err := os.WriteFile(filepath.Join(docsDir, "test.png"), testContent, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	handler := New(Config{
		DocsDir: docsDir,
		RootDir: "",
	})

	tests := []struct {
		name       string
		relPath    string
		setup      func(t *testing.T, s *site.Site)
		wantServed bool
		wantStatus int
		wantBody   string
	}{
		{"serve from docs dir", "test.png", nil, true, http.StatusOK, "test content"},
		{"non-existent file", "missing.png", nil, false, http.StatusNotFound, ""},
		{"not a static asset", "readme.md", nil, false, http.StatusNotFound, ""},
		{"path traversal attempt", "../etc/passwd", nil, false, http.StatusNotFound, ""},
		{"absolute path attempt", "/etc/passwd", nil, false, http.StatusNotFound, ""},
		{
			name:    "serve from GitHub when local asset is missing",
			relPath: "guide/test.png",
			setup: func(t *testing.T, s *site.Site) {
				mockClient := newMockGitHubAssetClient()
				mockClient.fetchContent["owner/repo/main/docs/guide/test.png"] = []byte("github asset")
				s.SetGitHubConfig(mockClient, "owner", "repo", "main", "docs")
			},
			wantServed: true,
			wantStatus: http.StatusOK,
			wantBody:   "github asset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSite, err := site.New(docsDir, "github", "", "")
			if err != nil {
				t.Fatalf("failed to create mock site: %v", err)
			}
			if tt.setup != nil {
				tt.setup(t, testSite)
			}

			req := httptest.NewRequest("GET", "/"+tt.relPath, nil)
			w := httptest.NewRecorder()

			served := handler.tryServeStaticAsset(w, req, tt.relPath, testSite)

			if served != tt.wantServed {
				t.Errorf("tryServeStaticAsset() served = %v, want %v", served, tt.wantServed)
			}

			if served {
				if w.Code != http.StatusOK {
					t.Errorf("tryServeStaticAsset() status = %d, want %d", w.Code, http.StatusOK)
				}
				// Check Content-Type header
				contentType := w.Header().Get("Content-Type")
				if contentType == "" {
					t.Error("tryServeStaticAsset() missing Content-Type header")
				}
				// Check Cache-Control header
				cacheControl := w.Header().Get("Cache-Control")
				if cacheControl != "public, max-age=31536000, immutable" {
					t.Errorf("tryServeStaticAsset() Cache-Control = %q, want %q", cacheControl, "public, max-age=31536000, immutable")
				}
				// Check Content-Length header
				contentLength := w.Header().Get("Content-Length")
				if contentLength == "" {
					t.Error("tryServeStaticAsset() missing Content-Length header")
				}
				if body := w.Body.String(); body != tt.wantBody {
					t.Errorf("tryServeStaticAsset() body = %q, want %q", body, tt.wantBody)
				}
			}
		})
	}
}

func TestServeStaticFileMIMETypes(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test")
	testContent := []byte("test content")

	handler := New(Config{
		DocsDir: tmpDir,
		RootDir: tmpDir,
	})

	tests := []struct {
		name        string
		ext         string
		wantMIME    string
		useOverride bool
	}{
		{"PNG image", ".png", "image/png", false},
		{"JPG image", ".jpg", "image/jpeg", false},
		{"ZIP file (override)", ".zip", "application/zip", true},
		{"XML file (override)", ".xml", "application/xml", true},
		{"JS file (override)", ".js", "application/javascript; charset=utf-8", true},
		{"WOFF font (fallback)", ".woff", "font/woff", true},
		{"WOFF2 font (fallback)", ".woff2", "font/woff2", true},
		{"TTF font (fallback)", ".ttf", "font/ttf", true},
		{"EOT font (fallback)", ".eot", "application/vnd.ms-fontobject", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := testFile + tt.ext
			if err := os.WriteFile(filePath, testContent, 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}
			defer os.Remove(filePath)

			file, err := os.Open(filePath)
			if err != nil {
				t.Fatalf("failed to open file: %v", err)
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				t.Fatalf("failed to stat file: %v", err)
			}

			w := httptest.NewRecorder()
			relPath := filepath.Base(filePath)

			handler.serveStaticFile(w, file, stat, relPath)

			contentType := w.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, tt.wantMIME) {
				t.Errorf("serveStaticFile() Content-Type = %q, want prefix %q", contentType, tt.wantMIME)
			}
		})
	}
}
