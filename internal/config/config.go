// Package config provides configuration handling for dorcs documentation server.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration options for the documentation site.
type Config struct {
	// Server port (default: 8080)
	Port int `json:"port" yaml:"port"`

	// Site metadata
	Site SiteConfig `json:"site" yaml:"site"`

	// Theme and styling
	Theme ThemeConfig `json:"theme" yaml:"theme"`

	// Navigation settings
	Nav NavConfig `json:"nav" yaml:"nav"`

	// Footer settings
	Footer FooterConfig `json:"footer" yaml:"footer"`

	// Authentication settings
	Auth AuthConfig `json:"auth" yaml:"auth"`

	// Languages configuration for multi-lingual support
	Languages LanguagesConfig `json:"languages" yaml:"languages"`

	// GitHub integration configuration
	GitHub GitHubConfig `json:"github" yaml:"github"`
}

// SiteConfig holds site-level metadata.
type SiteConfig struct {
	// Title shown in the header/brand area (top-left)
	Title string `json:"title" yaml:"title"`

	// Description for meta tags
	Description string `json:"description" yaml:"description"`

	// Logo URL (optional, replaces text title if set)
	Logo string `json:"logo" yaml:"logo"`

	// Favicon URL (optional, overrides default)
	Favicon string `json:"favicon" yaml:"favicon"`
}

// ThemeConfig holds theming and styling options.
type ThemeConfig struct {
	// Mode: "light", "dark", or "auto" (follows system preference)
	Mode string `json:"mode" yaml:"mode"`

	// Preset theme name: "default", "ocean", "forest", "sunset", "midnight", "lavender"
	Preset string `json:"preset" yaml:"preset"`

	// Custom colors (override preset if set)
	Colors ColorConfig `json:"colors" yaml:"colors"`

	// CodeTheme for syntax highlighting (deprecated: now determined by preset)
	// Kept for backward compatibility but ignored - preset determines code theme
	CodeTheme string

	// CustomCSS is a path to a custom CSS file to include
	CustomCSS string `json:"custom_css" yaml:"custom_css"`

	// FontFamily for body text (CSS font-family value)
	FontFamily string `json:"font_family" yaml:"font_family"`

	// MonoFontFamily for code (CSS font-family value)
	MonoFontFamily string `json:"mono_font_family" yaml:"mono_font_family"`
}

// ColorConfig allows custom color overrides.
type ColorConfig struct {
	// Light mode colors
	Light ColorScheme `json:"light" yaml:"light"`

	// Dark mode colors
	Dark ColorScheme `json:"dark" yaml:"dark"`
}

// ColorScheme defines a set of theme colors.
type ColorScheme struct {
	// Background color
	Background string `json:"background" yaml:"background"`

	// Foreground/text color
	Foreground string `json:"foreground" yaml:"foreground"`

	// Muted text color
	Muted string `json:"muted" yaml:"muted"`

	// Border color
	Border string `json:"border" yaml:"border"`

	// Link/accent color
	Accent string `json:"accent" yaml:"accent"`

	// Code background color
	CodeBackground string `json:"code_background" yaml:"code_background"`

	// Header background (optional, defaults to Background)
	HeaderBackground string `json:"header_background" yaml:"header_background"`

	// Sidebar background (optional, defaults to Background)
	SidebarBackground string `json:"sidebar_background" yaml:"sidebar_background"`
}

// NavConfig holds navigation configuration.
type NavConfig struct {
	// ShowSearch enables/disables the search box (default: true)
	ShowSearch *bool `json:"show_search" yaml:"show_search"`

	// ExpandAll keeps all folders expanded by default (default: false)
	ExpandAll bool `json:"expand_all" yaml:"expand_all"`

	// Links are additional navigation links shown in the header
	Links []NavLink `json:"links" yaml:"links"`
}

// NavLink represents a custom navigation link.
type NavLink struct {
	// Title of the link
	Title string `json:"title" yaml:"title"`

	// URL to link to
	URL string `json:"url" yaml:"url"`

	// External opens link in new tab if true
	External bool `json:"external" yaml:"external"`

	// Icon name (optional): "github", "twitter", "discord", "external"
	Icon string `json:"icon" yaml:"icon"`
}

