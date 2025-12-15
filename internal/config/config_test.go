package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Site.Title != "Documentation" {
		t.Errorf("expected default title 'Documentation', got %q", cfg.Site.Title)
	}
	if cfg.Theme.Mode != "auto" {
		t.Errorf("expected default theme mode 'auto', got %q", cfg.Theme.Mode)
	}
	if cfg.Theme.Preset != "default" {
		t.Errorf("expected default preset 'default', got %q", cfg.Theme.Preset)
	}
	if cfg.Nav.ShowSearch == nil || !*cfg.Nav.ShowSearch {
		t.Error("expected ShowSearch to default to true")
	}
	if cfg.Footer.ShowPoweredBy == nil || !*cfg.Footer.ShowPoweredBy {
		t.Error("expected ShowPoweredBy to default to true")
	}
}

func TestLoadYAML(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
site:
  title: "Test Docs"
  description: "A test description"
theme:
  mode: dark
  preset: ocean
nav:
  show_search: false
footer:
  copyright: "© 2024 Test"
  show_powered_by: false
`
	err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Site.Title != "Test Docs" {
		t.Errorf("expected title 'Test Docs', got %q", cfg.Site.Title)
	}
	if cfg.Site.Description != "A test description" {
		t.Errorf("expected description 'A test description', got %q", cfg.Site.Description)
	}
	if cfg.Theme.Mode != "dark" {
		t.Errorf("expected theme mode 'dark', got %q", cfg.Theme.Mode)
	}
	if cfg.Theme.Preset != "ocean" {
		t.Errorf("expected preset 'ocean', got %q", cfg.Theme.Preset)
	}
	if cfg.Nav.ShowSearch == nil || *cfg.Nav.ShowSearch {
		t.Error("expected ShowSearch to be false")
	}
	if cfg.Footer.Copyright != "© 2024 Test" {
		t.Errorf("expected copyright '© 2024 Test', got %q", cfg.Footer.Copyright)
	}
	if cfg.Footer.ShowPoweredBy == nil || *cfg.Footer.ShowPoweredBy {
		t.Error("expected ShowPoweredBy to be false")
	}
}

func TestLoadJSON(t *testing.T) {
	dir := t.TempDir()
	jsonContent := `{
  "site": {
    "title": "JSON Docs",
    "logo": "/logo.png"
  },
  "theme": {
    "mode": "light",
    "preset": "forest"
  }
}`
	err := os.WriteFile(filepath.Join(dir, "dorcs.json"), []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Site.Title != "JSON Docs" {
		t.Errorf("expected title 'JSON Docs', got %q", cfg.Site.Title)
	}
	if cfg.Site.Logo != "/logo.png" {
		t.Errorf("expected logo '/logo.png', got %q", cfg.Site.Logo)
	}
	if cfg.Theme.Mode != "light" {
		t.Errorf("expected theme mode 'light', got %q", cfg.Theme.Mode)
	}
	if cfg.Theme.Preset != "forest" {
		t.Errorf("expected preset 'forest', got %q", cfg.Theme.Preset)
	}
}

func TestLoadNoConfig(t *testing.T) {
	dir := t.TempDir()

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should return defaults
	if cfg.Site.Title != "Documentation" {
		t.Errorf("expected default title 'Documentation', got %q", cfg.Site.Title)
	}
	if cfg.Theme.Mode != "auto" {
		t.Errorf("expected default theme mode 'auto', got %q", cfg.Theme.Mode)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
site:
  title: "From File"
theme:
  preset: lavender
`
	configPath := filepath.Join(dir, "custom-config.yaml")
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if cfg.Site.Title != "From File" {
		t.Errorf("expected title 'From File', got %q", cfg.Site.Title)
	}
	if cfg.Theme.Preset != "lavender" {
		t.Errorf("expected preset 'lavender', got %q", cfg.Theme.Preset)
	}
}

