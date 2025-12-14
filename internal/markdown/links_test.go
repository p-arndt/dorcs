package markdown

import "testing"

func TestRewriteExtensionlessDocLinks(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		currentDirKey string
		basePath      string
		expected      string
	}{
		{
			name:          "simple link",
			input:         "[Getting Started](getting-started)",
			currentDirKey: "",
			expected:      "[Getting Started](/getting-started)",
		},
		{
			name:          "link with .md extension",
			input:         "[Guide](guide.md)",
			currentDirKey: "",
			expected:      "[Guide](/guide)",
		},
		{
			name:          "nested link",
			input:         "[API](api/endpoints)",
			currentDirKey: "",
			expected:      "[API](/api/endpoints)",
		},
		{
			name:          "absolute URL unchanged",
			input:         "[Example](https://example.com)",
			currentDirKey: "",
			expected:      "[Example](https://example.com)",
		},
		{
			name:          "anchor link unchanged",
			input:         "[Section](#section)",
			currentDirKey: "",
			expected:      "[Section](#section)",
		},
		{
			name:          "image link unchanged",
			input:         "![Alt](image.png)",
			currentDirKey: "",
			expected:      "![Alt](image.png)",
		},
		{
			name:          "relative link with ./",
			input:         "[Explain](./explain.md)",
			currentDirKey: "user",
			expected:      "[Explain](/user/explain)",
		},
		{
			name:          "folder index link",
			input:         "[Guide](guide/index.md)",
			currentDirKey: "",
			expected:      "[Guide](/guide)",
		},
		{
			name:          "root index link",
			input:         "[Home](index.md)",
			currentDirKey: "",
			expected:      "[Home](/)",
		},
		{
			name:          "multiple links",
			input:         "See [one](one.md) and [two](two.md).",
			currentDirKey: "",
			expected:      "See [one](/one) and [two](/two).",
		},
		{
			name:          "query string preserved",
			input:         "[Link](page.md?query=value)",
			currentDirKey: "",
			expected:      "[Link](page.md?query=value)",
		},
		{
			name:          "mailto link unchanged",
			input:         "[Email](mailto:test@example.com)",
			currentDirKey: "",
			expected:      "[Email](mailto:test@example.com)",
		},
		{
			name:          "root absolute link unchanged",
			input:         "[Page](/page)",
			currentDirKey: "",
			expected:      "[Page](/page)",
		},
		{
			name:          "relative from subdirectory",
			input:         "[Sibling](./sibling.md)",
			currentDirKey: "guide/advanced",
			expected:      "[Sibling](/guide/advanced/sibling)",
		},
		{
			name:          "parent directory traversal",
			input:         "[Parent](../parent.md)",
			currentDirKey: "guide/advanced",
			expected:      "[Parent](/guide/parent)",
		},
		{
			name:          "with basepath prefix",
			input:         "[Getting Started](getting-started)",
			currentDirKey: "",
			expected:      "[Getting Started](/docs/getting-started)",
			basePath:      "/docs",
		},
		{
			name:          "with basepath and nested path",
			input:         "[Guide](guide/intro.md)",
			currentDirKey: "",
			expected:      "[Guide](/docs/guide/intro)",
			basePath:      "/docs",
		},
		{
			name:          "with basepath and root index",
			input:         "[Home](index.md)",
			currentDirKey: "",
			expected:      "[Home](/docs/)",
			basePath:      "/docs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RewriteExtensionlessDocLinks(tt.input, tt.currentDirKey, tt.basePath)
			if result != tt.expected {
				t.Errorf("RewriteExtensionlessDocLinks(%q, %q, %q) = %q; want %q",
					tt.input, tt.currentDirKey, tt.basePath, result, tt.expected)
			}
		})
	}
}
