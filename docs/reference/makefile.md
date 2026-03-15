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


### boot

Launch Claude Code inside a brain directory. Works from anywhere.

```bash
cerveau boot MyApp
cerveau boot MyApp --resume
```

- First argument — brain name
- Remaining arguments are passed through to `claude`

Replaces the manual `cd ~/.cerveau/_brains_/myapp-brain && claude` workflow.

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

List available packages in the Cerveau marketplace, with optional filtering.

```bash
cerveau marketplace list                        # all packages
cerveau marketplace list theme                  # free-text search
cerveau marketplace list --tag design           # filter by tag
cerveau marketplace list --org studiowebux      # filter by org
cerveau marketplace list --org _local_          # show only local packages
```

Prints each package with its name, type, description, and tags. Filters are case-insensitive.

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

### backup

Create a backup archive of your Cerveau environment.

```bash
cerveau backup                        # backup everything (default)
cerveau backup --cerveau              # ~/.cerveau/ only
cerveau backup --mdplanner            # MDPlanner data only (~/.cerveau/data/)
cerveau backup --claude               # ~/.claude/ only
cerveau backup --cerveau --claude     # combine flags
cerveau backup --all -o /tmp/bk.tar.gz  # custom output path
```

The archive includes a `manifest.json` with metadata (timestamp, version, sections). The `cerveau` binary is excluded — reinstall via `cerveau update`.

For a consistent MDPlanner backup, stop the container first.

### restore

Restore from a backup archive.

```bash
cerveau restore backup-2026-03-15.tar.gz              # restore everything in archive
cerveau restore backup-2026-03-15.tar.gz --claude      # restore only ~/.claude/
cerveau restore backup-2026-03-15.tar.gz --mdplanner   # restore only MDPlanner data
```

Restore shows what will be overwritten and asks for confirmation before proceeding.

### dir

Print the absolute path to a brain or its codebase directory. Output is a single line with no decoration — designed for scripting and piping.

```bash
cerveau dir brain MyApp    # prints ~/.cerveau/_brains_/myapp-brain
cerveau dir code MyApp     # prints the codebase path from brains.json
```

### cd

Navigate to a brain or codebase directory. Requires the shell wrapper from `cerveau completion`.

```bash
cerveau cd brain MyApp     # cd to the brain directory
cerveau cd code MyApp      # cd to the codebase directory
```

Since a subprocess cannot change the parent shell's working directory, `cerveau cd` is implemented as a shell function that calls `cerveau dir` internally. The function is included in the completion script output — see `cerveau completion` below.

### completion

Output a shell completion script for tab-tab support.

```bash
eval "$(cerveau completion zsh)"     # add to .zshrc
eval "$(cerveau completion bash)"    # add to .bashrc
```

This enables:

- `cerveau <tab>` — list all commands
- `cerveau boot <tab>` — list brain names
- `cerveau cd <tab>` — complete `brain` or `code`, then brain names
- `cerveau marketplace <tab>` — list subcommands
- `cerveau marketplace list --tag <tab>` — list available tags

The completion script also installs a shell wrapper function that makes `cerveau cd` work.

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
cerveau boot MyApp

# Navigate to the codebase
cerveau cd code MyApp

# Add a marketplace package
cerveau marketplace install workflow-minimaldoc MyApp

# Update to the latest packages
cerveau update
```