// FooterConfig holds footer configuration.
type FooterConfig struct {
	// Text shown in the footer (supports basic markdown)
	Text string `json:"text" yaml:"text"`

	// ShowPoweredBy shows "Powered by dorcs" (default: true)
	ShowPoweredBy *bool `json:"show_powered_by" yaml:"show_powered_by"`

	// Copyright text (e.g., "© 2024 My Company")
	Copyright string `json:"copyright" yaml:"copyright"`

	// Links shown in footer
	Links []NavLink `json:"links" yaml:"links"`
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	// Enable authentication (default: false)
	Enabled bool `json:"enabled" yaml:"enabled"`

	// Username for login
	Username string `json:"username" yaml:"username"`

	// Password (will be hashed on first use)
	Password string `json:"password" yaml:"password"`

	// Password hash (set automatically, don't set manually)
	PasswordHash string `json:"password_hash,omitempty" yaml:"password_hash,omitempty"`

	// Sessions file path (default: .dorcs_sessions.json in docs dir)
	SessionsPath string `json:"sessions_path,omitempty" yaml:"sessions_path,omitempty"`
}

// LanguagesConfig holds multi-lingual configuration.
type LanguagesConfig struct {
	// Default language code (e.g., "en")
	Default string `json:"default" yaml:"default"`

	// Enabled languages list
	Enabled []Language `json:"enabled" yaml:"enabled"`
}

// Language represents a single language configuration.
type Language struct {
	// Code is the language code (e.g., "en", "de", "fr")
	Code string `json:"code" yaml:"code"`

	// Name is the display name (e.g., "English", "Deutsch", "Français")
	Name string `json:"name" yaml:"name"`
}

// GitHubConfig holds GitHub integration configuration.
type GitHubConfig struct {
	// Enabled enables/disables GitHub integration
	Enabled bool `json:"enabled" yaml:"enabled"`

	// Repository is the GitHub repository tree URL (e.g., "https://github.com/p-arndt/dorcs/tree/main/docs")
	Repository string `json:"repository" yaml:"repository"`

	// Token is the GitHub token for private repos (supports ${ENV_VAR} expansion)
	Token string `json:"token" yaml:"token"`

	// CacheTTL is the cache time-to-live (default: 1h)
	CacheTTL string `json:"cache_ttl" yaml:"cache_ttl"`
}

// Default returns the default configuration.
func Default() *Config {
	showSearch := true
	showPoweredBy := true

	return &Config{
		Site: SiteConfig{
			Title: "Documentation",
		},
		Theme: ThemeConfig{
			Mode:      "auto",
			Preset:    "default",
			CodeTheme: "github",
		},
		Nav: NavConfig{
			ShowSearch: &showSearch,
			ExpandAll:  false,
		},
		Footer: FooterConfig{
			ShowPoweredBy: &showPoweredBy,
		},
	}
}

// Load reads configuration from a file.
// It first looks in the current working directory, then in the given docs directory.
// It looks for dorcs.yaml, dorcs.yml, or dorcs.json in order.
// Returns default config if no config file is found.
func Load(docsDir string) (*Config, error) {
	// Load .env file first if it exists
	loadEnvFile()

	cfg := Default()

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		// If we can't get working directory, fall back to docs directory only
		wd = ""
	}

	// Try to find config file in working directory first, then in docs directory
	searchDirs := []string{}
	if wd != "" {
		searchDirs = append(searchDirs, wd)
	}
	// Only add docs directory if it's different from working directory
	if docsDir != "" && docsDir != wd {
		absDocsDir, err := filepath.Abs(docsDir)
		if err == nil && absDocsDir != wd {
			searchDirs = append(searchDirs, absDocsDir)
		}
	}

	// Try YAML first, then JSON
	for _, searchDir := range searchDirs {
		for _, name := range []string{"dorcs.yaml", "dorcs.yml"} {
			path := filepath.Join(searchDir, name)
			if data, err := os.ReadFile(path); err == nil {
				if err := yaml.Unmarshal(data, cfg); err != nil {
					return nil, err
				}
				applyDefaults(cfg)
				expandEnvVars(cfg)
				return cfg, nil
			}
		}

		// Try JSON
		path := filepath.Join(searchDir, "dorcs.json")
		if data, err := os.ReadFile(path); err == nil {
			if err := json.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
			applyDefaults(cfg)
			expandEnvVars(cfg)
			return cfg, nil
		}
	}

	return cfg, nil
}

