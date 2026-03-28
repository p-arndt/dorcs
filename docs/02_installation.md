---
title: "Installation"
description: "Download Dorcs or build it from source."
tags: [installation]
date: 2026-03-18
draft: false
---

# Installation

Dorcs is a single file — download it, and you're ready to go. No runtime, no package manager, no dependencies.

## Download a release {badge:RECOMMENDED}

Grab the latest release for your platform:

:::tabs
::tab Windows
Download `dorcs-windows-amd64.exe` from the [releases page](https://github.com/p-arndt/dorcs/releases) and run it:

```powershell
.\dorcs.exe
```
::tab Linux
Download `dorcs-linux-amd64` from the [releases page](https://github.com/p-arndt/dorcs/releases), make it executable, and run:

```bash
chmod +x dorcs
./dorcs
```
::tab macOS
macOS builds aren't pre-built yet — see the "Build from source" section below.
:::

## Build from source

If you prefer to build it yourself, you'll need **Go 1.25 or newer**.

```bash
go build -o dorcs ./cmd/dorcs
```

Or run it directly without building:

```bash
go run ./cmd/dorcs
```

:::accordion Cross-platform build commands
Building for a different OS? Use Go's cross-compilation:

```bash
# Build for Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dorcs ./cmd/dorcs

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o dorcs.exe ./cmd/dorcs
```
:::

## Verify it works

```bash
dorcs --help
```

If you see the help text, you're all set.

> [!TIP]
> Next step: follow the [Getting Started](./01_getting-started.md) guide to create your first site.
