#!/usr/bin/env node

import { spawn } from "node:child_process";
import { chmod } from "node:fs/promises";

const supported_platforms = new Map([
  ["win32-arm64", "@d-dev/changelog-win32-arm64"],
  ["win32-x64", "@d-dev/changelog-win32-x64"],
  ["darwin-arm64", "@d-dev/changelog-darwin-arm64"],
  ["darwin-x64", "@d-dev/changelog-darwin-x64"],
  ["linux-arm64", "@d-dev/changelog-linux-arm64"],
  ["linux-x64", "@d-dev/changelog-linux-x64"],
]);

async function run() {
  const os = process.platform;
  const cpu = process.arch;
  const key = `${os}-${cpu}`;
  if (!supported_platforms.has(key)) {
    throw new Error(`Platform ${key} is not supported`);
  }

  const args = process.argv.slice(2);
  const { default: getBinPath } = await import(supported_platforms.get(key));
  const binPath = getBinPath();
  await chmod(binPath, 0o774);
  const cp = spawn(binPath, args, { stdio: "inherit" });
  cp.on("close", (code) => {
    if (code !== 0) {
      console.error(`Process exited with code ${code}`);
      process.exit(code);
    }
  });
}

await run();
