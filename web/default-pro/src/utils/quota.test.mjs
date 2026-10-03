// quota.test.mjs quota.js 的 Node 内置测试
// Node built-in tests for quota.js.
// 版本: v0.0.24
// 日期: 2026-10-04
// 作者: opencode
import { test } from 'node:test'
import assert from 'node:assert/strict'
import {
  DEFAULT_QUOTA_PER_UNIT,
  resolveQuotaPerUnit,
  quotaToYuan,
  quotaToYuanExact,
  yuanInputPrecision,
  yuanToQuota,
  formatYuan,
  formatQuota,
} from './quota.js'

test('DEFAULT_QUOTA_PER_UNIT is 1_000_000 (1 CNY = 1e6 quota)', () => {
  assert.equal(DEFAULT_QUOTA_PER_UNIT, 1_000_000)
})

test('resolveQuotaPerUnit falls back on invalid input', () => {
  assert.equal(resolveQuotaPerUnit(0), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit(-1), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit(undefined), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit('abc'), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit(1000000), 1000000)
  assert.equal(resolveQuotaPerUnit('1000000'), 1000000)
})

test('quotaToYuan uses 1 CNY = quota_per_unit', () => {
  assert.equal(quotaToYuan(10000000, 1000000), 10)
  assert.equal(quotaToYuan(2500, 1000000), 0.0025)
})

test('yuanToQuota converts and rounds float error', () => {
  assert.equal(yuanToQuota(10, 1000000), 10000000)
  // 0.07 × 1e6 在浮点下为 69999.99…，需四舍五入为 70000
  assert.equal(yuanToQuota(0.07, 1000000), 70000)
  assert.equal(yuanToQuota(0, 1000000), 0)
  assert.equal(yuanToQuota(-1, 1000000), 0)
  assert.equal(yuanToQuota('abc', 1000000), 0)
})

test('quotaToYuanExact keeps sub-cent quota lossless', () => {
  // 500 额度 = ¥0.0005，旧实现 toFixed(2) 会截断为 0，导致保存时静默清零
  assert.equal(quotaToYuanExact(500, 1000000), 0.0005)
  assert.equal(quotaToYuanExact(1000, 1000000), 0.001)
  assert.equal(quotaToYuanExact(0, 1000000), 0)
  assert.equal(quotaToYuanExact(10000000, 1000000), 10)
  // 回填后再换算必须与原额度一致（无损往返）
  assert.equal(yuanToQuota(quotaToYuanExact(500, 1000000), 1000000), 500)
  assert.equal(yuanToQuota(quotaToYuanExact(123456, 1000000), 1000000), 123456)
})

test('yuanInputPrecision uses 2 decimals for whole cents, 6 otherwise', () => {
  assert.equal(yuanInputPrecision(10, 1000000), 2)
  assert.equal(yuanInputPrecision(0.07, 1000000), 2)
  assert.equal(yuanInputPrecision(0, 1000000), 2)
  assert.equal(yuanInputPrecision(0.0005, 1000000), 6)
  assert.equal(yuanInputPrecision(0.001, 1000000), 6)
  assert.equal(yuanInputPrecision(10.000001, 1000000), 6)
})

test('formatYuan adapts decimals to magnitude', () => {
  assert.equal(formatYuan(10), '¥10.00')
  assert.equal(formatYuan(0.0025), '¥0.002500')
  assert.equal(formatYuan(1234.5), '¥1,234.50')
})

test('formatQuota renders CNY and handles empty/invalid', () => {
  assert.equal(formatQuota(10000000, 1000000), '¥10.00')
  assert.equal(formatQuota(0, 1000000), '¥0.000000')
  assert.equal(formatQuota(null, 1000000), '-')
  assert.equal(formatQuota('', 1000000), '-')
})
