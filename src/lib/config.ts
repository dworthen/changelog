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
  return Bun.YAML.parse(text) as Config
}

export async function saveConfig(config: Config): Promise<void> {
  const text = Bun.YAML.stringify(config, null, 2)
  await Bun.write(CONFIG_PATH, text)
}