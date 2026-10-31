// title.test.mjs Tab 标题纯函数单元测试
// Unit tests for the browser tab title helpers.
// 版本: v0.0.27
// 日期: 2026-10-31
// 作者: opencode
// 运行方式: node web/default-pro/src/utils/title.test.mjs
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { SITE_NAME, buildDocumentTitle, resolvePageTitle } from './title.js'

// resolvePageTitle — 路由 meta.title 经 i18n 翻译为可见标题
test('resolvePageTitle: 命中翻译 → 返回页面标题', () => {
  const t = (key) => ({ 'route.token': '令牌', 'route.dashboard': '仪表盘' })[key] || key
  assert.equal(resolvePageTitle('route.token', t), '令牌')
  assert.equal(resolvePageTitle('route.dashboard', t), '仪表盘')
})

test('resolvePageTitle: 无 meta.title / 未翻译 / 空翻译 → 空串（走站点默认标题）', () => {
  const t = (key) => key
  assert.equal(resolvePageTitle(undefined, t), '')
  assert.equal(resolvePageTitle('', t), '')
  assert.equal(resolvePageTitle('route.unknownKey', t), '')
  assert.equal(resolvePageTitle('route.token', () => '   '), '')
  assert.equal(resolvePageTitle('route.token', undefined), '')
})

// buildDocumentTitle — 有页面标题时以页面标题为主，站点名固定不受后端 system_name 影响
test('buildDocumentTitle: 有页面标题 → 「页面—站点名」', () => {
  assert.equal(buildDocumentTitle('仪表盘', '企业级AI API 网关'), `仪表盘—${SITE_NAME}`)
  assert.equal(buildDocumentTitle('令牌', 'Enterprise AI API Gateway'), `令牌—${SITE_NAME}`)
})

test('buildDocumentTitle: 无页面标题 → 「站点名—站点默认标题」', () => {
  assert.equal(buildDocumentTitle('', '企业级AI API 网关'), `${SITE_NAME}—企业级AI API 网关`)
  assert.equal(buildDocumentTitle(undefined, '企业级AI API 网关'), `${SITE_NAME}—企业级AI API 网关`)
})

test('buildDocumentTitle: 空白裁剪与站点名回退', () => {
  assert.equal(buildDocumentTitle('  令牌  ', '默认'), `令牌—${SITE_NAME}`)
  assert.equal(buildDocumentTitle('令牌', '默认', '   '), `令牌—${SITE_NAME}`)
  assert.equal(buildDocumentTitle('', ''), SITE_NAME)
})

test('buildDocumentTitle: 结果只由路由标题与固定站点名组成', () => {
  // 后端 /api/status 的 system_name（默认可能是 one-api）不再参与标题拼装
  assert.equal(buildDocumentTitle('首页', '企业级AI API 网关'), `首页—${SITE_NAME}`)
  assert.equal(buildDocumentTitle('', '企业级AI API 网关'), `${SITE_NAME}—企业级AI API 网关`)
})
