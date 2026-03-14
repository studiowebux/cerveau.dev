---
title: Agents
---

# Agents

Agents are specialized sub-Claude instances that Claude Code can spawn for
focused tasks — code review, bug fixing, architecture planning, etc. They run
with their own system prompt and a restricted tool set.

## File Format

Agents live in `_packages_/studiowebux/minimaldoc/1.0.0/agents/`. For the agent file format and
available fields, see the [Claude Code agents documentation](https://docs.anthropic.com/en/docs/claude-code/sub-agents).

## Declaring Agents in a Brain

In `_configs_/brains.json`, list the packages that provide agents in the
`packages` array:

```json
{
  "name": "MyApp",
  "packages": ["studiowebux/core", "studiowebux/minimaldoc"]
}
```

After editing `brains.json`, run `cerveau rebuild MyApp` to update the
symlinks.

## Writing Agents with Claude

Ask Claude to generate an agent for a specific task:

```
Create a Claude Code agent for reviewing database migrations.
It should check for missing indexes, unsafe column drops, and
missing rollback steps. Restrict tools to Read and Glob.
Save it to _packages_/studiowebux/core/1.0.0/agents/migration-reviewer.md
```

Keep agent system prompts focused. An agent that does one thing well is more
reliable than a general-purpose one.

## Included Agent

The minimaldoc package ships with one agent as a reference:

- `minimaldoc-writer` — writes documentation in MinimalDoc format

Use it as a template when creating your own.

## Ideas to Get Started

Create agents for your stack. Some examples:

```
Create a Claude Code agent for reviewing database migrations.
It should check for missing indexes, unsafe column drops, and missing
rollback steps. Restrict tools to Read and Glob.
Save it to _packages_/studiowebux/core/1.0.0/agents/migration-reviewer.md

Create a Claude Code agent for Go architecture decisions.
It should enforce the repository pattern, OpenTelemetry instrumentation,
and the provider pattern for external services.
Restrict tools to Read, Grep, Glob.
Save it to _packages_/studiowebux/core/1.0.0/agents/golang-architect.md

Create a Claude Code agent for code review.
It should check for security issues, missing error handling, dead code,
and consistency with existing patterns.
Restrict tools to Read, Grep, Glob.
Save it to _packages_/studiowebux/core/1.0.0/agents/code-reviewer.md
```
