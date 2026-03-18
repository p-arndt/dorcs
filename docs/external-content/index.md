---
title: "External Content"
description: "Use GitHub as a content source for Dorcs."
tags: [github, external-content]
date: 2026-03-18
draft: false
---

# External Content

Dorcs can serve Markdown directly from GitHub instead of reading local docs files.

## Why use it

- deploy the binary without bundling docs
- read documentation from another repository
- keep the site tied to a branch in GitHub

## How it works

When GitHub integration is enabled:

- Dorcs discovers Markdown files from the configured repository path
- local Markdown content is ignored
- fetched content is cached

## Start here

```yaml
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: ${GITHUB_TOKEN}
  cache_ttl: "1h"
```

Use a token only when needed, especially for private repositories or higher API limits.

## Related pages

- [GitHub](./github.md)
- [Configuration](../03_configuration.md)
