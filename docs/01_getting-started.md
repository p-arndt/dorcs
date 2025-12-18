---
title: "Getting Started"
description: "Get your documentation site running in minutes."
tags: [getting-started, quickstart]
date: 2025-12-13
draft: false
after: "index"
---

# Getting Started

Get your documentation site running in under 5 minutes.

## Quick Start

Get your documentation site running in under 5 minutes. Choose your preferred method:

## Download 

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

After downloading:

**Windows:**
```pwsh
.\dorcs-windows-amd64.exe 
```

**Linux:**
```bash
# Make it executable
chmod +x dorcs-linux-amd64.exe

# Run dorcs
./dorcs-linux-amd64.exe 
```

Visit [http://localhost:8080](http://localhost:8080) to see your documentation!


## Your First Documentation

The easiest way to get started is using the `init` command, which sets up everything for you:

### 1. **Initialize your documentation site**:

```bash
./dorcs init
```

Also see the [Commands](./07_commands.md#init-command) page for more information.

### 2. **Start the development server**:

> [!TIP]
> Enable watch mode (`--watch`) for live reload during development

```bash
./dorcs --watch
```

### 3. **Open your browser**

Navigate to [http://localhost:8080](http://localhost:8080)

You'll see your new documentation site with a welcome page! Edit `docs/index.md` to customize your homepage.

## File Structure

> [!TIP]
> Put `index.md` inside folders to create section landing pages

dorcs uses **extensionless URLs** that map directly to your file structure:

| File Path                    | URL                   |
| ---------------------------- | --------------------- |
| `docs/index.md`              | `/`                   |
| `docs/getting-started.md`    | `/getting-started`    |
| `docs/guide/index.md`        | `/guide`              |
| `docs/guide/installation.md` | `/guide/installation` |

