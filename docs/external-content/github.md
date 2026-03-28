---
title: "GitHub Integration"
description: "Serve docs from GitHub or add Edit on GitHub links."
tags: [github]
date: 2026-03-18
draft: false
---

# GitHub Integration

Dorcs has two GitHub-related features. You can use either one or both together.

## Serve docs from GitHub {badge:NEW}

Pull your Markdown content directly from a GitHub repository:

```yaml
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: ${GITHUB_TOKEN}
  cache_ttl: "1h"
```

When this is active:

- Dorcs fetches Markdown files from the specified repo and path
- Local `docs/` content is ignored
- Content is cached under `.cache/github` for fast reloads
- If you leave out the branch, Dorcs uses the repo's default branch

**Supported URL formats:**

- `https://github.com/owner/repo/tree/branch/path`
- `https://github.com/owner/repo` (root of default branch)
- `github.com/owner/repo/tree/branch/path` (without https)

You can also skip config entirely and bootstrap from the command line:

```bash
dorcs --repo https://github.com/owner/repo/tree/main/docs
```

This loads both the docs and the config file from that repository.

## "Edit on GitHub" links

If your docs are served locally but the source lives on GitHub, add edit links to every page:

```yaml
github:
  edit_on_github:
    repository: "https://github.com/owner/repo/tree/main/docs"
```

Each page will show an "Edit on GitHub" link that takes the reader directly to the file in your repo.

> [!TIP]
> This is great when you want the speed of local serving but want contributors to use GitHub's pull request workflow for edits.

## Environment variables

The `token` field supports variable expansion:

- `${GITHUB_TOKEN}` — reads from the environment
- `${GITHUB_TOKEN:-my-fallback}` — uses a fallback if the variable isn't set

Dorcs automatically loads a `.env` file from your project root if one exists.

## When to use it

GitHub-backed docs are a good fit when:

- Your repo already defines the docs structure
- The runtime environment should stay read-only
- You want docs to update automatically when the repo changes (with cache TTL)
