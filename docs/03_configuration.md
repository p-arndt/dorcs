---
title: "Config Reference"
description: "Complete dorcs.yaml reference with all available settings."
tags: [configuration, reference]
date: 2026-03-18
draft: false
---

# Config Reference

This is the full reference for every `dorcs.yaml` setting. For practical guides on how to use these, see [Customize Your Site](./config/index.md).

> [!TIP]
> Dorcs looks for `dorcs.yaml`, `dorcs.yml`, or `dorcs.json` in your project root first, then inside `docs/`. Use `--config` to point to a specific file. CLI flags always override config values.

## Complete config example

```yaml
# Server port (default: 8080)
port: 8080

# Site metadata and branding
site:
  title: "My Docs"
  description: "Project documentation"
  logo: "/logo.png"
  favicon: "/favicon.ico"

# Theme and styling
theme:
  preset: ocean          # default, ocean, forest, sunset, midnight, lavender, rose, nord, gruvbox, dracula, solarized, mono, cyberpunk
  mode: auto             # light, dark, or auto
  custom_css: "custom.css"
  font_family: '"Inter", system-ui, sans-serif'
  mono_font_family: '"JetBrains Mono", monospace'
  colors:
    light:
      background: "#ffffff"
      foreground: "#1f2328"
      muted: "#57606a"
      border: "#d0d7de"
      accent: "#0969da"
      code_background: "#f6f8fa"
      header_background: "#ffffff"
      sidebar_background: "#fafbfc"
    dark:
      background: "#0d1117"
      foreground: "#e6edf3"
      muted: "#8b949e"
      border: "#30363d"
      accent: "#2f81f7"
      code_background: "#161b22"

# Navigation
nav:
  show_search: true
  expand_all: false

  # Option A: flat sidebar
  items:
    - Home: index.md
    - Guide:
        page: guide/index.md
        items:
          - Install: guide/install.md

  # Option B: section tabs (replaces items)
  sections:
    - title: "Getting Started"
      items:
        - Overview: index.md
    - title: "Reference"
      items:
        - CLI: commands.md

  # Header links
  links:
    - title: "GitHub"
      url: "https://github.com/you/repo"
      external: true
      icon: "github"    # github, twitter, discord, external

# Announcement banner
announcement:
  text: 'Check out <a href="/changelog">what is new</a>'
  dismissible: true

# Footer
footer:
  text: "Built with Dorcs"
  copyright: "© 2026 Your Name"
  show_powered_by: true
  links:
    - title: "Privacy"
      url: "/privacy"

# Multi-language support
languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"

# Doc versioning
versions:
  default: "latest"
  enabled:
    - id: "latest"
      name: "Latest"
    - id: "v1"
      name: "Version 1"

# GitHub integration
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: ${GITHUB_TOKEN}
  cache_ttl: "1h"
  edit_on_github:
    repository: "https://github.com/owner/repo/tree/main/docs"

# Browser-based editing
auth:
  enabled: true
  username: "admin"
  password: "change-me"
  sessions_path: ".dorcs_sessions.json"
```

## Config discovery order

1. `--config` flag (highest priority)
2. `--repo` flag (GitHub bootstrap)
3. `./dorcs.yaml` / `.yml` / `.json` in project root
4. Same filenames inside the `docs/` directory

## Related guides

- [Branding](./config/branding.md) — title, logo, footer, announcements
- [Themes](./05_themes.md) — presets and custom colors
- [Navigation](./config/navigation.md) — sidebar, section tabs, header links
- [Languages & Versions](./config/languages-versions.md) — multi-language and versioning
- [External Content](./external-content/index.md) — GitHub as a docs source
- [Edit Mode](./08_edit.md) — browser-based editing
