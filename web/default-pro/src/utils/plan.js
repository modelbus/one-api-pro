// plan 套餐相关纯函数（Node 单测友好，不依赖 Vue / Arco）
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode

// normalizeFeaturesRaw 把后端可能返回的 features 字段统一为字符串数组。
//
// 后端字段历史形态：
//   - 已是 array<string>          → 原样返回
//   - JSON 数组字符串 "[\"A\"]"   → JSON 解析
//   - 换行分隔纯文本 "A\nB"       → 按行拆分
//   - JSON 对象 {"A":true,"B":1}  → 取键名（值为 false/0/null 的键丢弃）
//   - 空 / null / 非法            → 返回 []
//
// 这是兼容旧数据的兜底；新 API 已统一返回 array<string>。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
export function normalizeFeaturesRaw(raw) {
  if (raw == null) return []
  if (Array.isArray(raw)) {
    return raw.map((s) => String(s ?? '').trim()).filter((s) => s.length > 0)
  }
  if (typeof raw !== 'string') return []
  const trimmed = raw.trim()
  if (!trimmed) return []
  // 主路径：JSON 数组
  try {
    const obj = JSON.parse(trimmed)
    if (Array.isArray(obj)) {
      return obj.map((s) => String(s ?? '').trim()).filter((s) => s.length > 0)
    }
    if (obj && typeof obj === 'object') {
      return Object.entries(obj)
        .filter(([, v]) => v === true || (v != null && v !== false && v !== 0))
        .map(([k]) => k.trim())
        .filter((k) => k.length > 0)
    }
    return []
  } catch {
    // 兜底：换行分隔纯文本
    return raw
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line.length > 0)
  }
}

// sanitizeFeaturesList 把表单中的 features 数组清洗为 API 接受的形态。
//
//   - 去掉空白条目
//   - 每项 trim
//   - 始终返回新数组（不修改入参）
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
export function sanitizeFeaturesList(list) {
  if (!Array.isArray(list)) return []
  return list
    .map((s) => (s == null ? '' : String(s)).trim())
    .filter((s) => s.length > 0)
}

// buildEmptyFeaturesForm 返回弹窗初始的空 features 表单（一行空字符串占位）。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
export function buildEmptyFeaturesForm() {
  return ['']
}

// featuresFromRecord 把从后端读到的 features 字段转成弹窗表单（数组，至少 1 行）。
// 用于编辑场景的回填；保证用户至少看到一个可输入框。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
export function featuresFromRecord(raw) {
  const list = normalizeFeaturesRaw(raw)
  return list.length > 0 ? list : buildEmptyFeaturesForm()
}