# Example Stack Rule

Add your stack-specific rules here. One file per stack.

Example filenames: `go.md`, `typescript.md`, `python.md`, `rust.md`

Declare which stacks a brain uses in `_configs_/brains.json`:

```json
{
  "stacks": ["go", "typescript"]
}
```

Only declared stacks get symlinked into the brain.
