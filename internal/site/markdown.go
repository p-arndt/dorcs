package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dorcs-v2/internal/markdown"
	"dorcs-v2/internal/syntax"
)

// New creates a Site serving markdown documents from rootDir.
// codeTheme specifies the Chroma syntax highlighting theme to use (e.g., "github", "dracula").
// basePath is the URL path prefix (e.g., "/docs"). Empty string means no prefix.
func New(rootDir string, codeTheme string, basePath string) (*Site, error) {
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
		RootDir:   abs,
		BasePath:  basePath,
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
