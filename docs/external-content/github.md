---
title: "GitHub"
description: "Serve docs from GitHub or add Edit on GitHub links."
tags: [github]
date: 2026-03-18
draft: false
---

# GitHub

Dorcs supports two GitHub-related modes.

## 1. GitHub as the content source

```yaml
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: ${GITHUB_TOKEN}
  cache_ttl: "1h"
```

Behavior:

- Dorcs fetches Markdown from the repo path
- if the branch is omitted, Dorcs resolves the repository default branch
- local Markdown content is ignored
- cache is stored under `.cache/github` when possible

Supported repository formats:

- `https://github.com/owner/repo/tree/branch/path`
- `https://github.com/owner/repo`
- `github.com/owner/repo/tree/branch/path`

## 2. Edit on GitHub for local docs

```yaml
github:
  edit_on_github:
    repository: "https://github.com/owner/repo/tree/main/docs"
```

Use this when content is local at runtime but stored in GitHub as the source repository. Dorcs adds an "Edit on GitHub" link per page.

## Environment variables

`token` supports:

- `${GITHUB_TOKEN}`
- `${GITHUB_TOKEN:-fallback}`

Dorcs loads `.env` automatically when present.

## Good fit

Choose GitHub-backed docs when the repository already defines the docs tree and the running environment should stay read-only.
