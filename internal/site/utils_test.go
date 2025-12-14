package site

import (
	"testing"
)

func TestIsIndexRel(t *testing.T) {
	tests := []struct {
		name     string
		rel      string
		expected bool
	}{
		{"root index", "index.md", true},
		{"folder index", "guide/index.md", true},
		{"deep folder index", "api/v1/index.md", true},
		{"regular file", "getting-started.md", false},
		{"nested file", "guide/intro.md", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIndexRel(tt.rel)
			if result != tt.expected {
				t.Errorf("isIndexRel(%q) = %v; want %v", tt.rel, result, tt.expected)
			}
		})
	}
}

func TestTitleFromIndexRel(t *testing.T) {
	tests := []struct {
		name     string
		rel      string
		expected string
	}{
		{"root index", "index.md", "Home"},
		{"folder index", "guide/index.md", "Guide"},
		{"deep folder index", "api/v1/index.md", "V1"},
		{"nested folder", "docs/guide/index.md", "Guide"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := titleFromIndexRel(tt.rel)
			if result != tt.expected {
				t.Errorf("titleFromIndexRel(%q) = %q; want %q", tt.rel, result, tt.expected)
			}
		})
	}
}

func TestExtractNumericPrefix(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected int
	}{
		{"with prefix", "01_getting-started.md", 1},
		{"double digit", "12_advanced.md", 12},
		{"triple digit", "123_deep.md", 123},
		{"with dash", "01-geting-started.md", 1},
		{"with space", "01 getting-started.md", 1},
		{"no prefix", "getting-started.md", 0},
		{"nested path", "guide/01_intro.md", 1},
		{"uppercase extension", "01_FILE.MD", 1},
		{"no number", "file.md", 0},
		{"number in middle", "file_01_name.md", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractNumericPrefix(tt.filename)
			if result != tt.expected {
				t.Errorf("extractNumericPrefix(%q) = %d; want %d", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool // whether parsing should succeed
	}{
		{"RFC3339", "2024-01-15T10:30:00Z", true},
		{"date only", "2024-01-15", true},
		{"date and time", "2024-01-15 10:30", true},
		{"full datetime", "2024-01-15 10:30:45", true},
		{"invalid", "not-a-date", false},
		{"empty", "", false},
		{"wrong format", "15/01/2024", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := parseDate(tt.input)
			if ok != tt.expected {
				t.Errorf("parseDate(%q) success = %v; want %v", tt.input, ok, tt.expected)
			}
		})
	}
}
