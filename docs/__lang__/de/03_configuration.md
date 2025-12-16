---
title: "Konfiguration"
description: "Vollständiger Leitfaden zur Konfiguration von dorcs mit dorcs.yaml und Befehlszeilenoptionen."
tags: [configuration, settings, customization]
date: 2025-12-13
draft: false
---

# Konfiguration

Passen Sie dorcs mit einer Konfigurationsdatei oder Befehlszeilen-Flags an.

## Konfigurationsdatei

Erstellen Sie eine `dorcs.yaml` (oder `dorcs.yml` oder `dorcs.json`) Datei. Der Server erkennt sie automatisch.

> [!TIP]
> Verwenden Sie YAML für die Konfigurationsdatei, da es das am besten lesbare Format ist.

### Dateispeicherort

dorcs sucht die Konfigurationsdatei in dieser Reihenfolge:
1. **Aktuelles Arbeitsverzeichnis** - `./dorcs.yaml`
2. **Dokumentationsverzeichnis** - `./docs/dorcs.yaml` (wenn anders als Arbeitsverzeichnis)
3. **Benutzerdefinierter Pfad** - Verwenden Sie `--config /pfad/zu/config.yaml`, um einen genauen Pfad anzugeben

### Server-Port

Konfigurieren Sie den Port, auf dem dorcs laufen soll:

```yaml
port: 8000  # Optional: Port angeben (Standard: 8080)
```

> [!IMPORTANT]
> Das `--addr` Befehlszeilen-Flag überschreibt die Port-Einstellung in der Konfigurationsdatei, falls angegeben.

## Konfigurationsabschnitte

### Site

```yaml
site:
  title: "My Docs"              # Seitentitel (wird im Header angezeigt)
  description: "..."            # Meta-Beschreibung
  logo: "/logo.png"             # Logo-Bild-URL (optional)
                                 # Platzieren Sie logo.png im Root-Verzeichnis, wo dorcs läuft
  favicon: "/favicon.ico"       # Benutzerdefiniertes Favicon (optional)
                                 # Platzieren Sie favicon.ico im Root-Verzeichnis, wo dorcs läuft
```

> [!NOTE]
> **Statische Dateien Speicherort**: Logo- und Favicon-Dateien sollten im **Root-Verzeichnis** platziert werden, wo Sie das dorcs-Executable ausführen (das aktuelle Arbeitsverzeichnis). Der Server prüft zuerst das Root-Verzeichnis auf statische Assets, dann fällt er auf das Dokumentationsverzeichnis zurück. Dateien sind unter ihrem Pfad erreichbar (z.B. `/logo.png`). Wenn Sie ein `--base-url` Präfix verwenden, wird der BasePath automatisch diesen URLs vorangestellt.

### Theme

**Presets:** `default`, `ocean`, `forest`, `sunset`, `midnight`, `lavender`, `rose` und mehr. Siehe [Themes](./05_themes.md) für alle Optionen.

```yaml
theme:
  preset: midnight              # Theme-Preset
  mode: auto                    # light, dark oder auto
  custom_css: "custom.css"      # Benutzerdefiniertes CSS (relativ zum docs-Verzeichnis)
  font_family: '"Inter", system-ui, sans-serif'
  mono_font_family: '"JetBrains Mono", monospace'
```

**Benutzerdefinierte Farben:**
```yaml
theme:
  preset: default
  colors:
    light:
      background: "#ffffff"
      foreground: "#1f2328"
      accent: "#0969da"
      # ... weitere Farboptionen
    dark:
      background: "#0d1117"
      foreground: "#e6edf3"
      accent: "#2f81f7"
      # ... weitere Farboptionen
```

**Hinweis:** Syntax-Hervorhebung wird automatisch durch das Preset-Theme bestimmt.

### Navigation

```yaml
nav:
  show_search: true             # Suchfeld aktivieren/deaktivieren
  expand_all: false             # Alle Ordner standardmäßig erweitern
  links:                        # Header-Navigationslinks
    - title: "GitHub"
      url: "https://github.com/..."
      external: true
      icon: "github"            # github, twitter, discord, external
```

### Footer

```yaml
footer:
  copyright: "© 2024 Your Name"
  text: "Zusätzlicher Footer-Text"
  show_powered_by: true         # "Powered by dorcs" anzeigen
  links:
    - title: "Datenschutz"
      url: "/privacy"
```

### Mehrsprachige Unterstützung

Konfigurieren Sie mehrere Sprachen für Ihre Dokumentation:

```yaml
languages:
  default: "en"                 # Standard-Sprachcode (wird unter Root-URL bereitgestellt)
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
    - code: "fr"
      name: "Français"
```

**Wie es funktioniert:**

- **Standardsprache**: Wird unter Root-URL (`/`) bereitgestellt - Dokumentation in `docs/` Ordner platzieren
- **Andere Sprachen**: Werden unter `/{lang}/` bereitgestellt - Dokumentation in `docs/{lang}/` Ordnern platzieren
- **Sprachumschalter**: Erscheint automatisch im Header, wenn mehrere Sprachen aktiviert sind
- **URL-Struktur**:
  - Standard: `/getting-started` → `docs/getting-started.md`
  - Deutsch: `/de/getting-started` → `docs/de/getting-started.md`
  - Französisch: `/fr/getting-started` → `docs/fr/getting-started.md`

**Dateistruktur-Beispiel:**

