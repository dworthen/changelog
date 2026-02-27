import { createCommand } from '@d-dev/roar'
import { $ } from 'bun'
import { loadConfig } from '../lib/config'

const NEXT_DIR = '.changelog/next'

export const checkCommand = createCommand(
  {
    usageName: 'changelog check',
    description:
      'Check that at least one changelog entry exists in .changelog/next for the current branch',
  },
  async () => {
    const config = await loadConfig()
    const mainBranch = config.git.mainBranch

    const mergeBase = (await $`git merge-base ${mainBranch} HEAD`.text()).trim()

    const files = (
      await $`git diff --name-only ${mergeBase} HEAD -- ${NEXT_DIR}`.text()
    )
      .trim()
      .split('\n')
      .filter((f) => f.length > 0)

    if (files.length === 0) {
      console.error(
        'Error: No changelog entries found in .changelog/next/ for this branch. Please add a changelog entry before merging.',
      )
      process.exit(1)
    }

    console.log(
      `Check passed: found ${files.length} changelog ${files.length === 1 ? 'entry' : 'entries'} in .changelog/next/`,
    )
  },
)