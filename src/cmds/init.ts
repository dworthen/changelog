import { mkdirSync } from 'node:fs'
import { createCommand } from '@d-dev/roar'
import { input } from '@inquirer/prompts'
import { saveConfig } from '../lib/config'
import { type Config } from '../lib/types'
// @ts-expect-error
import bodyContent from '../templates/body.eta' with { type: 'text' }
// @ts-expect-error
import footerContent from '../templates/footer.md' with { type: 'text' }
// @ts-expect-error
import headerContent from '../templates/header.md' with { type: 'text' }

export const initCommand = createCommand(
  {
    usageName: 'changelog init',
    description: 'Initialize a new .changelog directory',
  },
  async () => {
    const currentVersion = await input({
      message: 'Current version:',
      validate: (value) =>
        /^\d+\.\d+\.\d+$/.test(value) || 'Must be a valid semver (e.g., 0.0.0)',
    })

    const changelogFile = await input({
      message: 'Changelog file path:',
      default: 'CHANGELOG.md',
      validate: (value) => value.length > 0 || 'Required',
    })

    const mainBranch = await input({
      message: 'Main git branch:',
      default: 'main',
      validate: (value) => value.length > 0 || 'Required',
    })

    const postCommands: string[] = []
    console.log('Enter post-apply commands (leave blank to finish):')
    while (true) {
      const cmd = await input({
        message: `Post-apply command${postCommands.length > 0 ? ` (${postCommands.length} added)` : ''}:`,
      })
      if (cmd.trim() === '') break
      postCommands.push(cmd.trim())
    }

    const config: Config = {
      currentVersion,
      changelogFile,
      git: {
        mainBranch,
      },
      versionFiles: [
        {
          path: 'package.json',
          pattern: ':\\s*"(\\d+\\.\\d+\\.\\d+)"',
        },
        {
          path: 'pyproject.toml',
          pattern: 'version\\s*=\\s*"(\\d+\\.\\d+\\.\\d+)"',
        },
      ],
      apply: {
        ...(postCommands.length > 0 ? { postCommands } : {}),
      },
    }

    // Create directory structure
    mkdirSync('.changelog/templates', { recursive: true })
    mkdirSync('.changelog/next', { recursive: true })
    mkdirSync('.changelog/releases', { recursive: true })
    await Bun.write('.changelog/next/.gitkeep', '')
    await Bun.write('.changelog/releases/.gitkeep', '')

    await Bun.write('.changelog/templates/header.md', headerContent)

    await Bun.write('.changelog/templates/body.eta', bodyContent)

    await Bun.write('.changelog/templates/footer.md', footerContent)

    // Write config
    await saveConfig(config)

    console.log('✓ Initialized .changelog directory')
  },
)