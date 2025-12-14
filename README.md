# dorcs

A single-binary static documentation server for Markdown files with YAML front matter support.

## Features

- **Single binary** – no runtime dependencies, statically linkable
- **Extensionless URLs** – `/guide/getting-started` serves `docs/guide/getting-started.md`
- **YAML front matter** – extract title, description, date, tags, and draft status
- **Table of Contents** – auto-generated from H2/H3/H4 headings with scrollspy
- **Navigation sidebar** – built from your docs folder structure
- **Responsive design** – mobile-friendly with collapsible sidebar
- **Dark mode** – automatic based on system preference

## Quick Start

```sh
# Run from source
go run ./cmd/dorcs --dir ./docs --addr 127.0.0.1:8080

# Build and run
go build -o dorcs ./cmd/dorcs
./dorcs --dir ./docs

# Run with auto-reload and live browser refresh (for development)
./dorcs --dir ./docs --watch
```

Open http://localhost:8080 in your browser.

## Building

### Standard Build

```sh
go build -o dorcs ./cmd/dorcs
```

### Static Linux Binary (for containers/deployment)

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o dorcs \
  ./cmd/dorcs
```

### Windows

```sh
GOOS=windows GOARCH=amd64 go build -o dorcs.exe ./cmd/dorcs
```

## Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-dir` | `./docs` | Directory containing markdown documents |
| `-addr` | `:8080` | Listen address (e.g., `:8080`, `127.0.0.1:8080`) |
| `-base-url` | `""` | URL path prefix (e.g., `/docs`) |
| `-title` | `""` | Site title shown in header (overrides config file) |
| `-cache` | `true` | Cache rendered documents in memory |
| `-no-drafts` | `true` | Hide documents with `draft: true` front matter |
| `-config` | `""` | Path to config file (default: looks for `dorcs.yaml` in docs dir) |
| `-theme` | `""` | Theme preset: `default`, `ocean`, `forest`, `sunset`, `midnight`, `lavender`, `rose` |
| `-theme-mode` | `""` | Theme mode: `light`, `dark`, `auto` |
| `-watch` | `false` | Watch for file changes and automatically reload |

## Folder Structure

```
docs/
├── index.md              # Home page (served at /)
├── getting-started.md    # Served at /getting-started
├── guide/
│   ├── index.md          # Served at /guide
│   ├── installation.md   # Served at /guide/installation
│   └── configuration.md  # Served at /guide/configuration
└── api/
    ├── index.md          # Served at /api
    └── endpoints.md      # Served at /api/endpoints
```

### URL Routing Rules

- `docs/index.md` → `/`
- `docs/foo.md` → `/foo`
- `docs/guide/index.md` → `/guide`
- `docs/guide/intro.md` → `/guide/intro`

## Front Matter

Each markdown file can include YAML front matter:

```yaml
---
title: Getting Started
description: Learn how to set up and configure dorcs
date: 2024-01-15
tags:
  - tutorial
  - beginner
draft: false
---

# Your content starts here...
```

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Page title (used in nav, browser tab) |
| `description` | string | Meta description |
| `date` | string | Publication date (YYYY-MM-DD) |
| `tags` | list | Tags for categorization |
| `draft` | bool | If true, hidden when `-no-drafts` is set |

## Configuration File

dorcs supports a configuration file for advanced customization. Place a `dorcs.yaml` (or `dorcs.yml` or `dorcs.json`) file in your docs directory.

### Example Configuration

```yaml
# Site metadata
site:
  title: "My Documentation"
  description: "Documentation for my awesome project"
  # logo: "/static/logo.png"  # Optional: replaces text title
  # favicon: "/static/custom-favicon.ico"

# Theme and styling
theme:
  mode: auto          # light, dark, or auto (follows system)
  preset: default     # default, ocean, forest, sunset, midnight, lavender, rose
  # Note: Code syntax highlighting theme is automatically determined by the preset
  # custom_css: "custom.css"  # Path to custom CSS file
  # font_family: '"Inter", system-ui, sans-serif'
  # mono_font_family: '"JetBrains Mono", monospace'

# Navigation
nav:
  show_search: true
  expand_all: false
  links:
    - title: "GitHub"
      url: "https://github.com/yourusername/yourproject"
      external: true
      icon: "github"  # github, twitter, discord, external

# Footer
footer:
  copyright: "© 2024 Your Name"
  show_powered_by: true
  # links:
  #   - title: "Privacy"
  #     url: "/privacy"
```

### Theme Presets

| Preset | Description |
|--------|-------------|
| `default` | GitHub-inspired light/dark theme |
| `ocean` | Cool blue tones |
| `forest` | Natural green palette |
| `sunset` | Warm orange/brown tones |
| `midnight` | Deep purple/blue (Catppuccin-inspired) |
| `lavender` | Soft purple aesthetic |
| `rose` | Pink/rose accent colors |

### Custom Colors

Override any preset color by specifying custom values:

```yaml
theme:
  preset: default
  colors:
    light:
      background: "#ffffff"
      foreground: "#1f2328"
      muted: "#57606a"
      border: "#d0d7de"
      accent: "#0969da"
      code_background: "#f6f8fa"
    dark:
      background: "#0d1117"
      foreground: "#e6edf3"
      accent: "#2f81f7"
```

## Project Structure

```
dorcs-v2/
├── cmd/
│   └── dorcs/
│       ├── main.go           # Entry point
│       └── web/
│           ├── templates/    # HTML templates
│           │   ├── layout.html
│           │   ├── doc.html
│           │   └── index.html
│           └── static/
│               └── style.css
├── internal/
│   ├── config/              # Configuration loading and theme presets
│   ├── server/              # HTTP handler and middleware
│   ├── site/                # Markdown indexing and rendering
│   └── templates/           # Template utilities
├── docs/                    # Example documentation
├── go.mod
└── go.sum
```

## Development

```sh
# Run with auto-reload for markdown changes
go run ./cmd/dorcs --dir ./docs --watch

# The server will automatically detect changes to:
# - .md and .markdown files
# - New files and directories
# - File deletions and renames

# Run tests
go test ./...

# Check for issues
go vet ./...
```

### Watch Mode

When running with `--watch`, the server monitors your docs directory for changes:

- **Automatic rebuild**: Index is rebuilt when markdown files change
- **Smart content reload**: Updates page content without full refresh
- **Preserves scroll position**: Stay exactly where you were reading
- **Maintains sidebar state**: Keeps folders expanded/collapsed as you left them
- **No restart needed**: Server keeps running while you edit
- **Debounced (500ms)**: Multiple rapid changes are batched together - perfect for auto-save editors
- **Instant feedback**: See your changes ~500ms after saving

The live reload feature uses Server-Sent Events (SSE) to push reload notifications to connected browsers. When you save a markdown file, the browser intelligently updates just the content area while preserving your navigation state and scroll position. If the smart reload fails for any reason, it automatically falls back to a full page reload.

**Perfect for active editing**: The 500ms debounce means you can type and save rapidly without constant page flashing. Changes are batched and applied smoothly.

This is especially useful during development when you're frequently editing documentation.

## License

MIT