# Cerveau.dev

Multi-brain system for Claude Code. One protocol, many projects, zero
duplication.

Each project gets its own brain with selective rules, hooks, agents, and skills
loaded from a shared protocol. MDPlanner is the single source of truth for
tasks, notes, decisions, and progress.

Bug Tracker: https://github.com/studiowebux/cerveau.dev/issues
<br>
Discord: https://discord.gg/BG5Erm9fNv

## Features

| Feature                      | Description                                                                                                  |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------ |
| **Multi-Brain**              | One protocol, many projects. Each brain gets selective rules, hooks, and agents — no duplication, no drift.  |
| **Selective Loading**        | Declare exactly which stacks, practices, and workflows each brain needs. Only those rules load into context. |
| **MDPlanner Integration**    | Single source of truth for tasks, notes, decisions, and progress. Claude reads and writes it every session.  |
| **Hooks Enforcement**        | Boot reminders, commit validation, checkpoint checks, and progress gates. The protocol runs automatically.   |
| **Protocol-Driven Workflow** | Four phases. Boot → Work → Write Back → Close. Zero drift between sessions.                                  |
| **Zero Footprint**           | No files added to your project repos. The brain lives outside your code. Your codebase stays clean.          |
| **Bring Your Own Rules**     | No rules ship by default. Generate stack, practice, and workflow rules with a single Claude prompt.          |
| **Open Source**              | AGPL-3.0. Self-host everything. MDPlanner runs in a container. No external dependencies required.            |
| **Agent Support**            | Define custom agents in YAML. Declare which agents each brain loads. Agents live in the protocol.            |
| **Skills Support**           | Reusable skill definitions (slash commands) shared across all brains from the protocol.                      |

## Documentation

Full documentation at **https://cerveau.dev** — installation, quick start, guides, and reference.

## License

AGPL-3.0
