# Managing JS Projects

A step-by-step guide to using `changelog` in a Node.js project, following the recommended [Git Flow](/git-flow).

## 1. Create a Project

```bash
mkdir my-project && cd my-project
npm init -y
git init
```

## 2. Install Changelog

Install `@d-dev/changelog` as a dev dependency:

```bash
npm install -D @d-dev/changelog
```

## 3. Initialize Changelog

Run the init command:

```bash
npx changelog init
```

Walk through the prompts:

```
? Current version: 0.0.0
? Changelog file path: CHANGELOG.md
? Main git branch: main
Enter post-apply commands (leave blank to finish):
? Post-apply command: npm install
? Post-apply command (1 added):
✓ Initialized .changelog directory
```

Setting `npm install` as a post-apply command ensures `package-lock.json` stays in sync after the version in `package.json` is bumped.

## 4. Add package.json Scripts

Add scripts to your `package.json` for convenient access:

```json
{
  "scripts": {
    "changelog:add": "changelog add",
    "changelog:apply": "changelog apply",
    "changelog:check": "changelog check",
    "changelog:version": "changelog version",
    "changelog:view": "changelog view"
  }
}
```

## 5. Commit Initial Setup

```bash
git add .
git commit -m "chore: initial project setup with changelog"
```

## 6. Git Flow Walkthrough

### Create a Feature Branch

```bash
git checkout -b feature/add-greeting
```

Make some changes to your project — for example, add an `index.js`:

```bash
echo 'console.log("Hello, world!");' > index.js
```

### Add a Changelog Entry

```bash
npm run changelog:add
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

Open a PR on GitHub. CI will run `changelog check` to verify a changelog entry exists (see [CI Setup](#ci-setup) below).

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
npm run changelog:apply
```

```
✓ Updated version in package.json to 0.1.0
✓ Applied version 0.1.0
✓ Updated CHANGELOG.md
Running post-apply command: npm install
? Commit changes and tag the commit? Yes
✓ Committed and tagged version 0.1.0
```

### Push the Release

```bash
git push origin main --follow-tags
```

## CI Setup

Add a GitHub Actions workflow to enforce changelog entries on pull requests:

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
        run: npm run changelog:check

      - name: Lint
        run: npm run lint
```

!> **`fetch-depth: 0` is required.** The `changelog check` command uses `git merge-base` to find the branch fork point. Without full history, this will fail.

?> **Tip:** For PRs with no user-facing changes (refactors, CI tweaks, docs), use the `Internal` change type. It satisfies the check without adding to the public changelog.
