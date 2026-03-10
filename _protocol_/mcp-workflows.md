# MDPlanner MCP Workflows

Common patterns for using MCP tools together. Reference alongside
`mcp-reference.md` for parameter details.

**All queries below use `__PROJECT__` as the project scope. Replace with your
actual project name.**

---

## Session Start (every session, always)

```
1. list_tasks { section: "In Progress", project: "__PROJECT__" }  ← resume first
2. list_notes { search: "[project] __PROJECT__" }                  ← project context
   → get_note { id } for each result                               ← read full content
3. list_notes { search: "[architecture] __PROJECT__" }
   → get_note { id } for each result
4. list_notes { search: "[decision] __PROJECT__" }
   → get_note { id } for each result
5. list_notes { search: "[constraint] __PROJECT__" }
   → get_note { id } for each result
6. list_tasks { section: "Todo", project: "__PROJECT__" }          ← find next work
7. list_milestones { project: "__PROJECT__", status: "open" }      ← active milestone
```

If no `[project] __PROJECT__` note exists: stop, ask user to create one.
If MCP is unreachable: stop, tell user, do not proceed without context.

---

## Starting a Task

```
1. get_task { id }                     ← read full description
2. update_task {                       ← claim it
     id,
     section: "In Progress"
   }
3. list_milestones { project: "__PROJECT__", status: "open" }
   → if none: create_milestone { name: "...", project: "__PROJECT__" }
4. update_task { id, milestone: "..." } ← link to milestone (if not already)
```

---

## Completing a Task

```
1. add_task_comment {
     id,
     comment: "<what was done, commit hash if applicable>"
   }
2. update_task {
     id,
     section: "Done"
   }
```

Never set `completed: true` — that is the owner's action.

---

## Recording a Decision

After any non-obvious technical choice:

```
create_note {
  title: "[decision] __PROJECT__ — <short title>",
  content: "## Date\nYYYY-MM-DD\n\n## Status\n`closed`\n\n## Context\n...\n\n## Decision\n...\n\n## Rationale\n...\n\n## Consequences\n..."
}
```

---

## Recording Session Progress

At the end of every session or after a significant block of work:

```
create_note {
  title: "[progress] __PROJECT__ — YYYY-MM-DD <brief summary>",
  content: "## Done\n- ...\n\n## Commits\n- ...\n\n## Open\n- ...\n\n## Next\n- ..."
}
```

---

## Building a New Feature

```
1. list_notes { search: "[feature] __PROJECT__ — <name>" }
   → found → get_note { id } → read full spec
   → not found → ask user to fill template, or offer to create it

2. list_notes { search: "[architecture] __PROJECT__" }
   → get_note { id } for each → verify new code fits design

3. create_task {
     title: "Implement: <feature name>",
     section: "In Progress",
     project: "__PROJECT__",
     milestone: "...",
     tags: ["Feature"],
     priority: 2
   }

4. [build the feature]

5. add_task_comment { id, comment: "<what was done>" }
6. update_task { id, section: "Done" }
7. create_note { title: "[progress] __PROJECT__ — YYYY-MM-DD ...", content: "..." }
```

---

## Planning a Sprint / Milestone

```
1. list_milestones { project: "__PROJECT__", status: "open" }
2. list_tasks { section: "Todo", project: "__PROJECT__" }
3. list_tasks { section: "Backlog", project: "__PROJECT__" }
   → review priorities, pick batch for milestone
4. update_task × N { milestone: "<name>", priority: N }
5. create_note {
     title: "[progress] __PROJECT__ — sprint plan YYYY-MM-DD",
     content: "## Milestone: <name>\n## Tasks\n- ...\n## Goal\n..."
   }
```

---

## Retrospective After a Milestone

```
1. list_tasks { section: "Done", milestone: "<name>" }
2. list_tasks { section: "In Progress", project: "__PROJECT__" }  ← catch stragglers
3. create_retrospective {
     title: "__PROJECT__ — <milestone> retro",
     date: "<today>"
   }
4. update_retrospective {
     id,
     continue: ["<what worked>"],
     stop:     ["<what didn't>"],
     start:    ["<new practices>"],
     status:   "closed"
   }
5. update_milestone { id, status: "completed" }
```

---

## Weekly Review

```
1. list_tasks { section: "In Progress", project: "__PROJECT__" }
2. list_tasks { section: "Todo", project: "__PROJECT__", priority: 1 }
3. list_meetings { date_from: "<7 days ago>", open_actions_only: true }
4. list_notes { search: "[progress] __PROJECT__" }    ← recent logs
5. create_journal_entry {
     date: "<today>",
     title: "__PROJECT__ weekly review",
     body: "## Done\n...\n## Open\n...\n## Focus next week\n..."
   }
```

---

## Knowledge Base Maintenance

Periodically clean up and enrich the brain:

```
list_notes { search: "[architecture] __PROJECT__" }  → check still accurate
list_notes { search: "[decision] __PROJECT__" }      → verify still valid
list_notes { search: "[bug] __PROJECT__" }           → check if resolved
list_notes { search: "[feature] __PROJECT__" }       → mark implemented
get_project_config                                    → verify metadata
```

---

## Error Handling Patterns

```
# Entity not found
→ list_<entity> to find correct ID, then retry

# Duplicate milestone
→ list_milestones { project: "__PROJECT__" } → find existing → use its ID

# Search not available (no --cache flag)
→ use list_notes { search: "<keyword>" } as fallback

# GitHub token not configured
→ tell user to add PAT in Settings > Integrations > GitHub
→ do not retry; wait for confirmation

# Missing required field
→ re-read mcp-reference.md for exact required params
→ never guess field names
```
