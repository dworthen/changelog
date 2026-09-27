# AGENTS.md

## Project Overview

`changelog` (`@d-dev/changelog`) is a cross-platform command-line tool for managing
changelogs. Individual changelog entries are stored as files under `.changelog/`,
eliminating merge conflicts and enabling CI enforcement. It is built with Bun and
TypeScript and ships as self-contained, compiled binaries for Linux, macOS, and Windows
(x64 and arm64, including musl), and is published to npm (`@d-dev/changelog`), PyPI
(`changesets`), and GitHub Releases. The published repository is `dworthen/changelog`.

## Languages and Tooling

- **Language:** TypeScript (ESM, `strict` mode, targets ESNext)
- **Runtime:** Bun (configured via `bunfig.toml`)
- **CLI framework:** `@d-dev/roar` (`createCommand`)
- **Interactive prompts:** `@inquirer/prompts`
- **Templating:** `eta` (see `src/lib/eta.ts` and `src/templates/`)
- **Lint/format:** Biome (`@biomejs/biome`)
- **Release tooling:** `@d-dev/bin-upload` (binary packaging + npm/GitHub release
  publishing). This project manages its own changelog with itself — add an entry for
  user-facing changes with `bun run changelog:add`.

Prefer using Bun APIs over Node APIs when possible (https://bun.com/llms.txt). Here are some examples:

- Use Bun for file IO https://bun.com/docs/runtime/file-io
- Use Bun for globbing https://bun.com/docs/runtime/glob
- Use Bun for hashing https://bun.com/docs/runtime/hashing
- Use Bun for working with yaml https://bun.com/docs/runtime/yaml and serialize with `Bun.YAML.stringify(object, null, 2)`
- Use Bun server for web server https://bun.com/docs/runtime/http/server
- Use Bun Shell for running sub process commands https://bun.sh/docs/runtime/shell

## Setup Commands

- Install dependencies: `bun install`
- Run locally: `bun run src/index.ts <command>` (e.g. `bun run src/index.ts add`)
- Build all standalone binaries: `bun run build`
- Build a single target: `bun run scripts/build.ts <target>` (e.g. `bun-darwin-arm64`)
- Package + publish binaries (used in CI): `bun run pack`, then `bun run publish`
- Serve docs locally: `bun run docs`

## Code Styles and Linting Commands

- Formatter/linter: Biome (`biome.json`)
- Check formatting + lint: `bun run check`
- Auto-fix: `bun run fix`
- Conventions: single quotes, trailing commas everywhere, semicolons only as needed,
  imports auto-organized, inline `type` imports
- Honors `.editorconfig` and `.gitignore`; `noExplicitAny` and `noNonNullAssertion` are disabled

## Testing Instructions and Commands

- Use Bun's built-in test runner (https://bun.sh/docs/test): add `*.test.ts` files and run `bun test`.
- Focus on writing unit tests and using mocks. Do not test the CLI commands.
- CI quality gate (`.github/workflows/pr.yml`) runs on every PR: `bun install`,
  `bun run check:changelog`, and `bun run check`.

## Architecture Patterns

- Prefer functions over classes.
- Functions should implement types for composability.
- Required function parameters should be listed out while optional parameters are grouped into an `options` object.
- Here is an example of a proper function

```typescript
type SomeFunctionOptions = {
  someParam?: bool
}

type SomeFunctionResult =
  | {
      ok: false
      errors: string[]
    }
  | {
      ok: true
      result: string
    }

type SomeFunction = (arg1: string, arg2: bool, options?: SomeFunctionOptions) => Promise<SomeFunctionResult>

const someFunction: SomeFunction = (arg1, arg2, { someParam = false } = {}) {
  ...
}
```

## Considerations

- **Entry point:** `src/index.ts` builds the root CLI with `@d-dev/roar`, registers
  subcommands via `changelogCommand.addCommand(...)`, and centralizes error handling
  (graceful Ctrl+C for `@inquirer/prompts`).
- **Commands:** one file per command in `src/cmds/` (`init.ts`, `add.ts`, `apply.ts`,
  `check.ts`, `view.ts`, `version.ts`), each exporting a `createCommand({...}, handler)`.
- Commands should orchestrate testable functions/logic from `src/lib/`.
- **Domain types and logic:** Keep domain types (for example, `Config`, `ChangeType`,
  `ReleaseData`) in `src/lib/types.ts`, and their related loading, parsing, validation,
  mutation, and persistence functions in dedicated `src/lib/` modules (`config.ts`,
  `changelog.ts`, `git.ts`, `semver.ts`, `eta.ts`). Commands should import and orchestrate
  those APIs rather than declaring domain types or implementing persistence inline.
- **Templates:** changelog output is rendered from `eta`/Markdown templates in
  `src/templates/` (`body.eta`, `header.md`, `footer.md`); a project's own copies live
  under `.changelog/templates/`.
- **Data model:** unreleased entries live as YAML files in `.changelog/next/`, released
  entries in `.changelog/releases/`, and project settings in `.changelog/config.yaml`.
- **Build:** `scripts/build.ts` cross-compiles standalone executables for 8 targets
  (Linux/macOS/Windows × x64/arm64, including musl) into `bin/` via `Bun.build({ compile })`.
- **Distribution:** `bin-upload.config.yaml` archives binaries and publishes to npm
  (`@d-dev/changelog` plus per-platform packages) and GitHub Releases on `dworthen/changelog`.
- **CI:** `pr.yml` (changelog + checks) and `release.yml` (pack + publish on `v*` tags).

## Agent skills

### Issue tracker

Issues and specs are tracked as GitHub issues in `dworthen/changelog` via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Domain docs

Single-context: one `CONTEXT.md` and `docs/adr/` at the repo root. See `docs/agents/domain.md`.
