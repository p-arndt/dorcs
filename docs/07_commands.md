---
title: "Commands"
description: "CLI reference for server, init, and build commands."
tags: [commands, cli]
date: 2026-03-18
draft: false
---

# Commands

Dorcs has three main entry points:

- default server mode
- `init`
- `build`

## `dorcs`

Starts the docs server.

```bash
dorcs [flags]
```

Flags:

| Flag | Default | Purpose |
| --- | --- | --- |
| `--dir` | `./docs` | Docs directory |
| `--addr` | `:8080` | Listen address |
| `--base-url` | `""` | URL prefix such as `/docs` |
| `--title` | `""` | Override site title |
| `--cache` | `true` | Currently accepted by the CLI, but not behaviorally significant in the current server implementation |
| `--no-drafts` | `true` | Hide draft pages |
| `--config` | `""` | Explicit config file |
| `--theme` | `""` | Override theme preset |
| `--theme-mode` | `""` | Override theme mode |
| `--watch` | `false` | Enable file watching and live reload |

Example:

```bash
dorcs --dir ./docs --watch --addr 127.0.0.1:8080
```

## `dorcs init`

Creates a starter docs site.

```bash
dorcs init [flags]
```

Flags:

| Flag | Default | Purpose |
| --- | --- | --- |
| `--dir` | `./docs` | Directory to initialize |
| `--title` | `Documentation` | Starter site title |
| `--config` | `true` | Create a basic config file |

Example:

```bash
dorcs init --dir ./handbook --title "Team Handbook"
```

## `dorcs build`

Builds a static site.

```bash
dorcs build [flags]
```

Flags:

| Flag | Default | Purpose |
| --- | --- | --- |
| `--dir` | `./docs` | Docs directory |
| `--output` | `./dist` | Output directory |
| `--base-url` | `""` | URL prefix for deployment |
| `--title` | `""` | Override site title |
| `--no-drafts` | `true` | Hide draft pages |
| `--config` | `""` | Explicit config file |
| `--theme` | `""` | Override theme preset |
| `--theme-mode` | `""` | Override theme mode |

Example:

```bash
dorcs build --output ./dist --base-url /docs
```

## Notes

- The default server mode also exposes search and sitemap endpoints
- `--watch` is skipped for GitHub-backed docs
- `--base-url` is sanitized to avoid invalid path segments
- static builds generate `sitemap.xml`, but not a live search API
