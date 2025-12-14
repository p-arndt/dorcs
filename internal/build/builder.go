// Package build provides functionality for generating static HTML sites.
package build

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"dorcs-v2/internal/config"
	"dorcs-v2/internal/server"
	"dorcs-v2/internal/site"
)

// Builder generates a static HTML site from markdown documents.
type Builder struct {
	docsDir     string
	rootDir     string
	outputDir   string
	basePath    string
	site        *site.Site
	config      *config.Config
	template    *template.Template
	staticFS    embed.FS
	templatesFS embed.FS
}

// Config holds configuration for the builder.
type Config struct {
	DocsDir      string
	RootDir      string
	OutputDir    string
	BasePath     string
	Site         *site.Site
	SiteConfig   *config.Config
	DocumentTmpl *template.Template
	StaticFS     embed.FS
	TemplatesFS  embed.FS
}

// New creates a new builder.
func New(cfg Config) *Builder {
	return &Builder{
		docsDir:     cfg.DocsDir,
		rootDir:     cfg.RootDir,
		outputDir:   cfg.OutputDir,
		basePath:    cfg.BasePath,
		site:        cfg.Site,
		config:      cfg.SiteConfig,
		template:    cfg.DocumentTmpl,
		staticFS:    cfg.StaticFS,
		templatesFS: cfg.TemplatesFS,
	}
}

// Build generates the static site.
func (b *Builder) Build(includeDrafts bool) error {
	// Clean and create output directory
	if err := os.RemoveAll(b.outputDir); err != nil {
		return fmt.Errorf("remove output dir: %w", err)
	}
	if err := os.MkdirAll(b.outputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Copy static assets
	if err := b.copyStaticAssets(); err != nil {
		return fmt.Errorf("copy static assets: %w", err)
	}

	// Generate syntax highlighting CSS
	if err := b.writeSyntaxCSS(); err != nil {
		return fmt.Errorf("write syntax CSS: %w", err)
	}

	// Copy custom CSS if configured
	if err := b.copyCustomCSS(); err != nil {
		return fmt.Errorf("copy custom CSS: %w", err)
	}

	// Get all documents
	docs := b.site.ListDocs(includeDrafts)

	// Create a handler for rendering pages
	handler := server.New(server.Config{
		DocsDir:      b.docsDir,
		RootDir:      b.rootDir,
		SiteTitle:    b.getSiteTitle(),
		BasePath:     b.basePath,
		Cache:        false, // Not needed for static build
		HideDraft:    !includeDrafts,
		Site:         b.site,
		DocumentTmpl: b.template,
		SiteConfig:   b.config,
		Version:      "static",
	})

	// Render each document
	for _, doc := range docs {
		if err := b.renderDocument(handler, doc); err != nil {
			return fmt.Errorf("render document %s: %w", doc.Key, err)
		}
	}

	// Generate sitemap
	if err := b.generateSitemap(docs); err != nil {
		return fmt.Errorf("generate sitemap: %w", err)
	}

	// Copy static assets from docs directory (images, etc.)
	if err := b.copyDocsStaticAssets(); err != nil {
		return fmt.Errorf("copy docs static assets: %w", err)
	}

	return nil
}

// copyStaticAssets copies embedded static files to the output directory.
func (b *Builder) copyStaticAssets() error {
	staticSub, err := fs.Sub(b.staticFS, "web/static")
	if err != nil {
		return fmt.Errorf("static sub fs: %w", err)
	}

	return fs.WalkDir(staticSub, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Skip live-reload.js for static builds (not needed)
		if d.Name() == "live-reload.js" {
			return nil
		}

		// Read file from embedded FS
		src, err := staticSub.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		// Create destination path
		dstPath := filepath.Join(b.outputDir, "static", path)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		// Create destination file
		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy content
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	})
}

