---
title: "Branding"
description: "Set your site's title, logo, footer, and announcements."
tags: [configuration, branding]
date: 2026-03-18
draft: false
---

# Branding

Make your docs site feel like it belongs to your project.

## Site title and description

The title appears in the header, browser tab, and when someone shares a link to your docs on social media:

```yaml
site:
  title: "Acme Platform"
  description: "Developer documentation for the Acme Platform"
```

## Logo and favicon

Replace the text title with your logo, and add a favicon for browser tabs:

```yaml
site:
  logo: "/logo.png"
  favicon: "/favicon.ico"
```

> [!TIP]
> Put these files inside your `docs/` folder. The path starts with `/` because it's relative to the site root, not the config file.

In multi-language setups, shared assets like logos can live in the root `docs/` folder and be referenced from any language subfolder.

## Announcement banner

Need to tell visitors something? Show a banner at the top of every page:

```yaml
announcement:
  text: 'Version 2.0 is out! <a href="/changelog">See what changed</a>'
  dismissible: true
```

The banner supports basic HTML for links and formatting. When `dismissible` is `true` (the default), visitors can close it and it stays hidden — Dorcs saves their preference in the browser.

> [!NOTE]
> Great for release announcements, maintenance notices, or pointing people to a survey.

## Footer

Add copyright info, custom text, or a "Powered by Dorcs" link:

```yaml
footer:
  text: "Built with love by the Acme team"
  copyright: "© 2026 Acme Inc."
  show_powered_by: true
```

## Putting it all together

Here's a complete branding setup:

```yaml
site:
  title: "Acme Platform"
  description: "Developer docs for the Acme Platform"
  logo: "/logo.png"
  favicon: "/favicon.ico"

announcement:
  text: 'New: Check out our <a href="/guides/v2-migration">v2 migration guide</a>'
  dismissible: true

footer:
  copyright: "© 2026 Acme Inc."
  text: "Made with Dorcs"
  show_powered_by: true
```
