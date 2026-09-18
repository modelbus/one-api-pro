// docs/scripts/migrate-changelog.mjs
// Move CHANGELOG/vX.Y.Z.md -> docs/changelog/vX.Y.Z-zh.md and create -en stubs.

import { readdirSync, readFileSync, writeFileSync, existsSync, mkdirSync } from 'node:fs'
import { join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
// This script lives at docs/scripts/, so docs/ is one level up and the
// repository root is two levels up.
const DOCS = join(__dirname, '..')
const REPO = join(DOCS, '..')
const SRC = join(REPO, 'CHANGELOG')
const DST = join(DOCS, 'changelog')

const files = readdirSync(SRC).filter((f) => /^v\d+\.\d+\.\d+\.md$/.test(f)).sort()

let count = 0
for (const f of files) {
  const ver = f.replace(/\.md$/, '')
  const srcPath = join(SRC, f)
  const zhPath = join(DST, `${ver}-zh.md`)
  const enPath = join(DST, `${ver}-en.md`)

  const body = readFileSync(srcPath, 'utf8')
  // The changelog may contain GitHub Actions expressions like
  // `${{ secrets.GITHUB_TOKEN }}` which VitePress's Vue SSR tries to
  // evaluate. Wrap the body in `::: v-pre ... :::` to opt out of Vue
  // interpolation for the entire page, and prepend front-matter.
  const frontMatter =
    `---\ntitle: ${ver}\ndescription: One API Pro ${ver} 更新日志\n` +
    `category: changelog\norder: 0\n---\n\n`
  const content = frontMatter + '::: v-pre\n\n' + body.trim() + '\n\n:::\n'
  writeFileSync(zhPath, content, 'utf8')

  if (!existsSync(enPath)) {
    writeFileSync(enPath,
      `---\ntitle: ${ver}\ndescription: One API Pro ${ver} changelog\ncategory: changelog\norder: 0\n---\n\n# ${ver}\n\n> The Chinese version contains the full bilingual notes. Translation will follow.\n`,
      'utf8')
  }
  count++
}

// Write changelog index
// Note: links are written WITHOUT language suffix because the sync step
// renames `<slug>-<lang>.md` -> `<slug>.md` per language.
const versions = files.map((f) => f.replace(/\.md$/, '')).reverse()
const indexZh =
  '---\ntitle: 更新日志索引\ndescription: 所有历史版本的索引与升级要点\ncategory: changelog\norder: 1\n---\n\n' +
  '# 更新日志\n\n' +
  '按版本倒序排列。点击进入查看详情。\n\n' +
  versions.map((v) => `- [${v}](${v}.md)`).join('\n') +
  '\n'
writeFileSync(join(DST, 'index-zh.md'), indexZh, 'utf8')

const indexEn =
  '---\ntitle: Changelog Index\ndescription: Index of all releases with upgrade highlights\ncategory: changelog\norder: 1\n---\n\n' +
  '# Changelog\n\n' +
  'Ordered from newest to oldest.\n\n' +
  versions.map((v) => `- [${v}](${v}.md)`).join('\n') +
  '\n'
writeFileSync(join(DST, 'index-en.md'), indexEn, 'utf8')

console.log(`versions_migrated=${count}`)
