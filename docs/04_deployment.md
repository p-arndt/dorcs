---
title: "Deployment"
description: "Deploy dorcs to production on various platforms and services."
tags: [deployment, production, hosting]
date: 2025-12-13
draft: false
---

# Deployment

Deploy your dorcs documentation site to production.



## Build for Production

Build a static binary with no dependencies:

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

## Deployment Options

**Docker Compose:**

The project includes a `docker-compose.yml` file. For production, configure via `dorcs.yaml`:

```yaml
services:
  dorcs:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./docs:/docs
    restart: unless-stopped
```

## Monitoring

**Health check:**
```bash
curl http://localhost:8080/api/health
```


## Next Steps

- ⚙️ [Configuration](./03_configuration.md) - Production configuration
- 🎨 [Themes](./05_themes.md) - Customize appearance
