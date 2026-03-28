---
title: "Languages & Versions"
description: "Serve your docs in multiple languages and maintain multiple versions."
tags: [configuration, languages, versions]
date: 2026-03-18
draft: false
---

# Languages & Versions {badge:BETA}

Dorcs can serve your docs in multiple languages and maintain multiple versions side by side.

## Multiple languages

To offer your docs in more than one language, create a folder per language and tell Dorcs about them:

```text
docs/
├── en/
│   ├── index.md
│   └── guide.md
└── de/
    ├── index.md
    └── guide.md
```

```yaml
languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
```

Dorcs automatically adds a **language switcher** to the header. The default language is served at the root URL (`/guide`), other languages get a prefix (`/de/guide`).

> [!IMPORTANT]
> When languages are enabled, Markdown files in the root `docs/` folder are ignored. All content must live in language subfolders like `docs/en/` and `docs/de/`.

## Doc versioning

Keep older versions of your docs alongside the current ones:

```text
docs/
├── index.md        ← latest (default)
├── guide.md
└── v1/
    ├── index.md    ← version 1
    └── guide.md
```

```yaml
versions:
  default: "latest"
  enabled:
    - id: "latest"
      name: "Latest"
    - id: "v1"
      name: "Version 1"
```

The default version lives in the root docs folder — no subfolder needed. Other versions go in their own folders. Dorcs adds a **version switcher** to the header.

URLs:
- `/guide` — latest version
- `/v1/guide` — version 1

## Using both together

When you need languages AND versions, put languages first:

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

This gives you URLs like:

| Content | URL |
| --- | --- |
| English, latest | `/guide` |
| English, v1 | `/v1/guide` |
| German, latest | `/de/guide` |
| German, v1 | `/de/v1/guide` |
