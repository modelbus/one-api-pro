// plan_limits.test.mjs 套餐 model_limits 行式编辑纯函数的单元测试
// Unit tests for the plan model_limits row helpers
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
import { test } from 'node:test'
import assert from 'node:assert/strict'
import {
  BILLING_TOKEN,
  BILLING_REQUEST,
  DEFAULT_PERIOD_H,
  MAX_PERIOD_H,
  MIN_PERIOD_H,
  ROW_LIMIT_KEYS,
  createLimitRow,
  limitFieldOf,
  modelLimitsToRows,
  normalizeBillingType,
  normalizeLimit,
  normalizePeriodH,
  parseLimitsObject,
  rowsToModelLimits,
  validateLimitRows,
} from './plan_limits.js'

test('normalizeBillingType falls back to token', () => {
  assert.equal(normalizeBillingType(BILLING_REQUEST), BILLING_REQUEST)
  assert.equal(normalizeBillingType('nonsense'), BILLING_TOKEN)
  assert.equal(normalizeBillingType(undefined), BILLING_TOKEN)
})

test('limitFieldOf maps window + dimension to backend field', () => {
  assert.equal(limitFieldOf('week', BILLING_TOKEN), 'token_week')
  assert.equal(limitFieldOf('month', BILLING_REQUEST), 'request_month')
  assert.equal(limitFieldOf('bogus', BILLING_TOKEN), 'token_period')
})

test('normalizePeriodH clamps to 1-24 and falls back to 5', () => {
  assert.equal(normalizePeriodH(8), 8)
  assert.equal(normalizePeriodH(0), DEFAULT_PERIOD_H)
  assert.equal(normalizePeriodH(-3), DEFAULT_PERIOD_H)
  assert.equal(normalizePeriodH(99), MAX_PERIOD_H)
  assert.equal(normalizePeriodH('abc'), DEFAULT_PERIOD_H)
  assert.equal(MIN_PERIOD_H, 1)
})

test('normalizeLimit treats junk and negatives as 0', () => {
  assert.equal(normalizeLimit(12), 12)
  assert.equal(normalizeLimit(0), 0)
  assert.equal(normalizeLimit(-1), 0)
  assert.equal(normalizeLimit('x'), 0)
  assert.equal(normalizeLimit(undefined), 0)
})

test('createLimitRow returns one empty row with three windows', () => {
  // 限额默认 null：表单里保持空值以显示 placeholder，而不是预填 0
  assert.deepEqual(createLimitRow(), {
    model: '',
    periodH: DEFAULT_PERIOD_H,
    limitPeriod: null,
    limitWeek: null,
    limitMonth: null,
  })
  assert.deepEqual(createLimitRow('gpt-4o', { limitWeek: 20 }), {
    model: 'gpt-4o',
    periodH: DEFAULT_PERIOD_H,
    limitPeriod: null,
    limitWeek: 20,
    limitMonth: null,
  })
})

test('ROW_LIMIT_KEYS covers every window', () => {
  assert.deepEqual(ROW_LIMIT_KEYS, {
    period: 'limitPeriod',
    week: 'limitWeek',
    month: 'limitMonth',
  })
})

test('parseLimitsObject tolerates empty input and rejects junk', () => {
  assert.deepEqual(parseLimitsObject(''), { ok: true, limits: {}, reason: '' })
  assert.deepEqual(parseLimitsObject(null), { ok: true, limits: {}, reason: '' })
  assert.equal(parseLimitsObject('{oops}').reason, 'json')
  assert.equal(parseLimitsObject('[1,2]').reason, 'type')
  assert.equal(parseLimitsObject('{"gpt-4o":{"token_period":10}}').ok, true)
})

test('modelLimitsToRows keeps only the active billing dimension', () => {
  const raw = JSON.stringify({
    'gpt-4o': { period_h: 3, token_period: 100, token_month: 900, request_period: 7 },
  })
  const { ok, rows } = modelLimitsToRows(raw, BILLING_TOKEN)
  assert.equal(ok, true)
  assert.deepEqual(rows, [
    { model: 'gpt-4o', periodH: 3, limitPeriod: 100, limitWeek: null, limitMonth: 900 },
  ])
  const asRequest = modelLimitsToRows(raw, BILLING_REQUEST)
  assert.deepEqual(asRequest.rows, [
    { model: 'gpt-4o', periodH: 3, limitPeriod: 7, limitWeek: null, limitMonth: null },
  ])
})

