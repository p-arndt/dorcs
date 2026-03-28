---
title: "Front Matter"
description: "Add metadata to your pages with YAML front matter."
tags: [usage, metadata]
date: 2026-03-18
draft: false
---

# Front Matter

Every page can start with a YAML block that sets its title, description, and other metadata:

```yaml
---
title: "Getting Started"
description: "Set up the project and run it locally"
date: 2026-03-18
draft: false
---
```

Everything between the `---` markers is front matter. The rest of the file is your Markdown content.

## Available fields

| Field | What it does |
| --- | --- |
| `title` | Page title — shown in sidebar, browser tab, and social sharing |
| `description` | Short summary for meta tags and link previews |
| `date` | Display date |
| `tags` | A list of tags (metadata, not displayed) |
| `draft` | Set to `true` to hide from visitors |
| `order` | Number to control sidebar position (lower = higher) |
| `author` | Author name (metadata) |
| `after` | Place this page after another by filename |

> [!TIP]
> `title` is the most important field. If you skip it, Dorcs uses the first `# Heading` in the file instead.

## Draft pages

Mark a page as draft to keep it out of the sidebar and builds:

```yaml
---
title: "Work in Progress"
draft: true
---
```

Drafts are hidden by default. To preview them during development:

```bash
dorcs --no-drafts=false --watch
```

## Presentation fields

Turn a page into a slide deck by adding these:

```yaml
---
presentation: true
presentation_header: "Engineering Team"
presentation_footer: "Q1 2026"
---
```

See [Presentations](../07_presentations/index.md) for the full guide.
