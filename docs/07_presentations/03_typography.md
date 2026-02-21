---
title: "Typography"
description: "Typography blocks for slide decks"
tags: [presentations, slides, markdown]
date: 2026-02-21
draft: false
after: "layout"
---

# Typography Blocks

Composable blocks for different font sizes. Mix them on the same slide.

| Block         | Purpose                              |
| ------------- | ------------------------------------ |
| `::: hero`    | Large impactful text                 |
| `::: stat`    | Number (first line) + caption (rest) |
| `::: caption` | Small supporting text                |
| `::: label`   | Small uppercase label                |

```markdown
::: hero
35%
:::

::: caption
of an audience's retention rate is attributed to the visuals used.
:::
```

```markdown
::: stat
35%
of an audience's retention rate is attributed to the visuals used.
:::
```

```markdown
::: label
AGENDA
:::

1. Introducing yourself
2. Engaging the audience
```
