package commands

import (
	"embed"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/p-arndt/dorcs/internal/build"
	"github.com/p-arndt/dorcs/internal/github"
	"github.com/p-arndt/dorcs/internal/site"
	"github.com/p-arndt/dorcs/internal/templates"
)

// RunBuild handles the build subcommand for generating static sites.
func RunBuild(templatesFS, staticFS embed.FS) {
	buildFlags := flag.NewFlagSet("build", flag.ExitOnError)
	var (
		dir        = buildFlags.String("dir", "./docs", "Directory containing markdown documents")
		output     = buildFlags.String("output", "./dist", "Output directory for generated static site")
		baseURL    = buildFlags.String("base-url", "", "Optional base URL path prefix (e.g. /docs). No trailing slash.")
		title      = buildFlags.String("title", "", "Site title shown in header (overrides config file)")
		noDrafts   = buildFlags.Bool("no-drafts", true, "Hide documents with front matter draft: true")
		configFile = buildFlags.String("config", "", "Path to config file (default: looks for dorcs.yaml in current directory, then docs dir)")
		repo       = buildFlags.String("repo", "", "GitHub repository to bootstrap docs and config from")
		theme      = buildFlags.String("theme", "", "Theme preset: default, ocean, forest, sunset, midnight, lavender, rose")
		themeMode  = buildFlags.String("theme-mode", "", "Theme mode: light, dark, auto")
	)
	buildFlags.Parse(os.Args[2:]) // Skip "build" command
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
	// builds can run without a local docs directory.
	configuredAbsDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("resolve dir: %v", err)
	}

	// Resolve output directory
	absOutput, err := filepath.Abs(*output)
	if err != nil {
		log.Fatalf("resolve output dir: %v", err)
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
	} else if *repo != "" {
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
	s, err := site.New(absDir, codeTheme, prefix, "")
	if err != nil {
		log.Fatalf("init site: %v", err)
	}

	var ghClient github.ClientAPI
	var repoInfo *github.RepositoryInfo
	if cfg.GitHub.Enabled && cfg.GitHub.Repository != "" {
		ghClient, repoInfo, err = setupGitHubIntegration(absRootDir, cfg)
		if err != nil {
			log.Fatalf("setup GitHub integration: %v", err)
		}
		applyRepoToSite(s, ghClient, repoInfo, repoInfo.Path)
		log.Printf("dorcs: note: local docs directory will be ignored when GitHub integration is enabled")
	}
	s.SetExplicitNav(cfg.Nav.Items)
	s.SetSectionsConfigured(len(cfg.Nav.Sections) > 0)
	if err := s.BuildIndex(); err != nil {
		log.Fatalf("build index: %v", err)
	}

	// Parse templates (include presentation.html for slide decks)
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
		GitHubClient: ghClient,
		GitHubRepo:   repoInfo,
	})

	// Build static site
	log.Printf("dorcs: building static site from %s", absDir)
	log.Printf("dorcs: output directory: %s", absOutput)
	if err := builder.Build(!*noDrafts); err != nil {
		log.Fatalf("build failed: %v", err)
	}

	log.Printf("dorcs: static site built successfully!")
}
