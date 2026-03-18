// Package build provides functionality for generating static HTML sites.
package build

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/server"
	"github.com/p-arndt/dorcs/internal/site"
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
	parallelism int
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
	// Parallelism limits concurrent page renders. 0 = auto (NumCPU).
	Parallelism int
}

// New creates a new builder.
func New(cfg Config) *Builder {
	p := cfg.Parallelism
	if p == 0 {
		p = runtime.NumCPU()
	}
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
		parallelism: p,
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

	// Copy static assets (shared across all languages)
	if err := b.copyStaticAssets(); err != nil {
		return fmt.Errorf("copy static assets: %w", err)
	}

	// Generate syntax highlighting CSS (shared)
	if err := b.writeSyntaxCSS(); err != nil {
		return fmt.Errorf("write syntax CSS: %w", err)
	}

	// Copy custom CSS if configured (shared)
	if err := b.copyCustomCSS(); err != nil {
		return fmt.Errorf("copy custom CSS: %w", err)
	}

	// Check if versioning is enabled
	// MkDocs-style structure: language-first approach
	defaultVersion := ""
	defaultLang := ""
	if b.config != nil && b.config.IsMultiVersion() {
		defaultVersion = b.config.GetDefaultVersion()
	}
	if b.config != nil && b.config.IsMultiLingual() {
		defaultLang = b.config.GetDefaultLanguage()
	}

	if b.config != nil && b.config.IsMultiLingual() {
		// Multi-language: iterate languages first, then versions inside each language
		for _, lang := range b.config.Languages.Enabled {
			isDefaultLang := lang.Code == defaultLang
			langCodeForSite := ""
			if !isDefaultLang {
				langCodeForSite = lang.Code
			}

			if b.config.IsMultiVersion() {
				// Both languages and versions: docs/{lang}/{version}/
				for _, ver := range b.config.Versions.Enabled {
					isDefaultVersion := ver.ID == defaultVersion
					versionPath := ""
					if isDefaultVersion {
						versionPath = ver.ID
					}
					if err := b.buildVersionLanguage(ver.ID, versionPath, isDefaultVersion, langCodeForSite, isDefaultLang, lang.Code, includeDrafts); err != nil {
						return fmt.Errorf("build version %s language %s: %w", ver.ID, lang.Code, err)
					}
				}
			} else {
				// Languages only: docs/{lang}/
				if err := b.buildLanguage(langCodeForSite, isDefaultLang, lang.Code, includeDrafts); err != nil {
					return fmt.Errorf("build language %s: %w", lang.Code, err)
				}
			}
		}
	} else if b.config != nil && b.config.IsMultiVersion() {
		// Versions only: docs/{version}/
		for _, ver := range b.config.Versions.Enabled {
			isDefaultVersion := ver.ID == defaultVersion
			versionPath := ""
			if isDefaultVersion {
				versionPath = ver.ID
			}
			if err := b.buildVersionLanguage(ver.ID, versionPath, isDefaultVersion, "", true, "", includeDrafts); err != nil {
				return fmt.Errorf("build version %s: %w", ver.ID, err)
			}
		}
	} else {
		// Simple mode: no versioning, no languages - docs/ directly
		if err := b.buildLanguage("", true, "", includeDrafts); err != nil {
			return fmt.Errorf("build site: %w", err)
		}
	}

	// Copy static assets from docs directory (images, etc.)
	if err := b.copyDocsStaticAssets(); err != nil {
		return fmt.Errorf("copy docs static assets: %w", err)
	}

	return nil
}

