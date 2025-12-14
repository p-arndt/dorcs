package site

import (
	"html/template"
	"io/fs"

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

	// Read and parse markdown file
	raw, meta, hash, modTime, err := markdown.ReadMarkdownStripFrontMatter(doc.FilePath)
	if err != nil {
		return nil, err
	}

	// Preprocess markdown (strip front matter, rewrite links, convert alerts)
	raw = s.preprocessMarkdown(raw, doc)

	// Generate and process TOC
	toc, _, raw := s.generateAndProcessTOC(raw, doc, key)

	// Reconcile metadata with fresh read
	merged := s.reconcileMetadata(doc, meta, hash, modTime)

	// Convert markdown to HTML
	htmlOutput, err := s.convertMarkdownToHTML(raw)
	if err != nil {
		return nil, err
	}

	return &RenderedDoc{
		Doc:         merged,
		HTML:        template.HTML(htmlOutput),
		TocHTML:     toc,
		RawMarkdown: raw,
	}, nil
}