func TestGetColorScheme(t *testing.T) {
	cfg := Default()

	colors := cfg.GetColorScheme()
	if colors.Light.Background != "#ffffff" {
		t.Errorf("expected light background '#ffffff', got %q", colors.Light.Background)
	}
	if colors.Dark.Background != "#0d1117" {
		t.Errorf("expected dark background '#0d1117', got %q", colors.Dark.Background)
	}

	// Test with custom override
	cfg.Theme.Colors.Light.Background = "#fafafa"
	colors = cfg.GetColorScheme()
	if colors.Light.Background != "#fafafa" {
		t.Errorf("expected custom light background '#fafafa', got %q", colors.Light.Background)
	}
}

func TestGetColorSchemePresets(t *testing.T) {
	presets := []string{"default", "ocean", "forest", "sunset", "midnight", "lavender", "rose"}

	for _, preset := range presets {
		cfg := Default()
		cfg.Theme.Preset = preset

		colors := cfg.GetColorScheme()
		if colors.Light.Background == "" {
			t.Errorf("preset %q: expected light background to be set", preset)
		}
		if colors.Dark.Background == "" {
			t.Errorf("preset %q: expected dark background to be set", preset)
		}
		if colors.Light.Accent == "" {
			t.Errorf("preset %q: expected light accent to be set", preset)
		}
		if colors.Dark.Accent == "" {
			t.Errorf("preset %q: expected dark accent to be set", preset)
		}
	}
}

func TestGenerateThemeCSS(t *testing.T) {
	cfg := Default()
	css := cfg.GenerateThemeCSS()

	// Should contain :root declaration
	if !contains(css, ":root {") {
		t.Error("expected CSS to contain :root declaration")
	}

	// Should contain color variables
	if !contains(css, "--bg:") {
		t.Error("expected CSS to contain --bg variable")
	}
	if !contains(css, "--fg:") {
		t.Error("expected CSS to contain --fg variable")
	}
	if !contains(css, "--link:") {
		t.Error("expected CSS to contain --link variable")
	}

	// Auto mode should have prefers-color-scheme media query
	if !contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("expected CSS to contain dark mode media query for auto mode")
	}
}

func TestGenerateThemeCSSLightMode(t *testing.T) {
	cfg := Default()
	cfg.Theme.Mode = "light"
	css := cfg.GenerateThemeCSS()

	// Light mode should NOT have media query
	if contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("light mode should not contain dark mode media query")
	}
}

func TestGenerateThemeCSSDarkMode(t *testing.T) {
	cfg := Default()
	cfg.Theme.Mode = "dark"
	css := cfg.GenerateThemeCSS()

	// Dark mode should use dark colors in :root
	if !contains(css, "#0d1117") {
		t.Error("dark mode should contain dark background color")
	}

	// Should NOT have media query
	if contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("dark mode should not contain dark mode media query")
	}
}

func TestGenerateThemeCSSCustomFonts(t *testing.T) {
	cfg := Default()
	cfg.Theme.FontFamily = `"Inter", sans-serif`
	cfg.Theme.MonoFontFamily = `"JetBrains Mono", monospace`
	css := cfg.GenerateThemeCSS()

	if !contains(css, "--font-sans:") {
		t.Error("expected CSS to contain --font-sans variable")
	}
	if !contains(css, "--font-mono:") {
		t.Error("expected CSS to contain --font-mono variable")
	}
}

func TestNavLinks(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
nav:
  links:
    - title: "GitHub"
      url: "https://github.com/example"
      external: true
      icon: "github"
    - title: "Docs"
      url: "/docs"
      external: false
`
	err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Nav.Links) != 2 {
		t.Fatalf("expected 2 nav links, got %d", len(cfg.Nav.Links))
	}

	if cfg.Nav.Links[0].Title != "GitHub" {
		t.Errorf("expected first link title 'GitHub', got %q", cfg.Nav.Links[0].Title)
	}
	if cfg.Nav.Links[0].Icon != "github" {
		t.Errorf("expected first link icon 'github', got %q", cfg.Nav.Links[0].Icon)
	}
	if !cfg.Nav.Links[0].External {
		t.Error("expected first link to be external")
	}

	if cfg.Nav.Links[1].Title != "Docs" {
		t.Errorf("expected second link title 'Docs', got %q", cfg.Nav.Links[1].Title)
	}
	if cfg.Nav.Links[1].External {
		t.Error("expected second link to NOT be external")
	}
}

// Helper function
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
