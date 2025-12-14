---
title: "Installation"
description: "Detailed installation instructions for all platforms."
tags: [installation, setup]
date: 2025-12-13
draft: false
---

# Installation

Install dorcs using Docker, pre-built binaries, or build from source.

## Prerequisites

- **Docker and Docker Compose** (for Docker installation - recommended)
- **Go 1.25+** (if building from source)
- **No runtime dependencies** (single static binary)

## Installation Methods

### Docker Compose (Recommended)

**The easiest way - no Go installation needed.**

1. Clone the repository:
   ```bash
   git clone https://github.com/p-arndt/dorcs.git
   cd dorcs
   ```

2. Start with Docker Compose:
   ```bash
   docker-compose up
   ```

   Your site is running at `http://localhost:8080`.

### Docker (without Compose)

```bash
# Build the image
docker build -t dorcs .

# Run the container
docker run -p 8080:8080 -v $(pwd)/docs:/docs dorcs
```

### Pre-built Binary

1. Download from the [releases page](https://github.com/p-arndt/dorcs/releases)
2. Choose your platform:
   - **Linux**: `dorcs-linux-amd64`
   - **macOS**: `dorcs-darwin-amd64` or `dorcs-darwin-arm64`
   - **Windows**: `dorcs-windows-amd64.exe`
3. Make executable (Linux/macOS): `chmod +x dorcs`
4. Run: `./dorcs --dir ./docs`

### Build from Source

```bash
# Clone the repository
git clone https://github.com/p-arndt/dorcs.git
cd dorcs-v2

# Build
go build -o dorcs ./cmd/dorcs

# Verify
./dorcs --help
```

### Install via Go

```bash
go install github.com/p-arndt/dorcs/cmd/dorcs@latest
```

Installs to `$GOPATH/bin` or `$HOME/go/bin` (add to PATH if needed).

## Building Static Binaries

For deployment in containers or minimal environments:

**Linux:**
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o dorcs \
  ./cmd/dorcs
```

**Windows:**
```cmd
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -trimpath -ldflags="-s -w" -o dorcs.exe ./cmd/dorcs
```

**macOS:**
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o dorcs \
  ./cmd/dorcs
```

The `-s -w` flags strip debug symbols and reduce binary size.

## Verify Installation

**With Docker Compose:**
```bash
docker-compose up -d
docker-compose logs
# Visit http://localhost:8080
```

**With Binary:**
```bash
./dorcs --help
./dorcs --dir ./docs --addr :8080
# Visit http://localhost:8080
```

## Next Steps

- ✅ [Getting Started](./01_getting-started.md) - Run your first documentation site
- ⚙️ [Configuration](./03_configuration.md) - Customize your setup
- 🚀 [Deployment](./04_deployment.md) - Deploy to production
