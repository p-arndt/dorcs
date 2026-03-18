---
title: "Writing Your Docs"
description: "How to write your docs"
tags: [writing, docs]
date: 2025-12-14
draft: false
---

# Writing Your Docs

## How to write your docs

Markdown is a lightweight markup language with a simple syntax. It is designed to be easy to write and read. (Also See [Markdown Basics](../06_markdown/01_basics.md) for more information.)

Markdown is also highly understoon by most text editors and IDEs, so you can write your docs in your favorite editor. Another advantage is that it's also easy understandable for LLMs, so you can use them to generate your docs.

Here is an example of a simple Markdown file:

```markdown
---
title: "My First Doc"
---

# My First Doc

This is my first doc.
```



This will watch for changes in the `docs` directory and automatically reload the page.

## File Organization

### Single Language

For a single-language documentation site, simply place all your markdown files in the `docs/` directory:

```
docs/
  index.md
  getting-started.md
  guide/
    installation.md
    advanced.md
```

### Multi-lingual Documentation

If you've configured multiple languages in your `dorcs.yaml`, **define the default language in the config** and organize your documentation using **MkDocs-style structure**:

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

**Configuration:**

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

**Important:**
- **Markdown files**: When using multiple languages, markdown files in root `docs/` are **ignored**. All markdown files must be in language folders
- **Shared assets**: Static assets (images, PDFs, etc.) in root `docs/` are served as **language-independent shared assets**
- **MkDocs-style**: Direct language folders
- **Default language**: Defined in `dorcs.yaml` - default language files are in `docs/en/` but served at `/` (no prefix)
- Each language folder should mirror the structure of your default language
- **Relative paths**: Automatically handled - images and links work correctly
- See [Configuration](../03_configuration.md#multi-lingual-support) for setup instructions
- See [File Structure & Organization](./file-structure.md) for complete guide

> [!WARNING]
> **Markdown files in root are ignored**: When using multiple languages, any `.md` files in the root `docs/` folder are ignored. Only static assets in root are served as shared assets.

## YAML Front Matter

You can add metadata to your Markdown files using YAML front matter. This is a simple way to add metadata to your docs. See [Metadata](metadata.md) for more information.
