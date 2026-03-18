---
title: "Sample Presentation"
presentation: true
presentation_header: "Dorcs · Sample Presentation"
presentation_footer: "Sample Presentation built with Dorcs"
---

<!-- _layout: lead invert -->
<!-- _paginate: false -->

# Dorcs Presentations

Markdown-powered slide decks built in.

---

<!-- _layout: left -->

## Layout: Left

`<!-- _layout: left -->`

- Content aligned to the left
- Ideal for bullet lists
- Code snippets and dense content

---

<!-- _layout: right -->

## Layout: Right

`<!-- _layout: right -->`

Content aligned to the right.

Good for conclusions and emphasis.

---

<!-- _layout: cols -->

::: col

### Two Columns

`<!-- _layout: cols -->`

Use `::: col` blocks for each column.

:::

::: col

### Explicit Columns

- Full control over column content
- Each column is a separate block
- Cleaner than raw HTML divs

:::

---

<!-- _layout: columns-2 -->

## Auto Columns — `columns-2`

### Syntax

Use `<!-- _layout: columns-2 -->` and content flows automatically. No `::: col` needed.

### When to use

Best for lists, short paragraphs, or sections that distribute naturally across two columns.

---

<!-- _layout: columns-3 -->
<!-- _gap: loose -->

## Three Columns

**Features** — Single binary, Markdown, no external deps.

**Layout options** — lead, left, right, columns, split, timeline.

**Styling** — \_layout, \_gap, \_align, \_class.

---

<!-- _layout: big -->

35%

of an audience's retention rate is attributed to the visuals used.

---

<!-- _layout: quote -->

They may forget what you said, but they will never forget how you made them feel.

Carl W. Buechner

---

## Typography Blocks

Mix `::: hero`, `::: stat`, `::: caption`, `::: label` on one slide.

::: hero
35%
:::

::: caption
of an audience's retention rate is attributed to the visuals used.
:::

---

<!-- _layout: split -->
<!-- _align: start -->

## Split Layout

`<!-- _layout: split -->` — use `::: col` for each panel.

::: col

### Left Panel (1fr)

Title or visual placeholder. Smaller area.

:::

::: col

### Right Panel (2fr)

Main content goes here. Larger area for details, lists, or paragraphs.

:::

---

<!-- _layout: timeline -->

# Timeline Block

`<!-- _layout: timeline -->` with `::: timeline`

::: timeline

### Q1 · Discovery

**Research & Planning**
Initial phase and requirements.

### Q2 · Build

**Development**
Implementation and testing.

### Q3 · Launch

**Release**
Deploy and iterate.
:::

---

<!-- _layout: fit -->

## Fit Layout

`<!-- _layout: fit -->` — compact sizing

- Smaller headings and text
- Fits more content per slide
- Great for appendix slides or code-heavy content

---

<!-- _layout: columns-2 -->
<!-- _align: start -->

## Layout Primitives

| Directive | Effect                                                                           |
| --------- | -------------------------------------------------------------------------------- |
| `_layout` | Structure: lead, left, right, big, quote, columns-2, columns-3, split, cols, timeline |
| Blocks    | `::: col`, `::: hero`, `::: stat`, `::: caption`, `::: label`, `::: timeline`   |
| `_gap`    | Spacing: tight, normal, loose                                                    |
| `_align`  | Alignment: start, center, end                                                    |

---

<!-- _header: Custom Header -->
<!-- _footer: Custom Footer -->
<!-- _paginate: false -->

## Header & Footer

`<!-- _header: ... -->` and `<!-- _footer: ... -->`

Override presentation-wide header and footer for a single slide.

`<!-- _paginate: false -->` hides the slide number.

---

<!-- _paginate: false -->

# Thank You

**Dorcs** — Static docs and slides in one place.

→ / Space to navigate · ← to go back · F for fullscreen
