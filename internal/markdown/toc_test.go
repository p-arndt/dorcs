package markdown

import (
	"html/template"
	"strings"
	"testing"

	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	gmast "github.com/yuin/goldmark/ast"
)

func TestBuildTOC(t *testing.T) {
	md := NewRenderer("github")

	tests := []struct {
		name        string
		input       string
		contains    []string // substrings that should be in TOC
		notContains []string // substrings that should NOT be in TOC
		notEmpty    bool
	}{
		{
			name: "single h2",
			input: `# Main Title

## Section One
Content here.`,
			contains:    []string{"toc-list", "Section One"},
			notContains: []string{"Main Title"}, // Single H1 should be excluded
			notEmpty:    true,
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
			name: "single h1 excluded",
			input: `# H1 Title

## H2 Section
Content.`,
			contains:    []string{"H2 Section"},
			notContains: []string{"H1 Title"}, // Single H1 should be excluded
			notEmpty:    true,
		},
		{
			name: "multiple h1 included",
			input: `# First H1 Title

## H2 Section

# Second H1 Title

## Another H2 Section
Content.`,
			contains: []string{"First H1 Title", "Second H1 Title", "H2 Section", "Another H2 Section"},
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
				for _, substr := range tt.notContains {
					if strings.Contains(tocStr, substr) {
						t.Errorf("TOC should NOT contain %q, got %q", substr, tocStr)
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

func TestProcessTOCPlaceholder(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		tocHTML       string
		rootTOCHTML   string
		expected      string
		expectedFound map[string]bool
	}{
		{
			name:          "TOC on its own line",
			input:         "Some content\n[[TOC]]\nMore content",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "Some content\n<ul><li>TOC</li></ul>\nMore content",
			expectedFound: map[string]bool{"TOC": true},
		},
		{
			name:          "TOC at start of file",
			input:         "[[TOC]]\nContent here",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "<ul><li>TOC</li></ul>\nContent here",
			expectedFound: map[string]bool{"TOC": true},
		},
		{
			name:          "TOC at end of file",
			input:         "Content here\n[[TOC]]",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "Content here\n<ul><li>TOC</li></ul>",
			expectedFound: map[string]bool{"TOC": true},
		},
		{
			name:          "TOC-ROOT on its own line",
			input:         "Some content\n[[TOC-ROOT]]\nMore content",
			tocHTML:       "",
			rootTOCHTML:   "<ul><li>Root TOC</li></ul>",
			expected:      "Some content\n<ul><li>Root TOC</li></ul>\nMore content",
			expectedFound: map[string]bool{"TOC-ROOT": true},
		},
		{
			name:          "Both TOC and TOC-ROOT",
			input:         "Content\n[[TOC]]\nMore\n[[TOC-ROOT]]\nEnd",
			tocHTML:       "<ul><li>Page TOC</li></ul>",
			rootTOCHTML:   "<ul><li>Root TOC</li></ul>",
			expected:      "Content\n<ul><li>Page TOC</li></ul>\nMore\n<ul><li>Root TOC</li></ul>\nEnd",
			expectedFound: map[string]bool{"TOC": true, "TOC-ROOT": true},
		},
		{
			name:          "TOC inside fenced code block - should not be replaced",
			input:         "```\n[[TOC]]\n```",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "```\n[[TOC]]\n```",
			expectedFound: map[string]bool{},
		},
		{
			name:          "TOC-ROOT inside fenced code block - should not be replaced",
			input:         "```\n[[TOC-ROOT]]\n```",
			tocHTML:       "",
			rootTOCHTML:   "<ul><li>Root TOC</li></ul>",
			expected:      "```\n[[TOC-ROOT]]\n```",
			expectedFound: map[string]bool{},
		},
		{
			name:          "TOC inside fenced code block with language - should not be replaced",
			input:         "```markdown\n[[TOC]]\n```",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "```markdown\n[[TOC]]\n```",
			expectedFound: map[string]bool{},
		},
		{
			name:          "TOC with other text on line - should not be replaced",
			input:         "Use [[TOC]] in your markdown",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "Use [[TOC]] in your markdown",
			expectedFound: map[string]bool{},
		},
		{
			name:          "TOC with whitespace on its own line - should be replaced",
			input:         "  [[TOC]]  ",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "<ul><li>TOC</li></ul>",
			expectedFound: map[string]bool{"TOC": true},
		},
		{
			name:          "Multiple TOC placeholders on separate lines",
			input:         "[[TOC]]\nContent\n[[TOC]]",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "<ul><li>TOC</li></ul>\nContent\n<ul><li>TOC</li></ul>",
			expectedFound: map[string]bool{"TOC": true},
		},
		{
			name:          "TOC in tilde code block - should not be replaced",
			input:         "~~~\n[[TOC]]\n~~~",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "~~~\n[[TOC]]\n~~~",
			expectedFound: map[string]bool{},
		},
		{
			name:          "TOC outside code block after code block",
			input:         "```\ncode here\n```\n[[TOC]]\nMore content",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "```\ncode here\n```\n<ul><li>TOC</li></ul>\nMore content",
			expectedFound: map[string]bool{"TOC": true},
		},
		{
			name:          "Empty TOC HTML - should remove placeholder",
			input:         "Content\n[[TOC]]\nMore",
			tocHTML:       "",
			rootTOCHTML:   "",
			expected:      "Content\n\nMore",
			expectedFound: map[string]bool{"TOC": true},
		},
		{
			name:          "Nested code blocks",
			input:         "```\nouter\n```\n[[TOC]]\n```\ninner\n```",
			tocHTML:       "<ul><li>TOC</li></ul>",
			rootTOCHTML:   "",
			expected:      "```\nouter\n```\n<ul><li>TOC</li></ul>\n```\ninner\n```",
			expectedFound: map[string]bool{"TOC": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := ProcessTOCPlaceholder(tt.input, template.HTML(tt.tocHTML), template.HTML(tt.rootTOCHTML))
			if result != tt.expected {
				t.Errorf("ProcessTOCPlaceholder() = %q\nwant %q", result, tt.expected)
			}
			// Check found map
			if len(found) != len(tt.expectedFound) {
				t.Errorf("ProcessTOCPlaceholder() found = %v\nwant %v", found, tt.expectedFound)
			}
			for k, v := range tt.expectedFound {
				if found[k] != v {
					t.Errorf("ProcessTOCPlaceholder() found[%q] = %v\nwant %v", k, found[k], v)
				}
			}
		})
	}
}
