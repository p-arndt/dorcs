---
title: "Commands"
description: "Complete reference for all dorcs commands"
tags: [commands, cli, reference]
date: 2025-12-18
draft: false
---

# Commands

dorcs provides three main commands: `init`, `build`, and `server` (default). This guide covers all available options and usage examples.

## Overview

| Command  | Description                                         |
| -------- | --------------------------------------------------- |
| `init`   | Initialize a new documentation site                 |
| `build`  | Generate a static site for deployment               |
| `server` | Start the development server (default - not needed) |

## Init Command

The `init` command sets up a new documentation site with a basic structure.

### Usage

```bash
dorcs init [options]
```

### Options

| Flag       | Default         | Description                                                   |
| ---------- | --------------- | ------------------------------------------------------------- |
| `--dir`    | `./docs`        | Directory to initialize (will be created if it doesn't exist) |
| `--title`  | `Documentation` | Site title                                                    |
| `--config` | `true`          | Create a basic `dorcs.yaml` config file                       |

### Examples

```bash
# Basic initialization
dorcs init

# Custom directory and title
dorcs init --dir ./my-docs --title "My Awesome Docs"

# Skip config file creation
dorcs init --config=false
```

### What It Creates

- **`docs/` directory** (or your specified directory)
- **`index.md`** - Welcome page with sample content and frontmatter
- **`dorcs.yaml`** - Basic configuration file (if `--config=true`)

> [!NOTE]
> The init command is safe to run multiple times. It will warn you if files already exist and skip creating them.

## Build Command

The `build` command generates a static site from your markdown files, ready for deployment.

### Usage

```bash
dorcs build [options]
```

### Options

| Flag           | Default  | Description                                                                               |
| -------------- | -------- | ----------------------------------------------------------------------------------------- |
| `--dir`        | `./docs` | Directory containing markdown documents                                                   |
| `--output`     | `./dist` | Output directory for generated static site                                                |
| `--base-url`   | `""`     | Optional base URL path prefix (e.g. `/docs`). No trailing slash                           |
| `--title`      | `""`     | Site title shown in header (overrides config file)                                        |
| `--no-drafts`  | `true`   | Hide documents with front matter `draft: true`                                            |
| `--config`     | `""`     | Path to config file (default: looks for `dorcs.yaml` in current directory, then docs dir) |
| `--theme`      | `""`     | Theme preset: `default`, `ocean`, `forest`, `sunset`, `midnight`, `lavender`, `rose`      |
| `--theme-mode` | `""`     | Theme mode: `light`, `dark`, `auto`                                                       |

### Examples

```bash
# Basic build
dorcs build

# Custom output directory
dorcs build --output ./public

# Build with base URL for GitHub Pages
dorcs build --base-url /my-repo --output ./docs

# Override theme
dorcs build --theme ocean --theme-mode dark

# Include draft documents
dorcs build --no-drafts=false
```

### Output

The build command generates a complete static site in the output directory:

- All markdown files converted to HTML
- Static assets (CSS, JS, images)
- Syntax highlighting CSS
- Search index
- Sitemap

## Server Command (Default)

The server command starts a development server to preview your documentation locally. This is the default mode when no command is specified.

### Usage

```bash
dorcs [options]
# or explicitly:
dorcs server [options]
```

### Options

| Flag           | Default  | Description                                                                               |
| -------------- | -------- | ----------------------------------------------------------------------------------------- |
| `--dir`        | `./docs` | Directory containing markdown documents                                                   |
| `--addr`       | `:8080`  | Listen address (e.g. `:8080`, `127.0.0.1:8080`)                                           |
| `--base-url`   | `""`     | Optional base URL path prefix (e.g. `/docs`). No trailing slash                           |
| `--title`      | `""`     | Site title shown in header (overrides config file)                                        |
| `--cache`      | `true`   | Cache rendered documents in memory (mtime-based)                                          |
| `--no-drafts`  | `true`   | Hide documents with front matter `draft: true`                                            |
| `--config`     | `""`     | Path to config file (default: looks for `dorcs.yaml` in current directory, then docs dir) |
| `--theme`      | `""`     | Theme preset: `default`, `ocean`, `forest`, `sunset`, `midnight`, `lavender`, `rose`      |
| `--theme-mode` | `""`     | Theme mode: `light`, `dark`, `auto`                                                       |
| `--watch`      | `false`  | Watch for file changes and automatically reload                                           |

### Examples

```bash
# Start server on default port
dorcs

# Enable watch mode for live reload
dorcs --watch

# Custom port and directory
dorcs --addr :3000 --dir ./my-docs

# Disable caching (useful for development)
dorcs --cache=false --watch

# Custom theme
dorcs --theme ocean --theme-mode dark
```

### Watch Mode

When `--watch` is enabled:

- File changes trigger automatic page reload
- Config file changes are detected and applied
- Live reload uses Server-Sent Events (SSE)

> [!TIP]
> Always use `--watch` during development for the best experience!

## Getting Help

View help for any command:

```bash
# General help (shows all commands)
dorcs --help

# Command-specific help
dorcs init --help
dorcs build --help
```

## Common Workflows

### Initial Setup

```bash
# 1. Initialize your site
dorcs init --title "My Documentation"

# 2. Start development server
dorcs --watch
```

### Development

```bash
# Start server with watch mode
dorcs --watch

# Edit your markdown files - changes reload automatically!
```

### Production Build

```bash
# Build static site
dorcs build --output ./public

# Deploy the ./public directory to your hosting provider
```

### GitHub Pages

```bash
# Build with base URL matching your repository name
dorcs build --base-url /my-repo --output ./docs

# Commit and push the ./docs directory
```

## Configuration Priority

Command-line flags override configuration file settings:

1. **Config file** (`dorcs.yaml`) - Base configuration
2. **Command-line flags** - Override config file values

For example:

```bash
# Config file has theme: "default"
# Command line overrides it:
dorcs --theme ocean
```

