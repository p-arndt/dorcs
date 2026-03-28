package markdown

import (
	"strings"
	"testing"
)

func TestConvertVideoEmbedsInHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		absent   []string
	}{
		{
			name:  "YouTube watch URL",
			input: `<p>{video:https://www.youtube.com/watch?v=dQw4w9WgXcQ}</p>`,
			contains: []string{
				`class="video-embed"`,
				`src="https://www.youtube.com/embed/dQw4w9WgXcQ"`,
				"allowfullscreen",
			},
			absent: []string{"{video:"},
		},
		{
			name:  "YouTube short URL",
			input: `<p>{video:https://youtu.be/dQw4w9WgXcQ}</p>`,
			contains: []string{
				`src="https://www.youtube.com/embed/dQw4w9WgXcQ"`,
			},
		},
		{
			name:  "YouTube embed URL passthrough",
			input: `<p>{video:https://www.youtube.com/embed/dQw4w9WgXcQ}</p>`,
			contains: []string{
				`src="https://www.youtube.com/embed/dQw4w9WgXcQ"`,
			},
		},
		{
			name:  "Vimeo URL",
			input: `<p>{video:https://vimeo.com/123456789}</p>`,
			contains: []string{
				`src="https://player.vimeo.com/video/123456789"`,
			},
		},
		{
			name:  "generic URL",
			input: `<p>{video:https://example.com/video.mp4}</p>`,
			contains: []string{
				`src="https://example.com/video.mp4"`,
				`class="video-embed"`,
			},
		},
		{
			name:     "no video markers - unchanged",
			input:    "<p>Hello world</p>",
			contains: []string{"<p>Hello world</p>"},
			absent:   []string{"video-embed"},
		},
		{
			name:  "multiple videos",
			input: `<p>{video:https://youtu.be/abc12345678}</p><p>{video:https://vimeo.com/999}</p>`,
			contains: []string{
				"youtube.com/embed/abc12345678",
				"player.vimeo.com/video/999",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertVideoEmbedsInHTML(tt.input)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("result should contain %q, got:\n%s", want, result)
				}
			}
			for _, absent := range tt.absent {
				if strings.Contains(result, absent) {
					t.Errorf("result should not contain %q, got:\n%s", absent, result)
				}
			}
		})
	}
}

func TestExtractYouTubeID(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"https://youtu.be/dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"https://www.youtube.com/embed/dQw4w9WgXcQ", "dQw4w9WgXcQ"},
		{"https://youtube.com/watch?v=abc-123_456&t=10", "abc-123_456"},
		{"https://example.com/page", ""},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := extractYouTubeID(tt.url)
			if got != tt.want {
				t.Errorf("extractYouTubeID(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestExtractVimeoID(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://vimeo.com/123456789", "123456789"},
		{"https://www.vimeo.com/999", "999"},
		{"https://example.com/page", ""},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := extractVimeoID(tt.url)
			if got != tt.want {
				t.Errorf("extractVimeoID(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}
