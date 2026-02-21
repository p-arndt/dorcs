package markdown

import (
	"reflect"
	"testing"
)

func TestSplitSlides(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []string
	}{
		{
			name: "single slide",
			body: "# Hello\nWorld",
			want: []string{"# Hello\nWorld"},
		},
		{
			name: "two slides",
			body: "# Slide 1\n\nContent\n\n---\n\n# Slide 2\nMore content",
			want: []string{"# Slide 1\n\nContent", "# Slide 2\nMore content"},
		},
		{
			name: "three slides",
			body: "A\n\n---\n\nB\n\n---\n\nC",
			want: []string{"A", "B", "C"},
		},
		{
			name: "empty body",
			body: "",
			want: nil,
		},
		{
			name: "whitespace only",
			body: "   \n\n  ",
			want: nil,
		},
		{
			name: "underline separator",
			body: "One\n\n___\n\nTwo",
			want: []string{"One", "Two"},
		},
		{
			name: "asterisk separator",
			body: "X\n\n***\n\nY",
			want: []string{"X", "Y"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitSlides(tt.body)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitSlides() = %v, want %v", got, tt.want)
			}
		})
	}
}
