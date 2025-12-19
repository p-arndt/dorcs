---
title: "Dorcs"
description: "Welcome to dorcs - a single-binary static documentation server for Markdown files."
tags: [docs, markdown]
date: 2025-12-13
draft: false
---

# Dorcs

<div style="display: flex; justify-content: center; align-items: center;">
<img src="../logo.png" alt="Dorcs Logo" width="200" height="200" style="border: none;" />
</div>

Dorcs is a single-binary static documentation server for Markdown files. It is a simple and easy to use documentation server that allows you to create and host your documentation site in minutes.

## Get Started

Ready to get started? Check out the [Getting Started guide](./01_getting-started.md) to have your documentation site running in under 5 minutes.

Or if you ready to get started now, download the latest release for your platform:

<div style="display: flex; flex-wrap: wrap; justify-content: center; gap: 1.5em; margin: 2em 0;">
  <a href="https://github.com/p-arndt/dorcs/releases/latest/download/dorcs-windows-amd64.exe"
     style="
       display: inline-flex;
       align-items: center;
       background: linear-gradient(90deg, #3292df 0%, #0078d7 100%);
       color: #fff;
       font-weight: 600;
       font-size: 1.15em;
       padding: 0.8em 2.2em;
       border: none;
       border-radius: 50px;
       text-decoration: none;
       box-shadow: 0 2px 12px 0 rgba(50,146,223,0.12);
       transition: background 0.3s, box-shadow 0.3s;
       margin: 0.5em 0;
       width: 260px;
       justify-content: center;
     "
     download
  >
    <svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24" style="margin-right: 0.65em; flex: none;">
      <rect fill="#F25022" x="1" y="1" width="10" height="10"></rect>
      <rect fill="#7FBA00" x="13" y="1" width="10" height="10"></rect>
      <rect fill="#00A4EF" x="1" y="13" width="10" height="10"></rect>
      <rect fill="#FFB900" x="13" y="13" width="10" height="10"></rect>
    </svg>
    Windows
  </a>
  <a href="https://github.com/p-arndt/dorcs/releases/latest/download/dorcs-linux-amd64.exe"
     style="
       display: inline-flex;
       align-items: center;
       background: linear-gradient(90deg, #43D46A 0%, #1E8449 100%);
       color: #fff;
       font-weight: 600;
       font-size: 1.15em;
       padding: 0.8em 2.2em;
       border: none;
       border-radius: 50px;
       text-decoration: none;
       box-shadow: 0 2px 12px 0 rgba(67,212,106,0.12);
       transition: background 0.3s, box-shadow 0.3s;
       margin: 0.5em 0;
       width: 260px;
       justify-content: center;
     "
     download
  >
    <svg xmlns="http://www.w3.org/2000/svg" height="24" width="24" viewBox="0 0 24 24" style="margin-right: 0.65em; flex: none;">
      <g>
        <rect fill="#43D46A" x="2" y="2" width="20" height="20" rx="4"></rect>
        <path fill="#fff" d="M12 6c.512 0 .936.386.993.883l.007.117v4.586l1.793-1.793a1 1 0 0 1 1.497 1.32l-.083.094-3.5 3.5a1 1 0 0 1-1.32.083l-.094-.083-3.5-3.5a1 1 0 0 1 1.32-1.497l.094.083L11 11.586V7c0-.552.448-1 1-1z"></path>
      </g>
    </svg>
    Linux
  </a>
</div>


## Documentation

- 🚀 [Getting Started](./01_getting-started.md) - Complete quick start guide
- 📦 [Installation](./02_installation.md) - Detailed installation instructions
- ⚙️ [Configuration](./03_configuration.md) - Customize your site with `dorcs.yaml`
- 📁 [File Structure](./usage/file-structure.md) - Organize docs with languages and versions
- 🔗 [External Content](./external-content/index.md) - Serve content from external sources (GitHub, etc.)
- 🚢 [Deployment](./04_deployment.md) - Deploy to production
- 🎨 [Themes](./05_themes.md) - Browse all available themes
- 📝 [Markdown Features](./06_markdown/index.md) - Complete guide to markdown features

## Features

- **Single binary** – no runtime dependencies, statically linkable
- **Extensionless URLs** – `/guide/getting-started` serves `docs/guide/getting-started.md`
- **External Content Sources** – serve markdown files directly from GitHub repositories with automatic caching
- **Multi-lingual support** – serve documentation in multiple languages with automatic language switching
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

**Multi-lingual URLs:**

When multiple languages are configured, non-default languages use a language prefix:

| File Path                       | URL                      |
| ------------------------------- | ------------------------ |
| `docs/en/index.md`              | `/` or `/en/`            |
| `docs/de/index.md`              | `/de/`                   |
| `docs/de/getting-started.md`    | `/de/getting-started`    |
| `docs/fr/guide/installation.md` | `/fr/guide/installation` |

The default language is served at the root URL (no prefix), while other languages use `/{lang}/` prefixes.

**Versioned URLs:**

When versioning is configured, non-default versions use a version prefix:

| File Path                       | URL                      |
| ------------------------------- | ------------------------ |
| `docs/v1/index.md`              | `/v1/`                   |
| `docs/v1/getting-started.md`    | `/v1/getting-started`    |
| `docs/en/v1/getting-started.md` | `/en/v1/getting-started` |

**Combined (Languages + Versions):**

When both are configured, URLs use language-first structure:

| File Path                       | URL                      |
| ------------------------------- | ------------------------ |
| `docs/en/v1/getting-started.md` | `/en/v1/getting-started` |
| `docs/de/v1/getting-started.md` | `/de/v1/getting-started` |

See [File Structure & Organization](./usage/file-structure.md) for complete details.

It will also automatically build navigation from your structure and generate a sidebar and a table of contents for each page.
