package config

import (
	"os"
	"path/filepath"
	"strings"
)

// DetectStructure scans the docs directory and auto-detects language and version folders.
// Uses MkDocs-style structure: direct language folders (docs/en/, docs/de/) and version folders (docs/v1/, docs/latest/).
// Returns detected languages and versions, or empty if nothing is found.
func DetectStructure(docsDir string) (detectedLangs []Language, detectedVersions []Version) {
	absDir, err := filepath.Abs(docsDir)
	if err != nil {
		return nil, nil
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, nil
	}

	// Track what we find
	langMap := make(map[string]bool)
	versionMap := make(map[string]bool)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Check for MkDocs-style structure: direct language/version folders
		// Language folders: typically 2-5 character codes (en, de, fr, zh-CN, etc.)
		// Must be alphanumeric and NOT look like a version
		if len(name) >= 2 && len(name) <= 5 && isAlphanumeric(name) {
			// Exclude version-like names (v1, v2, etc.) and "latest" from language detection
			// Use isStrictVersionID to check if it's actually a version pattern
			if name != "latest" && !isStrictVersionID(name) {
				// Additional check: language codes are typically lowercase letters, possibly with dashes
				// Exclude anything that starts with "v" followed by numbers (version pattern)
				isLikelyLanguage := strings.ToLower(name) == name && !strings.HasPrefix(name, "v")
				// Also exclude if it contains numbers (versions have numbers, languages typically don't)
				if !strings.ContainsAny(name, "0123456789") || isLikelyLanguage {
					// Check if it looks like a language folder by checking if it contains docs
					langPath := filepath.Join(absDir, name)
					if hasMarkdownFiles(langPath) {
						langMap[name] = true
					}
				}
			}
		}

		// Version folders: "latest" or version-like strings (v1, v2, 1.0, 2.0.0, etc.)
		// Must be strict: only detect actual version patterns, not content folders
		// Versions must: start with "v" followed by numbers, be "latest", or be pure number patterns
		if name == "latest" {
			versionPath := filepath.Join(absDir, name)
			if hasMarkdownFiles(versionPath) {
				versionMap[name] = true
			}
		} else if isStrictVersionID(name) {
			// isStrictVersionID already ensures it's a version pattern (has numbers, v prefix, etc.)
			// So if it passes isStrictVersionID, it's definitely a version, not a language
			versionPath := filepath.Join(absDir, name)
			if hasMarkdownFiles(versionPath) {
				versionMap[name] = true
			}
		}
	}

	// Convert maps to slices
	for langCode := range langMap {
		detectedLangs = append(detectedLangs, Language{
			Code: langCode,
			Name: langCode, // Default name to code, can be overridden in config
		})
	}

	for versionID := range versionMap {
		detectedVersions = append(detectedVersions, Version{
			ID:   versionID,
			Name: versionID, // Default name to ID, can be overridden in config
		})
	}

	return detectedLangs, detectedVersions
}

// AutoDetectAndEnhanceConfig is deprecated - auto-detection has been removed.
// Languages and versions must now be explicitly configured in dorcs.yaml.
// This function is kept for backward compatibility but does nothing.
func AutoDetectAndEnhanceConfig(cfg *Config, docsDir string) {
	// Auto-detection has been removed - explicit configuration is required
	// This function is a no-op for backward compatibility
}

// Helper functions

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

func isValidVersionID(s string) bool {
	// Version IDs can be: v1, v2.0, 1.0, 2.0.0, latest, etc.
	// Allow alphanumeric, dots, dashes, underscores
	if len(s) == 0 || len(s) > 20 {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

// isStrictVersionID checks if a string looks like a version identifier.
// Much stricter than isValidVersionID - only matches actual version patterns:
// - Starts with "v" followed by numbers (v1, v2, v1.0, v2.0.0)
// - Pure number patterns (1, 1.0, 2.0.0)
// - Does NOT match content folders like "06_markdown" or "external-content"
func isStrictVersionID(s string) bool {
	if len(s) == 0 || len(s) > 20 {
		return false
	}

	// Must contain at least one digit to be a version
	if !strings.ContainsAny(s, "0123456789") {
		return false
	}

	// Pattern 1: Starts with "v" followed by numbers/dots (v1, v2, v1.0, v2.0.0)
	if strings.HasPrefix(s, "v") || strings.HasPrefix(s, "V") {
		rest := s[1:]
		// Rest must be numbers and dots only
		for _, r := range rest {
			if !((r >= '0' && r <= '9') || r == '.') {
				return false
			}
		}
		return len(rest) > 0 // Must have something after "v"
	}

	// Pattern 2: Pure number pattern (1, 1.0, 2.0.0) - only numbers and dots
	hasOnlyNumbersAndDots := true
	for _, r := range s {
		if !((r >= '0' && r <= '9') || r == '.') {
			hasOnlyNumbersAndDots = false
			break
		}
	}
	if hasOnlyNumbersAndDots {
		return true
	}

	// Reject anything with underscores, dashes, or letters mixed with numbers
	// (e.g., "06_markdown", "external-content" are NOT versions)
	return false
}

func hasMarkdownFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(strings.ToLower(name), ".md") {
				return true
			}
		}
	}
	return false
}
