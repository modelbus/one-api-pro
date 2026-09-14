// token.js 令牌提交/展示相关纯函数
// Pure helpers for token submission and rendering.
//
// 版本: v0.0.21
// 日期: 2026-09-14
// 作者: opencode

/**
 * 把表单的过期时间转换成后端约定的秒级 unix 时间戳。
 *
 * 与后端 sentinel 约定：
 *   -1 → 永不过期
 *   <=0 也视为永不过期（兜底：历史脏数据 / 旧前端 bug 写入 0）
 *   >0  → 具体的过期 unix 时间戳（秒）
 *
 * 注意：UI 上「永不过期」复选框或日期留空时，必须返回 -1 而不是 0。
 * 否则后端 ValidateUserToken 会把 0 视为 1970-01-01 而报"令牌已过期"。
 *
 * @param {object} form 表单状态对象
 * @param {boolean} form.never_expire 是否勾选「永不过期」
 * @param {number|null|undefined} form.expired_time 过期时间（毫秒时间戳或 null）
 * @returns {number} 后端约定的 expired_time（秒级 unix 时间戳或 -1）
 */
export function buildTokenExpiredTime(form) {
  const { never_expire, expired_time } = form || {}
  if (never_expire || !expired_time) return -1
  return Math.floor(expired_time / 1000)
}

/**
 * 渲染令牌过期时间展示文案。
 * 后端约定 <=0 都是"永不过期"，UI 一律显示「永不过期」。
 *
 * @param {number|null|undefined} ts 秒级 unix 时间戳（<=0 表示永不过期）
 * @returns {string} 展示文案
 */
export function formatExpiredTime(ts) {
  if (!ts) return '永不过期'
  const t = Number(ts)
  if (Number.isNaN(t) || t <= 0) return '永不过期'
  const d = new Date(t * 1000)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}