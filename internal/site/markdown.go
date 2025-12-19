package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/p-arndt/dorcs/internal/markdown"
	"github.com/p-arndt/dorcs/internal/syntax"
)

// New creates a Site serving markdown documents from rootDir.
// codeTheme specifies the Chroma syntax highlighting theme to use (e.g., "github", "dracula").
// basePath is the URL path prefix (e.g., "/docs"). Empty string means no prefix.
// language is the language code (e.g., "en", "de"). Empty string means default language (root docs folder).
// For non-default languages, documents are expected in rootDir/__lang__/{language}/ folder.
func New(rootDir string, codeTheme string, basePath string, language string) (*Site, error) {
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

	// For non-default languages, use __lang__/{language} folder structure
	actualRootDir := abs
	if language != "" {
		langDir := filepath.Join(abs, "__lang__", language)
		// Check if language directory exists
		if stat, err := os.Stat(langDir); err == nil && stat.IsDir() {
			actualRootDir = langDir
		}
		// If it doesn't exist, we'll use the base directory
		// This allows GitHub-only mode where language directories might not exist locally
	}

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
		RootDir:   actualRootDir, // Use actual language-specific directory
		BasePath:  basePath,
		Language:  language,
		md:        markdown.NewRenderer(codeTheme),
		index:     make(map[string]*Doc),
		nav:       &NavNode{Name: "", Key: "", IsDir: true},
		syntaxCSS: syntaxCSS,
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

// SetLanguage sets the language code for this site.
// This is useful when creating a site with GitHub integration where the local directory structure doesn't match.
func (s *Site) SetLanguage(language string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Language = language
}
