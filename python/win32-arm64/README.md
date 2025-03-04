# Changelog

Simple changelog management for projects using [semantic versioning](https://semver.org/) and git.

## Example

### Add

Add a changelog entry describing a set of changes.

```bash
changelog add
```

![changelog add](docs/images/changelog-add.gif)

This creates a timestamped file in the `.changelog` directory ready to be added to source control and applied at a later time. Depending on how changelog is configured, the `add` command also adds the timestamped file to the git staging area and then commit all staged files using the provided description as the commit message.

> [!NOTE]
> As a best practice, commit the timestamp file with the set of files the changelog entry describes. For example, run `git add <FILES...>` prior to running `changelog add`.

### Apply

The apply command gathers previously created change descriptions from the `.changelog` directory and prepends the descriptions to the `CHANGELOG.md` file. The command also bumps semantic version numbers in specified files.

```bash
changelog apply
```

![changelog apply](docs/images/changelog-apply.gif)

The resulting `CHANGELOG.md` file.

```md
## 0.2.0

### Minor Changes

- a07c8dd: Describe a minor change.
```

Running `git tag` shows that a new tag, `v0.2.0` has been applied.

## Features

- Store and track changelog descriptions per commit
- Apply changelog entries to `CHANGELOG.md`
- Sync version fields in specified JSON files such as the version field in a package.json file.

## TODO

- Add support for prereleases
- Add support for monorepos
- Add support for other version file formats (yaml, toml, etc.)

## Install

### NPM

```
npm install @d-dev/changelog -g
```

### Binaries

#### Windows

```powershell
curl -sSfL https://raw.githubusercontent.com/dworthen/changelog/main/scripts/install.ps1 | pwsh -Command -
```

This will install `changelog` to `~/bin`, be sure to export the location in the user or machine `$PATH` or use the `-to` flag to specify a download location.

Running the installer with additional flags:

```powershell
curl -sSfL https://raw.githubusercontent.com/dworthen/changelog/main/scripts/install.ps1 -o install.ps1 &&
pwsh -File install.ps1 -force -tag v0.0.1 -to ~/bin &&
rm install.ps1
```

#### Linux/Darwin

```bash
curl -sSfL https://raw.githubusercontent.com/dworthen/changelog/main/scripts/install.sh | bash
```

This will install `changelog` to `~/bin`, be sure to export the location in the user `$PATH or use the `-to` flag to specify a download location.

Running the installer with additional flags:

```bash
curl -sSfL https://raw.githubusercontent.com/dworthen/changelog/main/scripts/install.sh | bash -s -- --force --tag v0.0.1 --to ~/bin
```

### From Source with [Go](https://go.dev/)

```bash
go install github.com/dworthen/changelog@latest
```

## Getting Started

Within a project directory, run `changelog init`.

### 1. Current Version

![changelog version](docs/images/changelog-version.png)

Specifies the current version of the project. This version is tracked by `changelog` and considered the source of truth for project version. The version is bumped accordingly when `changelog` applies updates. The bumped version is then used to update specified version files. More on version files below.

### 2. Commit changelog entries

![changelog commit entries](docs/images/changelog-commit-entries.png)

If yes, `changelog add` creates a changelog entry, adds the file to the git staging area and then commits all staged files using the changelog description as the commit message. If using this option, be sure to run `git add <FILES...>` prior to running `changelog add` to associate the new changelog entry with the files it describes.

If no, then `changelog add` creates a changelog entry but does not commit the files or run any git commands.

### 3. Commit changelog apply

![changelog commit apply](docs/images/changelog-commit-apply.png)

`changelog apply` updates `CHANGELOG.md` and associated version files. If yes, then `changelog apply` will commit the changed files using the commit description `chore: apply changelog`.

### 4. Tag commit (not presented if previous answer == No)

![changelog tag](docs/images/changelog-tag.png)

If yes, `changelog apply` tags the commit created. This option is ignored if the answer to the previous question is No.

### 5. Tag format (not presented if either of the previous two answers == No)

![changelog tag format](docs/images/changelog-tag-format.png)

Specify the tag to apply. This string is evaluated using handlebars and has access to the new version. Defaults to `v{{version}}`. This option is ignored if either of the previous two answers is No.

### 6. Version files

![changelog version files](docs/images/changelog-version-files.png)

Specify JSON files containing fields that should match the project version. The JSON path can point to a nested field using dot notation. After `changelog apply` bumps the current package version it updates all the specified JSON files and fields.

> [!WARNING]
> Currently, only JSON files are supported.

### Config

`changelog init` creates a `.changelog` directory and the following two files.

**config.json**

Describes the config the `init` command walks through.

```json
{
  "version": "1.0.0",
  "bumpFiles": [
    {
      "file": "package.json",
      "path": "version"
    },
    {
      "file": "some-json-file.json",
      "path": "path.to.nested"
    }
  ],
  "onAdd": {
    "commitFiles": true
  },
  "onApply": {
    "commitFiles": true,
    "tagCommit": true,
    "tagFormat": "v{{version}}"
  }
}
```

**changelogTemplate.hbs**

A handlebars file describing the template to use when adding entries to `CHANGELOG.md`. Users have full control of how the changelog entries are formatted. The following handlebars variables are available.

- **version**: The new project version after `changelog apply` finishes.
- **oldVersion**: The previous project version prior to `changelog apply`.
- **patchChanges|minorChanges|majorChanges**: A list of changes grouped by the level of change (patch, minor, or major). The variable is available if there is a change of that type present in the current update, otherwise its null.
  - **sha**: The full git commit hash for the change.
  - **shortSha**: The short git commit hash for the change.
  - **description**: Changelog description.

```md
## {{version}}
{{#if majorChanges}}

### Major Changes

{{#each majorChanges}}
- {{shortSha}}: {{description}}
{{~/each}}
{{/if}}
{{#if minorChanges}}

### Minor Changes

{{#each minorChanges}}
- {{shortSha}}: {{description}}
{{~/each}}
{{/if}}
{{#if patchChanges}}

### Patch Changes

{{#each patchChanges}}
- {{shortSha}}: {{description}}
{{/each}}
{{~/if}}
```
