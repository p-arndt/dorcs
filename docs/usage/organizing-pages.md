---
title: "Organizing Your Pages"
description: "How to structure your docs folder, control page order, and link between pages."
tags: [usage, structure]
date: 2026-03-18
draft: false
---

# Organizing Your Pages

Your folder structure IS your site structure. Every Markdown file becomes a page, every folder becomes a section.

## Files become pages

```text
docs/
├── index.md              → /
├── getting-started.md    → /getting-started
└── guides/
    ├── index.md          → /guides
    ├── install.md        → /guides/install
    └── deploy.md         → /guides/deploy
```

That's the core idea — one file, one page, clean URLs.

> [!TIP]
> Add an `index.md` to any folder to give that section its own landing page.

## Controlling sidebar order

By default, Dorcs sorts pages alphabetically. To control the order, prefix filenames with numbers:

```text
docs/
├── 01_intro.md       ← first in sidebar
├── 02_install.md     ← second
└── 03_deploy.md      ← third
```

The numbers affect sidebar order but **don't appear in URLs** — `01_intro.md` becomes `/intro`.

You can also set order in [front matter](./front-matter.md):

```yaml
---
title: "Introduction"
order: 1
---
```

> [!NOTE]
> For larger sites, skip file prefixes entirely and define your sidebar in `dorcs.yaml` using `nav.items` or `nav.sections`. See [Navigation](../config/navigation.md) for details.

## Linking between pages

Use regular Markdown links with relative paths:

```markdown
Check out the [install guide](./guides/install.md) for next steps.
```

Dorcs automatically rewrites `.md` links to clean URLs. So `./guides/install.md` becomes `/guides/install` for the reader. Links work in your editor preview AND on the rendered site.

## Adding images

Put images anywhere in your `docs/` folder and reference them with relative paths:

```markdown
![Architecture diagram](./images/overview.png)
```

## Writing tips

- **One topic per page** — don't cram everything into one file
- **Use folders to group related pages** — they become sidebar sections automatically
- **Start with the outcome** — tell readers what they'll achieve, then show how
- **Keep sections short** — if it feels long, split it into another page
- **Use examples over explanations** — a code snippet beats a paragraph

## Live reload

Run Dorcs with `--watch` and your browser refreshes automatically every time you save a file — including when you change `dorcs.yaml`:

```bash
dorcs --watch
```

> [!IMPORTANT]
> File watching through Docker volumes can be unreliable on Windows/Mac. Run `dorcs` directly on your machine for the best experience.
