package markdown

import (
	"strings"
)

// StripYAMLFrontMatter removes a leading YAML front matter block from a markdown document.
// It only strips if the document starts with "---" on the first line and has a closing "---" on its own line.
// If no such block is found, the input is returned unchanged.
func StripYAMLFrontMatter(md string) string {
	// Normalize: allow leading UTF-8 BOM and whitespace-free start
	s := strings.TrimPrefix(md, "\uFEFF")

	// Must start with front matter delimiter on first line.
	if !strings.HasPrefix(s, "---\n") && s != "---" && !strings.HasPrefix(s, "---\r\n") {
		return md
	}

	// Find closing delimiter on its own line.
	// Accept both LF and CRLF.
	// We search starting after the first line.
	rest := s
	// Skip first delimiter line
	if i := strings.Index(rest, "\n"); i >= 0 {
		rest = rest[i+1:]
	} else {
		// Only "---" present
		return md
	}

	// We'll scan line-by-line to be safe.
	pos := 0
	for pos < len(rest) {
		// Extract next line
		nl := strings.IndexByte(rest[pos:], '\n')
		var line string
		var next int
		if nl == -1 {
			line = strings.TrimSuffix(rest[pos:], "\r")
			next = len(rest)
		} else {
			line = strings.TrimSuffix(rest[pos:pos+nl], "\r")
			next = pos + nl + 1
		}

		if line == "---" {
			// Strip everything up to and including this delimiter line.
			// If there's content immediately after (no newline), include it
			return rest[next:]
		}

		pos = next
	}

	// If we never found a closing delimiter, do not strip.
	return md
}
