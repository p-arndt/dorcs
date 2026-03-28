package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/markdown"
	"github.com/p-arndt/dorcs/internal/syntax"
)

// New creates a Site serving markdown documents from rootDir.
// codeTheme specifies the Chroma syntax highlighting theme to use (e.g., "github", "dracula").
// basePath is the URL path prefix (e.g., "/docs"). Empty string means no prefix.
// language is the language code (e.g., "en", "de"). Empty string means default language (root docs folder).
// For non-default languages, documents are expected in rootDir/{language}/ folder (MkDocs-style).
// version is the version identifier (e.g., "v1", "v2", "latest"). Empty string means default version.
// For non-default versions, documents are expected in rootDir/{version}/ or rootDir/{language}/{version}/ folder.
func New(rootDir string, codeTheme string, basePath string, language string) (*Site, error) {
	return NewWithVersion(rootDir, codeTheme, basePath, language, "")
}

// NewWithVersion creates a Site serving markdown documents from rootDir with version support.
// codeTheme specifies the Chroma syntax highlighting theme to use (e.g., "github", "dracula").
// basePath is the URL path prefix (e.g., "/docs"). Empty string means no prefix.
// language is the language code (e.g., "en", "de"). Empty string means default language (root docs folder).
// version is the version identifier (e.g., "v1", "v2", "latest"). Empty string means default version.
// Uses MkDocs-style structure: language-first (docs/{lang}/{version}/) or version-only (docs/{version}/).
func NewWithVersion(rootDir string, codeTheme string, basePath string, language string, version string) (*Site, error) {
	return NewWithVersionPath(rootDir, codeTheme, basePath, language, version, "")
}

// NewWithVersionPath creates a Site serving markdown documents from rootDir with version and custom path support.
// codeTheme specifies the Chroma syntax highlighting theme to use (e.g., "github", "dracula").
// basePath is the URL path prefix (e.g., "/docs"). Empty string means no prefix.
// language is the language code (e.g., "en", "de"). Empty string means default language (root docs folder).
// version is the version identifier (e.g., "v1", "v2", "latest"). Empty string means default version.
// versionPath is an optional custom path override for the version.
// Uses MkDocs-style structure: language-first (docs/{lang}/{version}/) or version-only (docs/{version}/).
func NewWithVersionPath(rootDir string, codeTheme string, basePath string, language string, version string, versionPath string) (*Site, error) {
	if strings.TrimSpace(rootDir) == "" {
		return nil, errors.New("rootDir is required")
	}
	if codeTheme == "" {
		codeTheme = "github" // fallback
	}
	abs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("abs rootDir: %w", err)
	}

	actualRootDir := abs

	// MkDocs-style structure: language-first approach
	// Structure: docs/{lang}/{version}/ or docs/{version}/ (version-only) or docs/ (default)

	// If language is specified, look for language folder first
	if language != "" {
		langDir := filepath.Join(abs, language)
		if stat, err := os.Stat(langDir); err == nil && stat.IsDir() {
			actualRootDir = langDir

			// If version or explicit versionPath is specified, look for version folder inside language folder.
			versionDirName := version
			if versionPath != "" {
				versionDirName = versionPath
			}
			if versionDirName != "" {
				versionDir := filepath.Join(langDir, versionDirName)
				if stat, err := os.Stat(versionDir); err == nil && stat.IsDir() {
					actualRootDir = versionDir
				}
				// If version folder doesn't exist, stay in language folder (default version for that language)
			}
		}
		// If language folder doesn't exist, we'll use the base directory
		// This allows GitHub-only mode where language directories might not exist locally
	} else if version != "" || versionPath != "" {
		// Version-only (no language): docs/{version}/ or explicit versionPath
		if versionPath != "" {
			// Use custom path
			versionDir := filepath.Join(abs, versionPath)
			if stat, err := os.Stat(versionDir); err == nil && stat.IsDir() {
				actualRootDir = versionDir
			}
		} else {
			// Direct version folder
			versionDir := filepath.Join(abs, version)
			if stat, err := os.Stat(versionDir); err == nil && stat.IsDir() {
				actualRootDir = versionDir
			}
		}
	}
	// If neither language nor version specified, use root (default language, default version)

	// Validate the base rootDir exists
	stat, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("stat rootDir: %w", err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("rootDir is not a directory: %s", abs)
	}

	// Generate syntax highlighting CSS
	syntaxCSS := syntax.GenerateCSS(codeTheme)

	return &Site{
		RootDir:         actualRootDir, // Use actual version/language-specific directory
		BasePath:        basePath,
		Language:        language,
		Version:         version,
		DefaultVersion:  "", // Will be set by SetDefaultVersion if needed
		DefaultLanguage: "", // Will be set by SetDefaultLanguage if needed
		md:              markdown.NewRenderer(codeTheme),
		index:           make(map[string]*Doc),
		nav:             &NavNode{Name: "", Key: "", IsDir: true},
		syntaxCSS:       syntaxCSS,
	}, nil
}

// SyntaxCSS returns the generated CSS for syntax highlighting.
// This should be served as a static CSS file or embedded in the page.
func (s *Site) SyntaxCSS() string {
	return s.syntaxCSS
}

// SetGitHubConfig configures GitHub integration for this site.
// client is the GitHub client, and owner, repo, branch, path specify the repository location.
func (s *Site) SetGitHubConfig(client interface {
	DiscoverMarkdownFiles(owner, repo, branch, rootPath string) ([]string, error)
	FetchMarkdown(owner, repo, branch, filePath string) ([]byte, error)
}, owner, repo, branch, repoPath string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.githubClient = client
	s.githubOwner = owner
	s.githubRepo = repo
	s.githubBranch = branch
	s.githubPath = repoPath
}

// FetchGitHubAsset fetches a non-markdown asset from the configured GitHub repository.
func (s *Site) FetchGitHubAsset(relPath string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.githubClient == nil {
		return nil, errors.New("GitHub client not configured")
	}

	fullPath := strings.Trim(relPath, "/")
	if s.githubPath != "" {
		fullPath = strings.Trim(s.githubPath+"/"+fullPath, "/")
	}

	return s.githubClient.FetchMarkdown(s.githubOwner, s.githubRepo, s.githubBranch, fullPath)
}

// SetLanguage sets the language code for this site.
// This is useful when creating a site with GitHub integration where the local directory structure doesn't match.
func (s *Site) SetLanguage(language string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Language = language
}

// SetVersion sets the version identifier for this site.
// This is useful when creating a site with GitHub integration where the local directory structure doesn't match.
func (s *Site) SetVersion(version string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Version = version
}

// SetDefaultVersion sets the default version identifier for this site.
// This is used for link rewriting to determine if version prefix should be added.
func (s *Site) SetDefaultVersion(defaultVersion string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DefaultVersion = defaultVersion
}

// SetDefaultLanguage sets the default language code for this site.
// This is used for link rewriting to determine if language prefix should be added.
func (s *Site) SetDefaultLanguage(defaultLanguage string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DefaultLanguage = defaultLanguage
}

// SetExplicitNav configures an optional explicit navigation tree from dorcs.yaml.
func (s *Site) SetExplicitNav(items config.NavItems) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.explicitNav = append(config.NavItems(nil), items...)
}

// SetSectionsConfigured marks whether nav.sections is active.
func (s *Site) SetSectionsConfigured(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sectionsConfigured = v
}

// ExplicitNav returns a snapshot of the current explicit navigation config.
func (s *Site) ExplicitNav() config.NavItems {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append(config.NavItems(nil), s.explicitNav...)
}
