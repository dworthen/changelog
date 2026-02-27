import { $ } from 'bun'

export async function getFileSha(filePath: string): Promise<string> {
  try {
    const result =
      await $`git log --diff-filter=A --format=%H -- ${filePath}`.text()
    const sha = result.trim()
    if (!sha) {
      console.warn(`Warning: No git commit found for ${filePath}`)
      return 'unknown'
    }
    return sha
  } catch {
    console.warn(`Warning: Could not get git SHA for ${filePath}`)
    return 'unknown'
  }
}

export function getShortSha(sha: string): string {
  return sha.slice(0, 7)
}