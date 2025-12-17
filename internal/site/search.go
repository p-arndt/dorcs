package site

import (
	"regexp"
	"sort"
	"strings"

	"github.com/p-arndt/dorcs/internal/markdown"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	gmast "github.com/yuin/goldmark/ast"
)

// headingMatch represents a matching heading in a document.
type headingMatch struct {
	id   string
	text string
	pos  int
}

// SearchDocs performs full-text search across all documents.
// It searches in titles, descriptions, tags, and content.
// Results are sorted by relevance (title matches > description matches > content matches).
func (s *Site) SearchDocs(query string, includeDraft bool, maxResults int) []SearchResult {
	if strings.TrimSpace(query) == "" {
		return nil
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))
	queryWords := strings.Fields(queryLower)

	s.mu.RLock()
	docs := make([]*Doc, 0, len(s.index))
	for _, d := range s.index {
		if !includeDraft && d.Draft {
			continue
		}
		docs = append(docs, d)
	}
	s.mu.RUnlock()

	results := make([]SearchResult, 0)

	for _, doc := range docs {
		// Read document content for searching
		raw, _, _, _, err := markdown.ReadMarkdownStripFrontMatter(doc.FilePath)
		if err != nil {
			continue
		}
		raw = markdown.StripYAMLFrontMatter(raw)
		contentLower := strings.ToLower(raw)

		// Check title matches (highest priority)
		titleLower := strings.ToLower(doc.Title)
		titleMatches := 0
		titleScore := 0
		for _, word := range queryWords {
			if strings.Contains(titleLower, word) {
				titleMatches++
				titleScore += 100 // High score for title matches
			}
		}
		if titleMatches == len(queryWords) {
			titleScore += 50 // Bonus for all words matching in title
		}

		// Check description matches
		descLower := strings.ToLower(doc.Description)
		descScore := 0
		for _, word := range queryWords {
			if strings.Contains(descLower, word) {
				descScore += 30
			}
		}

		// Check tag matches
		tagScore := 0
		for _, tag := range doc.Tags {
			tagLower := strings.ToLower(tag)
			for _, word := range queryWords {
				if strings.Contains(tagLower, word) {
					tagScore += 20
				}
			}
		}

		// Check content matches
		contentMatches := 0
		for _, word := range queryWords {
			if strings.Contains(contentLower, word) {
				contentMatches++
			}
		}

		// Calculate base score (title, description, tags, content)
		baseScore := titleScore + descScore + tagScore + contentMatches

		// Only process if there's at least one match
		if baseScore == 0 {
			continue
		}

		// Find ALL matching headings in this document
		matchingHeadings := findAllMatchingHeadings(s.md, raw, queryWords)

		// If we have matching headings, create a result for each
		if len(matchingHeadings) > 0 {
			for _, heading := range matchingHeadings {
				// Generate snippet around this heading
				snippet := generateSnippetAroundHeading(raw, heading.pos, queryWords, 150)

				// Build path with heading anchor
				path := "/"
				if doc.Key != "" {
					path = "/" + doc.Key
				}
				path += "#" + heading.id

				// Score for this specific heading match
				headingScore := baseScore
				if strings.Contains(strings.ToLower(heading.text), strings.ToLower(queryWords[0])) {
					headingScore += 40 // Bonus for heading match
				}

				results = append(results, SearchResult{
					Doc:         doc,
					Key:         doc.Key,
					Title:       doc.Title,
					Path:        path,
					Snippet:     snippet,
					Score:       headingScore,
					HeadingID:   heading.id,
					HeadingText: heading.text,
				})
			}
		} else {
			// No matching headings, but document matches - create one result
			snippet := generateSnippet(raw, queryWords, 150)

			// Build path
			path := "/"
			if doc.Key != "" {
				path = "/" + doc.Key
			}

			results = append(results, SearchResult{
				Doc:         doc,
				Key:         doc.Key,
				Title:       doc.Title,
				Path:        path,
				Snippet:     snippet,
				Score:       baseScore,
				HeadingID:   "",
				HeadingText: "",
			})
		}
	}

	// Sort by score (descending)
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		// If scores are equal, sort by title
		return strings.ToLower(results[i].Title) < strings.ToLower(results[j].Title)
	})

	// Limit results
	if maxResults > 0 && len(results) > maxResults {
		results = results[:maxResults]
	}

	return results
}

