# MinimalDoc Setup Workflow

One-shot setup for projects using MinimalDoc documentation + GitHub Pages.

## Prerequisites

- MinimalDoc binary available (clone + build from `studiowebux/minimaldoc`)
- GitHub Pages enabled in repo settings (source: GitHub Actions)
- Project has a `docs/` directory

## Step-by-Step

### 1. Create directory structure

```
docs/
├── config.yaml
├── TOC.md
├── index.md
├── __changelog__/
│   └── releases/
│       └── <version>.md
├── 01-getting-started/
│   ├── 01-installation.md
│   ├── 02-quick-start.md
│   ├── 03-configuration.md
│   └── 04-features.md
├── 02-guides/
│   └── ...
└── 03-reference/
    └── ...
```

### 2. Create `config.yaml`

Use the template from the `minimaldoc` practice rule. Required fields: `title`, `description`, `base_url`, `author`, `theme`, `dark_mode`, `enable_search`, `enable_llms`, `changelog`, `stale_warning`, `social_links`, `footer`.

Include landing page config if the project needs a marketing page:

```yaml
landing:
  enabled: true
  nav: [...]
  hero: { title, subtitle, buttons }
  features: { title, items }
  links: { title, items }
  opensource: { title, description, links }
```

### 3. Create `TOC.md`

Structured markdown list matching the nav order. No frontmatter. See `minimaldoc` practice rule for exact format.

### 4. Create `index.md`

Landing page content with frontmatter:

```yaml
---
title: Project Name
description: One-liner
---
```

### 5. Add frontmatter to every page

Every `.md` file (except `TOC.md`) needs:

```yaml
---
title: Page Title
---
```

### 6. Create changelog releases

One file per version in `docs/__changelog__/releases/<version>.md`:

```yaml
---
version: "1.0.0"
date: "2026-01-15T00:00:00Z"
---
```

### 7. Add GitHub Actions workflow

Copy the workflow template from the `minimaldoc-ci` workflow rule (`_protocol_/.claude/rules/workflow/minimaldoc-ci.md`). Key steps: checkout project, checkout minimaldoc, setup Go, build binary, build docs, deploy to Pages.

### 8. Add features table

Create `docs/01-getting-started/04-features.md` with a two-column table (bold feature name, one-liner description). Copy the same table to the project README under `## Features`. The docs page is canonical.

### 9. Verify locally

```bash
./minimaldoc build docs -o _site
./minimaldoc serve _site
```

Open `http://localhost:8080` and verify: landing page, search, dark mode toggle, changelog, all nav links work.

### 10. Push and verify deployment

Push to main. GitHub Actions builds and deploys to Pages. Verify the live URL matches local preview.

## Reference Implementations

Projects that follow this setup pattern have a `docs/` directory with a landing page, changelog, and GitHub Actions workflow. Use any such project's `docs/` as a structural reference.
