---
title: "External Content"
description: "Serve your docs directly from a GitHub repository."
tags: [github, external-content]
date: 2026-03-18
draft: false
---

# External Content

Instead of reading Markdown from your local `docs/` folder, Dorcs can pull content straight from GitHub.

## Why would you want this?

- **Deploy the binary without bundling docs** — keep your docs in the repo, serve them at runtime
- **Read docs from another repo** — great for monorepos or shared documentation
- **Keep your site tied to a branch** — always serve what's on `main`

## Quick setup

```yaml
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: ${GITHUB_TOKEN}
  cache_ttl: "1h"
```

> [!TIP]
> You only need a token for private repos or to avoid GitHub's rate limits. For public repos, you can leave it out.

When this is enabled, Dorcs ignores local Markdown files and serves everything from GitHub. Content is cached so it's fast after the first load.

## Learn more

- [GitHub Integration](./github.md) — full details on both GitHub modes
- [Configuration](../03_configuration.md) — the `github` config section
