---
title: "Dorcs"
description: "Welcome to dorcs - a single-binary static documentation server for Markdown files."
tags: [docs, markdown]
date: 2025-12-13
draft: false
---


# Dorcs

<div style="display: flex; justify-content: center; align-items: center;">
<img src="./logo.png" alt="Dorcs Logo" width="200" height="200" style="border: none;" />
</div>

Dorcs is a single-binary static documentation server for Markdown files. It is a simple and easy to use documentation server that allows you to create and host your documentation site in minutes.

## Get Started

Ready to get started? Check out the [Getting Started guide](./01_getting-started.md) to have your documentation site running in under 5 minutes.

## Documentation

- 🚀 [Getting Started](./01_getting-started.md) - Complete quick start guide
- 📦 [Installation](./02_installation.md) - Detailed installation instructions
- ⚙️ [Configuration](./03_configuration.md) - Customize your site with `dorcs.yaml`
- 🚢 [Deployment](./04_deployment.md) - Deploy to production
- 🎨 [Themes](./05_themes.md) - Browse all available themes
- 📝 [Markdown Features](./06_markdown/index.md) - Complete guide to markdown features

## Features

- **Single binary** – no runtime dependencies, statically linkable
- **Extensionless URLs** – `/guide/getting-started` serves `docs/guide/getting-started.md`
- **YAML front matter** – metadata support (title, description, date, tags, draft)
- **Table of Contents** – auto-generated from headings with scrollspy
- **Navigation sidebar** – built automatically from your folder structure
- **Responsive design** – mobile-friendly with collapsible sidebar
- **Dark mode** – automatic based on system preference
- **Live reload** – watch mode for development with smart content updates
- **Multiple themes** – choose from 20+ built-in themes
- **Search** – built-in search functionality
- **Edit Mode** – online editing with authentication (create, edit, delete files directly in the browser)

## How It Works

### URL Routing

dorcs uses extensionless URLs that map directly to your file structure:

| File Path                    | URL                   |
| ---------------------------- | --------------------- |
| `docs/index.md`              | `/`                   |
| `docs/getting-started.md`    | `/getting-started`    |
| `docs/guide/index.md`        | `/guide`              |
| `docs/guide/installation.md` | `/guide/installation` |

It will also automatically build navigation from your structure and generate a sidebar and a table of contents for each page.
