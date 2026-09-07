// plan 套餐工具函数单测（Node 内置 test runner）
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode

import { test } from 'node:test'
import assert from 'node:assert/strict'
import {
  normalizeFeaturesRaw,
  sanitizeFeaturesList,
  buildEmptyFeaturesForm,
  featuresFromRecord,
} from './plan.js'

test('normalizeFeaturesRaw: null/undefined → []', () => {
  assert.deepEqual(normalizeFeaturesRaw(null), [])
  assert.deepEqual(normalizeFeaturesRaw(undefined), [])
  assert.deepEqual(normalizeFeaturesRaw(''), [])
  assert.deepEqual(normalizeFeaturesRaw('   '), [])
})

test('normalizeFeaturesRaw: 已是数组 → 原样返回（去空白、去空）', () => {
  assert.deepEqual(normalizeFeaturesRaw(['A', 'B', 'C']), ['A', 'B', 'C'])
  assert.deepEqual(normalizeFeaturesRaw(['A', '', '  B  ']), ['A', 'B'])
  assert.deepEqual(normalizeFeaturesRaw([null, undefined, 'X']), ['X'])
})

test('normalizeFeaturesRaw: JSON 数组字符串 → 数组', () => {
  assert.deepEqual(normalizeFeaturesRaw('["A","B","C"]'), ['A', 'B', 'C'])
  assert.deepEqual(normalizeFeaturesRaw('["A", "B"]'), ['A', 'B'])
  assert.deepEqual(normalizeFeaturesRaw('[]'), [])
})

test('normalizeFeaturesRaw: 换行分隔纯文本 → 数组', () => {
  assert.deepEqual(normalizeFeaturesRaw('A\nB\nC'), ['A', 'B', 'C'])
  assert.deepEqual(normalizeFeaturesRaw('A\r\nB\r\nC'), ['A', 'B', 'C'])
  assert.deepEqual(normalizeFeaturesRaw('A\n\n  C  '), ['A', 'C'])
})

test('normalizeFeaturesRaw: JSON 对象 → 仅保留 truthy 键', () => {
  assert.deepEqual(
    normalizeFeaturesRaw('{"A":true,"B":false,"C":1,"D":0,"E":null}'),
    ['A', 'C']
  )
})

test('normalizeFeaturesRaw: 非法 JSON → 走换行兜底', () => {
  assert.deepEqual(normalizeFeaturesRaw('not-json'), ['not-json'])
  assert.deepEqual(normalizeFeaturesRaw('A;B;C'), ['A;B;C'])
})

test('normalizeFeaturesRaw: 非字符串非数组 → []', () => {
  assert.deepEqual(normalizeFeaturesRaw(123), [])
  assert.deepEqual(normalizeFeaturesRaw({}), [])
})

test('sanitizeFeaturesList: 过滤空白', () => {
  assert.deepEqual(sanitizeFeaturesList(['A', '', '  B  ', 'C']), ['A', 'B', 'C'])
  assert.deepEqual(sanitizeFeaturesList([null, undefined, '']), [])
})

test('sanitizeFeaturesList: 数字转字符串保留（值即文本）', () => {
  assert.deepEqual(sanitizeFeaturesList([0, 100, 'X']), ['0', '100', 'X'])
})

test('sanitizeFeaturesList: 非数组 → []', () => {
  assert.deepEqual(sanitizeFeaturesList(null), [])
  assert.deepEqual(sanitizeFeaturesList('A\nB'), [])
})

test('sanitizeFeaturesList: 返回新数组不修改入参', () => {
  const input = ['A', '', 'C']
  const result = sanitizeFeaturesList(input)
  assert.notEqual(result, input)
  assert.deepEqual(input, ['A', '', 'C'])
})

test('buildEmptyFeaturesForm: 始终返回包含一个空串的数组', () => {
  const form = buildEmptyFeaturesForm()
  assert.equal(form.length, 1)
  assert.equal(form[0], '')
})

test('featuresFromRecord: 空记录 → 空表单', () => {
  assert.deepEqual(featuresFromRecord(null), [''])
  assert.deepEqual(featuresFromRecord('[]'), [''])
  assert.deepEqual(featuresFromRecord(''), [''])
})

test('featuresFromRecord: 有值 → 原样回填', () => {
  assert.deepEqual(featuresFromRecord(['A', 'B']), ['A', 'B'])
  assert.deepEqual(featuresFromRecord('["A","B"]'), ['A', 'B'])
  assert.deepEqual(featuresFromRecord('A\nB\nC'), ['A', 'B', 'C'])
})