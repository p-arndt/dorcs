package markdown

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	katex "github.com/brickellis/goldmark-katex"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	mermaid "go.abhg.dev/goldmark/mermaid"
)

// mermaidProtectionTransformer marks mermaid code blocks with a "nohl" attribute
// to prevent the highlighting extension from processing them.
// Runs at priority 50 (before mermaid's transformer at 100).
type mermaidProtectionTransformer struct{}

func (t *mermaidProtectionTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		cb, ok := node.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		lang := cb.Language(source)
		if lang != nil && bytes.Equal(lang, []byte("mermaid")) {
			cb.SetAttribute([]byte("nohl"), []byte("true"))
		}

		return ast.WalkContinue, nil
	})
}

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
			// Mermaid must come before highlighting to prevent interference.
			// The mermaid transformer (priority 100) converts mermaid code blocks
			// to mermaid.Block nodes before highlighting (priority 200) processes them.
			&mermaid.Extender{RenderMode: mermaid.RenderModeClient},
			highlighting.NewHighlighting(
				highlighting.WithStyle(codeTheme),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
					chromahtml.WithLineNumbers(false),
				),
				highlighting.WithGuessLanguage(false),
			),
			&katex.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			// Register mermaid protection transformer at priority 50 to mark
			// mermaid blocks with "nohl" attribute before mermaid transformer runs.
			parser.WithASTTransformers(
				util.Prioritized(&mermaidProtectionTransformer{}, 50),
			),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithXHTML(),
			gmhtml.WithUnsafe(), // allow embedded HTML in markdown
		),
	)

	return md
}
