---
title: "Getting Started"
description: "Create your first Dorcs site in a few minutes."
tags: [getting-started, quickstart]
date: 2026-03-18
draft: false
---

# Getting Started

This guide gets a local docs site running with the fewest steps.

## 1. Create a docs site

```bash
dorcs init
```

This creates:

- `docs/index.md`
- a starter `dorcs.yaml` in your working directory if one does not already exist

## 2. Start the server

```bash
dorcs --watch
```

Open `http://localhost:8080`.

`--watch` rebuilds navigation and reloads the browser when Markdown or config files change.

## 3. Edit content

Start with `docs/index.md`:

```markdown
# Welcome

Your docs live in plain Markdown files.
```

Add a second page:

```text
docs/
├── index.md
└── guide.md
```

`docs/guide.md` becomes `/guide`.

## 4. Build static output

```bash
dorcs build
```

This writes a deployable site to `dist/` by default.

## 5. Know the routing model

| File | URL |
| --- | --- |
| `docs/index.md` | `/` |
| `docs/getting-started.md` | `/getting-started` |
| `docs/guide/index.md` | `/guide` |
| `docs/guide/install.md` | `/guide/install` |

## Next steps

- Use [Installation](./02_installation.md) if you still need the binary
- Use [Configuration](./03_configuration.md) to set title, theme, and navigation
- Use [Writing Your Docs](./usage/writing-your-docs.md) for structure and linking guidance
