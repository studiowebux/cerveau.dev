# Release Workflow

Execute a full release for this project. Follow every step in order.

## Steps

1. **Verify clean state.** Run `git status`. If there are uncommitted changes, commit them first with a descriptive message. Never release from a dirty working directory.

2. **Determine version.** Read the current version from the project's version file (check `deno.json`, `package.json`, `version.go`, or equivalent). Scan `git log` since the last tag to determine the bump: patch for fixes only, minor for features, major for breaking changes. **Confirm the version number with the user before proceeding.**

3. **Update CHANGELOG.md.** Move `[Unreleased]` items to a new version section with today's date. Follow Keep a Changelog format. Update comparison links at the bottom.

4. **Update MinimalDoc changelog (if docs/ uses MinimalDoc).** Check if `docs/config.yaml` exists and has `changelog: enabled: true`. If yes, create `docs/__changelog__/releases/<version>.md` following the format in the `minimaldoc` practice rule:
   - `version:` and `date:` must be **quoted strings**
   - `date:` must be **RFC3339** format: `"YYYY-MM-DDT00:00:00Z"` — plain `YYYY-MM-DD` causes ordering failures when multiple releases share the same calendar date
   - When multiple releases share the same calendar date, stagger timestamps by hour (lower version → lower hour offset)
   - Content: H2 headings (`## Added`, `## Changed`, `## Fixed`, `## Removed`, `## Security`) with bullet items drawn from CHANGELOG.md

5. **Bump version.** Update the version in all relevant files (deno.json, package.json, version constants, Info.plist, etc.).

5. **Build and verify.** Run the project's build command (`go build ./...`, `deno task build`, `npm run build`, etc.) and ensure it passes. Do not proceed if the build fails.

6. **Commit.** Stage all version-related changes and commit: `chore: release vX.Y.Z`

7. **Tag.** `git tag vX.Y.Z`

8. **Push.** `git push && git push --tags`

9. **GitHub Release.** Create a release from the tag with the changelog excerpt as body: `gh release create vX.Y.Z --title "vX.Y.Z" --notes "$(changelog excerpt)"`

10. **Update brain.** Write a `[progress]` note to `notes/` summarizing the release. Update the active milestone in `context.md` and `local-dev.md` if all tasks are done.

## Guards

- Never skip the build verification step
- Never tag without user confirmation on the version number
- Never release if CI is red
- Never release if there are uncommitted changes
