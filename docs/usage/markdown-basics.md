---
title: "Markdown Basics"
description: "Standard Markdown formatting that Dorcs supports."
tags: [markdown, basics]
date: 2026-03-18
draft: false
---

# Markdown Basics

If you've written Markdown before, you already know this. Here's a quick reference.

## Headings

```markdown
# Page Title
## Section
### Subsection
#### Sub-subsection
```

Dorcs builds the **table of contents** on the right side of the page from your headings. You can see it right now on this page.

## Text formatting

```markdown
**bold text**
*italic text*
~~strikethrough~~
`inline code`
```

## Lists

```markdown
- First item
- Second item
  - Nested item

1. Step one
2. Step two
3. Step three

- [x] Completed task
- [ ] Incomplete task
```

## Links

```markdown
[Another page](./other-page.md)
[External site](https://example.com)
```

Internal `.md` links are automatically converted to clean URLs.

## Images

```markdown
![Description](./images/screenshot.png)
```

## Code blocks

Wrap code in triple backticks. Add the language name for syntax highlighting:

````markdown
```python
def hello():
    print("Hello, world!")
```
````

```python
def hello():
    print("Hello, world!")
```

> [!TIP]
> Every code block gets a **copy button** in the top-right corner. Hover over the block above to see it.

## Tables

```markdown
| Name    | Role     |
| ------- | -------- |
| Alice   | Engineer |
| Bob     | Designer |
```

| Name    | Role     |
| ------- | -------- |
| Alice   | Engineer |
| Bob     | Designer |

## Blockquotes

```markdown
> This is a blockquote. Great for citing sources or adding aside text.
```

> This is a blockquote. Great for citing sources or adding aside text.

## Horizontal rules

```markdown
---
```

---

## Footnotes

```markdown
Dorcs supports footnotes for references and asides.[^1]

[^1]: Like this one — it appears at the bottom of the page.
```

Dorcs supports footnotes for references and asides.[^1]

[^1]: Like this one — it appears at the bottom of the page.

## What's next

Standard Markdown is just the start. Dorcs adds [callouts](./callouts.md) for highlighted messages and [extensions](./extensions.md) for tabs, badges, diagrams, and more.
