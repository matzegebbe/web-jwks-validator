# Release Please guidance

Release Please relies on Conventional Commits when generating release pull requests and
changelogs. To avoid parse failures:

1. Use **squash merge** (or rebase merge) for pull requests so the final commit message is
   a single Conventional Commit (e.g. `feat: add xyz`).
2. Ensure PR titles follow Conventional Commit format, since GitHub can use the PR title
   as the squash commit message.
3. Avoid merge commits like `Merge pull request ...`, which do not follow Conventional
   Commits and can cause parsing errors.

## Manual release (if needed)

If you need to cut a release manually, you can:

1. Create a tag on `main` (e.g. `v0.0.13`) and push it.
2. Create a GitHub Release from that tag with release notes.

This bypasses Release Please, but you should return to Conventional Commits afterward
so automated releases resume working.
