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

### Server-Port

Configure the port on which dorcs should run:

```yaml
port: 8000  # Optional: specify the port on which dorcs should run (default: 8080)
```

> [!IMPORTANT]
> The `--addr` command-line flag overrides the port setting in the configuration file if provided.

## Configuration Sections

### Site

```yaml
site:
  title: "My Docs"              # Site title (shown in header)
  description: "..."            # Meta description
  logo: "/logo.png"             # Logo image URL (optional)
                                 # Place logo.png in the root directory where dorcs is running
  favicon: "/favicon.ico"       # Custom favicon (optional)
                                 # Place favicon.ico in the root directory where dorcs is running
```

> [!NOTE]
> **Static Files Location**: Logo and favicon files should be placed in the **root directory** where you run the dorcs executable (the current working directory). The server will first check the root directory for static assets, then fall back to the docs directory. Files are accessible at their path (e.g., `/logo.png`). If you're using a `--base-url` prefix, the BasePath will be automatically prepended to these URLs.

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

### Multi-lingual Support

Configure multiple languages for your documentation:

```yaml
languages:
  default: "en"                 # Default language code (served at root URL)
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
    - code: "fr"
      name: "Français"
```

**How it works:**

- **Default language**: Served at root URL (`/`) - place docs in `docs/` folder
- **Other languages**: Served at `/{lang}/` - place docs in `docs/__lang__/{lang}/` folders
- **Language switcher**: Automatically appears in header when multiple languages are enabled
- **URL structure**:
  - Default: `/getting-started` → `docs/getting-started.md`
  - German: `/de/getting-started` → `docs/__lang__/de/getting-started.md`
  - French: `/fr/getting-started` → `docs/__lang__/fr/getting-started.md`

**File Structure Example:**

```
docs/
  index.md                    # Default language (English)
  getting-started.md
  guide/
    installation.md
  __lang__/                   # Language-specific folder
    de/                       # German language folder
      index.md
      getting-started.md
      guide/
        installation.md
    fr/                       # French language folder
      index.md
      getting-started.md
      guide/
        installation.md
```

**Notes:**

- Each language folder should have a complete copy of your documentation structure
- The `__lang__` folder keeps language-specific content separate and avoids conflicts with regular folders
- The default language stays in the root `docs/` folder
- Other languages go in `docs/__lang__/{lang}/` folders where `{lang}` is the language code
- The language switcher preserves the current page path when switching languages
- If only one language is configured (or none), the language switcher is hidden
- Language codes should follow ISO 639-1 standard (e.g., "en", "de", "fr", "es", "ja")

### Authentication & Edit Mode

Enable online editing with username/password authentication:

```yaml
auth:
  enabled: true                  # Enable authentication and edit mode
  username: "admin"              # Username for login
  password: "your-secure-password"  # Password (will be hashed automatically)
  sessions_path: ".dorcs_sessions.json"  # Optional: custom sessions file path
```

**How it works:**

- When enabled, a **Login** button appears in the footer
- After logging in, **Edit** and **Logout** buttons appear in the header
- Click **Edit** to open the edit mode panel where you can:
  - Browse and edit files
  - Create new files and folders
  - Delete files
  - Save changes directly to the filesystem

**Security Notes:**

- Passwords are automatically hashed using Argon2id on first use
- Sessions expire after 24 hours
- All edit operations require authentication
- The password hash is saved to the config file after first login

> [!WARNING]
> Only enable edit mode on trusted networks or behind proper authentication. The edit mode allows full file system access to your docs directory.

### GitHub (External Content)

Serve markdown files directly from a GitHub repository instead of local files. When enabled, local files in the docs directory are ignored and only GitHub files are served.

```yaml
github:
  enabled: true                                    # Enable GitHub content source
  repository: "https://github.com/owner/repo/tree/main/docs"  # Repository tree URL
  token: "${GITHUB_TOKEN}"                         # Optional: GitHub token (recommended)
  cache_ttl: "1h"                                 # Optional: Cache TTL (default: 1h)
```

> [!TIP]
> For detailed documentation on using GitHub as a content source, including setup instructions, troubleshooting, and best practices, see the [GitHub guide](./external-content/github.md).

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
  logo: "/logo.svg"        # Place logo.svg in the root directory where dorcs is running
  favicon: "/favicon.ico"  # Place favicon.ico in the root directory where dorcs is running

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

languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
    - code: "fr"
      name: "Français"

auth:
  enabled: true
  username: "admin"
  password: "secure-password-here"

github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: "${GITHUB_TOKEN}"  # Recommended: use token to avoid rate limiting
  cache_ttl: "1h"
```

> [!NOTE]
> See the [GitHub guide](./external-content/github.md) for detailed setup and usage instructions.
```

## Next Steps

- 🎨 [Themes](./05_themes.md) - Browse all available themes
- 🚀 [Deployment](./04_deployment.md) - Deploy to production