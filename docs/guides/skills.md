---
title: Skills
---

# Skills

Skills are slash commands — reusable, multi-step workflows Claude executes on
demand. They are plain markdown files that Claude reads and follows as instructions.

> See the official [Claude Code skills documentation](https://docs.anthropic.com/en/docs/claude-code/skills) for the full file format reference.

:::warning
**Never save skills into `_packages_/studiowebux/core/` or any downloaded package.** Those files are overwritten by `cerveau update`. Always save your skills to a `_local_` package.
:::

## File Format

Each skill lives in its own subdirectory with a `SKILL.md` file. The directory name is the slash command:

```
_packages_/_local_/my-skills/1.0.0/
  skills/
    deploy-staging/
      SKILL.md     ← invoked as /deploy-staging
    hotfix/
      SKILL.md     ← invoked as /hotfix
```

A `SKILL.md` is plain markdown describing what Claude should do step by step. No special syntax — just clear instructions.

## Included Skills

The core package ships with three skills:

- `/release` — full release workflow: version bump, changelog, build verification, tag, push, GitHub release, MDPlanner progress note
- `/update` — download and install the latest Cerveau protocol. Preserves `.env`, `_brains_/`, and `brains.json`. Reports the version before and after
- `/marketplace` — browse available packages and install them into a brain

Read the `SKILL.md` files in `~/.cerveau/_packages_/studiowebux/core/1.0.0/skills/` for the full procedures and use them as templates when writing your own.

## Writing Skills

Skills work best for workflows you run repeatedly that have a fixed sequence of steps:

<details>
<summary>Deploy to staging</summary>

```
Create a skill for deploying to staging.
Steps: build, run smoke tests, push image to registry, update the
deployment, verify health endpoint returns 200.
Save to ~/.cerveau/_packages_/_local_/my-skills/1.0.0/skills/deploy-staging/SKILL.md
```

</details>

<details>
<summary>Cut a hotfix</summary>

```
Create a skill for cutting a hotfix.
Steps: checkout main, create hotfix branch, confirm the fix is committed,
bump patch version, changelog entry, tag, push, PR.
Save to ~/.cerveau/_packages_/_local_/my-skills/1.0.0/skills/hotfix/SKILL.md
```

</details>

Keep skills focused. One workflow per skill. If a skill branches based on conditions, split it into two skills.

## Skills vs. Agents

| | Skills | Agents |
|---|---|---|
| Invocation | `/skill-name` by user | Called by Claude when needed |
| Purpose | Multi-step workflow | Specialized sub-instance |
| Format | Plain markdown steps | YAML frontmatter + system prompt |
| Tool access | Inherits all tools | Restricted to declared tools |
