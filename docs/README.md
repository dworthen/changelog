# Quick Start

Changelog is a CLI tool for managing changelogs. Individual changelog entries are stored as files, eliminating merge conflicts and enabling CI enforcement.

## Installation

<!-- tabs:start -->

#### **npm**

Install globally:

```bash
npm install -g @d-dev/changelog
```

Or as a dev dependency:

```bash
npm install -D @d-dev/changelog
```

Or run with `npx`

```bash
npx @d-dev/changelog
```

#### **PyPI**

Install with pip:

```bash
pip install changesets
```

Install with uv.

```bash
uv tool install changesets
```

Or as a dev dependency

```bash
uv add --dev changesets
```

Or run directly with uvx (no install required):

```bash
uvx --from changesets changelog <command>
```

#### **GitHub Releases**

Download pre-built binaries from [GitHub Releases](https://github.com/dworthen/changelog/releases).

Available platforms:

| Platform     | Architecture |
| ------------ | ------------ |
| Linux        | x64, arm64   |
| Linux (musl) | x64, arm64   |
| macOS        | x64, arm64   |
| Windows      | x64          |

<!-- tabs:end -->

## Initialize

Run `changelog init` in your project root to set up the `.changelog` directory:

```bash
changelog init
```

You'll be prompted for:

```
? Current version: 0.1.0
? Changelog file path: CHANGELOG.md
? Main git branch: main
Enter post-apply commands (leave blank to finish):
? Post-apply command: npm install
? Post-apply command (1 added):
✓ Initialized .changelog directory
```

This creates the following structure:

```
.changelog/
├── config.yaml
├── next/
│   └── .gitkeep
├── releases/
│   └── .gitkeep
└── templates/
    ├── header.md
    ├── body.eta
    └── footer.md
```

## Add an Entry

When you make a notable change, run `changelog add`:

```bash
changelog add
```

Select a change type and provide a description:

```
? Change type: Add - Add a new feature. Minor version bump.
? Description: Support for YAML configuration files
✓ Added changelog entry: .changelog/next/1740000000000.yaml
✓ Updated CHANGELOG.md
```

A YAML file is created in `.changelog/next/`:

```yaml
timestamp: 1740000000000
type: Add
description: Support for YAML configuration files
```

The changelog file is automatically regenerated with an **Unreleased** section so you can preview changes at any time.

## Apply a Release

When you're ready to release, run `changelog apply`:

```bash
changelog apply
```

This computes the version bump, creates the release, and updates your project:

```
✓ Updated version in package.json to 0.2.0
✓ Applied version 0.2.0
✓ Updated CHANGELOG.md
Running post-apply command: npm install
? Commit changes and tag the commit? Yes
✓ Committed and tagged version 0.2.0
```

The version bump is determined by the change types present:

| Change Types   | Bump  |
| -------------- | ----- |
| Change, Remove | Major |
| Add, Deprecate | Minor |
| Fix            | Patch |
| Internal       | None  |

After applying, push the commit and tag:

```bash
git push origin main --follow-tags
```

## Next Steps

- [How it Works](/how-it-works) — Understand the philosophy and mechanics behind changelog.
- [Git Flow](/git-flow) — Recommended workflow for teams using branches and pull requests.
- [Configuration](/configuration) — Customize version files, templates, and post-apply commands.
- [Managing JS Projects](/tutorials/js) — Step-by-step tutorial for Node.js projects.
- [Managing Python Projects](/tutorials/python) — Step-by-step tutorial for Python projects with uv.