// writeSyntaxCSS writes the syntax highlighting CSS to the output directory.
func (b *Builder) writeSyntaxCSS() error {
	css := b.site.SyntaxCSS()
	cssPath := filepath.Join(b.outputDir, "static", "syntax.css")
	if err := os.MkdirAll(filepath.Dir(cssPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(cssPath, []byte(css), 0644)
}

// copyCustomCSS copies custom CSS if configured.
func (b *Builder) copyCustomCSS() error {
	if b.config == nil || b.config.Theme.CustomCSS == "" {
		return nil
	}

	customCSSPath := b.config.Theme.CustomCSS
	if !filepath.IsAbs(customCSSPath) {
		customCSSPath = filepath.Join(b.docsDir, customCSSPath)
	}

	// Check if file exists
	if _, err := os.Stat(customCSSPath); os.IsNotExist(err) {
		return nil // File doesn't exist, skip silently
	}

	// Read source file
	src, err := os.Open(customCSSPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination
	dstPath := filepath.Join(b.outputDir, "static", "custom.css")
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy content
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// renderDocument renders a single document to HTML and writes it to the appropriate path.
func (b *Builder) renderDocument(handler *server.Handler, doc *site.Doc) error {
	// Build the URL path for this document
	var urlPath string
	if doc.Key == "" {
		urlPath = "/"
	} else {
		urlPath = "/" + doc.Key
	}

	// Add base path if configured
	if b.basePath != "" {
		urlPath = b.basePath + urlPath
	}

	// Determine output file path
	var outputPath string
	if doc.Key == "" {
		// Root index -> index.html
		outputPath = filepath.Join(b.outputDir, "index.html")
	} else {
		// Regular page -> create directory and index.html
		// e.g., "guide/getting-started" -> guide/getting-started/index.html
		outputPath = filepath.Join(b.outputDir, doc.Key, "index.html")
	}

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Create file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	// Create HTTP request
	req := httptest.NewRequest("GET", urlPath, nil)

	// Create response recorder that writes to file
	recorder := httptest.NewRecorder()

	// Render the document
	handler.ServeHTTP(recorder, req)

	// Check status code
	if recorder.Code != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", recorder.Code)
	}

	// Write response body to file
	if _, err := file.Write(recorder.Body.Bytes()); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// copyDocsStaticAssets copies static assets from the docs directory.
func (b *Builder) copyDocsStaticAssets() error {
	return filepath.WalkDir(b.docsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Skip certain files that shouldn't be copied
		name := d.Name()
		if name == ".dorcs_sessions.json" {
			return nil
		}

		// Check if it's a static asset
		if !isStaticAsset(path) {
			return nil
		}

		// Get relative path from docs directory
		rel, err := filepath.Rel(b.docsDir, path)
		if err != nil {
			return err
		}

		// Create destination path
		dstPath := filepath.Join(b.outputDir, rel)

		// Create destination directory
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		// Copy file
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	})
}

// generateSitemap generates a sitemap.xml file.
func (b *Builder) generateSitemap(docs []*site.Doc) error {
	sitemapPath := filepath.Join(b.outputDir, "sitemap.xml")
	file, err := os.Create(sitemapPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write XML header
	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	file.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	// Base URL (empty for relative URLs, or can be configured)
	baseURL := b.basePath
	if baseURL != "" && !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	for _, doc := range docs {
		// Build URL path
		urlPath := baseURL
		if doc.Key == "" {
			if !strings.HasSuffix(urlPath, "/") {
				urlPath += "/"
			}
		} else {
			if !strings.HasSuffix(urlPath, "/") {
				urlPath += "/"
			}
			// URL-encode path segments
			parts := strings.Split(doc.Key, "/")
			encodedParts := make([]string, len(parts))
			for i, part := range parts {
				encodedParts[i] = escapeXML(part)
			}
			urlPath += strings.Join(encodedParts, "/")
		}

		// Format last modified date
		lastmod := doc.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z")

		// Determine priority
		priority := "0.7"
		if doc.Key == "" {
			priority = "1.0"
		} else if !strings.Contains(doc.Key, "/") {
			priority = "0.9"
		}

		// Determine changefreq
		changefreq := "monthly"
		if doc.Key == "" {
			changefreq = "weekly"
		}

		// Write URL entry
		file.WriteString("  <url>\n")
		fmt.Fprintf(file, "    <loc>%s</loc>\n", urlPath)
		fmt.Fprintf(file, "    <lastmod>%s</lastmod>\n", lastmod)
		fmt.Fprintf(file, "    <changefreq>%s</changefreq>\n", changefreq)
		fmt.Fprintf(file, "    <priority>%s</priority>\n", priority)
		file.WriteString("  </url>\n")
	}

	file.WriteString("</urlset>\n")
	return nil
}

// getSiteTitle returns the site title from config or a default.
func (b *Builder) getSiteTitle() string {
	if b.config != nil && b.config.Site.Title != "" {
		return b.config.Site.Title
	}
	return "Documentation"
}

// isStaticAsset checks if a file path looks like a static asset.
func isStaticAsset(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	staticExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
		".svg": true, ".webp": true, ".ico": true, ".pdf": true,
		".zip": true, ".json": true, ".xml": true, ".txt": true,
		".css": true, ".js": true, ".woff": true, ".woff2": true,
		".ttf": true, ".eot": true,
	}
	return staticExts[ext]
}

// escapeXML escapes special XML characters.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
