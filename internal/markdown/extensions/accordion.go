package markdown

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// ConvertAccordionBlocksInMarkdown converts accordion blocks to blockquotes that goldmark can process.
// It converts `:::accordion Title\nContent\n:::` to `> **ACCORDION:Title**\n> Content...` format.
// It skips accordion syntax inside code blocks (both fenced and indented).
func ConvertAccordionBlocksInMarkdown(md string) string {
	lines := strings.Split(md, "\n")
	var result []string
	var inFencedCodeBlock bool
	var openingFenceChar string
	var openingFenceLength int
	var inIndentedCodeBlock bool

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Check for fenced code block start/end (``` or ~~~)
		if strings.HasPrefix(trimmed, "```") {
			if !inFencedCodeBlock {
				// Opening fence
				openingFenceChar = "`"
				openingFenceLength = 0
				for j := 0; j < len(trimmed) && trimmed[j] == '`'; j++ {
					openingFenceLength++
				}
				inFencedCodeBlock = true
			} else {
				// Check if closing fence
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
			result = append(result, line)
			i++
			continue
		}

		if strings.HasPrefix(trimmed, "~~~") {
			if !inFencedCodeBlock {
				// Opening fence
				openingFenceChar = "~"
				openingFenceLength = 0
				for j := 0; j < len(trimmed) && trimmed[j] == '~'; j++ {
					openingFenceLength++
				}
				inFencedCodeBlock = true
			} else if openingFenceChar == "~" {
				// Check if closing fence
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
			result = append(result, line)
			i++
			continue
		}

		// Check for indented code blocks (4+ spaces)
		if !inFencedCodeBlock && len(line) > 0 && strings.HasPrefix(line, "    ") && strings.TrimSpace(line) != "" {
			if !inIndentedCodeBlock {
				inIndentedCodeBlock = true
			}
		} else if inIndentedCodeBlock {
			// Check if we're leaving the indented code block
			if strings.TrimSpace(line) == "" {
				// Empty line might continue the code block, but if next line is not indented, we exit
				if i+1 < len(lines) {
					nextLine := lines[i+1]
					if len(nextLine) == 0 || !strings.HasPrefix(nextLine, "    ") {
						inIndentedCodeBlock = false
					}
				}
			} else if !strings.HasPrefix(line, "    ") {
				// Not indented, exit code block
				inIndentedCodeBlock = false
			}
		}

		// Skip processing if inside code blocks
		if inFencedCodeBlock || inIndentedCodeBlock {
			result = append(result, line)
			i++
			continue
		}

		// Check if this line starts an accordion block
		if strings.HasPrefix(trimmed, ":::accordion") {
			// Try to find the complete accordion block
			endIdx := -1

			// Look for closing :::
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) == ":::" {
					endIdx = j
					break
				}
			}

			if endIdx != -1 {
				// Found complete accordion block
				// Extract title (everything after "accordion" on first line)
				titleStart := strings.Index(trimmed, "accordion")
				if titleStart != -1 {
					title := strings.TrimSpace(trimmed[titleStart+9:]) // "accordion" is 9 chars

					// Extract content (lines between start and end)
					var contentLines []string
					for j := i + 1; j < endIdx; j++ {
						contentLines = append(contentLines, lines[j])
					}

					// Convert to blockquote format
					var blockquote strings.Builder
					blockquote.WriteString(fmt.Sprintf("> **ACCORDION:%s**\n", title))

					// Skip leading empty lines
					startIdx := 0
					for startIdx < len(contentLines) && strings.TrimSpace(contentLines[startIdx]) == "" {
						startIdx++
					}

					// Skip trailing empty lines
					contentEndIdx := len(contentLines)
					for contentEndIdx > startIdx && strings.TrimSpace(contentLines[contentEndIdx-1]) == "" {
						contentEndIdx--
					}

					// Process content lines (preserve empty lines in the middle for paragraph breaks)
					for j := startIdx; j < contentEndIdx; j++ {
						contentLine := contentLines[j]
						if strings.TrimSpace(contentLine) != "" {
							blockquote.WriteString("> " + contentLine + "\n")
						} else {
							blockquote.WriteString(">\n")
						}
					}
					result = append(result, blockquote.String())
					i = endIdx + 1 // endIdx is the line index of the closing :::
					continue
				}
			}
		}

		result = append(result, line)
		i++
	}

	return strings.Join(result, "\n")
}

// ConvertAccordionBlocksInHTML converts blockquotes that contain accordion markers to accordion HTML.
// It processes the HTML output from goldmark to find blockquotes with **ACCORDION:Title** markers.
func ConvertAccordionBlocksInHTML(htmlContent string) string {
	// Pattern to match blockquote with accordion marker - very flexible to handle different HTML structures
	accordionBlockquoteRE := regexp.MustCompile(`(?s)<blockquote[^>]*>\s*(?:<p[^>]*>\s*)?<strong>ACCORDION:(.*?)</strong>\s*(?:</p>)?\s*(.*?)</blockquote>`)

	htmlContent = accordionBlockquoteRE.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := accordionBlockquoteRE.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		title := strings.TrimSpace(submatches[1])
		content := strings.TrimSpace(submatches[2])

		if title == "" {
			title = "Click to expand"
		}

		// Clean up content - remove leading/trailing empty <p> tags
		content = regexp.MustCompile(`^\s*<p>\s*</p>\s*`).ReplaceAllString(content, "")
		content = regexp.MustCompile(`\s*<p>\s*</p>\s*$`).ReplaceAllString(content, "")
		// Remove leading/trailing <br> tags (often added by goldmark for empty blockquote lines)
		content = regexp.MustCompile(`^\s*<br\s*/?>\s*`).ReplaceAllString(content, "")
		content = regexp.MustCompile(`\s*<br\s*/?>\s*$`).ReplaceAllString(content, "")
		// Remove standalone <br> tags at the start of content (goldmark sometimes adds these)
		content = regexp.MustCompile(`^<br\s*/?>\s*`).ReplaceAllString(content, "")
		// Remove any remaining leading/trailing whitespace
		content = strings.TrimSpace(content)

		// Build accordion HTML - don't escape content as it's already HTML from goldmark
		return fmt.Sprintf(`<details class="dorcs-accordion"><summary>%s</summary><div class="dorcs-accordion-content">%s</div></details>`,
			html.EscapeString(title), content)
	})

	return htmlContent
}
