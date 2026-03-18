package markdown

import (
	"strings"
	"testing"
)

func TestConvertColBlocksInMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, got string)
	}{
		{
			name: "two col blocks converted to COLBLOCK markers",
			input: `::: col
### Left
Content A
:::

::: col
### Right
Content B
:::`,
			check: func(t *testing.T, got string) {
				if strings.Count(got, "**COLBLOCK**") != 2 {
					t.Errorf("expected 2 COLBLOCK markers, got:\n%s", got)
				}
				if strings.Contains(got, "::: col") {
					t.Errorf("expected ::: col to be converted, got:\n%s", got)
				}
			},
		},
		{
			name: "three col blocks",
			input: `::: col
A
:::

::: col
B
:::

::: col
C
:::`,
			check: func(t *testing.T, got string) {
				if strings.Count(got, "**COLBLOCK**") != 3 {
					t.Errorf("expected 3 COLBLOCK markers, got:\n%s", got)
				}
			},
		},
		{
			name: "col blocks inside fenced code block are preserved",
			input: "```markdown\n::: col\nContent\n:::\n```",
			check: func(t *testing.T, got string) {
				if strings.Contains(got, "COLBLOCK") {
					t.Errorf("col block inside fenced code should not be converted, got:\n%s", got)
				}
				if !strings.Contains(got, "::: col") {
					t.Errorf("col block inside fenced code should be preserved, got:\n%s", got)
				}
			},
		},
		{
			name: "unclosed col block preserved as-is",
			input: `::: col
No closing marker`,
			check: func(t *testing.T, got string) {
				if strings.Contains(got, "COLBLOCK") {
					t.Errorf("unclosed col block should not be converted, got:\n%s", got)
				}
				if !strings.Contains(got, "::: col") {
					t.Errorf("unclosed col block should be preserved, got:\n%s", got)
				}
			},
		},
		{
			name:  "non-col blocks are untouched",
			input: "::: hero\nBig text\n:::",
			check: func(t *testing.T, got string) {
				if strings.Contains(got, "COLBLOCK") {
					t.Errorf("non-col blocks should not produce COLBLOCK markers, got:\n%s", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertColBlocksInMarkdown(tt.input)
			tt.check(t, got)
		})
	}
}

func TestConvertColBlocksInHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, got string)
	}{
		{
			name:  "single COLBLOCK marker converted to div.col",
			input: `<blockquote><p><strong>COLBLOCK</strong></p><h3>Title</h3><p>Content</p></blockquote>`,
			check: func(t *testing.T, got string) {
				if !strings.Contains(got, `<div class="col">`) {
					t.Errorf("expected div.col, got:\n%s", got)
				}
				if strings.Contains(got, "COLBLOCK") {
					t.Errorf("COLBLOCK marker should be removed, got:\n%s", got)
				}
				if strings.Contains(got, "<blockquote>") {
					t.Errorf("blockquote should be replaced, got:\n%s", got)
				}
			},
		},
		{
			name: "two COLBLOCK markers both converted",
			input: `<blockquote>
<p><strong>COLBLOCK</strong></p>
<p>Left</p>
</blockquote>
<blockquote>
<p><strong>COLBLOCK</strong></p>
<p>Right</p>
</blockquote>`,
			check: func(t *testing.T, got string) {
				if strings.Count(got, `<div class="col">`) != 2 {
					t.Errorf("expected 2 col divs, got:\n%s", got)
				}
				if strings.Contains(got, "COLBLOCK") {
					t.Errorf("no COLBLOCK markers should remain, got:\n%s", got)
				}
			},
		},
		{
			name:  "non-COLBLOCK blockquote untouched",
			input: `<blockquote><p>Normal quote</p></blockquote>`,
			check: func(t *testing.T, got string) {
				if strings.Contains(got, `class="col"`) {
					t.Errorf("normal blockquote should not become div.col, got:\n%s", got)
				}
				if !strings.Contains(got, "<blockquote>") {
					t.Errorf("normal blockquote should be preserved, got:\n%s", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertColBlocksInHTML(tt.input)
			tt.check(t, got)
		})
	}
}
