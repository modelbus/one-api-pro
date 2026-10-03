// 充值模块纯函数工具（可在 Node 环境下单元测试，不依赖 Vue/Arco）
// Top-up module pure helpers (unit-testable in Node, no Vue/Arco dependency).
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode

import { resolveQuotaPerUnit, yuanToQuota } from './quota.js'

// formatNumber 简写大数字（>=10000 折算为 w；其余使用千分位）
//
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
export function formatNumber(n) {
  const v = Number(n) || 0
  if (v >= 10000) return (v / 10000).toFixed(2) + 'w'
  return v.toLocaleString()
}

// formatAmount 始终保留两位小数的金额字符串
//
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
export function formatAmount(n) {
  return Number(n || 0).toFixed(2)
}

// validateTopupPresets 前端预校验快捷金额，返回错误信息或 null。
// 规则：金额必须大于 0；不允许重复金额；到账额度不得低于支付金额折算额度
// （允许营销赠送，不允许缩水）。到账额度由「到账金额（元）」换算而来。
// 可选传入 vue-i18n 的 t 函数以输出本地化文案（默认中文，便于单测）。
//
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
export function validateTopupPresets(presets, t, quotaPerUnit) {
  const rate = resolveQuotaPerUnit(quotaPerUnit)
  const tr = typeof t === 'function'
    ? t
    : (key, params) => {
        if (key === 'amountPositive') return `第 ${params.n} 行金额必须大于 0`
        if (key === 'duplicate') return `快捷金额重复：${params.amt} 元已存在`
        if (key === 'bonusTooLow') return `第 ${params.n} 行到账金额不能低于支付金额`
        return key
      }
  const seen = new Set()
  for (let i = 0; i < presets.length; i++) {
    const amt = Number(presets[i].amount)
    if (!Number.isFinite(amt) || amt <= 0) {
      return tr('amountPositive', { n: i + 1 })
    }
    if (seen.has(amt)) {
      return tr('duplicate', { amt })
    }
    seen.add(amt)
    const bonus = Number(presets[i].bonus_quota)
    if (!Number.isFinite(bonus) || bonus < yuanToQuota(amt, rate)) {
      return tr('bonusTooLow', { n: i + 1 })
    }
  }
  return null
}
