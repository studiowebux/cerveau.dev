---
title: Skills
---

# Skills

Skills are slash commands — reusable, multi-step workflows Claude executes on
demand. They are plain markdown files that Claude reads and follows as
instructions.

## File Format

Skills live in `_protocol_/.claude/skills/<name>/SKILL.md`:

```
_protocol_/.claude/skills/
  release/
    SKILL.md
  import-project/
    SKILL.md
  update/
    SKILL.md
  marketplace/
    SKILL.md
```

The skill name is the directory name. Invoke it with `/release`, `/update`, etc.

A `SKILL.md` is a markdown document describing what Claude should do step by
step. No special syntax — just clear instructions.

## Included Skills

The protocol ships with four skills:

- `/release` — full release workflow: version bump, changelog, build
  verification, tag, push, GitHub release, MDPlanner progress note.
- `/import-project` — onboard a codebase into MDPlanner and spawn its brain
  in one automated flow. Accepts `NAME=MyApp PROJECT=/path/to/code` as
  arguments. Run this from the protocol session to bootstrap a new project.
- `/update` — download and install the latest Cerveau protocol. Preserves
  `.env`, `_brains_/`, and `brains.json`. Reports the version before and after.
- `/marketplace` — browse available packages and install them into a brain.
  Use `/marketplace` to list, or `/marketplace install <pkg> [brain]` to
  install a package.

Read the `SKILL.md` files in `_protocol_/.claude/skills/` for the full
procedures and use them as templates when writing your own.

## Writing Skills

Skills work best for workflows you run repeatedly that have a fixed sequence
of steps. Examples:

```
Create a skill for deploying to staging.
Steps: build, run smoke tests, push image to registry, update the
deployment, verify health endpoint returns 200.
Save it to _protocol_/.claude/skills/deploy-staging/SKILL.md
```

```
Create a skill for cutting a hotfix.
Steps: checkout main, create hotfix branch, confirm the fix is committed,
bump patch version, changelog entry, tag, push, PR.
Save it to _protocol_/.claude/skills/hotfix/SKILL.md
```

Keep skills focused. One workflow per skill. If a skill branches based on
conditions, split it into two skills.

## Skills vs. Agents

| | Skills | Agents |
|---|---|---|
| Invocation | `/skill-name` by user | Called by Claude when needed |
| Purpose | Multi-step workflow | Specialized sub-instance |
| Format | Plain markdown steps | YAML frontmatter + system prompt |
| Tool access | Inherits all tools | Restricted to declared tools |
