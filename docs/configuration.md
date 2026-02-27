# Configuration

The `changelog` tool is configured via `.changelog/config.yaml`, created during [`changelog init`](/commands?id=changelog-init).

## Structure

A complete example configuration:

```yaml
currentVersion: 1.0.0
changelogFile: CHANGELOG.md
git:
  mainBranch: main
versionFiles:
  - path: package.json
    pattern: ':\s*"(\d+\.\d+\.\d+)"'
  - path: pyproject.toml
    pattern: 'version\s*=\s*"(\d+\.\d+\.\d+)"'
apply:
  postCommands:
    - npm install
```

## Property Reference

### Top-level Properties

| Property         | Type     | Required | Default        | Description                                               |
| ---------------- | -------- | -------- | -------------- | --------------------------------------------------------- |
| `currentVersion` | `string` | Yes      | —              | Current semver version of the project (e.g., `1.0.0`)     |
| `changelogFile`  | `string` | Yes      | `CHANGELOG.md` | Path to the generated changelog file                      |
| `git`            | `object` | Yes      | —              | Git-related configuration                                 |
| `versionFiles`   | `array`  | No       | `[]`           | Files to update with the new version when running `apply` |
| `apply`          | `object` | No       | `{}`           | Configuration for the `apply` command                     |

### `git`

| Property     | Type     | Required | Default       | Description                                                                                                                |
| ------------ | -------- | -------- | ------------- | -------------------------------------------------------------------------------------------------------------------------- |
| `mainBranch` | `string` | Yes      | `origin/main` | The main/default git branch name. Used by the [`check`](/commands?id=changelog-check) command to determine the merge base. |

### `versionFiles`

An array of objects. Each entry defines a file and a regex pattern used to find and replace the version string during `apply`.

| Property  | Type     | Required | Description                                                                       |
| --------- | -------- | -------- | --------------------------------------------------------------------------------- |
| `path`    | `string` | Yes      | File path relative to the project root                                            |
| `pattern` | `string` | Yes      | Regex pattern with a capture group `(...)` matching the version string to replace |

The regex pattern must contain exactly one capture group that matches the version portion. During `apply`, the captured group is replaced with the new version while the rest of the match is preserved.

**Examples:**

```yaml
versionFiles:
  # Matches "version": "1.0.0" in JSON files
  - path: package.json
    pattern: ':\s*"(\d+\.\d+\.\d+)"'

  # Matches version = "1.0.0" in TOML files
  - path: pyproject.toml
    pattern: 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

### `apply`

| Property       | Type       | Required | Default | Description                                                                                                        |
| -------------- | ---------- | -------- | ------- | ------------------------------------------------------------------------------------------------------------------ |
| `postCommands` | `string[]` | No       | `[]`    | Shell commands to run sequentially after `apply` completes. Useful for updating lock files or running build steps. |

**Example:**

```yaml
apply:
  postCommands:
    - npm install
    - npm run build
```

## Templates

The `.changelog/templates/` directory contains three files that control how the changelog is rendered. These are created during `init` and can be customized.

### `header.md`

Static Markdown prepended to the top of the changelog. By default:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
```

### `body.eta`

An [Eta](https://eta.js.org/) template rendered once per release. The following variables are available:

| Variable      | Type       | Description                                                                                                                 |
| ------------- | ---------- | --------------------------------------------------------------------------------------------------------------------------- |
| `version`     | `string`   | The release version (e.g., `1.2.0`) or `Unreleased`                                                                         |
| `date`        | `string`   | Formatted date (`YYYY-MM-DD`)                                                                                               |
| `changeTypes` | `string[]` | Ordered list of change type names present in this release                                                                   |
| `changes`     | `object`   | Changes grouped by type. Keys are change type names, values are arrays of change objects with `shortSha` and `description`. |

Default template:

```
## <%= version %> - <%= date %>
<% changeTypes.forEach(changeType => { %>
### <%= changeType %>
<% changes[changeType].forEach(change => { %>
<%= change.shortSha %>: <%= change.description %>
<%- }) %>
<% }) %>
```

### `footer.md`

Static Markdown appended to the end of the changelog. Empty by default. Useful for preserving historical changelog content when [migrating](/migrating).

## Change Types

Each changelog entry has a type that determines the semver version bump when `apply` is run.

| Type      | Description                                        | Version Bump |
| --------- | -------------------------------------------------- | ------------ |
| Change    | Change an existing feature                         | Major        |
| Remove    | Remove a feature                                   | Major        |
| Add       | Add a new feature                                  | Minor        |
| Deprecate | Deprecate a feature                                | Minor        |
| Fix       | Fix a bug                                          | Patch        |
| Internal  | Internal change (not included in public changelog) | None         |

When multiple entry types are present, the highest-priority bump wins:

1. **Major** — if any `Change` or `Remove` entry exists
2. **Minor** — if any `Add` or `Deprecate` entry exists
3. **Patch** — if any `Fix` entry exists
4. **None** — if only `Internal` entries exist

?> `Internal` entries are tracked in `.changelog/next/` so they satisfy the [`changelog check`](/commands?id=changelog-check) CI gate, but they are filtered out of public releases.
