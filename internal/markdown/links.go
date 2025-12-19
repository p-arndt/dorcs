package markdown

import (
	"path"
	"regexp"
	"strings"
)

var mdInlineLinkRE = regexp.MustCompile(`\[(?P<text>[^\]]+)\]\((?P<href>[^)]+)\)`)
var mdImageRE = regexp.MustCompile(`!\[(?P<alt>[^\]]*)\]\((?P<src>[^)]+)\)`)

// LinkInfo represents a markdown link with its position.
type LinkInfo struct {
	Text   string // Link text
	Href   string // Link href (as written in markdown)
	Line   int    // Line number (1-indexed)
	Column int    // Column number (1-indexed)
}

// ExtractLinksWithLineNumbers extracts all markdown links from the given markdown content
// and returns them with their line and column positions.
// It ignores links found inside code blocks (fenced or indented) and inline code.
func ExtractLinksWithLineNumbers(md string) []LinkInfo {
	var links []LinkInfo
	lines := strings.Split(md, "\n")

	inFencedCodeBlock := false
	openingFenceChar := "" // "`" or "~"
	openingFenceLength := 0

	// Regex to match inline code (backticks)
	inlineCodeRE := regexp.MustCompile("`[^`]*`")

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for fenced code block start/end (``` or ~~~)
		if strings.HasPrefix(trimmed, "```") {
			if !inFencedCodeBlock {
				// Opening fence - count backticks
				openingFenceChar = "`"
				openingFenceLength = 0
				for j := 0; j < len(trimmed) && trimmed[j] == '`'; j++ {
					openingFenceLength++
				}
				inFencedCodeBlock = true
			} else {
				// Check if this is a closing fence (must have at least as many backticks)
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
			// Skip links in code block markers
			continue
		}

		if strings.HasPrefix(trimmed, "~~~") {
			if !inFencedCodeBlock {
				// Opening fence - count tildes
				openingFenceChar = "~"
				openingFenceLength = 0
				for j := 0; j < len(trimmed) && trimmed[j] == '~'; j++ {
					openingFenceLength++
				}
				inFencedCodeBlock = true
			} else if openingFenceChar == "~" {
				// Check if this is a closing fence (must have at least as many tildes)
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
			// Skip links in code block markers
			continue
		}

		// Skip links inside fenced code blocks
		if inFencedCodeBlock {
			continue
		}

		// Check for indented code blocks (4+ spaces at the start)
		// This is a simple heuristic - indented code blocks in markdown
		isIndentedCodeBlock := len(line) > 0 && strings.HasPrefix(line, "    ")

		// Find all matches in this line
		matches := mdInlineLinkRE.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			if len(match) >= 6 {
				textStart := match[2]
				textEnd := match[3]
				hrefStart := match[4]
				hrefEnd := match[5]

				// Skip if we're in an indented code block
				if isIndentedCodeBlock {
					continue
				}

				// Check if the link is inside inline code (backticks)
				// We need to check if the link position overlaps with any inline code spans
				linkInInlineCode := false
				inlineCodeMatches := inlineCodeRE.FindAllStringIndex(line, -1)
				for _, codeMatch := range inlineCodeMatches {
					codeStart := codeMatch[0]
					codeEnd := codeMatch[1]
					// Check if the link overlaps with this inline code span
					if (hrefStart >= codeStart && hrefStart < codeEnd) ||
						(hrefEnd > codeStart && hrefEnd <= codeEnd) ||
						(hrefStart <= codeStart && hrefEnd >= codeEnd) {
						linkInInlineCode = true
						break
					}
				}

				if linkInInlineCode {
					continue
				}

				text := line[textStart:textEnd]
				href := line[hrefStart:hrefEnd]

				// Check for @ignore comment on the same line or previous line
				// Format: <!-- @ignore --> or <!--@ignore--> or <!-- @ignore-link -->
				shouldIgnore := false

				// Check current line before the link
				beforeLink := line[:hrefStart]
				if strings.Contains(beforeLink, "<!--") {
					// Look for @ignore in HTML comments
					commentRE := regexp.MustCompile(`<!--\s*@ignore(-link)?\s*-->`)
					if commentRE.MatchString(beforeLink) {
						shouldIgnore = true
					}
				}

				// Check previous line for @ignore comment
				if !shouldIgnore && lineNum > 0 {
					prevLine := strings.TrimSpace(lines[lineNum-1])
					commentRE := regexp.MustCompile(`<!--\s*@ignore(-link)?\s*-->`)
					if commentRE.MatchString(prevLine) {
						shouldIgnore = true
					}
				}

				if shouldIgnore {
					continue
				}

				links = append(links, LinkInfo{
					Text:   text,
					Href:   href,
					Line:   lineNum + 1,   // 1-indexed
					Column: hrefStart + 1, // 1-indexed
				})
			}
		}
	}

	return links
}

