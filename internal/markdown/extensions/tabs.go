package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// ConvertTabBlocksInMarkdown converts tab blocks to blockquotes that goldmark can process.
// Syntax:
//
//	:::tabs
//	::tab Title1
//	Content for tab 1...
//	::tab Title2
//	Content for tab 2...
//	:::
//
// Converts to blockquotes with TABS markers for post-processing.
func ConvertTabBlocksInMarkdown(md string) string {
	lines := strings.Split(md, "\n")
	var result []string
	var inFencedCodeBlock bool
	var openingFenceLength int

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Track fenced code blocks
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			fenceChar := trimmed[0]
			fenceLen := 0
			for j := 0; j < len(trimmed) && trimmed[j] == fenceChar; j++ {
				fenceLen++
			}
			if !inFencedCodeBlock {
				inFencedCodeBlock = true
				openingFenceLength = fenceLen
			} else if fenceLen >= openingFenceLength {
				inFencedCodeBlock = false
			}
			result = append(result, line)
			i++
			continue
		}

		if inFencedCodeBlock {
			result = append(result, line)
			i++
			continue
		}

		// Check for :::tabs start
		if trimmed == ":::tabs" || strings.HasPrefix(trimmed, ":::tabs ") {
			// Collect all lines until closing :::
			var tabTitles []string
			var tabBodies [][]string
			var currentBody []string
			i++

			for i < len(lines) {
				innerLine := lines[i]
				innerTrimmed := strings.TrimSpace(innerLine)

				if innerTrimmed == ":::" {
					// Closing fence
					break
				}

				if strings.HasPrefix(innerTrimmed, "::tab ") {
					// Save previous tab body if any
					if len(tabTitles) > 0 {
						tabBodies = append(tabBodies, currentBody)
					}
					title := strings.TrimSpace(innerTrimmed[6:])
					tabTitles = append(tabTitles, title)
					currentBody = nil
					i++
					continue
				}

				currentBody = append(currentBody, innerLine)
				i++
			}

			// Save last tab body
			if len(tabTitles) > 0 {
				tabBodies = append(tabBodies, currentBody)
			}

			if len(tabTitles) == 0 {
				i++
				continue
			}

			// Convert to blockquote format with markers.
			// Each marker needs a blank blockquote line before/after to force
			// goldmark to put them in separate <p> tags.
			result = append(result, "> **DORCS_TABS_START**")
			result = append(result, ">")
			for idx, title := range tabTitles {
				result = append(result, fmt.Sprintf("> **DORCS_TAB:%s**", title))
				result = append(result, ">")
				if idx < len(tabBodies) {
					for _, bodyLine := range tabBodies[idx] {
						result = append(result, "> "+bodyLine)
					}
					result = append(result, ">")
				}
			}
			result = append(result, "> **DORCS_TABS_END**")
			result = append(result, "")
			i++
			continue
		}

		result = append(result, line)
		i++
	}

	return strings.Join(result, "\n")
}

var (
	tabsStartRe = regexp.MustCompile(`<p>\s*<strong>DORCS_TABS_START</strong>\s*</p>`)
	tabsEndRe   = regexp.MustCompile(`<p>\s*<strong>DORCS_TABS_END</strong>\s*</p>`)
	tabTitleRe  = regexp.MustCompile(`<p>\s*<strong>DORCS_TAB:([^<]+)</strong>\s*</p>`)
)

// ConvertTabBlocksInHTML converts blockquotes with DORCS_TABS markers into tabbed HTML.
func ConvertTabBlocksInHTML(htmlContent string) string {
	if !strings.Contains(htmlContent, "DORCS_TABS_START") {
		return htmlContent
	}

	tabGroupID := 0

	for {
		startLoc := tabsStartRe.FindStringIndex(htmlContent)
		if startLoc == nil {
			break
		}

		endLoc := tabsEndRe.FindStringIndex(htmlContent[startLoc[1]:])
		if endLoc == nil {
			break
		}

		// Adjust end location to be relative to full string
		endLoc[0] += startLoc[1]
		endLoc[1] += startLoc[1]

		// Extract the inner HTML between start and end markers
		inner := htmlContent[startLoc[1]:endLoc[0]]

		// Find the enclosing blockquote
		// Look backwards for <blockquote>
		bqStart := strings.LastIndex(htmlContent[:startLoc[0]], "<blockquote>")
		bqEnd := strings.Index(htmlContent[endLoc[1]:], "</blockquote>")

		if bqStart == -1 || bqEnd == -1 {
			break
		}
		bqEnd += endLoc[1] + len("</blockquote>")

		// Parse tabs from inner HTML
		type tabInfo struct {
			title   string
			content string
		}
		var tabs []tabInfo

		matches := tabTitleRe.FindAllStringSubmatchIndex(inner, -1)
		for i, match := range matches {
			title := inner[match[2]:match[3]]
			var content string
			if i+1 < len(matches) {
				content = inner[match[1]:matches[i+1][0]]
			} else {
				content = inner[match[1]:]
			}
			content = strings.TrimSpace(content)
			tabs = append(tabs, tabInfo{title: title, content: content})
		}

		if len(tabs) == 0 {
			break
		}

		// Build tabs HTML
		groupID := fmt.Sprintf("tabs-%d", tabGroupID)
		tabGroupID++

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf(`<div class="dorcs-tabs" data-tab-group="%s">`, groupID))
		sb.WriteString(`<div class="dorcs-tabs-header">`)
		for i, tab := range tabs {
			activeClass := ""
			if i == 0 {
				activeClass = " active"
			}
			sb.WriteString(fmt.Sprintf(`<button class="dorcs-tab-btn%s" data-tab="%d">%s</button>`, activeClass, i, tab.title))
		}
		sb.WriteString(`</div>`)
		for i, tab := range tabs {
			activeClass := ""
			if i == 0 {
				activeClass = " active"
			}
			sb.WriteString(fmt.Sprintf(`<div class="dorcs-tab-panel%s" data-tab="%d">%s</div>`, activeClass, i, tab.content))
		}
		sb.WriteString(`</div>`)

		htmlContent = htmlContent[:bqStart] + sb.String() + htmlContent[bqEnd:]
	}

	return htmlContent
}
