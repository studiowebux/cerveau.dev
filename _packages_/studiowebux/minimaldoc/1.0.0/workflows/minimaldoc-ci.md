---
paths:
  - ".github/workflows/**"
  - "**/docs/**"
---

# MinimalDoc Docs CI/CD

GitHub Actions workflow for building and deploying MinimalDoc-powered
documentation sites to GitHub Pages.

## Install method

Use `go install` with the full subpackage path and a pinned version:

```bash
go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@v1.7.0
```

## Workflow template

```yaml
name: Deploy Documentation

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/workflows/docs.yml'
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Install MinimalDoc
        run: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@v1.7.0

      - name: Build Documentation
        run: minimaldoc build docs -o public

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v4
        with:
          path: ./public/

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Notes

- Pin the minimaldoc version — do not use `@latest` in CI
- Build command: `minimaldoc build <source-dir> -o <output-dir>`
- Enable GitHub Pages in repo settings with source set to "GitHub Actions"
- The workflow needs `pages: write` and `id-token: write` permissions
