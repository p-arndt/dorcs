---
title: "File Structure & Organization"
description: "Complete guide to organizing your documentation files with languages and versions"
tags: [structure, organization, languages, versions]
date: 2025-12-14
draft: false
---

# File Structure & Organization

dorcs uses a **MkDocs-style structure** that minimizes boilerplate and requires explicit configuration for languages and versions. This guide explains how to structure your files for different scenarios.

## Core Principles

1. **Minimal boilerplate**: If you don't need languages or versions, just put files in `docs/`
2. **Explicit configuration**: Languages and versions must be explicitly configured in `dorcs.yaml`
3. **Language-first**: When using both languages and versions, use `docs/{lang}/{version}/` structure
4. **Default behavior**: Default language/version is served at root URL (no prefix)

## Basic Structure (No Languages, No Versions)

For simple, single-language documentation, just place your files directly in `docs/`:

```
docs/
  index.md
  getting-started.md
  guide/
    installation.md
    advanced.md
```

**URLs:**
- `/` → `docs/index.md`
- `/getting-started` → `docs/getting-started.md`
- `/guide/installation` → `docs/guide/installation.md`

## Multi-Language Structure

When you need multiple languages, **define the default language in `dorcs.yaml`** and place all files in language-specific folders:

```
docs/
  en/                         # English (default) - served at /
    index.md
    getting-started.md
    guide/
      installation.md
  de/                         # German - served at /de/
    index.md
    getting-started.md
    guide/
      installation.md
  fr/                         # French - served at /fr/
    index.md
    getting-started.md
    guide/
      installation.md
```

**Configuration in `dorcs.yaml`:**

```yaml
languages:
  default: "en"               # Define default language here
  enabled:                    # Explicit configuration (required)
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
    - code: "fr"
      name: "Français"
```

**Important Notes:**

- **Markdown files**: When using multiple languages, markdown files in root `docs/` are **ignored**. All markdown files must be in language folders (`docs/en/`, `docs/de/`, etc.)
- **Shared assets**: Static assets (images, PDFs, etc.) in root `docs/` are served as **language-independent shared assets** accessible from all languages
- **Default language**: Defined in `dorcs.yaml` - default language files are in `docs/en/` but served at `/` (no prefix)
- **Other languages**: Served at `/{lang}/` with their language prefix
- **Language codes**: Use ISO 639-1 codes (e.g., `en`, `de`, `fr`, `zh-CN`, `pt-BR`)
- **URL structure**: 
  - Default (English): `/getting-started` → `docs/en/getting-started.md`
  - German: `/de/getting-started` → `docs/de/getting-started.md`
  - French: `/fr/getting-started` → `docs/fr/getting-started.md`

> [!WARNING]
> **Markdown files in root are ignored**: When using multiple languages, any `.md` files in the root `docs/` folder are ignored. Only static assets (images, PDFs, etc.) in root are served as shared assets accessible from all languages.

## Versioning Structure

For versioned documentation, create direct version folders:

```
docs/
  index.md                    # Default version (e.g., latest) - served at /
  getting-started.md
  latest/                     # Latest version (explicit) - served at /latest/
    index.md
    getting-started.md
  v1/                         # Version 1 - served at /v1/
    index.md
    getting-started.md
  v2/                         # Version 2 - served at /v2/
    index.md
    getting-started.md
```

**Version identifiers:**
- Can be any string: `v1`, `v2`, `1.0`, `2.0.0`, `latest`, `stable`, etc.
- Default version is served at root URL (no prefix)
- Other versions use `/{version}/` prefix

**URLs:**
- `/getting-started` → `docs/getting-started.md` (default version)
- `/v1/getting-started` → `docs/v1/getting-started.md`
- `/v2/getting-started` → `docs/v2/getting-started.md`

## Combined: Languages + Versions

When using both languages and versions, use **language-first structure** and define defaults in `dorcs.yaml`:

```
docs/
  en/                         # English (default language)
    index.md                  # en, latest (default) - served at /
    getting-started.md
    latest/                   # en, latest (explicit) - served at /en/latest/
      index.md
    v1/                       # en, v1 - served at /en/v1/
      index.md
      getting-started.md
    v2/                       # en, v2 - served at /en/v2/
      index.md
  de/                         # German
    index.md                  # de, latest (default) - served at /de/
    getting-started.md
    v1/                       # de, v1 - served at /de/v1/
      index.md
      getting-started.md
    v2/                       # de, v2 - served at /de/v2/
      index.md
```

