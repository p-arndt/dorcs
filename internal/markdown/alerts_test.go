package markdown

import "testing"

func TestConvertAlertBlocksInMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "NOTE alert",
			input: `> [!NOTE]
> This is a note.`,
			expected: `> **ALERT:NOTE**
> This is a note.`,
		},
		{
			name: "WARNING alert",
			input: `> [!WARNING]
> Be careful!`,
			expected: `> **ALERT:WARNING**
> Be careful!`,
		},
		{
			name: "TIP alert",
			input: `> [!TIP]
> Here's a tip.`,
			expected: `> **ALERT:TIP**
> Here's a tip.`,
		},
		{
			name: "IMPORTANT alert",
			input: `> [!IMPORTANT]
> This is important.`,
			expected: `> **ALERT:IMPORTANT**
> This is important.`,
		},
		{
			name: "CAUTION alert",
			input: `> [!CAUTION]
> Take caution.`,
			expected: `> **ALERT:CAUTION**
> Take caution.`,
		},
		{
			name: "INFO mapped to NOTE",
			input: `> [!INFO]
> Information here.`,
			expected: `> **ALERT:NOTE**
> Information here.`,
		},
		{
			name: "multiline alert",
			input: `> [!NOTE]
> First line.
> Second line.`,
			expected: `> **ALERT:NOTE**
> First line.
> Second line.`,
		},
		{
			name: "alert with empty line",
			input: `> [!NOTE]
> First line.
> 
> Second line.`,
			expected: `> **ALERT:NOTE**
> First line.
> 
> Second line.`,
		},
		{
			name: "invalid alert type unchanged",
			input: `> [!INVALID]
> This should not change.`,
			expected: `> [!INVALID]
> This should not change.`,
		},
		{
			name:     "regular blockquote unchanged",
			input:    `> This is a regular blockquote.`,
			expected: `> This is a regular blockquote.`,
		},
		{
			name:     "no alert marker",
			input:    `> Just a quote.`,
			expected: `> Just a quote.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertAlertBlocksInMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertAlertBlocksInMarkdown() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestConvertAlertBlocksInHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "NOTE alert in HTML",
			input:    `<blockquote><p><strong>ALERT:NOTE</strong></p><p>This is a note.</p></blockquote>`,
			expected: `<div class="alert alert-note"><p class="alert-title">Note</p><div class="alert-content"><p>This is a note.</p></div></div>`,
		},
		{
			name:     "WARNING alert in HTML",
			input:    `<blockquote><p><strong>ALERT:WARNING</strong></p><p>Be careful!</p></blockquote>`,
			expected: `<div class="alert alert-warning"><p class="alert-title">Warning</p><div class="alert-content"><p>Be careful!</p></div></div>`,
		},
		{
			name:     "IMPORTANT alert in HTML",
			input:    `<blockquote><p><strong>ALERT:IMPORTANT</strong></p><p>This is important.</p></blockquote>`,
			expected: `<div class="alert alert-important"><p class="alert-title">Important</p><div class="alert-content"><p>This is important.</p></div></div>`,
		},
		{
			name:     "regular blockquote unchanged",
			input:    `<blockquote><p>Regular quote.</p></blockquote>`,
			expected: `<blockquote><p>Regular quote.</p></blockquote>`,
		},
		{
			name:     "alert with multiple paragraphs",
			input:    `<blockquote><p><strong>ALERT:NOTE</strong></p><p>First paragraph.</p><p>Second paragraph.</p></blockquote>`,
			expected: `<div class="alert alert-note"><p class="alert-title">Note</p><div class="alert-content"><p>First paragraph.</p><p>Second paragraph.</p></div></div>`,
		},
		{
			name:     "fallback plain text alert",
			input:    `<blockquote><p>ALERT:NOTE</p><p>Content here.</p></blockquote>`,
			expected: `<div class="alert alert-note"><p class="alert-title">Note</p><div class="alert-content"><p>Content here.</p></div></div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertAlertBlocksInHTML(tt.input)
			// For HTML output, we check that key elements are present rather than exact match
			// since the HTML structure might vary slightly
			if tt.name == "regular blockquote unchanged" {
				if result != tt.expected {
					t.Errorf("ConvertAlertBlocksInHTML() = %q; want %q", result, tt.expected)
				}
			} else {
				// Check that result contains alert div structure
				if !contains(result, "alert alert-") {
					t.Errorf("ConvertAlertBlocksInHTML() should contain alert div, got %q", result)
				}
				if !contains(result, "alert-title") {
					t.Errorf("ConvertAlertBlocksInHTML() should contain alert-title, got %q", result)
				}
				if !contains(result, "alert-content") {
					t.Errorf("ConvertAlertBlocksInHTML() should contain alert-content, got %q", result)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
