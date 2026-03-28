---
title: "Customize Your Site"
description: "Make your Dorcs site look and work exactly how you want."
tags: [configuration]
date: 2026-03-18
draft: false
---

# Customize Your Site

Everything in Dorcs is controlled through a single `dorcs.yaml` file. Here's how to make your site yours.

## The config file

Dorcs looks for `dorcs.yaml` (or `.yml` / `.json`) in your project root, then inside `docs/`. You can also point to a specific file:

```bash
dorcs --config ./my-config.yaml
```

> [!TIP]
> CLI flags always override config values. So `dorcs --theme ocean` beats whatever's in your yaml.

## What can you customize?

| Guide | What it covers |
| --- | --- |
| [Branding](./branding.md) | Site title, logo, favicon, description, footer, announcement banners |
| [Themes](../05_themes.md) | Color presets, light/dark mode, custom colors, fonts, CSS |
| [Navigation](./navigation.md) | Sidebar order, section tabs, header links, search |
| [Languages & Versions](./languages-versions.md) | Multi-language sites, doc versioning |
| [External Content](../external-content/index.md) | Serve docs from GitHub, "Edit on GitHub" links |
| [Edit Mode](../08_edit.md) | Browser-based editing with login |

## Quick start config

Here's a practical config that covers the essentials:

```yaml
site:
  title: "My Project"
  description: "Documentation for My Project"

theme:
  preset: ocean
  mode: auto

nav:
  show_search: true
  items:
    - Home: index.md
    - Getting Started: getting-started.md
    - Guides:
        page: guides/index.md
        items:
          - Install: guides/install.md
          - Deploy: guides/deploy.md

footer:
  text: "Built with Dorcs"
  show_powered_by: true
```

## What works without any config

Even with zero configuration, every page automatically gets:

- **Sidebar navigation** built from your folder structure
- **Table of contents** from your headings (you can see it on the right)
- **Breadcrumbs** showing where you are (like "Docs > Config > Branding")
- **Previous / Next links** at the bottom of every page
- **Code copy buttons** on every code block
- **Heading anchors** — hover any heading to get a shareable link
- **Back to top button** after scrolling down
- **Last updated date** from the file's modification time
- **Search** across all your pages
- **Sitemap** for search engines
- **Open Graph meta tags** for rich link previews on social media
- **Custom 404 page** matching your theme
