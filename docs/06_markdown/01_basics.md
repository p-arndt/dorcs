---
title: "Markdown Basics"
description: "Basic markdown features"
tags: [markdown, basics]
date: 2025-12-14
draft: false
---

# Markdown Basics

## Headings

Create headings using `#` symbols:

```markdown
# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6
```

**Result:**

## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6

### Emphasis

```markdown
*italic text*
_italic text_

**bold text**
__bold text__

***bold and italic***
___bold and italic___
```

**Result:**

*italic text*
_italic text_

**bold text**
__bold text__

***bold and italic***
___bold and italic___

### Lists

#### Unordered Lists

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

**Result:**

- Item 1
- Item 2
  - Nested item 2.1
  - Nested item 2.2
- Item 3

* Alternative bullet
* Another item

+ Another style
+ One more

#### Ordered Lists

```markdown
1. First item
2. Second item
3. Third item
   1. Nested item
   2. Another nested item
4. Fourth item
```

**Result:**

1. First item
2. Second item
3. Third item
   1. Nested item
   2. Another nested item
4. Fourth item

### Links

```md
[Link text](https://example.com)
[Link with title](https://example.com "Title text")
[Anchor link](#heading-id)

<https://example.com>
<email@example.com>
```

**Result:**

[Link text](https://example.com)
[Link with title](https://example.com "Title text")
[Anchor link](#heading-id)

<https://example.com>
<email@example.com>

### Images

```markdown
![Alt text](../logo.png)
![Alt text with title](../logo.png "Logo title")
![Alt text](https://example.com/image.png)
```

**Result:**

![Alt text](../logo.png)

### Blockquotes

```markdown
> This is a blockquote.
> It can span multiple lines.
>
> And include multiple paragraphs.

> Nested blockquote
>> Even deeper nesting
```

**Result:**

> This is a blockquote.
> It can span multiple lines.
>
> And include multiple paragraphs.

> Nested blockquote
>> Even deeper nesting

### Inline Code

```markdown
Use `inline code` in your text.
Use ``code with `backticks` inside`` for escaping.
```

**Result:**

Use `inline code` in your text.
Use ``code with `backticks` inside`` for escaping.

### Code Blocks

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

### Tables

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

**Result:**

~~This text is struck through~~
~~Multiple words~~ can be struck through.

### Task Lists

```markdown
- [x] Completed task
- [x] Another completed task
- [ ] Incomplete task
- [ ] Another incomplete task
  - [x] Nested completed task
  - [ ] Nested incomplete task
```

**Result:**

- [x] Completed task
- [x] Another completed task
- [ ] Incomplete task
- [ ] Another incomplete task
  - [x] Nested completed task
  - [ ] Nested incomplete task

### Autolinks

```markdown
https://example.com
www.example.com
email@example.com
```

**Result:**

https://example.com
www.example.com
email@example.com


## Code Blocks & Syntax Highlighting

dorcs uses Chroma for syntax highlighting with support for 200+ languages. Specify the language after the opening code fence:

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
> The syntax highlighting theme can be configured in `dorcs.yaml` using the `code_theme` option. See [Configuration](./../../03_configuration.md) for details.

## Footnotes

Create footnotes using the footnote syntax:

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


## HTML Support

dorcs allows embedded HTML in markdown files for advanced formatting:

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
> While HTML is supported, use it sparingly. Prefer markdown syntax when possible for better portability and consistency.

## Auto Heading IDs

All headings automatically get ID attributes based on their text, making it easy to create anchor links:

```markdown
## My Section Heading
```

This heading automatically gets the ID `my-section-heading`, which you can link to:

```markdown
[Link to section](#my-section-heading)
```

**Result:**

[Link to section](#my-section-heading)

The table of contents (TOC) is automatically generated from all headings on the page and uses these IDs for navigation.

## Typographer

The typographer extension automatically converts plain text to typographically correct characters:

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

## Combining Features

You can combine multiple features for rich formatting:

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

## LaTex Support

dorcs supports LaTex math rendering using [KaTeX](https://katex.org/). Simply enclose your math in `$$` tags:

```markdown
$$
f(x) = \int_{-\infty}^\infty
\hat f(\xi)\,e^{2 \pi i \xi x}
\,d\xi
$$
```

**Result:**

$$
f(x) = \int_{-\infty}^\infty \hat f(\xi)\,e^{2 \pi i \xi x} \,d\xi
$$

Or simply use inline math:

```markdown
A pi function $f(x) = x \cdot \pi$
```

**Result:**

A pi function $f(x) = x \cdot \pi$  
