---
title: Directory Structure
---

# Directory Structure

Full annotated layout of a Cerveau installation (`~/.cerveau/`).

```
~/.cerveau/
├── _protocol_/                          # Shared protocol (single source of truth)
│   ├── CLAUDE.md                        # Brain protocol (phases split into rule files)
│   ├── statusline.sh                    # Status bar script (install via cerveau install-statusline)
│   ├── templates/                       # Note templates
│   │   ├── architecture.md
│   │   ├── decision-record.md
│   │   ├── feature-spec.md
│   │   └── project-overview.md
│   └── .claude/
│       ├── settings.json.template       # Hooks config template (generated into each brain)
│       ├── hooks/                       # Hook scripts (wholesale symlink in brains)
│       │   ├── checkpoint-counter.sh
│       │   ├── commit-validator.sh
│       │   ├── context-warning.sh
│       │   ├── post-edit-reminder.sh
│       │   ├── pre-compact-handoff.sh
│       │   └── session-context.sh       # Boot reminder + version check
│       ├── agents/                      # Agent definitions (selective symlink in brains)
│       │   └── minimaldoc-writer.md     # Example agent (use as template)
│       ├── skills/                      # Skill definitions (wholesale symlink in brains)
│       │   ├── import-project/SKILL.md
│       │   ├── release/SKILL.md
│       │   ├── update/SKILL.md          # /update — pull latest protocol
│       │   └── marketplace/SKILL.md     # /marketplace — browse and install packages
│       └── rules/                       # Rules library (selective symlinks in brains)
│           ├── phase-boot.md            #   Phase 1 — Boot sequence (core, always loaded)
│           ├── phase-work.md            #   Phase 2 — Work and commit flow (core, always loaded)
│           ├── phase-close.md           #   Phase 3+4 — Write back and close (core, always loaded)
│           ├── code-discipline.md       #   Core (always loaded)
│           ├── goal-discipline.md       #   Core (always loaded)
│           ├── stack/                   #   Language/framework rules — add your own
│           ├── practices/               #   Development practice rules — add your own
│           └── workflow/                #   Process rules
│               ├── local-dev.md         #     Brain config table template
│               └── mdplanner-tasks.md   #     MDPlanner task workflow
│
├── _configs_/
│   ├── brains.json                      # Brain registry (declares what each brain loads)
│   └── registry.json                    # Marketplace package catalog
│
├── _brains_/                            # One directory per brain (created by cerveau spawn/onboard)
│   └── <brain-name>/
│       ├── templates/                   # Symlink → _protocol_/templates
│       └── .claude/
│           ├── CLAUDE.md                # Symlink → _protocol_/CLAUDE.md
│           ├── settings.json            # Generated — hooks + additionalDirectories (absolute path)
│           ├── settings.local.json      # Local overrides (not committed)
│           ├── hooks  -> _protocol_/.claude/hooks    # Wholesale symlink
│           ├── skills -> _protocol_/.claude/skills   # Wholesale symlink
│           ├── agents/                  # Selective symlinks (only declared agents)
│           └── rules/                   # Selective symlinks (only declared rules)
│               ├── phase-boot.md  -> ...
│               ├── phase-work.md  -> ...
│               ├── phase-close.md -> ...
│               ├── stack/               #   Only stacks declared in brains.json
│               ├── practices/           #   Only practices declared in brains.json
│               └── workflow/            #   Only workflows declared in brains.json
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
