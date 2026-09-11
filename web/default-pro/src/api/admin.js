// admin 管理员运营接口聚合模块
// Admin operations API module (dashboard overview + leaderboard)
// 版本: v0.0.16
// 日期: 2026-09-11
import api from './index'

// 仅追加有意义的过滤参数，避免 ?range= 这种空值污染 query。
// Only forward truthy params to keep query strings clean.
function buildParams(extras = {}) {
  const out = {}
  for (const [k, v] of Object.entries(extras)) {
    if (v === 0 || (v !== '' && v !== null && v !== undefined)) {
      out[k] = v
    }
  }
  return out
}

export const adminApi = {
  // 运营概览：KPI + 7 日趋势，一次性返回。
  // Operations overview: KPIs + 7-day trends, single-shot response.
  // range: today | 7d | 30d | all（默认 7d）
  overview: (range) => api.get('/api/admin/dashboard/overview', { params: buildParams({ range }) }),

  // 活跃用户排行榜。
  // Active-user leaderboard.
  // range: today | 7d | 30d | all（默认 7d）；limit 默认 20
  topUsers: (range, limit) =>
    api.get('/api/admin/dashboard/top-users', { params: buildParams({ range, limit }) }),
}

export default adminApi
