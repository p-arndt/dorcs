package markdown

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// ReadFrontMatterAndHash reads a markdown file and extracts front matter and content hash.
func ReadFrontMatterAndHash(path string) (FrontMatter, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return FrontMatter{}, "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return FrontMatter{}, "", err
	}

	sum := sha256.Sum256(b)
	hash := hex.EncodeToString(sum[:])

	// Best-effort: parse front matter using goldmark-meta.
	// If anything fails, return empty metadata but still provide a content hash.
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithXHTML(),
			gmhtml.WithUnsafe(),
		),
	)

	ctx := parser.NewContext()
	var out bytes.Buffer
	if err := md.Convert(b, &out, parser.WithContext(ctx)); err != nil {
		return FrontMatter{}, hash, err
	}

	data := meta.Get(ctx)
	fm := FrontMatter{}

	if v, ok := data["title"].(string); ok {
		fm.Title = v
	}
	if v, ok := data["description"].(string); ok {
		fm.Description = v
	}
	switch v := data["date"].(type) {
	case string:
		fm.Date = v
	case time.Time:
		fm.Date = v.Format(time.RFC3339)
	}
	if v, ok := data["draft"].(bool); ok {
		fm.Draft = v
	}
	// tags may come through as []any; best-effort conversion to []string
	if tv, ok := data["tags"]; ok {
		switch t := tv.(type) {
		case []string:
			fm.Tags = append([]string(nil), t...)
		case []any:
			var tags []string
			for _, x := range t {
				if s, ok := x.(string); ok && strings.TrimSpace(s) != "" {
					tags = append(tags, s)
				}
			}
			fm.Tags = tags
		}
	}
	// order may come through as int or float64 (YAML can parse numbers as float64)
	if ov, ok := data["order"]; ok {
		switch o := ov.(type) {
		case int:
			fm.Order = o
		case int64:
			fm.Order = int(o)
		case float64:
			fm.Order = int(o)
		}
	}
	if v, ok := data["author"].(string); ok {
		fm.Author = v
	}
	if v, ok := data["after"].(string); ok {
		fm.After = v
	}

	return fm, hash, nil
}

// ReadMarkdownStripFrontMatter reads a markdown file and returns the raw content, front matter, hash, and modtime.
func ReadMarkdownStripFrontMatter(path string) (raw string, metaOut FrontMatter, hash string, modTime time.Time, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", FrontMatter{}, "", time.Time{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", FrontMatter{}, "", time.Time{}, err
	}

	sum := sha256.Sum256(b)
	hash = hex.EncodeToString(sum[:])

	// Parse front matter using goldmark-meta.
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithXHTML(),
			gmhtml.WithUnsafe(),
		),
	)

	ctx := parser.NewContext()
	var out bytes.Buffer
	if err := md.Convert(b, &out, parser.WithContext(ctx)); err != nil {
		// If parsing fails, treat the whole file as markdown and return no metadata.
		return string(b), FrontMatter{}, hash, stat.ModTime(), nil
	}

	data := meta.Get(ctx)
	fm := FrontMatter{}

	if v, ok := data["title"].(string); ok {
		fm.Title = v
	}
	if v, ok := data["description"].(string); ok {
		fm.Description = v
	}
	switch v := data["date"].(type) {
	case string:
		fm.Date = v
	case time.Time:
		fm.Date = v.Format(time.RFC3339)
	}
	if v, ok := data["draft"].(bool); ok {
		fm.Draft = v
	}
	if tv, ok := data["tags"]; ok {
		switch t := tv.(type) {
		case []string:
			fm.Tags = append([]string(nil), t...)
		case []any:
			var tags []string
			for _, x := range t {
				if s, ok := x.(string); ok && strings.TrimSpace(s) != "" {
					tags = append(tags, s)
				}
			}
			fm.Tags = tags
		}
	}
	// order may come through as int or float64 (YAML can parse numbers as float64)
	if ov, ok := data["order"]; ok {
		switch o := ov.(type) {
		case int:
			fm.Order = o
		case int64:
			fm.Order = int(o)
		case float64:
			fm.Order = int(o)
		}
	}
	if v, ok := data["author"].(string); ok {
		fm.Author = v
	}
	if v, ok := data["after"].(string); ok {
		fm.After = v
	}

	// NOTE: raw markdown is returned unchanged for now (may still include front matter).
	return string(b), fm, hash, stat.ModTime(), nil
}

// ParseFrontMatterFromContent parses front matter from markdown content bytes.
// Returns the front matter, content hash, and any error.
func ParseFrontMatterFromContent(content []byte) (*FrontMatter, string, error) {
	sum := sha256.Sum256(content)
	hash := hex.EncodeToString(sum[:])

	// Parse front matter using goldmark-meta.
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithXHTML(),
			gmhtml.WithUnsafe(),
		),
	)

	ctx := parser.NewContext()
	var out bytes.Buffer
	if err := md.Convert(content, &out, parser.WithContext(ctx)); err != nil {
		return nil, hash, err
	}

	data := meta.Get(ctx)
	fm := &FrontMatter{}

	if v, ok := data["title"].(string); ok {
		fm.Title = v
	}
	if v, ok := data["description"].(string); ok {
		fm.Description = v
	}
	switch v := data["date"].(type) {
	case string:
		fm.Date = v
	case time.Time:
		fm.Date = v.Format(time.RFC3339)
	}
	if v, ok := data["draft"].(bool); ok {
		fm.Draft = v
	}
	if tv, ok := data["tags"]; ok {
		switch t := tv.(type) {
		case []string:
			fm.Tags = append([]string(nil), t...)
		case []any:
			var tags []string
			for _, x := range t {
				if s, ok := x.(string); ok && strings.TrimSpace(s) != "" {
					tags = append(tags, s)
				}
			}
			fm.Tags = tags
		}
	}
	if ov, ok := data["order"]; ok {
		switch o := ov.(type) {
		case int:
			fm.Order = o
		case int64:
			fm.Order = int(o)
		case float64:
			fm.Order = int(o)
		}
	}
	if v, ok := data["author"].(string); ok {
		fm.Author = v
	}
	if v, ok := data["after"].(string); ok {
		fm.After = v
	}

	return fm, hash, nil
}
