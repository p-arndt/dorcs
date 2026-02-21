---
title: "Edit"
description: "Edit your documentation site in the browser"
tags: [edit, browser, markdown]
date: 2026-02-21
after: "07_commands"
---

# Edit


Edit mode allows you to edit your documentation site in the browser. It is a simple and easy to use editor that lets you create, edit, and delete files directly in the browser—all changes are saved to your local filesystem.

## Authentication

Edit mode requires authentication when enabled. This protects your documentation from unauthorized changes.

### Enabling Authentication

Configure authentication in your `dorcs.yaml`:

```yaml
auth:
  enabled: true
  username: "admin"
  password: "your-secure-password"
  sessions_path: ".dorcs_sessions.json"  # Optional, default: .dorcs_sessions.json
```

See [Configuration](./03_configuration.md#authentication--edit-mode) for full configuration options.

### Login Flow

1. **Login** – When authentication is enabled but you're not logged in, a **Login** button appears in the footer. Click it to open the login dialog and enter your credentials.
2. **Edit** – After logging in, **Edit** and **Logout** buttons appear in the header. Click **Edit** to open the edit mode panel.
3. **Logout** – Click **Logout** to end your session.

### Session Behavior

- Sessions persist for **24 hours**
- Session data is stored in the file specified by `sessions_path` (default: `.dorcs_sessions.json` in your docs directory)
- Passwords are automatically hashed using Argon2id on first use; the plain text password is removed from the config after hashing

### Security Notes

> [!WARNING]
> Only enable edit mode on trusted networks or behind additional protection. Edit mode grants full file system access to your docs directory—users can create, edit, and delete any file.

- All edit operations (list files, read, save, create, delete) require authentication
- Use a strong password for production deployments
