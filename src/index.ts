#!/usr/bin/env node

import { createCommand } from '@d-dev/roar'
import pkg from '../package.json'
import { addCommand } from './cmds/add.js'
import { applyCommand } from './cmds/apply.js'
import { checkCommand } from './cmds/check.js'
import { initCommand } from './cmds/init.js'
import { versionCommand } from './cmds/version.js'
import { viewCommand } from './cmds/view.js'

const changelogCommand = createCommand({
  usageName: 'changelog',
  description: pkg.description,
  version: pkg.version,
  versionFlag: 'version',
})

changelogCommand.addCommand('init', initCommand)
changelogCommand.addCommand('add', addCommand)
changelogCommand.addCommand('apply', applyCommand)
changelogCommand.addCommand('check', checkCommand)
changelogCommand.addCommand('view', viewCommand)
changelogCommand.addCommand('version', versionCommand)

await changelogCommand.run(process.argv.slice(2))