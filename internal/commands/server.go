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
	"strings"
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
		repo       = flag.String("repo", "", "GitHub repository to bootstrap docs and config from")
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
	configuredDir := *dir

	// Get root directory (where dorcs is running) for static assets like logo/favicon
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get working directory: %v", err)
	}
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("resolve root dir: %v", err)
	}

	// Resolve docs directory path. Validation happens after config load so GitHub-only
	// setups can run without a local docs directory.
	configuredAbsDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("resolve dir: %v", err)
	}

	prefix := SanitizeBasePrefix(*baseURL)

	// Load configuration
	bootstrap, err := loadConfigWithBootstrap(absRootDir, configuredDir, *configFile, *repo)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	cfg := bootstrap.Config
	if bootstrap.Source != "" && bootstrap.Source != "defaults" {
		log.Printf("dorcs: loaded config from %s", bootstrap.Source)
	} else if strings.TrimSpace(*repo) != "" {
		log.Printf("dorcs: no remote config found in %s; using defaults", *repo)
	}

	absDir, cleanupDocsDir, err := ResolveDocsDir(configuredAbsDir, cfg.GitHub.Enabled)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer cleanupDocsDir()

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
	var ghClient github.ClientAPI
	var repoInfo *github.RepositoryInfo
	if cfg.GitHub.Enabled && cfg.GitHub.Repository != "" {
		ghClient, repoInfo, err = setupGitHubIntegration(absRootDir, cfg)
		if err != nil {
			log.Fatalf("setup GitHub integration: %v", err)
		}
		log.Printf("dorcs: note: local docs directory will be ignored when GitHub integration is enabled")
	}

	// Create sites for each language if multi-lingual is enabled
	var defaultSite *site.Site
	sites := make(map[string]*site.Site)
	versionSites := make(map[string]*site.Site)

	// MkDocs-style structure: language-first approach
	// Structure: docs/{lang}/{version}/ or docs/{version}/ (version-only) or docs/ (default)
	defaultVersion := ""
	defaultLang := ""
	if cfg.IsMultiVersion() {
		defaultVersion = cfg.GetDefaultVersion()
	}
	if cfg.IsMultiLingual() {
		defaultLang = cfg.GetDefaultLanguage()
	}

	// Create sites: language-first iteration
	// When using multiple languages, markdown files in root docs/ are ignored.
	// In version-only mode, root docs/ serves the default version.
	// Static assets (images, etc.) in root are served as shared assets in all modes.
	if cfg.IsMultiLingual() {
		rootHasMarkdown := hasMarkdownFilesInDir(absDir)
		if rootHasMarkdown {
			msg := "when using multiple languages, markdown files in root docs/ folder are ignored"
			if cfg.IsMultiVersion() {
				msg += ". Please move all markdown files to language/version folders (docs/en/v1/, docs/de/v1/, etc.)"
			} else {
				msg += ". Please move all markdown files to language-specific folders (docs/en/, docs/de/, etc.)"
			}
			msg += ". Static assets (images, etc.) in root are served as shared assets."
			log.Printf("dorcs: warning: %s Found markdown files in %s.", msg, absDir)
		}
	}

	if cfg.IsMultiLingual() {

		// Multi-language: iterate languages first, then versions inside each language
		for _, lang := range cfg.Languages.Enabled {
			// For default language: always use language folder (markdown files must be in language folders)
			langCodeForSite := lang.Code

			if cfg.IsMultiVersion() {
				// Both languages and versions: docs/{lang}/{version}/
				for _, ver := range cfg.Versions.Enabled {
					versionID := ""
					versionPath := ""
					if ver.ID != defaultVersion {
						versionID = ver.ID
					} else {
						versionPath = ver.ID
					}

					// Create site: language-first structure
					verLangSite, err := site.NewWithVersionPath(absDir, codeTheme, prefix, langCodeForSite, versionID, versionPath)
					if err != nil {
						if ghClient == nil {
							log.Printf("dorcs: warning: language %s version %s folder not found, skipping", lang.Code, ver.ID)
							continue
						}
						log.Fatalf("init site for language %s version %s: %v", lang.Code, ver.ID, err)
					}

					// Set GitHub config if enabled
					if ghClient != nil && repoInfo != nil {
						githubPath := github.ContentPath(repoInfo.Path, langCodeForSite, ver.ID, versionPath, defaultVersion)
						applyRepoToSite(verLangSite, ghClient, repoInfo, githubPath)
					}

					verLangSite.SetExplicitNav(cfg.Nav.Items)
					verLangSite.SetSectionsConfigured(len(cfg.Nav.Sections) > 0)
					verLangSite.SetDefaultVersion(defaultVersion)
					if cfg.IsMultiLingual() {
						verLangSite.SetDefaultLanguage(defaultLang)
					}

					if err := verLangSite.BuildIndex(); err != nil {
						log.Fatalf("build index for language %s version %s: %v", lang.Code, ver.ID, err)
					}

					// Store in versionSites map: key format "{version}:{language}" or "{version}:" for default language
					// Default language uses empty language in key (even if it uses its folder)
					siteKey := ver.ID + ":"
					if lang.Code != defaultLang {
						siteKey = ver.ID + ":" + lang.Code
					}
					versionSites[siteKey] = verLangSite
					// Also store with explicit language code for /en/v1/ access when default uses folder
					if lang.Code == defaultLang && langCodeForSite != "" {
						versionSites[ver.ID+":"+lang.Code] = verLangSite
					}
				}
			} else {
				// Languages only: docs/{lang}/
				langSite, err := site.New(absDir, codeTheme, prefix, langCodeForSite)
				if err != nil {
					if ghClient == nil {
						log.Printf("dorcs: warning: language %s folder not found, skipping", lang.Code)
						continue
					}
					log.Fatalf("init site for language %s: %v", lang.Code, err)
				}

				// Set GitHub config if enabled
				if ghClient != nil && repoInfo != nil {
					githubPath := github.ContentPath(repoInfo.Path, langCodeForSite, "", "", "")
					applyRepoToSite(langSite, ghClient, repoInfo, githubPath)
				}

				langSite.SetExplicitNav(cfg.Nav.Items)
				langSite.SetSectionsConfigured(len(cfg.Nav.Sections) > 0)
				langSite.SetDefaultLanguage(defaultLang)

				if err := langSite.BuildIndex(); err != nil {
					log.Fatalf("build index for language %s: %v", lang.Code, err)
				}

				// Store in sites map: default language uses empty key (even if it uses its folder),
				// others use their code. This ensures / routes to default language correctly.
				if lang.Code == defaultLang {
					sites[""] = langSite
					defaultSite = langSite
					// Also store with language code for explicit /en/ access
					if langCodeForSite != "" {
						sites[lang.Code] = langSite
					}
				} else {
					sites[lang.Code] = langSite
				}
			}
		}
	} else if cfg.IsMultiVersion() {
		// Versions only: docs/{version}/
		for _, ver := range cfg.Versions.Enabled {
			versionID := ""
			versionPath := ""
			if ver.ID != defaultVersion {
				versionID = ver.ID
			} else {
				versionPath = ver.ID
			}

			verSite, err := site.NewWithVersionPath(absDir, codeTheme, prefix, "", versionID, versionPath)
			if err != nil {
				if ghClient == nil {
					log.Printf("dorcs: warning: version %s folder not found, skipping", ver.ID)
					continue
				}
				log.Fatalf("init site for version %s: %v", ver.ID, err)
			}

			// Set GitHub config if enabled
			if ghClient != nil && repoInfo != nil {
				githubPath := github.ContentPath(repoInfo.Path, "", ver.ID, versionPath, defaultVersion)
				applyRepoToSite(verSite, ghClient, repoInfo, githubPath)
			}

			verSite.SetExplicitNav(cfg.Nav.Items)
			verSite.SetSectionsConfigured(len(cfg.Nav.Sections) > 0)
			verSite.SetDefaultVersion(defaultVersion)

			if err := verSite.BuildIndex(); err != nil {
				log.Fatalf("build index for version %s: %v", ver.ID, err)
			}

			// Store in versionSites map: key format "{version}:"
			versionSites[ver.ID+":"] = verSite
		}
	} else {
		// Simple mode: no versioning, no languages - docs/ directly
		singleSite, err := site.New(absDir, codeTheme, prefix, "")
		if err != nil {
			log.Fatalf("init site: %v", err)
		}

		// Set GitHub config if enabled
		if ghClient != nil && repoInfo != nil {
			applyRepoToSite(singleSite, ghClient, repoInfo, repoInfo.Path)
		}

		singleSite.SetExplicitNav(cfg.Nav.Items)
		singleSite.SetSectionsConfigured(len(cfg.Nav.Sections) > 0)
		if err := singleSite.BuildIndex(); err != nil {
			log.Fatalf("build index: %v", err)
		}
		defaultSite = singleSite
	}

	// Set default site if not already set
	if defaultSite == nil {
		if len(versionSites) > 0 {
			// Find default version site
			defaultVerKey := defaultVersion + ":"
			if s, ok := versionSites[defaultVerKey]; ok {
				defaultSite = s
			} else {
				// Use first available
				for _, s := range versionSites {
					defaultSite = s
					break
				}
			}
		} else if len(sites) > 0 {
			// Use default language site
			if s, ok := sites[""]; ok {
				defaultSite = s
			} else {
				// Use first available
				for _, s := range sites {
					defaultSite = s
					break
				}
			}
		}
	}

	if defaultSite == nil {
		log.Fatalf("no site instances available")
	}
	s := defaultSite

	// Parse templates
	// Note: doc.html must be parsed AFTER index.html so its "content" block overrides index.html's
	// Include presentation.html for slide decks
	tmplDoc, err := templates.ParseFS(templatesFS, "doc",
		"web/templates/layout.html",
		"web/templates/doc.html",
		"web/templates/404.html",
		"web/templates/presentation.html",
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
		VersionSites:      versionSites,
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

				if err := rebuildSitesWithExplicitNav(allSites(s, sites, versionSites), newCfg.Nav.Items); err != nil {
					return err
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

				if err := rebuildSitesWithExplicitNav(allSites(s, sites, versionSites), newCfg.Nav.Items); err != nil {
					return err
				}

				// Update handler config
				handler.UpdateConfig(newCfg)

				return nil
			}
		}

		// Rebuild all sites when any file changes (handles multi-lang/version)
		rebuildAll := func() error {
			seen := make(map[*site.Site]struct{})
			for _, site := range []*site.Site{s} {
				if site != nil {
					if _, ok := seen[site]; !ok {
						seen[site] = struct{}{}
						if err := site.BuildIndex(); err != nil {
							return err
						}
					}
				}
			}
			for _, site := range sites {
				if site != nil {
					if _, ok := seen[site]; !ok {
						seen[site] = struct{}{}
						if err := site.BuildIndex(); err != nil {
							return err
						}
					}
				}
			}
			for _, site := range versionSites {
				if site != nil {
					if _, ok := seen[site]; !ok {
						seen[site] = struct{}{}
						if err := site.BuildIndex(); err != nil {
							return err
						}
					}
				}
			}
			return nil
		}

		cleanup, err := s.StartWatcher(absDir, reloadBroadcaster, configReload, configPath, rebuildAll)
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

