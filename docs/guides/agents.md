---
title: Agents
---

# Agents

Agents are specialized sub-Claude instances that Claude Code can spawn for
focused tasks — code review, bug fixing, architecture planning, etc. They run
with their own system prompt and a restricted tool set.

## File Format

Agents are markdown files with YAML frontmatter. For the full format and available fields, see the [Claude Code agents documentation](https://docs.anthropic.com/en/docs/claude-code/sub-agents).

:::warning
**Never save agents into `_packages_/studiowebux/core/` or any downloaded package.** Those files are overwritten by `cerveau update`. Always save your agents to a `_local_` package.
:::

## Writing Agents with Claude

Ask Claude to generate an agent for a specific task:

<details>
<summary>Database migration reviewer</summary>

```
Create a Claude Code agent for reviewing database migrations.
It should check for missing indexes, unsafe column drops, and
missing rollback steps. Restrict tools to Read and Glob.
Save to ~/.cerveau/_packages_/_local_/my-agents/1.0.0/agents/migration-reviewer.md
```

</details>

<details>
<summary>Go architecture reviewer (reusing the golang-stack example)</summary>

```
Create a Claude Code agent for Go architecture decisions.
It should enforce the repository pattern, OpenTelemetry instrumentation,
and the provider pattern for external services.
Restrict tools to Read, Grep, Glob.
Save to ~/.cerveau/_packages_/_local_/golang-stack/1.0.0/agents/go-architect.md
```

</details>

<details>
<summary>Code reviewer</summary>

```
Create a Claude Code agent for code review.
It should check for security issues, missing error handling, dead code,
and consistency with existing patterns.
Restrict tools to Read, Grep, Glob.
Save to ~/.cerveau/_packages_/_local_/my-agents/1.0.0/agents/code-reviewer.md
```

</details>

Keep agent system prompts focused. An agent that does one thing well is more reliable than a general-purpose one.

## Registering and Loading

Register your agents in `~/.cerveau/_configs_/registry.local.json`:

```json
{
  "version": "1.0.0",
  "packages": [
    {
      "name": "my-agents",
      "org": "_local_",
      "version": "1.0.0",
      "path": "_packages_/_local_/my-agents/1.0.0",
      "description": "My custom agents",
      "files": [
        { "name": "migration-reviewer.md", "type": "agents" },
        { "name": "code-reviewer.md",      "type": "agents" }
      ],
      "tags": ["local", "agents"]
    }
  ]
}
```

Add the package to your brain in `brains.json`, then rebuild:

```bash
cerveau rebuild myapp
```

## Included Agent

The `studiowebux/minimaldoc` package ships with one agent as a reference:

- `minimaldoc-writer` — writes documentation in MinimalDoc format

Use it as a template when creating your own.
