package markdown

import (
	"strings"
	"testing"

	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	gmast "github.com/yuin/goldmark/ast"
)

func TestBuildTOC(t *testing.T) {
	md := NewRenderer("github")

	tests := []struct {
		name     string
		input    string
		contains []string // substrings that should be in TOC
		notEmpty bool
	}{
		{
			name: "single h2",
			input: `# Main Title

## Section One
Content here.`,
			contains: []string{"toc-list", "Section One"},
			notEmpty: true,
		},
		{
			name: "multiple h2",
			input: `## First Section
Content.

## Second Section
More content.`,
			contains: []string{"toc-list", "First Section", "Second Section"},
			notEmpty: true,
		},
		{
			name: "nested h3",
			input: `## Main Section

### Subsection
Content.`,
			contains: []string{"toc-nested", "Subsection"},
			notEmpty: true,
		},
		{
			name: "h2, h3, h4 hierarchy",
			input: `## Section

### Subsection

#### Sub-subsection
Content.`,
			contains: []string{"toc-list", "Section"},
			notEmpty: true,
		},
		{
			name:     "no headings",
			input:    `Just some content without headings.`,
			notEmpty: false,
		},
		{
			name: "h1 ignored",
			input: `# H1 Title

## H2 Section
Content.`,
			contains: []string{"H2 Section"},
			notEmpty: true,
		},
		{
			name: "h5 ignored",
			input: `## H2 Section

##### H5 Section
Content.`,
			contains: []string{"H2 Section"},
			notEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toc := BuildTOC(md, tt.input)

			if tt.notEmpty {
				if string(toc) == "" {
					t.Error("TOC should not be empty")
					return
				}

				tocStr := string(toc)
				for _, substr := range tt.contains {
					if !strings.Contains(tocStr, substr) {
						t.Errorf("TOC should contain %q, got %q", substr, tocStr)
					}
				}
			} else {
				if string(toc) != "" {
					t.Errorf("TOC should be empty, got %q", toc)
				}
			}
		})
	}
}

func TestExtractNodeText(t *testing.T) {
	md := NewRenderer("github")

	tests := []struct {
		name     string
		input    string
		contains string // text that should be extracted
	}{
		{
			name:     "simple heading",
			input:    "## Simple Heading",
			contains: "Simple Heading",
		},
		{
			name:     "heading with bold",
			input:    "## **Bold** Heading",
			contains: "Bold",
		},
		{
			name:     "heading with code",
			input:    "## Code `example`",
			contains: "Code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := []byte(tt.input)
			ctx := parser.NewContext()
			reader := text.NewReader(src)
			doc := md.Parser().Parse(reader, parser.WithContext(ctx))

			// Find first heading
			var heading *gmast.Heading
			_ = gmast.Walk(doc, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
				if entering {
					if h, ok := n.(*gmast.Heading); ok {
						heading = h
						return gmast.WalkStop, nil
					}
				}
				return gmast.WalkContinue, nil
			})

			if heading == nil {
				t.Fatal("No heading found")
			}

			text := ExtractNodeText(heading, src)
			if !strings.Contains(text, tt.contains) {
				t.Errorf("ExtractNodeText() = %q; should contain %q", text, tt.contains)
			}
		})
	}
}
