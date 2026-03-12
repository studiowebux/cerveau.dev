# MDPlanner Task Workflow — Unique Patterns

Supplements the brain CLAUDE.md. Contains only patterns NOT already covered there.

## get_context_pack Fields

| Field | Extract |
| ----- | ------- |
| `people.agents` | Find "Claude" entry → Claude's person ID |
| `people.owner` | Owner's person ID |
| `milestone` | Active milestone name + ID |
| `inProgress` | Tasks to resume first |
| `todo` | Top 10 ready tasks (sorted by priority) |
| `recentProgress` | Most recent `[progress]` note excerpt |
| `decisions` / `architecture` / `constraints` | Note titles — `get_note` only when relevant |

## Version Check (every session)

1. Read version from the version file (path from `local-dev.md`)
2. `GET <server-base-url>/api/version`
3. Versions differ or backend unreachable → stop and tell user before doing anything else

## Ready Filter

```text
list_tasks { section: "Todo", project: "<name>", ready: true }
```

Returns only tasks with all `blocked_by` resolved. No need to manually check blockers.

## Batch Operations

Replaces N `update_task` + N `add_task_comment` calls:

```text
batch_update_tasks {
  updates: [
    { id: "...", section: "In Progress", assignee: "<claude-id>", milestone: "vX.Y.Z" },
    { id: "...", section: "Pending Review", assignee: "<owner-id>",  comment: "[vX.Y.Z] Fixed in abc1234 — summary" }
  ]
}
```

The `comment` field in batch updates replaces a separate `add_task_comment` call.

## Multi-Agent Claiming

```text
claim_task {
  id: "<task_id>",
  assignee: "<agent-id>",
  expected_section: "Todo",
  expected_revision: <current_revision>
}
```

| Error | Meaning | Action |
| ----- | ------- | ------ |
| `CLAIM_CONFLICT` | Another agent claimed first | Skip, pick next — do not retry |
| `CLAIM_GUARD` | Task held by different agent | Do not force-update |
| `REVISION_CONFLICT` | Stale revision | Re-fetch with `get_task`, retry |

## Sweep Stale Claims

Run at session start to recover tasks from crashed agents:

```text
sweep_stale_claims { ttl_minutes: 30 }
```

## Common Error Recovery

| Error | Action |
| ----- | ------ |
| Entity not found | `list_<entity>` to find correct ID, then retry |
| Duplicate milestone | `list_milestones { project }` → find existing → use its ID |
| `REVISION_CONFLICT` | Re-fetch with `get_task`, retry with current revision |
| `CLAIM_CONFLICT` | Skip, pick next task — do not retry |
| `CLAIM_GUARD` | Task held by different agent — do not force-update |

## Periodic Checkpoints

Every 15–20 tool calls: write a `[progress]` note with current status and next step.
Write `[decision]` or `[bug]` notes for any discoveries immediately.

Enforced by the `checkpoint-counter.sh` hook.
