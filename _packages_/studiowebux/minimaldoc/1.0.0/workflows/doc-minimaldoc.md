---
paths:
  - "**/docs/**"
  - "**/docs/config.yaml"
---

# MinimalDoc Setup Workflow

One-shot setup for projects using MinimalDoc documentation + GitHub Pages.

> Full MinimalDoc reference: https://minimaldoc.com/llms.txt

## Prerequisites

- MinimalDoc binary: `go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@v1.6.0`
- GitHub Pages enabled in repo settings (source: GitHub Actions)
- Project has a `docs/` directory

## Step-by-Step

### 1. Create directory structure

No numeric prefixes on directories or filenames — `TOC.md` controls navigation order.

```
docs/
├── config.yaml
├── TOC.md
├── __changelog__/
│   └── releases/
│       └── <version>.md
├── getting-started/
│   ├── installation.md
│   ├── quick-start.md
│   ├── configuration.md
│   └── features.md
├── guides/
│   └── ...
└── reference/
    └── ...
```

### 2. Create `config.yaml`

Use the template from the `minimaldoc` practice rule. Required fields: `title`, `description`, `base_url`, `author`, `theme`, `dark_mode`, `enable_search`, `enable_llms`, `changelog`, `stale_warning`, `social_links`, `footer`.

Add a `landing:` block if the project needs a marketing-style landing page — this replaces the need for an `index.md`.

### 3. Create `TOC.md`

Structured markdown list matching the nav order. No frontmatter. See `minimaldoc` practice rule for exact format.

### 4. Add frontmatter to every page

Every `.md` file (except `TOC.md`) needs:

```yaml
---
title: Page Title
---
```

### 5. Create changelog releases

One file per version in `docs/__changelog__/releases/<version>.md`:

```yaml
---
version: "1.0.0"
date: "2026-01-15T00:00:00Z"
---
```

### 6. Add GitHub Actions workflow

Use the workflow template from the `minimaldoc-ci` workflow rule. Key steps: setup Go, install MinimalDoc via `go install`, build docs, deploy to Pages.

### 7. Add features table

Create `docs/getting-started/features.md` with a two-column table (bold feature name, one-liner description). Copy the same table to the project README under `## Features`. The docs page is canonical.

### 8. Verify locally

```bash
minimaldoc build docs -o _site
minimaldoc serve _site
```

Open `http://localhost:8080` and verify: landing page, search, dark mode toggle, changelog, all nav links work.

### 9. Push and verify deployment

Push to main. GitHub Actions builds and deploys to Pages. Verify the live URL matches local preview.
