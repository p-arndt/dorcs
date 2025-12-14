---
title: "Special Markdown"
description: "Special markdown features"
tags: [markdown, special]
date: 2025-12-14
draft: false
---

# Special Markdown

## Table of Contents

Create a table of contents using the placeholder:

```markdown
[[TOC]]
```

**Result:**

[[TOC]]

## Table of Contents Root Navigation

Create a table of contents for the root navigation using the placeholder:

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

This will only show headings up to level 2 (H1 and H2). You can specify any depth from 1 to 6.

**Result:**

[[TOC:1]]

## Alert Blocks

dorcs supports GitHub-style alert blocks for callouts and important information:

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




## Breadcrumb Navigation

Display breadcrumb navigation showing the path to the current page:

```markdown
[[BREADCRUMBS]]
```

**Result:**

[[BREADCRUMBS]]

## Children Pages

Display all children pages and folders in the current section:

```markdown
[[CHILDREN]]
```

This shows all direct children (pages and subfolders) of the current page/folder, with descriptions and indicators for directories.

**Result:**

[[CHILDREN]]

## Sibling Pages

Display links to sibling pages (pages at the same level):

```markdown
[[SIBLINGS]]
```

**Result:**

[[SIBLINGS]]

## Related Pages

Display related pages based on shared tags:

```markdown
[[RELATED]]
```

**Result:**

[[RELATED]]

## Recently Updated Pages

Display a list of recently updated pages:

```markdown
[[RECENT]]
```

**Result:**

[[RECENT]]

## Page Tags

Display tags for the current page:

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

**Result:**

[[INDEX]]

## Publication Date

Display the publication date from front matter:

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

**Result:**

[[LAST-UPDATED]]

## Author Information

Display author information from front matter:

```markdown
[[AUTHOR]]
```

**Result:**

[[AUTHOR]]

> [!NOTE]
> Add `author: "Author Name"` to your front matter to use this feature.

## Page Summary

Display the page summary (from description or first paragraph):

```markdown
[[SUMMARY]]
```

**Result:**

[[SUMMARY]]

## Pages by Tag

Display all pages grouped by their tags:

```markdown
[[PAGES-BY-TAG]]
```

**Result:**

[[PAGES-BY-TAG]]