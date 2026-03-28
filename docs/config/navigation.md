---
title: "Navigation"
description: "Control the sidebar, section tabs, and header links."
tags: [configuration, navigation]
date: 2026-03-18
draft: false
---

# Navigation

Dorcs builds your sidebar automatically from the folder structure. That works great for small sites — but when you need more control, there are two more powerful options.

## Automatic sidebar (default)

With no config at all, Dorcs scans your `docs/` folder and builds the sidebar from your files and folders. The order is determined by:

1. Numeric prefixes (`01_intro.md` comes before `02_setup.md`)
2. Front matter `order` field
3. Alphabetical as a fallback

This is enough for personal projects and small sites. Once you need custom labels or a specific structure, read on.

## Custom sidebar with `nav.items`

Define exactly what appears in the sidebar and in what order:

```yaml
nav:
  items:
    - Home: index.md
    - Getting Started: getting-started.md
    - Guides:
        page: guides/index.md
        items:
          - Installation: guides/install.md
          - Deployment: guides/deploy.md
    - API Reference: api.md
```

This gives you:

- **Custom labels** — "Getting Started" instead of "01_getting-started"
- **Exact ordering** — no more renaming files to reorder
- **Nested groups** — folders with their own section pages

Each entry can be written in several ways:

```yaml
# Simple: label → file
- Home: index.md

# With a landing page and children
- Guides:
    page: guides/index.md
    items:
      - Install: guides/install.md

# Group without a landing page
- Guides:
    items:
      - Install: guides/install.md
```

## Section tabs for large sites {badge:RECOMMENDED}

When your docs grow beyond a single sidebar, **section tabs** split them into top-level categories. Each tab gets its own sidebar. You're looking at this feature right now on this very site — notice the tabs in the header.

```yaml
nav:
  sections:
    - title: "Getting Started"
      items:
        - Overview: index.md
        - Installation: installation.md

    - title: "Guides"
      items:
        - Writing Docs: guides/writing.md
        - Markdown: guides/markdown.md

    - title: "Reference"
      items:
        - CLI Commands: reference/commands.md
        - Config Options: reference/config.md
```

How it works:

- The header shows a **row of clickable tabs** (Getting Started, Guides, Reference)
- Clicking a tab shows **only that section's pages** in the sidebar
- The active tab is **detected automatically** from the current page
- Each section's `items` uses the same syntax as `nav.items`

> [!IMPORTANT]
> When `nav.sections` is set, `nav.items` is ignored — sections contain their own items.

## Header links

Add external links to the header — great for linking to your GitHub repo, Discord, or main website:

```yaml
nav:
  links:
    - title: "GitHub"
      url: "https://github.com/you/your-repo"
      external: true
      icon: "github"

    - title: "Discord"
      url: "https://discord.gg/your-server"
      external: true
      icon: "discord"
```

Available icons: `github`, `twitter`, `discord`, `external`.

## Search

Search is enabled by default. It covers all your pages and provides instant results. You can toggle it:

```yaml
nav:
  show_search: true
```

> [!NOTE]
> Search works in server mode. Static builds (`dorcs build`) include a sitemap but not the live search API.
