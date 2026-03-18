package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsMultiLingual(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name:     "no languages configured",
			config:   Default(),
			expected: false,
		},
		{
			name: "single language",
			config: &Config{
				Languages: LanguagesConfig{
					Default: "en",
					Enabled: []Language{
						{Code: "en", Name: "English"},
					},
				},
			},
			expected: true,
		},
		{
			name: "multiple languages",
			config: &Config{
				Languages: LanguagesConfig{
					Default: "en",
					Enabled: []Language{
						{Code: "en", Name: "English"},
						{Code: "de", Name: "Deutsch"},
					},
				},
			},
			expected: true,
		},
		{
			name: "three languages",
			config: &Config{
				Languages: LanguagesConfig{
					Default: "en",
					Enabled: []Language{
						{Code: "en", Name: "English"},
						{Code: "de", Name: "Deutsch"},
						{Code: "fr", Name: "Français"},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsMultiLingual()
			if result != tt.expected {
				t.Errorf("IsMultiLingual() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDefaultLanguage(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name:     "no languages configured",
			config:   Default(),
			expected: "",
		},
		{
			name: "with default language",
			config: &Config{
				Languages: LanguagesConfig{
					Default: "en",
					Enabled: []Language{
						{Code: "en", Name: "English"},
						{Code: "de", Name: "Deutsch"},
					},
				},
			},
			expected: "en",
		},
		{
			name: "different default language",
			config: &Config{
				Languages: LanguagesConfig{
					Default: "de",
					Enabled: []Language{
						{Code: "en", Name: "English"},
						{Code: "de", Name: "Deutsch"},
					},
				},
			},
			expected: "de",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetDefaultLanguage()
			if result != tt.expected {
				t.Errorf("GetDefaultLanguage() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestIsLanguageEnabled(t *testing.T) {
	cfg := &Config{
		Languages: LanguagesConfig{
			Default: "en",
			Enabled: []Language{
				{Code: "en", Name: "English"},
				{Code: "de", Name: "Deutsch"},
				{Code: "fr", Name: "Français"},
			},
		},
	}

	tests := []struct {
		name     string
		langCode string
		expected bool
	}{
		{"enabled language - en", "en", true},
		{"enabled language - de", "de", true},
		{"enabled language - fr", "fr", true},
		{"disabled language - es", "es", false},
		{"disabled language - ja", "ja", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.IsLanguageEnabled(tt.langCode)
			if result != tt.expected {
				t.Errorf("IsLanguageEnabled(%q) = %v; want %v", tt.langCode, result, tt.expected)
			}
		})
	}
}

func TestLoadMultiLingualConfig(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
    - code: "fr"
      name: "Français"
`
	err := os.WriteFile(filepath.Join(dir, "dorcs.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !cfg.IsMultiLingual() {
		t.Error("expected IsMultiLingual() to return true")
	}

	if cfg.GetDefaultLanguage() != "en" {
		t.Errorf("expected default language 'en', got %q", cfg.GetDefaultLanguage())
	}

	if len(cfg.Languages.Enabled) != 3 {
		t.Fatalf("expected 3 enabled languages, got %d", len(cfg.Languages.Enabled))
	}

	if !cfg.IsLanguageEnabled("en") {
		t.Error("expected 'en' to be enabled")
	}
	if !cfg.IsLanguageEnabled("de") {
		t.Error("expected 'de' to be enabled")
	}
	if !cfg.IsLanguageEnabled("fr") {
		t.Error("expected 'fr' to be enabled")
	}
	if cfg.IsLanguageEnabled("es") {
		t.Error("expected 'es' to NOT be enabled")
	}
}
