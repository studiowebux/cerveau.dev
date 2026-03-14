---
title: CLI Reference
---

# CLI Reference

The `cerveau` binary is installed at `~/.cerveau/bin/cerveau`. The installer adds it to your `PATH`.

```bash
cerveau help
```

## Commands

### spawn

Create a new brain and wire `.claude` into it.

```bash
cerveau spawn MyApp /absolute/path/to/your/code
```

- First argument — brain name (creates `~/.cerveau/_brains_/myapp-brain/`)
- Second argument — absolute path to your code repo

What it does:

- Creates `_brains_/myapp-brain/` with all templates and symlinks
- Replaces `__PROJECT__` placeholders with the brain name
- Generates `settings.json` with `additionalDirectories` pointing to the absolute project path
- Adds an entry to `_configs_/brains.json`
- Auto-wires MCP globally from `~/.cerveau/.env` if present (no extra step needed)


### rebuild

Rebuild selective symlinks from `brains.json` for a brain.

```bash
cerveau rebuild MyApp
```

Reads the brain's entry in `_configs_/brains.json` and creates selective symlinks for rules and agents. Run this after editing `brains.json`.

### update

Download and apply the latest Cerveau packages. Preserves `.env`, `_brains_/`, and `brains.json`.

```bash
cerveau update
```

Or from a brain session: `/update`

Safe to run at any time — your config and brains are never overwritten.

### marketplace list

List all available packages in the Cerveau marketplace.

```bash
cerveau marketplace list
```

Prints each package with its name, type, description, and tags.

### marketplace install

Install a marketplace package into a brain.

```bash
cerveau marketplace install workflow-minimaldoc MyApp
```

- First argument — package name from `marketplace list`
- Second argument — brain name to install into

Adds the package to the brain's `brains.json` `packages` array, then rebuilds rules automatically.

Or from a brain session: `/marketplace install workflow-minimaldoc MyApp`

### install-statusline

Install the status line script to `~/.claude/statusline.sh`.

```bash
cerveau install-statusline
```

Copies the statusline script from the core package and makes it executable. Run once after installing Cerveau.

### status

Show install status for a brain.

```bash
cerveau status MyApp
```

Reports: symlink status for rules/hooks/agents, settings.json validity,
presence of CLAUDE.md files.

### list

List all existing brains.

```bash
cerveau list
```

Scans `_brains_/` and prints all `*-brain` directories with their paths.

### validate

Check a brain has no remaining `__PROJECT__` placeholders.

```bash
cerveau validate MyApp
# expected: OK: No __PROJECT__ placeholders found
```

Run this after `cerveau spawn` to confirm the template was fully substituted.

### diff

Show what changed between the package templates and a brain.

```bash
cerveau diff MyApp
```

Useful for reviewing customizations made to a brain after spawning.

### help

Print all available commands with descriptions.

```bash
cerveau help
```

## Workflow

```bash
# Install once
curl -fsSL https://cerveau.dev/install.sh | bash

# Create a brain
cerveau spawn MyApp /path/to/myapp

# Launch the brain
cd ~/.cerveau/_brains_/myapp-brain && claude

# Add a marketplace package
cerveau marketplace install workflow-minimaldoc MyApp

# Update to the latest packages
cerveau update
```
