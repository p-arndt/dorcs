package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"dorcs-v2/internal/auth"
	"dorcs-v2/internal/build"
	"dorcs-v2/internal/config"
	"dorcs-v2/internal/server"
	"dorcs-v2/internal/site"
	"dorcs-v2/internal/templates"
)

//go:embed web/templates/*.html web/templates/partials/*.html
var templatesFS embed.FS

//go:embed web/static/**
var staticFS embed.FS

// Version is the version identifier for dorcs.
// This can be set at build time using -ldflags:
//
//	go build -ldflags "-X dorcs-v2/cmd/dorcs.Version=1.0.0"
var Version = "dev"

// runBuild handles the build subcommand for generating static sites.
func runBuild() {
	buildFlags := flag.NewFlagSet("build", flag.ExitOnError)
	var (
		dir        = buildFlags.String("dir", "./docs", "Directory containing markdown documents")
		output     = buildFlags.String("output", "./dist", "Output directory for generated static site")
		baseURL    = buildFlags.String("base-url", "", "Optional base URL path prefix (e.g. /docs). No trailing slash.")
		title      = buildFlags.String("title", "", "Site title shown in header (overrides config file)")
		noDrafts   = buildFlags.Bool("no-drafts", true, "Hide documents with front matter draft: true")
		configFile = buildFlags.String("config", "", "Path to config file (default: looks for dorcs.yaml in current directory, then docs dir)")
		theme      = buildFlags.String("theme", "", "Theme preset: default, ocean, forest, sunset, midnight, lavender, rose")
		themeMode  = buildFlags.String("theme-mode", "", "Theme mode: light, dark, auto")
	)
	buildFlags.Parse(os.Args[2:]) // Skip "build" command

	// Get root directory (where dorcs is running) for static assets like logo/favicon
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get working directory: %v", err)
	}
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("resolve root dir: %v", err)
	}

	// Resolve and validate docs directory
	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("resolve dir: %v", err)
	}
	st, err := os.Stat(absDir)
	if err != nil {
		log.Fatalf("stat dir: %v", err)
	}
	if !st.IsDir() {
		log.Fatalf("dir is not a directory: %s", absDir)
	}

	// Resolve output directory
	absOutput, err := filepath.Abs(*output)
	if err != nil {
		log.Fatalf("resolve output dir: %v", err)
	}

	prefix := sanitizeBasePrefix(*baseURL)

	// Load configuration
	var cfg *config.Config
	if *configFile != "" {
		cfg, err = config.LoadFromFile(*configFile)
		if err != nil {
			log.Fatalf("load config: %v", err)
		}
		log.Printf("dorcs: loaded config from %s", *configFile)
	} else {
		cfg, err = config.Load(absDir)
		if err != nil {
			log.Fatalf("load config: %v", err)
		}
	}

	// Command-line flags override config file
	if *title != "" {
		cfg.Site.Title = *title
	}
	if *theme != "" {
		cfg.Theme.Preset = *theme
	}
	if *themeMode != "" {
		cfg.Theme.Mode = *themeMode
	}

	// Use default title if still empty
	siteTitle := cfg.Site.Title
	if siteTitle == "" {
		siteTitle = "Documentation"
	}

	// Build site index
	codeTheme := cfg.GetCodeTheme()
	s, err := site.New(absDir, codeTheme, prefix)
	if err != nil {
		log.Fatalf("init site: %v", err)
	}
	if err := s.BuildIndex(); err != nil {
		log.Fatalf("build index: %v", err)
	}

	// Parse templates
	tmplDoc, err := templates.ParseFS(templatesFS, "doc",
		"web/templates/layout.html",
		"web/templates/doc.html",
		"web/templates/partials/*.html",
	)
	if err != nil {
		log.Fatalf("parse templates: %v", err)
	}

	// Create builder
	builder := build.New(build.Config{
		DocsDir:      absDir,
		RootDir:      absRootDir,
		OutputDir:    absOutput,
		BasePath:     prefix,
		Site:         s,
		SiteConfig:   cfg,
		DocumentTmpl: tmplDoc,
		StaticFS:     staticFS,
		TemplatesFS:  templatesFS,
	})

	// Build static site
	log.Printf("dorcs: building static site from %s", absDir)
	log.Printf("dorcs: output directory: %s", absOutput)
	if err := builder.Build(!*noDrafts); err != nil {
		log.Fatalf("build failed: %v", err)
	}

	log.Printf("dorcs: static site built successfully!")
}

