---
title: "Edit"
description: "Browser-based editing in Dorcs."
tags: [edit, auth]
date: 2026-03-18
draft: false
---

# Edit

Dorcs can expose a browser-based editing API for local docs content.

## Enable edit mode

```yaml
auth:
  enabled: true
  username: "admin"
  password: "change-me"
```

On first run, Dorcs hashes the password and clears the plain-text value in memory.

## What edit mode does

Authenticated users can:

- browse the docs directory
- read files
- save changes
- create files
- delete files

Hidden files are not shown, and file access is constrained to the docs directory.

## Session storage

By default, sessions are stored in:

```text
docs/.dorcs_sessions.json
```

Override it if needed:

```yaml
auth:
  sessions_path: ".dorcs_sessions.json"
```

## Security notes

- Use strong credentials
- Put Dorcs behind HTTPS in shared or public environments
- Treat edit mode as an internal authoring feature, not anonymous public editing

## When not to use it

If your source of truth is Git and edits should always happen through pull requests, prefer local editing plus [Edit on GitHub](./external-content/github.md).
