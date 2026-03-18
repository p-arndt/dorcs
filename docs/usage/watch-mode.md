---
title: "Watch Mode"
description: "Live reload and file watching in Dorcs."
tags: [usage, watch]
date: 2026-03-18
draft: false
---

# Watch Mode

Watch mode is the fastest way to write docs locally.

## Start it

```bash
dorcs --watch
```

> [!IMPORTANT]
> File watching through Docker volumes may not work reliably on Windows/Mac. For best results, run dorcs directly on your host: `./dorcs --dir ./docs --watch`

## What does watch mode do?

Watch mode watches the `docs` directory for changes and automatically reloads the page. 
It also watches the `dorcs.yaml` file for changes. So if you change the configuration, the server will automatically reload.

## Limits

When GitHub content mode is enabled, file watching is skipped because content is not sourced from the local docs directory.