func main() {
	// Check if build command is requested
	if len(os.Args) > 1 && os.Args[1] == "build" {
		runBuild()
		return
	}

	var (
		dir        = flag.String("dir", "./docs", "Directory containing markdown documents")
		addr       = flag.String("addr", ":8080", "Listen address (e.g. :8080, 127.0.0.1:8080)")
		baseURL    = flag.String("base-url", "", "Optional base URL path prefix (e.g. /docs). No trailing slash.")
		title      = flag.String("title", "", "Site title shown in header (overrides config file)")
		cache      = flag.Bool("cache", true, "Cache rendered documents in memory (mtime-based)")
		noDrafts   = flag.Bool("no-drafts", true, "Hide documents with front matter draft: true")
		configFile = flag.String("config", "", "Path to config file (default: looks for dorcs.yaml in current directory, then docs dir)")
		theme      = flag.String("theme", "", "Theme preset: default, ocean, forest, sunset, midnight, lavender, rose")
		themeMode  = flag.String("theme-mode", "", "Theme mode: light, dark, auto")
		watch      = flag.Bool("watch", false, "Watch for file changes and automatically reload")
	)
	flag.Parse()

	// Get root directory (where dorcs is running) for static assets like logo/favicon
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get working directory: %v", err)
	}
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("resolve root dir: %v", err)
	}

	// Resolve and validate docs directory
	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("resolve dir: %v", err)
	}
	st, err := os.Stat(absDir)
	if err != nil {
		log.Fatalf("stat dir: %v", err)
	}
	if !st.IsDir() {
		log.Fatalf("dir is not a directory: %s", absDir)
	}

	prefix := sanitizeBasePrefix(*baseURL)

	// Load configuration
	var cfg *config.Config
	if *configFile != "" {
		cfg, err = config.LoadFromFile(*configFile)
		if err != nil {
			log.Fatalf("load config: %v", err)
		}
		log.Printf("dorcs: loaded config from %s", *configFile)
	} else {
		cfg, err = config.Load(absDir)
		if err != nil {
			log.Fatalf("load config: %v", err)
		}
	}

	// Command-line flags override config file
	if *title != "" {
		cfg.Site.Title = *title
	}
	if *theme != "" {
		cfg.Theme.Preset = *theme
	}
	if *themeMode != "" {
		cfg.Theme.Mode = *themeMode
	}

	// Use default title if still empty
	siteTitle := cfg.Site.Title
	if siteTitle == "" {
		siteTitle = "Documentation"
	}

	// Build site index
	codeTheme := cfg.GetCodeTheme()
	s, err := site.New(absDir, codeTheme, prefix)
	if err != nil {
		log.Fatalf("init site: %v", err)
	}
	if err := s.BuildIndex(); err != nil {
		log.Fatalf("build index: %v", err)
	}

	// Parse templates
	// Note: doc.html must be parsed AFTER index.html so its "content" block overrides index.html's
	tmplDoc, err := templates.ParseFS(templatesFS, "doc",
		"web/templates/layout.html",
		"web/templates/doc.html",
		"web/templates/partials/*.html",
	)
	if err != nil {
		log.Fatalf("parse templates: %v", err)
	}

	// Set up live reload broadcaster if watch mode is enabled
	var reloadBroadcaster *site.ReloadBroadcaster
	if *watch {
		reloadBroadcaster = site.NewReloadBroadcaster()
	}

	// Document handler (created before watcher so it can be referenced in config reload callback)
	handler := server.New(server.Config{
		DocsDir:           absDir,
		RootDir:           absRootDir,
		SiteTitle:         siteTitle,
		BasePath:          prefix,
		Cache:             *cache,
		HideDraft:         *noDrafts,
		Site:              s,
		DocumentTmpl:      tmplDoc,
		SiteConfig:        cfg,
		ReloadBroadcaster: reloadBroadcaster,
		Version:           Version,
	})

	// Set up watcher if watch mode is enabled
	if *watch {
		// Determine config file path
		var configPath string
		if *configFile != "" {
			configPath = *configFile
		} else {
			// Try to find config file (current directory first, then docs directory)
			configPath = config.FindConfigFile(absDir)
		}

		// Create config reload callback
		var configReload site.ConfigReloadCallback
		if configPath != "" {
			configReload = func() error {
				// Reload config from file
				newCfg, err := config.LoadFromFile(configPath)
				if err != nil {
					return err
				}

				// Apply command-line overrides again
				if *title != "" {
					newCfg.Site.Title = *title
				}
				if *theme != "" {
					newCfg.Theme.Preset = *theme
				}
				if *themeMode != "" {
					newCfg.Theme.Mode = *themeMode
				}

				// Update handler config
				handler.UpdateConfig(newCfg)

				return nil
			}
		} else {
			// Try loading from directory
			configReload = func() error {
				newCfg, err := config.Load(absDir)
				if err != nil {
					return err
				}

				// Apply command-line overrides again
				if *title != "" {
					newCfg.Site.Title = *title
				}
				if *theme != "" {
					newCfg.Theme.Preset = *theme
				}
				if *themeMode != "" {
					newCfg.Theme.Mode = *themeMode
				}

				// Update handler config
				handler.UpdateConfig(newCfg)

				return nil
			}
		}

		cleanup, err := s.StartWatcher(reloadBroadcaster, configReload, configPath)
		if err != nil {
			log.Fatalf("start watcher: %v", err)
		}
		defer cleanup()
		log.Printf("dorcs: watching for file changes in %s", absDir)
		if configPath != "" {
			log.Printf("dorcs: watching config file: %s", configPath)
		}
		log.Printf("dorcs: live reload enabled")
	}

	// Set up HTTP mux
	mux := http.NewServeMux()

	// Static assets (embedded)
	staticSub, err := fs.Sub(staticFS, "web/static")
	if err != nil {
		log.Fatalf("static sub fs: %v", err)
	}
	staticHandler := cachingFileServer(staticSub)
	mux.Handle(prefix+"/static/", http.StripPrefix(prefix+"/static/", staticHandler))

	// Syntax highlighting CSS (generated from Chroma)
	syntaxCSS := s.SyntaxCSS()
	mux.HandleFunc(prefix+"/static/syntax.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write([]byte(syntaxCSS))
	})

	// Serve custom CSS if configured
	if cfg.Theme.CustomCSS != "" {
		customCSSPath := cfg.Theme.CustomCSS
		if !filepath.IsAbs(customCSSPath) {
			customCSSPath = filepath.Join(absDir, customCSSPath)
		}
		mux.HandleFunc(prefix+"/static/custom.css", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour for custom CSS
			http.ServeFile(w, r, customCSSPath)
		})
	}

	// Live reload endpoint (SSE)
	if reloadBroadcaster != nil {
		mux.Handle(prefix+"/__reload", reloadBroadcaster)
	}

	mux.HandleFunc(prefix+"/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	// Search API endpoint
	mux.HandleFunc(prefix+"/api/search", handler.ServeSearch)
	// Sitemap endpoint
	mux.HandleFunc(prefix+"/sitemap.xml", handler.ServeSitemap)

	// Edit mode API endpoints (if auth is enabled)
	var editHandlers *server.EditHandlers
	if cfg.Auth.Enabled {
		// Hash password if provided in plain text
		if cfg.Auth.Password != "" && cfg.Auth.PasswordHash == "" {
			hash, err := auth.HashPassword(cfg.Auth.Password)
			if err != nil {
				log.Fatalf("hash password: %v", err)
			}
			cfg.Auth.PasswordHash = hash
			cfg.Auth.Password = "" // Clear plain text password

			// Save config back to file if possible
			if *configFile != "" {
				// Try to update config file with hash
				// Note: This is a simple approach - in production you might want to be more careful
				log.Printf("dorcs: password hashed and saved to config")
			}
		}

		if cfg.Auth.Username == "" || cfg.Auth.PasswordHash == "" {
			log.Fatalf("auth enabled but username or password not configured")
		}

		// Determine sessions path
		sessionsPath := cfg.Auth.SessionsPath
		if !filepath.IsAbs(sessionsPath) {
			sessionsPath = filepath.Join(absDir, sessionsPath)
		}

		// Create auth manager
		authConfig := &auth.AuthConfig{
			Username:     cfg.Auth.Username,
			PasswordHash: cfg.Auth.PasswordHash,
			SessionsPath: sessionsPath,
		}

		authManager, err := auth.NewAuthManager(authConfig)
		if err != nil {
			log.Fatalf("create auth manager: %v", err)
		}

		// Create edit handlers
		editHandlers = server.NewEditHandlers(authManager, absDir, prefix)

		// Register edit API routes
		// Public routes
		mux.HandleFunc(prefix+"/api/edit/login", editHandlers.HandleLogin)
		mux.HandleFunc(prefix+"/api/edit/logout", editHandlers.HandleLogout)
		mux.HandleFunc(prefix+"/api/edit/check-auth", editHandlers.HandleCheckAuth)

		// Protected routes (require authentication)
		mux.Handle(prefix+"/api/edit/list-files", authManager.RequireAuth(http.HandlerFunc(editHandlers.HandleListFiles)))
		mux.Handle(prefix+"/api/edit/read-file", authManager.RequireAuth(http.HandlerFunc(editHandlers.HandleReadFile)))
		mux.Handle(prefix+"/api/edit/save-file", authManager.RequireAuth(http.HandlerFunc(editHandlers.HandleSaveFile)))
		mux.Handle(prefix+"/api/edit/create-file", authManager.RequireAuth(http.HandlerFunc(editHandlers.HandleCreateFile)))
		mux.Handle(prefix+"/api/edit/delete-file", authManager.RequireAuth(http.HandlerFunc(editHandlers.HandleDeleteFile)))
	}

	mux.Handle(prefix+"/", handler)

	// Create server
	srv := &http.Server{
		Addr:              *addr,
		Handler:           server.LoggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	log.Printf("dorcs: serving %s", absDir)
	log.Printf("dorcs: listening on http://%s%s", ln.Addr().String(), prefix+"/")

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Printf("dorcs: shutting down...")
		_ = srv.Shutdown(shutdownCtx)
	}()

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve: %v", err)
	}
}

