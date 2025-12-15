---
title: "Dorcs"
description: "Willkommen bei dorcs - ein Single-Binary statischer Dokumentationsserver für Markdown-Dateien."
tags: [docs, markdown]
date: 2025-12-13
draft: false
---


# Dorcs

<div style="display: flex; justify-content: center; align-items: center;">
<img src="./logo.png" alt="Dorcs Logo" width="200" height="200" style="border: none;" />
</div>

Dorcs ist ein Single-Binary statischer Dokumentationsserver für Markdown-Dateien. Es ist ein einfacher und benutzerfreundlicher Dokumentationsserver, mit dem Sie Ihre Dokumentationsseite in wenigen Minuten erstellen und hosten können.

## Loslegen

Bereit loszulegen? Schauen Sie sich die [Erste Schritte Anleitung](./01_getting-started.md) an, um Ihre Dokumentationsseite in weniger als 5 Minuten zum Laufen zu bringen.

## Dokumentation

- 🚀 [Erste Schritte](./01_getting-started.md) - Vollständige Schnellstart-Anleitung
- 📦 [Installation](./02_installation.md) - Detaillierte Installationsanweisungen
- ⚙️ [Konfiguration](./03_configuration.md) - Passen Sie Ihre Seite mit `dorcs.yaml` an
- 🚢 [Deployment](./04_deployment.md) - In Produktion bereitstellen
- 🎨 [Themes](./05_themes.md) - Durchsuchen Sie alle verfügbaren Themes
- 📝 [Markdown-Funktionen](./06_markdown/index.md) - Vollständiger Leitfaden zu Markdown-Funktionen

## Funktionen

- **Single Binary** – keine Laufzeitabhängigkeiten, statisch verlinkbar
- **Erweiterungslose URLs** – `/guide/getting-started` liefert `docs/guide/getting-started.md`
- **YAML Front Matter** – Metadaten-Unterstützung (Titel, Beschreibung, Datum, Tags, Entwurf)
- **Inhaltsverzeichnis** – automatisch aus Überschriften generiert mit Scrollspy
- **Navigations-Sidebar** – automatisch aus Ihrer Ordnerstruktur erstellt
- **Responsives Design** – mobilfreundlich mit zusammenklappbarer Sidebar
- **Dunkler Modus** – automatisch basierend auf Systemeinstellung
- **Live Reload** – Watch-Modus für die Entwicklung mit intelligenten Inhaltsaktualisierungen
- **Mehrere Themes** – wählen Sie aus über 20 integrierten Themes
- **Suche** – integrierte Suchfunktion
- **Bearbeitungsmodus** – Online-Bearbeitung mit Authentifizierung (Dateien direkt im Browser erstellen, bearbeiten, löschen)

## Wie es funktioniert

### URL-Routing

dorcs verwendet erweiterungslose URLs, die direkt auf Ihre Dateistruktur abgebildet werden:

| Dateipfad                      | URL                   |
| ------------------------------ | --------------------- |
| `docs/index.md`                | `/`                   |
| `docs/getting-started.md`      | `/getting-started`    |
| `docs/guide/index.md`          | `/guide`              |
| `docs/guide/installation.md`   | `/guide/installation` |

Es erstellt auch automatisch die Navigation aus Ihrer Struktur und generiert eine Sidebar und ein Inhaltsverzeichnis für jede Seite.