// ResolveLinkToDocKey resolves a markdown link href to a document key.
// It follows the same logic as RewriteExtensionlessDocLinks but returns the target key instead of rewriting.
// Returns the resolved document key and true if this is a doc link that should be checked, false otherwise.
func ResolveLinkToDocKey(href string, currentDirKey string) (string, bool) {
	href = strings.TrimSpace(href)
	if href == "" {
		return "", false
	}
	lower := strings.ToLower(href)

	// Skip absolute URLs / mailto / etc.
	if strings.Contains(lower, "://") || strings.HasPrefix(lower, "mailto:") {
		return "", false
	}
	// Skip anchors (standalone anchors)
	if strings.HasPrefix(href, "#") {
		return "", false
	}
	// Extract anchor if present (e.g., "file.md#anchor" -> "file.md" and "#anchor")
	if anchorIdx := strings.Index(href, "#"); anchorIdx != -1 {
		href = href[:anchorIdx]
		lower = strings.ToLower(href)
	}
	// Handle root-absolute links (strip leading / and treat as relative to root)
	isRootAbsolute := strings.HasPrefix(href, "/")
	if isRootAbsolute {
		href = strings.TrimPrefix(href, "/")
		currentDirKey = "" // Root-absolute links are always relative to root
		lower = strings.ToLower(href)
	}
	// Skip query-ish links
	if strings.Contains(href, "?") || strings.Contains(href, "&") {
		return "", false
	}

	// Handle .md extension: strip it and continue
	hasMdExt := strings.HasSuffix(lower, ".md")
	if hasMdExt {
		href = href[:len(href)-3] // remove ".md"
	}

	// Skip if it has a non-.md extension (e.g. ".png", ".html")
	if ext := path.Ext(href); ext != "" {
		return "", false
	}

	// Resolve relative paths based on the current document's directory
	hrefClean := strings.ReplaceAll(href, "\\", "/")

	// Build the full path from current directory
	var fullPath string
	if currentDirKey == "" {
		// Document is at root
		fullPath = hrefClean
	} else {
		// Document is in a subdirectory
		fullPath = currentDirKey + "/" + hrefClean
	}
	// Clean to resolve .. and . components
	clean := path.Clean("/" + fullPath)
	clean = strings.TrimPrefix(clean, "/")
	if clean == "" || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
		return "", false
	}

	// Handle index.md references: folder/index -> folder, index -> "" (root)
	if clean == "index" {
		return "", true // root index
	}
	if strings.HasSuffix(clean, "/index") {
		clean = strings.TrimSuffix(clean, "/index")
		if clean == "" {
			return "", true // root index
		}
		return clean, true
	}

	return clean, true
}

