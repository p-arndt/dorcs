---
title: "Ihre Dokumentation schreiben"
description: "Wie Sie Ihre Dokumentation schreiben"
tags: [writing, docs]
date: 2025-12-14
draft: false
after: "index"
---

# Ihre Dokumentation schreiben

## Wie Sie Ihre Dokumentation schreiben

Markdown ist eine leichtgewichtige Auszeichnungssprache mit einer einfachen Syntax. Sie ist darauf ausgelegt, einfach zu schreiben und zu lesen zu sein. (Siehe auch [Markdown-Grundlagen](../06_markdown/01_basics.md) für weitere Informationen.)

Markdown wird auch von den meisten Texteditoren und IDEs gut verstanden, sodass Sie Ihre Dokumentation in Ihrem bevorzugten Editor schreiben können. Ein weiterer Vorteil ist, dass es auch für LLMs leicht verständlich ist, sodass Sie diese verwenden können, um Ihre Dokumentation zu generieren.

Hier ist ein Beispiel für eine einfache Markdown-Datei:

```markdown
---
title: "Mein erstes Dokument"
---

# Mein erstes Dokument

Dies ist mein erstes Dokument.
```



Dies überwacht Änderungen im `docs` Verzeichnis und lädt die Seite automatisch neu.

## Dateiorganisation

### Einzelsprache

Für eine einsprachige Dokumentationsseite platzieren Sie einfach alle Ihre Markdown-Dateien im `docs/` Verzeichnis:

```
docs/
  index.md
  getting-started.md
  guide/
    installation.md
    advanced.md
```

### Mehrsprachige Dokumentation

Wenn Sie mehrere Sprachen in Ihrer `dorcs.yaml` konfiguriert haben, organisieren Sie Ihre Dokumentation wie folgt:

```
docs/
  index.md                    # Standardsprache (z.B. Englisch)
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

**Wichtig:**
- Jeder Sprachordner sollte die Struktur Ihrer Standardsprache widerspiegeln
- Die Standardsprache bleibt im Root-`docs/` Ordner
- Andere Sprachen gehen in `docs/{lang}/` Ordner, wobei `{lang}` der Sprachcode ist
- Siehe [Konfiguration](../03_configuration.md#multi-lingual-support) für Setup-Anweisungen

## YAML Front Matter

Sie können Metadaten zu Ihren Markdown-Dateien mit YAML Front Matter hinzufügen. Dies ist eine einfache Möglichkeit, Metadaten zu Ihren Dokumenten hinzuzufügen. Siehe [Metadaten](metadata.md) für weitere Informationen.
