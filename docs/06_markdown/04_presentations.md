---
title: "Presentations"
description: "Create slide decks from Markdown with layout and styling options"
tags: [markdown, presentation, slides]
date: 2025-02-21
draft: false
---

# Presentations

Dorcs can render Markdown files as full-screen slide decks, similar to Marp but built-in. Use `---` to separate slides and add layout directives for styling.

## Quick Start

Add `presentation: true` to your front matter:

```markdown
---
title: "My Talk"
presentation: true
---

# Welcome
First slide content.

---

## Second Slide
More content here.
```

## Slide Directives

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

### Spot vs inherited directives

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

### Layout Classes

- **`lead`** – Title slide: larger headings, more dramatic typography
- **`left`** – Left-aligned content
- **`right`** – Right-aligned content
- **`fit`** – Compact layout, fits more content
- **`invert`** – Inverted colors (dark background, light text)
- **`two-columns`** – Two-column layout (use `<div class="col">` for each column)

### Layout Control

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

### Timeline Block

Use the `::: timeline` block for roadmaps and process steps. Each `###` or `####` heading starts a new step; the heading text becomes the date/label marker.

```markdown
<!-- _layout: timeline -->

# Project Roadmap

::: timeline
### 2024 · Q1
**Discovery**
Initial research and planning.

### 2024 · Q2
**Build**
Development and testing.

### 2024 · Q3
**Launch**
Release and iterate.
:::
```

## Front Matter Options

| Field | Description |
|-------|-------------|
| `presentation: true` | Enable slide deck mode |
| `presentation_header` | Header text shown on every slide |
| `presentation_footer` | Footer text shown on every slide |

```yaml
---
title: "My Presentation"
presentation: true
presentation_header: "Conference 2025"
presentation_footer: "© My Name"
---
```

## Keyboard Navigation

- **→ / ↓ / Space** – Next slide
- **← / ↑** – Previous slide
- **Home** – First slide
- **End** – Last slide
- **F** – Toggle fullscreen (presentation mode)

## URL Hash

Append `#2` to link to a specific slide, e.g. `/talks/my-talk#2` opens on slide 2.
