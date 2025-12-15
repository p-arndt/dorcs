---
title: "Erste Schritte"
description: "Bringen Sie Ihre Dokumentationsseite in wenigen Minuten zum Laufen."
tags: [getting-started, quickstart]
date: 2025-12-13
draft: false
after: "index"
---

# Erste Schritte

Bringen Sie Ihre Dokumentationsseite in weniger als 5 Minuten zum Laufen.

## Schnellstart

Bringen Sie Ihre Dokumentationsseite in weniger als 5 Minuten zum Laufen. Wählen Sie Ihre bevorzugte Methode:

### Option 1: Docker Compose (Empfohlen)

```bash
git clone https://github.com/p-arndt/dorcs.git
cd dorcs
docker-compose up
```

### Option 2: Vorkompiliertes Binary

1. Laden Sie die neueste Version von der [Releases-Seite](https://github.com/p-arndt/dorcs/releases) herunter
2. Machen Sie es ausführbar (Linux/macOS): `chmod +x dorcs`
3. Ausführen: `./dorcs --dir ./docs`

### Option 3: Aus dem Quellcode bauen

```bash
git clone https://github.com/p-arndt/dorcs.git
cd dorcs
go build -o dorcs ./cmd/dorcs
./dorcs --dir ./docs
```

Besuchen Sie `http://localhost:8080`, um Ihre Dokumentation zu sehen!

## Ihre erste Dokumentation

> [!TIP]
> Aktivieren Sie den Watch-Modus (`--watch`) für Live-Reload während der Entwicklung

1. **Erstellen Sie Ihr Dokumentationsverzeichnis** (oder verwenden Sie den vorhandenen `./docs` Ordner)

2. **Fügen Sie Ihre erste Markdown-Datei hinzu**:

```bash
cat > docs/index.md << 'EOF'
---
title: "Willkommen"
description: "Meine Dokumentationsseite"
---

# Willkommen

Dies ist meine Dokumentationsseite, betrieben von dorcs!
EOF
```

3. **Starten Sie den Server**:

```bash
# Mit Docker Compose
docker-compose up

# Oder mit Binary
./dorcs --dir ./docs
```

4. **Öffnen Sie Ihren Browser**: Navigieren Sie zu `http://localhost:8080`

## Dateistruktur

> [!TIP]
> Platzieren Sie `index.md` in Ordnern, um Abschnitts-Landingpages zu erstellen

dorcs verwendet **erweiterungslose URLs**, die direkt auf Ihre Dateistruktur abgebildet werden:

| Dateipfad                      | URL                   |
| ------------------------------ | --------------------- |
| `docs/index.md`                | `/`                   |
| `docs/getting-started.md`      | `/getting-started`    |
| `docs/guide/index.md`          | `/guide`              |
| `docs/guide/installation.md`   | `/guide/installation` |

## Front Matter

Fügen Sie Metadaten zu Ihren Markdown-Dateien mit YAML Front Matter hinzu:

```yaml
---
title: "Erste Schritte"
description: "Lernen Sie, wie Sie loslegen"
date: 2025-12-13
tags: [tutorial, anfänger]
draft: false
---
# Erste Schritte

Ihr Inhalt hier...
```

**Verfügbare Felder:**

- `title` - Seitentitel (wird in Navigation und Browser-Tab verwendet)
- `description` - Meta-Beschreibung für SEO
- `date` - Veröffentlichungsdatum (YYYY-MM-DD)
- `tags` - Liste von Tags zur Kategorisierung
- `draft` - Auf `true` setzen, um aus der Navigation auszublenden (bei Verwendung von `--no-drafts`)
- `order` - Numerischer Wert zum Sortieren von Seiten in der Navigation (niedrigere Zahlen erscheinen zuerst)
- `author` - Autorenname (wird mit [[AUTHOR]] Platzhalter angezeigt)
- `after` - Schlüssel des Elements, nach dem dies in der Navigation erscheinen soll (verwenden Sie `"index"`, um nach index.md zu platzieren)

Siehe [Metadaten](./usage/metadata.md) für detaillierte Dokumentation zu allen Front Matter Feldern, einschließlich Navigationssortierung mit `after` und `order`.

## Watch-Modus

Wenn Sie das Flag `--watch` verwenden, überwacht dorcs Änderungen in Ihrem Dokumentationsverzeichnis und lädt die Seite automatisch neu.

**Mit Binary:**

```bash
./dorcs --dir ./docs --watch
```

**Mit Docker Compose:**

> [!IMPORTANT]
> Dateiüberwachung über Docker-Volumes funktioniert möglicherweise nicht zuverlässig unter Windows/Mac. Für beste Ergebnisse führen Sie dorcs direkt auf Ihrem Host aus: `./dorcs --dir ./docs --watch`

Kommentieren Sie die `command` Zeile in `docker-compose.yml` aus:

```yaml
command: ["/app/dorcs", "--dir", "/docs", "--addr", "0.0.0.0:8080", "--watch"]
```

## Häufige Befehle

**Binary:**

```bash
./dorcs --dir ./docs                    # Grundlegende Verwendung
./dorcs --dir ./docs --addr :3000       # Benutzerdefinierter Port
./dorcs --dir ./docs --watch            # Entwicklungsmodus
./dorcs --dir ./docs --theme midnight   # Benutzerdefiniertes Theme
```

## Bearbeitungsmodus (Optional)

Aktivieren Sie die Online-Bearbeitung, um Ihre Dokumentation direkt aus dem Browser zu verwalten:

1. **In der Konfiguration aktivieren** (`dorcs.yaml`):
```yaml
auth:
  enabled: true
  username: "admin"
  password: "your-secure-password"
```

2. **Server neu starten** und **Login** in der Fußzeile klicken
3. **Bearbeiten** im Header klicken, um das Bearbeitungsmodus-Panel zu öffnen
4. **Dateien durchsuchen, bearbeiten, erstellen und löschen** direkt in Ihrem Browser

Siehe [Konfiguration](./03_configuration.md#authentication--edit-mode) für weitere Details.

## Nächste Schritte

- ⚙️ [Konfiguration](./03_configuration.md) - Passen Sie Ihre Seite mit `dorcs.yaml` an
- 🎨 [Themes](./05_themes.md) - Wählen Sie aus integrierten Themes
- 🚀 [Deployment](./04_deployment.md) - In Produktion bereitstellen

## Fehlerbehebung

**Port bereits in Verwendung?**

```bash
./dorcs --dir ./docs --addr :8081
```

**Dateien werden nicht angezeigt?**

- Stellen Sie sicher, dass Dateien die Endung `.md` oder `.markdown` haben
- Überprüfen Sie, dass Dateien nicht als `draft: true` markiert sind (außer bei Verwendung von `--no-drafts=false`)
- Überprüfen Sie, ob der `--dir` Pfad korrekt ist
