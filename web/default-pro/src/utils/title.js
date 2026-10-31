// 浏览器 Tab 标题纯函数工具（可在 Node 环境下单元测试，不依赖 Vue/Arco）
// Browser tab title helpers (unit-testable in Node, no Vue/Arco dependency).
// 版本: v0.0.27
// 日期: 2026-10-31
// 作者: opencode

// SITE_NAME 站点名固定来自前端，避免被后端 system_name（默认可能是 one-api）覆盖
// The site name is fixed on the frontend so backend system_name cannot reset the tab title.
export const SITE_NAME = 'ONE-API-PRO'

// TITLE_SEPARATOR 页面标题与站点名之间的分隔符（与 index.html 静态标题保持一致）
export const TITLE_SEPARATOR = '—'

// isUntranslated 判断 i18n 翻译结果是否为「未命中的 key」（形如 route.xxx）
// Detects whether an i18n lookup fell through and returned the key itself.
function isUntranslated(value, key) {
  if (typeof value !== 'string' || !value.trim()) return true
  return !!key && value === key
}

// resolvePageTitle 把路由 meta.title 的 i18n key 解析为可见标题
// 返回空串表示该路由没有可用的页面标题，标题应由站点默认文案兜底。
// Resolves a route meta.title i18n key into visible text; '' means "no page title".
export function resolvePageTitle(metaTitle, translate) {
  if (!metaTitle || typeof translate !== 'function') return ''
  const translated = translate(metaTitle)
  return isUntranslated(translated, metaTitle) ? '' : translated.trim()
}

// buildDocumentTitle 拼装最终的 document.title
// 有页面标题 → 「页面—站点名」；无页面标题 → 「站点名—站点默认标题」
// Composes the final document.title; falls back to the site default copy.
export function buildDocumentTitle(pageTitle, baseTitle, siteName = SITE_NAME) {
  const site = (siteName || '').trim() || SITE_NAME
  const page = typeof pageTitle === 'string' ? pageTitle.trim() : ''
  if (page) return `${page}${TITLE_SEPARATOR}${site}`
  const base = typeof baseTitle === 'string' ? baseTitle.trim() : ''
  return base ? `${site}${TITLE_SEPARATOR}${base}` : site
}
