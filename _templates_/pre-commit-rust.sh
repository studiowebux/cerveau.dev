#!/usr/bin/env bash
# pre-commit: run CI checks locally before every commit.
# Skip with: git commit --no-verify

set -euo pipefail

ROOT=$(git rev-parse --show-toplevel)
cd "$ROOT"

if [ ! -f Cargo.toml ]; then
  echo "pre-commit: no Cargo.toml at repo root, skipping Rust checks"
  exit 0
fi

echo "→ cargo fmt --check"
cargo fmt --check

echo "→ cargo clippy"
cargo clippy --all-targets -- -D warnings

echo "→ cargo test"
cargo test

echo "✓ all checks passed"
