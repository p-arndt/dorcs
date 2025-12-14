package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
)

func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name      string
		codeTheme string
		check     func(*testing.T, goldmark.Markdown)
	}{
		{
			name:      "github theme",
			codeTheme: "github",
			check: func(t *testing.T, md goldmark.Markdown) {
				if md == nil {
					t.Error("Renderer should not be nil")
				}
			},
		},
		{
			name:      "empty theme defaults to github",
			codeTheme: "",
			check: func(t *testing.T, md goldmark.Markdown) {
				if md == nil {
					t.Error("Renderer should not be nil")
				}
			},
		},
		{
			name:      "dracula theme",
			codeTheme: "dracula",
			check: func(t *testing.T, md goldmark.Markdown) {
				if md == nil {
					t.Error("Renderer should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := NewRenderer(tt.codeTheme)
			tt.check(t, md)
		})
	}
}

func TestRendererConvertsMarkdown(t *testing.T) {
	md := NewRenderer("github")

	tests := []struct {
		name     string
		input    string
		contains []string // substrings that should be in output
	}{
		{
			name:     "heading",
			input:    "# Hello World",
			contains: []string{"h1", "Hello World"},
		},
		{
			name:     "paragraph",
			input:    "This is a paragraph.",
			contains: []string{"<p>", "This is a paragraph"},
		},
		{
			name:     "bold text",
			input:    "**bold** text",
			contains: []string{"<strong>", "bold"},
		},
		{
			name:     "code block",
			input:    "```go\nfunc main() {}\n```",
			contains: []string{"chroma", "func", "main"},
		},
		{
			name:     "link",
			input:    "[Link](https://example.com)",
			contains: []string{"<a", "href", "example.com"},
		},
		{
			name:     "list",
			input:    "- Item 1\n- Item 2",
			contains: []string{"<ul>", "<li>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.input), &buf); err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			output := buf.String()
			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("Output should contain %q, got %q", substr, output)
				}
			}
		})
	}
}
