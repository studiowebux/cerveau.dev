---
title: Makefile Targets
---

# Makefile Targets

All targets run from `~/brains/_protocol_/`:

```bash
cd ~/brains/_protocol_
make help
```

## Targets

### spawn

Create a new brain and wire `.claude` into it.

```bash
make spawn NAME=MyApp PROJECT=/absolute/path/to/your/code
```

- `NAME` — brain name (creates `_brains_/myapp-brain/`)
- `PROJECT` — absolute path to your code repo

What it does:
- Creates `_brains_/myapp-brain/` with all templates
- Replaces `__PROJECT__` placeholders with `NAME`
- Creates selective symlinks for rules, agents, hooks
- Generates `settings.json` with `additionalDirectories` pointing to `PROJECT`
- Adds an entry to `_configs_/brains.json`

### status

Show install status for a brain.

```bash
make status NAME=MyApp
```

Reports: symlink status for rules/hooks/agents, settings.json validity,
presence of CLAUDE.md files.

### list

List all existing brains.

```bash
make list
```

Scans `_brains_/` and prints all `*-brain` directories with their paths.

### validate

Check a brain has no remaining `__PROJECT__` placeholders.

```bash
make validate NAME=MyApp
# expected: OK: No __PROJECT__ placeholders found
```

Run this after `make spawn` to confirm the template was fully substituted.

### diff

Show what changed between the protocol template and a brain.

```bash
make diff NAME=MyApp
```

Useful for reviewing customizations made to a brain after spawning.

### sync-shareable

Copy changed protocol files to `_shareable_/`. Only updates files that differ.

```bash
make sync-shareable
```

Run this from the monorepo `_protocol_/` (not from `_shareable_/_protocol_/`).
Then `cd _shareable_` and commit + push to GitHub.

## Workflow

```bash
# 1. Create a brain
make spawn NAME=MyApp PROJECT=/path/to/myapp

# 2. Connect MCP
cd ../_brains_/myapp-brain && claude mcp add --transport http mdplanner \
  http://localhost:8003/mcp --header "Authorization: Bearer <token>"

# 3. Rebuild selective rules
cd ~/brains && ./_scripts_/rebuild-brain-rules.sh MyApp

# 4. Launch
cd ~/brains/_brains_/myapp-brain && claude
```
