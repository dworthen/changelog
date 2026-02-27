# Managing Python Projects

A step-by-step guide to using `changelog` (published as `changesets` on PyPI) in a Python project using [uv](https://docs.astral.sh/uv/) and [Poe the Poet](https://poethepoet.natn.io/), following the recommended [Git Flow](/git-flow).

## 1. Create a Project

Use uv to scaffold a new Python project:

```bash
uv init my-project
cd my-project
git init
```

This creates a `pyproject.toml`, `src/` directory, and other standard files.

## 2. Install Changelog

Install `changesets` as a global tool with uv:

```bash
uv tool install changesets
```

This makes the `changesets` command available globally. Alternatively, you can run it on-demand without installing:

```bash
uvx changesets <command>
```

## 3. Install Poe the Poet

Add `poethepoet` as a dev dependency to define project tasks:

```bash
uv add --dev poethepoet
```

## 4. Initialize Changelog

Run the init command:

```bash
changesets init
```

Walk through the prompts:

```
? Current version: 0.0.0
? Changelog file path: CHANGELOG.md
? Main git branch: main
Enter post-apply commands (leave blank to finish):
? Post-apply command: uv lock
? Post-apply command (1 added):
✓ Initialized .changelog directory
```

Setting `uv lock` as a post-apply command ensures `uv.lock` stays in sync after the version in `pyproject.toml` is bumped.

## 5. Define Poe Tasks

Add task definitions to your `pyproject.toml`:

```toml
[tool.poe.tasks]
changelog-add.cmd = "changesets add"
changelog-add.help = "Add a new changelog entry"

changelog-apply.cmd = "changesets apply"
changelog-apply.help = "Apply pending changelog entries and bump version"

changelog-check.cmd = "changesets check"
changelog-check.help = "Check that at least one changelog entry exists"

changelog-version.cmd = "changesets version"
changelog-version.help = "Print the current version"

changelog-view.cmd = "changesets view"
changelog-view.help = "View the changelog"
```

You can now run tasks with `poe`:

```bash
poe changelog-add
poe changelog-apply
poe changelog-check
```

## 6. Commit Initial Setup

```bash
git add .
git commit -m "chore: initial project setup with changelog"
```

## 7. Git Flow Walkthrough

### Create a Feature Branch

```bash
git checkout -b feature/add-greeting
```

Make some changes — for example, add a module:

```bash
cat > src/my_project/greet.py << 'EOF'
def greet(name: str) -> str:
    return f"Hello, {name}!"
EOF
```

### Add a Changelog Entry

```bash
poe changelog-add
```

```
? Change type: Add - Add a new feature. Minor version bump.
? Description: Add greeting module
✓ Added changelog entry: .changelog/next/1740000000000.yaml
✓ Updated CHANGELOG.md
```

### Commit and Push

```bash
git add .
git commit -m "feat: add greeting module"
git push origin feature/add-greeting
```

### Open a Pull Request

Open a PR on GitHub. CI will run `changesets check` to verify a changelog entry exists (see [CI Setup](#ci-setup) below).

```
✓ Check passed: found 1 changelog entry in .changelog/next/
```

### Merge the PR

Once approved, merge the PR into main.

### Cut a Release

On main, pull the latest changes and apply:

```bash
git checkout main
git pull
poe changelog-apply
```

```
✓ Updated version in pyproject.toml to 0.1.0
✓ Applied version 0.1.0
✓ Updated CHANGELOG.md
Running post-apply command: uv lock
? Commit changes and tag the commit? Yes
✓ Committed and tagged version 0.1.0
```

### Push the Release

```bash
git push origin main --follow-tags
```

## CI Setup

Add a GitHub Actions workflow to enforce changelog entries on pull requests. Using `uvx`, you don't need to install `changesets` as a project dependency in CI:

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

      - name: Install uv
        uses: astral-sh/setup-uv@v7

      - name: Check Changelog
        run: uvx changesets check
```

!> **`fetch-depth: 0` is required.** The `changesets check` command uses `git merge-base` to find the branch fork point. Without full history, this will fail.

?> **Tip:** For PRs with no user-facing changes (refactors, CI tweaks, docs), use the `Internal` change type. It satisfies the check without adding to the public changelog.
