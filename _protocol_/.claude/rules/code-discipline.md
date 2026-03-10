# Code Discipline Rules

## Core Principle

Every line of code must be production-ready. No exceptions.

## Forbidden

**No partial implementation.** Never write half-finished functions, TODO comments, stub methods, or placeholder components. If you cannot complete it now, do not start it.

**No code sprawl.** Never add features "just in case", create abstractions for single use cases, leave commented-out code, or keep unused imports/variables/functions.

**No complexity spirals.** Stop when adding complexity to fix complexity. Stop when you need workarounds for workarounds. If fixing a bug requires 3+ new conditions, the design is wrong. Refactor instead of patching.

**No premature architecture.** Never add versioning to POCs, plugin systems for 2 implementations, or configuration for hardcoded values. Start simple, refactor when patterns emerge.

## Required

**Complete what you start.** Before starting: can I finish this fully? Do I understand all requirements? If either answer is no, break it down or research first.

**Clean as you go.** After every function: remove unused code. After every feature: delete scaffolding. After every refactor: eliminate the old path. Before every commit: scan for debris.

**Refactor over patch.** When code becomes unwieldy: stop adding conditions, identify the design flaw, refactor the structure, simplify.

**Minimal viable implementation.** Write the least code that solves the problem. Use existing utilities before creating new ones. Inline single-use helpers. Extract only at the Rule of Three.

## Decision Framework

For every piece of code:

1. Is this complete? If no, do not write it.
2. Is this the simplest solution? If no, simplify it.
3. Will this create debt? If yes, redesign it.
4. Is this actually needed? If unsure, do not add it.

## Completion Checklist

- No TODOs in code
- No stubs — every function has real implementation
- No dead code — no commented or unreachable code
- No unused imports, variables, or functions
- No partial features — everything works end-to-end
- Solution matches problem scale

## Consistency — Read Before You Write

Before implementing anything that might already exist elsewhere in the codebase — input handling, error patterns, API calls, UI interactions — grep for it first.

**Read every existing implementation of the same pattern. Match it exactly.**

If implementations are inconsistent, normalize them all in the same commit. Never add a third variant when two already exist.

This applies universally: copy/paste handling, form inputs, network calls, state mutations, event handlers, modal lifecycle — everything. "I'll do it quickly" is how inconsistency accumulates. The cost of grepping first is zero. The cost of a codebase with three ways to do the same thing compounds forever.

## Red Flags — Stop Immediately

"I'll clean this up later." "This is just temporary." "Let me add a TODO." "Just need one more workaround." "This will work for now." "I'll just do it the quick way this time."

If you think these thoughts: stop. Fix it properly or do not do it.
