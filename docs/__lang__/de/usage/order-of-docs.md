---
title: "Reihenfolge der Dokumente"
description: "Wie Sie Ihre Dokumente ordnen"
tags: [ordering, docs]
date: 2025-12-14
draft: false
---

# Reihenfolge der Dokumente

## Verwendung von `order` und numerischen Präfixen

Standardmäßig werden Seiten sortiert nach:

1. Numerischen Präfixen in Dateinamen (z.B. `01_`, `02_`)
2. `order` Feld im Front Matter
3. Alphabetisch nach Titel

Dateien mit numerischen Präfixen (wie `01_installation.md`) erscheinen vor Dateien ohne Präfixe, unabhängig von ihrem `order`-Wert.

## Verwendung von `after` für relative Positionierung

Das `after` Feld ermöglicht es Ihnen, eine Seite direkt nach einer anderen Seite in der Navigation zu platzieren, mit **absoluter Priorität** über numerische Präfixe und `order`-Werte.

**Eine Seite direkt nach dem Index platzieren:**

```yaml
---
title: "Erste Schritte"
after: "index"
---
# Erste Schritte
```

Dies platziert "Erste Schritte" direkt nach der Root-`index.md` oder Ordner-`index.md`, vor allen nummerierten Dateien.

**Eine Seite nach einer bestimmten Seite platzieren:**

```yaml
---
title: "Erweiterte Themen"
after: "getting-started"
---
# Erweiterte Themen
```

Dies platziert "Erweiterte Themen" direkt nach der Seite mit dem Schlüssel `getting-started`.

**Hinweis:** Das `after` Feld hat absolute Priorität über alle anderen Sortiermechanismen (numerische Präfixe, `order` Felder, etc.). Wenn mehrere Seiten `after: "index"` verwenden, werden sie untereinander mit den normalen Sortierregeln sortiert.
