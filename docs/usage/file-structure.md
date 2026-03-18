---
title: "File Structure & Organization"
description: "How Dorcs maps folders to routes, languages, and versions."
tags: [usage, structure]
date: 2026-03-18
draft: false
---

# File Structure & Organization

Dorcs uses your folder tree as the primary source of routing and navigation.

## Standard layout

```text
docs/
├── index.md
├── guide.md
└── api/
    ├── index.md
    └── auth.md
```

## URL mapping

| File | URL |
| --- | --- |
| `docs/index.md` | `/` |
| `docs/guide.md` | `/guide` |
| `docs/api/index.md` | `/api` |
| `docs/api/auth.md` | `/api/auth` |

## Ordering by filename

Numeric prefixes help when you use automatic navigation:

```text
docs/
├── 01_intro.md
├── 02_install.md
└── 03_deploy.md
```

The numbers influence order but are not part of the URL.

## Languages

```text
docs/
├── en/
│   ├── index.md
│   └── guide.md
└── de/
    ├── index.md
    └── guide.md
```

With:

```yaml
languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
```

URLs become:

- `/guide` for default language content in `docs/en/guide.md`
- `/de/guide` for German content

## Versions

```text
docs/
├── index.md
├── guide.md
└── v1/
    ├── index.md
    └── guide.md
```

With:

```yaml
versions:
  default: "latest"
  enabled:
    - id: "latest"
      name: "Latest"
    - id: "v1"
      name: "Version 1"
```

URLs become:

- `/guide`
- `/v1/guide`

Default-version content is served from the root docs folder. Non-default versions use folders.

If you configure `latest` as the default version, `/` serves the content from the root docs folder, not from `docs/latest/`.

## Languages and versions together

Use language-first folders:

```text
docs/
├── en/
│   ├── index.md
│   ├── guide.md
│   └── v1/
│       ├── index.md
│       └── guide.md
└── de/
    ├── index.md
    ├── guide.md
    └── v1/
        ├── index.md
        └── guide.md
```

Dorcs serves paths like `/de/v1/...`.
