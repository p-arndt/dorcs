package site

import (
	"fmt"
	"html/template"
	"io/fs"
	"time"

	"github.com/p-arndt/dorcs/internal/markdown"
)

// RenderDoc reads and renders the markdown backing the given key.
// It strips front matter and returns HTML.
// If the file changed since indexing, it uses the latest content.
func (s *Site) RenderDoc(key string) (*RenderedDoc, error) {
	doc, ok := s.GetDoc(key)
	if !ok {
		return nil, fs.ErrNotExist
	}

	var raw []byte
	var meta markdown.FrontMatter
	var hash string
	var modTime time.Time
	var err error

	// Handle GitHub-sourced documents
	if doc.IsGitHub {
		raw, err = s.fetchGitHubMarkdown(doc)
		if err != nil {
			// Return a user-friendly error that can be displayed in the page
			return nil, fmt.Errorf("failed to load content from GitHub: %w. The file may have been moved or deleted, or there may be a network issue.", err)
		}

		// Parse front matter from content
		metaPtr, hashVal, _ := markdown.ParseFrontMatterFromContent(raw)
		if metaPtr != nil {
			meta = *metaPtr
		}
		hash = hashVal
		modTime = time.Now() // GitHub docs use current time
	} else {
		// Read and parse markdown file from local filesystem
		var rawStr string
		rawStr, meta, hash, modTime, err = markdown.ReadMarkdownStripFrontMatter(doc.FilePath)
		if err != nil {
			return nil, err
		}
		raw = []byte(rawStr)
	}

	// Preprocess markdown (strip front matter, rewrite links, convert alerts)
	rawStr := string(raw)
	rawStr = s.preprocessMarkdown(rawStr, doc)

	// Generate and process TOC
	toc, _, rawStr := s.generateAndProcessTOC(rawStr, doc, key)

	// Reconcile metadata with fresh read
	merged := s.reconcileMetadata(doc, meta, hash, modTime)

	// Convert markdown to HTML
	htmlOutput, err := s.convertMarkdownToHTML(rawStr)
	if err != nil {
		return nil, err
	}

	return &RenderedDoc{
		Doc:         merged,
		HTML:        template.HTML(htmlOutput),
		TocHTML:     toc,
		RawMarkdown: rawStr,
	}, nil
}

// fetchGitHubMarkdown fetches markdown content from GitHub.
func (s *Site) fetchGitHubMarkdown(doc *Doc) ([]byte, error) {
	if s.githubClient == nil {
		return nil, fmt.Errorf("GitHub client not configured")
	}

	// Use the full GitHub path stored in the doc
	fullPath := doc.GitHubPath
	if fullPath == "" {
		// Fallback: construct from components
		if s.githubPath != "" {
			fullPath = s.githubPath + "/" + doc.RelPath
		} else {
			fullPath = doc.RelPath
		}
	}

	return s.githubClient.FetchMarkdown(s.githubOwner, s.githubRepo, s.githubBranch, fullPath)
}
