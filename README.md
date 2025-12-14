# Dorcs

<div style="display: flex; justify-content: center; align-items: center; padding: 20px;">
<img src="./docs/logo.png" alt="Dorcs Logo" width="200" height="200" style="border: none;" />
</div>

Dorcs is a single-binary static documentation server for Markdown files. It is a simple and easy to use documentation server that allows you to create and host your documentation site in minutes.

## Documentation 

Check out the [documentation](https://dorcs.allthing.eu) for more information and seeing it in action.

## Features

- **Single binary** – no runtime dependencies, statically linkable
- **Extensionless URLs** – `/guide/getting-started` serves `docs/guide/getting-started.md`
- **YAML front matter** – metadata support (title, description, date, tags, draft)
- **Table of Contents** – auto-generated from headings with scrollspy
- **Navigation sidebar** – built automatically from your folder structure
- **Responsive design** – mobile-friendly with collapsible sidebar
- **Dark mode** – automatic based on system preference
- **Live reload** – watch mode for development with smart content updates
- **Multiple themes** – choose from 20+ built-in themes
- **Search** – built-in search functionality

## Quick Start

### Option 1: Use pre-built binary

1. Download the latest release from the [releases page](https://github.com/p-arndt/dorcs/releases)
2. Make it executable (Linux/macOS): `chmod +x dorcs`
3. Run: `./dorcs --dir ./docs`

### Option 2: Build from source

```sh
# Run from source
go run ./cmd/dorcs

# Build and run
go build -o dorcs ./cmd/dorcs
./dorcs

# Run with auto-reload and live browser refresh (for development)
./dorcs --watch
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

| Flag          | Default  | Description                                                                          |
| ------------- | -------- | ------------------------------------------------------------------------------------ |
| `-dir`        | `./docs` | Directory containing markdown documents                                              |
| `-addr`       | `:8080`  | Listen address (e.g., `:8080`, `127.0.0.1:8080`)                                     |
| `-base-url`   | `""`     | URL path prefix (e.g., `/docs`)                                                      |
| `-title`      | `""`     | Site title shown in header (overrides config file)                                   |
| `-cache`      | `true`   | Cache rendered documents in memory                                                   |
| `-no-drafts`  | `true`   | Hide documents with `draft: true` front matter                                       |
| `-config`     | `""`     | Path to config file (default: looks for `dorcs.yaml` in docs dir)                    |
| `-theme`      | `""`     | Theme preset: `default`, `ocean`, `forest`, `sunset`, `midnight`, `lavender`, `rose` |
| `-theme-mode` | `""`     | Theme mode: `light`, `dark`, `auto`                                                  |
| `-watch`      | `false`  | Watch for file changes and automatically reload                                      |

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

| Field         | Type   | Description                              |
| ------------- | ------ | ---------------------------------------- |
| `title`       | string | Page title (used in nav, browser tab)    |
| `description` | string | Meta description                         |
| `date`        | string | Publication date (YYYY-MM-DD)            |
| `tags`        | list   | Tags for categorization                  |
| `draft`       | bool   | If true, hidden when `-no-drafts` is set |

## Configuration File

dorcs supports a configuration file for advanced customization. Place a `dorcs.yaml` (or `dorcs.yml` or `dorcs.json`) file in your docs directory.


### Theme Presets

| Preset     | Description                            |
| ---------- | -------------------------------------- |
| `default`  | GitHub-inspired light/dark theme       |
| `ocean`    | Cool blue tones                        |
| `forest`   | Natural green palette                  |
| `sunset`   | Warm orange/brown tones                |
| `midnight` | Deep purple/blue (Catppuccin-inspired) |
| `lavender` | Soft purple aesthetic                  |
| `rose`     | Pink/rose accent colors                |

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

## Development

```sh
# Run with auto-reload for markdown changes
go run ./cmd/dorcs --dir ./docs --watch
```

The server will automatically detect changes to:
- .md and .markdown files
- New files and directories
- File deletions and renames

# Run tests
```sh
go test ./...
```

## License

MIT