// RewriteExtensionlessDocLinks rewrites relative markdown links into extensionless doc routes.
//
// It handles:
// - extensionless links: [Getting Started](getting-started) -> [Getting Started](/getting-started) or [Getting Started](/basepath/getting-started) if basepath is set
// - .md links: [Getting Started](getting-started.md) -> [Getting Started](/getting-started) or [Getting Started](/basepath/getting-started) if basepath is set
// - nested paths: [Guide](guide/intro.md) -> [Guide](/guide/intro) or [Guide](/basepath/guide/intro) if basepath is set
//
// Links are NOT rewritten if they are:
// - absolute URLs (http/https)
// - anchors (#...)
// - root-absolute (/...) - these are left as-is (templates will handle basepath)
// - other file extensions (e.g. ".png", ".html", ".pdf")
// - query strings (?...)
//
// basePath is the URL path prefix (e.g., "/docs"). Empty string means no prefix.
// language is the language code (e.g., "de"). Empty string means default language (no prefix).
// version is the version identifier (e.g., "v1"). Empty string means default version (no prefix).
// defaultVersion is the default version identifier. If version equals defaultVersion, no version prefix is added.
// defaultLanguage is the default language code. If language equals defaultLanguage, no language prefix is added.
func RewriteExtensionlessDocLinks(md string, currentDirKey string, basePath string, language string, version string, defaultVersion string, defaultLanguage string) string {
	return mdInlineLinkRE.ReplaceAllStringFunc(md, func(m string) string {
		sub := mdInlineLinkRE.FindStringSubmatch(m)
		if len(sub) != 3 {
			return m
		}
		text := sub[1]
		href := strings.TrimSpace(sub[2])

		if href == "" {
			return m
		}
		lower := strings.ToLower(href)

		// Skip absolute URLs / mailto / etc.
		if strings.Contains(lower, "://") || strings.HasPrefix(lower, "mailto:") {
			return m
		}
		// Skip anchors (standalone anchors)
		if strings.HasPrefix(href, "#") {
			return m
		}
		// Handle root-absolute links (strip .md extension and anchor if present)
		if strings.HasPrefix(href, "/") {
			// Extract anchor if present
			var anchor string
			if anchorIdx := strings.Index(href, "#"); anchorIdx != -1 {
				anchor = href[anchorIdx:]
				href = href[:anchorIdx]
			}
			// Strip .md extension if present
			lower := strings.ToLower(href)
			if strings.HasSuffix(lower, ".md") {
				href = href[:len(href)-3]
			}
			return "[" + text + "](" + href + anchor + ")"
		}
		// Skip query-ish links (keep conservative)
		if strings.Contains(href, "?") || strings.Contains(href, "&") {
			return m
		}

		// Extract anchor if present (e.g., "file.md#anchor" -> "file.md" and "#anchor")
		var anchor string
		if anchorIdx := strings.Index(href, "#"); anchorIdx != -1 {
			anchor = href[anchorIdx:]
			href = href[:anchorIdx]
			lower = strings.ToLower(href)
		}

		// Handle .md extension: strip it and continue
		hasMdExt := strings.HasSuffix(lower, ".md")
		if hasMdExt {
			href = href[:len(href)-3] // remove ".md"
		}

		// Skip if it has a non-.md extension (e.g. ".png", ".html")
		if ext := path.Ext(href); ext != "" {
			return m
		}

		// Resolve relative paths based on the current document's directory
		// All paths without a leading / are treated as relative
		hrefClean := strings.ReplaceAll(href, "\\", "/")

		// Build the full path from current directory
		var fullPath string
		if currentDirKey == "" {
			// Document is at root
			fullPath = hrefClean
		} else {
			// Document is in a subdirectory
			fullPath = currentDirKey + "/" + hrefClean
		}
		// Clean to resolve .. and . components
		clean := path.Clean("/" + fullPath)
		clean = strings.TrimPrefix(clean, "/")
		if clean == "" || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
			return m
		}

		// Build version and language prefix (language-first structure)
		// Only add language prefix if language is set and not the default
		var prefixBuilder strings.Builder
		if language != "" && language != defaultLanguage {
			prefixBuilder.WriteByte('/')
			prefixBuilder.WriteString(language)
		}
		// Only add version prefix if version is set and not the default
		if version != "" && version != defaultVersion {
			prefixBuilder.WriteByte('/')
			prefixBuilder.WriteString(version)
		}
		prefix := prefixBuilder.String()

		// Handle index.md references: folder/index -> /folder, index -> /
		if clean == "index" {
			newHref := prefix + "/"
			if basePath != "" {
				newHref = basePath + newHref
			}
			return "[" + text + "](" + newHref + anchor + ")"
		}
		if strings.HasSuffix(clean, "/index") {
			clean = strings.TrimSuffix(clean, "/index")
			if clean == "" {
				newHref := prefix + "/"
				if basePath != "" {
					newHref = basePath + newHref
				}
				return "[" + text + "](" + newHref + anchor + ")"
			}
			newHref := prefix + "/" + clean
			if basePath != "" {
				newHref = basePath + newHref
			}
			return "[" + text + "](" + newHref + anchor + ")"
		}

		newHref := prefix + "/" + clean
		if basePath != "" {
			newHref = basePath + newHref
		}
		return "[" + text + "](" + newHref + anchor + ")"
	})
}

// RewriteRelativeImagePaths rewrites relative image paths to absolute paths.
// This ensures images work correctly when default language uses its folder but URL has no prefix.
// For example: ./logo.png in docs/en/index.md (served at /) -> /logo.png
func RewriteRelativeImagePaths(md string, currentDirKey string, basePath string, language string, version string, defaultVersion string, defaultLanguage string) string {
	return mdImageRE.ReplaceAllStringFunc(md, func(m string) string {
		sub := mdImageRE.FindStringSubmatch(m)
		if len(sub) != 3 {
			return m
		}
		alt := sub[1]
		src := strings.TrimSpace(sub[2])

		if src == "" {
			return m
		}
		lower := strings.ToLower(src)

		// Skip absolute URLs
		if strings.Contains(lower, "://") {
			return m
		}
		// Skip if already absolute (starts with /)
		if strings.HasPrefix(src, "/") {
			// Add basePath if needed
			if basePath != "" && !strings.HasPrefix(src, basePath) {
				return "![" + alt + "](" + basePath + src + ")"
			}
			return m
		}

		// Resolve relative path
		srcClean := strings.ReplaceAll(src, "\\", "/")
		var fullPath string
		if currentDirKey == "" {
			fullPath = srcClean
		} else {
			fullPath = currentDirKey + "/" + srcClean
		}
		// Clean to resolve .. and . components
		clean := path.Clean("/" + fullPath)
		clean = strings.TrimPrefix(clean, "/")
		if clean == "" || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
			return m
		}

		// Build absolute path (no language/version prefix for images - they're served from site root)
		newSrc := "/" + clean
		if basePath != "" {
			newSrc = basePath + newSrc
		}
		return "![" + alt + "](" + newSrc + ")"
	})
}
