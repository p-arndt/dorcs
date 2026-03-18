---
title: "Installation"
description: "Install Dorcs from a release or build it from source."
tags: [installation]
date: 2026-03-18
draft: false
---

# Installation

Dorcs is a Go application distributed as a single executable.

## Use a release binary

Download a release from the GitHub releases page and run it directly.

Linux:

```bash
chmod +x dorcs
./dorcs
```

Windows:

```powershell
.\dorcs.exe
```

macOS is currently a build-from-source path in this repository.

## Build from source

Requirements:

- Go 1.25+

Build:

```bash
go build -o dorcs ./cmd/dorcs
```

Run from source:

```bash
go run ./cmd/dorcs
```

## Cross-platform builds

Linux:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dorcs ./cmd/dorcs
```

Windows:

```bash
GOOS=windows GOARCH=amd64 go build -o dorcs.exe ./cmd/dorcs
```

## Verify the install

```bash
dorcs --help
```

Then continue with [Getting Started](./01_getting-started.md).