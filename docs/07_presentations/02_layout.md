---
title: "Layout"
description: "Presentation layouts supported by Dorcs."
tags: [presentations, layout]
date: 2026-03-18
draft: false
---

# Layout

Use `_layout` as the primary directive for all slide appearance. Use `_class` only for custom CSS that doesn't map to a built-in layout.

## Available layouts

| Layout | Use case |
| --- | --- |
| `default` | Standard slide |
| `lead` | Title slide |
| `big` | One strong statement or metric |
| `quote` | Quote plus attribution |
| `left` | Left-aligned content |
| `right` | Right-aligned content |
| `columns-2` | Two flowing columns (CSS columns, no divs needed) |
| `columns-3` | Three flowing columns |
| `cols` | Two explicit side-by-side panels (use `::: col`) |
| `split` | Narrow left + wide right panels (use `::: col`) |
| `timeline` | Horizontal timeline layout |
| `fit` | Compact sizing for dense content |
| `invert` | Dark background, light text |

> **Note:** `two-columns` is a deprecated alias for `cols`. Use `cols` in new slides.

## Example: lead

```markdown
<!-- _layout: lead -->

# Dorcs
Markdown docs and slides in one place
```

## Example: cols

```markdown
<!-- _layout: cols -->

::: col

## Problem
Too many disconnected docs.

:::

::: col

## Solution
One Markdown source for docs and decks.

:::
```

## Example: split

```markdown
<!-- _layout: split -->

::: col

## Context
Narrow panel, 1fr.

:::

::: col

## Detail
Wide panel, 2fr. Main content goes here.

:::
```

## Example: auto columns

Content after the heading flows automatically into columns. Put section sub-headings (`###`) directly after the slide title — no intro paragraph between them.

```markdown
<!-- _layout: columns-2 -->

## Key Points

### What it does

Content flows into the left column automatically.

### How it works

CSS multi-column balances the remaining content.
```

> **Don't mix auto columns with `::: col`.** `columns-2`/`columns-3` flow content
> automatically — they don't use `::: col` blocks. If you want explicit
> side-by-side panels you control, use `cols` (or `split`) instead. The slide
> title (`#`/`##`) always spans full width above the columns; use `###` and below
> for sub-headings inside a column.

## Combining layout and class

Pass space-separated values to `_layout` to combine a layout with extra classes:

```markdown
<!-- _layout: lead invert -->
```

This is equivalent to `_layout: lead` + `_class: invert`.
