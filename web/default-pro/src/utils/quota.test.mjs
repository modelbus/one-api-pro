// quota.test.mjs quota.js 的 Node 内置测试
// Node built-in tests for quota.js.
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode
//
// 自 v0.0.25 起，后端在所有涉及金额的 API 上统一返回「元」，
// 前端只做格式化，不再有任何 1e6 换算，因此这里只覆盖显示相关函数。
import { test } from 'node:test'
import assert from 'node:assert/strict'
import {
  DEFAULT_TOKEN_QUOTA_YUAN,
  YUAN_INPUT_PRECISION,
  YUAN_DEFAULT_PRECISION,
  yuanInputPrecision,
  formatYuan,
  formatQuota,
} from './quota.js'

test('DEFAULT_TOKEN_QUOTA_YUAN is 1 CNY (same value as the old 1e6 micro-quota default)', () => {
  assert.equal(DEFAULT_TOKEN_QUOTA_YUAN, 1)
})

test('precision constants stay 6 (max) and 2 (default)', () => {
  assert.equal(YUAN_INPUT_PRECISION, 6)
  assert.equal(YUAN_DEFAULT_PRECISION, 2)
})

test('yuanInputPrecision uses 2 decimals for whole cents, 6 otherwise', () => {
  assert.equal(yuanInputPrecision(10), 2)
  assert.equal(yuanInputPrecision(0.07), 2)
  assert.equal(yuanInputPrecision(0), 2)
  assert.equal(yuanInputPrecision(0.0005), 6)
  assert.equal(yuanInputPrecision(0.001), 6)
  assert.equal(yuanInputPrecision(10.000001), 6)
  // 非法输入回退到常规精度
  assert.equal(yuanInputPrecision('abc'), 2)
  assert.equal(yuanInputPrecision(null), 2)
})

test('formatYuan adapts decimals to magnitude', () => {
  assert.equal(formatYuan(10), '¥10.00')
  assert.equal(formatYuan(0.0025), '¥0.002500')
  assert.equal(formatYuan(1234.5), '¥1,234.50')
  // 负余额（无限额度令牌的哨兵值）也按量级展示
  assert.equal(formatYuan(-0.000001), '¥-0.000001')
})

test('formatQuota renders already-yuan amounts and handles empty/invalid', () => {
  // 入参已是「元」：10 → ¥10.00，不再做任何换算
  assert.equal(formatQuota(10), '¥10.00')
  assert.equal(formatQuota(0), '¥0.000000')
  assert.equal(formatQuota(0.0005), '¥0.000500')
  assert.equal(formatQuota(null), '-')
  assert.equal(formatQuota(undefined), '-')
  assert.equal(formatQuota(''), '-')
  assert.equal(formatQuota('abc'), 'abc')
})