// LoadFromFile reads configuration from a specific file path.
func LoadFromFile(path string) (*Config, error) {
	// Load .env file first if it exists
	loadEnvFile()

	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	default:
		// Try YAML first, then JSON
		if err := yaml.Unmarshal(data, cfg); err != nil {
			if err := json.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
		}
	}

	applyDefaults(cfg)
	expandEnvVars(cfg)
	return cfg, nil
}

// loadEnvFile loads environment variables from a .env file if it exists.
// Uses godotenv library for proper .env file parsing.
func loadEnvFile() {
	// Try to load .env file (ignore errors if file doesn't exist)
	_ = godotenv.Load()
}

// expandEnvVars expands environment variables in config strings.
// Supports ${VAR} and ${VAR:-default} syntax.
func expandEnvVars(cfg *Config) {
	cfg.GitHub.Token = expandEnvVar(cfg.GitHub.Token)
}

var envVarRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

// expandEnvVar expands environment variables in a string.
func expandEnvVar(s string) string {
	return envVarRegex.ReplaceAllStringFunc(s, func(match string) string {
		// Remove ${ and }
		varName := match[2 : len(match)-1]

		// Check for default value syntax: ${VAR:-default}
		if idx := strings.Index(varName, ":-"); idx != -1 {
			envVar := varName[:idx]
			defaultVal := varName[idx+2:]
			if val := os.Getenv(envVar); val != "" {
				return val
			}
			return defaultVal
		}

		// Simple ${VAR} syntax
		if val := os.Getenv(varName); val != "" {
			return val
		}
		return match // Return original if not found
	})
}

