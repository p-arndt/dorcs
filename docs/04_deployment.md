---
title: "Deployment"
description: "Build your docs and deploy them anywhere."
tags: [deployment]
date: 2026-03-18
draft: false
---

# Deployment

When your docs are ready, Dorcs can export them as plain HTML files you can host anywhere.

## Build your site

```bash
dorcs build --dir ./docs --output ./dist
```

This creates a `dist/` folder with everything ready to upload.

> [!TIP]
> You can also override the theme or config during build: `dorcs build --theme ocean --config ./dorcs.yaml`

## Where to deploy

The `dist/` folder works with any static hosting provider:

:::tabs
::tab GitHub Pages
```bash
dorcs build --output ./dist --base-url /my-project
```

Then publish the `dist/` folder to the `gh-pages` branch or use a GitHub Actions workflow.

> [!NOTE]
> The `--base-url` flag is important for GitHub Pages since your site lives under a subpath like `/my-project`.
::tab Netlify / Vercel
```bash
dorcs build --output ./dist
```

Point your build command to `dorcs build` and the publish directory to `dist/`. No base URL needed — these platforms serve from the root.
::tab Nginx / Caddy
```bash
dorcs build --output ./dist
```

Serve the `dist/` folder as static files. If you're running Dorcs as a live server behind a reverse proxy instead, make sure to:

- Forward `X-Forwarded-Proto` when terminating TLS upstream
- Keep `--base-url` consistent with the mounted path
::tab S3 + CDN
```bash
dorcs build --output ./dist
```

Upload the contents of `dist/` to your S3 bucket and put a CDN like CloudFront in front.
:::

## Subpath deployments

If your docs live under a subpath (like `example.com/docs` instead of the root), tell Dorcs about it:

```bash
dorcs build --base-url /docs
```

This ensures all internal links and assets use the correct prefix.

## Before you ship

> [!IMPORTANT]
> Quick checklist before deploying:
> - Build with the correct `--base-url` if needed
> - Check that your logo and favicon paths work
> - Drafts are hidden by default — that's usually what you want
> - Static builds include a `sitemap.xml` but **not** the live search API (search only works in server mode)
