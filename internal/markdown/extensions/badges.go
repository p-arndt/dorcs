package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

var badgePattern = regexp.MustCompile(`\{badge:([^}]+)\}`)

// knownBadgeTypes maps known badge names to CSS classes.
var knownBadgeTypes = map[string]string{
	"new":          "badge-new",
	"beta":         "badge-beta",
	"deprecated":   "badge-deprecated",
	"experimental": "badge-experimental",
	"required":     "badge-required",
}

// ConvertBadgesInHTML replaces {badge:TYPE} patterns in HTML with styled badge spans.
// Supports: NEW, BETA, DEPRECATED, EXPERIMENTAL, REQUIRED and custom labels.
func ConvertBadgesInHTML(htmlContent string) string {
	if !strings.Contains(htmlContent, "{badge:") {
		return htmlContent
	}

	return badgePattern.ReplaceAllStringFunc(htmlContent, func(match string) string {
		label := match[7 : len(match)-1] // strip {badge: and }
		cssClass := "badge"
		lower := strings.ToLower(strings.TrimSpace(label))
		if cls, ok := knownBadgeTypes[lower]; ok {
			cssClass += " " + cls
		} else {
			cssClass += " badge-new" // default style for unknown badges
		}
		return fmt.Sprintf(`<span class="%s">%s</span>`, cssClass, strings.TrimSpace(label))
	})
}
