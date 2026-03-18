---
title: "Layout"
description: "Layout options for slide decks"
tags: [presentations, slides, markdown]
date: 2026-02-21
draft: false
---

# Layout Classes

- **`lead`** – Title slide: larger headings, more dramatic typography
- **`left`** – Left-aligned content
- **`right`** – Right-aligned content
- **`fit`** – Compact layout, fits more content
- **`invert`** – Inverted colors (dark background, light text)
- **`two-columns`** – Two-column layout (use `<div class="col">` for each column)

## Layout Control

Use `_layout` for structure and `_class` for extra styling (e.g. `_class: invert`). Layout primitives compose:

| Layout | Description |
|--------|-------------|
| `default` | Centered, single column |
| `lead` | Large title, subtitle |
| `big` | Large centered text — stats, single statement (no heading required) |
| `quote` | Large centered quote; last paragraph = smaller attribution |
| `left` | Left-aligned content |
| `right` | Right-aligned content |
| `columns-2` | Two equal columns (auto flow, no `div.col` needed) |
| `columns-3` | Three equal columns |
| `split` | 1fr 2fr grid: use `<div class="col">` for each panel |
| `timeline` | Horizontal timeline axis (for `::: timeline` block) |

Combine with `_gap`, `_align`, or `_columns`:

```markdown
<!-- _layout: columns-2 -->
<!-- _gap: loose -->

## Compare & Contrast
Left column content flows here.

Right column continues automatically.
```

# Examples

## Example: Split Layout

```markdown
<!-- _layout: split -->

<div class="col">

## Title or visual
Smaller area on the left.

</div>

<div class="col">

### Main content
Larger area on the right for detailed content.

</div>
```

### Example: Lead Slide

```markdown
<!-- _class: lead -->

# Welcome to Dorcs
Build docs and slides in one place.
```

### Example: Big Text (stat or statement)

```markdown
<!-- _layout: big -->

35%

of an audience's retention rate is attributed to the visuals used.
```

### Example: Quote

```markdown
<!-- _layout: quote -->

They may forget what you said, but they will never forget how you made them feel.

– Carl W. Buechner
```

### Example: Left-Aligned Slide

```markdown
<!-- _class: left -->

## Features
- Single binary
- Markdown slides
- No external deps
```

### Example: Custom Background

```markdown
<!-- _class: lead -->
<!-- _backgroundColor: #16213e -->

# Dark Theme Slide
Custom background color.
```

### Example: Background Image + Gradient

```markdown
<!-- _backgroundImage: "linear-gradient(to bottom, #67b8e3, #0288d1)" -->
<!-- _color: white -->

# Gradient Background
White text on blue gradient.

---

<!-- _backgroundColor: black -->
<!-- _color: white -->

# Black + White
High-contrast slide.
```

### Example: Hide Slide Number

```markdown
<!-- _paginate: false -->

# Title Slide
No slide number shown.
```

### Example: Two Columns

```markdown
<!-- _class: two-columns -->

<div class="col">

### Left Column
- Item 1
- Item 2

</div>

<div class="col">

### Right Column
- Item A
- Item B

</div>
```

### Example: Auto Columns (no div.col)

```markdown
<!-- _layout: columns-2 -->

### Left column content
Paragraphs and lists flow automatically into columns.

### Right column content
No manual wrapping needed.
```
