#!/bin/sh
# Pre-commit hook — blocks commits that fail fmt, lint, or type check.
# Runs from the repo root.

set -e

echo "[pre-commit] Scanning for secrets..."
PATTERNS='(sk-[a-zA-Z0-9]{20,}|ghp_[a-zA-Z0-9]{36,}|AKIA[0-9A-Z]{16}|password\s*=\s*["'"'"'][^"'"'"']+["'"'"'])'
MATCHES=$(grep -rEn "$PATTERNS" --include='*.ts' --include='*.js' --include='*.json' --include='*.yml' --include='*.yaml' --include='*.env' --exclude-dir=node_modules --exclude-dir=.git . || true)
if [ -n "$MATCHES" ]; then
  echo "[pre-commit] FAILED: Potential secrets detected:"
  echo "$MATCHES"
  exit 1
fi

echo "[pre-commit] Running deno fmt --check..."
deno fmt --check
if [ $? -ne 0 ]; then
  echo "[pre-commit] FAILED: deno fmt --check. Run 'deno fmt' to fix."
  exit 1
fi

echo "[pre-commit] Running deno lint..."
deno lint
if [ $? -ne 0 ]; then
  echo "[pre-commit] FAILED: deno lint. Fix lint errors before committing."
  exit 1
fi

echo "[pre-commit] Running deno check v2/bin.ts..."
deno check v2/bin.ts
if [ $? -ne 0 ]; then
  echo "[pre-commit] FAILED: deno check. Fix type errors before committing."
  exit 1
fi

echo "[pre-commit] Running deno task test..."
deno task test
if [ $? -ne 0 ]; then
  echo "[pre-commit] FAILED: deno task test. Fix failing tests before committing."
  exit 1
fi

echo "[pre-commit] All checks passed."
