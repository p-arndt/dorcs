package markdown

import "testing"

func TestStripYAMLFrontMatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with front matter",
			input:    "---\ntitle: Test\n---\n# Content",
			expected: "# Content",
		},
		{
			name:     "no front matter",
			input:    "# Just Content",
			expected: "# Just Content",
		},
		{
			name:     "empty front matter",
			input:    "---\n---\n# Content",
			expected: "# Content",
		},
		{
			name:     "unclosed front matter",
			input:    "---\ntitle: Test\n# Content",
			expected: "---\ntitle: Test\n# Content",
		},
		{
			name:     "front matter with CRLF",
			input:    "---\r\ntitle: Test\r\n---\r\n# Content",
			expected: "# Content",
		},
		{
			name:     "with BOM",
			input:    "\uFEFF---\ntitle: Test\n---\n# Content",
			expected: "# Content",
		},
		{
			name:     "only delimiters",
			input:    "---",
			expected: "---",
		},
		{
			name:     "multiline front matter",
			input:    "---\ntitle: Test\ndescription: A test\n---\n# Content",
			expected: "# Content",
		},
		{
			name:     "front matter with empty lines",
			input:    "---\n\ntitle: Test\n\n---\n# Content",
			expected: "# Content",
		},
		{
			name:     "no newline after closing delimiter",
			input:    "---\ntitle: Test\n---\n# Content",
			expected: "# Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripYAMLFrontMatter(tt.input)
			if result != tt.expected {
				t.Errorf("StripYAMLFrontMatter() = %q; want %q", result, tt.expected)
			}
		})
	}
}
