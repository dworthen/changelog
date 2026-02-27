import { resolve } from 'node:path'
import { argv } from 'bun'

const builds: Record<string, string> = {
  'bun-linux-arm64': 'bin/linux-arm64/changelog',
  'bun-linux-arm64-musl': 'bin/linux-arm64-musl/changelog',
  'bun-linux-x64-modern': 'bin/linux-x64/changelog',
  'bun-linux-x64-musl-modern': 'bin/linux-x64-musl/changelog',
  'bun-windows-x64-modern': 'bin/win-x64/changelog',
  'bun-darwin-arm64': 'bin/darwin-arm64/changelog',
  'bun-darwin-x64': 'bin/darwin-x64/changelog',
}

async function buildTarget(target: string, outFile: string): Promise<void> {
  console.log(`Building for target: ${target}...`)
  await Bun.build({
    entrypoints: ['./src/index.ts'],
    compile: {
      // @ts-expect-error
      target,
      outfile: resolve(outFile),
      autoloadTsConfig: false,
      autoloadPackageJson: false,
      autoloadBunConfig: false,
      autoloadDotEnv: true,
    },
    minify: true,
    sourcemap: 'linked',
    env: 'disable',
  })
  console.log(`Built ${target} successfully! Output: ${outFile}`)
}

const target = argv[2]
if (target && builds[target]) {
  await buildTarget(target, builds[target])
} else {
  await Promise.all(
    Object.entries(builds).map(
      async ([target, outFile]) => await buildTarget(target, outFile),
    ),
  )
}