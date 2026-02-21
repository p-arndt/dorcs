package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// textBlockTypes are the supported typography block names.
var textBlockTypes = map[string]bool{
	"hero": true, "stat": true, "caption": true, "label": true,
}

// isTextBlockStart returns true if the line starts a known text block (::: hero, ::: stat, etc.).
func isTextBlockStart(trimmed string) (blockType string, ok bool) {
	if !strings.HasPrefix(trimmed, ":::") {
		return "", false
	}
	rest := strings.TrimSpace(trimmed[3:])
	if rest == "" {
		return "", false
	}
	// Get first word
	parts := strings.Fields(rest)
	if len(parts) == 0 {
		return "", false
	}
	bt := strings.ToLower(parts[0])
	if textBlockTypes[bt] {
		return bt, true
	}
	return "", false
}

// ConvertTextBlocksInMarkdown converts typography blocks to blockquotes that goldmark can process.
// Supports: ::: hero, ::: stat, ::: caption, ::: label
func ConvertTextBlocksInMarkdown(md string) string {
	lines := strings.Split(md, "\n")
	var result []string
	var inFencedCodeBlock bool
	var inIndentedCodeBlock bool

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Check for fenced code blocks
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFencedCodeBlock = !inFencedCodeBlock
			result = append(result, line)
			i++
			continue
		}

		if len(line) > 0 && strings.HasPrefix(line, "    ") && strings.TrimSpace(line) != "" {
			inIndentedCodeBlock = true
		} else if inIndentedCodeBlock {
			if strings.TrimSpace(line) == "" && i+1 < len(lines) && (len(lines[i+1]) == 0 || !strings.HasPrefix(lines[i+1], "    ")) {
				inIndentedCodeBlock = false
			} else if len(line) > 0 && !strings.HasPrefix(line, "    ") {
				inIndentedCodeBlock = false
			}
		}

		if inFencedCodeBlock || inIndentedCodeBlock {
			result = append(result, line)
			i++
			continue
		}

		if blockType, ok := isTextBlockStart(trimmed); ok {
			endIdx := -1
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) == ":::" {
					endIdx = j
					break
				}
			}

			if endIdx != -1 {
				var contentLines []string
				for j := i + 1; j < endIdx; j++ {
					contentLines = append(contentLines, lines[j])
				}

				converted := convertTextBlockContent(blockType, contentLines)
				result = append(result, converted)
				i = endIdx + 1
				continue
			}
		}

		result = append(result, line)
		i++
	}

	return strings.Join(result, "\n")
}

func convertTextBlockContent(blockType string, contentLines []string) string {
	content := trimTextBlockLines(contentLines)
	if len(content) == 0 {
		return ""
	}

	if blockType == "stat" {
		// First line = value, rest = caption (or first line only if no rest)
		value := content[0]
		var caption []string
		if len(content) > 1 {
			caption = content[1:]
		}
		var out strings.Builder
		out.WriteString("> **TEXTBLOCK_STAT**\n")
		out.WriteString("> " + value + "\n")
		for _, c := range caption {
			if strings.TrimSpace(c) != "" {
				out.WriteString("> " + c + "\n")
			} else {
				out.WriteString(">\n")
			}
		}
		return strings.TrimSuffix(out.String(), "\n")
	}

	// hero, caption, label - simple wrapper
	var out strings.Builder
	out.WriteString(fmt.Sprintf("> **TEXTBLOCK:%s**\n", blockType))
	for _, c := range content {
		if strings.TrimSpace(c) != "" {
			out.WriteString("> " + c + "\n")
		} else {
			out.WriteString(">\n")
		}
	}
	return strings.TrimSuffix(out.String(), "\n")
}

func trimTextBlockLines(lines []string) []string {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return lines[start:end]
}

// textBlockRE matches blockquote with TEXTBLOCK:type marker
var textBlockRE = regexp.MustCompile(`(?s)<blockquote[^>]*>\s*(?:<p[^>]*>\s*)?<strong>TEXTBLOCK:(hero|caption|label)</strong>\s*(?:</p>)?\s*(.*?)</blockquote>`)

// textBlockStatRE matches blockquote with TEXTBLOCK_STAT
var textBlockStatRE = regexp.MustCompile(`(?s)<blockquote[^>]*>\s*(?:<p[^>]*>\s*)?<strong>TEXTBLOCK_STAT</strong>\s*(?:</p>)?\s*(.*?)</blockquote>`)

// ConvertTextBlocksInHTML converts blockquotes with TEXTBLOCK markers into typography block HTML.
func ConvertTextBlocksInHTML(htmlContent string) string {
	// Process simple blocks (hero, caption, label)
	htmlContent = textBlockRE.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := textBlockRE.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		blockType := submatches[1]
		content := cleanBlockContent(submatches[2])
		return fmt.Sprintf(`<div class="text-block text-block-%s">%s</div>`, blockType, content)
	})

	// Process stat blocks
	htmlContent = textBlockStatRE.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := textBlockStatRE.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		inner := cleanBlockContent(submatches[1])
		// First <p> = value, rest = caption
		value, caption := splitStatContent(inner)
		return fmt.Sprintf(`<div class="text-block text-block-stat"><div class="text-block-stat-value">%s</div><div class="text-block-stat-caption">%s</div></div>`,
			value, caption)
	})

	return htmlContent
}

func cleanBlockContent(content string) string {
	content = regexp.MustCompile(`^\s*<p>\s*</p>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*<p>\s*</p>\s*$`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`^\s*<br\s*/?>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*<br\s*/?>\s*$`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`^<br\s*/?>\s*`).ReplaceAllString(content, "")
	return strings.TrimSpace(content)
}

// splitStatContent splits HTML content into value (first p) and caption (rest)
func splitStatContent(html string) (value, caption string) {
	// Find first <p>...</p>
	pOpen := strings.Index(html, "<p>")
	if pOpen == -1 {
		return html, ""
	}
	pClose := strings.Index(html[pOpen+3:], "</p>")
	if pClose == -1 {
		return html, ""
	}
	pClose += pOpen + 3
	value = strings.TrimSpace(html[pOpen+3 : pClose])
	caption = strings.TrimSpace(html[pClose+4:])
	return value, caption
}
