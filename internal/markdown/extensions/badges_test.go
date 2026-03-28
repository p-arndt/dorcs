package markdown

import (
	"strings"
	"testing"
)

func TestConvertBadgesInHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		absent   []string
	}{
		{
			name:  "NEW badge",
			input: `<h2>Feature {badge:NEW}</h2>`,
			contains: []string{
				`class="badge badge-new"`,
				">NEW</span>",
			},
			absent: []string{"{badge:NEW}"},
		},
		{
			name:  "BETA badge",
			input: `<p>This API is {badge:BETA} and may change.</p>`,
			contains: []string{
				`class="badge badge-beta"`,
				">BETA</span>",
			},
		},
		{
			name:  "DEPRECATED badge",
			input: `<h3>Old Method {badge:DEPRECATED}</h3>`,
			contains: []string{
				`class="badge badge-deprecated"`,
			},
		},
		{
			name:  "EXPERIMENTAL badge",
			input: `<p>{badge:EXPERIMENTAL}</p>`,
			contains: []string{
				`class="badge badge-experimental"`,
			},
		},
		{
			name:  "case insensitive",
			input: `<p>{badge:beta}</p>`,
			contains: []string{
				`class="badge badge-beta"`,
			},
		},
		{
			name:  "custom badge",
			input: `<p>{badge:PREVIEW}</p>`,
			contains: []string{
				`class="badge badge-new"`,
				">PREVIEW</span>",
			},
		},
		{
			name:     "no badges - unchanged",
			input:    "<p>Hello world</p>",
			contains: []string{"<p>Hello world</p>"},
			absent:   []string{"badge"},
		},
		{
			name:  "multiple badges",
			input: `<p>{badge:NEW} {badge:BETA}</p>`,
			contains: []string{
				"badge-new",
				"badge-beta",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertBadgesInHTML(tt.input)
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
