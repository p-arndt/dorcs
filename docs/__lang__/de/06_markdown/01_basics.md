---
title: "Markdown-Grundlagen"
description: "Grundlegende Markdown-Funktionen"
tags: [markdown, basics]
date: 2025-12-14
draft: false
---

# Markdown-Grundlagen

## Überschriften

Erstellen Sie Überschriften mit `#` Symbolen:

```markdown
# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6
```

**Ergebnis:**

## Überschrift 2
### Überschrift 3
#### Überschrift 4
##### Überschrift 5
###### Überschrift 6

### Hervorhebung

```markdown
*italic text*
_italic text_

**bold text**
__bold text__

***bold and italic***
___bold and italic___
```

**Ergebnis:**

*kursiver Text*
_kursiver Text_

**fetter Text**
__fetter Text__

***fett und kursiv***
___fett und kursiv___

### Listen

#### Ungeordnete Listen

```markdown
- Item 1
- Item 2
  - Nested item 2.1
  - Nested item 2.2
- Item 3

* Alternative bullet
* Another item

+ Another style
+ One more
```

**Ergebnis:**

- Element 1
- Element 2
  - Verschachteltes Element 2.1
  - Verschachteltes Element 2.2
- Element 3

* Alternative Aufzählung
* Weitere Elemente

+ Weitere Stil
+ Noch eines

#### Geordnete Listen

```markdown
1. First item
2. Second item
3. Third item
   1. Nested item
   2. Another nested item
4. Fourth item
```

**Ergebnis:**

1. Erstes Element
2. Zweites Element
3. Drittes Element
   1. Verschachteltes Element
   2. Weitere verschachtelte Elemente
4. Viertes Element

### Links

```md
[Link text](https://example.com)
[Link with title](https://example.com "Title text")
[Anchor link](#heading-id)

<https://example.com>
<email@example.com>
```

**Ergebnis:**

