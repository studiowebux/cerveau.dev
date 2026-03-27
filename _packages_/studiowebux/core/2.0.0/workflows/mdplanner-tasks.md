# MDPlanner Task Workflow — Unique Patterns

Supplements the brain CLAUDE.md. Contains only patterns NOT already covered there.
Claude discovers tool parameters from MCP schemas — this doc covers **when** and **why**, not parameter lists.

## get_context_pack

Single boot call. Extract these fields:

| Field | Use |
| ----- | --- |
| `people.agents` | Find "Claude" entry → Claude's person ID |
| `people.owner` | Owner's person ID |
| `milestone` | Active milestone name + ID |
| `inProgress` | Tasks to resume first |
| `todo` | Top 10 ready tasks (sorted by priority) |
| `recentProgress` | Most recent `[progress]` note excerpt |
| `decisions` / `architecture` / `constraints` | Note titles — `get_note` only when relevant |

If versions differ, notify the user. If backend unreachable, stop and tell user.

## Key Patterns

**Ready filter** — Use `ready: true` on `list_tasks` to get only tasks with all `blocked_by` resolved. No need to manually check blockers.

**Batch operations** — Use `batch_update_tasks` instead of N individual `update_task` + `add_task_comment` calls. The inline `comment` field replaces a separate `add_task_comment` call.

**Claiming** — Always use `claim_task` instead of `update_task` to pick up tasks. Only `id` and `assignee` are required. Run `sweep_stale_claims` at session start to recover tasks from crashed agents.

## Error Recovery

| Error | Action |
| ----- | ------ |
| Entity not found | `list_<entity>` to find correct ID, then retry |
| Duplicate milestone | `list_milestones` → use existing ID |
| `REVISION_CONFLICT` | Re-fetch with `get_task`, retry with current revision |
| `CLAIM_CONFLICT` | Another agent claimed first — skip, pick next |
| `CLAIM_GUARD` | Task held by different agent — do not force-update |

## Checkpoints

Every 15-20 tool calls: write a `[progress]` note. Write `[decision]` or `[bug]` notes immediately on discovery. Enforced by `checkpoint-counter.sh` hook.