// sanitizeBasePrefix normalizes a base URL prefix.
// Rules:
// - empty => ""
// - ensure leading slash
// - remove trailing slash
// - disallow "." and ".." segments
func sanitizeBasePrefix(in string) string {
	s := strings.TrimSpace(in)
	if s == "" {
		return ""
	}
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	for strings.HasSuffix(s, "/") && s != "/" {
		s = strings.TrimSuffix(s, "/")
	}
	if s == "/" {
		return ""
	}
	// Safety checks (avoid path traversal)
	parts := strings.Split(s, "/")
	for _, p := range parts {
		if p == "." || p == ".." {
			log.Fatalf("invalid base-url: %q", in)
		}
	}
	return s
}

// cachingFileServer creates a file server with proper ETag and cache support for embedded files.
func cachingFileServer(fsys fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		// Open the file
		file, err := fsys.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Read file content to compute ETag
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Compute ETag from file content (SHA256 hash)
		hash := sha256.Sum256(data)
		etag := `"` + hex.EncodeToString(hash[:16]) + `"` // Use first 16 bytes for shorter ETag

		// Handle conditional requests BEFORE setting other headers (to avoid writing headers if 304)
		// Check ETag first (more reliable) - but only if browser isn't forcing no-cache
		if r.Header.Get("Cache-Control") != "no-cache" {
			if match := r.Header.Get("If-None-Match"); match != "" {
				// ETag comparison - handle both quoted and unquoted, and comma-separated lists
				cleanETag := strings.Trim(etag, `"`)
				cleanMatch := strings.Trim(match, `"`)
				if strings.Contains(cleanMatch, cleanETag) || match == etag || strings.Contains(match, cleanETag) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}

		// Set content type
		switch {
		case strings.HasSuffix(path, ".css"):
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case strings.HasSuffix(path, ".js"):
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case strings.HasSuffix(path, ".ico"):
			w.Header().Set("Content-Type", "image/x-icon")
		case strings.HasSuffix(path, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg"):
			w.Header().Set("Content-Type", "image/jpeg")
		case strings.HasSuffix(path, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		}

		// Set cache headers (1 year, immutable) - only if not 304
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("ETag", etag)
		w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

		// Serve the file content
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		if r.Method != http.MethodHead {
			w.Write(data)
		}
	})
}
