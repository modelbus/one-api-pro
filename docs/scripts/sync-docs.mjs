// docs/scripts/sync-docs.mjs
// Sync docs/<category>/<slug>-{zh|en}.md into the VitePress doc-site.
//
// Locale layout (VitePress `root` locale convention):
//   zh (default) -> docs/doc-site/docs/<category>/<slug>.md   + .../docs/index.md
//   en           -> docs/doc-site/docs/en/<category>/<slug>.md + .../docs/en/index.md
//
// Usage (run from the doc-site package via `npm run sync`, or directly):
//   node docs/scripts/sync-docs.mjs          # sync both languages
//   node docs/scripts/sync-docs.mjs zh       # sync only zh
//   node docs/scripts/sync-docs.mjs en       # sync only en
//
// Behavior:
//   - Strips the `-<lang>` suffix from file names
//   - Wipes only the language's own category dirs + index before copying,
//     so removed pages disappear without touching public/ or the other locale
//   - Copies docs/assets -> docs/doc-site/docs/public/assets

import {
  rmSync,
  mkdirSync,
  readdirSync,
  copyFileSync,
  statSync,
  existsSync,
  writeFileSync,
  readFileSync,
} from 'node:fs'
import { join, dirname, basename } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
// This script lives at docs/scripts/, so docs/ is one level up.
const DOCS = join(__dirname, '..')
const SITE = join(DOCS, 'doc-site')
const SITE_DOCS = join(SITE, 'docs')
const SITE_PUBLIC = join(SITE_DOCS, 'public')

const only = process.argv[2] // 'zh' | 'en' | undefined
const LANGS = only ? [only] : ['zh', 'en']

// Language-specific categories under docs/. Skip non-content dirs:
// "assets" (static files), "doc-site" (the site project) and "scripts".
const NON_CATEGORY_DIRS = new Set(['assets', 'doc-site', 'scripts'])
const CATEGORIES = readdirSync(DOCS).filter((name) => {
  const full = join(DOCS, name)
  return statSync(full).isDirectory() && !NON_CATEGORY_DIRS.has(name)
})

// zh lives at the site root; en lives under /en/.
function langBase(lang) {
  return lang === 'zh' ? SITE_DOCS : join(SITE_DOCS, 'en')
}

function cleanDest(lang) {
  const base = langBase(lang)
  for (const cat of CATEGORIES) {
    const target = join(base, cat)
    if (existsSync(target)) rmSync(target, { recursive: true, force: true })
  }
  const index = join(base, 'index.md')
  if (existsSync(index)) rmSync(index, { force: true })
}

function syncCategory(lang) {
  const suffix = `-${lang}.md`
  const base = langBase(lang)
  for (const cat of CATEGORIES) {
    const catDir = join(DOCS, cat)
    const targetDir = join(base, cat)
    mkdirSync(targetDir, { recursive: true })
    for (const file of readdirSync(catDir)) {
      if (!file.endsWith(suffix)) continue
      const slug = basename(file, suffix)
      const raw = readFileSync(join(catDir, file), 'utf8')
      writeFileSync(join(targetDir, `${slug}.md`), escapeAngleBracketsInTables(raw), 'utf8')
    }
  }
}

function syncRoot(lang) {
  const src = join(DOCS, `index-${lang}.md`)
  const dst = join(langBase(lang), 'index.md')
  if (!existsSync(src)) return
  mkdirSync(dirname(dst), { recursive: true })
  copyFileSync(src, dst)
}

// Escape raw `<` inside markdown tables so VitePress's Vue parser
// does not interpret tokens like `array<string>` as unclosed HTML.
function escapeAngleBracketsInTables(text) {
  const lines = text.split('\n')
  let inFence = false
  const out = []
  for (const line of lines) {
    if (/^```/.test(line)) {
      inFence = !inFence
      out.push(line)
      continue
    }
    if (inFence || !/^\s*\|/.test(line)) {
      out.push(line)
      continue
    }
    out.push(line.replace(/<(?!\/?[a-zA-Z][\w-]*(?:\s+[^<>]*)?>)/g, '&lt;'))
  }
  return out.join('\n')
}

function syncAssets() {
  if (!existsSync(join(DOCS, 'assets'))) return
  mkdirSync(SITE_PUBLIC, { recursive: true })
  copyDir(join(DOCS, 'assets'), join(SITE_PUBLIC, 'assets'))
  // Expose the project logo at a stable public path for the home hero.
  const logo = join(DOCS, 'logo.png')
  if (existsSync(logo)) copyFileSync(logo, join(SITE_PUBLIC, 'logo.png'))
}

function copyDir(src, dst) {
  mkdirSync(dst, { recursive: true })
  for (const entry of readdirSync(src)) {
    const s = join(src, entry)
    const d = join(dst, entry)
    if (statSync(s).isDirectory()) copyDir(s, d)
    else copyFileSync(s, d)
  }
}

for (const lang of LANGS) {
  cleanDest(lang)
  syncCategory(lang)
  syncRoot(lang)
}
syncAssets()

console.log(`synced_langs=${LANGS.join(',')}`)
