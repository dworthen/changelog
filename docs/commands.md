# Commands

Reference for all CLI commands provided by `changelog`.

## `changelog init`

Interactive setup wizard that initializes a `.changelog` directory in your project.

**Usage:**

```bash
changelog init
```

**What it does:**

Creates the following directory structure:

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

**Interactive prompts:**

| Prompt              | Description                                              | Default        |
| ------------------- | -------------------------------------------------------- | -------------- |
| Current version     | The current semver version of your project               | —              |
| Changelog file path | Path to the generated changelog file                     | `CHANGELOG.md` |
| Main git branch     | The main/default branch name                             | `origin/main`  |
| Post-apply commands | Shell commands to run after `apply` (repeat until blank) | —              |

**Example:**

```bash
$ changelog init
? Current version: 0.1.0
? Changelog file path: CHANGELOG.md
? Main git branch: main
Enter post-apply commands (leave blank to finish):
? Post-apply command: npm install
? Post-apply command (1 added):
✓ Initialized .changelog directory
```

## `changelog add`

Interactively add a new changelog entry.

**Usage:**

```bash
changelog add
```

**What it does:**

1. Prompts for a change type and description.
2. Writes a timestamped YAML file to `.changelog/next/`.
3. Regenerates the changelog file with an "Unreleased" section.

**Interactive prompts:**

| Prompt      | Description                                           |
| ----------- | ----------------------------------------------------- |
| Change type | One of: Add, Change, Deprecate, Remove, Fix, Internal |
| Description | A short description of the change                     |

**Change types:**

| Type      | Description                                       | Version Bump |
| --------- | ------------------------------------------------- | ------------ |
| Change    | Change an existing feature                        | Major        |
| Remove    | Remove a feature                                  | Major        |
| Add       | Add a new feature                                 | Minor        |
| Deprecate | Deprecate a feature                               | Minor        |
| Fix       | Fix a bug                                         | Patch        |
| Internal  | Internal change, not included in public changelog | None         |

**Example:**

```bash
$ changelog add
? Change type: Add - Add a new feature. Minor version bump.
? Description: Support for YAML configuration files
✓ Added changelog entry: .changelog/next/1740000000000.yaml
✓ Updated CHANGELOG.md
```

The created entry file (`.changelog/next/1740000000000.yaml`):

```yaml
timestamp: 1740000000000
type: Add
description: Support for YAML configuration files
```

## `changelog apply`

Apply pending changelog entries and bump the project version.

**Usage:**

```bash
changelog apply
```

**What it does:**

1. Reads all pending entries from `.changelog/next/`.
2. Computes the new semver version based on the change types present.
3. Creates a versioned release file in `.changelog/releases/`.
4. Removes processed entry files from `.changelog/next/`.
5. Updates the version string in all configured `versionFiles`.
6. Updates `currentVersion` in `.changelog/config.yaml`.
7. Regenerates the changelog file from all releases.
8. Runs any configured `postCommands`.
9. Prompts to commit changes and create a git tag.

**Version bump rules:**

- If any entry is `Change` or `Remove` → **major** bump.
- Else if any entry is `Add` or `Deprecate` → **minor** bump.
- Else if any entry is `Fix` → **patch** bump.
- `Internal` entries do not trigger a version bump.

**Example:**

```bash
$ changelog apply
✓ Updated version in package.json to 1.1.0
✓ Applied version 1.1.0
✓ Updated CHANGELOG.md
Running post-apply command: npm install
? Commit changes and tag the commit? Yes
✓ Committed and tagged version 1.1.0
```

## `changelog view`

View the changelog or a specific release.

**Usage:**

```bash
changelog view [argument]
```

**Arguments:**

| Argument    | Description                                          |
| ----------- | ---------------------------------------------------- |
| _(none)_    | Full changelog preview, including unreleased entries |
| `latest`    | Render the latest release only                       |
| `next`      | Render unreleased entries only                       |
| `<version>` | Render a specific version (e.g., `1.2.0`)            |

**Examples:**

```bash
# Full changelog with unreleased entries
changelog view

# Latest release only
changelog view latest

# Unreleased entries
changelog view next

# Specific version
changelog view 1.2.0
```

## `changelog check`

CI/CD gate that verifies the current branch includes at least one changelog entry.

**Usage:**

```bash
changelog check
```

**What it does:**

1. Loads the `git.mainBranch` from `.changelog/config.yaml`.
2. Finds the merge base between the main branch and `HEAD` using `git merge-base`.
3. Checks `git diff` for any files added or changed in `.changelog/next/` since the fork point.
4. Exits with code `0` if entries are found, or code `1` with an error message if none are found.

**Example (passing):**

```bash
$ changelog check
Check passed: found 2 changelog entries in .changelog/next/
```

**Example (failing):**

```bash
$ changelog check
Error: No changelog entries found in .changelog/next/ for this branch. Please add a changelog entry before merging.
```

!> The `check` command requires full git history to compute the merge base. In CI, make sure to use `fetch-depth: 0` when checking out the repository.

## `changelog version`

Print the current project version.

**Usage:**

```bash
changelog version
```

**What it does:**

Reads `currentVersion` from `.changelog/config.yaml` and prints it to stdout.

**Example:**

```bash
$ changelog version
1.2.0
```
