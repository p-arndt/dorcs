---
title: "Slide Directives"
description: "Marpit-style slide directives supported by Dorcs."
tags: [presentations, directives]
date: 2026-03-18
draft: false
---

# Slide Directives

Add directives as HTML comments at the start of a slide.

## Example

```markdown
<!-- _layout: lead -->
<!-- _backgroundColor: #0f172a -->
<!-- _color: white -->

# Welcome
```

## Common directives

| Directive | Purpose |
| --- | --- |
| `_class` | Extra slide classes |
| `_layout` | Layout preset |
| `_color` | Text color |
| `_backgroundColor` | Background color |
| `_backgroundImage` | Background image |
| `_backgroundPosition` | Background position |
| `_backgroundSize` | Background sizing |
| `_backgroundRepeat` | Background repeat |
| `_header` | Slide header |
| `_footer` | Slide footer |
| `_paginate` | Hide or skip slide number |
| `_gap` | Tight, normal, or loose spacing |
| `_align` | Start, center, or end alignment |
| `_columns` | Two, three, or four columns |

## Spot vs inherited directives

- `_name`: applies only to the current slide
- `name`: applies to the current and following slides

```markdown
<!-- backgroundColor: #111827 -->
<!-- color: white -->

# Slide one

---

# Slide two

---

<!-- _backgroundColor: white -->
<!-- _color: black -->

# Slide three only
```
