package markdown

import (
	"path"
	"regexp"
	"strings"
)

var mdInlineLinkRE = regexp.MustCompile(`\[(?P<text>[^\]]+)\]\((?P<href>[^)]+)\)`)

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
func RewriteExtensionlessDocLinks(md string, currentDirKey string, basePath string) string {
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
		// Skip anchors
		if strings.HasPrefix(href, "#") {
			return m
		}
		// Skip root-absolute links
		if strings.HasPrefix(href, "/") {
			return m
		}
		// Skip query-ish links (keep conservative)
		if strings.Contains(href, "?") || strings.Contains(href, "&") {
			return m
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

		// Handle index.md references: folder/index -> /folder, index -> /
		if clean == "index" {
			newHref := "/"
			if basePath != "" {
				newHref = basePath + "/"
			}
			return "[" + text + "](" + newHref + ")"
		}
		if strings.HasSuffix(clean, "/index") {
			clean = strings.TrimSuffix(clean, "/index")
			if clean == "" {
				newHref := "/"
				if basePath != "" {
					newHref = basePath + "/"
				}
				return "[" + text + "](" + newHref + ")"
			}
			newHref := "/" + clean
			if basePath != "" {
				newHref = basePath + newHref
			}
			return "[" + text + "](" + newHref + ")"
		}

		newHref := "/" + clean
		if basePath != "" {
			newHref = basePath + newHref
		}
		return "[" + text + "](" + newHref + ")"
	})
}
