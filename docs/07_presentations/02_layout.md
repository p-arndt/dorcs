---
title: "Layout"
description: "Presentation layouts supported by Dorcs."
tags: [presentations, layout]
date: 2026-03-18
draft: false
---

# Layout

Use `_layout` for structure and `_class` for additional visual treatment.

## Available layouts

| Layout | Use case |
| --- | --- |
| `default` | Standard slide |
| `lead` | Title slide |
| `big` | One strong statement or metric |
| `quote` | Quote plus attribution |
| `left` | Left-aligned content |
| `right` | Right-aligned content |
| `columns-2` | Two flowing columns |
| `columns-3` | Three flowing columns |
| `split` | Two explicit panels |
| `timeline` | Horizontal timeline layout |

## Example: lead

```markdown
<!-- _layout: lead -->

# Dorcs
Markdown docs and slides in one place
```

## Example: split

```markdown
<!-- _layout: split -->

<div class="col">

## Problem
Too many disconnected docs.

</div>

<div class="col">

## Solution
One Markdown source for docs and decks.

</div>
```

## Example: auto columns

```markdown
<!-- _layout: columns-2 -->
<!-- _gap: loose -->

## Benefits

- Simple
- Fast
- Portable
- Easy to deploy
```

## Extra classes

Useful classes include:

- `lead`
- `left`
- `right`
- `fit`
- `invert`
- `two-columns`

Prefer `_layout` for structure. Use `_class` when you only want additional styling.
