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

If you start Dorcs with `--repo`, config bootstrap changes:

1. `dorcs.yaml`, `dorcs.yml`, `dorcs.json` at the GitHub repo root
2. the same filenames at the repo path from `--repo` such as `docs/`

If no remote config file exists, Dorcs keeps running with defaults.

Precedence is:

1. `--config`
2. `--repo`
3. local discovery

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

Supported `nav.items` forms:

- `Label: path.md`
- `Label: { page: path.md }`
- `Label: { items: [...] }`
- `Label: { page: path.md, items: [...] }`

#### Section tabs

For larger documentation sites, use `nav.sections` to add a second header row with section tabs. Each section groups its own sidebar navigation. Clicking a tab shows only that section's pages in the sidebar.

```yaml
nav:
  sections:
    - title: "Getting Started"
      items:
        - Overview: index.md
        - Installation: installation.md
    - title: "Configuration"
      items:
        - Options: configuration.md
        - Themes: themes.md
    - title: "Reference"
      items:
        - CLI: commands.md
        - API: api.md
```

When `nav.sections` is configured:

- The header displays a second row of clickable section tabs
- The sidebar shows only the active section's items
- The active section is detected automatically from the current page
- `nav.items` is ignored (sections contain their own items)

Each section's `items` follows the same syntax as `nav.items`.

### `announcement`

Show a banner at the top of every page:

```yaml
announcement:
  text: 'Version 2.0 is here! <a href="/changelog">See what changed</a>'
  dismissible: true
```

- `text` supports basic HTML (links, bold, etc.)
- `dismissible` (default: `true`) shows a close button; dismissed state is saved in localStorage

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

You can also skip local config mounting entirely and bootstrap from a repo:

```bash
dorcs --repo https://github.com/owner/repo/tree/main/docs
```

In repo mode, Dorcs loads both docs and config from that repository. The fetched config does not need to repeat `github.enabled` or `github.repository`.

## Built-in page features

Dorcs automatically adds several UI features to every page. These require no configuration and work out of the box.

### Breadcrumbs

A navigation trail is shown above the page content (e.g., "Docs / Guide / Intro"). Breadcrumbs are derived from the sidebar navigation structure. When `nav.sections` is configured, breadcrumbs reflect the active section's hierarchy.

### Previous / Next navigation

Sequential page links appear at the bottom of every page. The order follows the sidebar navigation — either the configured `nav.items` / `nav.sections` order or the auto-detected file order. This lets readers move through docs linearly without returning to the sidebar.

### Code copy button

Every fenced code block gets a copy button in the top-right corner. It appears on hover (always visible on mobile), copies the code content to the clipboard, and shows a checkmark for 2 seconds to confirm.

### Heading anchor links

Hovering over any heading (h1–h6) reveals a `#` link on the right. Clicking it copies the full permalink URL to the clipboard and updates the browser address bar. This makes it easy to share deep links to specific sections.

### Back to top button

A floating button appears in the bottom-right corner after scrolling down 400 pixels. Clicking it smoothly scrolls back to the top of the page.

### Last updated date

The file modification time of each Markdown file is shown at the bottom of the content area as "Last updated: January 2, 2026". This updates automatically when the source file changes.

### Custom 404 page

When a page is not found, Dorcs renders a styled 404 page that matches your site's theme. It includes a link to the homepage and a button to open the search overlay, so users can quickly find what they were looking for.

### Sitemap

Dorcs automatically serves a `sitemap.xml` at `/sitemap.xml` in server mode. In static build mode (`dorcs build`), it generates a `sitemap.xml` file in the output directory. The sitemap includes all non-draft pages with last-modified dates, priority, and change frequency.

### Open Graph meta tags

Dorcs generates `og:title`, `og:description`, `og:site_name`, and `twitter:card` meta tags on every page. These are populated from front matter (`title`, `description`) and the site title from config. When you share a doc URL on Slack, Twitter, or other platforms, the link preview will show the page title and description.

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
