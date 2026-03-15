---
title: Writing Rules
---

# Writing Rules

The core package ships with no rules — you generate them for your stack. Claude
Code can write them for you.

## File Types

When registering files in `registry.local.json` (or any package), the `type` field controls where the file is installed in the brain and how Claude Code loads it.

| Type | Installed to (inside brain) | Description |
|---|---|---|
| `rules` | `.claude/rules/` | Always-loaded rules — apply every session regardless of open files |
| `stacks` | `.claude/rules/stack/` | Language or framework conventions (Go, TypeScript, Python…) |
| `practices` | `.claude/rules/practices/` | How you work — testing, code review, security, architecture |
| `workflows` | `.claude/rules/workflow/` | Process rules — git flow, commit format, release steps |
| `hooks` | `.claude/hooks/` | Shell scripts triggered by Claude Code events (PreToolUse, PostToolUse…) |
| `skills` | `.claude/skills/` | Slash commands Claude can execute on demand (`/my-skill`) |
| `agents` | `.claude/agents/` | Subagents Claude can spawn for specialized tasks |
| `templates` | `templates/` | Markdown templates for notes, decisions, specs, and reports |

## Scoping Rules to File Patterns

> See the official [Claude Code memory documentation](https://docs.anthropic.com/en/docs/claude-code/memory) for the full reference on rule scoping and frontmatter.

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

Open Claude Code and use these prompts. Each one includes the save path so Claude writes directly to the right `_local_` package location.

### Stack rules

<details>
<summary>Go stack rules</summary>

```
Create a Claude Code rule file for Go development.
Enforce: explicit error handling, table-driven tests, slog for logging,
no global mutable state, go fmt before every commit.
Keep under 80 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/golang-stack/1.0.0/stacks/go-stack.md
```

</details>

<details>
<summary>TypeScript / Deno stack rules</summary>

```
Create a Claude Code rule file for TypeScript with Deno.
Enforce: Zod for validation, explicit return types on exports,
no any, Deno.test for tests, deno fmt before commits.
Keep under 80 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/ts-stack/1.0.0/stacks/ts-stack.md
```

</details>

### Practice rules

<details>
<summary>Code review practices</summary>

```
Create a Claude Code rule file for code review practices.
Enforce: security-first review, error handling coverage, flag dead code,
require PR descriptions with test plans.
Keep under 80 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/my-practices/1.0.0/practices/code-review.md
```

</details>

<details>
<summary>Testing practices</summary>

```
Create a Claude Code rule file for testing practices.
Enforce: unit tests for all exported functions, table-driven tests,
no test-only logic in production code, mocks over real dependencies in unit tests.
Keep under 80 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/my-practices/1.0.0/practices/testing.md
```

</details>

### Workflow rules

<details>
<summary>Git workflow</summary>

```
Create a Claude Code rule file for git workflow.
Enforce: feature branches, conventional commits (feat/fix/chore/docs/test),
no force push to main, PR-based merges only.
Keep under 80 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/my-practices/1.0.0/workflows/git.md
```

</details>

### registry.local.json

Once your files are generated, register all three packages so `cerveau rebuild` can find them:

<details>
<summary>~/.cerveau/_configs_/registry.local.json</summary>

```json
{
  "version": "1.0.0",
  "packages": [
    {
      "name": "golang-stack",
      "org": "_local_",
      "version": "1.0.0",
      "path": "_packages_/_local_/golang-stack/1.0.0",
      "description": "Go language conventions",
      "files": [
        { "name": "go-stack.md", "type": "stacks" }
      ],
      "tags": ["local", "go"]
    },
    {
      "name": "ts-stack",
      "org": "_local_",
      "version": "1.0.0",
      "path": "_packages_/_local_/ts-stack/1.0.0",
      "description": "TypeScript / Deno language conventions",
      "files": [
        { "name": "ts-stack.md", "type": "stacks" }
      ],
      "tags": ["local", "typescript", "deno"]
    },
    {
      "name": "my-practices",
      "org": "_local_",
      "version": "1.0.0",
      "path": "_packages_/_local_/my-practices/1.0.0",
      "description": "Code review, testing, and git workflow practices",
      "files": [
        { "name": "code-review.md", "type": "practices" },
        { "name": "testing.md",     "type": "practices" },
        { "name": "git.md",         "type": "workflows" }
      ],
      "tags": ["local", "practices"]
    }
  ]
}
```

</details>

## Declare in brains.json

Add each brain with its packages to `~/.cerveau/_configs_/brains.json`. Multiple brains can share the same `_local_` packages — each one only loads what it declares.

```json
{
  "brains": [
    {
      "name": "myapp",
      "path": "_brains_/myapp-brain",
      "codebase": "/home/user/projects/myapp",
      "packages": ["studiowebux/core", "_local_/golang-stack", "_local_/my-practices"]
    },
    {
      "name": "myapp-frontend",
      "path": "_brains_/myapp-frontend-brain",
      "codebase": "/home/user/projects/myapp-frontend",
      "packages": ["studiowebux/core", "_local_/ts-stack", "_local_/my-practices"]
    },
    {
      "name": "myapp-infra",
      "path": "_brains_/myapp-infra-brain",
      "codebase": "/home/user/projects/myapp-infra",
      "packages": ["studiowebux/core", "_local_/my-practices"]
    }
  ]
}
```

Then rebuild whichever brain you updated:

```bash
cerveau rebuild myapp
cerveau rebuild myapp-frontend
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
