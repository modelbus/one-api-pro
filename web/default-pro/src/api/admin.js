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
  // 运营概览：纯 KPI（不含图表数据）。
  // Operations overview: KPI-only (chart data lives in charts()).
  // range: today | 7d | 30d | all（默认 7d）
  overview: (range) => api.get('/api/admin/dashboard/overview', { params: buildParams({ range }) }),

  // 全站图表数据：day × model 聚合，结构与 /api/user/dashboard 一致。
  // 一份数据同时驱动：请求量/额度/Token 折线图 + 模型分布 + 使用明细。
  // Whole-site chart data: day×model aggregate, same shape as /api/user/dashboard.
  // One payload drives the 3 line charts + model distribution + usage details.
  // range: today | 7d | 30d | all（默认 7d）
  charts: (range) => api.get('/api/admin/dashboard/charts', { params: buildParams({ range }) }),

  // 活跃用户排行榜。
  // Active-user leaderboard.
  // range: today | 7d | 30d | all（默认 7d）；limit 默认 20
  topUsers: (range, limit) =>
    api.get('/api/admin/dashboard/top-users', { params: buildParams({ range, limit }) }),
}

export default adminApi
