---
title: "Writing Your Docs"
description: "How to write your docs"
tags: [writing, docs]
date: 2025-12-14
draft: false
after: "index"
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

If you've configured multiple languages in your `dorcs.yaml`, organize your documentation like this:

```
docs/
  index.md                    # Default language (e.g., English)
  getting-started.md
  guide/
    installation.md
  __lang__/                   # Language-specific folder
    de/                       # German language folder
      index.md
      getting-started.md
      guide/
        installation.md
    fr/                       # French language folder
      index.md
      getting-started.md
      guide/
        installation.md
```

**Important:**
- Each language folder should mirror the structure of your default language
- The default language stays in the root `docs/` folder
- Other languages go in `docs/__lang__/{lang}/` folders where `{lang}` is the language code
- The `__lang__` folder keeps language-specific content separate and avoids conflicts with regular folders
- See [Configuration](../03_configuration.md#multi-lingual-support) for setup instructions

## YAML Front Matter

You can add metadata to your Markdown files using YAML front matter. This is a simple way to add metadata to your docs. See [Metadata](metadata.md) for more information.