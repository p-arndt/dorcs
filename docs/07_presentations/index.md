---
title: "Presentations"
description: "Create slide decks from Markdown with layout and styling options"
tags: [presentations, slides, markdown]
date: 2026-02-21
draft: false
after: "06_markdown"
---

# Presentations

Dorcs can render Markdown files as full-screen slide decks, similar to Marp but built-in. Use `---` to separate slides and add layout directives for styling.

> [!NOTE]
> Checkout the [Sample Presentation](./05_sample_presentation.md) for a working example.

## Quick Start

Add `presentation: true` to your front matter:

```markdown
---
title: "My Talk"
presentation: true
---

# Welcome

First slide content.
```

## Front Matter Options

| Field                 | Description                      |
| --------------------- | -------------------------------- |
| `presentation: true`  | Enable slide deck mode           |
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

## See Also

- [Slide Directives](./01_directives.md)
- [Layout](./02_layout.md)
