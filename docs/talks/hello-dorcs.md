---
title: "Dorcs Presentations"
presentation: true
presentation_header: "Dorcs · Feature Overview"
presentation_footer: "Build docs and slides in one place"
---

<!-- _class: lead -->
<!-- _paginate: false -->

# Dorcs Presentations

Markdown-powered slide decks built in.

---

<!-- _class: left -->

## Layout: Left

`<!-- _class: left -->`

- Content aligned to the left
- Ideal for bullet lists
- Code snippets and dense content

---

<!-- _class: right -->

## Layout: Right

`<!-- _class: right -->`

Content aligned to the right.

Good for conclusions and emphasis.

---

<!-- _class: two-columns -->

<div class="col">

### Two Columns

`<!-- _class: two-columns -->`

Use `<div class="col">` for each column.

</div>

<div class="col">

### Explicit Columns

- Full control over column content
- Each column is a separate block
- Marpit-compatible

</div>

---

<!-- _layout: columns-2 -->

## Auto Columns

`<!-- _layout: columns-2 -->`

Content flows automatically into columns. No `div.col` needed.

### Left Column

Paragraphs and lists distribute naturally.

### Right Column

Headings and content flow into the second column.

---

<!-- _layout: columns-3 -->
<!-- _gap: loose -->

## Three Columns

**Features** — Single binary, Markdown, no external deps.

**Layout options** — lead, left, right, columns, split, timeline.

**Styling** — _class, _layout, _gap, _align, _columns.

---

<!-- _layout: big -->

35%

of an audience's retention rate is attributed to the visuals used.

---

<!-- _layout: quote -->

They may forget what you said, but they will never forget how you made them feel.

– Carl W. Buechner

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

`<!-- _layout: split -->` — use `<div class="col">` for each panel.

<div class="col">

### Left Panel (1fr)

Title or visual placeholder. Smaller area.

</div>

<div class="col">

### Right Panel (2fr)

Main content goes here. Larger area for details, lists, or paragraphs.

</div>

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

<!-- _class: fit -->

## Fit Layout

`<!-- _class: fit -->` — compact sizing

- Smaller headings and text
- Fits more content per slide
- Great for appendix slides or code-heavy content

---

<!-- _layout: columns-2 -->
<!-- _align: start -->

## Layout Primitives

| Directive | Effect |
|-----------|--------|
| `_layout` | Structure: lead, left, right, big, quote, columns-2, columns-3, split, timeline |
| Blocks | `::: hero`, `::: stat`, `::: caption`, `::: label`, `::: timeline` |
| `_gap` | Spacing: tight, normal, loose |
| `_align` | Alignment: start, center, end |
| `_columns` | Column count: 2, 3, 4 |

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
