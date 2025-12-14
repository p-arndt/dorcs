---
title: "Configuration"
description: "Complete guide to configuring dorcs with dorcs.yaml and command-line options."
tags: [configuration, settings, customization]
date: 2025-12-13
draft: false
---

# Configuration

Customize dorcs using a configuration file or command-line flags.

## Configuration File

Create a `dorcs.yaml` (or `dorcs.yml` or `dorcs.json`) file. The server automatically detects it.

> [!TIP]
> Use YAML for the configuration file as it is the most human-readable format.

### File Location

dorcs looks for the configuration file in this order:
1. **Current working directory** - `./dorcs.yaml`
2. **Docs directory** - `./docs/dorcs.yaml` (if different from working directory)
3. **Custom path** - Use `--config /path/to/config.yaml` to specify an exact path

## Configuration Sections

### Site

```yaml
site:
  title: "My Docs"              # Site title (shown in header)
  description: "..."            # Meta description
  logo: "/static/logo.png"      # Logo image URL (optional)
  favicon: "/static/favicon.ico" # Custom favicon (optional)
```

### Theme

**Presets:** `default`, `ocean`, `forest`, `sunset`, `midnight`, `lavender`, `rose`, and more. See [Themes](./05_themes.md) for all options.

```yaml
theme:
  preset: midnight              # Theme preset
  mode: auto                    # light, dark, or auto
  custom_css: "custom.css"      # Custom CSS (relative to docs dir)
  font_family: '"Inter", system-ui, sans-serif'
  mono_font_family: '"JetBrains Mono", monospace'
```

**Custom Colors:**
```yaml
theme:
  preset: default
  colors:
    light:
      background: "#ffffff"
      foreground: "#1f2328"
      accent: "#0969da"
      # ... more color options
    dark:
      background: "#0d1117"
      foreground: "#e6edf3"
      accent: "#2f81f7"
      # ... more color options
```

**Note:** Syntax highlighting is automatically determined by the preset theme.

### Navigation

```yaml
nav:
  show_search: true             # Enable/disable search box
  expand_all: false             # Expand all folders by default
  links:                        # Header navigation links
    - title: "GitHub"
      url: "https://github.com/..."
      external: true
      icon: "github"            # github, twitter, discord, external
```

### Footer

```yaml
footer:
  copyright: "© 2024 Your Name"
  text: "Additional footer text"
  show_powered_by: true         # Show "Powered by dorcs"
  links:
    - title: "Privacy"
      url: "/privacy"
```

## Command-Line Flags

All configuration options can be overridden via command-line flags:

| Flag | Description | Example |
|------|-------------|---------|
| `--dir` | Docs directory | `--dir ./docs` |
| `--addr` | Listen address | `--addr :8080` |
| `--base-url` | URL path prefix | `--base-url /docs` |
| `--title` | Site title | `--title "My Docs"` |
| `--config` | Config file path | `--config /path/to/config.yaml` |
| `--theme` | Theme preset | `--theme midnight` |
| `--theme-mode` | Theme mode | `--theme-mode dark` |
| `--cache` | Enable caching | `--cache=true` |
| `--no-drafts` | Hide drafts | `--no-drafts=true` |
| `--watch` | Watch mode | `--watch` |

**Note:** Command-line flags override configuration file settings.

## Examples

### Minimal

```yaml
site:
  title: "My Docs"
```

### Full Customization

```yaml
site:
  title: "My Awesome Documentation"
  description: "Complete guide to my project"
  logo: "/static/logo.svg"
  favicon: "/static/favicon.ico"

theme:
  mode: auto
  preset: midnight
  custom_css: "overrides.css"
  font_family: '"SF Pro Display", system-ui, sans-serif'
  mono_font_family: '"SF Mono", monospace'

nav:
  show_search: true
  expand_all: true
  links:
    - title: "GitHub"
      url: "https://github.com/user/repo"
      external: true
      icon: "github"
    - title: "Discord"
      url: "https://discord.gg/server"
      external: true
      icon: "discord"

footer:
  copyright: "© 2024 My Company"
  text: "Made with ❤️ using dorcs"
  show_powered_by: false
  links:
    - title: "Privacy"
      url: "/privacy"
    - title: "Terms"
      url: "/terms"
```

## Next Steps

- 🎨 [Themes](./05_themes.md) - Browse all available themes
- 🚀 [Deployment](./04_deployment.md) - Deploy to production
