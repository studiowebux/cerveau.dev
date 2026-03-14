---
title: Packages Overview
---

# Packages

Everything in Cerveau is a package. Rules, workflows, hooks, skills, agents,
templates — all live in `_packages_/` and are installed into brains via the
registry.

## Structure

```
_packages_/
  studiowebux/
    core/
      1.0.0/
        rules/
        workflows/
        hooks/
        skills/
        agents/
        templates/
        claude/
    minimaldoc/
      1.0.0/
        workflows/
        practices/
        agents/
  _local_/
    my-pkg/
      0.1.0/
        workflows/
```

Path pattern: `_packages_/{org}/{package-name}/{version}/{type}/[files]`

## Org Namespacing

Every package belongs to an org. The org prevents name collisions — two orgs
can each publish their own `core` package without conflict.

- **`studiowebux`** — official packages, updated from upstream via `cerveau update`
- **`_local_`** — your private packages, never touched by updates

## Package Types

The `type` field on each file determines where it gets symlinked in the brain:

| Type | Package directory | Brain destination |
|---|---|---|
| `rules` | `rules/` | `.claude/rules/` |
| `workflows` | `workflows/` | `.claude/rules/workflow/` |
| `practices` | `practices/` | `.claude/rules/practices/` |
| `stacks` | `stacks/` | `.claude/rules/stack/` |
| `hooks` | `hooks/` | `.claude/hooks/` |
| `skills` | `skills/` | `.claude/skills/` |
| `agents` | `agents/` | `.claude/agents/` |
| `templates` | `templates/` | `templates/` |
| `claude` | `claude/` | `.claude/` |

## Installing Packages

### On spawn

```bash
cerveau spawn MyApp /path/to/code
cerveau spawn MyApp /path/to/code --packages studiowebux/core,studiowebux/minimaldoc
```

Without `--packages`, the default is `studiowebux/core`.

### After spawn

```bash
cerveau marketplace install studiowebux/minimaldoc MyApp
```

### From Claude Code

```
/marketplace install studiowebux/minimaldoc MyApp
```

## Removing Packages

```bash
cerveau marketplace uninstall studiowebux/minimaldoc MyApp
```

This removes the package from `brains.json` and cleans up symlinks on rebuild.

## Browsing

```bash
cerveau marketplace list
cerveau marketplace info studiowebux/core
```

`list` shows all packages from both `registry.json` and `registry.local.json`.
`info` shows the full file list with types.

## Customizing Package Files

Every file is symlinked by default. To customize a file from an installed
package:

1. Delete the symlink in the brain (e.g. `rm brain/.claude/hooks/commit-validator.sh`)
2. Place your own version as a real file at the same path
3. Run `cerveau rebuild` — it preserves real files and skips the symlink

## Real Files

Some files need to be copied instead of symlinked because they contain
per-brain data. These are marked `realFile: true` in the registry. Example:
`local-dev.md` gets copied and has placeholders substituted with brain-specific
values.

Rebuild never overwrites an existing real file.

## Private Packages

Create packages under `_packages_/_local_/`:

```
_packages_/_local_/my-hooks/0.1.0/hooks/my-hook.sh
```

Register them in `_configs_/registry.local.json`:

```json
{
  "version": "1.0.0",
  "packages": [
    {
      "name": "my-hooks",
      "org": "_local_",
      "version": "0.1.0",
      "path": "_packages_/_local_/my-hooks/0.1.0",
      "description": "My custom hooks",
      "files": [
        {"name": "my-hook.sh", "type": "hooks"}
      ],
      "tags": []
    }
  ]
}
```

Local packages are never overwritten by `cerveau update`. Entries in
`registry.local.json` must use the `_local_` org — others are ignored.

## Updates

`cerveau update` downloads the latest upstream packages. It:

- Overwrites upstream orgs (e.g. `studiowebux/`)
- Never touches `_packages_/_local_/`
- Preserves `.env`, `brains.json`, `registry.local.json`
- Compares registries — if upstream removed a file your brains depend on,
  it warns and asks for confirmation before applying
- Auto-rebuilds all brains after a successful update
