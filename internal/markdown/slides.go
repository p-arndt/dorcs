package markdown

import (
	"regexp"
	"strings"
)

// slideSeparator matches a horizontal rule on its own line.
// Marp-compatible: ---, ___, ***, or - - - with optional surrounding whitespace.
var slideSeparator = regexp.MustCompile(`(?m)^\s*(---|___|\*\*\*|- - -)\s*$`)

// directivePattern matches Marpit directives: <!-- key: value --> or <!-- _key: value -->
// _ prefix = spot (this slide only). No prefix = inherited (this and following slides).
var directivePattern = regexp.MustCompile(`^\s*<!--\s*(_?)(\w+):\s*(.+)\s*-->\s*$`)

// knownLayouts is the set of recognized layout names used for auto-promotion and normalization.
var knownLayouts = map[string]bool{
	"lead": true, "left": true, "right": true, "big": true, "quote": true,
	"columns-2": true, "columns-3": true, "cols": true, "split": true,
	"timeline": true, "fit": true, "invert": true,
}

// SlideMetadata holds parsed directives for a single slide.
// Marpit-compatible: https://marpit.marp.app/directives
type SlideMetadata struct {
	Class              string // HTML class for slide (extra styling, e.g. invert)
	Layout             string // Layout preset: default, lead, left, right, columns-2, columns-3, timeline, split, cols
	Gap                string // Spacing: tight, normal, loose
	Align              string // Content alignment: start, center, end
	Color              string // CSS color (text)
	BackgroundColor    string
	BackgroundImage    string
	BackgroundPosition string // default: center
	BackgroundRepeat   string // default: no-repeat
	BackgroundSize     string // default: cover
	Header             string
	Footer             string
	Paginate           string // "true", "false", "hold", "skip" - we use true/false
}

func parseDirective(meta *SlideMetadata, key, val string) {
	keyLower := strings.ToLower(key)
	switch keyLower {
	case "class":
		meta.Class = val
	case "layout":
		// Space-separated layout values: first known layout token → Layout, rest → Class.
		parts := strings.Fields(val)
		if len(parts) == 0 {
			return
		}
		layout := normalizeLayout(parts[0])
		meta.Layout = layout
		if len(parts) > 1 {
			extra := strings.Join(parts[1:], " ")
			if meta.Class != "" {
				meta.Class = meta.Class + " " + extra
			} else {
				meta.Class = extra
			}
		}
	case "gap":
		meta.Gap = val
	case "align":
		meta.Align = val
	case "color":
		meta.Color = val
	case "backgroundcolor":
		meta.BackgroundColor = val
	case "backgroundimage":
		meta.BackgroundImage = val
	case "backgroundposition":
		meta.BackgroundPosition = val
	case "backgroundrepeat":
		meta.BackgroundRepeat = val
	case "backgroundsize":
		meta.BackgroundSize = val
	case "header":
		meta.Header = val
	case "footer":
		meta.Footer = val
	case "paginate":
		meta.Paginate = val
	}
}

// normalizeLayout maps layout aliases to canonical names.
func normalizeLayout(layout string) string {
	if layout == "two-columns" {
		return "cols"
	}
	return layout
}

// ParseSlideDirectives extracts Marpit-style directives from the start of a slide chunk.
// Supports: class, color, backgroundColor, backgroundImage, backgroundPosition,
// backgroundRepeat, backgroundSize, header, footer, paginate.
// _ prefix = spot (single slide). No prefix = inherited by following slides.
// Returns (metadata for this slide, remaining content, inherited meta for next slide).
func ParseSlideDirectives(chunk string, inherited SlideMetadata) (SlideMetadata, string, SlideMetadata) {
	meta := inherited
	nextInherited := inherited
	lines := strings.Split(chunk, "\n")
	var kept []string
	inDirectiveBlock := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if inDirectiveBlock {
			if m := directivePattern.FindStringSubmatch(line); m != nil {
				underscore, key, val := m[1], strings.TrimSpace(m[2]), strings.TrimSpace(m[3])
				// Strip surrounding quotes (YAML-style) from value
				if len(val) >= 2 && (val[0] == '"' && val[len(val)-1] == '"' || val[0] == '\'' && val[len(val)-1] == '\'') {
					val = val[1 : len(val)-1]
				}
				spot := underscore == "_"
				parseDirective(&meta, key, val)
				if !spot {
					parseDirective(&nextInherited, key, val)
				}
				continue
			}
			if trimmed == "" {
				continue
			}
			inDirectiveBlock = false
		}
		kept = append(kept, line)
	}

	// _class auto-promotion: if Layout is not set and Class matches a known layout, promote it.
	if meta.Layout == "" && knownLayouts[meta.Class] {
		meta.Layout = meta.Class
	}

	return meta, strings.TrimSpace(strings.Join(kept, "\n")), nextInherited
}

// SplitSlides splits markdown content into slide chunks by horizontal ruler (---).
// The input should have front matter already stripped; the first --- in the document
// starts front matter, so callers must pass body-only content (e.g. after StripYAMLFrontMatter
// or preprocessMarkdown which strips it).
// Marp-compatible: uses ---, ___, ***, or - - - as slide separators.
func SplitSlides(body string) []string {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil
	}

	// Split by slide separator, keeping the separator positions
	chunks := slideSeparator.Split(body, -1)

	var slides []string
	for _, c := range chunks {
		trimmed := strings.TrimSpace(c)
		// Include empty chunks as blank slides (e.g. for transitions)
		slides = append(slides, trimmed)
	}
	return slides
}
