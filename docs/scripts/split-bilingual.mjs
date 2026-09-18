// docs/scripts/split-bilingual.mjs
// Strip bilingual echoes from docs/**/*-zh.md and docs/**/*-en.md so each
// file is single-language per the project convention.
//
// Patterns handled:
//   1. Heading suffix ` / English` -> drop the English part
//      e.g. `## 1. 准备 / Prepare` -> `## 1. 准备`
//   2. Table cell `Chinese / English` -> keep Chinese side
//      e.g. `| 变量 / Var | 默认 / Default |` -> `| 变量 | 默认 |`
//   3. Consecutive blockquote pair where one line is the other-language
//      echo -> drop the echo.
//      e.g.
//        > 中文简介。
//        > English summary.
//      keeps only the Chinese for -zh.md and the English for -en.md.
//   4. Changelog files: drop `## English` section from `*-zh.md`,
//      drop `## 中文` section from `*-en.md`.
//   5. A Chinese paragraph immediately followed by an English-only
//      paragraph (separated by a blank line) -> drop the echo.
//   6. A Chinese list item with an indented English continuation -> drop
//      the continuation (and mirror for en mode).
//
// Notes:
//   - Inside fenced code blocks nothing is rewritten (so that bilingual
//     comments like `1. 准备 release 分支\n   Prepare a release branch`
//     inside `text` fences stay intact for technical accuracy).
//   - Inside table cells only the `Chinese / English` separator pattern
//     is rewritten; we never touch the rest of the cell text.
//
// Usage:
//   node docs/scripts/split-bilingual.mjs          # process all
//   node docs/scripts/split-bilingual.mjs zh       # zh files only
//   node docs/scripts/split-bilingual.mjs en       # en files only
//   node docs/scripts/split-bilingual.mjs --dry    # preview to stdout

