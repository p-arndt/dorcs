package markdown

import "testing"

func TestRewriteExtensionlessDocLinksWithLanguage(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		currentDirKey string
		basePath      string
		language      string
		expected      string
	}{
		{
			name:          "default language - no prefix",
			input:         "[Getting Started](getting-started)",
			currentDirKey: "",
			basePath:      "",
			language:      "",
			expected:      "[Getting Started](/getting-started)",
		},
		{
			name:          "German language - add prefix",
			input:         "[Getting Started](getting-started)",
			currentDirKey: "",
			basePath:      "",
			language:      "de",
			expected:      "[Getting Started](/de/getting-started)",
		},
		{
			name:          "German language - root index",
			input:         "[Home](index.md)",
			currentDirKey: "",
			basePath:      "",
			language:      "de",
			expected:      "[Home](/de/)",
		},
		{
			name:          "German language - nested path",
			input:         "[Guide](guide/intro.md)",
			currentDirKey: "",
			basePath:      "",
			language:      "de",
			expected:      "[Guide](/de/guide/intro)",
		},
		{
			name:          "German language - with basepath",
			input:         "[Getting Started](getting-started)",
			currentDirKey: "",
			basePath:      "/docs",
			language:      "de",
			expected:      "[Getting Started](/docs/de/getting-started)",
		},
		{
			name:          "German language - folder index with basepath",
			input:         "[Guide](guide/index.md)",
			currentDirKey: "",
			basePath:      "/docs",
			language:      "de",
			expected:      "[Guide](/docs/de/guide)",
		},
		{
			name:          "default language - with basepath",
			input:         "[Getting Started](getting-started)",
			currentDirKey: "",
			basePath:      "/docs",
			language:      "",
			expected:      "[Getting Started](/docs/getting-started)",
		},
		{
			name:          "absolute URL unchanged with language",
			input:         "[Example](https://example.com)",
			currentDirKey: "",
			basePath:      "",
			language:      "de",
			expected:      "[Example](https://example.com)",
		},
		{
			name:          "anchor link unchanged with language",
			input:         "[Section](#section)",
			currentDirKey: "",
			basePath:      "",
			language:      "de",
			expected:      "[Section](#section)",
		},
		{
			name:          "root absolute link unchanged with language",
			input:         "[Page](/page)",
			currentDirKey: "",
			basePath:      "",
			language:      "de",
			expected:      "[Page](/page)",
		},
		{
			name:          "French language - nested deep path",
			input:         "[API](api/v2/endpoints.md)",
			currentDirKey: "",
			basePath:      "",
			language:      "fr",
			expected:      "[API](/fr/api/v2/endpoints)",
		},
		{
			name:          "German language - relative from subdirectory",
			input:         "[Sibling](./sibling.md)",
			currentDirKey: "guide/advanced",
			basePath:      "",
			language:      "de",
			expected:      "[Sibling](/de/guide/advanced/sibling)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RewriteExtensionlessDocLinks(tt.input, tt.currentDirKey, tt.basePath, tt.language)
			if result != tt.expected {
				t.Errorf("RewriteExtensionlessDocLinks(%q, %q, %q, %q) = %q; want %q",
					tt.input, tt.currentDirKey, tt.basePath, tt.language, result, tt.expected)
			}
		})
	}
}