func collectSites(siteMap map[string]*site.Site) []*site.Site {
	out := make([]*site.Site, 0, len(siteMap))
	for _, s := range siteMap {
		out = append(out, s)
	}
	return out
}

func allSites(defaultSite *site.Site, sites map[string]*site.Site, versionSites map[string]*site.Site) []*site.Site {
	out := make([]*site.Site, 0, 1+len(sites)+len(versionSites))
	out = append(out, defaultSite)
	out = append(out, collectSites(sites)...)
	out = append(out, collectSites(versionSites)...)
	return out
}

func rebuildSitesWithExplicitNav(targets []*site.Site, items config.NavItems) error {
	seen := make(map[*site.Site]struct{}, len(targets))
	previousNav := make(map[*site.Site]config.NavItems, len(targets))
	updated := make([]*site.Site, 0, len(targets))

	for _, currentSite := range targets {
		if currentSite == nil {
			continue
		}
		if _, ok := seen[currentSite]; ok {
			continue
		}
		seen[currentSite] = struct{}{}

		previousNav[currentSite] = currentSite.ExplicitNav()
		currentSite.SetExplicitNav(items)
		if err := currentSite.BuildIndex(); err != nil {
			if rollbackErr := rollbackExplicitNav(updated, previousNav); rollbackErr != nil {
				return fmt.Errorf("apply explicit nav: %w (rollback failed: %v)", err, rollbackErr)
			}
			return err
		}
		updated = append(updated, currentSite)
	}

	return nil
}

func rollbackExplicitNav(targets []*site.Site, snapshots map[*site.Site]config.NavItems) error {
	for i := len(targets) - 1; i >= 0; i-- {
		currentSite := targets[i]
		currentSite.SetExplicitNav(snapshots[currentSite])
		if err := currentSite.BuildIndex(); err != nil {
			return err
		}
	}
	return nil
}

// hasMarkdownFilesInDir checks if a directory contains any markdown files.
func hasMarkdownFilesInDir(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(strings.ToLower(name), ".md") {
				return true
			}
		}
	}
	return false
}
