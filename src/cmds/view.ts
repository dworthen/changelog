import { readdirSync } from 'node:fs'
import { basename, join } from 'node:path'
import { createCommand } from '@d-dev/roar'
import {
  buildReleaseData,
  generateChangelog,
  loadBodyTemplate,
  loadNextEntries,
  loadReleases,
  renderRelease,
} from '../lib/changelog'
import { type ReleaseData } from '../lib/types'

export async function viewHandler(argument?: string): Promise<string | null> {
  if (!argument) {
    // Case 1: No argument — full changelog preview
    const releases = await loadReleases()
    const entries = await loadNextEntries()

    let allReleases = [...releases]
    if (entries.length > 0) {
      const unreleased = await buildReleaseData('Unreleased', entries)
      allReleases = [unreleased, ...releases]
    }

    const changelog = await generateChangelog(allReleases)
    return changelog
  }

  if (argument === 'latest') {
    // Case 2: "latest" — render latest release
    const releases = await loadReleases()
    if (releases.length === 0) {
      console.error('Error: No releases found in .changelog/releases/')
      process.exit(1)
    }

    const latest = releases[0]!
    const bodyTemplate = await loadBodyTemplate()
    const rendered = renderRelease(bodyTemplate, latest)
    return rendered
  }

  if (argument === 'next') {
    // Case: "next" — render unreleased changelog entries
    const entries = await loadNextEntries()
    if (entries.length === 0) {
      console.log('No pending changelog entries found.')
      return null
    }

    const unreleased = await buildReleaseData('Unreleased', entries)
    const bodyTemplate = await loadBodyTemplate()
    const rendered = renderRelease(bodyTemplate, unreleased)
    return rendered
  }

  // Case 3: Specific version
  const releasesDir = '.changelog/releases'
  let files: string[]
  try {
    files = readdirSync(releasesDir).filter((f) => f.endsWith('.yaml'))
  } catch {
    throw new Error(
      `Error: No releases directory found. Run 'changelog init' first.`,
    )
  }

  const matchFile = files.find((f) => basename(f, '.yaml') === argument)

  if (!matchFile) {
    throw new Error(
      `Error: Release "${argument}" not found in .changelog/releases/`,
    )
  }
  const filePath = join(releasesDir, matchFile)
  const text = await Bun.file(filePath).text()
  const release = Bun.YAML.parse(text) as ReleaseData
  const bodyTemplate = await loadBodyTemplate()
  const rendered = renderRelease(bodyTemplate, release)
  return rendered
}

export const viewCommand = createCommand(
  {
    usageName: 'changelog view',
    description: 'View the changelog',
  },
  async (cli) => {
    const argument = cli.input[0]
    console.log(await viewHandler(argument))
  },
)