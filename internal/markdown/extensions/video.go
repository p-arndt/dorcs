package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

var videoPattern = regexp.MustCompile(`\{video:\s*(https?://[^}\s]+)\s*\}`)

// ConvertVideoEmbedsInHTML replaces {video:URL} patterns with responsive iframe embeds.
// Supports YouTube and Vimeo URLs. Other URLs are embedded as generic iframes.
func ConvertVideoEmbedsInHTML(htmlContent string) string {
	if !strings.Contains(htmlContent, "{video:") {
		return htmlContent
	}

	return videoPattern.ReplaceAllStringFunc(htmlContent, func(match string) string {
		url := match[7 : len(match)-1] // strip {video: and }
		url = strings.TrimSpace(url)

		embedURL := toEmbedURL(url)
		if embedURL == "" {
			return match // Unknown format, leave as-is
		}

		return fmt.Sprintf(
			`<div class="video-embed"><iframe src="%s" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen loading="lazy"></iframe></div>`,
			embedURL,
		)
	})
}

// toEmbedURL converts a video URL to its embeddable iframe URL.
func toEmbedURL(url string) string {
	// YouTube: youtube.com/watch?v=ID, youtu.be/ID, youtube.com/embed/ID
	if id := extractYouTubeID(url); id != "" {
		return "https://www.youtube.com/embed/" + id
	}

	// Vimeo: vimeo.com/ID
	if id := extractVimeoID(url); id != "" {
		return "https://player.vimeo.com/video/" + id
	}

	// Generic URL — assume it's directly embeddable
	return url
}

var (
	ytWatchRe = regexp.MustCompile(`(?:youtube\.com/watch\?.*v=|youtube\.com/embed/|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	vimeoRe   = regexp.MustCompile(`vimeo\.com/(\d+)`)
)

func extractYouTubeID(url string) string {
	m := ytWatchRe.FindStringSubmatch(url)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func extractVimeoID(url string) string {
	m := vimeoRe.FindStringSubmatch(url)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}
