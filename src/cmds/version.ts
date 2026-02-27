import { createCommand } from '@d-dev/roar'
import { loadConfig } from '../lib/config'

export async function versionHandler(): Promise<string> {
  const config = await loadConfig()
  return config.currentVersion
}

export const versionCommand = createCommand(
  {
    usageName: 'changelog version',
    description: 'Print the current version of the project',
  },
  async () => {
    console.log(await versionHandler())
  },
)