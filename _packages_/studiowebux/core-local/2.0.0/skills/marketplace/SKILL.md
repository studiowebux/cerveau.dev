# /marketplace — Browse and Install Cerveau Packages

Browse the Cerveau marketplace and install packages into a brain.

## Usage

- `/marketplace` — list all available packages
- `/marketplace install <pkg> [brain]` — install a package into a brain

## Steps

### List packages

Run:
```bash
cerveau marketplace list
```

Display the output to the user.

### Install a package

If the user provides a package name (and optionally a brain name):

1. If no brain name given, ask the user which brain to install into (or run `cerveau list` to show available brains).

2. Run:
   ```bash
   cerveau marketplace install <pkg-name> <brain-name>
   ```

3. Report what was added and confirm the rules were rebuilt.

## Notes

- Package types: `workflow`, `practice`, `agent`, `stack` — each maps to the corresponding array in `brains.json`.
- Installing a package that is already present is safe (no-op).
- After install, the rebuild runs automatically — no manual step needed.