// findAllMatchingHeadings finds all headings that contain any of the query words.
func findAllMatchingHeadings(md goldmark.Markdown, markdownSource string, queryWords []string) []headingMatch {
	src := []byte(markdownSource)
	ctx := parser.NewContext()
	reader := text.NewReader(src)
	doc := md.Parser().Parse(reader, parser.WithContext(ctx))

	var matches []headingMatch

	_ = gmast.Walk(doc, func(n gmast.Node, entering bool) (gmast.WalkStatus, error) {
		if !entering {
			return gmast.WalkContinue, nil
		}
		h, ok := n.(*gmast.Heading)
		if !ok {
			return gmast.WalkContinue, nil
		}

		// Get heading ID
		idVal, _ := h.AttributeString("id")
		var idStr string
		switch v := idVal.(type) {
		case []byte:
			idStr = strings.TrimSpace(string(v))
		case string:
			idStr = strings.TrimSpace(v)
		}
		if idStr == "" {
			return gmast.WalkContinue, nil
		}

		// Get heading text
		headingText := strings.ToLower(markdown.ExtractNodeText(h, src))
		headingTextOriginal := strings.TrimSpace(markdown.ExtractNodeText(h, src))

		// Check if any query word matches
		for _, word := range queryWords {
			if strings.Contains(headingText, word) {
				matches = append(matches, headingMatch{
					id:   idStr,
					text: headingTextOriginal,
					pos:  int(h.Lines().At(0).Start),
				})
				break // Found a match, no need to check other words for this heading
			}
		}

		return gmast.WalkContinue, nil
	})

	return matches
}

