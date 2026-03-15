---
title: Quick Start
---

# Quick Start

Five steps from zero to a running brain session.

---

## Step 1 — Install

```bash
# Default install
curl -fsSL https://cerveau.dev/install.sh | bash

# Custom home directory or port
CERVEAU_HOME=/opt/cerveau MCP_PORT=9000 curl -fsSL https://cerveau.dev/install.sh | bash
```

This installs Cerveau to `~/.cerveau/`, starts MDPlanner, and registers the MCP globally. See [Installation](installation.md) for prerequisites.

Verify:

```bash
curl -s http://localhost:8003/health
# expected: {"status":"ok"}
```

**Set up MDPlanner people.** Open the MDPlanner UI at `http://localhost:8003` and create at least two people: yourself (the project owner) and Claude (the AI agent). Their IDs are used to populate the People Registry in `local-dev.md` during first boot — without them, Claude cannot assign tasks or track ownership correctly. See [MDPlanner Setup](../guides/mdplanner.md) for details.

---

## Step 2 — Write Your Rules (Optional)

:::question
Not sure what rules you need yet? Skip this step. After spawning the brain in Step 3, ask Claude to read your codebase and propose rules, agents, and practices. Refine them together, save the results into a `_local_` package as described here, register it in `registry.local.json`, add it to your brain in `brains.json`, then run `cerveau rebuild <name>` to apply.
:::

Your own rules go in a `_local_` package — never inside a downloaded package like `studiowebux/core`. Local packages follow the same structure as community packages.

**1. Create the package directory:**

```bash
mkdir -p ~/.cerveau/_packages_/_local_/golang-stack/1.0.0/{stacks,practices,workflows,agents}
```

**2. Generate your rule files** — open Claude Code and ask for each:

<details>
<summary>go-stack.md — language conventions</summary>

```
Create a Claude Code rule file for Go stack conventions.
Enforce: explicit error handling, slog for logging, no global mutable state,
context propagation, no init() functions.
Keep under 80 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/golang-stack/1.0.0/stacks/go-stack.md
```

</details>

<details>
<summary>go-practices.md — testing practices</summary>

```
Create a Claude Code rule file for Go testing practices.
Enforce: table-driven tests, testify for assertions, test file naming _test.go,
subtests with t.Run, no t.Parallel() unless proven safe.
Keep under 80 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/golang-stack/1.0.0/practices/go-practices.md
```

</details>

<details>
<summary>go-checks.md — pre-commit workflow checks</summary>

```
Create a Claude Code rule file for Go pre-commit checks.
Enforce: go fmt, go vet, go mod tidy, staticcheck before every commit.
Keep under 40 lines. Rules only — no examples, no prose.
Save to ~/.cerveau/_packages_/_local_/golang-stack/1.0.0/workflows/go-checks.md
```

</details>

<details>
<summary>go-architect.md — architecture review agent</summary>

```
Create a Claude Code agent file for Go architecture review.
The agent reviews Go package structure, dependency direction, interface design,
and flags violations of clean architecture. Output a short bullet list of findings.
Save to ~/.cerveau/_packages_/_local_/golang-stack/1.0.0/agents/go-architect.md
```

</details>

Save the files into the package directory:

```
~/.cerveau/_packages_/_local_/golang-stack/1.0.0/
  stacks/
    go-stack.md        ← language conventions
  practices/
    go-practices.md    ← testing practices
  workflows/
    go-checks.md       ← pre-commit workflow checks
  agents/
    go-architect.md    ← architecture review agent
```

**3. Register the package** — create `~/.cerveau/_configs_/registry.local.json`:

```json
{
  "version": "1.0.0",
  "packages": [
    {
      "name": "golang-stack",
      "org": "_local_",
      "version": "1.0.0",
      "path": "_packages_/_local_/golang-stack/1.0.0",
      "description": "My Go language rules, testing practices, and pre-commit checks",
      "files": [
        { "name": "go-stack.md",     "type": "stacks" },
        { "name": "go-practices.md", "type": "practices" },
        { "name": "go-checks.md",    "type": "workflows" },
        { "name": "go-architect.md", "type": "agents" }
      ],
      "tags": ["local", "go", "golang"]
    }
  ]
}
```

Available file types: `rules`, `stacks`, `practices`, `workflows`, `hooks`, `skills`, `agents`, `templates`.

See [Writing Rules](../guides/writing-rules.md) for more prompts.

---

## Step 3 — Spawn a Brain

```bash
# Core only (default — omitting --packages also defaults to studiowebux/core)
cerveau spawn MyApp /absolute/path/to/your/code --packages studiowebux/core

# Core + your local package
cerveau spawn MyApp /absolute/path/to/your/code --packages studiowebux/core,_local_/golang-stack
```

:::warning
`studiowebux/core` must always be included. It provides the session protocol, hooks, and boot rules that everything else depends on. A brain without it won't function correctly.
:::

This spawns the brain, wires MCP (already global from the install), and rebuilds selective rules in one step.

Verify no placeholders remain:

```bash
cerveau validate MyApp
```

:::info
`~/.cerveau/_packages_/` contains the shared packages. All project work happens from the brain session in `~/.cerveau/_brains_/myapp-brain/`.
:::

:::question
Added a local package after spawning? Register it in `registry.local.json`, add it to `brains.json`, then run `cerveau rebuild MyApp` to sync. Never edit the brain's `.claude/` files directly — they are managed by rebuild and will be overwritten.
:::

---

## Step 4 — Launch the Brain Session

```bash
cerveau boot MyApp
```

This launches Claude Code inside the brain directory — works from anywhere. To enable tab completion for brain names and other commands, add this to your shell config:

```bash
eval "$(cerveau completion zsh)"    # .zshrc
eval "$(cerveau completion bash)"   # .bashrc
```

Or manually without the `boot` command:

```bash
cd ~/.cerveau/_brains_/myapp-brain && claude
```

---

## Step 5 — Boot

Inside the brain session, type `boot`. Claude will read the rules,
explore the codebase, fill `local-dev.md`, and set up MDPlanner state.

Type `boot` at the start of every session to trigger Phase 1.

:::warning
**`local-dev.md` is the pointer between the brain and your codebase.** It tells Claude where the code lives, who owns the project, and what commands to run. Keep it concise and accurate — the rest (tasks, decisions, progress) comes from MDPlanner, covered in the next section. An incomplete or stale `local-dev.md` leads to wrong paths and wasted context; an over-stuffed one duplicates what MDPlanner already tracks.
:::

---

## Next

→ [First Brain](first-brain.md) — what happens inside the first session
