package markdown

import (
	"html/template"
	"regexp"
	"strings"
)

// PlaceholderMap holds all placeholder replacements
type PlaceholderMap map[string]template.HTML

// ProcessPlaceholders replaces all placeholders in markdown with provided HTML.
// Only processes placeholders when they're on their own line and not inside code blocks.
// Supports placeholders like [[TOC]], [[TOC:2]], [[BREADCRUMBS]], etc.
func ProcessPlaceholders(markdownSource string, placeholders PlaceholderMap) string {
	lines := strings.Split(markdownSource, "\n")
	var result strings.Builder
	inFencedCodeBlock := false
	openingFenceChar := "" // "`" or "~"
	openingFenceLength := 0

	// Regex to match placeholders with optional parameters: [[PLACEHOLDER]] or [[PLACEHOLDER:param]]
	placeholderRegex := regexp.MustCompile(`^\[\[([A-Za-z-]+)(?::([^\]]+))?\]\]$`)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for fenced code block start/end
		if strings.HasPrefix(trimmed, "```") {
			if !inFencedCodeBlock {
				// Opening fence - count backticks
				openingFenceChar = "`"
				openingFenceLength = 0
				for j := 0; j < len(trimmed) && trimmed[j] == '`'; j++ {
					openingFenceLength++
				}
				inFencedCodeBlock = true
			} else {
				// Check if this is a closing fence (must have at least as many backticks)
				closingFenceLength := 0
				for j := 0; j < len(trimmed) && trimmed[j] == '`'; j++ {
					closingFenceLength++
				}
				if closingFenceLength >= openingFenceLength {
					inFencedCodeBlock = false
					openingFenceChar = ""
					openingFenceLength = 0
				}
			}
			result.WriteString(line)
			if i < len(lines)-1 {
				result.WriteString("\n")
			}
			continue
		}

		if strings.HasPrefix(trimmed, "~~~") {
			if !inFencedCodeBlock {
				// Opening fence - count tildes
				openingFenceChar = "~"
				openingFenceLength = 0
				for j := 0; j < len(trimmed) && trimmed[j] == '~'; j++ {
					openingFenceLength++
				}
				inFencedCodeBlock = true
			} else if openingFenceChar == "~" {
				// Check if this is a closing fence (must have at least as many tildes)
				closingFenceLength := 0
				for j := 0; j < len(trimmed) && trimmed[j] == '~'; j++ {
					closingFenceLength++
				}
				if closingFenceLength >= openingFenceLength {
					inFencedCodeBlock = false
					openingFenceChar = ""
					openingFenceLength = 0
				}
			}
			result.WriteString(line)
			if i < len(lines)-1 {
				result.WriteString("\n")
			}
			continue
		}

		// If we're inside a code block, don't process placeholders
		if inFencedCodeBlock {
			result.WriteString(line)
			if i < len(lines)-1 {
				result.WriteString("\n")
			}
			continue
		}

		// Check if this line contains a placeholder
		matches := placeholderRegex.FindStringSubmatch(trimmed)
		if matches != nil {
			placeholderName := strings.ToUpper(matches[1])
			param := matches[2]

			// Build key for lookup (with or without parameter)
			lookupKey := placeholderName
			if param != "" {
				lookupKey = placeholderName + ":" + param
			}

			// Try exact match first, then fallback to name without param
			replacement, found := placeholders[lookupKey]
			if !found {
				replacement, found = placeholders[placeholderName]
			}

			if found {
				// Replace with placeholder HTML
				if replacement != "" {
					result.WriteString(string(replacement))
				}
				// If replacement is empty, we skip the line entirely
				if i < len(lines)-1 {
					result.WriteString("\n")
				}
				continue
			}
		}

		// Line doesn't have a placeholder - keep as is
		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
