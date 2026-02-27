# How it Works

## Philosophy

Traditional changelog management asks developers to edit a single shared file — typically `CHANGELOG.md`. This approach creates frequent merge conflicts when multiple branches are in flight, and it's easy to forget to update the changelog entirely.

`changelog` takes a different approach: **each changelog entry is its own file**. Instead of editing a shared document, developers run `changelog add` to create a small YAML file describing their change. These individual files live in `.changelog/next/` and are merged alongside the code changes that produce them. Since each entry is a separate file, **merge conflicts are virtually eliminated**.

When it's time to cut a release, `changelog apply` collects all pending entries, computes a version bump, and produces the final changelog — automatically.

## Directory Structure

After running `changelog init`, your project contains a `.changelog/` directory:

```
.changelog/
├── config.yaml          # Project configuration
├── next/                # Pending changelog entries
│   └── .gitkeep
├── releases/            # Versioned release data
│   └── .gitkeep
└── templates/           # Changelog rendering templates
    ├── header.md
    ├── body.eta
    └── footer.md
```

| Directory/File | Purpose                                                                                                                        |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `config.yaml`  | Stores the current version, git settings, version file patterns, and post-apply commands. See [Configuration](/configuration). |
| `next/`        | Holds pending changelog entry files (one per change). Emptied on `apply`.                                                      |
| `releases/`    | Stores one YAML file per released version, containing all changes for that release.                                            |
| `templates/`   | Controls how the changelog file is rendered.                                                                                   |

## The `add` Workflow

Running [`changelog add`](/commands?id=changelog-add) creates a timestamped YAML file in `.changelog/next/`:

```yaml
timestamp: 1740000000000
type: Add
description: Support for YAML configuration files
```

Each file captures:

- **timestamp** — When the entry was created (milliseconds since epoch).
- **type** — The category of change (`Add`, `Change`, `Deprecate`, `Remove`, `Fix`, or `Internal`).
- **description** — A human-readable summary of the change.

After writing the entry, the changelog file is regenerated with an **"Unreleased"** section so that your working changelog is always up to date — even before a release is cut.

## The `apply` Workflow

Running [`changelog apply`](/commands?id=changelog-apply) performs the release process:

1. **Collect entries** — All `.yaml` files in `.changelog/next/` are read.
2. **Compute version** — The new semver version is determined from the change types present. `Change`/`Remove` trigger a major bump, `Add`/`Deprecate` a minor bump, and `Fix` a patch bump. The highest-priority bump wins.
3. **Create release file** — A YAML file is written to `.changelog/releases/<version>.yaml` containing the release timestamp, version, change types, and all change descriptions with their git SHAs.
4. **Clean up** — All `.yaml` files in `.changelog/next/` are deleted.
5. **Update version files** — The version string is updated in all files listed in `config.versionFiles` (e.g., `package.json`, `pyproject.toml`).
6. **Update config** — `currentVersion` in `config.yaml` is set to the new version.
7. **Regenerate changelog** — The full changelog is rebuilt from all release files.
8. **Run post-commands** — Any commands in `config.apply.postCommands` are executed (e.g., `npm install` to update `package-lock.json`).
9. **Commit and tag** — You're prompted to commit all changes and create a `v<version>` annotated git tag.

## Changelog Generation

The changelog file is **generated on the fly** from the data in `.changelog/releases/`. It is never the source of truth — the release YAML files are. This means you can regenerate the changelog at any time without data loss.

The generation process concatenates three parts:

1. **Header** (`templates/header.md`) — Static content at the top of the file.
2. **Body** (`templates/body.eta`) — An [Eta](https://eta.js.org/) template rendered once per release, in reverse chronological order (newest first).
3. **Footer** (`templates/footer.md`) — Static content at the bottom. Useful for preserving historical changelog content when [migrating](/migrating).

Because the changelog is derived from structured data, you can customize the output format by editing the body template without touching historical entries.

## Internal Changes

The `Internal` change type serves a special purpose. Entries marked as `Internal`:

- **Are tracked** in `.changelog/next/` so they satisfy the [`changelog check`](/commands?id=changelog-check) CI gate.
- **Are filtered out** of public releases — they do not appear in the generated changelog.
- **Do not trigger a version bump** on their own.

This is useful for changes that don't affect end users — refactors, CI configuration updates, documentation improvements, or dependency bumps. Developers can still pass the CI changelog check without adding noise to the public changelog.
