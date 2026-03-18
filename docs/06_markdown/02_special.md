---
title: "Special Markdown"
description: "Documentation-oriented Markdown patterns."
tags: [markdown]
date: 2026-03-18
draft: false
---

# Special Markdown

These patterns are especially useful in docs content.

## Front matter

Put YAML front matter at the top of the file:

```yaml
---
title: "API Overview"
description: "What the API covers"
draft: false
---
```

## Relative links

Write links to local Markdown files naturally:

```markdown
[Auth guide](./auth.md)
```

Dorcs serves the target as a clean extensionless URL.

## GitHub-style callouts

Dorcs supports GitHub alert syntax:

```markdown
> [!NOTE]
> This is a note.

> [!TIP]
> This is a tip.
```

Supported types:

- `NOTE`
- `TIP`
- `IMPORTANT`
- `WARNING`
- `CAUTION`

## Images

```markdown
![Diagram](./images/diagram.png)
```

Keep image paths relative to the current file unless you intentionally want shared root-based assets.
