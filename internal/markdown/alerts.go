package markdown

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// ConvertAlertBlocksInMarkdown converts GitHub-style alert blocks to a format goldmark can process.
// It converts `> [!NOTE]` to `> **ALERT:NOTE**` so we can identify and convert them later in HTML.
func ConvertAlertBlocksInMarkdown(md string) string {
	lines := strings.Split(md, "\n")
	var result []string
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Check if this line starts an alert block: > [!TYPE]
		if strings.HasPrefix(line, "> [!") && strings.Contains(line, "]") {
			// Extract alert type
			start := strings.Index(line, "[!")
			end := strings.Index(line, "]")
			if start != -1 && end != -1 && end > start+2 {
				alertType := strings.ToUpper(line[start+2 : end])
				// Map INFO to NOTE
				if alertType == "INFO" {
					alertType = "NOTE"
				}

				// Valid alert types
				validTypes := map[string]bool{
					"NOTE":      true,
					"TIP":       true,
					"IMPORTANT": true,
					"WARNING":   true,
					"CAUTION":   true,
				}

				if validTypes[alertType] {
					// Replace the type line with a marker
					result = append(result, fmt.Sprintf("> **ALERT:%s**", alertType))
					i++
					// Keep the rest of the blockquote lines as-is
					for i < len(lines) {
						if strings.HasPrefix(lines[i], "> ") {
							result = append(result, lines[i])
							i++
						} else if strings.TrimSpace(lines[i]) == "" && i+1 < len(lines) {
							// Check if next line continues the alert block
							if i+1 < len(lines) && strings.HasPrefix(lines[i+1], "> ") {
								result = append(result, lines[i])
								i++
								continue
							}
							break
						} else {
							break
						}
					}
					continue
				}
			}
		}

		result = append(result, line)
		i++
	}

	return strings.Join(result, "\n")
}

// ConvertAlertBlocksInHTML converts blockquotes that contain alert markers to alert divs.
// It processes the HTML output from goldmark to find blockquotes with **ALERT:TYPE** markers.
func ConvertAlertBlocksInHTML(htmlContent string) string {
	validTypes := map[string]string{
		"NOTE":      "note",
		"TIP":       "tip",
		"IMPORTANT": "important",
		"WARNING":   "warning",
		"CAUTION":   "caution",
	}

	// First, try to match blockquotes with the ALERT: marker in strong tags
	// Pattern to match blockquote with alert marker - very flexible to handle different HTML structures
	alertBlockquoteRE := regexp.MustCompile(`(?s)<blockquote[^>]*>\s*(?:<p[^>]*>\s*)?<strong>ALERT:(?P<type>NOTE|TIP|IMPORTANT|WARNING|CAUTION)</strong>\s*(?:</p>)?\s*(?P<content>.*?)</blockquote>`)

	htmlContent = alertBlockquoteRE.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := alertBlockquoteRE.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		alertType := strings.ToUpper(submatches[1])
		content := submatches[2]

		cssClass, ok := validTypes[alertType]
		if !ok {
			return match
		}

		return buildAlertHTML(alertType, cssClass, content)
	})

	// Fallback: also catch blockquotes that contain "ALERT:NOTE" etc as plain text
	// This handles cases where the markdown preprocessing didn't work
	fallbackRE := regexp.MustCompile(`(?s)<blockquote[^>]*>\s*<p[^>]*>\s*ALERT:(?P<type>NOTE|TIP|IMPORTANT|WARNING|CAUTION)\s*</p>\s*(?P<content>.*?)</blockquote>`)

	htmlContent = fallbackRE.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := fallbackRE.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		alertType := strings.ToUpper(submatches[1])
		content := submatches[2]

		cssClass, ok := validTypes[alertType]
		if !ok {
			return match
		}

		return buildAlertHTML(alertType, cssClass, content)
	})

	return htmlContent
}

// buildAlertHTML constructs the HTML for an alert block
func buildAlertHTML(alertType, cssClass, content string) string {
	// Get title
	title := alertType
	if alertType == "IMPORTANT" {
		title = "Important"
	} else {
		// Capitalize first letter
		if len(alertType) > 0 {
			title = strings.ToUpper(alertType[:1]) + strings.ToLower(alertType[1:])
		}
	}

	// Clean up content - remove leading/trailing whitespace
	content = strings.TrimSpace(content)

	// Remove leading/trailing <br> tags
	content = regexp.MustCompile(`^\s*<br\s*/?>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*<br\s*/?>\s*$`).ReplaceAllString(content, "")

	// Remove any empty <p> tags
	content = regexp.MustCompile(`\s*<p>\s*</p>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`^\s*<p>\s*</p>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*<p>\s*</p>\s*$`).ReplaceAllString(content, "")

	// Remove leading/trailing closing </p> tags without opening
	content = regexp.MustCompile(`^\s*</p>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*</p>\s*$`).ReplaceAllString(content, "")

	// Remove multiple consecutive <br> tags (more than 2)
	content = regexp.MustCompile(`(<br\s*/?>\s*){3,}`).ReplaceAllString(content, "<br><br>")

	// Clean up any remaining whitespace
	content = strings.TrimSpace(content)

	// If content is empty, add a placeholder
	if content == "" {
		content = "<p></p>"
	}

	// Build the alert HTML with proper structure
	return fmt.Sprintf(`<div class="alert alert-%s"><p class="alert-title">%s</p><div class="alert-content">%s</div></div>`,
		cssClass, html.EscapeString(title), content)
}
