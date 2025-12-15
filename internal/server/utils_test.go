package server

import (
	"strings"
	"testing"
	"time"
)

func TestCleanKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple path", "getting-started", "getting-started"},
		{"nested path", "guide/getting-started", "guide/getting-started"},
		{"with leading slash", "/getting-started", "getting-started"},
		{"with trailing slash", "getting-started/", "getting-started"},
		{"with both slashes", "/getting-started/", "getting-started"},
		{"with backslashes", "guide\\getting-started", "guide/getting-started"},
		{"with spaces", "  getting-started  ", "getting-started"},
		{"empty string", "", ""},
		{"just slashes", "///", ""},
		{"path traversal dot", ".", ""},
		{"path traversal dot dot", "..", ""},
		{"path traversal prefix", "../etc/passwd", "etc/passwd"},       // path.Clean normalizes to "etc/passwd" (security handled elsewhere)
		{"path traversal middle", "guide/../etc/passwd", "etc/passwd"}, // path.Clean normalizes to "etc/passwd" (security handled elsewhere)
		{"clean path", "guide/advanced/topics", "guide/advanced/topics"},
		{"root path", "/", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanKey(tt.input)
			if got != tt.expected {
				t.Errorf("cleanKey(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"valid date", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), "2024-01-15"},
		{"zero time", time.Time{}, ""},
		{"another date", time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "2023-12-25"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDate(tt.input)
			if got != tt.expected {
				t.Errorf("formatDate(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no special chars", "Hello World", "Hello World"},
		{"ampersand", "A & B", "A &amp; B"},
		{"less than", "A < B", "A &lt; B"},
		{"greater than", "A > B", "A &gt; B"},
		{"double quote", `Say "Hello"`, `Say &quot;Hello&quot;`},
		{"single quote", "It's great", "It&apos;s great"},
		{"all special chars", `<tag attr="value">Text & More</tag>`, `&lt;tag attr=&quot;value&quot;&gt;Text &amp; More&lt;/tag&gt;`},
		{"empty string", "", ""},
		{"URL with query", "https://example.com?q=test&lang=en", "https://example.com?q=test&amp;lang=en"},
		{"mixed content", `John's "favorite" <tag>`, `John&apos;s &quot;favorite&quot; &lt;tag&gt;`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeXML(tt.input)
			if got != tt.expected {
				t.Errorf("escapeXML(%q) = %q, want %q", tt.input, got, tt.expected)
			}
			// Verify the result doesn't contain unescaped special characters
			if strings.Contains(got, "&") && !strings.Contains(got, "&amp;") && !strings.Contains(got, "&lt;") && !strings.Contains(got, "&gt;") && !strings.Contains(got, "&quot;") && !strings.Contains(got, "&apos;") {
				// This is a bit tricky - we need to check that standalone & is escaped
				if strings.Count(got, "&") > strings.Count(got, "&amp;")+strings.Count(got, "&lt;")+strings.Count(got, "&gt;")+strings.Count(got, "&quot;")+strings.Count(got, "&apos;") {
					t.Errorf("escapeXML(%q) contains unescaped &", tt.input)
				}
			}
		})
	}
}

func TestEscapeXMLRoundTrip(t *testing.T) {
	// Test that escaped XML can be safely used in XML context
	testCases := []string{
		"Simple text",
		"Text with & ampersand",
		"Text with < less than",
		"Text with > greater than",
		`Text with " quotes`,
		"Text with ' apostrophe",
		"Complex: <tag attr=\"value\">Text & More</tag>",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			escaped := escapeXML(tc)
			// Verify escaped string doesn't contain unescaped XML special chars
			if strings.Contains(escaped, "<") && !strings.Contains(escaped, "&lt;") {
				t.Errorf("escapeXML(%q) contains unescaped <", tc)
			}
			if strings.Contains(escaped, ">") && !strings.Contains(escaped, "&gt;") {
				t.Errorf("escapeXML(%q) contains unescaped >", tc)
			}
		})
	}
}
