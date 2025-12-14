package site

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"dorcs-v2/internal/markdown"
)

// preprocessMarkdown prepares markdown for rendering by:
// - Stripping YAML front matter from body
// - Rewriting extensionless doc links
// - Converting GitHub-style alert blocks
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
	raw = markdown.RewriteExtensionlessDocLinks(raw, docDir, s.BasePath)

	// Convert GitHub-style alert blocks in markdown (pre-process for goldmark)
	raw = markdown.ConvertAlertBlocksInMarkdown(raw)

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

// convertMarkdownToHTML converts markdown to HTML and processes alert blocks.
func (s *Site) convertMarkdownToHTML(raw string) (string, error) {
	var buf bytes.Buffer
	if err := s.md.Convert([]byte(raw), &buf); err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}

	// Convert GitHub-style alert blocks in the HTML output
	htmlOutput := markdown.ConvertAlertBlocksInHTML(buf.String())

	return htmlOutput, nil
}
