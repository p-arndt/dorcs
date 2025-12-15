---
title: "Deployment"
description: "Stellen Sie dorcs in Produktion auf verschiedenen Plattformen und Diensten bereit."
tags: [deployment, production, hosting]
date: 2025-12-13
draft: false
---

# Deployment

Stellen Sie Ihre dorcs-Dokumentationsseite in Produktion bereit.



## Für Produktion bauen

Bauen Sie ein statisches Binary ohne Abhängigkeiten:

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

Die Flags `-s -w` entfernen Debug-Symbole und reduzieren die Binary-Größe.

## Deployment-Optionen

**Docker Compose:**

Das Projekt enthält eine `docker-compose.yml` Datei. Für die Produktion konfigurieren Sie über `dorcs.yaml`:

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

**Health Check:**
```bash
curl http://localhost:8080/api/health
```


## Nächste Schritte

- ⚙️ [Konfiguration](./03_configuration.md) - Produktionskonfiguration
- 🎨 [Themes](./05_themes.md) - Erscheinungsbild anpassen
