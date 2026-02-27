export type ChangeType =
  | 'Add'
  | 'Change'
  | 'Deprecate'
  | 'Remove'
  | 'Fix'
  | 'Internal'

export const CHANGE_TYPES: Record<ChangeType, string> = {
  Change: 'Change an existing feature. Major version bump.',
  Remove: 'Remove a feature. Major version bump.',
  Add: 'Add a new feature. Minor version bump.',
  Deprecate: 'Get ready to remove a feature. Minor version bump.',
  Fix: 'Fix a bug. Patch version bump.',
  Internal:
    'Internal change that does not affect users, e.g., refactors or doc updates. No version bump. An entry is added to .changelog/next in order to pass `changelog check` CI checks but these entries are not added to the public changelog.',
}

export type ChangeDescription = {
  type: ChangeType
  timestamp: number
  sha: string
  shortSha: string
  description: string
}

export type ReleaseData = {
  timestamp: number
  version: string
  changeTypes: ChangeType[]
  changes: ChangeDescription[]
}

export type Config = {
  currentVersion: string
  changelogFile: string
  git: {
    mainBranch: string
  }
  versionFiles: {
    path: string
    pattern: string
  }[]
  apply: {
    postCommands?: string[]
  }
}