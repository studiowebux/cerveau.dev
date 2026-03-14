---
title: Directory Structure
---

# Directory Structure

Full annotated layout of a Cerveau installation (`~/.cerveau/`).

```
~/.cerveau/
├── _packages_/                          # Shared packages (source of truth)
│   ├── studiowebux/
│   │   ├── core/
│   │   │   └── 1.0.0/
│   │   │       ├── CLAUDE.md            # Brain protocol (phases split into rule files)
│   │   │       ├── statusline.sh        # Status bar script (install via cerveau install-statusline)
│   │   │       ├── templates/           # Note templates
│   │   │       │   ├── architecture.md
│   │   │       │   ├── decision-record.md
│   │   │       │   ├── feature-spec.md
│   │   │       │   └── project-overview.md
│   │   │       ├── hooks/               # Hook scripts (wholesale symlink in brains)
│   │   │       │   ├── checkpoint-counter.sh
│   │   │       │   ├── commit-validator.sh
│   │   │       │   ├── context-warning.sh
│   │   │       │   ├── post-edit-reminder.sh
│   │   │       │   ├── pre-compact-handoff.sh
│   │   │       │   └── session-context.sh
│   │   │       ├── skills/              # Skill definitions (wholesale symlink in brains)
│   │   │       │   ├── release/SKILL.md
│   │   │       │   ├── update/SKILL.md
│   │   │       │   └── marketplace/SKILL.md
│   │   │       ├── rules/               # Rules library (selective symlinks in brains)
│   │   │       │   ├── phase-boot.md
│   │   │       │   ├── phase-work.md
│   │   │       │   ├── phase-close.md
│   │   │       │   ├── code-discipline.md
│   │   │       │   ├── goal-discipline.md
│   │   │       │   ├── stack/           #   Language/framework rules — add your own
│   │   │       │   ├── practices/       #   Development practice rules — add your own
│   │   │       │   └── workflow/        #   Process rules
│   │   │       │       ├── local-dev.md
│   │   │       │       └── mdplanner-tasks.md
│   │   │       └── settings.json.template
│   │   └── minimaldoc/
│   │       └── 1.0.0/
│   │           └── agents/              # Agent definitions (selective symlink in brains)
│   │               └── minimaldoc-writer.md
│
├── _configs_/
│   ├── brains.json                      # Brain registry (declares packages each brain loads)
│   └── registry.json                    # Package catalog
│
├── _brains_/                            # One directory per brain (created by cerveau spawn)
│   └── <brain-name>/
│       ├── templates/                   # Copied from package on spawn
│       └── .claude/
│           ├── CLAUDE.md                # Symlink → package CLAUDE.md
│           ├── settings.json            # Generated — hooks + additionalDirectories (absolute path)
│           ├── settings.local.json      # Local overrides (not committed)
│           ├── hooks  -> _packages_/.../hooks     # Wholesale symlink
│           ├── skills -> _packages_/.../skills    # Wholesale symlink
│           ├── agents/                  # Selective symlinks (only from declared packages)
│           └── rules/                   # Selective symlinks (only from declared packages)
│               ├── phase-boot.md  -> ...
│               ├── phase-work.md  -> ...
│               ├── phase-close.md -> ...
│               ├── stack/
│               ├── practices/
│               └── workflow/
│                   └── local-dev.md     #   Real file (not symlink) — brain-specific config
│
├── _scripts_/
│   └── backup-claude.sh                 # Archives Claude session logs
│
├── bin/
│   └── cerveau                          # CLI binary (spawn, rebuild, update, marketplace, etc.)
│
├── .env                                 # MCP token + config (preserved across updates)
├── version.txt                          # Installed Cerveau version
├── cerveau-package.json                 # Version manifest
├── install.sh                           # Installer script
└── docker-compose.yml                   # MDPlanner container
```

## Key Distinctions

**`_brains_/<brain>/.claude/rules/workflow/local-dev.md` is the only real file**
in a brain's rules directory. Everything else is a symlink. It holds
brain-specific config: MCP project name, server URL, person IDs, active
milestone, and Brain Memory.

**`stack/` and `practices/` ship empty.** You generate rules for your own
stack and practices — see [Writing Rules](../02-guides/02-writing-rules.md).

**Your code lives anywhere.** The brain's `settings.json` uses an absolute path
in `additionalDirectories` to point at your project repo — git submodule,
separate clone, or any directory on disk.
