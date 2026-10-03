// quota.js 额度与货币换算纯函数（可在 Node 环境下单测，不依赖 Vue/Arco）
// Quota <-> currency conversion helpers, unit-testable in Node.
// 版本: v0.0.24
// 日期: 2026-10-04
// 作者: opencode

// 1 元对应的额度基准兜底值，与后端 common/config 的 QuotaPerUnit 常量一致。
// 语义同微信支付的「分」：1 元 = 1_000_000 额度（1 额度 = 1e-6 元）。
// Fallback quota-per-yuan base, aligned with backend common/config.QuotaPerUnit.
//
// 版本: v0.0.24
// 日期: 2026-10-03
export const DEFAULT_QUOTA_PER_UNIT = 1_000_000

// resolveQuotaPerUnit 归一 quota_per_unit：非正数/非法值回退到兜底基准。
// Normalize quota_per_unit, falling back to the default when invalid.
//
// 版本: v0.0.23
// 日期: 2026-10-03
export function resolveQuotaPerUnit(quotaPerUnit) {
  const n = Number(quotaPerUnit)
  return Number.isFinite(n) && n > 0 ? n : DEFAULT_QUOTA_PER_UNIT
}

// quotaToYuan 将额度换算为元（1 元 = quota_per_unit 额度）。
// Convert quota to CNY using the 1 CNY = quota_per_unit base.
//
// 版本: v0.0.23
// 日期: 2026-10-03
export function quotaToYuan(quota, quotaPerUnit) {
  return Number(quota || 0) / resolveQuotaPerUnit(quotaPerUnit)
}

// yuanToQuota 将金额（元）换算为额度；用 Math.round 消除浮点截断误差
// （与后端 model.YuanToQuota 保持一致，例如 0.07 元 → 70000 额度）。
// Convert CNY amount to quota with rounding, mirroring backend model.YuanToQuota.
//
// 版本: v0.0.24
// 日期: 2026-10-03
export function yuanToQuota(amount, quotaPerUnit) {
  const n = Number(amount)
  if (!Number.isFinite(n) || n <= 0) return 0
  return Math.round(n * resolveQuotaPerUnit(quotaPerUnit))
}

// YUAN_INPUT_PRECISION 「元」输入框的最大小数位，等于 1 额度对应的位数
// （quota_per_unit=1e6 时为 6 位，1e-6 元 = 1 额度，可无损表示任意 quota）。
// Max decimals for CNY inputs; equals the digits needed for 1 quota (6 when quota_per_unit=1e6).
//
// 版本: v0.0.24
// 日期: 2026-10-04
export const YUAN_INPUT_PRECISION = 6

// YUAN_DEFAULT_PRECISION 「元」输入框的常规小数位（金额 ≥ 分时使用，保持常见观感）。
// Default decimals for CNY inputs (used when the amount is a whole number of cents).
//
// 版本: v0.0.24
// 日期: 2026-10-04
export const YUAN_DEFAULT_PRECISION = 2

// quotaToYuanExact 额度 → 元，保留 6 位小数（无损：1e-6 元 = 1 额度）。
// 用于管理端「元」输入框回填，避免 toFixed(2) 把小于 ¥0.01 的额度静默截断为 0。
// Lossless quota → CNY conversion (6 decimals), for round-tripping admin CNY inputs.
//
// 版本: v0.0.24
// 日期: 2026-10-04
export function quotaToYuanExact(quota, quotaPerUnit) {
  return Number(quotaToYuan(quota, quotaPerUnit).toFixed(YUAN_INPUT_PRECISION))
}

// yuanInputPrecision 依金额自适应输入精度：为「分」的整数倍时用 2 位，否则用 6 位。
// 这样常规金额保持 2 位小数观感，而 ¥0.0005 这类小额仍能无损输入/回填。
// Pick input precision by magnitude: 2 decimals for whole cents, 6 otherwise.
//
// 版本: v0.0.24
// 日期: 2026-10-04
export function yuanInputPrecision(amount, quotaPerUnit) {
  const centQuota = resolveQuotaPerUnit(quotaPerUnit) / 100
  return yuanToQuota(amount, quotaPerUnit) % centQuota === 0
    ? YUAN_DEFAULT_PRECISION
    : YUAN_INPUT_PRECISION
}

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

// formatQuota 将额度按当前基准格式化为元；空值返回 '-'，便于表格直接渲染。
// Format quota as CNY; returns '-' for empty input so tables can render it directly.
//
// 版本: v0.0.23
// 日期: 2026-10-03
export function formatQuota(quota, quotaPerUnit) {
  if (quota == null || quota === '') return '-'
  const n = Number(quota)
  if (!Number.isFinite(n)) return String(quota)
  return formatYuan(quotaToYuan(n, quotaPerUnit))
}
