import { parse, stringify } from 'yaml'
import { type Config } from './types'

const CONFIG_PATH = '.changelog/config.yaml'

export async function loadConfig(): Promise<Config> {
  const file = Bun.file(CONFIG_PATH)
  const exists = await file.exists()
  if (!exists) {
    console.error(
      "Error: .changelog/config.yaml not found. Run 'changelog init' first.",
    )
    process.exit(1)
  }
  const text = await file.text()
  return parse(text) as Config
}

export async function saveConfig(config: Config): Promise<void> {
  const text = stringify(config)
  await Bun.write(CONFIG_PATH, text)
}