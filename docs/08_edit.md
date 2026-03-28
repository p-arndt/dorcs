---
title: "Edit Mode"
description: "Edit your docs right in the browser."
tags: [edit, auth]
date: 2026-03-18
draft: false
---

# Edit Mode

Dorcs includes a built-in browser editor so you can make quick changes without opening a code editor.

## Set it up

Add this to your `dorcs.yaml`:

```yaml
auth:
  enabled: true
  username: "admin"
  password: "change-me"
```

> [!TIP]
> Dorcs automatically hashes your password on first run, so it's never stored in plain text after that.

## What you can do

Once logged in, you can:

- Browse all files in the docs folder
- Edit and save existing pages
- Create new pages
- Delete pages

> [!NOTE]
> Hidden files are never shown, and you can only access files inside the `docs/` directory.

## Session storage

Sessions are saved in `docs/.dorcs_sessions.json` by default. You can change the path:

```yaml
auth:
  sessions_path: "/custom/path/sessions.json"
```

## Security

> [!WARNING]
> Edit mode is designed as an **internal authoring tool**, not for public-facing anonymous editing. If you're exposing Dorcs publicly:
> - Use strong credentials
> - Put Dorcs behind HTTPS
> - Consider restricting access at the network level

## When to use something else

If your docs live in a Git repo and you want changes to go through pull requests, skip edit mode and use the [Edit on GitHub](./external-content/github.md) integration instead. That gives each page an "Edit on GitHub" link that takes readers straight to the file in your repo.
