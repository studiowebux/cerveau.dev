# Import Project into Cerveau.dev

Onboard a codebase into MD Planner and spawn its brain in one automated flow.

## Arguments

The skill accepts optional inline arguments: `/import-project NAME=MyApp PROJECT=/path/to/code`

If NAME or PROJECT are missing, ask the user before proceeding.

---

## Steps

### 1. Resolve inputs

- Extract `NAME` and `PROJECT` from the skill arguments if provided.
- If either is missing, ask the user:
  - `NAME` — short PascalCase name for the brain (e.g. `HodeiVault`)
  - `PROJECT` — absolute path to the codebase directory
- Verify the directory exists before continuing.

### 2. Spawn the brain

Run from `_protocol_/` (resolve its path relative to cerveau.dev root):

```bash
cd <cerveau-root>/_protocol_ && make onboard NAME=<NAME> PROJECT=<PROJECT>
```

This creates the brain directory, generates `settings.json`, symlinks rules/hooks/agents, connects MCP, and rebuilds selective rules. Stop and report if this fails.

**STOP HERE.** Once `make onboard` succeeds, report to the user:

```
Brain created at: <cerveau-root>/_brains_/<name>-brain

Next step — launch Claude from the brain directory:

  cd <cerveau-root>/_brains_/<name>-brain && claude

Then run /import-project inside that session to complete the setup
(explore codebase, create MD Planner items, fill local-dev.md).
```

Do not proceed to Step 3. The remaining steps run in the brain session, not here.

---

## Brain session steps (run after `cd <brain> && claude`)

### 3. Explore the codebase

Read the project to build a mental model. At minimum:
- `README.md` — purpose, setup, usage
- Any `PLANNER.md`, `ARCHITECTURE.md`, `DESIGN.md`, `ROADMAP.md`, or equivalent
- `package.json`, `pubspec.yaml`, `deno.json`, `go.mod`, `Cargo.toml`, or equivalent — name, version, dependencies
- Top-level directory structure (`find <PROJECT> -maxdepth 2 -not -path '*/.git/*'`)
- Any existing TODO/issue files

Extract:
- What the project does (1–3 sentences)
- Technology stack (language, framework, key dependencies)
- Target platforms
- External services / APIs
- Security model (if relevant)
- Open TODO items
- Logical milestones (current + upcoming)

### 4. Get MD Planner context

```
get_project_config
```

Note the MD Planner project name — this is the `<mcp-project>` for all subsequent calls.

### 5. Create MD Planner items

Create all items scoped to `project: "<NAME>"`. Run independent creates in parallel.

**Portfolio item** — `create_portfolio_item` (always first — other items reference the project by name)
- `name`: project name (must match the `project:` field used on tasks/milestones)
- `description`: 2–3 sentence summary of what the project does
- `status`: `"active"` for ongoing, `"completed"` for finished, `"archived"` for inactive
- `category`: e.g. `"Mobile / Desktop App"`, `"Web App"`, `"CLI Tool"`, `"Library"`, `"API"`
- `github_repo`: `owner/repo` format if found in git remote
- `tech_stack`: array of key technologies inferred from dependencies and config files
- `urls`: bridge server, API endpoints, docs, or any relevant URLs found in the codebase
- `brain_managed`: `true`
- `start_date`: today's date
- `progress`: rough estimate (0–100) based on how complete the project appears

**Brief** — `create_brief`
- `title`: project name
- `summary`: 4–6 bullet points covering purpose, platforms, architecture highlights
- `mission`: 1–2 bullets on the core goal
- `guidingPrinciples`: key design/security/quality principles extracted from docs
- `highLevelTimeline`: phases derived from current state + open items

**Architecture note** — `create_note`
- `title`: `[architecture] <NAME> — Overview`
- `project`: `<NAME>`
- `content`: Full markdown covering system diagram (ASCII), key flows, tech stack table, security model, build targets, external services. Use all details gathered in Step 3.

**Risk entry** — `create_risk`
- `title`: `<NAME> — Risk Analysis`

**Milestones** — `create_milestone` (one per logical phase)
- Mark clearly completed phases as `status: "completed"`
- Open phases as `status: "open"`
- Always set `project: "<NAME>"`

**Tasks** — `create_task` for each open TODO item found in the codebase
- Set `project: "<NAME>"`
- Link to the appropriate milestone
- Write a useful `description` (not just a copy of the TODO line — expand it)
- Set `priority` (1 = critical, 5 = low)
- Add relevant `tags`

### 6. Report

Summarize what was created:
- Brain path + launch command
- MD Planner items created (brief, note, milestones, tasks)
- Any manual follow-up needed (e.g. stacks/practices/workflows to declare in `brains.json`)

---

## Guards

- Never proceed past Step 2 if `make onboard` fails — fix the error first.
- Never create duplicate milestones — check with `list_milestones { project: "<NAME>" }` first if unsure.
- Never invent tasks that aren't grounded in actual code or documentation findings.
- If the project has no README or docs, do a shallow code scan to infer purpose before creating items.
- Keep note content factual — only include what was actually found in the codebase.
