package markdown

import (
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// NewRenderer creates a configured goldmark instance for parsing and rendering markdown.
func NewRenderer(codeTheme string) goldmark.Markdown {
	if codeTheme == "" {
		codeTheme = "github" // fallback
	}

	chromaStyle := styles.Get(codeTheme)
	if chromaStyle == nil {
		chromaStyle = styles.Fallback
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // tables, strikethrough, task list, etc.
			extension.Footnote,
			extension.Typographer,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle(codeTheme),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
					chromahtml.WithLineNumbers(false),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithXHTML(),
			gmhtml.WithUnsafe(), // allow embedded HTML in markdown
		),
	)

	return md
}