test('modelLimitsToRows keeps one row per model even without matching limits', () => {
  const raw = JSON.stringify({ 'gpt-4o': { period_h: 5, request_period: 20 } })
  const { rows } = modelLimitsToRows(raw, BILLING_TOKEN)
  assert.deepEqual(rows, [
    { model: 'gpt-4o', periodH: 5, limitPeriod: null, limitWeek: null, limitMonth: null },
  ])
})

test('modelLimitsToRows falls back to the default window hours', () => {
  const raw = JSON.stringify({ 'gpt-4o': { token_week: 50 } })
  const { rows } = modelLimitsToRows(raw, BILLING_TOKEN)
  assert.equal(rows[0].periodH, DEFAULT_PERIOD_H)
})

test('modelLimitsToRows reports failure for malformed JSON', () => {
  const { ok, rows, reason } = modelLimitsToRows('{oops}', BILLING_TOKEN)
  assert.equal(ok, false)
  assert.equal(reason, 'json')
  assert.deepEqual(rows, [])
})

test('rowsToModelLimits writes every filled window', () => {
  const rows = [createLimitRow('gpt-4o', { periodH: 4, limitPeriod: 100, limitWeek: 200, limitMonth: 1000 })]
  assert.deepEqual(JSON.parse(rowsToModelLimits(rows, BILLING_TOKEN).json), {
    'gpt-4o': { period_h: 4, token_period: 100, token_week: 200, token_month: 1000 },
  })
  assert.deepEqual(JSON.parse(rowsToModelLimits(rows, BILLING_REQUEST).json), {
    'gpt-4o': { period_h: 4, request_period: 100, request_week: 200, request_month: 1000 },
  })
})

test('rowsToModelLimits skips blank models and rows without any limit', () => {
  const rows = [
    createLimitRow('', { limitPeriod: 10 }),
    createLimitRow('gpt-4o'),
    createLimitRow('gpt-4o-mini', { limitWeek: 30 }),
  ]
  assert.deepEqual(JSON.parse(rowsToModelLimits(rows, BILLING_TOKEN).json), {
    'gpt-4o-mini': { period_h: DEFAULT_PERIOD_H, token_week: 30 },
  })
})

test('rowsToModelLimits round-trips through modelLimitsToRows', () => {
  const rows = [
    createLimitRow('gpt-4o', { periodH: 5, limitPeriod: 100 }),
    createLimitRow('gpt-4o-mini', { periodH: 24, limitMonth: 5000 }),
  ]
  const { json } = rowsToModelLimits(rows, BILLING_REQUEST)
  const back = modelLimitsToRows(json, BILLING_REQUEST)
  assert.deepEqual(back.rows, rows)
})

test('validateLimitRows rejects the empty list', () => {
  assert.equal(validateLimitRows([]).code, 'errLimitsEmpty')
  assert.equal(validateLimitRows(undefined).code, 'errLimitsEmpty')
})

test('validateLimitRows rejects rows without a model', () => {
  assert.equal(validateLimitRows([createLimitRow('  ', { limitPeriod: 10 })]).code, 'errLimitModel')
})

test('validateLimitRows rejects a row with no window limit at all', () => {
  const err = validateLimitRows([createLimitRow('gpt-4o')])
  assert.equal(err.code, 'errLimitValue')
  assert.deepEqual(err.params, { n: 1 })
})

test('validateLimitRows accepts a row filling only one window', () => {
  assert.equal(validateLimitRows([createLimitRow('gpt-4o', { limitMonth: 10 })]), null)
})

test('validateLimitRows rejects duplicate models', () => {
  const rows = [
    createLimitRow('gpt-4o', { limitWeek: 10 }),
    createLimitRow('gpt-4o', { limitMonth: 20 }),
  ]
  const err = validateLimitRows(rows)
  assert.equal(err.code, 'errLimitDuplicate')
  assert.deepEqual(err.params, { n: 2, model: 'gpt-4o' })
})
