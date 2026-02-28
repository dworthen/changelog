import { readdirSync } from 'node:fs'
import { join } from 'node:path'

import { renderString } from './eta'
import { getFileSha, getShortSha } from './git'
import {
  CHANGE_TYPES,
  type ChangeDescription,
  type ChangeType,
  type ReleaseData,
} from './types'

export type NextEntry = {
  filePath: string
  timestamp: number
  type: ChangeType
  description: string
}

export async function loadNextEntries(): Promise<NextEntry[]> {
  const dir = '.changelog/next'
  let files: string[]
  try {
    files = readdirSync(dir).filter((f) => f.endsWith('.yaml'))
  } catch {
    return []
  }

  const entries: NextEntry[] = []
  for (const file of files) {
    const filePath = join(dir, file)
    try {
      const text = await Bun.file(filePath).text()
      const data = Bun.YAML.parse(text) as NextEntry
      entries.push({
        filePath,
        timestamp: data.timestamp,
        type: data.type,
        description: data.description,
      })
    } catch (err) {
      console.warn(`Warning: Could not parse ${filePath}: ${err}`)
    }
  }

  return entries.filter((e) => e.type !== 'Internal')
}

export async function buildReleaseData(
  version: string,
  entries: NextEntry[],
): Promise<ReleaseData> {
  const changes: ChangeDescription[] = []

  for (const entry of entries) {
    const sha = await getFileSha(entry.filePath)
    const shortSha = getShortSha(sha)
    changes.push({
      type: entry.type,
      timestamp: entry.timestamp,
      sha,
      shortSha,
      description: entry.description,
    })
  }

  const changeTypeSet = new Set(changes.map((c) => c.type))
  const availableChangeTypes: ChangeType[] = Object.keys(
    CHANGE_TYPES,
  ) as ChangeType[]
  const changeTypes: ChangeType[] = []
  for (const ct of availableChangeTypes) {
    if (changeTypeSet.has(ct)) {
      changeTypes.push(ct)
    }
  }

  return {
    timestamp: Date.now(),
    version,
    changeTypes,
    changes,
  }
}

export function groupChangesByType(
  changes: ChangeDescription[],
): Record<string, ChangeDescription[]> {
  const grouped: Record<string, ChangeDescription[]> = {}
  for (const change of changes) {
    if (!grouped[change.type]) {
      grouped[change.type] = []
    }
    grouped[change.type]!.push(change)
  }
  return grouped
}

export function formatDate(timestamp: number): string {
  const d = new Date(timestamp)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export async function loadReleases(): Promise<ReleaseData[]> {
  const dir = '.changelog/releases'
  let files: string[]
  try {
    files = readdirSync(dir).filter((f) => f.endsWith('.yaml'))
  } catch {
    return []
  }

  const releases: ReleaseData[] = []
  for (const file of files) {
    const filePath = join(dir, file)
    try {
      const text = await Bun.file(filePath).text()
      const data = Bun.YAML.parse(text) as ReleaseData
      releases.push(data)
    } catch (err) {
      console.warn(`Warning: Could not parse ${filePath}: ${err}`)
    }
  }

  releases.sort((a, b) => b.timestamp - a.timestamp)
  return releases
}

export async function generateChangelog(
  releases: ReleaseData[],
): Promise<string> {
  const header = await Bun.file('.changelog/templates/header.md').text()
  const bodyTemplate = await Bun.file('.changelog/templates/body.eta').text()
  const footer = await Bun.file('.changelog/templates/footer.md').text()

  const renderedBodies: string[] = []

  for (const release of releases) {
    const groupedChanges = groupChangesByType(release.changes)
    const rendered = renderString(bodyTemplate, {
      version: release.version,
      date: formatDate(release.timestamp),
      changeTypes: release.changeTypes,
      changes: groupedChanges,
    })
    renderedBodies.push(rendered)
  }

  let output = `${header.trim()}\n\n${renderedBodies.join('\n').trim()}`
  if (footer) {
    output += `\n\n${footer.trim()}`
  }
  return output
}

export async function loadBodyTemplate(): Promise<string> {
  return Bun.file('.changelog/templates/body.eta').text()
}

export function renderRelease(
  bodyTemplate: string,
  release: ReleaseData,
): string {
  const groupedChanges = groupChangesByType(release.changes)
  return renderString(bodyTemplate, {
    version: release.version,
    date: formatDate(release.timestamp),
    changeTypes: release.changeTypes,
    changes: groupedChanges,
  })
}