**Configuration in `dorcs.yaml`:**

```yaml
languages:
  default: "en"               # Default language
  enabled:                    # Explicit configuration (required)
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"

versions:
  default: "latest"           # Default version
  enabled:                    # Explicit configuration (required)
    - id: "latest"
      name: "Latest"
    - id: "v1"
      name: "Version 1"
    - id: "v2"
      name: "Version 2"
```

**URL Structure (Language-First):**
- Default (en, latest): `/getting-started` → `docs/en/getting-started.md`
- English v1: `/en/v1/getting-started` → `docs/en/v1/getting-started.md`
- German v1: `/de/v1/getting-started` → `docs/de/v1/getting-started.md`
- German default: `/de/getting-started` → `docs/de/getting-started.md`

## Relative Paths & Images

### How Relative Paths Work

When using multiple languages, the default language files are in `docs/en/` but served at `/` (no prefix). Relative paths are automatically handled:

**Example:**
- File: `docs/en/index.md` (served at `/`)
- Image: `docs/en/logo.png`
- Markdown: `![Logo](./logo.png)`

**Result:** The image path is automatically rewritten to `/logo.png`, which correctly resolves to `docs/en/logo.png`.

This ensures that relative paths work correctly even when the default language uses its folder but the URL has no language prefix.

### Best Practices for Images

1. **Use relative paths**: `./logo.png` or `../images/logo.png`
2. **Language-specific images**: Place images in language folders alongside markdown files
3. **Shared assets**: Place language-independent assets (logos, common images) in root `docs/` folder
4. **Absolute paths work too**: `/logo.png` works from any page

**Example structure:**
```
docs/
  shared-logo.png            # Shared asset - accessible from all languages at /shared-logo.png
  en/
    index.md
    logo.png                  # Language-specific - referenced as ./logo.png
    guide/
      installation.md
      images/
        screenshot.png        # Referenced as ./images/screenshot.png
  de/
    index.md
    logo.png                  # German-specific logo
```

**Shared assets in root:**
- Files like `docs/shared-logo.png` are accessible from all languages at `/shared-logo.png`
- Useful for logos, common images, or other language-independent assets
- Markdown files in root are **ignored** when using multiple languages

### Relative Links in Markdown

Relative links to other documentation pages are automatically rewritten:

- `[Getting Started](./getting-started.md)` → `/getting-started` (for default language)
- `[Installation](../guide/installation.md)` → `/guide/installation` (for default language)
- Links automatically include language/version prefixes when needed

## Configuration

Languages and versions must be explicitly configured in `dorcs.yaml`:

```yaml
languages:
  default: "en"
  enabled:                    # Explicit list (required)
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"

versions:
  default: "latest"
  enabled:                    # Explicit list (required)
    - id: "latest"
      name: "Latest"
    - id: "v1"
      name: "Version 1"
```

**Note:** Languages and versions must be explicitly configured. There is no auto-detection - this ensures predictable behavior and prevents content folders from being misidentified.

## Summary

| Scenario | Structure | Default URL | Other URLs | Config Required |
|----------|-----------|-------------|------------|----------------|
| Simple | `docs/` | `/` | - | No |
| Languages | `docs/en/`, `docs/de/` | `/` (default lang) | `/de/` | Yes - explicit config required |
| Versions | `docs/v1/`, `docs/v2/` | `/` (default version) | `/v1/` | Yes - explicit config required |
| Both | `docs/en/v1/`, `docs/de/v1/` | `/` (default lang/version) | `/en/v1/`, `/de/v1/` | Yes - explicit config required |

**Key Rules:**
- **Single language**: Files can go directly in `docs/`
- **Multiple languages**: 
  - Markdown files in root `docs/` are **ignored** - all markdown must be in `docs/{lang}/` folders
  - Static assets in root `docs/` are served as **shared assets** accessible from all languages
- **Default language**: Must be defined in `dorcs.yaml` when using multiple languages
- **Relative paths**: Automatically handled for all scenarios

Remember: **Minimal boilerplate, maximum flexibility!**
