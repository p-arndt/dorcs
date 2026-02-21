---
title: "Slide Directives"
description: "Marpit-compatible directives for slide layout and styling"
tags: [presentations, slides, markdown]
date: 2026-02-21
draft: false
after: "presentations"
---

# Slide Directives

Dorcs supports [Marpit-compatible directives](https://marpit.marp.app/directives). Add HTML comments at the **start** of a slide:

| Directive | Example | Description |
|-----------|---------|-------------|
| `_class` | `<!-- _class: lead -->` | Layout class(es) for the slide |
| `_color` | `<!-- _color: white -->` | Text color |
| `_backgroundColor` | `<!-- _backgroundColor: #1a1a2e -->` | Background color |
| `_backgroundImage` | `<!-- _backgroundImage: url(bg.jpg) -->` | Background image (defaults: cover, center, no-repeat) |
| `_backgroundPosition` | `<!-- _backgroundPosition: top -->` | Background position |
| `_backgroundSize` | `<!-- _backgroundSize: contain -->` | Background size |
| `_backgroundRepeat` | `<!-- _backgroundRepeat: repeat -->` | Background repeat |
| `_header` | `<!-- _header: Slide title -->` | Slide-specific header |
| `_footer` | `<!-- _footer: Confidential -->` | Slide-specific footer |
| `_paginate` | `<!-- _paginate: false -->` | Hide slide number (`false` or `skip`) |
| `_layout` | `<!-- _layout: columns-2 -->` | Layout preset (see Layout Control) |
| `_gap` | `<!-- _gap: loose -->` | Spacing between columns/items: `tight`, `normal`, `loose` |
| `_align` | `<!-- _align: start -->` | Content alignment: `start`, `center`, `end` |
| `_columns` | `<!-- _columns: 3 -->` | Number of equal columns: `2`, `3`, `4` |

## Spot vs inherited directives

- **`_` prefix** (spot): Applies only to the current slide.
- **No prefix**: Applies to the current slide and all following slides (inherited).

```markdown
<!-- backgroundColor: aqua -->
# Slide 1 - has aqua background

---
# Slide 2 - also has aqua background (inherited)

---
<!-- _backgroundColor: white -->
# Slide 3 - white only (spot), next slide inherits aqua again
```