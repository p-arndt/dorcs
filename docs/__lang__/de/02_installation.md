---
title: "Installation"
description: "Detaillierte Installationsanweisungen für alle Plattformen."
tags: [installation, setup]
date: 2025-12-13
draft: false
---

# Installation

Installieren Sie dorcs mit Docker, vorkompilierten Binaries oder bauen Sie es aus dem Quellcode.

## Voraussetzungen

- **Docker und Docker Compose** (für Docker-Installation - empfohlen)
- **Go 1.25+** (wenn aus dem Quellcode gebaut wird)
- **Keine Laufzeitabhängigkeiten** (einzelnes statisches Binary)

## Installationsmethoden

### Docker Compose (Empfohlen)

**Der einfachste Weg - keine Go-Installation erforderlich.**

1. Repository klonen:
   ```bash
   git clone https://github.com/p-arndt/dorcs.git
   cd dorcs
   ```

2. Mit Docker Compose starten:
   ```bash
   docker-compose up
   ```

   Ihre Seite läuft unter `http://localhost:8080`.

### Docker (ohne Compose)

```bash
# Image bauen
docker build -t dorcs .

# Container ausführen
docker run -p 8080:8080 -v $(pwd)/docs:/docs dorcs
```

### Vorkompiliertes Binary

1. Laden Sie von der [Releases-Seite](https://github.com/p-arndt/dorcs/releases) herunter
2. Wählen Sie Ihre Plattform:
   - **Linux**: `dorcs-linux-amd64`
   - **macOS**: `dorcs-darwin-amd64` oder `dorcs-darwin-arm64`
   - **Windows**: `dorcs-windows-amd64.exe`
3. Ausführbar machen (Linux/macOS): `chmod +x dorcs`
4. Ausführen: `./dorcs --dir ./docs`

### Aus dem Quellcode bauen

```bash
# Repository klonen
git clone https://github.com/p-arndt/dorcs.git
cd dorcs-v2

# Bauen
go build -o dorcs ./cmd/dorcs

# Überprüfen
./dorcs --help
```

### Über Go installieren

```bash
go install github.com/p-arndt/dorcs/cmd/dorcs@latest
```

Installiert nach `$GOPATH/bin` oder `$HOME/go/bin` (zu PATH hinzufügen, falls nötig).

## Statische Binaries bauen

Für Deployment in Containern oder minimalen Umgebungen:

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

## Installation überprüfen

**Mit Docker Compose:**
```bash
docker-compose up -d
docker-compose logs
# Besuchen Sie http://localhost:8080
```

**Mit Binary:**
```bash
./dorcs --help
./dorcs --dir ./docs --addr :8080
# Besuchen Sie http://localhost:8080
```

## Nächste Schritte

- ✅ [Erste Schritte](./01_getting-started.md) - Starten Sie Ihre erste Dokumentationsseite
- ⚙️ [Konfiguration](./03_configuration.md) - Passen Sie Ihre Einrichtung an
- 🚀 [Deployment](./04_deployment.md) - In Produktion bereitstellen
