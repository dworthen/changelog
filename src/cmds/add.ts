import { createCommand } from '@d-dev/roar'
import { input, select } from '@inquirer/prompts'
import {
  buildReleaseData,
  generateChangelog,
  loadNextEntries,
  loadReleases,
} from '../lib/changelog'
import { loadConfig } from '../lib/config'
import { CHANGE_TYPES, type ChangeType } from '../lib/types'

export const addCommand = createCommand(
  {
    usageName: 'changelog add',
    description: 'Add a new changelog entry',
  },
  async () => {
    const config = await loadConfig()

    const type = await select<ChangeType>({
      message: 'Change type:',
      choices: Object.entries(CHANGE_TYPES).map(([name, desc]) => ({
        name: `${name} - ${desc}`,
        value: name as ChangeType,
      })),
    })

    const description = await input({
      message: 'Description:',
      validate: (value) => value.length > 0 || 'Required',
    })

    const timestamp = Date.now()
    const filename = `${timestamp}.yaml`
    const entryYaml = Bun.YAML.stringify(
      { timestamp, type, description },
      null,
      2,
    )

    await Bun.write(`.changelog/next/${filename}`, entryYaml)
    console.log(`✓ Added changelog entry: .changelog/next/${filename}`)

    // Regenerate changelog with "Unreleased" version
    const entries = await loadNextEntries()
    const unreleased = await buildReleaseData('Unreleased', entries)
    const releases = await loadReleases()
    const allReleases = [unreleased, ...releases]
    const changelog = await generateChangelog(allReleases)
    await Bun.write(config.changelogFile, changelog)

    console.log(`✓ Updated ${config.changelogFile}`)
  },
)