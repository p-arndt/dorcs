package site

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	"dorcs-v2/internal/markdown"
)

// RenderDoc reads and renders the markdown backing the given key.
// It strips front matter and returns HTML.
// If the file changed since indexing, it uses the latest content.
func (s *Site) RenderDoc(key string) (*RenderedDoc, error) {
	doc, ok := s.GetDoc(key)
	if !ok {
		return nil, fs.ErrNotExist
	}

	raw, meta, hash, modTime, err := markdown.ReadMarkdownStripFrontMatter(doc.FilePath)
	if err != nil {
		return nil, err
	}

	// Strip YAML front matter from the markdown body so metadata is not rendered as page content.
	// (We still use the parsed metadata fields from `meta` below.)
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

	// Generate a table of contents from headings (h2/h3 by default).
	toc := markdown.BuildTOC(s.md, raw)

	// Reconcile metadata with fresh read.
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

	var buf bytes.Buffer
	if err := s.md.Convert([]byte(raw), &buf); err != nil {
		return nil, fmt.Errorf("render markdown: %w", err)
	}

	// Convert GitHub-style alert blocks in the HTML output
	htmlOutput := markdown.ConvertAlertBlocksInHTML(buf.String())

	return &RenderedDoc{
		Doc:         &merged,
		HTML:        template.HTML(htmlOutput),
		TocHTML:     toc,
		RawMarkdown: raw,
	}, nil
}
