package markdown

import (
	"strings"
	"testing"
)

func TestConvertTextBlocksInMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "hero block",
			input: `::: hero
35%
:::`,
			expected: "TEXTBLOCK:hero",
		},
		{
			name: "stat block",
			input: `::: stat
35%
of an audience's retention rate.
:::`,
			expected: "TEXTBLOCK_STAT",
		},
		{
			name: "caption block",
			input: `::: caption
Supporting text here.
:::`,
			expected: "TEXTBLOCK:caption",
		},
		{
			name: "label block",
			input: `::: label
AGENDA
:::`,
			expected: "TEXTBLOCK:label",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTextBlocksInMarkdown(tt.input)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("ConvertTextBlocksInMarkdown() = %q; want to contain %q", result, tt.expected)
			}
		})
	}
}

func TestConvertTextBlocksInHTML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantHero    bool
		wantStat    bool
		wantCaption bool
	}{
		{
			name:     "hero block",
			input:    `<blockquote><p><strong>TEXTBLOCK:hero</strong></p><p>35%</p></blockquote>`,
			wantHero: true,
		},
		{
			name:     "stat block",
			input:    `<blockquote><p><strong>TEXTBLOCK_STAT</strong></p><p>35%</p><p>caption</p></blockquote>`,
			wantStat: true,
		},
		{
			name:        "caption block",
			input:       `<blockquote><p><strong>TEXTBLOCK:caption</strong></p><p>text</p></blockquote>`,
			wantCaption: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTextBlocksInHTML(tt.input)
			if tt.wantHero && !strings.Contains(result, "text-block-hero") {
				t.Errorf("expected text-block-hero, got %q", result)
			}
			if tt.wantStat && !strings.Contains(result, "text-block-stat") {
				t.Errorf("expected text-block-stat, got %q", result)
			}
			if tt.wantCaption && !strings.Contains(result, "text-block-caption") {
				t.Errorf("expected text-block-caption, got %q", result)
			}
		})
	}
}
