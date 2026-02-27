# Git Flow

Recommended workflow for using `changelog` with a branch-based git flow.

## Overview

The core idea is simple: changelog entries are added **on feature branches** alongside the code changes that produce them, and releases are cut **on the main branch** after merging.

```
main        ──●──────────────●── apply ──●──▶
              │              │
feature       └──● add ──● ──┘
```

## Workflow

### 1. Branch

Create a feature or fix branch from main:

```bash
git checkout -b feature/my-new-feature
```

### 2. Add Entries

After making your changes, record them with `changelog add`:

```bash
changelog add
```

```
? Change type: Add - Add a new feature. Minor version bump.
? Description: Add support for custom templates
✓ Added changelog entry: .changelog/next/1740000000000.yaml
✓ Updated CHANGELOG.md
```

Commit the entry file alongside your code changes:

```bash
git add .
git commit -m "feat: add support for custom templates"
```

You can run `changelog add` multiple times on the same branch if you make several notable changes.

### 3. Pull Request

Open a pull request. CI runs [`changelog check`](/commands?id=changelog-check) to verify that at least one changelog entry exists on the branch:

```bash
changelog check
# Check passed: found 1 changelog entry in .changelog/next/
```

If no entries are found, the check fails and the PR is blocked.

### 4. Merge

Once the PR is approved and CI passes, merge it into main. The changelog entry files in `.changelog/next/` accumulate as branches are merged.

### 5. Release

When you're ready to release, run `changelog apply` on main:

```bash
git checkout main
git pull
changelog apply
```

This bumps the version, generates the changelog, and prompts to commit and tag:

```
✓ Applied version 1.1.0
✓ Updated CHANGELOG.md
? Commit changes and tag the commit? Yes
✓ Committed and tagged version 1.1.0
```

### 6. Push

Push the commit and tag to the remote:

```bash
git push origin main --follow-tags
```

## CI: Enforcing Changelog Entries on PRs

Use a GitHub Actions workflow to run `changelog check` on every pull request. This ensures that no PR is merged without a changelog entry.

Here is an example workflow:

```yaml
# .github/workflows/pr.yml
name: PR
on:
  pull_request:
    branches:
      - main

permissions:
  contents: read
  pull-requests: read
  checks: write

jobs:
  pr:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "22"

      - name: Install dependencies
        run: npm ci

      - name: Check Changelog
        run: npx changelog check
```

!> **`fetch-depth: 0` is critical.** The `changelog check` command uses `git merge-base` to find where the branch diverged from main. Without full git history, this will fail.

## Handling Non-User-Facing Changes

Not every PR contains changes that belong in a public changelog. For refactors, CI updates, documentation changes, or dependency bumps, use the **Internal** change type:

```bash
changelog add
```

```
? Change type: Internal - Internal change that does not affect users...
? Description: Refactor authentication module
```

`Internal` entries:

- ✅ Satisfy the `changelog check` CI gate
- ✅ Are tracked in `.changelog/next/`
- ❌ Do **not** appear in the public changelog
- ❌ Do **not** trigger a version bump on their own

?> **Tip:** `Internal` is the escape hatch for PRs that must pass CI but have no user-facing impact.
