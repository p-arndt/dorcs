---
title: "Writing Your Docs"
description: "Practical writing patterns for Dorcs documentation."
tags: [usage, writing]
date: 2026-03-18
draft: false
---

# Writing Your Docs

Dorcs works best when the file tree mirrors the way readers should navigate the docs.

## Recommended pattern

- keep one topic per page
- use folders for sections
- add `index.md` to section folders
- keep titles reader-facing, not file-facing

## Example

```text
docs/
├── index.md
├── getting-started.md
└── guides/
    ├── index.md
    ├── install.md
    └── deploy.md
```

This produces:

- `/`
- `/getting-started`
- `/guides`
- `/guides/install`
- `/guides/deploy`

## Linking pages

Use normal Markdown links:

```markdown
[Install guide](./guides/install.md)
```

Dorcs rewrites internal Markdown links to clean URLs.

## Assets

Store images and other shared files in `docs/` and reference them with relative paths:

```markdown
![Architecture](./images/overview.png)
```

## Writing style

- start with the outcome
- keep sections short
- use examples instead of long explanation blocks
- put full reference detail on reference pages, not onboarding pages
