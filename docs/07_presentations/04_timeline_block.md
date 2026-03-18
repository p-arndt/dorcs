---
title: "Timeline Block"
description: "Timeline block for slide decks"
tags: [presentations, slides, markdown]
date: 2026-02-21
---

# Timeline Block

Use the `::: timeline` block for roadmaps and process steps. Each `###` or `####` heading starts a new step; the heading text becomes the date/label marker.

```markdown
<!-- _layout: timeline -->

# Project Roadmap

::: timeline

### 2024 · Q1

**Discovery**
Initial research and planning.

### 2024 · Q2

**Build**
Development and testing.

### 2024 · Q3

**Launch**
Release and iterate.
:::
```
