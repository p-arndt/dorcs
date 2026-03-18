---
title: "Typography"
description: "Simple typography rules for better Dorcs presentations."
tags: [presentations, typography]
date: 2026-03-18
draft: false
---

# Typography

Presentation slides need less text and stronger hierarchy than normal docs pages.

## Good defaults

- one idea per slide
- short headings
- three to five bullets at most
- large numbers or statements on `big` slides

## Good slide text

```markdown
<!-- _layout: big -->

42%

faster onboarding after consolidating product docs
```

## Avoid

- long paragraphs
- nested lists
- dense code samples on title or summary slides

## Typography blocks

Dorcs also supports small presentation-oriented block helpers:

- `::: hero`
- `::: stat`
- `::: caption`
- `::: label`

Example:

```markdown
::: hero
35%
:::

::: caption
faster onboarding after restructuring the docs
:::
```

## Use docs for detail

Keep the full explanation in a normal Dorcs page and use the deck to present the main story.
