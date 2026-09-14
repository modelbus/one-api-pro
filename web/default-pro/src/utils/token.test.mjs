// token.test.mjs 令牌提交/展示工具函数单元测试
// Unit tests for token submission / rendering helpers.
//
// 版本: v0.0.21
// 日期: 2026-09-14
// 作者: opencode
// 运行方式: node web/default-pro/src/utils/token.test.mjs
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { buildTokenExpiredTime, formatExpiredTime } from './token.js'

// buildTokenExpiredTime — 提交 sentinel 必须与后端约定一致
test('buildTokenExpiredTime: never_expire=true → -1', () => {
  assert.equal(
    buildTokenExpiredTime({ never_expire: true, expired_time: null }),
    -1,
  )
  assert.equal(
    buildTokenExpiredTime({ never_expire: true, expired_time: 1700000000000 }),
    -1,
    'never_expire 优先于过期时间字段',
  )
})

test('buildTokenExpiredTime: expired_time 为空/null/0 → -1', () => {
  assert.equal(buildTokenExpiredTime({ never_expire: false, expired_time: null }), -1)
  assert.equal(buildTokenExpiredTime({ never_expire: false, expired_time: 0 }), -1)
  assert.equal(buildTokenExpiredTime({ never_expire: false, expired_time: undefined }), -1)
})

test('buildTokenExpiredTime: never_expire=false 且有日期 → 折算秒', () => {
  const ms = 1700000000000
  assert.equal(
    buildTokenExpiredTime({ never_expire: false, expired_time: ms }),
    Math.floor(ms / 1000),
  )
})

test('buildTokenExpiredTime: form 为空对象 → -1', () => {
  assert.equal(buildTokenExpiredTime({}), -1)
  assert.equal(buildTokenExpiredTime(null), -1)
})

// 回归测试：旧前端 bug 把"永不过期"写成 0，导致后端报"令牌已过期"。
// 这次修复后必须永远不再提交 0。
test('buildTokenExpiredTime 回归：never_expire=true 永不返回 0', () => {
  for (const exp of [null, undefined, 0, '0', '', 1700000000000]) {
    assert.notEqual(
      buildTokenExpiredTime({ never_expire: true, expired_time: exp }),
      0,
      `输入 expired_time=${exp} 时不应返回 0`,
    );
  }
})

// formatExpiredTime — 展示与 sentinel 一致
test('formatExpiredTime: null/undefined/0/负数 → 永不过期', () => {
  assert.equal(formatExpiredTime(null), '永不过期')
  assert.equal(formatExpiredTime(undefined), '永不过期')
  assert.equal(formatExpiredTime(0), '永不过期')
  assert.equal(formatExpiredTime(-1), '永不过期')
})

test('formatExpiredTime: 合法 unix 时间戳 → 格式化日期', () => {
  // 2024-01-15 12:34:56 UTC
  // 1705322096 seconds = 2024-01-15T12:34:56Z; rendered in local time
  const out = formatExpiredTime(1705322096)
  assert.match(out, /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/)
})

test('formatExpiredTime: 非数字字符串 → 永不过期', () => {
  assert.equal(formatExpiredTime('abc'), '永不过期')
})