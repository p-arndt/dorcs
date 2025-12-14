package markdown

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	gmast "github.com/yuin/goldmark/ast"
)

// tocItem represents a single table of contents item.
type tocItem struct {
	Level int
	ID    string
	Text  string
}

// BuildTOC parses markdown headings from the source and generates a nested <ul>...</ul> table of contents.
// It includes h2/h3/h4 with visual hierarchy (h2 = top-level, h3 = nested, h4 = more nested).
// If no headings are found, it returns empty.
func BuildTOC(md goldmark.Markdown, markdownSource string) template.HTML {
	src := []byte(markdownSource)

	// Parse to AST so we can inspect headings and their auto IDs.
	ctx := parser.NewContext()
	reader := text.NewReader(src)
	doc := md.Parser().Parse(reader, parser.WithContext(ctx))

	var items []tocItem

	// Walk the AST, collect headings
	_ = gmast.Walk(doc, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
		if !entering {
			return gmast.WalkContinue, nil
		}
		h, ok := n.(*gmast.Heading)
		if !ok {
			return gmast.WalkContinue, nil
		}

		// Include H2, H3, H4 for better hierarchy
		if h.Level < 2 || h.Level > 4 {
			return gmast.WalkContinue, nil
		}

		// goldmark stores attributes as interface{}; for heading ids it's typically []byte.
		idVal, _ := h.AttributeString("id")
		var idStr string
		switch v := idVal.(type) {
		case []byte:
			idStr = strings.TrimSpace(string(v))
		case string:
			idStr = strings.TrimSpace(v)
		default:
			idStr = ""
		}
		if idStr == "" {
			// AutoHeadingID should provide this; if not, skip.
			return gmast.WalkContinue, nil
		}

		txt := strings.TrimSpace(ExtractNodeText(h, src))
		if txt == "" {
			return gmast.WalkContinue, nil
		}

		items = append(items, tocItem{
			Level: h.Level,
			ID:    idStr,
			Text:  txt,
		})
		return gmast.WalkContinue, nil
	})

	if len(items) == 0 {
		return ""
	}

	// Build properly nested list HTML.
	// h2 items are top-level, h3 items nest inside the preceding h2, h4 inside h3, etc.
	var b strings.Builder

	b.WriteString(`<ul class="toc-list">`)

	// Stack tracks how many nested <ul> are open beyond the root.
	// Level 2 = root list (already open), level 3 = one nested, level 4 = two nested.
	baseLevel := 2
	currentLevel := baseLevel

	for i, it := range items {
		level := it.Level
		if level < baseLevel {
			level = baseLevel
		}

		// Close lists if going back up
		for currentLevel > level {
			b.WriteString(`</li></ul>`)
			currentLevel--
		}

		// Open nested lists if going deeper
		for currentLevel < level {
			b.WriteString(`<ul class="toc-nested">`)
			currentLevel++
		}

		// Close previous sibling li if same level (except first item)
		if i > 0 && level == items[i-1].Level {
			b.WriteString(`</li>`)
		}

		// Write the item
		b.WriteString(`<li class="toc-item toc-h`)
		b.WriteString(fmt.Sprintf("%d", it.Level))
		b.WriteString(`"><a href="#`)
		b.WriteString(html.EscapeString(it.ID))
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(it.Text))
		b.WriteString(`</a>`)
		// Don't close </li> yet — a nested list might follow
	}

	// Close all remaining open items and lists
	for currentLevel >= baseLevel {
		b.WriteString(`</li>`)
		if currentLevel > baseLevel {
			b.WriteString(`</ul>`)
		}
		currentLevel--
	}
	b.WriteString(`</ul>`)

	return template.HTML(b.String())
}

// ExtractNodeText extracts plain text from a goldmark AST node using the source buffer.
func ExtractNodeText(n gmast.Node, source []byte) string {
	var b strings.Builder
	var walk func(node gmast.Node)
	walk = func(node gmast.Node) {
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			switch t := c.(type) {
			case *gmast.Text:
				seg := t.Segment
				b.Write(seg.Value(source))
			case *gmast.CodeSpan:
				// approximate: include raw text inside code span by walking children
				walk(t)
			default:
				walk(c)
			}
		}
	}
	walk(n)
	return b.String()
}

// ProcessTOCPlaceholder replaces [[TOC]] placeholders in markdown with the provided TOC HTML.
// The TOC HTML should already be properly escaped HTML.
func ProcessTOCPlaceholder(markdownSource string, tocHTML template.HTML) string {
	if tocHTML == "" {
		// If no TOC, just remove the placeholder
		return strings.ReplaceAll(markdownSource, "[[TOC]]", "")
	}
	// Replace [[TOC]] with the HTML TOC
	// We'll insert it as raw HTML which goldmark will render (since we have WithUnsafe enabled)
	return strings.ReplaceAll(markdownSource, "[[TOC]]", string(tocHTML))
}
