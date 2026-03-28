---
title: "Dorcs"
description: "Turn your Markdown files into a beautiful documentation site with a single binary."
tags: [docs, overview]
date: 2026-03-18
draft: false
---

# Dorcs

<div style="display: flex; justify-content: center; align-items: center;">
<img src="./logo.png" alt="Dorcs Logo" width="200" height="200" style="border: none;" />
</div>

**Dorcs** turns a folder of Markdown files into a beautiful documentation website. One binary, zero dependencies — just write and ship.

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

## Three commands. That's it.

```bash
dorcs init        # scaffold a docs site
dorcs --watch     # start writing with live reload
dorcs build       # export static HTML
```

## What you get out of the box

Dorcs isn't just a Markdown renderer — it's a complete documentation platform:

:::tabs
::tab Look & Feel
- **Beautiful themes** with light and dark mode — [see them all](./05_themes.md)
- Automatically generated **sidebar navigation** from your folder structure
- **Section tabs** to organize large sites into top-level categories
- **Breadcrumbs**, **previous/next links**, and a **back-to-top** button
- **Custom 404 page** that matches your theme
::tab Writing Features
- **Tabs** for platform-specific content (like this one!)
- **Callouts** for tips, warnings, and important notes
- **Badges** like {badge:NEW} and {badge:BETA} to label features
- **Accordions** for collapsible content
- **Video embeds** from YouTube and Vimeo
- **Mermaid diagrams** and **KaTeX math**
::tab Developer Experience
- **Live reload** — edit a file, browser updates instantly
- **Clean URLs** — `docs/guide/install.md` becomes `/guide/install`
- **Built-in search** with instant results
- **Code copy buttons** on every code block
- **Heading anchor links** for easy sharing
- **Sitemap** and **Open Graph tags** for SEO
::tab Advanced
- **Multiple languages** with automatic language switcher
- **Doc versioning** — maintain v1, v2, latest side by side
- **GitHub integration** — serve docs directly from a repo
- **Browser-based editing** with authentication
- **Slide decks** — turn any Markdown page into a presentation
- **Static export** — deploy anywhere as plain HTML
:::

## Where to start

> [!TIP]
> New to Dorcs? Start with [Getting Started](./01_getting-started.md) — you'll have a working site in under a minute.

Already set up? Jump to what you need:

- **[Writing Docs](./usage/index.md)** — how to structure your content and use all the Markdown features
- **[Customize Your Site](./config/index.md)** — themes, branding, navigation, and more
- **[Deploy](./04_deployment.md)** — ship your site to GitHub Pages, Netlify, or any host
- **[Commands](./07_commands.md)** — full CLI reference
