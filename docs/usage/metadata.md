---
title: "Metadata"
description: "Front matter fields supported by Dorcs."
tags: [usage, metadata]
date: 2026-03-18
draft: false
---

# Metadata

Dorcs reads YAML front matter from Markdown files.

## Example

```yaml
---
title: "Getting Started"
description: "Set up the project and run it locally"
date: 2026-03-18
tags:
  - guide
  - setup
draft: false
order: 10
author: "Team Docs"
after: "installation"
presentation: false
presentation_header: "Engineering"
presentation_footer: "Internal"
---
```

## Common fields

| Field | Type | Purpose |
| --- | --- | --- |
| `title` | string | Page title |
| `description` | string | Meta description |
| `date` | string | Display date |
| `tags` | list | Tag metadata |
| `draft` | bool | Hide when drafts are disabled |
| `order` | number | Manual order hint |
| `author` | string | Author metadata |
| `after` | string | Relative ordering hint |

## Presentation fields

| Field | Type | Purpose |
| --- | --- | --- |
| `presentation` | bool | Render page as a slide deck |
| `presentation_header` | string | Default slide header |
| `presentation_footer` | string | Default slide footer |

## Draft behavior

Draft pages are hidden by default because `--no-drafts` defaults to `true` in both server and build mode.
