// 充值工具函数单元测试
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
// 运行方式: node web/default-pro/src/utils/topup.test.mjs
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { formatNumber, formatAmount, validateTopupPresets } from './topup.js'

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
  assert.equal(validateTopupPresets([], null, 1000000), null)
})
test('validateTopupPresets 1:1 金额通过', () => {
  assert.equal(validateTopupPresets([
    { amount: 10, bonus_quota: 10000000 },
    { amount: 50, bonus_quota: 50000000 },
  ], null, 1000000), null)
})
test('validateTopupPresets 允许赠送（充 10 得 15）', () => {
  assert.equal(validateTopupPresets([
    { amount: 10, bonus_quota: 15000000 },
  ], null, 1000000), null)
})
test('validateTopupPresets 拒绝到账低于支付', () => {
  const err = validateTopupPresets([
    { amount: 10, bonus_quota: 9999999 },
  ], null, 1000000)
  assert.match(err, /到账金额不能低于支付金额/)
})
test('validateTopupPresets 拒绝重复金额', () => {
  const err = validateTopupPresets([
    { amount: 10, bonus_quota: 10000000 },
    { amount: 10, bonus_quota: 20000000 },
  ], null, 1000000)
  assert.match(err, /快捷金额重复/)
  assert.match(err, /10/)
})
test('validateTopupPresets 拒绝 0 金额', () => {
  const err = validateTopupPresets([{ amount: 0, bonus_quota: 10000000 }], null, 1000000)
  assert.match(err, /必须大于 0/)
})
test('validateTopupPresets 拒绝负数', () => {
  const err = validateTopupPresets([{ amount: -5, bonus_quota: 10000000 }], null, 1000000)
  assert.match(err, /必须大于 0/)
})
test('validateTopupPresets 拒绝 NaN/字符串', () => {
  const err = validateTopupPresets([{ amount: 'abc', bonus_quota: 10000000 }], null, 1000000)
  assert.match(err, /必须大于 0/)
})
test('validateTopupPresets 区分大小数（1 与 1.0 视为相同）', () => {
  const err = validateTopupPresets([
    { amount: 1, bonus_quota: 1000000 },
    { amount: 1.0, bonus_quota: 2000000 },
  ], null, 1000000)
  assert.match(err, /快捷金额重复/)
})