```
docs/
  index.md                    # Standardsprache (Englisch)
  getting-started.md
  guide/
    installation.md
  de/                         # Deutscher Sprachordner
    index.md
    getting-started.md
    guide/
      installation.md
  fr/                         # Französischer Sprachordner
    index.md
    getting-started.md
    guide/
      installation.md
```

**Hinweise:**

- Jeder Sprachordner sollte eine vollständige Kopie Ihrer Dokumentationsstruktur haben
- Der Sprachumschalter behält den aktuellen Seitenpfad beim Wechseln der Sprache bei
- Wenn nur eine Sprache konfiguriert ist (oder keine), wird der Sprachumschalter ausgeblendet
- Sprachcodes sollten dem ISO 639-1 Standard folgen (z.B. "en", "de", "fr", "es", "ja")

### Authentifizierung & Bearbeitungsmodus

Aktivieren Sie die Online-Bearbeitung mit Benutzername/Passwort-Authentifizierung:

```yaml
auth:
  enabled: true                  # Authentifizierung und Bearbeitungsmodus aktivieren
  username: "admin"              # Benutzername für Login
  password: "your-secure-password"  # Passwort (wird automatisch gehasht)
  sessions_path: ".dorcs_sessions.json"  # Optional: benutzerdefinierter Pfad für Sessions-Datei
```

**Wie es funktioniert:**

- Wenn aktiviert, erscheint ein **Login**-Button in der Fußzeile
- Nach dem Einloggen erscheinen **Bearbeiten**- und **Abmelden**-Buttons im Header
- Klicken Sie auf **Bearbeiten**, um das Bearbeitungsmodus-Panel zu öffnen, wo Sie können:
  - Dateien durchsuchen und bearbeiten
  - Neue Dateien und Ordner erstellen
  - Dateien löschen
  - Änderungen direkt im Dateisystem speichern

**Sicherheitshinweise:**

- Passwörter werden automatisch mit Argon2id beim ersten Gebrauch gehasht
- Sessions laufen nach 24 Stunden ab
- Alle Bearbeitungsvorgänge erfordern Authentifizierung
- Der Passwort-Hash wird nach dem ersten Login in der Konfigurationsdatei gespeichert

> [!WARNING]
> Aktivieren Sie den Bearbeitungsmodus nur in vertrauenswürdigen Netzwerken oder hinter ordnungsgemäßer Authentifizierung. Der Bearbeitungsmodus ermöglicht vollen Dateisystemzugriff auf Ihr Dokumentationsverzeichnis.

## Befehlszeilen-Flags

Alle Konfigurationsoptionen können über Befehlszeilen-Flags überschrieben werden:

| Flag | Beschreibung | Beispiel |
|------|-------------|---------|
| `--dir` | Dokumentationsverzeichnis | `--dir ./docs` |
| `--addr` | Listen-Adresse (überschreibt Config-Port) | `--addr :8080` |
| `--base-url` | URL-Pfad-Präfix | `--base-url /docs` |
| `--title` | Seitentitel | `--title "My Docs"` |
| `--config` | Konfigurationsdatei-Pfad | `--config /pfad/zu/config.yaml` |
| `--theme` | Theme-Preset | `--theme midnight` |
| `--theme-mode` | Theme-Modus | `--theme-mode dark` |
| `--cache` | Caching aktivieren | `--cache=true` |
| `--no-drafts` | Entwürfe ausblenden | `--no-drafts=true` |
| `--watch` | Watch-Modus | `--watch` |

**Hinweis:** Befehlszeilen-Flags überschreiben Konfigurationsdatei-Einstellungen. Das `--addr` Flag überschreibt die `port` Einstellung in der Konfigurationsdatei.

## Beispiele

### Minimal

```yaml
port: 8080  # Optional: Standard ist 8080

site:
  title: "My Docs"
```

### Vollständige Anpassung

```yaml
port: 8000  # Optional: Port angeben (Standard: 8080)

site:
  title: "My Awesome Documentation"
  description: "Vollständiger Leitfaden zu meinem Projekt"
  logo: "/logo.svg"        # Platzieren Sie logo.svg im Root-Verzeichnis, wo dorcs läuft
  favicon: "/favicon.ico"  # Platzieren Sie favicon.ico im Root-Verzeichnis, wo dorcs läuft

theme:
  mode: auto
  preset: midnight
  custom_css: "overrides.css"
  font_family: '"SF Pro Display", system-ui, sans-serif'
  mono_font_family: '"SF Mono", monospace'

nav:
  show_search: true
  expand_all: true
  links:
    - title: "GitHub"
      url: "https://github.com/user/repo"
      external: true
      icon: "github"
    - title: "Discord"
      url: "https://discord.gg/server"
      external: true
      icon: "discord"

footer:
  copyright: "© 2024 My Company"
  text: "Erstellt mit ❤️ mit dorcs"
  show_powered_by: false
  links:
    - title: "Datenschutz"
      url: "/privacy"
    - title: "Nutzungsbedingungen"
      url: "/terms"

languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
    - code: "fr"
      name: "Français"

auth:
  enabled: true
  username: "admin"
  password: "secure-password-here"
```

## Nächste Schritte

- 🎨 [Themes](./05_themes.md) - Durchsuchen Sie alle verfügbaren Themes
- 🚀 [Deployment](./04_deployment.md) - In Produktion bereitstellen
