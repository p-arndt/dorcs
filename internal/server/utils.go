// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"path"
	"strings"
	"time"
)

func cleanKey(in string) string {
	s := strings.TrimSpace(in)
	s = strings.ReplaceAll(s, "\\", "/")
	s = path.Clean("/" + s)
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, "/")
	if s == "." || s == ".." || strings.HasPrefix(s, "../") || strings.Contains(s, "/../") {
		return ""
	}
	return s
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// escapeXML escapes special XML characters in a string.
// Optimized to use strings.Builder for better performance.
func escapeXML(s string) string {
	var b strings.Builder
	b.Grow(len(s) + len(s)/4) // Pre-allocate with some extra capacity for escaped chars

	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
