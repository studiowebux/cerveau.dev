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

### Core Packages

Cerveau ships two base protocol packages — pick one per brain:

- **`studiowebux/core`** — full brain protocol with MDPlanner integration. Tasks, notes, and milestones are stored in MDPlanner (requires server + MCP).
- **`studiowebux/core-local`** — same session phases, discipline rules, hooks, skills, and templates, but stores everything as local markdown files. No server, no container, fully offline.

Both provide identical code discipline, goal discipline, commit flow, and session management. The only difference is where task state lives.

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

## Versioning

Each package can have multiple versions in the registry. Brains pin exactly one
version per package — installing a different version replaces the previous one
automatically.

Package references use the `org/name[@version]` format:

- `studiowebux/core` — resolves to the first available version
- `studiowebux/core@1.0.0` — pins to an exact version

Brains store versioned refs (e.g. `studiowebux/core@1.0.0`) in `brains.json`.

## Installing Packages

### On spawn

```bash
cerveau spawn MyApp /path/to/code
cerveau spawn MyApp /path/to/code --packages studiowebux/core,studiowebux/minimaldoc
cerveau spawn MyApp /path/to/code --packages studiowebux/core@1.0.0
```

Without `--packages`, the default is `studiowebux/core`.

### After spawn

```bash
cerveau marketplace install studiowebux/minimaldoc MyApp
cerveau marketplace install studiowebux/core@2.0.0 MyApp          # upgrade to v2
cerveau marketplace install studiowebux/minimaldoc,studiowebux/github MyApp  # multiple
```

Install accepts comma-separated package refs. If another version of the same
package is already installed, it is replaced automatically.

### From Claude Code

```
/marketplace install studiowebux/minimaldoc MyApp
```

## Removing Packages

```bash
cerveau marketplace uninstall studiowebux/minimaldoc MyApp
cerveau marketplace uninstall studiowebux/minimaldoc,studiowebux/github MyApp  # multiple
```

Uninstall matches by base `org/name` regardless of installed version. It removes
the package from `brains.json` and cleans up symlinks on rebuild.

## Browsing

```bash
cerveau marketplace list                        # all packages (grouped by name, all versions shown)
cerveau marketplace list theme                  # free-text search
cerveau marketplace list --tag design           # filter by tag
cerveau marketplace list --org studiowebux      # filter by org
cerveau marketplace list --org _local_          # show only local packages
cerveau marketplace info studiowebux/core       # show package details (latest)
cerveau marketplace info studiowebux/core@1.0.0 # show specific version
```

`list` groups packages by `org/name` and shows all available versions in brackets.
`info` shows the full file list with types and lists all available versions.

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
