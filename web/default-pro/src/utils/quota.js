// quota.js 金额显示纯函数（可在 Node 环境下单测，不依赖 Vue/Arco）
// Currency display helpers, unit-testable in Node.
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode
//
// 单位口径（v0.0.25 起）：
//   - 后端在所有涉及金额的 API 上统一返回「元」（浮点），数据库内部仍存微元；
//   - 前端只负责「显示格式化」，不再做任何 1e6 换算；
//   - 提交给后端的金额（quota / remain_quota / bonus_quota）同样直接使用「元」。
//
// Unit convention (since v0.0.25):
//   - Every amount-bearing API returns CNY yuan (micro-quota stays the storage unit);
//   - the frontend only formats values and never divides/multiplies by 1e6;
//   - amounts sent to the backend (quota / remain_quota / bonus_quota) are also yuan.

// YUAN_INPUT_PRECISION 「元」输入框的最大小数位数。
// 1 微元 = 1e-6 元，因此 6 位可无损往返后端传来的任意金额。
// Max decimals for CNY inputs; 1 micro-quota = 1e-6 CNY, so 6 digits round-trip losslessly.
//
// 版本: v0.0.25
// 日期: 2026-10-04
export const YUAN_INPUT_PRECISION = 6

// YUAN_DEFAULT_PRECISION 「元」输入框的常规小数位（金额恰为「分」的整数倍时使用）。
// Default decimals for CNY inputs (used when the amount is a whole number of cents).
//
// 版本: v0.0.25
// 日期: 2026-10-04
export const YUAN_DEFAULT_PRECISION = 2

// DEFAULT_TOKEN_QUOTA_YUAN 新建令牌时的默认剩余额度（元）。
// 与旧口径的 DEFAULT_QUOTA_PER_UNIT（1e6 微元）等值，均为 1 元。
// Default remain quota (CNY) for a new token; equals the previous 1e6 micro-quota default.
//
// 版本: v0.0.25
// 日期: 2026-10-04
export const DEFAULT_TOKEN_QUOTA_YUAN = 1

// formatYuan 格式化为 ¥ 金额字符串；小数位随量级自适应（小额保留更多位，避免显示为 ¥0.00）。
// Format a CNY amount; decimal places adapt to magnitude so small values aren't shown as ¥0.00.
//
// 版本: v0.0.23
// 日期: 2026-10-03
export function formatYuan(amount) {
  const v = Number(amount) || 0
  const abs = Math.abs(v)
  const digits = abs >= 1 ? 2 : abs >= 0.01 ? 4 : 6
  const text = v.toLocaleString('en-US', {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  })
  return `¥${text}`
}

// formatQuota 将金额格式化为 ¥ 字符串；空值返回 '-'，便于表格直接渲染。
// 入参已是「元」（后端口径），不再做任何换算。
// Format an already-yuan amount; returns '-' for empty input so tables can render it directly.
//
// 版本: v0.0.25
// 日期: 2026-10-04
export function formatQuota(amount) {
  if (amount == null || amount === '') return '-'
  const n = Number(amount)
  if (!Number.isFinite(n)) return String(amount)
  return formatYuan(n)
}

// yuanInputPrecision 依金额自适应输入精度：恰为「分」的整数倍时用 2 位，否则用 6 位。
// 这样常规金额保持 2 位小数观感，而 ¥0.0005 这类小额仍能无损输入/回填。
// Pick input precision by magnitude: 2 decimals for whole cents, 6 otherwise.
//
// 版本: v0.0.25
// 日期: 2026-10-04
export function yuanInputPrecision(amount) {
  const n = Number(amount)
  if (!Number.isFinite(n)) return YUAN_DEFAULT_PRECISION
  // 用「乘以 100 后是否仍为整数」判断是否为「分」的整数倍；
  // 加 1e-9 容差以吸收浮点表示误差（如 0.07 * 100 = 7.000000000000001）。
  // Use a small epsilon because 0.07 * 100 === 7.000000000000001 in IEEE 754.
  const cents = n * 100
  return Math.abs(cents - Math.round(cents)) < 1e-9 ? YUAN_DEFAULT_PRECISION : YUAN_INPUT_PRECISION
}
