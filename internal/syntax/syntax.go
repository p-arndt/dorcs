package syntax

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

// GenerateCSS creates CSS for both light and dark modes.
func GenerateCSS(codeTheme string) string {
	if codeTheme == "" {
		codeTheme = "github" // fallback
	}

	chromaStyle := styles.Get(codeTheme)
	if chromaStyle == nil {
		chromaStyle = styles.Fallback
	}

	return generateSyntaxCSS(chromaStyle, codeTheme)
}

// getDarkThemeVariant returns a dark variant theme name for a given light theme.
func getDarkThemeVariant(lightThemeName string) string {
	// Map light themes to their dark variants
	variantMap := map[string]string{
		"github":          "github-dark",
		"solarized-light": "solarized-dark",
		"dracula":         "dracula",
		"nord":            "nord",
		"monokai":         "monokai",
		"onedark":         "onedark",
	}

	if variant, ok := variantMap[lightThemeName]; ok {
		return variant
	}
	// If no specific variant, try appending "-dark"
	if dark := styles.Get(lightThemeName + "-dark"); dark != nil {
		return lightThemeName + "-dark"
	}
	// Fallback
	return "github-dark"
}

// generateSyntaxCSS creates CSS for both light and dark modes.
func generateSyntaxCSS(lightStyle *chroma.Style, lightThemeName string) string {
	var buf bytes.Buffer

	// Light mode CSS
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	_ = formatter.WriteCSS(&buf, lightStyle)

	// Add default color override to ensure visibility
	// This fixes issues where Chroma themes don't set proper default colors
	// Use CSS variables so it adapts to the theme
	// Note: Specific token colors from Chroma will override these defaults
	buf.WriteString("\n/* Ensure default text is visible - uses theme foreground color */\n")
	buf.WriteString("/* Chroma token-specific colors will override these defaults */\n")
	buf.WriteString(".chroma { color: var(--fg, #24292f); }\n")
	buf.WriteString(".chroma .cl { color: var(--fg, #24292f); }\n")
	buf.WriteString(".chroma .line { color: var(--fg, #24292f); }\n")
	// Fix github theme: nx (identifiers) and p (punctuation) often have poor contrast
	buf.WriteString(".chroma .nx { color: var(--fg, #24292f); }\n")
	buf.WriteString(".chroma .p { color: var(--fg, #24292f); }\n")

	// Dark mode CSS - try to find a dark variant of the light theme
	darkStyleName := getDarkThemeVariant(lightThemeName)
	darkStyle := styles.Get(darkStyleName)
	if darkStyle == nil {
		// Fallback to github-dark or monokai
		darkStyle = styles.Get("github-dark")
		if darkStyle == nil {
			darkStyle = styles.Get("monokai")
		}
		if darkStyle == nil {
			darkStyle = styles.Fallback
		}
	}

	buf.WriteString("\n@media (prefers-color-scheme: dark) {\n")
	var darkBuf bytes.Buffer
	_ = formatter.WriteCSS(&darkBuf, darkStyle)
	// Indent the dark mode CSS
	for _, line := range strings.Split(darkBuf.String(), "\n") {
		if strings.TrimSpace(line) != "" {
			buf.WriteString("  " + line + "\n")
		}
	}
	// Add default color override for dark mode - uses theme foreground color
	// Note: Specific token colors from Chroma will override these defaults
	buf.WriteString("  .chroma { color: var(--fg, #e6edf3); }\n")
	buf.WriteString("  .chroma .cl { color: var(--fg, #e6edf3); }\n")
	buf.WriteString("  .chroma .line { color: var(--fg, #e6edf3); }\n")
	// Fix github-dark theme: nx (identifiers) and p (punctuation) often have poor contrast
	buf.WriteString("  .chroma .nx { color: var(--fg, #e6edf3); }\n")
	buf.WriteString("  .chroma .p { color: var(--fg, #e6edf3); }\n")
	buf.WriteString("  .chroma .ge { color: var(--fg, #e6edf3); }\n")
	buf.WriteString("  .chroma .na { color: var(--fg, #e6edf3); }\n")
	buf.WriteString("}\n")

	return buf.String()
}
