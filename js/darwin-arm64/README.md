# Changelog

Git-based changelog manager for JavaScript, Python, and Go projects using [semantic versioning](https://semver.org/) and git.

## Usage

```shell
$ changelog help

Git-based changelog manager for JavaScript, Python, and Go projects.

Usage:
  changelog [command]

Available Commands:
  add         Add a changelog entry.
  apply       Apply changelog entries.
  check       Check that the current branch or last commit contains a changelog entry. Useful for CI workflows to enforce the presence of changelog entries.
  help        Help about any command
  init        Initialize project to use changelog.
  update      Update to the latest version of changelog
  version     Print the current version of changelog

Flags:
      --cwd string   project directory for changelog. (default ".")
  -h, --help         help for changelog
      --verbose      enable verbose logging.

Use "changelog [command] --help" for more information about a command.
```

## Init

```shell
changelog init
```

Walks through creating a `.changelog` directory with the following files.

### `config.yaml`

Changelog configuration:

```yaml
# Example config.yaml

# Current version of the project.
# This gets bumped, along with the files listed below,
# when running `changelog apply`
version: 0.1.0
changelogFile: CHANGELOG.md # Changelog file to manage

# A list of files containing versions to bump when running
# `changelog apply`.
# Uses RegExp patterns for matching and replacing the actual version.
# The RegExp must contain 1 capture group with the version portion to replace.
# View https://pkg.go.dev/regexp/syntax@go1.24.1 for RegExp syntax
# Use https://regex101.com/ with GoLang selected to test patterns.
# Examples for Node package.json and python pyproject.toml below.
files:
    - path: package.json
      pattern: '"version":\s*"(\d+\.\d+\.\d+)"'
    - path: pyproject.toml
      pattern: 'version\s*=\s*"(\d+\.\d+\.\d+)"'

# Configures `changelog add` command
onAdd:
    # commit staged files + added changelog entry
    # Uses the provided description as the commit message.
    commitFiles: true

# Configure `changelog apply` command
onApply:
    commitFiles: true
    tagCommit: true
    tagFormat: v{{version}}
    # A list of commands to run after bumping version files
    # and before committing and tagging.
    # Often useful to run install/sync commands that may update
    # lock files.
    commands:
        - npm install
        - uv sync

```

### `changelogTemplate.hbs`

A [Handlebars](https://handlebarsjs.com/) template used for adding new entries to the changelog file specified in the `.changelog/config.yaml` file.

Default template:

```
## {{version}}
{{#if majorChanges}}

### Major Changes

{{#each majorChanges}}
- {{shortSha}}: {{description}}
{{/each}}
{{/if}}
{{#if minorChanges}}

### Minor Changes

{{#each minorChanges}}
- {{shortSha}}: {{description}}
{{/each}}
{{/if}}
{{#if patchChanges}}

### Patch Changes

{{#each patchChanges}}
- {{shortSha}}: {{description}}
{{/each}}
{{/if}}


```

**Available Variables**

- version: The new version of the project.
- oldVersion: The previous version of the project.
- majorChanges, minorChanges, patchChanges: A list of change descriptions.
  - **Change Description:**
    - Sha: The git Sha associated with the change commit.
    - shortSha: The git short Sha associated with the change commit.
    - change: `patch|minor|major`
    - description: The description of the change

## Add

```shell
changelog add
```

![changelog add](docs/images/changelog-add.gif)

Creates a timestamped `.md` changelog entry file in the `.changelog` directory ready to be added to source control and applied at a later time. The changelog entry contains the type of change (`patch|minor|major`) along with the provided description.

Example changelog entry:

```
---
change: "major"
---
A major change.
```

Depending on how changelog is configured, the `add` command also adds the timestamped file to the git staging area and then commit all staged files using the provided description as the commit message.

> [!NOTE]
> As a best practice, commit the timestamp file with the set of files the changelog entry describes. For example, run `git add <FILES...>` prior to running `changelog add`.

## Apply

The apply command...

- Runs in dry mode by default. The command shows a preview of the changes and ask for confirmation prior to applying any changes.
- gathers previously created changelog entries from the `.changelog` directory and prepends the info to the `CHANGELOG.md` file according to the `.changelog/changelogTemplate.hbs` template file.
- Bumps semantic version numbers in specified files according to `.changelog/config.yaml`.
- run any commands listed in `.changelog/config.yaml`.

```bash
changelog apply
```

![changelog apply](docs/images/changelog-apply.gif)

If `.changelog/config.yaml` is configured to commit and tag files when running `changelog apply` then `git log` and `git tag` should show the commit and associated tag.

## TODO

- Add support for prereleases

## Install

### NPM

https://www.npmjs.com/package/@d-dev/changelog

```shell
npm install @d-dev/changelog
# or globally
npm install @d-dev/changelog -g
```

### Python

https://pypi.org/project/changesets/

```shell
pip install changesets
# Or with UV
uv add changesets --dev
# package is called changesets but command is still changelog
```

### Golang

https://pkg.go.dev/github.com/dworthen/changelog

With go version `1.24.x` you can add the package as a tool to your project with

```shell
go get -tool github.com/dworthen/changelog@latest
```

and use with

```shell
go tool changelog help
```

Prior to 1.24

```
go install github.com/dworthen/changelog
```

### Binaries

#### Windows

Requires the newer [powershell core](https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows?view=powershell-7.5).

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

This will install `changelog` to `~/bin`, be sure to export the location in the user `$PATH` or use the `-to` flag to specify a download location.

Running the installer with additional flags:

```bash
curl -sSfL https://raw.githubusercontent.com/dworthen/changelog/main/scripts/install.sh | bash -s -- --force --tag v0.0.1 --to ~/bin
```

### Prebuilt Binaries

https://github.com/dworthen/changelog/releases

### From Source with [Go](https://go.dev/)

```bash
git clone https://github.com/dworthen/changelog.git
cd changelog
go mod tidy
go build ./main.go
```
