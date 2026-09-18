// docs/scripts/split-changelog-en.mjs
// Split root CHANGELOG/vX.Y.Z.md (bilingual) into docs/changelog/vX.Y.Z-en.md
// by extracting the ## English section onward and rewriting the lead/title.
//
// Usage: node docs/scripts/split-changelog-en.mjs

import { readdirSync, readFileSync, writeFileSync, existsSync } from 'node:fs'
import { join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const REPO = join(__dirname, '..', '..')
const SRC = join(REPO, 'CHANGELOG')
const DST = join(__dirname, '..', 'changelog')

const files = readdirSync(SRC).filter((f) => /^v\d+\.\d+\.\d+\.md$/.test(f)).sort()

let written = 0
for (const f of files) {
  const ver = f.replace(/\.md$/, '')
  const srcPath = join(SRC, f)
  const enPath = join(DST, `${ver}-en.md`)
  const text = readFileSync(srcPath, 'utf8')
  const lines = text.split('\n')

  // Locate `## English`
  const enIdx = lines.findIndex((l) => /^## English\s*$/.test(l))
  if (enIdx < 0) {
    console.warn(`skip ${ver}: no ## English section`)
    continue
  }

  // First H1 "# vX.Y.Z — title" — split the Chinese vs English subtitle.
  // The opening blockquote (lines immediately following the H1) usually has:
  //   > <Chinese line>
  //   > <English line>
  // We extract the English line if present.
  let titleLine = lines[0]
  let enBlockquoteLine = null
  for (let i = 1; i < Math.min(lines.length, 8); i++) {
    const m = lines[i].match(/^>\s*(.*)$/)
    if (!m) break
    if (enBlockquoteLine === null) {
      // First blockquote line = Chinese
      enBlockquoteLine = m[1]
    } else {
      // Second blockquote line = English (override)
      enBlockquoteLine = m[1]
    }
  }

  // Build English content: skip the leading `## English` H2 (the whole page
  // is English) and take the rest. Wrap in `::: v-pre` so VitePress's Vue
  // SSR does not evaluate GitHub Actions expressions like `${{ secrets.X }}`
  // that appear in some changelogs.
  const enBodyRaw = lines.slice(enIdx + 1).join('\n').trimEnd()
  const enBody = '::: v-pre\n\n' + enBodyRaw + '\n\n:::'

  const frontMatter =
    `---\ntitle: ${ver}\ndescription: One API Pro ${ver} changelog (English)\n` +
    `category: changelog\norder: 0\n---\n\n`

  const content =
    frontMatter +
    '# ' + ver + '\n\n' +
    (enBlockquoteLine ? `> ${enBlockquoteLine}\n\n` : '') +
    enBody + '\n'

  writeFileSync(enPath, content, 'utf8')
  written++
}

console.log(`en_files_written=${written}`)