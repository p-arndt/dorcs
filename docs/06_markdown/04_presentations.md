---
title: "Presentations"
description: "Use Markdown pages as slide decks."
tags: [markdown, presentations]
date: 2026-03-18
draft: false
---

# Presentations

A normal Markdown page can become a slide deck by setting front matter:

```yaml
---
title: "Quarterly Review"
presentation: true
presentation_header: "Engineering"
presentation_footer: "Q1 2026"
---
```

Separate slides with horizontal-rule style separators:

```markdown
# Slide one

---

# Slide two
```

Supported separators:

- `---`
- `___`
- `***`
- `- - -`

For deck-specific features, continue with the [Presentations](../07_presentations/index.md) section.
