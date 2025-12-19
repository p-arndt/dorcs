---
title: "GitHub"
description: "Serve markdown files directly from GitHub repositories with automatic caching"
tags: [github, external, content, caching]
date: 2025-12-18
draft: false
---

# GitHub

Serve markdown files directly from a GitHub repository instead of local files. When enabled, local files in the docs directory are ignored and only GitHub files are served.

## Configuration

Add the following to your `dorcs.yaml`:

```yaml
github:
  enabled: true                                    # Enable GitHub content source
  repository: "https://github.com/owner/repo/tree/main/docs"  # Repository tree URL
  token: "${GITHUB_TOKEN}"                         # Optional: GitHub token (recommended)
```

See the [Configuration](../03_configuration.md#github-external-content) page for more details on configuration options.

## Repository URL Formats

The repository URL can be in several formats:
- `https://github.com/owner/repo/tree/branch/path` - Full GitHub tree URL
- `https://github.com/owner/repo/tree/main/docs` - Specific branch and path
- `https://github.com/owner/repo` - Root of repository (uses default branch)
- `owner/repo/tree/branch/path` - Without `https://github.com/` prefix (added automatically)

## How It Works

- When enabled, dorcs automatically discovers all `.md` files in the specified GitHub directory tree
- Files are fetched from GitHub and cached locally in `.cache/github/` directory
- The URL structure mirrors the GitHub directory structure
- Local files in the docs directory are completely ignored when GitHub mode is enabled

## Example URL Mapping

| GitHub Path                    | Served URL              |
| ------------------------------ | ----------------------- |
| `docs/index.md`                | `/`                     |
| `docs/getting-started.md`      | `/getting-started`      |
| `docs/guide/installation.md`   | `/guide/installation`   |

## Multi-lingual Support

GitHub content source works seamlessly with multi-lingual setups:
- **Default language**: Uses the base repository path (e.g., `docs/`)
- **Other languages**: Automatically uses `docs/__lang__/{lang}/` paths
- Language switching works seamlessly with GitHub-sourced content

### Example with Languages

```yaml
languages:
  default: "en"
  enabled:
    - code: "en"
      name: "English"
    - code: "de"
      name: "Deutsch"

github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
```

This configuration will:
- Serve English files from `docs/` in GitHub
- Serve German files from `docs/__lang__/de/` in GitHub
- Language switcher will work correctly

## GitHub Token (Recommended)

> [!IMPORTANT]
> **Using a GitHub token is highly recommended** to avoid rate limiting issues. Without a token, GitHub API has strict rate limits (60 requests/hour for unauthenticated requests). With a token, you get 5,000 requests/hour.

### Creating a GitHub Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate a new token with `public_repo` scope (for public repos) or `repo` scope (for private repos)
3. Use the token in your config:

```yaml
github:
  enabled: true
  repository: "https://github.com/owner/repo/tree/main/docs"
  token: "${GITHUB_TOKEN}"  # Supports environment variable expansion
```

### Environment Variable Expansion

The token supports environment variable expansion:
- `${GITHUB_TOKEN}` - Uses the `GITHUB_TOKEN` environment variable
- `${VAR:-default}` - Uses `VAR` if set, otherwise uses `default`

## Caching

- Content is cached both in-memory and on disk (in `.cache/github/` directory)
- Cache persists across server restarts
- Default cache TTL is 1 hour (configurable)
- Cache directory is automatically created in the working directory where dorcs runs
- Cache files are automatically cleaned up when expired

### Cache Directory

The cache is stored in `.cache/github/` in the working directory where dorcs is running. This directory is automatically created and is ignored by git (already in `.gitignore`).

## Important Notes

- When GitHub is enabled, the local docs directory is completely ignored
- File watcher is disabled when GitHub mode is enabled (no local files to watch)
- The `__lang__` folder structure in GitHub is automatically handled for multi-lingual setups
- Large files (>= 1MB) are automatically handled via GitHub's download URL API
- If GitHub fetch fails, the server will log errors but continue serving cached content if available

## Troubleshooting

### Rate Limit Errors

If you see rate limit errors, make sure to:
1. Configure a GitHub token in your `dorcs.yaml`
2. Use environment variables for the token (don't commit tokens to git)
3. Check that the token has the correct permissions (`public_repo` for public repos, `repo` for private repos)

### Files Not Appearing

- Verify the repository URL is correct and accessible
- Check that the branch/tag exists
- Ensure the path in the URL points to a directory containing markdown files
- Check server logs for any error messages

### Authentication Errors

- Verify your token is valid and not expired
- Check that the token has the correct scopes for your repository type
- For private repositories, ensure the token has `repo` scope

## Related Documentation

- [Configuration](../03_configuration.md) - Complete configuration reference
- [Multi-lingual Support](../03_configuration.md#multi-lingual-support) - Language configuration