// FindConfigFile looks for a config file in the current working directory first,
// then in the given docs directory. Returns the path if found, empty string otherwise.
func FindConfigFile(docsDir string) string {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	// Try to find config file in working directory first, then in docs directory
	searchDirs := []string{}
	if wd != "" {
		searchDirs = append(searchDirs, wd)
	}
	// Only add docs directory if it's different from working directory
	if docsDir != "" && docsDir != wd {
		absDocsDir, err := filepath.Abs(docsDir)
		if err == nil && absDocsDir != wd {
			searchDirs = append(searchDirs, absDocsDir)
		}
	}

	// Try YAML first, then JSON
	for _, searchDir := range searchDirs {
		for _, name := range []string{"dorcs.yaml", "dorcs.yml", "dorcs.json"} {
			path := filepath.Join(searchDir, name)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// applyDefaults ensures required fields have sensible defaults.
func applyDefaults(cfg *Config) {
	if cfg.Site.Title == "" {
		cfg.Site.Title = "Documentation"
	}
	if cfg.Theme.Mode == "" {
		cfg.Theme.Mode = "auto"
	}
	if cfg.Theme.Preset == "" {
		cfg.Theme.Preset = "default"
	}
	// CodeTheme is now determined by preset automatically
	// We still set it for backward compatibility, but it's ignored
	cfg.Theme.CodeTheme = getCodeThemeForPreset(cfg.Theme.Preset)
	if cfg.Nav.ShowSearch == nil {
		showSearch := true
		cfg.Nav.ShowSearch = &showSearch
	}
	if cfg.Footer.ShowPoweredBy == nil {
		showPoweredBy := true
		cfg.Footer.ShowPoweredBy = &showPoweredBy
	}
	// Auth defaults
	if cfg.Auth.SessionsPath == "" {
		cfg.Auth.SessionsPath = ".dorcs_sessions.json"
	}
	// Languages defaults - if default is set but no enabled list, add default to enabled
	if cfg.Languages.Default != "" && len(cfg.Languages.Enabled) == 0 {
		cfg.Languages.Enabled = []Language{{Code: cfg.Languages.Default, Name: cfg.Languages.Default}}
	}
}

// GetColorScheme returns the effective color scheme based on config.
func (c *Config) GetColorScheme() ColorConfig {
	// Start with preset
	colors, ok := ThemePresets[c.Theme.Preset]
	if !ok {
		colors = ThemePresets["default"]
	}

	// Override with custom colors if set
	if c.Theme.Colors.Light.Background != "" {
		colors.Light.Background = c.Theme.Colors.Light.Background
	}
	if c.Theme.Colors.Light.Foreground != "" {
		colors.Light.Foreground = c.Theme.Colors.Light.Foreground
	}
	if c.Theme.Colors.Light.Muted != "" {
		colors.Light.Muted = c.Theme.Colors.Light.Muted
	}
	if c.Theme.Colors.Light.Border != "" {
		colors.Light.Border = c.Theme.Colors.Light.Border
	}
	if c.Theme.Colors.Light.Accent != "" {
		colors.Light.Accent = c.Theme.Colors.Light.Accent
	}
	if c.Theme.Colors.Light.CodeBackground != "" {
		colors.Light.CodeBackground = c.Theme.Colors.Light.CodeBackground
	}
	if c.Theme.Colors.Light.HeaderBackground != "" {
		colors.Light.HeaderBackground = c.Theme.Colors.Light.HeaderBackground
	}
	if c.Theme.Colors.Light.SidebarBackground != "" {
		colors.Light.SidebarBackground = c.Theme.Colors.Light.SidebarBackground
	}

	if c.Theme.Colors.Dark.Background != "" {
		colors.Dark.Background = c.Theme.Colors.Dark.Background
	}
	if c.Theme.Colors.Dark.Foreground != "" {
		colors.Dark.Foreground = c.Theme.Colors.Dark.Foreground
	}
	if c.Theme.Colors.Dark.Muted != "" {
		colors.Dark.Muted = c.Theme.Colors.Dark.Muted
	}
	if c.Theme.Colors.Dark.Border != "" {
		colors.Dark.Border = c.Theme.Colors.Dark.Border
	}
	if c.Theme.Colors.Dark.Accent != "" {
		colors.Dark.Accent = c.Theme.Colors.Dark.Accent
	}
	if c.Theme.Colors.Dark.CodeBackground != "" {
		colors.Dark.CodeBackground = c.Theme.Colors.Dark.CodeBackground
	}
	if c.Theme.Colors.Dark.HeaderBackground != "" {
		colors.Dark.HeaderBackground = c.Theme.Colors.Dark.HeaderBackground
	}
	if c.Theme.Colors.Dark.SidebarBackground != "" {
		colors.Dark.SidebarBackground = c.Theme.Colors.Dark.SidebarBackground
	}

	return colors
}

// GenerateThemeCSS generates CSS custom properties based on config.
func (c *Config) GenerateThemeCSS() string {
	colors := c.GetColorScheme()
	css := ""

	// Helper to add color if set
	addColor := func(name, value string) string {
		if value != "" {
			return "    " + name + ": " + value + ";\n"
		}
		return ""
	}

	// Define light mode variables (default)
	css += ":root {\n"
	css += addColor("--bg", colors.Light.Background)
	css += addColor("--fg", colors.Light.Foreground)
	css += addColor("--muted", colors.Light.Muted)
	css += addColor("--border", colors.Light.Border)
	css += addColor("--link", colors.Light.Accent)
	css += addColor("--code-bg", colors.Light.CodeBackground)
	if colors.Light.HeaderBackground != "" {
		css += addColor("--header-bg", colors.Light.HeaderBackground)
	}
	if colors.Light.SidebarBackground != "" {
		css += addColor("--sidebar-bg", colors.Light.SidebarBackground)
	}
	if c.Theme.FontFamily != "" {
		css += "    --font-sans: " + c.Theme.FontFamily + ";\n"
	}
	if c.Theme.MonoFontFamily != "" {
		css += "    --font-mono: " + c.Theme.MonoFontFamily + ";\n"
	}
	// Border radius is fixed at 12px (convention over configuration)
	css += "    --radius: 12px;\n"
	css += "    --radius-sm: 10px;\n"
	css += "}\n\n"

	// Define dark mode variables as CSS custom properties for manual toggle
	css += ":root {\n"
	css += "  --bg-dark: " + colors.Dark.Background + ";\n"
	css += "  --fg-dark: " + colors.Dark.Foreground + ";\n"
	css += "  --muted-dark: " + colors.Dark.Muted + ";\n"
	css += "  --border-dark: " + colors.Dark.Border + ";\n"
	css += "  --link-dark: " + colors.Dark.Accent + ";\n"
	css += "  --code-bg-dark: " + colors.Dark.CodeBackground + ";\n"
	if colors.Dark.HeaderBackground != "" {
		css += "  --header-bg-dark: " + colors.Dark.HeaderBackground + ";\n"
	}
	if colors.Dark.SidebarBackground != "" {
		css += "  --sidebar-bg-dark: " + colors.Dark.SidebarBackground + ";\n"
	}
	css += "}\n\n"

	// Apply dark mode based on theme mode setting
	switch c.Theme.Mode {
	case "dark":
		// Force dark mode by default
		css += ":root {\n"
		css += addColor("--bg", colors.Dark.Background)
		css += addColor("--fg", colors.Dark.Foreground)
		css += addColor("--muted", colors.Dark.Muted)
		css += addColor("--border", colors.Dark.Border)
		css += addColor("--link", colors.Dark.Accent)
		css += addColor("--code-bg", colors.Dark.CodeBackground)
		if colors.Dark.HeaderBackground != "" {
			css += addColor("--header-bg", colors.Dark.HeaderBackground)
		}
		if colors.Dark.SidebarBackground != "" {
			css += addColor("--sidebar-bg", colors.Dark.SidebarBackground)
		}
		css += "}\n\n"
		// Allow manual override to light mode
		css += ":root.light-mode {\n"
		css += addColor("--bg", colors.Light.Background)
		css += addColor("--fg", colors.Light.Foreground)
		css += addColor("--muted", colors.Light.Muted)
		css += addColor("--border", colors.Light.Border)
		css += addColor("--link", colors.Light.Accent)
		css += addColor("--code-bg", colors.Light.CodeBackground)
		if colors.Light.HeaderBackground != "" {
			css += addColor("--header-bg", colors.Light.HeaderBackground)
		}
		if colors.Light.SidebarBackground != "" {
			css += addColor("--sidebar-bg", colors.Light.SidebarBackground)
		}
		css += "}\n"
	case "light":
		// Already set as default, allow manual override to dark mode
		css += ":root.dark-mode {\n"
		css += addColor("--bg", colors.Dark.Background)
		css += addColor("--fg", colors.Dark.Foreground)
		css += addColor("--muted", colors.Dark.Muted)
		css += addColor("--border", colors.Dark.Border)
		css += addColor("--link", colors.Dark.Accent)
		css += addColor("--code-bg", colors.Dark.CodeBackground)
		if colors.Dark.HeaderBackground != "" {
			css += addColor("--header-bg", colors.Dark.HeaderBackground)
		}
		if colors.Dark.SidebarBackground != "" {
			css += addColor("--sidebar-bg", colors.Dark.SidebarBackground)
		}
		css += "}\n"
	default:
		// Auto mode - use prefers-color-scheme by default, but allow manual override
		css += "@media (prefers-color-scheme: dark) {\n"
		css += "  :root:not(.light-mode) {\n"
		css += "  " + addColor("--bg", colors.Dark.Background)
		css += "  " + addColor("--fg", colors.Dark.Foreground)
		css += "  " + addColor("--muted", colors.Dark.Muted)
		css += "  " + addColor("--border", colors.Dark.Border)
		css += "  " + addColor("--link", colors.Dark.Accent)
		css += "  " + addColor("--code-bg", colors.Dark.CodeBackground)
		if colors.Dark.HeaderBackground != "" {
			css += "  " + addColor("--header-bg", colors.Dark.HeaderBackground)
		}
		if colors.Dark.SidebarBackground != "" {
			css += "  " + addColor("--sidebar-bg", colors.Dark.SidebarBackground)
		}
		css += "  }\n"
		css += "}\n\n"
		// Manual dark mode override (works even in light system preference)
		css += ":root.dark-mode {\n"
		css += addColor("--bg", colors.Dark.Background)
		css += addColor("--fg", colors.Dark.Foreground)
		css += addColor("--muted", colors.Dark.Muted)
		css += addColor("--border", colors.Dark.Border)
		css += addColor("--link", colors.Dark.Accent)
		css += addColor("--code-bg", colors.Dark.CodeBackground)
		if colors.Dark.HeaderBackground != "" {
			css += addColor("--header-bg", colors.Dark.HeaderBackground)
		}
		if colors.Dark.SidebarBackground != "" {
			css += addColor("--sidebar-bg", colors.Dark.SidebarBackground)
		}
		css += "}\n\n"
		// Manual light mode override (works even in dark system preference)
		css += ":root.light-mode {\n"
		css += addColor("--bg", colors.Light.Background)
		css += addColor("--fg", colors.Light.Foreground)
		css += addColor("--muted", colors.Light.Muted)
		css += addColor("--border", colors.Light.Border)
		css += addColor("--link", colors.Light.Accent)
		css += addColor("--code-bg", colors.Light.CodeBackground)
		if colors.Light.HeaderBackground != "" {
			css += addColor("--header-bg", colors.Light.HeaderBackground)
		}
		if colors.Light.SidebarBackground != "" {
			css += addColor("--sidebar-bg", colors.Light.SidebarBackground)
		}
		css += "}\n"
	}

	return css
}

// getCodeThemeForPreset returns the appropriate Chroma code theme for a given preset.
// This maps theme presets to matching syntax highlighting themes.
func getCodeThemeForPreset(preset string) string {
	// Map presets to appropriate Chroma themes
	presetToCodeTheme := map[string]string{
		"default":     "github",
		"ocean":       "github",
		"forest":      "github",
		"sunset":      "github",
		"midnight":    "github",
		"lavender":    "github",
		"rose":        "github",
		"nord":        "nord",
		"gruvbox":     "github",
		"dracula":     "dracula",
		"solarized":   "solarized-light",
		"mono":        "github",
		"cyberpunk":   "github",
		"desert":      "github",
		"ice":         "github",
		"coffee":      "github",
		"emerald":     "github",
		"amber":       "github",
		"matrix":      "github",
		"vscode-dark": "github",
		"carbon":      "github",
		"sakura":      "github",
		"terminal":    "github",
	}

	if theme, ok := presetToCodeTheme[preset]; ok {
		return theme
	}
	return "github" // fallback
}

// GetCodeTheme returns the code theme to use, determined by the preset.
func (c *Config) GetCodeTheme() string {
	return getCodeThemeForPreset(c.Theme.Preset)
}

// IsMultiLingual returns true if multiple languages are enabled.
func (c *Config) IsMultiLingual() bool {
	return len(c.Languages.Enabled) > 1
}

// GetLanguage returns the Language config for a given code, or nil if not found.
func (c *Config) GetLanguage(code string) *Language {
	for _, lang := range c.Languages.Enabled {
		if lang.Code == code {
			return &lang
		}
	}
	return nil
}

// IsLanguageEnabled returns true if the given language code is enabled.
func (c *Config) IsLanguageEnabled(code string) bool {
	return c.GetLanguage(code) != nil
}

// GetDefaultLanguage returns the default language code, or empty string if not configured.
func (c *Config) GetDefaultLanguage() string {
	if c.Languages.Default != "" {
		return c.Languages.Default
	}
	// If no default set but languages are enabled, use first one
	if len(c.Languages.Enabled) > 0 {
		return c.Languages.Enabled[0].Code
	}
	return ""
}
