---
title: "Metadaten"
description: "Metadaten-Funktionen"
tags: [metadata, frontmatter]
date: 2025-12-14
draft: false
---

# Metadaten

Fügen Sie Metadaten zu Ihren Markdown-Dateien mit YAML Front Matter am Anfang der Datei hinzu:

```yaml
---
title: "Seitentitel"
description: "Seitenbeschreibung für SEO"
date: 2025-12-13
tags: [tag1, tag2, tag3]
draft: false
order: 1
---
# Erste Schritte

Ihr Inhalt hier...
```

**Verfügbare Felder:**

- `title` - Seitentitel (wird in Navigation und Browser-Tab verwendet)
- `description` - Meta-Beschreibung für SEO
- `date` - Veröffentlichungsdatum (YYYY-MM-DD Format)
- `tags` - Liste von Tags zur Kategorisierung
- `draft` - Auf `true` setzen, um aus der Navigation auszublenden (bei Verwendung von `--no-drafts`)
- `order` - Numerischer Wert zum Sortieren von Seiten in der Navigation (niedrigere Zahlen erscheinen zuerst)
- `author` - Autorenname (wird mit `[[AUTHOR]]` Platzhalter angezeigt)
- `after` - Schlüssel des Elements, nach dem dies in der Navigation erscheinen soll (verwenden Sie `"index"`, um nach index.md zu platzieren) (Siehe [Reihenfolge der Dokumente](./order-of-docs.md) für weitere Informationen.)
