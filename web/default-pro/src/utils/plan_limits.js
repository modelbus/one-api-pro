// plan_limits.js 套餐 model_limits 的行式编辑纯函数（行 ↔ JSON 互转 + 校验）
// Pure helpers converting a plan's model_limits JSON to/from editable rows
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode

// 计费维度常量（与后端 model.BillingType* 保持一致）
export const BILLING_TOKEN = 'token'
export const BILLING_REQUEST = 'request'

// 窗口类型常量（与后端 model.WindowType* 保持一致）
export const WINDOW_PERIOD = 'period'
export const WINDOW_WEEK = 'week'
export const WINDOW_MONTH = 'month'

// 三种窗口，顺序即表格列顺序（滚动窗口 / 每周 / 每月）
export const WINDOW_TYPES = [WINDOW_PERIOD, WINDOW_WEEK, WINDOW_MONTH]

// 行对象上各窗口对应的字段名（一行同时承载三种窗口的限额）
export const ROW_LIMIT_KEYS = {
  [WINDOW_PERIOD]: 'limitPeriod',
  [WINDOW_WEEK]: 'limitWeek',
  [WINDOW_MONTH]: 'limitMonth',
}

// 窗口小时数取值范围：后端 GetWindowDurationSeconds 在 period_h <= 0 时回退 5 小时，
// 前端下拉固定 1-24 小时，默认 5。
export const DEFAULT_PERIOD_H = 5
export const MIN_PERIOD_H = 1
export const MAX_PERIOD_H = 24

// normalizeBillingType 把任意输入归一为 token / request，非法值回退 token。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function normalizeBillingType(billingType) {
  return billingType === BILLING_REQUEST ? BILLING_REQUEST : BILLING_TOKEN
}

// limitFieldOf 返回 (窗口, 计费维度) 对应的后端 JSON 字段名，如 token_week。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function limitFieldOf(windowType, billingType) {
  const window = WINDOW_TYPES.includes(windowType) ? windowType : WINDOW_PERIOD
  return `${normalizeBillingType(billingType)}_${window}`
}

// normalizePeriodH 把窗口小时数夹到 [1, 24]；非法值回退默认 5 小时。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function normalizePeriodH(value) {
  const n = Math.floor(Number(value))
  if (!Number.isFinite(n) || n <= 0) return DEFAULT_PERIOD_H
  return Math.min(MAX_PERIOD_H, Math.max(MIN_PERIOD_H, n))
}

// normalizeLimit 把任意输入归一为非负整数限额；非法值回退 0（表示该窗口不限制）。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function normalizeLimit(value) {
  const n = Math.floor(Number(value))
  return Number.isFinite(n) && n > 0 ? n : 0
}

// createLimitRow 创建一个限制行（一个模型一行，三个窗口限额并列）。
//
// 三个限额默认 null（不限制）：表单里保持空值以显示 placeholder，
// 避免新增行时把 0 当成已填值，保存时才由 normalizeLimit 归一到 0/正整数。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function createLimitRow(model = '', overrides = {}) {
  return {
    model,
    periodH: DEFAULT_PERIOD_H,
    limitPeriod: null,
    limitWeek: null,
    limitMonth: null,
    ...overrides,
  }
}

// parseLimitsObject 宽松解析 model_limits（字符串 / 对象 / 空值）。
// 返回 { ok, limits, reason }：reason 为 'json'（语法错误）或 'type'（不是对象）。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function parseLimitsObject(raw) {
  if (raw === null || raw === undefined) return { ok: true, limits: {}, reason: '' }
  let value = raw
  if (typeof raw === 'string') {
    const text = raw.trim()
    if (!text) return { ok: true, limits: {}, reason: '' }
    try {
      value = JSON.parse(text)
    } catch (e) {
      return { ok: false, limits: {}, reason: 'json' }
    }
  }
  if (value === null || typeof value !== 'object' || Array.isArray(value)) {
    return { ok: false, limits: {}, reason: 'type' }
  }
  return { ok: true, limits: value, reason: '' }
}

// modelLimitsToRows 把 model_limits 转成可编辑的限制行（一个模型一行）。
//
// 只读取当前计费维度对应的字段（token_* 或 request_*）：后端 WeightedUsage 同样
// 只按 billing_type 取限额，另一维度的字段不参与计算，因此不展示、也不回写。
// 三个窗口限额全为空的模型仍保留一行，避免编辑时静默丢模型（保存时会被校验拦下）。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function modelLimitsToRows(raw, billingType = BILLING_TOKEN) {
  const billing = normalizeBillingType(billingType)
  const { ok, limits, reason } = parseLimitsObject(raw)
  if (!ok) return { ok: false, rows: [], reason }
  const rows = []
  for (const model of Object.keys(limits)) {
    const rule = limits[model]
    if (rule === null || typeof rule !== 'object' || Array.isArray(rule)) continue
    const row = createLimitRow(model, { periodH: normalizePeriodH(rule.period_h) })
    for (const windowType of WINDOW_TYPES) {
      // 后端 0 / 缺失都表示该窗口不限制，表单里统一呈现为空值（placeholder）。
      row[ROW_LIMIT_KEYS[windowType]] = normalizeLimit(rule[limitFieldOf(windowType, billing)]) || null
    }
    rows.push(row)
  }
  return { ok: true, rows, reason: '' }
}

// rowsToModelLimits 把限制行合成为后端 model_limits JSON 字符串。
//
// 某行三个窗口限额都未填时整行跳过（由 validateLimitRows 负责提示用户）。
// period_h 为模型级字段（仅滚动窗口使用），取值夹到 1-24。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function rowsToModelLimits(rows, billingType = BILLING_TOKEN) {
  const billing = normalizeBillingType(billingType)
  const limits = {}
  for (const row of Array.isArray(rows) ? rows : []) {
    const model = typeof row?.model === 'string' ? row.model.trim() : ''
    if (!model) continue
    const rule = { period_h: normalizePeriodH(row?.periodH) }
    let hasLimit = false
    for (const windowType of WINDOW_TYPES) {
      const value = normalizeLimit(row?.[ROW_LIMIT_KEYS[windowType]])
      if (value <= 0) continue
      rule[limitFieldOf(windowType, billing)] = value
      hasLimit = true
    }
    if (!hasLimit) continue
    limits[model] = rule
  }
  return { ok: true, json: JSON.stringify(limits) }
}

// validateLimitRows 校验限制行；通过返回 null，否则返回 { code, params }。
//
// code 是 settingPage.plan.* 下的 i18n key 后缀，由调用方翻译。
//
// 版本: v0.0.25
// 日期: 2026-10-05
export function validateLimitRows(rows) {
  const list = Array.isArray(rows) ? rows : []
  if (list.length === 0) return { code: 'errLimitsEmpty', params: {} }
  const seen = new Set()
  for (let i = 0; i < list.length; i += 1) {
    const row = list[i] || {}
    const n = i + 1
    const model = typeof row.model === 'string' ? row.model.trim() : ''
    if (!model) return { code: 'errLimitModel', params: { n } }
    if (seen.has(model)) return { code: 'errLimitDuplicate', params: { n, model } }
    seen.add(model)
    const hasLimit = WINDOW_TYPES.some((w) => normalizeLimit(row[ROW_LIMIT_KEYS[w]]) > 0)
    if (!hasLimit) return { code: 'errLimitValue', params: { n } }
  }
  return null
}
