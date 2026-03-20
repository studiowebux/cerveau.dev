# Phase 3 — Writing Back

Every note file name MUST include `<project>`. Use templates from `templates/`.

| When | File name format |
|------|-----------------|
| Decision made | `notes/decision-YYYYMMDD-<slug>.md` |
| Bug root cause found | `notes/bug-YYYYMMDD-<slug>.md` |
| Session end or milestone | `notes/progress-YYYYMMDD-<slug>.md` |
| Architecture established | `notes/architecture-YYYYMMDD-<slug>.md` |
| Feature specced | `notes/feature-YYYYMMDD-<slug>.md` |
| Hard limit confirmed | `notes/constraint-YYYYMMDD-<slug>.md` |
| Investigation paused | `notes/investigation-YYYYMMDD-<slug>.md` |

Note title (H1 inside the file) follows: `[type] <project> — <title>`

---

# Phase 4 — Session Close

1. Write a `[progress]` note file summarizing what was accomplished.
2. Update `context.md` — refresh In Progress, Todo, Recent Progress, and Key Notes sections.
3. Leave unfinished tasks In Progress — Phase 1 Boot resumes them automatically next session.
4. Create discovered tasks in `tasks/` with `section: Backlog`.
