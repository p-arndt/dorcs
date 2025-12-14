---
title: "Order of Docs"
description: "How to order your docs"
tags: [ordering, docs]
date: 2025-12-14
draft: false
---

# Order of Docs

## Using `order` and Numeric Prefixes

By default, pages are sorted by:

1. Numeric prefixes in filenames (e.g., `01_`, `02_`)
2. `order` field in front matter
3. Alphabetical by title

Files with numeric prefixes (like `01_installation.md`) will appear before files without prefixes, regardless of their `order` value.

## Using `after` for Relative Positioning

The `after` field allows you to place a page directly after another page in the navigation, taking **absolute precedence** over numeric prefixes and `order` values.

**Place a page right after the index:**

```yaml
---
title: "Getting Started"
after: "index"
---
# Getting Started
```

This will place "Getting Started" immediately after the root `index.md` or folder `index.md`, before any numbered files.

**Place a page after a specific page:**

```yaml
---
title: "Advanced Topics"
after: "getting-started"
---
# Advanced Topics
```

This will place "Advanced Topics" directly after the page with key `getting-started`.

**Note:** The `after` field takes absolute precedence over all other sorting mechanisms (numeric prefixes, `order` fields, etc.). If multiple pages use `after: "index"`, they will be sorted among themselves using the normal sorting rules.
