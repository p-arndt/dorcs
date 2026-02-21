package markdown

import (
	"strings"
	"testing"
)

func TestConvertTimelineBlocksInMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple timeline with two steps",
			input: `::: timeline
### 2024 · Jan
**Step 1 Title**
Lorem ipsum dolor sit amet.

### 2024 · Jun
**Step 2 Title**
Consectetuer adipiscing elit.
:::`,
			expected: "TIMELINE_STEP",
		},
		{
			name: "timeline with h4 headings",
			input: `::: timeline
#### Q1 2024
First quarter.

#### Q2 2024
Second quarter.
:::`,
			expected: "TIMELINE_STEP",
		},
		{
			name:     "no timeline block",
			input:    "# Just a heading\n\nSome content.",
			expected: "# Just a heading",
		},
		{
			name: "timeline in code block is preserved",
			input: "```\n::: timeline\n### Step\nContent\n:::\n```",
			expected: "```",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTimelineBlocksInMarkdown(tt.input)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("ConvertTimelineBlocksInMarkdown() = %q; want to contain %q", result, tt.expected)
			}
		})
	}
}

func TestConvertTimelineBlocksInHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTimeline bool
		wantStep bool
	}{
		{
			name: "converts TIMELINE_STEP blockquote",
			input: `<blockquote><p><strong>TIMELINE_STEP:2024 · Jan</strong></p><p><strong>Step 1</strong></p><p>Lorem ipsum.</p></blockquote>`,
			wantTimeline: true,
			wantStep: true,
		},
		{
			name: "multiple steps",
			input: `<blockquote><p><strong>TIMELINE_STEP:A</strong></p><p>Content A</p></blockquote>
<blockquote><p><strong>TIMELINE_STEP:B</strong></p><p>Content B</p></blockquote>`,
			wantTimeline: true,
			wantStep: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTimelineBlocksInHTML(tt.input)
			if tt.wantTimeline && !strings.Contains(result, `class="timeline"`) {
				t.Errorf("ConvertTimelineBlocksInHTML() should contain timeline div, got %q", result)
			}
			if tt.wantStep && !strings.Contains(result, `class="timeline-step"`) {
				t.Errorf("ConvertTimelineBlocksInHTML() should contain timeline-step div, got %q", result)
			}
		})
	}
}
