---
title: "Order of Docs"
description: "How to order your docs"
tags: [ordering, docs]
date: 2025-12-14
draft: false
---

# Order of Docs

## Two Navigation Modes

dorcs supports two ways to control navigation:

1. Explicit navigation in `dorcs.yaml` with `nav.items`
2. Automatic navigation from folders and files

If `nav.items` is present, it fully defines the sidebar order and grouping.

```yaml
nav:
  items:
    - Home: index.md
    - Getting Started: 01_getting-started.md
    - Usage:
        page: usage/index.md
        items:
          - Writing Your Docs: usage/writing-your-docs.md
          - Metadata: usage/metadata.md
```

## Automatic Navigation

If `nav.items` is not configured, dorcs builds navigation automatically from your docs tree.

Automatic ordering prefers:

1. Numeric prefixes in filenames (for example `01_`, `02_`)
2. `order` field in front matter
3. Alphabetical by title

Files with numeric prefixes such as `01_installation.md` appear before files without prefixes.