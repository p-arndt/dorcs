---
title: "Configuration"
description: "Configure Dorcs with dorcs.yaml, dorcs.yml, or dorcs.json."
tags: [configuration]
date: 2026-03-18
draft: false
---

# Configuration

Dorcs works with defaults, but most sites should add a config file.

## Where config is loaded from

Dorcs checks in this order:

1. `./dorcs.yaml`
2. `./dorcs.yml`
3. `./dorcs.json`
4. the same filenames inside the docs directory

You can override discovery with `--config`.

## Minimal config

```yaml
site:
  title: "My Docs"

theme:
  preset: default
  mode: auto
```

## Override rules

- CLI flags override config values
- If `port` is set and `--addr` is left at its default, Dorcs listens on that port
- Theme preset controls the syntax highlighting theme automatically

## Main sections

### `site`

```yaml
site:
  title: "My Docs"
  description: "Internal platform documentation"
  logo: "/logo.png"
  favicon: "/favicon.ico"
```

Use this for branding and page metadata.
Store these assets in the docs tree. In multilingual setups, shared assets can live in the root `docs/` folder.

### `theme`

```yaml
theme:
  preset: ocean
  mode: auto
  custom_css: "custom.css"
  font_family: '"IBM Plex Sans", sans-serif'
  mono_font_family: '"JetBrains Mono", monospace'
```

You can also override colors:

```yaml
theme:
  preset: default
  colors:
    light:
      background: "#ffffff"
      foreground: "#101828"
      accent: "#0f766e"
    dark:
      background: "#09111f"
      foreground: "#e5eefc"
      accent: "#7dd3fc"
```

Available color keys:

- `background`
- `foreground`
- `muted`
- `border`
- `accent`
- `code_background`
- `header_background`
- `sidebar_background`

### `nav`

If you do nothing, Dorcs builds the sidebar from the folder tree.

Use `nav.items` when you need explicit order or labels:

```yaml
nav:
  show_search: true
  expand_all: false
  items:
    - Home: index.md
    - Getting Started: 01_getting-started.md
    - Usage:
        page: usage/index.md
        items:
          - Writing: usage/writing-your-docs.md
          - Metadata: usage/metadata.md
  links:
    - title: "GitHub"
      url: "https://github.com/example/repo"
      external: true
      icon: "github"
```

Current implementation note:

- `nav.items` and `nav.links` are active
- `show_search` is accepted in config, but the current templates always render the search trigger
- `expand_all` is accepted in config, but the current UI does not use it yet

Supported `nav.items` forms:

- `Label: path.md`
- `Label: { page: path.md }`
- `Label: { items: [...] }`
- `Label: { page: path.md, items: [...] }`

### `footer`

```yaml
footer:
  text: "Platform docs"
  copyright: "© 2026 Example"
  show_powered_by: true
```

### `auth`

Enable browser-based edit mode:

```yaml
auth:
  enabled: true
  username: "admin"
  password: "change-me"
```

Dorcs hashes the password on first use and stores sessions in `.dorcs_sessions.json` inside the docs directory by default.

### `languages`

```yaml
languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
```

When languages are enabled, Markdown files in the root `docs/` folder are ignored. Put content in language folders such as `docs/en/` and `docs/de/`.

### `versions`

```yaml
versions:
  default: "latest"
  enabled:
    - id: "latest"
      name: "Latest"
    - id: "v1"
      name: "Version 1"
```

With versioning enabled:

- default-version content lives in the root docs folder
- non-default versions live in folders like `docs/v1/`
- when languages are also enabled, default-version content lives in the default language folder and non-default versions live in `docs/<lang>/<version>/`

### `github`

```yaml
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: ${GITHUB_TOKEN}
  cache_ttl: "1h"
  edit_on_github:
    repository: "https://github.com/owner/repo/tree/main/docs"
```

Use this when content should be read from GitHub instead of the local docs directory, or when you want page-level "Edit on GitHub" links for local content.

## Example config

```yaml
port: 8080

site:
  title: "Dorcs"
  description: "Project documentation"

theme:
  preset: midnight
  mode: auto

nav:
  show_search: true
  expand_all: false

footer:
  text: "Built with Dorcs"
  show_powered_by: true
```

## Related pages

- [File Structure & Organization](./usage/file-structure.md)
- [Themes](./05_themes.md)
- [External Content](./external-content/index.md)
- [Edit](./08_edit.md)
