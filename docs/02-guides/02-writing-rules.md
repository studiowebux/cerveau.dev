---
title: Writing Rules
---

# Writing Rules

The protocol ships with no rules — you generate them for your stack. Claude
Code can write them for you.

## Rule Types

| Type | Directory | Loaded when |
|---|---|---|
| **Stack** | `_protocol_/.claude/rules/stack/` | Brain declares the stack |
| **Practice** | `_protocol_/.claude/rules/practices/` | Brain declares the practice |
| **Workflow** | `_protocol_/.claude/rules/workflow/` | Brain declares the workflow |
| **Core** | `_protocol_/.claude/rules/` | Always — every brain |

## Scoping Rules to File Patterns

By default a rule loads every session. Add `paths:` frontmatter to scope it so
Claude Code only loads it when a matching file is open in context:

```yaml
---
paths:
  - "**/*_test.go"
  - "**/*.test.ts"
  - "**/tests/**"
---

Write useful tests. Do not over-test...
```

Rules without `paths:` load in every session regardless of what files are open.
Rules with `paths:` load only when a matching file is active.

Use path-scoped rules for concerns tied to specific file types: testing rules for
test files, CI rules for workflow files, keybinding rules for keymap files.
Keep general practices (error handling, architecture, security) as always-loaded.

## Generate with Claude

Open Claude Code inside `cerveau.dev/` and use these prompts:

### Stack rules

```
Create a Claude Code rule file for Go development.
Enforce: explicit error handling, table-driven tests, slog for logging,
no global mutable state, go fmt before every commit.
Keep under 80 lines. Rules only — no examples, no prose.
```

```
Create a Claude Code rule file for TypeScript with Deno.
Enforce: Zod for validation, explicit return types on exports,
no any, Deno.test for tests, deno fmt before commits.
Keep under 80 lines. Rules only — no examples, no prose.
```

### Practice rules

```
Create a Claude Code rule file for code review practices.
Enforce: security-first review, error handling coverage, flag dead code,
require PR descriptions with test plans.
Keep under 80 lines. Rules only — no examples, no prose.
```

```
Create a Claude Code rule file for testing practices.
Enforce: unit tests for all exported functions, table-driven tests,
no test-only logic in production code, mocks over real dependencies in unit tests.
Keep under 80 lines. Rules only — no examples, no prose.
```

### Workflow rules

```
Create a Claude Code rule file for git workflow.
Enforce: feature branches, conventional commits (feat/fix/chore/docs/test),
no force push to main, PR-based merges only.
Keep under 80 lines. Rules only — no examples, no prose.
```

## Declare in brains.json

Add the rule name (filename without `.md`) to the brain's array:

```json
{
  "name": "MyApp",
  "path": "_brains_/myapp-brain",
  "codebase": "_projects_/myapp",
  "isCore": false,
  "stacks": ["go"],
  "practices": ["testing", "code-review"],
  "workflows": ["git", "mdplanner-tasks", "local-dev"],
  "agents": ["goal-planner"]
}
```

Then rebuild:

```bash
./_scripts_/rebuild-brain-rules.sh MyApp
```

## Evolving Rules

As your project evolves, ask Claude to refine:

```
Read my codebase and create a stack rule that captures the patterns and
conventions we're using. Focus on what's unique to this project.
```

```
We keep making the same mistakes in PRs. Create a code-review practice rule
that catches [specific issues].
```

## Best Practices

- Keep rules short — under 80 lines per file
- Rules only — no examples, no tutorials, no prose
- One file per concern — many small rules beats one large generic rule
- Language-agnostic practices — put language tips in stack rules
- Every line costs tokens — eliminate anything that doesn't enforce a rule
