// pages.js 汇总各页面 i18n 模块并合并到主语言包
// Aggregates per-page i18n modules and merges them into the main messages.
// 版本: v0.0.21
// 日期: 2026-09-17
// 作者: opencode
const modules = import.meta.glob('./pages/*.js', { eager: true })

// deepMerge 递归合并对象，避免不同模块共用同一命名空间时相互覆盖
// Recursively merges objects so modules sharing a namespace do not clobber each other.
// 版本: v0.0.21
// 日期: 2026-09-17
function deepMerge(target, source) {
  for (const key of Object.keys(source || {})) {
    const sv = source[key]
    if (sv && typeof sv === 'object' && !Array.isArray(sv)) {
      if (!target[key] || typeof target[key] !== 'object') target[key] = {}
      deepMerge(target[key], sv)
    } else {
      target[key] = sv
    }
  }
  return target
}

const pageMessages = { zh: {}, en: {} }

for (const path in modules) {
  const mod = modules[path].default || {}
  for (const locale of ['zh', 'en']) {
    deepMerge(pageMessages[locale], mod[locale] || {})
  }
}

export default pageMessages