import { execFileSync } from 'node:child_process'
import { readFileSync, writeFileSync } from 'node:fs'
import { join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const DOCS = join(__dirname, '..')

const args = process.argv.slice(2)
const dryRun = args.includes('--dry')
const only = args.find((a) => a === 'zh' || a === 'en')

function exec(cmd, args) {
  return execFileSync(cmd, args, { cwd: join(__dirname, '..', '..'), encoding: 'utf8' })
}

// ---------------------------------------------------------------------------
// Heuristics
// ---------------------------------------------------------------------------
const CJK = /[\u3400-\u9fff\uF900-\uFAFF]/
const LATIN = /[A-Za-z]/
const isChinese = (s) => CJK.test(s)
const isEnglishOnly = (s) => !CJK.test(s) && LATIN.test(s)
const isBlank = (s) => /^\s*$/.test(s)
const isHeading = (s) => /^#{1,6}\s+/.test(s)
const isList = (s) => /^\s*([-*+]|\d+\.)\s+/.test(s)
const isCodeFence = (s) => /^```/.test(s)
const isBlockquote = (s) => /^>\s?/.test(s)

// ---------------------------------------------------------------------------
// Cleaners
// ---------------------------------------------------------------------------

// Strip ` / Anything` from heading lines: `## 1. 准备 / Prepare` -> `## 1. 准备`.
// Keep only the first segment (left of the `/`) so we keep the language we
// are targeting (Chinese for zh files, English for en files).
function stripHeadingBilingual(text) {
  return text.replace(
    /^(#{1,6}\s+)(.+?)\s+\/\s+.+?$/gm,
    (_m, hash, head) => `${hash}${head.trimEnd()}`
  )
}

// Strip ` / English` from each table cell on a line.
function stripTableCellBilingual(text) {
  return text
    .split('\n')
    .map((line) => {
      if (!/^\s*\|/.test(line)) return line
      return line
        .split('|')
        .map((cell) => {
          const i = cell.search(/\s+\/\s+[A-Za-z]/)
          return i < 0 ? cell : cell.slice(0, i).trimEnd()
        })
        .join('|')
    })
    .join('\n')
}

// Drop `## English` (or `## 中文`) section for changelog files.
function stripChangelogSection(text, keepChinese) {
  const lines = text.split('\n')
  const out = []
  let inFence = false
  let skip = false
  for (const line of lines) {
    if (isCodeFence(line)) {
      inFence = !inFence
      if (!skip) out.push(line)
      continue
    }
    if (inFence) {
      if (!skip) out.push(line)
      continue
    }
    if (line.startsWith('## ')) {
      if (/^## 中文\s*$/.test(line)) skip = !keepChinese
      else if (/^## English\s*$/.test(line)) skip = keepChinese
      else skip = false
    }
    if (!skip) out.push(line)
  }
  return out.join('\n')
}

// Drop the echo line in a consecutive pair of blockquote lines, keep the
// line in the target language.
function stripBlockquotePairLang(text, keepChinese) {
  const lines = text.split('\n')
  const out = []
  let i = 0
  while (i < lines.length) {
    const cur = lines[i]
    const next = lines[i + 1]
    if (isBlockquote(cur) && next && isBlockquote(next)) {
      const a = cur.replace(/^>\s?/, '')
      const b = next.replace(/^>\s?/, '')
      const aL = isChinese(a) ? 'zh' : isEnglishOnly(a) ? 'en' : 'other'
      const bL = isChinese(b) ? 'zh' : isEnglishOnly(b) ? 'en' : 'other'
      if (aL === 'zh' && bL === 'en') out.push(keepChinese ? cur : next)
      else if (aL === 'en' && bL === 'zh') out.push(keepChinese ? next : cur)
      else {
        out.push(cur)
        out.push(next)
      }
      i += 2
      continue
    }
    out.push(cur)
    i += 1
  }
  return out.join('\n')
}

// Drop echo paragraph when a paragraph is immediately followed (single blank
// line between) by its single-language translation. Preserves non-echo
// paragraphs and code fences verbatim.
function dropEchoParagraphs(text, keepChinese) {
  const lines = text.split('\n')
  const segments = []
  let cur = []
  let inFence = false
  let inContainer = false // ::: xxx ... :::
  let blanksSinceLast = 0
  for (const line of lines) {
    if (isCodeFence(line)) {
      if (cur.length) {
        segments.push({ lines: cur, sep: '\n' })
        cur = []
      }
      blanksSinceLast = 0
      inFence = !inFence
      cur.push(line)
      continue
    }
    if (inFence) {
      cur.push(line)
      continue
    }
    // VitePress markdown-it containers: `::: v-pre`, `::: details`, etc.
    // Opening `::: <name>` (no `:::` at end) and closing `:::`.
    if (/^:::\s/.test(line)) {
      if (cur.length) {
        segments.push({ lines: cur, sep: '\n' })
        cur = []
      }
      blanksSinceLast = 0
      inContainer = true
      cur.push(line)
      continue
    }
    if (inContainer) {
      cur.push(line)
      if (/^:::\s*$/.test(line)) {
        segments.push({ lines: cur, sep: '\n' })
        cur = []
        inContainer = false
      }
      continue
    }
    if (isBlank(line)) {
      blanksSinceLast++
      if (blanksSinceLast === 1) {
        if (cur.length) segments.push({ lines: cur, sep: '\n' })
        else segments.push({ lines: [], sep: '\n' })
        cur = []
        if (segments.length > 0) segments[segments.length - 1].sep = '\n\n'
      } else {
        segments.push({ lines: [], sep: '\n' })
        if (segments.length > 0) segments[segments.length - 1].sep = '\n\n'
      }
      continue
    }
    blanksSinceLast = 0
    cur.push(line)
  }
  if (cur.length) segments.push({ lines: cur, sep: '\n' })

  const out = []
  for (let i = 0; i < segments.length; i++) {
    const s = segments[i]
    if (s.lines.length === 0) {
      out.push(s)
      continue
    }
    const txt = s.lines.join('\n')
    const isStructural = s.lines.every(
      (l) => isList(l) || isHeading(l) || isBlockquote(l)
    )
    if (isStructural) {
      out.push(s)
      continue
    }
    let nextIdx = i + 1
    while (nextIdx < segments.length && segments[nextIdx].lines.length === 0)
      nextIdx++
    if (nextIdx < segments.length) {
      const next = segments[nextIdx]
      const nextTxt = next.lines.join('\n')
      const nextIsStructural = next.lines.every(
        (l) => isList(l) || isHeading(l) || isBlockquote(l)
      )
      if (!nextIsStructural) {
        const currIsChinese = isChinese(txt)
        const nextIsChinese = isChinese(nextTxt)
        const currIsEnglish = isEnglishOnly(txt)
        const nextIsEnglish = isEnglishOnly(nextTxt)
        if (keepChinese) {
          if (currIsChinese && nextIsEnglish) {
            out.push(s)
            segments[nextIdx] = { ...next, _skip: true }
            continue
          }
        } else {
          if (currIsEnglish && nextIsChinese) {
            out.push(s)
            segments[nextIdx] = { ...next, _skip: true }
            continue
          }
        }
      }
    }
    out.push(s)
  }
  const finalSegs = out.filter((s) => !s._skip)
  let result = ''
  for (const s of finalSegs) {
    result += s.lines.join('\n')
    result += s.sep
  }
  return result
}

// Drop the indented echo line under a list item.
function dropEchoListItems(text, keepChinese) {
  const lines = text.split('\n')
  const out = []
  let inFence = false
  for (let i = 0; i < lines.length; i++) {
    const cur = lines[i]
    if (isCodeFence(cur)) {
      inFence = !inFence
      out.push(cur)
      continue
    }
    if (inFence) {
      out.push(cur)
      continue
    }
    if (i === 0) {
      out.push(cur)
      continue
    }
    const prev = lines[i - 1]
    if (
      /^\s{2,}\S/.test(cur) &&
      isList(prev) &&
      (keepChinese ? isChinese(prev) : isEnglishOnly(prev)) &&
      (keepChinese ? isEnglishOnly(cur) : isChinese(cur)) &&
      !isList(cur)
    ) {
      continue
    }
    out.push(cur)
  }
  return out.join('\n')
}

// ---------------------------------------------------------------------------
// Pipeline
// ---------------------------------------------------------------------------
function processFile(text, keepChinese, isChangelog) {
  let t = text
  if (isChangelog) t = stripChangelogSection(t, keepChinese)
  t = stripHeadingBilingual(t)
  t = stripTableCellBilingual(t)
  t = stripBlockquotePairLang(t, keepChinese)
  t = dropEchoParagraphs(t, keepChinese)
  t = dropEchoListItems(t, keepChinese)
  t = t.replace(/\n{3,}/g, '\n\n')
  return t
}

// ---------------------------------------------------------------------------
// Driver
// ---------------------------------------------------------------------------
function listFiles(suffix) {
  return exec('find', ['docs', '-name', `*-${suffix}.md`])
    .split('\n')
    .filter(Boolean)
    .map((p) => p.replace(/^docs\//, ''))
}

let total = 0
let changed = 0

for (const lang of ['zh', 'en']) {
  if (only && only !== lang) continue
  const keepChinese = lang === 'zh'
  const files = listFiles(lang)
  for (const rel of files) {
    const full = join(DOCS, rel)
    const before = readFileSync(full, 'utf8')
    const after = processFile(before, keepChinese, rel.startsWith('changelog/'))
    total++
    if (after !== before) {
      if (dryRun) {
        process.stdout.write(`--- would change ${rel} (${before.length} -> ${after.length} bytes)\n`)
      } else {
        writeFileSync(full, after, 'utf8')
      }
      changed++
    }
  }
}

console.log(`${dryRun ? 'preview' : 'processed'} files=${total} changed=${changed}`)