// stripMarkdownSyntax removes markdown syntax and HTML from text to make it plain text.
func stripMarkdownSyntax(text string) string {
	// First, remove all HTML tags and their attributes (including style attributes)
	// This handles cases like <div style="display: flex;...">content</div>
	htmlTagRE := regexp.MustCompile("(?s)<[^>]*>")
	text = htmlTagRE.ReplaceAllString(text, " ")

	// Remove code blocks (```language ... ```)
	codeBlockRE := regexp.MustCompile("(?s)```[a-z]*\\n.*?```")
	text = codeBlockRE.ReplaceAllString(text, "")

	// Remove inline code (`code`)
	inlineCodeRE := regexp.MustCompile("`[^`]+`")
	text = inlineCodeRE.ReplaceAllString(text, "")

	// Remove headings (##, ###, etc.)
	headingRE := regexp.MustCompile(`(?m)^#{1,6}\s+`)
	text = headingRE.ReplaceAllString(text, "")

	// Remove bold/italic (**text**, *text*)
	boldRE := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = boldRE.ReplaceAllString(text, "$1")
	italicRE := regexp.MustCompile(`\*([^*]+)\*`)
	text = italicRE.ReplaceAllString(text, "$1")
	boldUnderscoreRE := regexp.MustCompile("__([^_]+)__")
	text = boldUnderscoreRE.ReplaceAllString(text, "$1")
	italicUnderscoreRE := regexp.MustCompile("_([^_]+)_")
	text = italicUnderscoreRE.ReplaceAllString(text, "$1")

	// Remove links [text](url) -> text
	linkRE := regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
	text = linkRE.ReplaceAllString(text, "$1")

	// Remove images ![alt](url)
	imageRE := regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`)
	text = imageRE.ReplaceAllString(text, "$1")

	// Remove horizontal rules (---, ***)
	hrRE := regexp.MustCompile(`(?m)^[-*]{3,}\s*$`)
	text = hrRE.ReplaceAllString(text, "")

	// Remove blockquotes (> text)
	blockquoteRE := regexp.MustCompile(`(?m)^>\s+`)
	text = blockquoteRE.ReplaceAllString(text, "")

	// Remove list markers (-, *, +, 1.)
	listRE := regexp.MustCompile(`(?m)^\s*[-*+]\s+`)
	text = listRE.ReplaceAllString(text, "")
	numberedListRE := regexp.MustCompile(`(?m)^\s*\d+\.\s+`)
	text = numberedListRE.ReplaceAllString(text, "")

	// Clean up multiple spaces/newlines
	newlineRE := regexp.MustCompile(`\n{3,}`)
	text = newlineRE.ReplaceAllString(text, "\n\n")
	spaceRE := regexp.MustCompile("[ \t]+")
	text = spaceRE.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// generateSnippetAroundHeading creates a snippet around a specific heading position.
func generateSnippetAroundHeading(content string, headingPos int, queryWords []string, maxLen int) string {
	if len(content) == 0 {
		return ""
	}

	// Extract content around the heading
	start := headingPos - maxLen/3
	if start < 0 {
		start = 0
	} else {
		// Try to start at word boundary
		for start > 0 && start < len(content) && !strings.ContainsRune(" \n\t", rune(content[start])) {
			start++
		}
		if start > headingPos {
			start = headingPos - maxLen/3
			if start < 0 {
				start = 0
			}
		}
	}

	end := headingPos + maxLen*2/3
	if end > len(content) {
		end = len(content)
	} else {
		// Try to end at word boundary
		originalEnd := end
		for end < len(content) && !strings.ContainsRune(" \n\t", rune(content[end])) {
			end++
		}
		if end-originalEnd > 20 {
			end = originalEnd
		}
	}

	snippet := content[start:end]
	snippet = strings.TrimSpace(snippet)

	// Strip markdown syntax
	snippet = stripMarkdownSyntax(snippet)

	// Add ellipsis
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	// Truncate if still too long
	if len(snippet) > maxLen+10 {
		cutPos := maxLen
		for cutPos > 0 && cutPos < len(snippet) && !strings.ContainsRune(" \n\t", rune(snippet[cutPos])) {
			cutPos--
		}
		if cutPos > maxLen*3/4 {
			snippet = snippet[:cutPos] + "..."
		} else {
			snippet = snippet[:maxLen] + "..."
		}
	}

	return snippet
}

// generateSnippet creates a text excerpt from content that includes query terms.
func generateSnippet(content string, queryWords []string, maxLen int) string {
	if len(content) == 0 {
		return ""
	}

	contentLower := strings.ToLower(content)

	// Find the first occurrence of any query word
	bestPos := -1
	bestWord := ""
	for _, word := range queryWords {
		pos := strings.Index(contentLower, word)
		if pos >= 0 && (bestPos < 0 || pos < bestPos) {
			bestPos = pos
			bestWord = word
		}
	}

	if bestPos < 0 {
		// No match found, return beginning of content
		snippet := strings.TrimSpace(content)
		// Strip markdown syntax
		snippet = stripMarkdownSyntax(snippet)
		if len(snippet) > maxLen {
			// Try to cut at word boundary
			cutPos := maxLen
			for cutPos > 0 && cutPos < len(snippet) && !strings.ContainsRune(" \n\t", rune(snippet[cutPos])) {
				cutPos--
			}
			if cutPos > maxLen*3/4 {
				snippet = snippet[:cutPos] + "..."
			} else {
				snippet = snippet[:maxLen] + "..."
			}
		}
		return snippet
	}

	// Extract snippet around the match
	start := bestPos - maxLen/3
	if start < 0 {
		start = 0
	} else {
		// Try to start at word boundary
		for start > 0 && start < len(content) && !strings.ContainsRune(" \n\t", rune(content[start])) {
			start++
		}
		if start > bestPos {
			start = bestPos - maxLen/3
			if start < 0 {
				start = 0
			}
		}
	}

	end := bestPos + len(bestWord) + maxLen*2/3
	if end > len(content) {
		end = len(content)
	} else {
		// Try to end at word boundary
		originalEnd := end
		for end < len(content) && !strings.ContainsRune(" \n\t", rune(content[end])) {
			end++
		}
		// If we went too far, revert
		if end-originalEnd > 20 {
			end = originalEnd
		}
	}

	snippet := content[start:end]
	snippet = strings.TrimSpace(snippet)

	// Strip markdown syntax
	snippet = stripMarkdownSyntax(snippet)

	// Add ellipsis
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	// Truncate if still too long (shouldn't happen often, but be safe)
	if len(snippet) > maxLen+10 {
		cutPos := maxLen
		for cutPos > 0 && cutPos < len(snippet) && !strings.ContainsRune(" \n\t", rune(snippet[cutPos])) {
			cutPos--
		}
		if cutPos > maxLen*3/4 {
			snippet = snippet[:cutPos] + "..."
		} else {
			snippet = snippet[:maxLen] + "..."
		}
	}

	return snippet
}
