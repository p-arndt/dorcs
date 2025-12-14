package syntax

import (
	"strings"
	"testing"
)

func TestGenerateCSS(t *testing.T) {
	tests := []struct {
		name     string
		theme    string
		contains []string // substrings that should be in CSS
	}{
		{
			name:     "github theme",
			theme:    "github",
			contains: []string{".chroma", "color", "@media"},
		},
		{
			name:     "dracula theme",
			theme:    "dracula",
			contains: []string{".chroma"},
		},
		{
			name:     "empty theme defaults",
			theme:    "",
			contains: []string{".chroma"},
		},
		{
			name:     "invalid theme falls back",
			theme:    "nonexistent-theme",
			contains: []string{".chroma"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			css := GenerateCSS(tt.theme)

			if css == "" {
				t.Error("CSS should not be empty")
			}

			for _, substr := range tt.contains {
				if !strings.Contains(css, substr) {
					t.Errorf("CSS should contain %q", substr)
				}
			}

			// Should contain dark mode media query
			if !strings.Contains(css, "@media (prefers-color-scheme: dark)") {
				t.Error("CSS should contain dark mode media query")
			}

			// Should contain chroma class
			if !strings.Contains(css, ".chroma") {
				t.Error("CSS should contain .chroma class")
			}
		})
	}
}

func TestGenerateCSSThemes(t *testing.T) {
	themes := []string{"github", "dracula", "monokai", "nord", "solarized-light"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			css := GenerateCSS(theme)

			if css == "" {
				t.Errorf("CSS for theme %q should not be empty", theme)
			}

			// Should have reasonable length (at least 100 chars for basic CSS)
			if len(css) < 100 {
				t.Errorf("CSS for theme %q seems too short: %d chars", theme, len(css))
			}
		})
	}
}

func TestGenerateCSSDarkMode(t *testing.T) {
	css := GenerateCSS("github")

	// Should contain dark mode styles
	if !strings.Contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("CSS should contain dark mode media query")
	}

	// Should contain both light and dark color overrides
	lightCount := strings.Count(css, "var(--fg, #24292f)")
	darkCount := strings.Count(css, "var(--fg, #e6edf3)")

	if lightCount == 0 {
		t.Error("CSS should contain light mode color overrides")
	}

	if darkCount == 0 {
		t.Error("CSS should contain dark mode color overrides")
	}
}
