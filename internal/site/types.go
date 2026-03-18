package site

import (
	"html/template"
	"sync"
	"time"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/yuin/goldmark"
)

// Site indexes a directory of markdown documents and renders them to HTML.
// It is safe for concurrent use.
type Site struct {
	RootDir string

	// BasePath is the URL path prefix (e.g., "/docs"). Empty string means no prefix.
	BasePath string

	// Language is the language code this site instance serves (empty for default language)
	Language string

	// Version is the version identifier this site instance serves (empty for default version)
	Version string

	// DefaultVersion is the default version identifier (used for link rewriting)
	DefaultVersion string

	// DefaultLanguage is the default language code (used for link rewriting - no prefix for default)
	DefaultLanguage string

	md goldmark.Markdown

	mu    sync.RWMutex
	index map[string]*Doc // key: normalized doc key ("" for root index), posix-style

	// nav is a cached navigation tree derived from index.
	nav *NavNode

	// explicitNav holds an optional user-defined navigation tree from config.
	explicitNav config.NavItems

	// syntaxCSS holds the generated Chroma CSS for syntax highlighting.
	syntaxCSS string

	// GitHub client and configuration (optional, for GitHub integration)
	githubClient interface {
		DiscoverMarkdownFiles(owner, repo, branch, rootPath string) ([]string, error)
		FetchMarkdown(owner, repo, branch, filePath string) ([]byte, error)
	}
	githubOwner  string
	githubRepo   string
	githubBranch string
	githubPath   string
}

// Doc represents a single markdown document on disk.
type Doc struct {
	// Key is the normalized URL key.
	//
	// Routing rules (index.md):
	// - "<root>/index.md"           -> Key = ""        (served at "/")
	// - "<root>/guide/index.md"     -> Key = "guide"   (served at "/guide")
	// - "<root>/guide/intro.md"     -> Key = "guide/intro" (served at "/guide/intro")
	Key string

	// FilePath is the absolute path to the backing .md file.
	FilePath string

	// RelPath is the path relative to Site.RootDir using "/" separators, including ".md".
	RelPath string

	// DirKey is the Key of the directory (without trailing slash) for grouping, e.g. "getting-started".
	DirKey string

	Title              string
	Description        string
	Date               time.Time
	Tags               []string
	Draft              bool
	Order              int    // Order for sorting (lower numbers appear first)
	Author             string // Author name
	After              string // Key of item this should appear after (use "index" for index.md)
	Presentation       bool   // If true, render as Marp-style slide deck
	PresentationHeader string // Header text on each slide
	PresentationFooter string // Footer text on each slide

	// UpdatedAt is the file modtime.
	UpdatedAt time.Time

	// ContentHash is a hash of the file contents for caching/invalidation.
	ContentHash string

	// IsGitHub indicates if this document is sourced from GitHub.
	IsGitHub bool

	// GitHubPath is the path to the file in the GitHub repository (for GitHub-sourced docs).
	GitHubPath string

	// GitHubCacheKey is the cache key for GitHub content (owner/repo/branch/path).
	GitHubCacheKey string
}

// NavNode is a directory tree node for sidebar navigation.
type NavNode struct {
	// Name is the display name (e.g. "guide", "API").
	Name string

	// ExplicitTitle indicates Name came from explicit navigation config and should override Page.Title.
	ExplicitTitle bool

	// Key is the URL key for this folder/page.
	// For folders, this is the folder key ("guide"); for pages, it is the doc key ("guide/getting-started").
	Key string

	// IsDir indicates whether this node is a directory.
	IsDir bool

	// Page points to the folder landing page (index.md) when IsDir is true and such a doc exists.
	Page *Doc

	// Children contains subdirectories (as NavNode with IsDir=true) and pages (IsDir=false).
	Children []*NavNode
}

// RenderedSlide is a single slide with Marpit-compatible layout metadata.
type RenderedSlide struct {
	HTML               template.HTML
	Class              string
	Layout             string // Layout preset: default, lead, left, right, columns-2, columns-3, timeline, split, cols
	Gap                string // Spacing: tight, normal, loose
	Align              string // Content alignment: start, center, end
	Color              string // CSS text color
	BackgroundColor    string
	BackgroundImage    string
	BackgroundPosition string
	BackgroundRepeat   string
	BackgroundSize     string
	Header             string
	Footer             string
	Paginate           string // "true", "false" - show/hide slide number
}

// RenderedDoc is the result of rendering a markdown file to HTML.
type RenderedDoc struct {
	Doc         *Doc
	HTML        template.HTML
	TocHTML     template.HTML   // reserved for future TOC support
	RawMarkdown string          // markdown prior to conversion (may include front matter)
	Slides      []RenderedSlide // Populated when Doc.Presentation is true
}

// SearchResult represents a single search result.
type SearchResult struct {
	Doc         *Doc
	Key         string
	Title       string
	Path        string
	Snippet     string // Excerpt from content showing the match
	Score       int    // Relevance score (higher is better)
	HeadingID   string // ID of the heading where the match was found (for anchor links)
	HeadingText string // Text of the heading where the match was found
}
