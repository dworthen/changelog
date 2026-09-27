import { createCommand } from '@d-dev/roar'
import { input, select } from '@inquirer/prompts'
import {
  buildReleaseData,
  generateChangelog,
  loadNextEntries,
  loadReleases,
} from '../lib/changelog'
import { parseChangeType } from '../lib/changeType'
import { loadConfig } from '../lib/config'
import { CHANGE_TYPES, type ChangeType } from '../lib/types'

export const addCommand = createCommand(
  {
    usageName: 'changelog add',
    description: 'Add a new changelog entry',
    flags: {
      type: {
        type: 'string',
        shortFlag: 't',
        description: 'Change type. Skips prompts when used with --message.',
      },
      message: {
        type: 'string',
        shortFlag: 'm',
        description: 'Entry description. Skips prompts when used with --type.',
      },
    },
  },
  async ({ flags }) => {
    const config = await loadConfig()

    let type: ChangeType
    let description: string

    if (flags.type != null || flags.message != null) {
      if (flags.type == null || flags.message == null) {
        console.error('Error: --type and --message must be used together')
        process.exit(1)
      }

      const parsedType = parseChangeType(flags.type)
      if (!parsedType.ok) {
        console.error(`Error: ${parsedType.error}`)
        process.exit(1)
      }

      if (flags.message.trim().length === 0) {
        console.error('Error: --message must not be empty')
        process.exit(1)
      }

      type = parsedType.value
      description = flags.message
    } else {
      type = await select<ChangeType>({
        message: 'Change type:',
        choices: Object.entries(CHANGE_TYPES).map(([name, desc]) => ({
          name: `${name} - ${desc}`,
          value: name as ChangeType,
        })),
      })

      description = await input({
        message: 'Description:',
        validate: (value) => value.length > 0 || 'Required',
      })
    }

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