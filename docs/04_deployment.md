---
title: "Deployment"
description: "Build and deploy a static Dorcs site."
tags: [deployment]
date: 2026-03-18
draft: false
---

# Deployment

For production, Dorcs is usually used in static build mode.

## Build the site

```bash
dorcs build --dir ./docs --output ./dist
```

Common flags:

- `--base-url /docs` for subpath deployments
- `--config ./dorcs.yaml` for an explicit config file
- `--theme` and `--theme-mode` to override config during build

## Deploy the output

Upload the contents of `dist/` to any static host:

- GitHub Pages
- Netlify
- Vercel
- S3 + CDN
- Nginx or Caddy

## Subpath deployments

If the site is served below the domain root, set the base path during build:

```bash
dorcs build --base-url /docs
```

That ensures internal links, assets, and APIs are generated with the correct prefix.
That ensures internal links and assets are generated with the correct prefix.

## GitHub Pages example

```bash
dorcs build --output ./dist --base-url /my-project
```

Publish `dist/` to the `gh-pages` branch or your Pages artifact workflow.

## Reverse proxy setup

If you run Dorcs as a live server behind a proxy instead of exporting static files:

- forward `X-Forwarded-Proto` when terminating TLS upstream
- keep the configured `--base-url` consistent with the mounted path

## Production checklist

- build with the right `--base-url`
- verify logo and favicon paths
- disable drafts if you do not want preview content published
- check external links and "Edit on GitHub" targets
- note that static builds do not include the live `/api/search` endpoint
