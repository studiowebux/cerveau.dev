# Claude Code — Session Rules

## Mandatory Gates

1. **Boot first.** Run **Phase 1 — Boot** from the brain's `CLAUDE.md` before responding to the user. No exceptions.
2. **Ticket before work.** No code changes, no subagents, no deep exploration without an mdplanner task. Light reads to write the ticket description are the only exception. See Phase 2 in the brain's `CLAUDE.md`.

## Where to Find Context

Do NOT guess or assume — read the relevant document first.

| When you need... | Read this |
|---|---|
| **Session phases & protocol** | Brain's `CLAUDE.md` (via additionalDirectories) |
| **MCP tool reference** | Brain's `mcp-reference.md` |
| **MCP workflow patterns** | Brain's `mcp-workflows.md` |
| **Code quality** | `@.claude/rules/code-discipline.md` |
| **Goal management** | `@.claude/rules/goal-discipline.md` |
| **Stack rules** | `@.claude/rules/stack/*.md` |
| **Practice rules** | `@.claude/rules/practices/*.md` |
| **Workflow rules** | `@.claude/rules/workflow/*.md` |
| **Project decisions** | mdplanner: `list_notes { search: "[decision] __PROJECT__" }` |
| **Architecture** | mdplanner: `list_notes { search: "[architecture] __PROJECT__" }` |

## Source of Truth

mdplanner notes > local files > memory. Never guess — search mdplanner first, ask user second, guess last. Decisions are append-only: never edit a `[decision]` note, create a superseding one.
