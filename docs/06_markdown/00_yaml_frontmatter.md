---
title: "YAML Front Matter"
description: "YAML front matter features"
tags: [yaml, frontmatter]
date: 2025-12-14
draft: false
---
# YAML Front Matter

Add metadata to your markdown files using YAML front matter at the top of the file:

```yaml
---
title: "Page Title"
description: "Page description for SEO"
date: 2025-12-13
tags: [tag1, tag2, tag3]
draft: false
order: 1
---
```

**Available fields:**

- `title` - Page title (used in navigation and browser tab)
- `description` - Meta description for SEO
- `date` - Publication date (YYYY-MM-DD format)
- `tags` - List of tags for categorization
- `draft` - Set to `true` to hide from navigation (when using `--no-drafts`)
- `order` - Numeric value for sorting pages in navigation (lower numbers appear first)
- `author` - Author name (displayed with [[AUTHOR]] placeholder)

**Example:**

```yaml
---
title: "Getting Started"
description: "Learn how to get started with dorcs"
date: 2025-12-13
tags: [tutorial, beginner, quickstart]
draft: false
order: 1
---

# Getting Started

Your content here...
```