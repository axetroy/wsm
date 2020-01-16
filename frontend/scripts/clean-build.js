// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
const fs = require('fs')
const path = require('path')

const dir = path.join(process.cwd(), '.nuxt')

const files = fs.readdirSync(dir)

for (const filename of files) {
  const filepath = path.join(dir, filename)

  if (filename !== 'dist') {
    const stat = fs.statSync(filepath)
    if (stat.isDirectory()) {
      fs.rmdirSync(filepath, { recursive: true })
    } else {
      fs.unlinkSync(filepath)
    }
  }
}