[Link-Text](https://example.com)
[Link mit Titel](https://example.com "Titel-Text")
[Anker-Link](#heading-id)

<https://example.com>
<email@example.com>

### Bilder

```markdown
![Alt text](../logo.png)
![Alt text with title](../logo.png "Logo title")
![Alt text](https://example.com/image.png)
```

**Ergebnis:**

![Alt-Text](../logo.png)

### Blockzitate

```markdown
> This is a blockquote.
> It can span multiple lines.
>
> And include multiple paragraphs.

> Nested blockquote
>> Even deeper nesting
```

**Ergebnis:**

> Dies ist ein Blockzitat.
> Es kann mehrere Zeilen umfassen.
>
> Und mehrere Absätze enthalten.

> Verschachteltes Blockzitat
>> Noch tiefere Verschachtelung

### Inline-Code

```markdown
Use `inline code` in your text.
Use ``code with `backticks` inside`` for escaping.
```

**Ergebnis:**

Verwenden Sie `Inline-Code` in Ihrem Text.
Verwenden Sie ``Code mit `Backticks` darin`` zum Escapen.

### Code-Blöcke

```markdown
```
Plain code block
with multiple lines
```
```

**Result:**

```
Plain code block
with multiple lines
```

### Horizontal Rules

```markdown
---

***

___
```

**Result:**

---

***

___

## GitHub Flavored Markdown (GFM)

### Tabellen

```markdown
| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |

| Left Aligned | Center Aligned | Right Aligned |
|:-------------|:--------------:|--------------:|
| Left         | Center         | Right         |
| More left    | More center    | More right    |
```

**Result:**

| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |

| Left Aligned | Center Aligned | Right Aligned |
|:-------------|:--------------:|--------------:|
| Left         | Center         | Right         |
| More left    | More center    | More right    |

### Strikethrough

```markdown
~~This text is struck through~~
~~Multiple words~~ can be struck through.
```

**Ergebnis:**

~~Dieser Text ist durchgestrichen~~
~~Mehrere Wörter~~ können durchgestrichen werden.

### Aufgabenlisten

```markdown
- [x] Completed task
- [x] Another completed task
- [ ] Incomplete task
- [ ] Another incomplete task
  - [x] Nested completed task
  - [ ] Nested incomplete task
```

**Ergebnis:**

- [x] Abgeschlossene Aufgabe
- [x] Weitere abgeschlossene Aufgabe
- [ ] Unvollständige Aufgabe
- [ ] Weitere unvollständige Aufgabe
  - [x] Verschachtelte abgeschlossene Aufgabe
  - [ ] Verschachtelte unvollständige Aufgabe

### Auto-Links

```markdown
https://example.com
www.example.com
email@example.com
```

**Result:**

https://example.com
www.example.com
email@example.com


## Code-Blöcke & Syntax-Hervorhebung

dorcs verwendet Chroma für Syntax-Hervorhebung mit Unterstützung für über 200 Sprachen. Geben Sie die Sprache nach dem öffnenden Code-Fence an:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

```python
def hello_world():
    print("Hello, World!")
    return True
```

```javascript
function helloWorld() {
    console.log("Hello, World!");
    return true;
}
```

```bash
#!/bin/bash
echo "Hello, World!"
```

```yaml
name: example
version: 1.0.0
features:
  - syntax highlighting
  - multiple languages
```

```json
{
  "name": "example",
  "version": "1.0.0",
  "features": ["syntax", "highlighting"]
}
```

```html
<!DOCTYPE html>
<html>
<head>
    <title>Example</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>
```

```css
body {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 20px;
}
```

```sql
SELECT name, email
FROM users
WHERE active = true
ORDER BY created_at DESC;
```


**Result:**

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```


```python
def hello_world():
    print("Hello, World!")
    return True
```

```javascript
function helloWorld() {
    console.log("Hello, World!");
    return true;
}
```

```bash
#!/bin/bash
echo "Hello, World!"
```

```yaml
name: example
version: 1.0.0
features:
  - syntax highlighting
  - multiple languages
```

```json
{
  "name": "example",
  "version": "1.0.0",
  "features": ["syntax", "highlighting"]
}
```

```html
<!DOCTYPE html>
<html>
<head>
    <title>Example</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>
```

```css
body {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 20px;
}
```

```sql
SELECT name, email
FROM users
WHERE active = true
ORDER BY created_at DESC;
```

> [!TIP]
> Das Syntax-Hervorhebungs-Theme kann in `dorcs.yaml` mit der Option `code_theme` konfiguriert werden. Siehe [Konfiguration](./../../03_configuration.md) für Details.

## Fußnoten

Erstellen Sie Fußnoten mit der Fußnoten-Syntax:

```markdown
Here's a sentence with a footnote[^1].

Another sentence with multiple footnotes[^2][^3].

[^1]: This is the first footnote.
[^2]: This is the second footnote.
[^3]: This is the third footnote with **bold text** and `code`.
```

**Result:**

Here's a sentence with a footnote[^1].

Another sentence with multiple footnotes[^2][^3].

[^1]: This is the first footnote.
[^2]: This is the second footnote.
[^3]: This is the third footnote with **bold text** and `code`.


## HTML-Unterstützung

dorcs erlaubt eingebettetes HTML in Markdown-Dateien für erweiterte Formatierung:

```markdown
<div style="text-align: center;">
  <h2>Centered Heading</h2>
</div>

<details>
  <summary>Click to expand</summary>
  This is hidden content that can be expanded.
</details>

<kbd>Ctrl</kbd> + <kbd>C</kbd> to copy

<mark>Highlighted text</mark>

<sub>Subscript</sub> and <sup>Superscript</sup>
```

**Result:**

<div style="text-align: center;">
  <h2>Centered Heading</h2>
</div>

<details>
  <summary>Click to expand</summary>
  This is hidden content that can be expanded.
</details>

<kbd>Ctrl</kbd> + <kbd>C</kbd> to copy

<mark>Highlighted text</mark>

<sub>Subscript</sub> and <sup>Superscript</sup>

> [!WARNING]
> Während HTML unterstützt wird, verwenden Sie es sparsam. Bevorzugen Sie Markdown-Syntax, wenn möglich, für bessere Portabilität und Konsistenz.

## Automatische Überschriften-IDs

Alle Überschriften erhalten automatisch ID-Attribute basierend auf ihrem Text, was das Erstellen von Anker-Links erleichtert:

```markdown
## My Section Heading
```

Diese Überschrift erhält automatisch die ID `my-section-heading`, auf die Sie verlinken können:

```markdown
[Link zum Abschnitt](#my-section-heading)
```

**Ergebnis:**

[Link zum Abschnitt](#my-section-heading)

Das Inhaltsverzeichnis (TOC) wird automatisch aus allen Überschriften auf der Seite generiert und verwendet diese IDs für die Navigation.

## Typograf

Die Typograf-Erweiterung konvertiert automatisch einfachen Text in typografisch korrekte Zeichen:

```markdown
"Smart quotes" and 'smart quotes'
-- en dash
--- em dash
... ellipsis
(c) copyright
(r) registered
(tm) trademark
```

**Result:**

"Smart quotes" and 'smart quotes'
-- en dash
--- em dash
... ellipsis
(c) copyright
(r) registered
(tm) trademark

## Hard Wraps

Line breaks in markdown are preserved as hard breaks (like `<br>` tags), so you can format text with line breaks:

```markdown
This is line one.
This is line two.
This is line three.
```

**Result:**

This is line one.
This is line two.
This is line three.

## Funktionen kombinieren

Sie können mehrere Funktionen für reichhaltige Formatierung kombinieren:

```markdown
> [!TIP]
> Use **bold text** and `code` in alert blocks.
> 
> - Lists work too
> - Multiple items
> 
> [Links](https://example.com) are supported.

| Feature | Status | Notes |
|---------|--------|-------|
| Tables | ✅ | Fully supported |
| Code blocks | ✅ | With syntax highlighting |
| Alerts | ✅ | 5 types available |

Check out the [footnote](#footnotes) section[^example] for more details.

[^example]: This is an example footnote in a combined example.
```

**Result:**

> [!TIP]
> Use **bold text** and `code` in alert blocks.
> 
> - Lists work too
> - Multiple items
> 
> [Links](https://example.com) are supported.

| Feature | Status | Notes |
|---------|--------|-------|
| Tables | ✅ | Fully supported |
| Code blocks | ✅ | With syntax highlighting |
| Alerts | ✅ | 5 types available |

Check out the [footnote](#footnotes) section[^example] for more details.

[^example]: This is an example footnote in a combined example.
