---
title: "Commands"
description: "CLI reference for all Dorcs commands."
tags: [commands, cli]
date: 2026-03-18
draft: false
---

# Commands

Dorcs has three commands: the default **server**, **init** to scaffold a new site, and **build** to export static HTML.

## `dorcs` — Start the server {badge:DEFAULT}

```bash
dorcs [flags]
```

Starts a local documentation server.

| Flag | Default | What it does |
| --- | --- | --- |
| `--dir` | `./docs` | Where your Markdown files are |
| `--addr` | `:8080` | Address to listen on |
| `--base-url` | `""` | URL prefix (e.g., `/docs`) |
| `--title` | `""` | Override the site title |
| `--config` | `""` | Path to a specific config file |
| `--repo` | `""` | Bootstrap docs from a GitHub repo |
| `--theme` | `""` | Override theme preset |
| `--theme-mode` | `""` | Override theme mode |
| `--watch` | `false` | Live reload on file changes |
| `--no-drafts` | `true` | Hide draft pages |
| `--cache` | `true` | Accepted by CLI (no current effect) |

**Examples:**

```bash
# Start with live reload
dorcs --dir ./docs --watch

# Serve from a GitHub repo
dorcs --repo https://github.com/owner/repo/tree/main/docs
```

## `dorcs init` — Create a new site

```bash
dorcs init [flags]
```

Scaffolds a starter docs site with an `index.md` and a `dorcs.yaml`.

| Flag | Default | What it does |
| --- | --- | --- |
| `--dir` | `./docs` | Where to create the docs folder |
| `--title` | `Documentation` | Starter site title |
| `--config` | `true` | Also create a config file |

**Example:**

```bash
dorcs init --dir ./handbook --title "Team Handbook"
```

## `dorcs build` — Export static HTML

```bash
dorcs build [flags]
```

Builds a deployable static site.

| Flag | Default | What it does |
| --- | --- | --- |
| `--dir` | `./docs` | Where your Markdown files are |
| `--output` | `./dist` | Where to write the built site |
| `--base-url` | `""` | URL prefix for deployment |
| `--title` | `""` | Override the site title |
| `--config` | `""` | Path to a specific config file |
| `--repo` | `""` | Bootstrap docs from a GitHub repo |
| `--theme` | `""` | Override theme preset |
| `--theme-mode` | `""` | Override theme mode |
| `--no-drafts` | `true` | Hide draft pages |

**Example:**

```bash
dorcs build --output ./dist --base-url /docs
```

## Good to know

> [!NOTE]
> - The server mode includes live search (`/api/search`) and sitemap generation — static builds only get the sitemap
> - `--watch` is automatically skipped when serving from GitHub
> - `--config` takes priority over `--repo`, which takes priority over auto-discovery
> - `--base-url` is sanitized to prevent invalid path segments
