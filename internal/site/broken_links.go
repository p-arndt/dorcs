package site

import (
	"fmt"
	"os"
	"strings"

	"dorcs-v2/internal/markdown"
)

// BrokenLink represents a broken link found in a document.
type BrokenLink struct {
	DocPath   string // Path to the document containing the broken link
	DocKey    string // Key of the document containing the broken link
	LinkText  string // Text of the link
	LinkHref  string // Href of the link (as written)
	TargetKey string // Resolved target document key (what was expected)
	Line      int    // Line number (1-indexed)
	Column    int    // Column number (1-indexed)
}

// CheckBrokenLinks checks all documents in the site for broken internal links.
// It returns a list of broken links found.
func (s *Site) CheckBrokenLinks() []BrokenLink {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var brokenLinks []BrokenLink

	// Iterate through all documents
	for _, doc := range s.index {
		// Read the markdown file
		content, err := os.ReadFile(doc.FilePath)
		if err != nil {
			// Skip files that can't be read
			continue
		}

		originalContent := string(content)

		// Calculate front matter line offset before stripping
		// Front matter format: ---\n...content...\n---\n
		// StripYAMLFrontMatter removes everything up to and including the closing ---
		frontMatterLineOffset := 0
		originalLines := strings.Split(originalContent, "\n")
		if len(originalLines) > 0 && strings.TrimSpace(originalLines[0]) == "---" {
			// Find the closing --- delimiter
			for i := 1; i < len(originalLines); i++ {
				if strings.TrimSpace(originalLines[i]) == "---" {
					// Front matter ends at line i (0-indexed)
					// StripYAMLFrontMatter removes lines 0 through i (inclusive), which is i+1 lines
					// Content in stripped content starts at line 0 (0-indexed), which is line i+1 (0-indexed) in original
					// So offset is i+1. But since ExtractLinksWithLineNumbers returns 1-indexed line numbers,
					// and we need to convert to original file line numbers, we add the offset.
					// However, if the link is at line L in stripped (1-indexed), it's at line L-1 (0-indexed) in stripped,
					// which is at line (L-1)+(i+1) = L+i (0-indexed) in original = line L+i+1 (1-indexed).
					// So we need offset = i+1. But user reports it's one too much, so let's use i.
					frontMatterLineOffset = i
					break
				}
			}
		}

		// Strip front matter to get just the markdown content
		mdContent := markdown.StripYAMLFrontMatter(originalContent)

		// Extract all links with line numbers
		links := markdown.ExtractLinksWithLineNumbers(mdContent)

		// Determine the document's directory key for resolving relative links
		docDir := doc.DirKey
		if isIndexRel(doc.RelPath) {
			// This is an index.md file, so the document's directory is its Key
			docDir = doc.Key
		}

		// Check each link
		for _, link := range links {
			// Resolve the link to a document key
			targetKey, shouldCheck := markdown.ResolveLinkToDocKey(link.Href, docDir)
			if !shouldCheck {
				// This is not a doc link (external URL, anchor, etc.), skip it
				continue
			}

			// Check if the original href contains suspicious path traversal
			// Count how many ../ segments would be needed to go outside docs directory
			hrefClean := strings.ReplaceAll(link.Href, "\\", "/")
			parentTraversals := strings.Count(hrefClean, "../")
			// Calculate how many directory levels deep we are
			dirLevels := 0
			if docDir != "" {
				dirLevels = strings.Count(docDir, "/") + 1
			}
			// If traversals exceed directory levels, it goes outside docs root
			suspiciousPath := parentTraversals > dirLevels

			// Check if the target document exists
			_, exists := s.index[targetKey]
			if !exists || suspiciousPath {
				// Adjust line number to account for front matter
				adjustedLine := link.Line + frontMatterLineOffset

				// If it exists but has suspicious path, still report as warning
				if exists && suspiciousPath {
					// Treat as warning (path goes outside docs directory but resolves correctly)
					brokenLinks = append(brokenLinks, BrokenLink{
						DocPath:   doc.RelPath,
						DocKey:    doc.Key,
						LinkText:  link.Text,
						LinkHref:  link.Href,
						TargetKey: targetKey + " (warning: path goes outside docs directory)",
						Line:      adjustedLine,
						Column:    link.Column,
					})
				} else if !exists {
					// Actually broken link
					brokenLinks = append(brokenLinks, BrokenLink{
						DocPath:   doc.RelPath,
						DocKey:    doc.Key,
						LinkText:  link.Text,
						LinkHref:  link.Href,
						TargetKey: targetKey,
						Line:      adjustedLine,
						Column:    link.Column,
					})
				}
			}
		}
	}

	return brokenLinks
}

// ReportBrokenLinks prints broken links to the console in a formatted way.
func (s *Site) ReportBrokenLinks() {
	brokenLinks := s.CheckBrokenLinks()

	if len(brokenLinks) == 0 {
		return
	}

	fmt.Printf("\n⚠️  Found %d broken link(s):\n\n", len(brokenLinks))

	for _, bl := range brokenLinks {
		// Format: doc.md:line:column: broken link [text](href) -> target-key
		targetDisplay := bl.TargetKey
		if targetDisplay == "" {
			targetDisplay = "(root index)"
		}
		fmt.Printf("  %s:%d:%d: broken link [%s](%s) -> %s\n",
			bl.DocPath, bl.Line, bl.Column, bl.LinkText, bl.LinkHref, targetDisplay)
	}

	fmt.Println()
}
