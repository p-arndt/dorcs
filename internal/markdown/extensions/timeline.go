package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// headingPattern matches ### or #### at start of line (with optional spaces).
var timelineHeadingPattern = regexp.MustCompile(`(?m)^(#{3,4})\s+(.+)$`)

// isTimelineStart returns true if the line starts a timeline block (::: timeline or :::timeline).
func isTimelineStart(trimmed string) bool {
	return trimmed == "::: timeline" || trimmed == ":::timeline" ||
		strings.HasPrefix(trimmed, "::: timeline ") || strings.HasPrefix(trimmed, ":::timeline ")
}

// ConvertTimelineBlocksInMarkdown converts timeline blocks to blockquotes that goldmark can process.
// Syntax: ::: timeline
// ### 2024 · Jan
// **Step 1 Title**
// Description text.
//
// ### 2024 · Jun
// **Step 2 Title**
// More text.
// :::
//
// Each ### or #### heading starts a new timeline step. The heading text becomes the step label (e.g. date).
func ConvertTimelineBlocksInMarkdown(md string) string {
	lines := strings.Split(md, "\n")
	var result []string
	var inFencedCodeBlock bool
	var inIndentedCodeBlock bool

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Check for fenced code blocks (same logic as accordion)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			if !inFencedCodeBlock {
				inFencedCodeBlock = true
			} else {
				inFencedCodeBlock = false
			}
			result = append(result, line)
			i++
			continue
		}

		if len(line) > 0 && strings.HasPrefix(line, "    ") && strings.TrimSpace(line) != "" {
			inIndentedCodeBlock = true
		} else if inIndentedCodeBlock && strings.TrimSpace(line) == "" {
			if i+1 < len(lines) && (len(lines[i+1]) == 0 || !strings.HasPrefix(lines[i+1], "    ")) {
				inIndentedCodeBlock = false
			}
		} else if inIndentedCodeBlock && len(line) > 0 && !strings.HasPrefix(line, "    ") {
			inIndentedCodeBlock = false
		}

		if inFencedCodeBlock || inIndentedCodeBlock {
			result = append(result, line)
			i++
			continue
		}

		// Check for ::: timeline
		if isTimelineStart(trimmed) {
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

				// Split content by ### or #### headings
				converted := convertTimelineContent(contentLines)
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

// convertTimelineContent splits content by h3/h4 headings and converts each step to a blockquote.
func convertTimelineContent(contentLines []string) string {
	var steps []struct {
		label   string
		content []string
	}

	var currentLabel string
	var currentContent []string

	for _, line := range contentLines {
		if matches := timelineHeadingPattern.FindStringSubmatch(line); matches != nil {
			if currentLabel != "" || len(currentContent) > 0 {
				content := trimEmptyLines(currentContent)
				if currentLabel == "" {
					currentLabel = "Step"
				}
				steps = append(steps, struct {
					label   string
					content []string
				}{currentLabel, content})
			}
			currentLabel = strings.TrimSpace(matches[2])
			currentContent = nil
		} else {
			currentContent = append(currentContent, line)
		}
	}

	if currentLabel != "" || len(currentContent) > 0 {
		content := trimEmptyLines(currentContent)
		if currentLabel == "" {
			currentLabel = "Step"
		}
		steps = append(steps, struct {
			label   string
			content []string
		}{currentLabel, content})
	}

	if len(steps) == 0 {
		return ""
	}

	var out strings.Builder
	for _, step := range steps {
		out.WriteString(fmt.Sprintf("> **TIMELINE_STEP:%s**\n", step.label))
		for _, c := range step.content {
			if strings.TrimSpace(c) != "" {
				out.WriteString("> " + c + "\n")
			} else {
				out.WriteString(">\n")
			}
		}
		out.WriteString("\n")
	}
	return strings.TrimSuffix(out.String(), "\n\n")
}

func trimEmptyLines(lines []string) []string {
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

// timelineBlockquoteRE matches a blockquote starting with TIMELINE_STEP:label
var timelineBlockquoteRE = regexp.MustCompile(`(?s)<blockquote[^>]*>\s*(?:<p[^>]*>\s*)?<strong>TIMELINE_STEP:(.*?)</strong>\s*(?:</p>)?\s*(.*?)</blockquote>`)

// ConvertTimelineBlocksInHTML converts blockquotes with TIMELINE_STEP markers into timeline HTML structure.
// Consecutive TIMELINE_STEP blockquotes are wrapped in <div class="timeline"> with each step in <div class="timeline-step">.
func ConvertTimelineBlocksInHTML(htmlContent string) string {
	var result strings.Builder
	pos := 0

	for pos < len(htmlContent) {
		sub := htmlContent[pos:]
		loc := timelineBlockquoteRE.FindStringIndex(sub)
		if loc == nil {
			result.WriteString(sub)
			break
		}

		// Copy content before this match
		result.WriteString(sub[:loc[0]])

		// Collect consecutive TIMELINE_STEP blockquotes
		var steps []struct {
			label   string
			content string
		}
		searchStart := pos + loc[0]
		for searchStart < len(htmlContent) {
			sub := htmlContent[searchStart:]
			loc := timelineBlockquoteRE.FindStringIndex(sub)
			if loc == nil {
				break
			}
			submatches := timelineBlockquoteRE.FindStringSubmatch(sub)
			if len(submatches) < 3 {
				searchStart += loc[1]
				break
			}
			steps = append(steps, struct {
				label   string
				content string
			}{
				label:   strings.TrimSpace(submatches[1]),
				content: cleanTimelineContent(submatches[2]),
			})
			searchStart += loc[1]

			// Skip whitespace to find next potential blockquote
			for searchStart < len(htmlContent) && (htmlContent[searchStart] == ' ' || htmlContent[searchStart] == '\t' || htmlContent[searchStart] == '\n') {
				searchStart++
			}
		}

		// Build timeline HTML
		result.WriteString(`<div class="timeline">`)
		for _, step := range steps {
			result.WriteString(`<div class="timeline-step">`)
			result.WriteString(`<div class="timeline-marker">`)
			result.WriteString(escapeHTML(step.label))
			result.WriteString(`</div>`)
			result.WriteString(`<div class="timeline-content">`)
			result.WriteString(step.content)
			result.WriteString(`</div>`)
			result.WriteString(`</div>`)
		}
		result.WriteString(`</div>`)

		pos = searchStart
	}

	return result.String()
}

func cleanTimelineContent(content string) string {
	content = regexp.MustCompile(`^\s*<p>\s*</p>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*<p>\s*</p>\s*$`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`^\s*<br\s*/?>\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\s*<br\s*/?>\s*$`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`^<br\s*/?>\s*`).ReplaceAllString(content, "")
	return strings.TrimSpace(content)
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
