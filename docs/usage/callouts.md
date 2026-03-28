---
title: "Callouts"
description: "Highlight important information with styled callout boxes."
tags: [markdown, callouts]
date: 2026-03-18
draft: false
---

# Callouts

Callouts draw attention to important information. Dorcs uses the same syntax as GitHub:

```markdown
> [!NOTE]
> Background information the reader should know.
```

## All five types

> [!NOTE]
> Useful background information that supplements the main content.

> [!TIP]
> Helpful advice to make things easier or better.

> [!IMPORTANT]
> Key information the reader needs to be aware of.

> [!WARNING]
> Something that could cause problems if ignored.

> [!CAUTION]
> Something that could cause data loss or a serious issue.

## How to write them

Each callout starts with `> [!TYPE]` on its own line, followed by the content as a blockquote:

```markdown
> [!TIP]
> You can write **multiple lines** in a callout.
> Just keep using the `>` prefix.
>
> You can even have paragraphs, code blocks, and lists inside.
```

> [!TIP]
> You can write **multiple lines** in a callout.
> Just keep using the `>` prefix.
>
> You can even have paragraphs, code blocks, and lists inside.

## When to use which

| Type | Use for |
| --- | --- |
| NOTE | "By the way" information, context that helps understanding |
| TIP | Shortcuts, best practices, things that save time |
| IMPORTANT | Must-know info that affects how something works |
| WARNING | Things that could go wrong or cause unexpected behavior |
| CAUTION | Destructive actions, data loss risks, security concerns |

> [!TIP]
> Don't overuse callouts — if everything is highlighted, nothing stands out. Save them for information that truly deserves extra attention.
