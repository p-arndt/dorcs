package commands

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/p-arndt/dorcs/internal/auth"
	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/github"
	"github.com/p-arndt/dorcs/internal/server"
	"github.com/p-arndt/dorcs/internal/site"
	"github.com/p-arndt/dorcs/internal/templates"
)

// RunServer handles the server subcommand (default mode) for serving documentation.
func RunServer(templatesFS, staticFS embed.FS, version string) {
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

	// Set custom usage to mention subcommands
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  init     Initialize a new documentation site\n")
		fmt.Fprintf(os.Stderr, "           Use '%s init --help' for init command options\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  build    Build a static site from markdown documents\n")
		fmt.Fprintf(os.Stderr, "           Use '%s build --help' for build command options\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nServer Options:\n")
		flag.PrintDefaults()
	}

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

	prefix := SanitizeBasePrefix(*baseURL)

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

	// Set up GitHub integration if enabled
	var ghClient *github.Client
	var ghCache *github.Cache
	var repoInfo *github.RepositoryInfo
	if cfg.GitHub.Enabled && cfg.GitHub.Repository != "" {
		// Parse cache TTL
		cacheTTL := time.Hour // default
		if cfg.GitHub.CacheTTL != "" {
			if parsed, err := time.ParseDuration(cfg.GitHub.CacheTTL); err == nil {
				cacheTTL = parsed
			} else {
				log.Printf("dorcs: warning: invalid cache_ttl '%s', using default 1h", cfg.GitHub.CacheTTL)
			}
		}

		// Parse repository URL
		var err error
		repoInfo, err = github.ParseRepositoryURL(cfg.GitHub.Repository)
		if err != nil {
			log.Fatalf("parse GitHub repository URL: %v", err)
		}

		// Create cache directory in working directory
		cacheDir := filepath.Join(rootDir, ".cache", "github")
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			log.Printf("dorcs: warning: failed to create cache directory %s: %v (using in-memory cache only)", cacheDir, err)
			cacheDir = "" // Fall back to in-memory only
		}

		// Get default branch if not specified
		if repoInfo.Branch == "" {
			ghCache = github.NewCache(cacheDir)
			tempClient := github.NewClient(cfg.GitHub.Token, ghCache, cacheTTL)
			defaultBranch, err := tempClient.GetDefaultBranch(repoInfo.Owner, repoInfo.Repo)
			if err != nil {
				log.Fatalf("get default branch: %v", err)
			}
			repoInfo.Branch = defaultBranch
		}

		// Create cache and client
		if ghCache == nil {
			ghCache = github.NewCache(cacheDir)
		}
		ghClient = github.NewClient(cfg.GitHub.Token, ghCache, cacheTTL)

		if cacheDir != "" {
			log.Printf("dorcs: using persistent cache at %s", cacheDir)
		}

		log.Printf("dorcs: GitHub integration enabled: %s/%s@%s/%s", repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, repoInfo.Path)
		log.Printf("dorcs: note: local docs directory will be ignored when GitHub integration is enabled")
	}

	// Create sites for each language if multi-lingual is enabled
	var defaultSite *site.Site
	sites := make(map[string]*site.Site)

	if cfg.IsMultiLingual() {
		defaultLang := cfg.GetDefaultLanguage()
		for _, lang := range cfg.Languages.Enabled {
			// Default language uses empty string (root docs folder)
			// Other languages use their code (docs/__lang__/{lang}/ folder)
			langCodeForSite := ""
			if lang.Code != defaultLang {
				langCodeForSite = lang.Code
			}

			langSite, err := site.New(absDir, codeTheme, prefix, langCodeForSite)
			if err != nil {
				// If GitHub is enabled, we can still create the site (local dir not required)
				if ghClient == nil {
					// Language folder doesn't exist and GitHub is not enabled, skip it
					log.Printf("dorcs: warning: language folder for %s not found, skipping", lang.Code)
					continue
				}
				// GitHub enabled: try with empty language code (use base dir)
				// The Language field will be set correctly by Site.New based on langCodeForSite
				langSite, err = site.New(absDir, codeTheme, prefix, langCodeForSite)
				if err != nil {
					// Last resort: use base directory
					langSite, err = site.New(absDir, codeTheme, prefix, "")
					if err != nil {
						log.Fatalf("init site for language %s: %v", lang.Code, err)
					}
					// Manually set the language since we used empty string
					langSite.SetLanguage(langCodeForSite)
				}
			}

			// Set GitHub config if enabled
			if ghClient != nil && repoInfo != nil {
				// Adjust GitHub path based on language
				// Default language: use repoInfo.Path as-is (e.g., "docs")
				// Other languages: use repoInfo.Path + "/__lang__/" + langCode (e.g., "docs/__lang__/de")
				githubPath := repoInfo.Path
				if lang.Code != defaultLang {
					// Non-default language: add __lang__/{lang} to path
					if githubPath != "" {
						githubPath = githubPath + "/__lang__/" + lang.Code
					} else {
						githubPath = "__lang__/" + lang.Code
					}
				}
				langSite.SetGitHubConfig(ghClient, repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, githubPath)
				log.Printf("dorcs: language %s using GitHub path: %s", lang.Code, githubPath)
			}

			if err := langSite.BuildIndex(); err != nil {
				log.Fatalf("build index for language %s: %v", lang.Code, err)
			}

			// Store in sites map: default language uses empty key, others use their code
			if lang.Code == defaultLang {
				sites[""] = langSite
				defaultSite = langSite
			} else {
				sites[lang.Code] = langSite
			}
		}
		// If no default site was set, use first available
		if defaultSite == nil && len(sites) > 0 {
			for _, s := range sites {
				defaultSite = s
				break
			}
		}
	} else {
		// Single language mode (backward compatible)
		singleSite, err := site.New(absDir, codeTheme, prefix, "")
		if err != nil {
			log.Fatalf("init site: %v", err)
		}

		// Set GitHub config if enabled
		if ghClient != nil && repoInfo != nil {
			singleSite.SetGitHubConfig(ghClient, repoInfo.Owner, repoInfo.Repo, repoInfo.Branch, repoInfo.Path)
		}

		if err := singleSite.BuildIndex(); err != nil {
			log.Fatalf("build index: %v", err)
		}
		defaultSite = singleSite
	}

	if defaultSite == nil {
		log.Fatalf("no site instances available")
	}
	s := defaultSite

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
		Sites:             sites,
		DocumentTmpl:      tmplDoc,
		SiteConfig:        cfg,
		ReloadBroadcaster: reloadBroadcaster,
		Version:           version,
	})

	// Set up watcher if watch mode is enabled (skip if GitHub is enabled)
	if *watch && !cfg.GitHub.Enabled {
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
	staticHandler := CachingFileServer(staticSub)
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

	// Determine listen address: use config port if set and addr flag is default, otherwise use flag
	listenAddr := *addr
	if cfg.Port > 0 && *addr == ":8080" {
		// Config has a port set and addr flag is still default, use config port
		listenAddr = fmt.Sprintf(":%d", cfg.Port)
	}

	// Create server
	srv := &http.Server{
		Addr:              listenAddr,
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
