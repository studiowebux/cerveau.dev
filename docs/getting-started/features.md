---
title: Features
---

# Features

| Feature                      | Description                                                                                                  |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------ |
| **Multi-Brain**              | One protocol, many projects. Each brain gets selective rules, hooks, and agents — no duplication, no drift.  |
| **Selective Loading**        | Declare exactly which stacks, practices, and workflows each brain needs. Only those rules load into context. |
| **MDPlanner Integration**    | Single source of truth for tasks, notes, decisions, and progress. Claude reads and writes it every session.  |
| **Hooks Enforcement**        | Boot reminders, commit validation, checkpoint checks, and progress gates. The protocol runs automatically.   |
| **Protocol-Driven Workflow** | Four phases. Boot → Work → Write Back → Close. Zero drift between sessions.                                  |
| **Zero Footprint**           | No files added to your project repos. The brain lives outside your code. Your codebase stays clean.          |
| **Bring Your Own Rules**     | No rules ship by default. Generate stack, practice, and workflow rules with a single Claude prompt.          |
| **Marketplace**              | Browse and install community packages (workflows, practices, agents) into any brain with one command. Filter by text, tag, or org. |
| **Shell Completions**        | Tab-tab support for all commands, brain names, packages, tags, and orgs. Includes `cerveau cd` shell wrapper. |
| **Boot from Anywhere**       | `cerveau boot <name>` launches Claude Code in any brain from any directory. No `cd` needed.                  |
| **Backup & Restore**         | `cerveau backup` archives `~/.cerveau/`, `~/.claude/`, and MDPlanner data. Selective with `--cerveau`, `--mdplanner`, `--claude`. |
| **Auto-Update**              | Version check on session start. Update the protocol in place — `.env`, brains, and config are preserved.     |
| **One-Line Install**         | `curl -fsSL https://cerveau.dev/install.sh \| bash` — installs protocol, starts MDPlanner, wires MCP.         |
| **Open Source**              | AGPL-3.0. Self-host everything. MDPlanner runs in a container. No external dependencies required.            |
| **Agent Support**            | Define custom agents in YAML. Declare which agents each brain loads. Agents live in the protocol.            |
| **Skills Support**           | Reusable skill definitions (slash commands) shared across all brains from the protocol.                      |
