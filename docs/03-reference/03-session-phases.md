---
title: Session Phases
---

# Session Phases

Every brain session follows four phases defined in the brain CLAUDE.md.

## Phases

| Phase | What happens |
|---|---|
| **Phase 1 — Boot** | Load context from MDPlanner: most recent progress note, open tasks, architecture, decisions |
| **Phase 2 — Work** | Ticket before work, one task at a time, implement, commit, update task, move to Done |
| **Phase 3 — Write Back** | Record decisions, bugs, learnings as MDPlanner notes |
| **Phase 4 — Close** | Write progress note, move unfinished tasks back to Todo |

## Phase 1 — Boot

The `session-context` hook fires on `SessionStart` and reminds Claude to run
Phase 1 before any work.

Claude reads from MDPlanner:
1. Most recent progress note (sorted by `updatedAt`)
2. Open tasks for the active milestone
3. Architecture notes
4. Recent decisions

Then checks git status in the project repo.

## Phase 2 — Work

### Commit Flow

After implementing a task:

1. Stage changes: `git add <files>`
2. Commit: `git commit -m "type: subject"`
3. Push: `git push`
4. Add progress comment to the MDPlanner task
5. Move task to Done in MDPlanner

The `commit-validator` hook blocks commits that:
- Don't follow conventional commit format
- Have staged files containing secret patterns

### Conventional Commits

```
feat: add login endpoint
fix: resolve nil pointer in auth handler
chore: update dependencies
docs: add API reference
test: add unit tests for validator
refactor: extract shared config loader
```

## Phase 3 — Write Back

After completing a task or making a significant decision:

- **Decision notes**: `[decision] <title> — <rationale>` — never edit, supersede
- **Bug notes**: `[bug] <title> — <description>` — discovered during work
- **Architecture updates**: update `[architecture]` notes when structure changes

## Phase 4 — Close

Before ending the session:

1. Write a progress note: `[progress] <date> — <summary of what was done>`
2. Move any unfinished In Progress tasks back to Todo
3. The `stop-progress-check` hook verifies a progress note was written

## Hard Rules

| Rule | Description |
|---|---|
| Ticket before work | No code changes without an MDPlanner task |
| One task at a time | Only one task in In Progress |
| Branch before commit | Create a feature branch before starting work |
| Push after commit | Always push immediately after committing |
| Never mark complete | Only the human sets `completed: true` |
| Never edit decisions | Create a superseding note instead |
| Progress before close | Must write a progress note before session ends |
| No unilateral deferrals | Deferring tasks requires owner approval |
| Continue means work | `continue` resumes the current task, not exploration |
| Backlog ownership | Claude owns the backlog — keeps it clean and up to date |
| Brain stays external | Brain directory never lives inside a project repo |
| One-task-at-a-time | Never work on multiple tasks simultaneously |
| No force push to main | Protected branch — PR workflow only |
| Context compaction note | Write handoff note before context is compacted |
| Secrets gate | Commit validator scans for secrets — blocked if found |
| Checkpoint acknowledgement | Respond to checkpoint hooks before continuing |
