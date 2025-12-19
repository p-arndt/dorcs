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
		{
			name:          "link with anchor",
			input:         "[Section](page#section)",
			currentDirKey: "",
			expected:      "[Section](/page#section)",
		},
		{
			name:          "link with .md extension and anchor",
			input:         "[Section](page.md#section)",
			currentDirKey: "",
			expected:      "[Section](/page#section)",
		},
		{
			name:          "root absolute link with .md extension and anchor",
			input:         "[Section](/page.md#section)",
			currentDirKey: "",
			expected:      "[Section](/page#section)",
		},
		{
			name:          "root absolute link with anchor",
			input:         "[Section](/page#section)",
			currentDirKey: "",
			expected:      "[Section](/page#section)",
		},
		{
			name:          "nested link with .md extension and anchor",
			input:         "[Section](guide/page.md#section)",
			currentDirKey: "",
			expected:      "[Section](/guide/page#section)",
		},
		{
			name:          "relative link with .md extension and anchor",
			input:         "[Section](./page.md#section)",
			currentDirKey: "guide",
			expected:      "[Section](/guide/page#section)",
		},
		{
			name:          "parent directory link with .md extension and anchor",
			input:         "[Section](../page.md#section)",
			currentDirKey: "guide/advanced",
			expected:      "[Section](/guide/page#section)",
		},
		{
			name:          "index link with .md extension and anchor",
			input:         "[Home](index.md#intro)",
			currentDirKey: "",
			expected:      "[Home](/#intro)",
		},
		{
			name:          "folder index link with .md extension and anchor",
			input:         "[Guide](guide/index.md#intro)",
			currentDirKey: "",
			expected:      "[Guide](/guide#intro)",
		},
		{
			name:          "with basepath and anchor",
			input:         "[Section](page.md#section)",
			currentDirKey: "",
			expected:      "[Section](/docs/page#section)",
			basePath:      "/docs",
		},
		{
			name:          "root absolute with basepath and anchor",
			input:         "[Section](/page.md#section)",
			currentDirKey: "",
			expected:      "[Section](/page#section)",
			basePath:      "/docs",
		},
		{
			name:          "multiple anchors in same document",
			input:         "See [Section 1](page.md#section1) and [Section 2](page.md#section2).",
			currentDirKey: "",
			expected:      "See [Section 1](/page#section1) and [Section 2](/page#section2).",
		},
		{
			name:          "anchor with complex name",
			input:         "[Init Command](./07_commands.md#init-command)",
			currentDirKey: "",
			expected:      "[Init Command](/07_commands#init-command)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RewriteExtensionlessDocLinks(tt.input, tt.currentDirKey, tt.basePath, "", "", "", "")
			if result != tt.expected {
				t.Errorf("RewriteExtensionlessDocLinks(%q, %q, %q) = %q; want %q",
					tt.input, tt.currentDirKey, tt.basePath, result, tt.expected)
			}
		})
	}
}

func TestResolveLinkToDocKey(t *testing.T) {
	tests := []struct {
		name          string
		href          string
		currentDirKey string
		expectedKey   string
		shouldCheck   bool
	}{
		{
			name:          "simple link",
			href:          "getting-started",
			currentDirKey: "",
			expectedKey:   "getting-started",
			shouldCheck:   true,
		},
		{
			name:          "link with .md extension",
			href:          "guide.md",
			currentDirKey: "",
			expectedKey:   "guide",
			shouldCheck:   true,
		},
		{
			name:          "link with anchor",
			href:          "page#section",
			currentDirKey: "",
			expectedKey:   "page",
			shouldCheck:   true,
		},
		{
			name:          "link with .md extension and anchor",
			href:          "page.md#section",
			currentDirKey: "",
			expectedKey:   "page",
			shouldCheck:   true,
		},
		{
			name:          "root absolute link with .md extension and anchor",
			href:          "/page.md#section",
			currentDirKey: "",
			expectedKey:   "page",
			shouldCheck:   true,
		},
		{
			name:          "nested link with .md extension and anchor",
			href:          "guide/page.md#section",
			currentDirKey: "",
			expectedKey:   "guide/page",
			shouldCheck:   true,
		},
		{
			name:          "relative link with .md extension and anchor",
			href:          "./page.md#section",
			currentDirKey: "guide",
			expectedKey:   "guide/page",
			shouldCheck:   true,
		},
		{
			name:          "parent directory link with .md extension and anchor",
			href:          "../page.md#section",
			currentDirKey: "guide/advanced",
			expectedKey:   "guide/page",
			shouldCheck:   true,
		},
		{
			name:          "index link with .md extension and anchor",
			href:          "index.md#intro",
			currentDirKey: "",
			expectedKey:   "",
			shouldCheck:   true,
		},
		{
			name:          "folder index link with .md extension and anchor",
			href:          "guide/index.md#intro",
			currentDirKey: "",
			expectedKey:   "guide",
			shouldCheck:   true,
		},
		{
			name:          "anchor with complex name",
			href:          "./07_commands.md#init-command",
			currentDirKey: "",
			expectedKey:   "07_commands",
			shouldCheck:   true,
		},
		{
			name:          "standalone anchor",
			href:          "#section",
			currentDirKey: "",
			expectedKey:   "",
			shouldCheck:   false,
		},
		{
			name:          "absolute URL",
			href:          "https://example.com/page",
			currentDirKey: "",
			expectedKey:   "",
			shouldCheck:   false,
		},
		{
			name:          "image link",
			href:          "image.png",
			currentDirKey: "",
			expectedKey:   "",
			shouldCheck:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, shouldCheck := ResolveLinkToDocKey(tt.href, tt.currentDirKey)
			if shouldCheck != tt.shouldCheck {
				t.Errorf("ResolveLinkToDocKey(%q, %q) shouldCheck = %v; want %v",
					tt.href, tt.currentDirKey, shouldCheck, tt.shouldCheck)
			}
			if shouldCheck && key != tt.expectedKey {
				t.Errorf("ResolveLinkToDocKey(%q, %q) key = %q; want %q",
					tt.href, tt.currentDirKey, key, tt.expectedKey)
			}
		})
	}
}
