package markdown

import (
	"reflect"
	"testing"
)

func TestParseSlideDirectives(t *testing.T) {
	tests := []struct {
		name        string
		chunk       string
		inherited   SlideMetadata
		wantMeta    SlideMetadata
		wantContent string
	}{
		{
			name:        "no directives",
			chunk:       "# Title\nContent",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{},
			wantContent: "# Title\nContent",
		},
		{
			name:        "class directive",
			chunk:       "<!-- _class: lead -->\n\n# Title",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Class: "lead", Layout: "lead"},
			wantContent: "# Title",
		},
		{
			name:        "multiple directives",
			chunk:       "<!-- _class: lead -->\n<!-- _backgroundColor: #1a1a2e -->\n\n# Title",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Class: "lead", Layout: "lead", BackgroundColor: "#1a1a2e"},
			wantContent: "# Title",
		},
		{
			name:        "class with multiple values",
			chunk:       "<!-- _class: lead left -->\n\nContent",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Class: "lead left"},
			wantContent: "Content",
		},
		{
			name:        "footer directive",
			chunk:       "<!-- _footer: Confidential -->\n\nSlide text",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Footer: "Confidential"},
			wantContent: "Slide text",
		},
		{
			name:        "inheritance",
			chunk:       "<!-- backgroundColor: aqua -->\n\nPage 1",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{BackgroundColor: "aqua"},
			wantContent: "Page 1",
		},
		{
			name:        "inherited value applied",
			chunk:       "# Slide 2",
			inherited:   SlideMetadata{BackgroundColor: "aqua", Header: "My Talk"},
			wantMeta:    SlideMetadata{BackgroundColor: "aqua", Header: "My Talk"},
			wantContent: "# Slide 2",
		},
		{
			name:        "color and backgroundImage",
			chunk:       "<!-- _color: white -->\n<!-- _backgroundImage: url(bg.jpg) -->\n\nDark slide",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Color: "white", BackgroundImage: "url(bg.jpg)"},
			wantContent: "Dark slide",
		},
		{
			name:        "layout primitives",
			chunk:       "<!-- _layout: columns-2 -->\n<!-- _gap: loose -->\n<!-- _align: start -->\n\nContent",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Layout: "columns-2", Gap: "loose", Align: "start"},
			wantContent: "Content",
		},
		{
			name:        "layout inheritance",
			chunk:       "# Slide 2",
			inherited:   SlideMetadata{Layout: "lead", Gap: "tight"},
			wantMeta:    SlideMetadata{Layout: "lead", Gap: "tight"},
			wantContent: "# Slide 2",
		},
		{
			name:        "space-separated layout with extra class",
			chunk:       "<!-- _layout: lead invert -->\n\n# Title",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Layout: "lead", Class: "invert"},
			wantContent: "# Title",
		},
		{
			name:        "_class promotes to layout when no layout set",
			chunk:       "<!-- _class: lead -->\n\n# Title",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Class: "lead", Layout: "lead"},
			wantContent: "# Title",
		},
		{
			name:        "_class does not promote when layout already set",
			chunk:       "<!-- _layout: big -->\n<!-- _class: invert -->\n\n# Title",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Layout: "big", Class: "invert"},
			wantContent: "# Title",
		},
		{
			name:        "two-columns alias normalizes to cols",
			chunk:       "<!-- _layout: two-columns -->\n\nContent",
			inherited:   SlideMetadata{},
			wantMeta:    SlideMetadata{Layout: "cols"},
			wantContent: "Content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMeta, gotContent, _ := ParseSlideDirectives(tt.chunk, tt.inherited)
			if !reflect.DeepEqual(gotMeta, tt.wantMeta) {
				t.Errorf("ParseSlideDirectives() meta = %+v, want %+v", gotMeta, tt.wantMeta)
			}
			if gotContent != tt.wantContent {
				t.Errorf("ParseSlideDirectives() content = %q, want %q", gotContent, tt.wantContent)
			}
		})
	}
}
