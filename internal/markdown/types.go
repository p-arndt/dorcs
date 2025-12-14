package markdown

// FrontMatter represents supported metadata at the top of the markdown file.
// YAML front matter is expected by default but the parser supports multiple formats.
type FrontMatter struct {
	Title       string   `yaml:"title" toml:"title" json:"title"`
	Description string   `yaml:"description" toml:"description" json:"description"`
	Date        string   `yaml:"date" toml:"date" json:"date"`
	Tags        []string `yaml:"tags" toml:"tags" json:"tags"`
	Draft       bool     `yaml:"draft" toml:"draft" json:"draft"`
	Order       int      `yaml:"order" toml:"order" json:"order"`    // Order for sorting (lower numbers appear first)
	Author      string   `yaml:"author" toml:"author" json:"author"` // Author name
	After       string   `yaml:"after" toml:"after" json:"after"`    // Key of item this should appear after (use "index" for index.md)
}