// buildLanguage builds the static site for a specific language.
// langCodeForSite is the language code to pass to site.New (empty for default)
// isDefault indicates if this is the default language
// langCodeForURL is the language code for URLs (used in output paths)
func (b *Builder) buildLanguage(langCodeForSite string, isDefault bool, langCodeForURL string, includeDrafts bool) error {
	// Create site instance for this language
	codeTheme := b.config.GetCodeTheme()
	langSite, err := site.New(b.docsDir, codeTheme, b.basePath, langCodeForSite)
	if err != nil {
		// If language folder doesn't exist, skip it
		return nil
	}
	langSite.SetExplicitNav(b.config.Nav.Items)
	if err := langSite.BuildIndex(); err != nil {
		return fmt.Errorf("build index for language %s: %w", langCodeForURL, err)
	}

	// Get all documents for this language
	docs := langSite.ListDocs(includeDrafts)

	// Determine output directory for this language
	var langOutputDir string
	if isDefault {
		langOutputDir = b.outputDir
	} else {
		langOutputDir = filepath.Join(b.outputDir, langCodeForURL)
		if err := os.MkdirAll(langOutputDir, 0755); err != nil {
			return fmt.Errorf("create language output dir: %w", err)
		}
	}

	// Create sites map for handler
	sites := make(map[string]*site.Site)
	if isDefault {
		sites[""] = langSite
	} else {
		sites[langCodeForURL] = langSite
	}

	// Create a handler for rendering pages
	handler := server.New(server.Config{
		DocsDir:      b.docsDir,
		RootDir:      b.rootDir,
		SiteTitle:    b.getSiteTitle(),
		BasePath:     b.basePath,
		Cache:        false, // Not needed for static build
		HideDraft:    !includeDrafts,
		Site:         langSite,
		Sites:        sites,
		DocumentTmpl: b.template,
		SiteConfig:   b.config,
		Version:      "static",
	})

	// Render each document concurrently using a limited worker pool.
	// This significantly speeds up static generation on multi-core machines.
	if len(docs) > 0 {
		g, _ := errgroup.WithContext(context.Background())
		sem := make(chan struct{}, b.parallelism)
		for _, doc := range docs {
			d := doc // capture
			sem <- struct{}{}
			g.Go(func() error {
				defer func() { <-sem }()
				if err := b.renderDocumentForLanguage(handler, d, langCodeForURL, isDefault, langOutputDir); err != nil {
					return fmt.Errorf("render document %s: %w", d.Key, err)
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	}

	// Generate sitemap for this language
	if err := b.generateSitemapForLanguage(docs, langCodeForURL, isDefault, langOutputDir); err != nil {
		return fmt.Errorf("generate sitemap: %w", err)
	}

	return nil
}

// buildVersionLanguage builds the static site for a specific version and language.
// versionID is the version identifier (e.g., "v1", "v2")
// versionPath is the path to the version folder (empty for default version)
// isDefaultVersion indicates if this is the default version
// langCodeForSite is the language code to pass to site.New (empty for default)
// isDefaultLang indicates if this is the default language
// langCodeForURL is the language code for URLs (used in output paths)
func (b *Builder) buildVersionLanguage(versionID string, versionPath string, isDefaultVersion bool, langCodeForSite string, isDefaultLang bool, langCodeForURL string, includeDrafts bool) error {
	// Create site instance for this version and language
	codeTheme := b.config.GetCodeTheme()
	verLangSite, err := site.NewWithVersionPath(b.docsDir, codeTheme, b.basePath, langCodeForSite, versionID, versionPath)
	if err != nil {
		// If version/language folder doesn't exist, skip it
		return nil
	}
	verLangSite.SetExplicitNav(b.config.Nav.Items)
	if err := verLangSite.BuildIndex(); err != nil {
		return fmt.Errorf("build index for version %s language %s: %w", versionID, langCodeForURL, err)
	}

	// Get all documents for this version and language
	docs := verLangSite.ListDocs(includeDrafts)

	// Determine output directory for this version and language (language-first structure)
	var outputDir string
	if isDefaultLang {
		// Default language: output to root or version subdirectory
		if isDefaultVersion {
			outputDir = b.outputDir
		} else {
			outputDir = filepath.Join(b.outputDir, versionID)
		}
	} else {
		// Non-default language: output to language subdirectory, then version if needed
		if isDefaultVersion {
			outputDir = filepath.Join(b.outputDir, langCodeForURL)
		} else {
			outputDir = filepath.Join(b.outputDir, langCodeForURL, versionID)
		}
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Create sites map for handler
	sites := make(map[string]*site.Site)
	versionSites := make(map[string]*site.Site)
	if isDefaultLang {
		sites[""] = verLangSite
		versionSites[versionID+":"] = verLangSite
	} else {
		sites[langCodeForURL] = verLangSite
		versionSites[versionID+":"+langCodeForURL] = verLangSite
	}

	// Create a handler for rendering pages
	handler := server.New(server.Config{
		DocsDir:      b.docsDir,
		RootDir:      b.rootDir,
		SiteTitle:    b.getSiteTitle(),
		BasePath:     b.basePath,
		Cache:        false, // Not needed for static build
		HideDraft:    !includeDrafts,
		Site:         verLangSite,
		Sites:        sites,
		VersionSites: versionSites,
		DocumentTmpl: b.template,
		SiteConfig:   b.config,
		Version:      "static",
	})

	// Render each document concurrently using a limited worker pool.
	if len(docs) > 0 {
		g, _ := errgroup.WithContext(context.Background())
		sem := make(chan struct{}, b.parallelism)
		for _, doc := range docs {
			d := doc // capture
			sem <- struct{}{}
			g.Go(func() error {
				defer func() { <-sem }()
				if err := b.renderDocumentForVersionLanguage(handler, d, versionID, isDefaultVersion, langCodeForURL, isDefaultLang, outputDir); err != nil {
					return fmt.Errorf("render document %s: %w", d.Key, err)
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	}

	// Generate sitemap for this version and language
	if err := b.generateSitemapForVersionLanguage(docs, versionID, isDefaultVersion, langCodeForURL, isDefaultLang, outputDir); err != nil {
		return fmt.Errorf("generate sitemap: %w", err)
	}

	return nil
}

// renderDocumentForVersionLanguage renders a single document for a specific version and language.
func (b *Builder) renderDocumentForVersionLanguage(handler *server.Handler, doc *site.Doc, versionID string, isDefaultVersion bool, langCode string, isDefaultLang bool, outputDir string) error {
	// Build the URL path for this document (language-first structure)
	var urlPath string
	if doc.Key == "" {
		urlPath = "/"
	} else {
		urlPath = "/" + doc.Key
	}

	// Add language prefix if not default (language-first)
	if !isDefaultLang && langCode != "" {
		if urlPath == "/" {
			urlPath = "/" + langCode + "/"
		} else {
			urlPath = "/" + langCode + urlPath
		}
	}

	// Add version prefix if not default (after language)
	if !isDefaultVersion {
		if urlPath == "/" {
			urlPath = "/" + versionID + "/"
		} else if !isDefaultLang && langCode != "" {
			// Language already in path: /en/... -> /en/v1/...
			urlPath = strings.Replace(urlPath, "/"+langCode+"/", "/"+langCode+"/"+versionID+"/", 1)
			if !strings.HasSuffix(urlPath, "/") && doc.Key == "" {
				urlPath += "/"
			}
		} else {
			// No language: /... -> /v1/...
			urlPath = "/" + versionID + urlPath
		}
	}

	// Add base path if configured
	if b.basePath != "" {
		urlPath = b.basePath + urlPath
	}

	// Determine output file path
	var outputPath string
	if doc.Key == "" {
		// Root index -> index.html
		outputPath = filepath.Join(outputDir, "index.html")
	} else {
		// Regular page -> create directory and index.html
		outputPath = filepath.Join(outputDir, doc.Key, "index.html")
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

// generateSitemapForVersionLanguage generates a sitemap.xml file for a specific version and language.
func (b *Builder) generateSitemapForVersionLanguage(docs []*site.Doc, versionID string, isDefaultVersion bool, langCode string, isDefaultLang bool, outputDir string) error {
	sitemapPath := filepath.Join(outputDir, "sitemap.xml")
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
		// Build URL path (language-first structure)
		urlPath := baseURL
		if !isDefaultLang && langCode != "" {
			urlPath += langCode + "/"
		}
		if !isDefaultVersion {
			urlPath += versionID + "/"
		}
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

// renderDocumentForLanguage renders a single document for a specific language.
func (b *Builder) renderDocumentForLanguage(handler *server.Handler, doc *site.Doc, langCode string, isDefault bool, langOutputDir string) error {
	// Build the URL path for this document
	var urlPath string
	if doc.Key == "" {
		urlPath = "/"
	} else {
		urlPath = "/" + doc.Key
	}

	// Add language prefix if not default
	if !isDefault && langCode != "" {
		urlPath = "/" + langCode + urlPath
	}

	// Add base path if configured
	if b.basePath != "" {
		urlPath = b.basePath + urlPath
	}

	// Determine output file path
	var outputPath string
	if doc.Key == "" {
		// Root index -> index.html
		outputPath = filepath.Join(langOutputDir, "index.html")
	} else {
		// Regular page -> create directory and index.html
		// e.g., "guide/getting-started" -> guide/getting-started/index.html
		outputPath = filepath.Join(langOutputDir, doc.Key, "index.html")
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

// generateSitemapForLanguage generates a sitemap.xml file for a specific language.
func (b *Builder) generateSitemapForLanguage(docs []*site.Doc, langCode string, isDefault bool, langOutputDir string) error {
	sitemapPath := filepath.Join(langOutputDir, "sitemap.xml")
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
		if !isDefault && langCode != "" {
			urlPath += langCode + "/"
		}
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
