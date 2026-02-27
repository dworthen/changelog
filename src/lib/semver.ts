import { type ChangeType } from './types'

export function bumpVersion(
  current: string,
  changeTypes: ChangeType[],
): string {
  const [major, minor, patch] = current.split('.').map(Number)

  if (changeTypes.includes('Change') || changeTypes.includes('Remove')) {
    return `${major! + 1}.0.0`
  }
  if (changeTypes.includes('Add') || changeTypes.includes('Deprecate')) {
    return `${major}.${minor! + 1}.0`
  }
  if (changeTypes.includes('Fix')) {
    return `${major}.${minor}.${patch! + 1}`
  }
  return `${major}.${minor}.${patch}`
}