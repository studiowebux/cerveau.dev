---
title: Directory Structure
---

# Directory Structure

Full annotated layout of a cerveau.dev installation.

```
cerveau.dev/
├── _protocol_/                          # Shared protocol (single source of truth)
│   ├── CLAUDE.md                        # Brain protocol (phases split into rule files)
│   ├── statusline.sh                    # Status bar script (install via make install)
│   ├── Makefile                         # onboard, spawn, install, status, validate, list, diff, help
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
│       │   ├── session-context.sh
│       │   └── stop-progress-check.sh
│       ├── agents/                      # Agent definitions (selective symlink in brains)
│       │   └── minimaldoc-writer.md     # Example agent (use as template)
│       ├── skills/                      # Skill definitions (wholesale symlink in brains)
│       │   ├── import-project/SKILL.md
│       │   └── release/SKILL.md
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
│   └── brains.json                      # Brain registry (declares what each brain loads)
│
├── _brains_/                            # One directory per brain (created by make spawn/onboard)
│   └── <brain-name>/
│       ├── templates/                   # Copied from protocol on spawn
│       ├── setup/                       # Copied from protocol on spawn
│       └── .claude/
│           ├── CLAUDE.md                # Symlink → _protocol_/CLAUDE.md
│           ├── settings.json            # Generated — hooks + additionalDirectories
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
│   ├── rebuild-brain-rules.sh           # Rebuilds selective symlinks from brains.json
│   └── backup-claude.sh                 # Archives Claude session logs
│
├── docker-compose.yml                   # MDPlanner container
└── .env.example                         # Environment variable template
```

## Key Distinctions

**`_brains_/<brain>/.claude/rules/workflow/local-dev.md` is the only real file**
in a brain's rules directory. Everything else is a symlink. It holds
brain-specific config: MCP project name, server URL, person IDs, active
milestone, and Brain Memory.

**`stack/` and `practices/` ship empty.** You generate rules for your own
stack and practices — see [Writing Rules](../02-guides/02-writing-rules.md).

**`_projects_/` is optional.** Your code can live anywhere — git submodule,
absolute path, separate clone. What matters is that `brains.json` has the
correct `codebase` path pointing to it.
