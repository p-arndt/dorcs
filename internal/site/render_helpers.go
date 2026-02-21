package site

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/p-arndt/dorcs/internal/markdown"
	markdownext "github.com/p-arndt/dorcs/internal/markdown/extensions"
)

// preprocessMarkdown prepares markdown for rendering by:
// - Stripping YAML front matter from body
// - Rewriting extensionless doc links
// - Converting GitHub-style alert blocks
// - Converting accordion blocks
// - Running plugin preprocessing hooks
func (s *Site) preprocessMarkdown(raw string, doc *Doc) string {
	// Strip YAML front matter from the markdown body so metadata is not rendered as page content.
	raw = markdown.StripYAMLFrontMatter(raw)

	// Support extensionless doc links in markdown:
	// [Getting Started](getting-started) -> /getting-started (or /basepath/getting-started if basepath is set)
	// Also resolves relative links like ./explain.md based on the current document's directory
	// For index.md files, use the document's Key as the directory (e.g., guide/index.md uses "guide")
	// For regular files, use DirKey (e.g., guide/intro.md uses "guide")
	docDir := doc.DirKey
	if isIndexRel(doc.RelPath) {
		// This is an index.md file, so the document's directory is its Key
		docDir = doc.Key
	}

	// For GitHub documents, use the directory of the GitHub file path for relative path resolution
	if doc.IsGitHub && doc.GitHubPath != "" {
		githubDir := dirKeyFromKey(keyFromRel(doc.GitHubPath))
		// Use GitHub path directory for relative resolution
		docDir = githubDir
	}

	defaultVersion := s.DefaultVersion
	if defaultVersion == "" {
		// If DefaultVersion is not set, assume current version is default (backward compatibility)
		defaultVersion = s.Version
	}
	defaultLanguage := s.DefaultLanguage
	// For link rewriting, use empty string for language if it's the default (so URLs don't have /en/ prefix)
	languageForURL := s.Language
	if languageForURL == defaultLanguage {
		languageForURL = ""
	}
	raw = markdown.RewriteExtensionlessDocLinks(raw, docDir, s.BasePath, languageForURL, s.Version, defaultVersion, defaultLanguage)

	// Rewrite relative image paths
	// For GitHub-sourced docs: resolve to raw.githubusercontent.com URLs (images live in the repo)
	// For local docs: resolve to absolute site paths
	if doc.IsGitHub && s.githubClient != nil && doc.GitHubPath != "" {
		githubDir := path.Dir(doc.GitHubPath)
		raw = markdown.RewriteRelativeImagePathsForGitHub(raw, githubDir, s.githubOwner, s.githubRepo, s.githubBranch)
	} else {
		raw = markdown.RewriteRelativeImagePaths(raw, docDir, s.BasePath, languageForURL, s.Version, defaultVersion, defaultLanguage)
	}

	// Convert GitHub-style alert blocks in markdown (pre-process for goldmark)
	raw = markdownext.ConvertAlertBlocksInMarkdown(raw)

	// Convert accordion blocks in markdown (pre-process for goldmark)
	raw = markdownext.ConvertAccordionBlocksInMarkdown(raw)

	// Convert timeline blocks in markdown (pre-process for goldmark)
	raw = markdownext.ConvertTimelineBlocksInMarkdown(raw)

	// Convert typography text blocks in markdown (pre-process for goldmark)
	raw = markdownext.ConvertTextBlocksInMarkdown(raw)

	return raw
}

// reconcileMetadata merges front matter metadata with the document.
func (s *Site) reconcileMetadata(doc *Doc, meta markdown.FrontMatter, hash string, modTime time.Time) *Doc {
	merged := *doc

	if t := strings.TrimSpace(meta.Title); t != "" {
		merged.Title = t
	}
	if ds := strings.TrimSpace(meta.Description); ds != "" {
		merged.Description = ds
	}
	if len(meta.Tags) > 0 {
		merged.Tags = append([]string(nil), meta.Tags...)
	}
	merged.Draft = meta.Draft
	merged.Order = meta.Order
	merged.Presentation = meta.Presentation
	merged.PresentationHeader = meta.PresentationHeader
	merged.PresentationFooter = meta.PresentationFooter
	if a := strings.TrimSpace(meta.Author); a != "" {
		merged.Author = a
	}
	if ds := strings.TrimSpace(meta.Date); ds != "" {
		if t, ok := parseDate(ds); ok {
			merged.Date = t
		}
	}
	merged.UpdatedAt = modTime
	merged.ContentHash = hash
	if merged.Title == "" {
		merged.Title = titleFromKey(merged.Key)
	}

	return &merged
}

// convertMarkdownToHTML converts markdown to HTML and processes alert blocks, accordion blocks, and plugin post-processing.
func (s *Site) convertMarkdownToHTML(raw string, doc *Doc) (string, error) {
	var buf bytes.Buffer
	if err := s.md.Convert([]byte(raw), &buf); err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}

	htmlOutput := buf.String()

	// Convert GitHub-style alert blocks in the HTML output
	htmlOutput = markdownext.ConvertAlertBlocksInHTML(htmlOutput)

	// Convert accordion blocks in the HTML output
	htmlOutput = markdownext.ConvertAccordionBlocksInHTML(htmlOutput)

	// Convert timeline blocks in the HTML output
	htmlOutput = markdownext.ConvertTimelineBlocksInHTML(htmlOutput)

	// Convert typography text blocks in the HTML output
	htmlOutput = markdownext.ConvertTextBlocksInHTML(htmlOutput)

	return htmlOutput, nil
}
