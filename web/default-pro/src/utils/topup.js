// 充值模块纯函数工具（可在 Node 环境下单元测试，不依赖 Vue/Arco）
// Top-up module pure helpers (unit-testable in Node, no Vue/Arco dependency).
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode

// 说明：后端自 v0.0.25 起在充值设置 API 上统一返回「元」，
// 因此 amount 与 bonus_quota 可直接比较，无需再做额度换算。
// Since v0.0.25 the top-up settings API returns CNY yuan, so amount and bonus_quota
// are directly comparable.

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
// 规则：金额必须大于 0；不允许重复金额；到账金额不得低于支付金额
// （允许营销赠送，不允许缩水）。两个字段均为「元」口径，直接比较。
// 可选传入 vue-i18n 的 t 函数以输出本地化文案（默认中文，便于单测）。
//
// 注意：调用方传入的 t 必须是「已绑定命名空间」的函数
// （如 key => t('settingPage.topup.' + key)），否则 vue-i18n 找不到 key 会
// 原样返回 key 字符串，用户界面会出现 "bonusTooLow" 这样的原始 key。
//
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode
export function validateTopupPresets(presets, t) {
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
    if (!Number.isFinite(bonus) || bonus < amt) {
      return tr('bonusTooLow', { n: i + 1 })
    }
  }
  return null
}
