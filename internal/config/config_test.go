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

func TestNavItemsLoad(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
nav:
  items:
    - Home: index.md
    - Getting Started: 01_getting-started.md
    - Usage:
        page: usage/index.md
        items:
          - Writing: usage/writing-your-docs.md
          - Metadata: usage/metadata.md
    - External:
        items:
          - GitHub: external-content/github.md
`
	err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Nav.Items) != 4 {
		t.Fatalf("expected 4 nav items, got %d", len(cfg.Nav.Items))
	}
	if cfg.Nav.Items[0].Label != "Home" || cfg.Nav.Items[0].Page != "index.md" {
		t.Fatalf("unexpected first nav item: %+v", cfg.Nav.Items[0])
	}
	if cfg.Nav.Items[2].Label != "Usage" || cfg.Nav.Items[2].Page != "usage/index.md" {
		t.Fatalf("unexpected Usage item: %+v", cfg.Nav.Items[2])
	}
	if len(cfg.Nav.Items[2].Items) != 2 {
		t.Fatalf("expected Usage to have 2 children, got %d", len(cfg.Nav.Items[2].Items))
	}
	if cfg.Nav.Items[3].Page != "" || len(cfg.Nav.Items[3].Items) != 1 {
		t.Fatalf("unexpected External item: %+v", cfg.Nav.Items[3])
	}
}

func TestNavItemsValidation(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
nav:
  items:
    - Empty: {}
`
	err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err = Load(dir)
	if err == nil {
		t.Fatal("expected invalid nav.items config to fail")
	}
	if !contains(err.Error(), "must define page, items, or both") {
		t.Fatalf("unexpected error: %v", err)
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

func TestNavSectionsLoad(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
nav:
  sections:
    - title: "Getting Started"
      items:
        - Overview: index.md
        - Quickstart: getting-started.md
    - title: "Reference"
      items:
        - API: api.md
        - CLI:
            page: cli/index.md
            items:
              - Commands: cli/commands.md
`
	err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Nav.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(cfg.Nav.Sections))
	}
	if cfg.Nav.Sections[0].Title != "Getting Started" {
		t.Fatalf("unexpected first section title: %q", cfg.Nav.Sections[0].Title)
	}
	if len(cfg.Nav.Sections[0].Items) != 2 {
		t.Fatalf("expected first section to have 2 items, got %d", len(cfg.Nav.Sections[0].Items))
	}
	if cfg.Nav.Sections[1].Title != "Reference" {
		t.Fatalf("unexpected second section title: %q", cfg.Nav.Sections[1].Title)
	}
	if len(cfg.Nav.Sections[1].Items) != 2 {
		t.Fatalf("expected second section to have 2 items, got %d", len(cfg.Nav.Sections[1].Items))
	}
	// Check nested items in second section
	refCLI := cfg.Nav.Sections[1].Items[1]
	if refCLI.Label != "CLI" || refCLI.Page != "cli/index.md" || len(refCLI.Items) != 1 {
		t.Fatalf("unexpected CLI item: %+v", refCLI)
	}
}

func TestNavSectionsValidation(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr string
	}{
		{
			name: "empty section title",
			yaml: `
nav:
  sections:
    - title: ""
      items:
        - Page: index.md
`,
			wantErr: "title cannot be empty",
		},
		{
			name: "section with no items",
			yaml: `
nav:
  sections:
    - title: "Empty"
      items: []
`,
			wantErr: "must have at least one item",
		},
		{
			name: "section with invalid item",
			yaml: `
nav:
  sections:
    - title: "Bad"
      items:
        - Empty: {}
`,
			wantErr: "must define page, items, or both",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(tt.yaml), 0644)
			if err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}
			_, err = Load(dir)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestAnnouncementConfig(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
announcement:
  text: "New release v2.0!"
  dismissible: true
`
	err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Announcement == nil {
		t.Fatal("expected announcement to be set")
	}
	if cfg.Announcement.Text != "New release v2.0!" {
		t.Fatalf("unexpected announcement text: %q", cfg.Announcement.Text)
	}
	if cfg.Announcement.Dismissible == nil || !*cfg.Announcement.Dismissible {
		t.Fatal("expected dismissible to be true")
	}
}

func TestLoadEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create .env file
	envContent := `GITHUB_TOKEN=test-token-from-env
OTHER_VAR=test-value
# This is a comment
EMPTY_VAR=
QUOTED_VAR="quoted value"
`
	if err := os.WriteFile(".env", []byte(envContent), 0644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("OTHER_VAR")

	// Load config (which should load .env file)
	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify that environment variable was loaded
	if os.Getenv("GITHUB_TOKEN") != "test-token-from-env" {
		t.Errorf("expected GITHUB_TOKEN to be loaded from .env, got %q", os.Getenv("GITHUB_TOKEN"))
	}

	// Test that config can use the env var
	yamlContent := `
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: "${GITHUB_TOKEN}"
`
	if err := os.WriteFile("dorcs.yaml", []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err = Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.GitHub.Token != "test-token-from-env" {
		t.Errorf("expected token to be expanded from .env, got %q", cfg.GitHub.Token)
	}
}

func TestExplicitConfigurationRequired(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}

	// Create language and version folders
	os.MkdirAll(filepath.Join(docsDir, "en"), 0755)
	os.MkdirAll(filepath.Join(docsDir, "de"), 0755)
	os.MkdirAll(filepath.Join(docsDir, "v1"), 0755)
	os.MkdirAll(filepath.Join(docsDir, "v2"), 0755)

	// Create a markdown file in each folder
	os.WriteFile(filepath.Join(docsDir, "en", "index.md"), []byte("# English"), 0644)
	os.WriteFile(filepath.Join(docsDir, "de", "index.md"), []byte("# Deutsch"), 0644)
	os.WriteFile(filepath.Join(docsDir, "v1", "index.md"), []byte("# V1"), 0644)
	os.WriteFile(filepath.Join(docsDir, "v2", "index.md"), []byte("# V2"), 0644)

	// Load config - no auto-detection, explicit configuration required
	cfg, err := Load(docsDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should NOT auto-detect (auto-detection has been removed)
	if len(cfg.Languages.Enabled) > 0 {
		t.Errorf("expected no languages to be auto-detected, got %d", len(cfg.Languages.Enabled))
	}
	if len(cfg.Versions.Enabled) > 0 {
		t.Errorf("expected no versions to be auto-detected, got %d", len(cfg.Versions.Enabled))
	}
}

func TestLoadAddsDefaultVersionToEnabledList(t *testing.T) {
	tmpDir := t.TempDir()

	yamlContent := `
versions:
  default: "latest"
  enabled:
    - id: "v1"
      name: "V1"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "dorcs.yaml"), []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if got := cfg.GetDefaultVersion(); got != "latest" {
		t.Fatalf("default version = %q, want %q", got, "latest")
	}
	if len(cfg.Versions.Enabled) != 2 {
		t.Fatalf("expected 2 enabled versions, got %d", len(cfg.Versions.Enabled))
	}
	if cfg.Versions.Enabled[0].ID != "latest" {
		t.Fatalf("first enabled version = %q, want %q", cfg.Versions.Enabled[0].ID, "latest")
	}
	if cfg.Versions.Enabled[1].ID != "v1" {
		t.Fatalf("second enabled version = %q, want %q", cfg.Versions.Enabled[1].ID, "v1")
	}
}
