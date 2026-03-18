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
  enabled:                      # Explicit configuration (required)
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"
    - code: "fr"
      name: "Français"
```

**How it works:**

- **Default language**: Defined in `dorcs.yaml` - files go in `docs/{default-lang}/` folder, served at root URL (`/`)
- **Other languages**: Served at `/{lang}/` - place docs in `docs/{lang}/` folders
- **Language switcher**: Automatically appears in header when multiple languages are enabled
- **Explicit configuration required**: Languages must be explicitly configured in `enabled` list
- **URL structure**:
  - Default (English): `/getting-started` → `docs/en/getting-started.md`
  - German: `/de/getting-started` → `docs/de/getting-started.md`
  - French: `/fr/getting-started` → `docs/fr/getting-started.md`

**File Structure Example (MkDocs-style):**

```
docs/
  en/                         # English (default) - served at /
    index.md
    getting-started.md
    guide/
      installation.md
  de/                         # German - served at /de/
    index.md
    getting-started.md
    guide/
      installation.md
  fr/                         # French - served at /fr/
    index.md
    getting-started.md
    guide/
      installation.md
```

**Notes:**

- **MkDocs-style structure**: Direct language folders
- **Markdown files**: When using multiple languages, markdown files in root `docs/` are **ignored**. All markdown files must be in language folders
- **Shared assets**: Static assets (images, PDFs, etc.) in root `docs/` are served as **language-independent shared assets** accessible from all languages
- **Default language**: Defined in `dorcs.yaml` - default language files are in `docs/en/` but served at `/` (no prefix)
- **Relative paths**: Automatically handled - images and links work correctly even when default language uses its folder
- Each language folder should have a complete copy of your documentation structure
- The language switcher preserves the current page path when switching languages
- If only one language is configured (or none), the language switcher is hidden
- Language codes should follow ISO 639-1 standard (e.g., "en", "de", "fr", "es", "ja")
- See [File Structure & Organization](./usage/file-structure.md) for detailed examples

> [!WARNING]
> **Markdown files in root are ignored**: When using multiple languages, any `.md` files in the root `docs/` folder are ignored. Only static assets (images, PDFs, etc.) in root are served as shared assets accessible from all languages.

### Versioning

Configure multiple versions of your documentation:

```yaml
versions:
  default: "latest"            # Default version identifier (served at root URL)
  enabled:                     # Explicit configuration (required)
    - id: "latest"
      name: "Latest"
    - id: "v1"
      name: "Version 1"
    - id: "v2"
      name: "Version 2"
```

**How it works:**

- **Default version**: Served at root URL (`/`) - can be in `docs/` or `docs/{version}/` folder
- **Other versions**: Served at `/{version}/` - place docs in `docs/{version}/` folders
- **Version switcher**: Automatically appears in header when multiple versions are enabled
- **Explicit configuration required**: Versions must be explicitly configured in `enabled` list
- **URL structure**:
  - Default: `/getting-started` → `docs/getting-started.md` or `docs/latest/getting-started.md`
  - v1: `/v1/getting-started` → `docs/v1/getting-started.md`
  - v2: `/v2/getting-started` → `docs/v2/getting-started.md`

**File Structure Example (MkDocs-style):**

```
docs/
  index.md                    # Default version (latest) - served at /
  getting-started.md
  latest/                     # Latest (explicit) - served at /latest/
    index.md
  v1/                         # Version 1 - served at /v1/
    index.md
    getting-started.md
  v2/                         # Version 2 - served at /v2/
    index.md
    getting-started.md
```

**Combined with Languages (Language-First):**

```
docs/
  en/                         # English
    index.md                  # en, latest (default) - served at /en/
    v1/                       # en, v1 - served at /en/v1/
      index.md
  de/                         # German
    index.md                  # de, latest (default) - served at /de/
    v1/                       # de, v1 - served at /de/v1/
      index.md
```

**URLs with both languages and versions:**
- `/en/v1/getting-started` → `docs/en/v1/getting-started.md`
- `/de/v1/getting-started` → `docs/de/v1/getting-started.md`

**Notes:**

- **MkDocs-style structure**: Direct version folders
- **Version identifiers**: Can be any string (e.g., `v1`, `v2`, `1.0`, `latest`, `stable`)
- **Language-first**: When using both, structure is `docs/{lang}/{version}/`
- **Relative paths**: Automatically handled for all versions
- See [File Structure & Organization](./usage/file-structure.md) for detailed examples

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

With a GitHub token, dorcs can also serve documentation from private repositories, including linked assets such as images and PDFs.

```yaml
github:
  enabled: true                                    # Enable GitHub content source
  repository: "https://github.com/owner/repo/tree/main/docs"  # Repository tree URL
  token: "${GITHUB_TOKEN}"                         # Optional: GitHub token (recommended)
  cache_ttl: "1h"                                 # Optional: Cache TTL (default: 1h)
  edit_on_github:                                 # Optional: "Edit on GitHub" for local docs
    repository: "https://github.com/owner/repo/tree/main/docs"
```

**Edit on GitHub link:** When docs are sourced from GitHub (`enabled: true`), an "Edit on GitHub" link appears on each page—clicking it opens the source file in GitHub's editor. For local docs (e.g. deployed to GitHub Pages), add `edit_on_github.repository` to show the link.

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