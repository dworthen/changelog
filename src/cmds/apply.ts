import { readdirSync, unlinkSync } from 'node:fs'
import { join } from 'node:path'
import { createCommand } from '@d-dev/roar'
import { confirm } from '@inquirer/prompts'
import { $ } from 'bun'
import {
  buildReleaseData,
  generateChangelog,
  loadNextEntries,
  loadReleases,
} from '../lib/changelog'
import { loadConfig, saveConfig } from '../lib/config'
import { bumpVersion } from '../lib/semver'
import { viewHandler } from './view'

export const applyCommand = createCommand(
  {
    usageName: 'changelog apply',
    description: 'Apply pending changelog entries and bump version',
  },
  async () => {
    const config = await loadConfig()

    const entries = await loadNextEntries()
    if (entries.length === 0) {
      console.warn('No pending changelog entries found.')
      return
    }

    // Collect unique change types
    const changeTypes = [...new Set(entries.map((e) => e.type))]

    // Compute new version
    const newVersion = bumpVersion(config.currentVersion, changeTypes)

    // Build release data
    const releaseData = await buildReleaseData(newVersion, entries)

    // Write release file
    const releaseYaml = Bun.YAML.stringify(
      {
        timestamp: releaseData.timestamp,
        version: releaseData.version,
        changeTypes: releaseData.changeTypes,
        changes: releaseData.changes,
      },
      null,
      2,
    )
    await Bun.write(`.changelog/releases/${newVersion}.yaml`, releaseYaml)

    // Delete all .yaml files in .changelog/next
    const nextDir = '.changelog/next'
    const nextFiles = readdirSync(nextDir).filter((f) => f.endsWith('.yaml'))
    for (const file of nextFiles) {
      unlinkSync(join(nextDir, file))
    }

    // Update version in versionFiles
    for (const vf of config.versionFiles) {
      const file = Bun.file(vf.path)
      if (await file.exists()) {
        const content = await file.text()
        const regex = new RegExp(vf.pattern, 'g')
        const updated = content.replaceAll(regex, (full, captured) => {
          return full.replace(captured, newVersion)
        })
        if (updated !== content) {
          await Bun.write(vf.path, updated)
          console.log(`✓ Updated version in ${vf.path} to ${newVersion}`)
        }
      }
    }

    // Update config with new version
    config.currentVersion = newVersion
    await saveConfig(config)

    // Regenerate changelog
    const releases = await loadReleases()
    const changelog = await generateChangelog(releases)
    await Bun.write(config.changelogFile, changelog)

    console.log(`✓ Applied version ${newVersion}`)
    console.log(`✓ Updated ${config.changelogFile}`)

    if (config.apply.postCommands && config.apply.postCommands.length > 0) {
      for (const cmd of config.apply.postCommands) {
        console.log(`Running post-apply command: ${cmd}`)
        try {
          await $`${{ raw: cmd }}`
        } catch (err) {
          console.error(`Warning: Post-apply command failed: ${err}`)
        }
      }
    }

    const changelogEntry = await viewHandler('latest')

    const shouldCommit = await confirm({
      message: 'Commit changes and tag the commit?',
      default: true,
    })

    if (shouldCommit) {
      try {
        await $`git add .`
        await $`git commit -m "${changelogEntry}"`
        await $`git tag -a v${newVersion} -m "v${newVersion}"`
        console.log(`✓ Committed and tagged version ${newVersion}`)
      } catch (err) {
        console.error(`Error: Failed to commit and tag: ${err}`)
      }
    }
  },
)