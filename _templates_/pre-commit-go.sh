#!/usr/bin/env bash
# pre-commit: run CI checks locally before every commit.
# Mirrors .github/workflows/checks.yml exactly.
# Skip with: git commit --no-verify

set -euo pipefail

GOPATH=$(go env GOPATH)
ROOT=$(git rev-parse --show-toplevel)
cd "$ROOT"

echo "→ go build"
go build -trimpath -ldflags="-s -w" -o /dev/null . 2>&1

echo "→ go vet"
go vet ./... 2>&1

echo "→ staticcheck"
"$GOPATH/bin/staticcheck" ./... 2>&1

echo "→ gosec"
"$GOPATH/bin/gosec" -severity medium ./... 2>&1

echo "→ govulncheck"
"$GOPATH/bin/govulncheck" ./... 2>&1

echo "✓ all checks passed"
