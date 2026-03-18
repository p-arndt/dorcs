---
title: "Dorcs"
description: "Single-binary documentation server and static site builder for Markdown."
tags: [docs, overview]
date: 2026-03-18
draft: false
---

# Dorcs

<div style="display: flex; justify-content: center; align-items: center;">
<img src="./logo.png" alt="Dorcs Logo" width="200" height="200" style="border: none;" />
</div>

Dorcs turns a folder of Markdown files into a documentation site and can also export it as static HTML.

## Why dorcs

- Single binary, no runtime stack
- Clean extensionless URLs
- Generated navigation, table of contents, and live search in server mode
- Live reload in development
- Static builds for deployment
- Optional themes, GitHub-backed docs, versions, languages, presentations, and edit mode

## Start here


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

1. Read [Getting Started](./01_getting-started.md) for the fastest path to a running site.
2. Use [Configuration](./03_configuration.md) to customize branding, navigation, and behavior.
3. See [Writing Your Docs](./usage/writing-your-docs.md) for content structure and authoring patterns.
4. Use [Deployment](./04_deployment.md) when you are ready to ship a static build.

## Core workflow

```bash
dorcs init
dorcs --watch
dorcs build
```

## What the docs cover

- [Installation](./02_installation.md): binaries and building from source
- [Configuration](./03_configuration.md): `dorcs.yaml` and override rules
- [Usage](./usage/index.md): file layout, metadata, ordering, watch mode
- [Markdown](./06_markdown/index.md): supported syntax and extensions
- [Presentations](./07_presentations/index.md): slide decks from Markdown
- [Commands](./07_commands.md): CLI reference
- [External Content](./external-content/index.md): GitHub as a docs source
- [Edit Mode](./08_edit.md): browser-based editing with authentication
