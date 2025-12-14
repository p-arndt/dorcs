package site

import (
	"path"
	"path/filepath"
	"strings"
	"time"
)

// keyFromRel converts a relative path like "a/b/c.md" to a URL key.
//
// Special rule for index.md (folder landing pages):
// - "index.md"          -> ""
// - "guide/index.md"    -> "guide"
// - "guide/intro.md"    -> "guide/intro"
func keyFromRel(rel string) string {
	rel = strings.TrimSpace(rel)
	if rel == "" {
		return ""
	}
	rel = filepath.ToSlash(rel)
	if strings.HasPrefix(rel, "/") || strings.Contains(rel, "..") {
		return ""
	}
	if !strings.HasSuffix(strings.ToLower(rel), ".md") {
		return ""
	}

	// Remove ".md"
	key := rel[:len(rel)-3]
	key = strings.TrimPrefix(key, "./")
	key = normalizeKey(key)

	// Apply index.md routing rule
	// If key ends with "/index" -> drop trailing "/index"
	// If key is exactly "index" -> root index -> empty key
	if key == "index" {
		return ""
	}
	if strings.HasSuffix(key, "/index") {
		key = strings.TrimSuffix(key, "/index")
		key = strings.TrimSuffix(key, "/")
		return normalizeKey(key)
	}

	return key
}

// normalizeKey normalizes a URL key by cleaning up separators and removing unsafe patterns.
func normalizeKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, "\\", "/")
	key = strings.TrimPrefix(key, "/")
	key = strings.TrimSuffix(key, "/")
	for strings.Contains(key, "//") {
		key = strings.ReplaceAll(key, "//", "/")
	}

	// Allow empty key, which represents root index (index.md at RootDir).
	if key == "" {
		return ""
	}

	if strings.Contains(key, "..") || strings.HasPrefix(key, ".") {
		// Disallow traversal or hidden-style keys.
		return ""
	}
	return key
}

// titleFromKey generates a title from a URL key.
func titleFromKey(key string) string {
	base := key
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	base = strings.TrimSpace(base)
	if base == "" {
		return "Untitled"
	}
	// Title-case-ish without locale complexity.
	return strings.ToUpper(base[:1]) + base[1:]
}

// dirKeyFromKey extracts the directory key from a document key.
func dirKeyFromKey(key string) string {
	if idx := strings.LastIndex(key, "/"); idx >= 0 {
		return key[:idx]
	}
	return ""
}

// isIndexRel checks if a relative path represents an index.md file.
func isIndexRel(rel string) bool {
	rel = filepath.ToSlash(strings.TrimSpace(rel))
	return rel == "index.md" || strings.HasSuffix(rel, "/index.md")
}

// titleFromIndexRel generates a title from an index.md relative path.
func titleFromIndexRel(rel string) string {
	rel = filepath.ToSlash(strings.TrimSpace(rel))
	rel = strings.TrimSuffix(rel, "/index.md")
	rel = strings.TrimSuffix(rel, "index.md")
	rel = strings.Trim(rel, "/")
	if rel == "" {
		return "Home"
	}
	// Use last folder name.
	base := path.Base(rel)
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	base = strings.TrimSpace(base)
	if base == "" {
		return "Home"
	}
	return strings.ToUpper(base[:1]) + base[1:]
}

// parseDate attempts to parse a date string using common formats.
func parseDate(s string) (time.Time, bool) {
	// Common layouts users tend to put in front matter.
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// extractNumericPrefix extracts a numeric prefix from a filename.
// For example: "01_getting-started.md" -> 1, "02_advanced.md" -> 2
// Returns 0 if no numeric prefix is found.
func extractNumericPrefix(filename string) int {
	// Remove path and extension
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, ".md")
	base = strings.TrimSuffix(base, ".MD")

	// Try to match patterns like "01_", "01-", "01 ", or just "01" at start
	// Match one or more digits at the beginning
	var numStr strings.Builder
	for i, r := range base {
		if r >= '0' && r <= '9' {
			numStr.WriteRune(r)
		} else if i == 0 {
			// First char is not a digit, no prefix
			return 0
		} else {
			// Found a non-digit after digits, stop
			break
		}
	}

	if numStr.Len() == 0 {
		return 0
	}

	// Parse the number
	var num int
	for _, r := range numStr.String() {
		num = num*10 + int(r-'0')
	}
	return num
}
