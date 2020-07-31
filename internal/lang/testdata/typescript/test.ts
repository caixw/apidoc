// SPDX-License-Identifier: MIT

const x = "//\""

/// line1

const y = '/* xx\' */'
const z = `/* xx${x} */`
const c = 'c'
const t = /test\\a/.test(z)

/**
 * line1
 * line2
 * line3
 */
