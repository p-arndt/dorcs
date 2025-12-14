---
title: "Watch Mode"
description: "How to use watch mode"
tags: [watch-mode, development]
date: 2025-12-14
draft: false
---

# Watch Mode

You can use Dorcs with watch mode for development. This means that you can see your changes immediately in the browser without having to restart the server.

It is as easy as adding the `--watch` flag to the command when starting the server.

```bash
./dorcs --watch
```

> [!IMPORTANT]
> File watching through Docker volumes may not work reliably on Windows/Mac. For best results, run dorcs directly on your host: `./dorcs --dir ./docs --watch`

## What does watch mode do?

Watch mode watches the `docs` directory for changes and automatically reloads the page. 
It also watches the `dorcs.yaml` file for changes. So if you change the configuration, the server will automatically reload.


