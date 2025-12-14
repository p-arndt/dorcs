---
title: "Getting Started"
description: "Get your documentation site running in minutes."
tags: [getting-started, quickstart]
date: 2025-12-13
draft: false
---

# Getting Started

Get your documentation site running in under 5 minutes.

## Quick Start

Get your documentation site running in under 5 minutes. Choose your preferred method:

### Option 1: Docker Compose (Recommended)

```bash
git clone https://github.com/p-arndt/dorcs.git
cd dorcs
docker-compose up
```

### Option 2: Pre-built Binary

1. Download the latest release from the [releases page](https://github.com/p-arndt/dorcs/releases)
2. Make it executable (Linux/macOS): `chmod +x dorcs`
3. Run: `./dorcs --dir ./docs`

### Option 3: Build from Source

```bash
git clone https://github.com/p-arndt/dorcs.git
cd dorcs
go build -o dorcs ./cmd/dorcs
./dorcs --dir ./docs
```

Visit `http://localhost:8080` to see your documentation!

## Your First Documentation

> [!TIP]
> Enable watch mode (`--watch`) for live reload during development

1. **Create your docs directory** (or use the existing `./docs` folder)

2. **Add your first markdown file**:

```bash
cat > docs/index.md << 'EOF'
---
title: "Welcome"
description: "My documentation site"
---

# Welcome

This is my documentation site powered by dorcs!
EOF
```

3. **Start the server**:

```bash
# With Docker Compose
docker-compose up

# Or with binary
./dorcs --dir ./docs
```

4. **Open your browser**: Navigate to `http://localhost:8080`

## File Structure

> [!TIP]
> Put `index.md` inside folders to create section landing pages

dorcs uses **extensionless URLs** that map directly to your file structure:

| File Path                    | URL                   |
| ---------------------------- | --------------------- |
| `docs/index.md`              | `/`                   |
| `docs/getting-started.md`    | `/getting-started`    |
| `docs/guide/index.md`        | `/guide`              |
| `docs/guide/installation.md` | `/guide/installation` |

## Front Matter

Add metadata to your markdown files using YAML front matter:

```yaml
---
title: "Getting Started"
description: "Learn how to get started"
date: 2025-12-13
tags: [tutorial, beginner]
draft: false
---
# Getting Started

Your content here...
```

**Available fields:**

- `title` - Page title (used in navigation and browser tab)
- `description` - Meta description for SEO
- `date` - Publication date (YYYY-MM-DD)
- `tags` - List of tags for categorization
- `draft` - Set to `true` to hide from navigation (when using `--no-drafts`)

## Watch Mode

When using the `--watch` flag, dorcs will watch for changes in your docs directory and automatically reload the page.

**With binary:**

```bash
./dorcs --dir ./docs --watch
```

**With Docker Compose:**

> [!IMPORTANT]
> File watching through Docker volumes may not work reliably on Windows/Mac. For best results, run dorcs directly on your host: `./dorcs --dir ./docs --watch`

Uncomment the `command` line in `docker-compose.yml`:

```yaml
command: ["/app/dorcs", "--dir", "/docs", "--addr", "0.0.0.0:8080", "--watch"]
```

## Common Commands

**Binary:**

```bash
./dorcs --dir ./docs                    # Basic usage
./dorcs --dir ./docs --addr :3000       # Custom port
./dorcs --dir ./docs --watch            # Development mode
./dorcs --dir ./docs --theme midnight   # Custom theme
```

## Next Steps

- ⚙️ [Configuration](./03_configuration.md) - Customize your site with `dorcs.yaml`
- 🎨 [Themes](./05_themes.md) - Choose from built-in themes
- 🚀 [Deployment](./04_deployment.md) - Deploy to production

## Troubleshooting

**Port already in use?**

```bash
./dorcs --dir ./docs --addr :8081
```

**Files not showing up?**

- Ensure files have `.md` or `.markdown` extension
- Check that files aren't marked as `draft: true` (unless using `--no-drafts=false`)
- Verify the `--dir` path is correct
