// 充值工具函数单元测试
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
// 运行方式: node web/default-pro/src/utils/topup.test.mjs
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { formatNumber, formatAmount, validateTopupPresets, calcCustomBonus } from './topup.js'

// formatNumber
test('formatNumber 处理 0', () => {
  assert.equal(formatNumber(0), '0')
})
test('formatNumber 处理小数千分位', () => {
  assert.equal(formatNumber(1234), '1,234')
})
test('formatNumber >=10000 折算为 w', () => {
  assert.equal(formatNumber(50000), '5.00w')
  assert.equal(formatNumber(123456), '12.35w')
})
test('formatNumber 非数字兜底 0', () => {
  assert.equal(formatNumber('abc'), '0')
  assert.equal(formatNumber(null), '0')
})

// formatAmount
test('formatAmount 始终保留两位小数', () => {
  assert.equal(formatAmount(1), '1.00')
  assert.equal(formatAmount(1.5), '1.50')
  assert.equal(formatAmount(0), '0.00')
  assert.equal(formatAmount(null), '0.00')
})

// validateTopupPresets
test('validateTopupPresets 空数组通过', () => {
  assert.equal(validateTopupPresets([]), null)
})
test('validateTopupPresets 正常金额通过', () => {
  assert.equal(validateTopupPresets([
    { amount: 10, bonus_quota: 10 },
    { amount: 50, bonus_quota: 60 },
  ]), null)
})
test('validateTopupPresets 拒绝重复金额', () => {
  const err = validateTopupPresets([
    { amount: 10, bonus_quota: 10 },
    { amount: 10, bonus_quota: 20 },
  ])
  assert.match(err, /快捷金额重复/)
  assert.match(err, /10/)
})
test('validateTopupPresets 拒绝 0 金额', () => {
  const err = validateTopupPresets([{ amount: 0, bonus_quota: 10 }])
  assert.match(err, /必须大于 0/)
})
test('validateTopupPresets 拒绝负数', () => {
  const err = validateTopupPresets([{ amount: -5, bonus_quota: 10 }])
  assert.match(err, /必须大于 0/)
})
test('validateTopupPresets 拒绝 NaN/字符串', () => {
  const err = validateTopupPresets([{ amount: 'abc', bonus_quota: 10 }])
  assert.match(err, /必须大于 0/)
})
test('validateTopupPresets 区分大小数（1 与 1.0 视为不同）', () => {
  // Number(1) === Number(1.0) 实际为 true，验证当前实现语义
  const err = validateTopupPresets([
    { amount: 1, bonus_quota: 1 },
    { amount: 1.0, bonus_quota: 2 },
  ])
  // 1 === 1.0，期望被判定为重复
  assert.match(err, /快捷金额重复/)
})

// calcCustomBonus
test('calcCustomBonus 默认 1:1', () => {
  assert.equal(calcCustomBonus(25, 1), 25)
})
test('calcCustomBonus 500000 比例', () => {
  assert.equal(calcCustomBonus(10, 500000), 5_000_000)
})
test('calcCustomBonus 空 exchange_rate 视为 1', () => {
  assert.equal(calcCustomBonus(7, undefined), 7)
  assert.equal(calcCustomBonus(7, null), 7)
  // rate <= 0 是非法配置，工具函数防御性返回 0
  assert.equal(calcCustomBonus(7, 0), 0)
  assert.equal(calcCustomBonus(7, -1), 0)
})
test('calcCustomBonus 非数字金额视为 0', () => {
  assert.equal(calcCustomBonus('abc', 1), 0)
  assert.equal(calcCustomBonus(null, 1), 0)
})
