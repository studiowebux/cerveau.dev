# Local Task Workflow

Supplements the brain CLAUDE.md. Contains patterns for managing tasks, notes, and milestones as local files.

## File Formats

### Task File (`tasks/<id>.md`)

```markdown
---
title: Short task title
section: Backlog | Todo | In Progress | Pending Review | Done
priority: 1-5
milestone: milestone-name
assignee: Claude | Owner
tags: [tag1, tag2]
created: YYYY-MM-DD
updated: YYYY-MM-DD
blocked_by: [task-id-1, task-id-2]
---
## Description

What needs to be done.

## Comments

- YYYY-MM-DD Author: Comment text...
```

**Task ID convention:** `task-YYYYMMDD-short-slug` (e.g. `task-20260320-add-login-page`)

### Note File (`notes/<type>-YYYYMMDD-slug.md`)

```markdown
---
type: progress | decision | bug | architecture | feature | constraint | investigation
project: <project>
created: YYYY-MM-DD
---
# [type] <project> — Title

Content...
```

### Milestone File (`milestones/<name>.md`)

```markdown
---
name: milestone-name
status: open | closed
created: YYYY-MM-DD
description: Short description
---
# milestone-name

## Description

What this milestone delivers.
```

## Key Patterns

**Reading tasks** — Use `Glob` to find task files, `Read` to get full content. Filter by section using `Grep` on frontmatter.

```
Glob: tasks/*.md
Grep: "section: Todo" in tasks/
Grep: "section: In Progress" in tasks/
```

**Creating a task** — Write a new file to `tasks/<id>.md` with frontmatter and description.

**Updating a task** — Edit the frontmatter fields (section, assignee, etc.) and add a comment entry.

**Finding ready tasks** — Grep for `section: Todo`, then check each for `blocked_by`. A task is ready if `blocked_by` is empty or all referenced tasks have `section: Done`.

**Moving a task** — Edit the `section:` field in frontmatter. Update `context.md` immediately after.

**Boot sequence** — Read `context.md` instead of calling `get_context_pack`. If `context.md` is stale or missing, rebuild it by scanning `tasks/`, `notes/`, and `milestones/`.

## Updating context.md

After any of these actions, update `context.md`:
- Task status change (section change)
- New task created
- New note written
- Milestone status change

Update only the relevant section — do not rewrite the entire file for a single task move.

## Checkpoints

Every 15-20 tool calls: write a `[progress]` note to `notes/`. Write `[decision]` or `[bug]` notes immediately on discovery.
