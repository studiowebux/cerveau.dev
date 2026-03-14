# /update — Update Cerveau Protocol

Update the installed Cerveau protocol to the latest version.

## Steps

1. Run the update command:
   ```bash
   cerveau update
   ```

2. Report the result — show the version before and after.

3. If the update succeeded, remind the user to restart Claude Code for changes to take effect.

## Notes

- `.env` and `brains.json` are preserved during update — no credentials or brain config is lost.
- `_brains_/` content (brain directories) is not affected by the update.
- If the download fails, tell the user to check their internet connection and try again.
