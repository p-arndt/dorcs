---
title: "Getting Started"
description: "Build a real documentation site from scratch in 5 minutes."
tags: [getting-started, quickstart]
date: 2026-03-18
draft: false
---

# Getting Started

By the end of this page, you'll have a working documentation site with multiple pages, a custom theme, and organized navigation. Let's go.

## Step 1: Create your site

```bash
dorcs init
dorcs --watch
```

Open [http://localhost:8080](http://localhost:8080). You've got a docs site. It has one page (`docs/index.md`) and uses the default theme.

> [!TIP]
> The `--watch` flag gives you live reload — every time you save a file, the browser updates automatically. Keep it running while you work.

## Step 2: Add some pages

Your site needs more than a homepage. Create a few Markdown files:

```text
docs/
├── index.md
├── 01_getting-started.md
├── 02_features.md
└── guides/
    ├── index.md
    └── first-steps.md
```

Each file becomes a page with a clean URL:

| File | URL |
| --- | --- |
| `docs/index.md` | `/` |
| `docs/01_getting-started.md` | `/getting-started` |
| `docs/02_features.md` | `/features` |
| `docs/guides/index.md` | `/guides` |
| `docs/guides/first-steps.md` | `/guides/first-steps` |

The sidebar navigation updates automatically — every page you add shows up. The `01_`, `02_` prefixes control the order but don't appear in URLs.

## Step 3: Write real content

Open any page and write Markdown. Here's a taste of what you can do — beyond the basics:

```markdown
> [!TIP]
> Use callouts to highlight important information.

:::tabs
::tab macOS
\`\`\`bash
brew install my-tool
\`\`\`
::tab Linux
\`\`\`bash
apt install my-tool
\`\`\`
:::

This feature is {badge:NEW} in version 2.0.
```

That gives you styled callouts, tabbed content, and inline badges — all from plain Markdown. See [Extensions](./usage/extensions.md) for everything available.

## Step 4: Pick a theme

Open `dorcs.yaml` and change the theme:

```yaml
site:
  title: "My Project Docs"

theme:
  preset: ocean
  mode: auto
```

Try `ocean`, `forest`, `sunset`, `midnight`, `lavender`, `rose`, or any of the other presets. Your site updates instantly thanks to `--watch`. Browse them all on the [Themes](./05_themes.md) page.

## Step 5: Organize your navigation

For a small site, the automatic sidebar is fine. But once you have more pages, take control with `nav.items` in your config:

```yaml
nav:
  items:
    - Home: index.md
    - Getting Started: 01_getting-started.md
    - Features: 02_features.md
    - Guides:
        page: guides/index.md
        items:
          - First Steps: guides/first-steps.md
```

This gives you custom labels and exact ordering. For even bigger sites, you can use **section tabs** — see [Navigation](./config/navigation.md).

## Step 6: Ship it

Build static HTML and deploy anywhere:

```bash
dorcs build
```

Upload the `dist/` folder to GitHub Pages, Netlify, Vercel, or any web server. See [Deployment](./04_deployment.md) for platform-specific guides.

## What you've built

In just a few minutes, you now have:

- A multi-page docs site with automatic navigation
- A custom theme with light/dark mode
- Clean URLs, search, code copy buttons, and a table of contents
- A static build ready for any host

## Keep going

:::tabs
::tab I want to customize more
Learn how to add your logo, announcement banners, and custom colors in [Customize Your Site](./config/index.md).
::tab I want to learn all the Markdown features
See [Writing Docs](./usage/index.md) for tabs, badges, diagrams, math, video embeds, and more.
::tab I need the full CLI reference
Check [Commands](./07_commands.md) for every flag and option.
:::
