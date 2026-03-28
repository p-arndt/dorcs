package markdown

import (
	"strings"
	"testing"
)

func TestConvertTabBlocksInMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "basic tabs",
			input: `:::tabs
::tab macOS
Install with brew
::tab Linux
Install with apt
:::`,
			contains: []string{
				"DORCS_TABS_START",
				"DORCS_TAB:macOS",
				"DORCS_TAB:Linux",
				"Install with brew",
				"Install with apt",
				"DORCS_TABS_END",
			},
		},
		{
			name: "tabs inside code block are ignored",
			input: "```\n:::tabs\n::tab A\ncontent\n:::\n```",
			contains: []string{
				":::tabs",
			},
		},
		{
			name:  "no tabs",
			input: "Just regular markdown\n\nWith paragraphs",
			contains: []string{
				"Just regular markdown",
			},
		},
		{
			name: "multiple tab groups",
			input: `:::tabs
::tab A
Content A
:::

:::tabs
::tab B
Content B
:::`,
			contains: []string{
				"DORCS_TAB:A",
				"DORCS_TAB:B",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTabBlocksInMarkdown(tt.input)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("result should contain %q, got:\n%s", want, result)
				}
			}
		})
	}
}

func TestConvertTabBlocksInHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:  "converts tab markers to tab HTML",
			input: `<blockquote><p><strong>DORCS_TABS_START</strong></p><p><strong>DORCS_TAB:macOS</strong></p><p>Install with brew</p><p><strong>DORCS_TAB:Linux</strong></p><p>Install with apt</p><p><strong>DORCS_TABS_END</strong></p></blockquote>`,
			contains: []string{
				`class="dorcs-tabs"`,
				`class="dorcs-tab-btn active"`,
				`class="dorcs-tab-panel active"`,
				"macOS",
				"Linux",
				"Install with brew",
				"Install with apt",
			},
		},
		{
			name:  "no markers - unchanged",
			input: "<p>Hello world</p>",
			contains: []string{
				"<p>Hello world</p>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTabBlocksInHTML(tt.input)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("result should contain %q, got:\n%s", want, result)
				}
			}
		})
	}
}
