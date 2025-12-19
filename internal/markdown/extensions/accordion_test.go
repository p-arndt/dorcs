package markdown

import "testing"

func TestConvertAccordionBlocksInMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple accordion",
			input: `:::accordion Click to expand
This is the content.
:::`,
			expected: `> **ACCORDION:Click to expand**
> This is the content.
`,
		},
		{
			name: "accordion with multiline content",
			input: `:::accordion Title
First line.
Second line.
Third line.
:::`,
			expected: `> **ACCORDION:Title**
> First line.
> Second line.
> Third line.
`,
		},
		{
			name: "accordion with empty lines",
			input: `:::accordion Title
First paragraph.

Second paragraph.
:::`,
			expected: `> **ACCORDION:Title**
> First paragraph.
>
> Second paragraph.
`,
		},
		{
			name: "accordion with markdown formatting",
			input: `:::accordion Title
This has **bold** and *italic* text.

- List item 1
- List item 2
:::`,
			expected: `> **ACCORDION:Title**
> This has **bold** and *italic* text.
>
> - List item 1
> - List item 2
`,
		},
		{
			name:     "accordion in code block should not be processed",
			input:    "```\n:::accordion Title\nContent\n:::\n```",
			expected: "```\n:::accordion Title\nContent\n:::\n```",
		},
		{
			name: "accordion in indented code block should not be processed",
			input: `    :::accordion Title
    Content
    :::`,
			expected: `    :::accordion Title
    Content
    :::`,
		},
		{
			name: "multiple accordions",
			input: `:::accordion First
Content 1
:::

:::accordion Second
Content 2
:::`,
			expected: `> **ACCORDION:First**
> Content 1


> **ACCORDION:Second**
> Content 2
`,
		},
		{
			name: "accordion without title",
			input: `:::accordion
Content here.
:::`,
			expected: `> **ACCORDION:**
> Content here.
`,
		},
		{
			name: "incomplete accordion (no closing) should not be processed",
			input: `:::accordion Title
Content without closing`,
			expected: `:::accordion Title
Content without closing`,
		},
		{
			name: "accordion with code block inside",
			input: `:::accordion Title
Here's some code:

` + "```go" + `
func main() {}
` + "```" + `
:::`,
			expected: `> **ACCORDION:Title**
> Here's some code:
>
> ` + "```go" + `
> func main() {}
> ` + "```" + `
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertAccordionBlocksInMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertAccordionBlocksInMarkdown() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestConvertAccordionBlocksInHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple accordion in HTML",
			input:    `<blockquote><p><strong>ACCORDION:Click to expand</strong></p><p>This is the content.</p></blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>Click to expand</summary><div class="dorcs-accordion-content"><p>This is the content.</p></div></details>`,
		},
		{
			name:     "accordion with multiple paragraphs",
			input:    `<blockquote><p><strong>ACCORDION:Title</strong></p><p>First paragraph.</p><p>Second paragraph.</p></blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>Title</summary><div class="dorcs-accordion-content"><p>First paragraph.</p><p>Second paragraph.</p></div></details>`,
		},
		{
			name:     "accordion without title",
			input:    `<blockquote><p><strong>ACCORDION:</strong></p><p>Content here.</p></blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>Click to expand</summary><div class="dorcs-accordion-content"><p>Content here.</p></div></details>`,
		},
		{
			name:     "accordion with list content",
			input:    `<blockquote><p><strong>ACCORDION:List</strong></p><ul><li>Item 1</li><li>Item 2</li></ul></blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>List</summary><div class="dorcs-accordion-content"><ul><li>Item 1</li><li>Item 2</li></ul></div></details>`,
		},
		{
			name:     "regular blockquote unchanged",
			input:    `<blockquote><p>Regular quote.</p></blockquote>`,
			expected: `<blockquote><p>Regular quote.</p></blockquote>`,
		},
		{
			name:     "accordion with empty paragraphs",
			input:    `<blockquote><p><strong>ACCORDION:Title</strong></p><p></p><p>Content</p><p></p></blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>Title</summary><div class="dorcs-accordion-content"><p>Content</p></div></details>`,
		},
		{
			name:     "accordion without p tag around strong",
			input:    `<blockquote><strong>ACCORDION:Title</strong><p>Content</p></blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>Title</summary><div class="dorcs-accordion-content"><p>Content</p></div></details>`,
		},
		{
			name:     "accordion with complex content",
			input:    `<blockquote><p><strong>ACCORDION:Complex</strong></p><p>Paragraph with <strong>bold</strong> and <em>italic</em>.</p><h2>Heading</h2><p>More content.</p></blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>Complex</summary><div class="dorcs-accordion-content"><p>Paragraph with <strong>bold</strong> and <em>italic</em>.</p><h2>Heading</h2><p>More content.</p></div></details>`,
		},
		{
			name:     "accordion with leading br tag",
			input:    `<blockquote><p><strong>ACCORDION:Click to expand</strong></p><br>This is the accordion content.<br>It can contain <strong>markdown</strong> formatting.</blockquote>`,
			expected: `<details class="dorcs-accordion"><summary>Click to expand</summary><div class="dorcs-accordion-content">This is the accordion content.<br>It can contain <strong>markdown</strong> formatting.</div></details>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertAccordionBlocksInHTML(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertAccordionBlocksInHTML() = %q; want %q", result, tt.expected)
			}
		})
	}
}
