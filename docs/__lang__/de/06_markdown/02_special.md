---
title: "Spezielle Markdown-Funktionen"
description: "Spezielle Markdown-Funktionen"
tags: [markdown, special]
date: 2025-12-14
draft: false
---

# Spezielle Markdown-Funktionen

## Inhaltsverzeichnis

Erstellen Sie ein Inhaltsverzeichnis mit dem Platzhalter:

```markdown
[[TOC]]
```

**Ergebnis:**

[[TOC]]

## Inhaltsverzeichnis Root-Navigation

Erstellen Sie ein Inhaltsverzeichnis für die Root-Navigation mit dem Platzhalter:

```markdown
[[TOC-ROOT]]
```

**Result:**

[[TOC-ROOT]]

## Filtered Table of Contents

Create a table of contents with a maximum depth limit:

```markdown
[[TOC:1]]
```

Dies zeigt nur Überschriften bis zur Ebene 2 (H1 und H2) an. Sie können jede Tiefe von 1 bis 6 angeben.

**Ergebnis:**

[[TOC:1]]

## Warnblöcke

dorcs unterstützt GitHub-Style-Warnblöcke für Hinweise und wichtige Informationen:

```markdown
> [!NOTE]
> This is a note alert block.
> It can contain multiple lines.

> [!TIP]
> This is a tip alert block.
> Use it for helpful suggestions.

> [!IMPORTANT]
> This is an important alert block.
> Use it for critical information.

> [!WARNING]
> This is a warning alert block.
> Use it to warn users about potential issues.

> [!CAUTION]
> This is a caution alert block.
> Use it for dangerous or risky operations.
```

**Result:**

> [!NOTE]
> This is a note alert block.
> It can contain multiple lines.


> [!TIP]
> This is a tip alert block.
> Use it for helpful suggestions.


> [!IMPORTANT]
> This is an important alert block.
> Use it for critical information.


> [!WARNING]
> This is a warning alert block.
> Use it to warn users about potential issues.


> [!CAUTION]
> This is a caution alert block.
> Use it for dangerous or risky operations.




## Breadcrumb-Navigation

Zeigen Sie Breadcrumb-Navigation an, die den Pfad zur aktuellen Seite zeigt:

```markdown
[[BREADCRUMBS]]
```

**Ergebnis:**

[[BREADCRUMBS]]

## Unterseiten

Zeigen Sie alle Unterseiten und Ordner im aktuellen Abschnitt an:

```markdown
[[CHILDREN]]
```

Dies zeigt alle direkten Kinder (Seiten und Unterordner) der aktuellen Seite/des aktuellen Ordners an, mit Beschreibungen und Indikatoren für Verzeichnisse.

**Ergebnis:**

[[CHILDREN]]

## Geschwister-Seiten

Zeigen Sie Links zu Geschwister-Seiten (Seiten auf derselben Ebene) an:

```markdown
[[SIBLINGS]]
```

**Ergebnis:**

[[SIBLINGS]]

## Verwandte Seiten

Zeigen Sie verwandte Seiten basierend auf gemeinsamen Tags an:

```markdown
[[RELATED]]
```

**Ergebnis:**

[[RELATED]]

## Kürzlich aktualisierte Seiten

Zeigen Sie eine Liste kürzlich aktualisierter Seiten an:

```markdown
[[RECENT]]
```

**Ergebnis:**

[[RECENT]]

## Seiten-Tags

Zeigen Sie Tags für die aktuelle Seite an:

```markdown
[[TAGS]]
```

**Result:**

[[TAGS]]

## Site Index

Display a full alphabetical index of all pages:

```markdown
[[INDEX]]
```

**Ergebnis:**

[[INDEX]]

## Veröffentlichungsdatum

Zeigen Sie das Veröffentlichungsdatum aus dem Front Matter an:

```markdown
[[DATE]]
```

or

```markdown
[[PUBLISHED]]
```

**Result:**

[[DATE]]

## Last Updated Date

Display the last modified date of the page:

```markdown
[[LAST-UPDATED]]
```

**Ergebnis:**

[[LAST-UPDATED]]

## Autoreninformationen

Zeigen Sie Autoreninformationen aus dem Front Matter an:

```markdown
[[AUTHOR]]
```

**Ergebnis:**

[[AUTHOR]]

> [!NOTE]
> Fügen Sie `author: "Autorenname"` zu Ihrem Front Matter hinzu, um diese Funktion zu verwenden.

## Seiten-Zusammenfassung

Zeigen Sie die Seiten-Zusammenfassung an (aus Beschreibung oder erstem Absatz):

```markdown
[[SUMMARY]]
```

**Ergebnis:**

[[SUMMARY]]

## Seiten nach Tag

Zeigen Sie alle Seiten gruppiert nach ihren Tags an:

```markdown
[[PAGES-BY-TAG]]
```

**Ergebnis:**

[[PAGES-BY-TAG]]