// quota.test.mjs quota.js 的 Node 内置测试
// Node built-in tests for quota.js.
// 版本: v0.0.23
// 日期: 2026-10-03
// 作者: opencode
import { test } from 'node:test'
import assert from 'node:assert/strict'
import {
  DEFAULT_QUOTA_PER_UNIT,
  resolveQuotaPerUnit,
  quotaToYuan,
  formatYuan,
  formatQuota,
} from './quota.js'

test('resolveQuotaPerUnit falls back on invalid input', () => {
  assert.equal(resolveQuotaPerUnit(0), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit(-1), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit(undefined), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit('abc'), DEFAULT_QUOTA_PER_UNIT)
  assert.equal(resolveQuotaPerUnit(1000000), 1000000)
  assert.equal(resolveQuotaPerUnit('1000000'), 1000000)
})

test('quotaToYuan uses 1 CNY = quota_per_unit', () => {
  assert.equal(quotaToYuan(5000000, 500000), 10)
  assert.equal(quotaToYuan(1250, 500000), 0.0025)
})

test('formatYuan adapts decimals to magnitude', () => {
  assert.equal(formatYuan(10), '¥10.00')
  assert.equal(formatYuan(0.0025), '¥0.002500')
  assert.equal(formatYuan(1234.5), '¥1,234.50')
})

test('formatQuota renders CNY and handles empty/invalid', () => {
  assert.equal(formatQuota(5000000, 500000), '¥10.00')
  assert.equal(formatQuota(0, 500000), '¥0.000000')
  assert.equal(formatQuota(null, 500000), '-')
  assert.equal(formatQuota('', 500000), '-')
})
