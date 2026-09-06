// 充值模块纯函数工具（可在 Node 环境下单元测试，不依赖 Vue/Arco）
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode

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
// 规则：金额必须大于 0；不允许重复金额。
//
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
export function validateTopupPresets(presets) {
  const seen = new Set()
  for (let i = 0; i < presets.length; i++) {
    const amt = Number(presets[i].amount)
    if (!Number.isFinite(amt) || amt <= 0) {
      return `第 ${i + 1} 行金额必须大于 0`
    }
    if (seen.has(amt)) {
      return `快捷金额重复：${amt} 元已存在`
    }
    seen.add(amt)
  }
  return null
}

// calcCustomBonus 自定义金额按 exchange_rate 算出到账 quota
// amount 非法时按 0 计算；rate 非法时按 1 兜底；rate <= 0 时按 0 兜底（防配置错误）。
//
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
export function calcCustomBonus(amount, exchangeRate) {
  const amt = Number(amount) || 0
  let rate
  if (exchangeRate === undefined || exchangeRate === null || exchangeRate === '') {
    rate = 1
  } else {
    rate = Number(exchangeRate)
    if (!Number.isFinite(rate) || rate <= 0) rate = 0
  }
  return Math.round(amt * rate)
}
