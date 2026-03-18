---
title: "Order of Docs"
description: "How Dorcs determines page order in navigation."
tags: [usage, navigation]
date: 2026-03-18
draft: false
---

# Order of Docs

Dorcs can order pages automatically or follow an explicit navigation tree.

## Best option for most sites

Use `nav.items` when you need stable labels and exact ordering across folders.

## Automatic ordering

Without `nav.items`, Dorcs uses the file tree and common ordering hints:

1. numeric prefixes such as `01_`, `02_`
2. front matter hints like `order`
3. alphabetical fallback

Example:

```text
docs/
├── 01_intro.md
├── 02_install.md
└── 03_deploy.md
```

## Explicit ordering

```yaml
nav:
  items:
    - Home: index.md
    - Guides:
        page: guides/index.md
        items:
          - Install: guides/install.md
          - Deploy: guides/deploy.md
```

Use this when:

- filenames should stay technical
- labels should be more readable
- multiple sections need custom grouping

## Recommendation

Use filename prefixes for small sites. Use `nav.items` once the docs become product-facing or multi-section.
