// scripts/split-api-doc.mjs
// Split legacy docs/API.md into docs/api/*.md by `## N. Title (Code)` headers.
// English placeholders are regenerated as before; Chinese gets the real content.

import { readFileSync, writeFileSync, existsSync } from 'node:fs'
import { join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
// This script lives at docs/scripts/, so docs/ is one level up.
const ROOT = join(__dirname, '..')
const API_MD = join(ROOT, 'API.md')

// Map: legacy section number -> [categoryFile, sectionTitleZh, sectionTitleEn]
// Sections in the new layout that don't have a dedicated placeholder
// (e.g. misc public endpoints) are appended to api/misc-zh.md.
const SECTION_MAP = {
  1:  ['model-price', '模型定价 API', 'Model Price API'],
  2:  ['group-price', '分组折扣 API', 'Group Price API'],
  3:  null,           // merged into misc
  4:  null,           // merged into misc
  5:  ['channel', '渠道 API', 'Channel API'],
  6:  ['token', '令牌 API', 'Token API'],
  7:  ['user', '用户 API', 'User API'],
  8:  ['log', '日志 API', 'Log API'],
  9:  ['redemption', '兑换码 API', 'Redemption API'],
  10: ['plan', '套餐 API', 'Plan API'],
  11: ['subscription', '订阅 API', 'Subscription API'],
  12: ['openai-compatible', 'OpenAI 兼容接口', 'OpenAI Compatible API'],
  13: ['cluster', '集群 API', 'Cluster API'],
  14: null,           // merged into misc
}

// Read source
const raw = readFileSync(API_MD, 'utf8')
const lines = raw.split('\n')

// Find section boundaries
const starts = []
lines.forEach((line, idx) => {
  const m = line.match(/^## (\d+)\. (.+)$/)
  if (m) starts.push({ idx, num: parseInt(m[1], 10), title: m[2] })
})

// Appendices (A..D) -> api/README-zh.md
const appendixStarts = []
lines.forEach((line, idx) => {
  if (/^## 附录/.test(line)) appendixStarts.push(idx)
})

// Build section bodies
const sections = {}
for (let i = 0; i < starts.length; i++) {
  const cur = starts[i]
  const nextIdx = i + 1 < starts.length ? starts[i + 1].idx : lines.length
  const body = lines.slice(cur.idx, nextIdx).join('\n').trimEnd()
  sections[cur.num] = { title: cur.title, body }
}

// Build appendix body
let appendixBody = ''
if (appendixStarts.length > 0) {
  const firstAppendix = appendixStarts[0]
  // Stop at first `## 1.` (top-level section)
  const firstSection = starts[0]?.idx ?? lines.length
  appendixBody = lines.slice(firstAppendix, firstSection).join('\n').trimEnd()
}

// Detect "## 14. 其他公共接口" content (everything after the last numbered section)
const lastNum = Math.max(...Object.keys(sections).map(Number))
const lastIdx = starts.find((s) => s.num === lastNum).idx
// All remaining lines after section 14 header itself if any, else after section 13
let miscBody = ''
const miscStart = starts.find((s) => s.num === 14)
if (miscStart) {
  miscBody = lines
    .slice(miscStart.idx, lines.length)
    .join('\n')
    .trimEnd()
}

// ---------------------------------------------------------------------------
// Write files
// ---------------------------------------------------------------------------
function writeTarget(slug, titleZh, titleEn, body, order, descZh, descEn) {
  const zhPath = join(ROOT, 'api', `${slug}-zh.md`)
  const enPath = join(ROOT, 'api', `${slug}-en.md`)

  const fm = (title, desc) =>
    [
      '---',
      `title: ${title}`,
      `description: ${desc}`,
      `category: api`,
      `order: ${order}`,
      '---',
      '',
    ].join('\n')

  const zhContent = fm(titleZh, descZh) + body + '\n'
  writeFileSync(zhPath, zhContent, 'utf8')

  // English placeholder only if not exists
  if (!existsSync(enPath)) {
    const enPlaceholder =
      fm(titleEn, descEn) +
      `# ${titleEn}\n\n` +
      `> ${descEn}\n\n` +
      '## Overview\n\nThis page mirrors the Chinese version. Translate or open a PR.\n'
    writeFileSync(enPath, enPlaceholder, 'utf8')
  }
}

let count = 0
let nextOrder = 11
for (const [numStr, mapping] of Object.entries(SECTION_MAP)) {
  if (!mapping) continue
  const num = parseInt(numStr, 10)
  const sec = sections[num]
  if (!sec) continue
  const [slug, titleZh, titleEn] = mapping
  const descZh = `${titleZh}：${sec.title}`
  const descEn = `${titleEn}: ${sec.title}`
  writeTarget(slug, titleZh, titleEn, sec.body, nextOrder++, descZh, descEn)
  count++
}

// Misc bucket: append sections 3, 4, 14 (and anything that mentions them)
function appendSection(target, num) {
  const sec = sections[num]
  if (!sec) return
  target += '\n\n' + sec.body + '\n'
  return target
}

let misc = '# 其他公共 API\n\n'
misc = appendSection(misc, 3) ?? misc
misc = appendSection(misc, 4) ?? misc
if (miscBody) misc += '\n\n' + miscBody + '\n'
writeTarget('misc', '其他公共 API', 'Misc Public API', misc, nextOrder,
  '分组列表、系统选项与其他公共端点。',
  'Group list, system options and other public endpoints.')
count++

// README: appendices + auth + permissions
const readmeBody =
  '# API 参考 · 总览\n\n' +
  '本文档汇总 One API Pro 的所有 `/api/*` 与 `/v1/*` 接口，以及鉴权、权限、计费等公共约定。\n\n' +
  appendixBody +
  '\n'

// Preserve README front-matter placeholder; replace body only
const readmeZhPath = join(ROOT, 'api', 'README-zh.md')
const readmeEnPath = join(ROOT, 'api', 'README-en.md')
writeFileSync(readmeZhPath,
  '---\ntitle: API 参考 · 总览\ndescription: 鉴权机制、响应格式、权限等级与术语表\ncategory: api\norder: 1\n---\n\n' +
  readmeBody,
  'utf8')

if (!existsSync(readmeEnPath)) {
  writeFileSync(readmeEnPath,
    '---\ntitle: API Reference · Overview\ndescription: Auth, response format, permission tiers and glossary\ncategory: api\norder: 1\n---\n\n' +
    '# API Reference · Overview\n\n> Translation pending. Refer to the Chinese version.\n',
    'utf8')
}
count++

console.log(`sections_written=${count}`)
