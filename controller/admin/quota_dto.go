// quota_dto.go 管理看板 API 边界的额度单位换算（微元 <-> 元）
// API-boundary quota unit conversion (micro-quota <-> CNY yuan) for the admin dashboard
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode
package admin

import (
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/model"
)

// 口径约定同 controller/quota_dto.go：存储层（DB）一律微元整数，API 边界一律「元」浮点。
//
// admin 是独立包，无法复用 controller 包内的私有换算函数；此处各自实现一份
// （仅一个除法），以避免为此引入跨包依赖或新增公共包。
//
// Same convention as controller/quota_dto.go: micro-quota in storage, CNY yuan at the
// API boundary. This package is standalone and cannot reuse controller's private helpers,
// so it carries its own one-line conversion instead of adding a cross-package dependency.

// quotaToYuan 将存储层的微元额度换算为「元」，用于 API 输出。
// Convert storage-layer micro-quota to CNY yuan for API responses.
func quotaToYuan(q int64) float64 {
	return float64(q) / config.QuotaPerUnit
}

// adminQuotaDTO 管理看板的额度聚合输出：全部以「元」为单位。
// Admin dashboard quota aggregate response, all expressed in CNY yuan.
type adminQuotaDTO struct {
	Today float64 `json:"today"`
	Week  float64 `json:"week"`
	Month float64 `json:"month"`
	Total float64 `json:"total"`
}

// adminOverviewDTO 管理看板概览输出：quota 聚合以「元」为单位，其余字段透传。
// Admin dashboard overview response: the quota aggregate is expressed in CNY yuan;
// every other field passes through unchanged.
type adminOverviewDTO struct {
	*model.AdminDashboardOverview
	Quota adminQuotaDTO `json:"quota"`
}

// toAdminOverviewDTO 把概览转换为「元」口径的输出结构。
// Convert the overview to the yuan-denominated response shape.
func toAdminOverviewDTO(o *model.AdminDashboardOverview) *adminOverviewDTO {
	if o == nil {
		return nil
	}
	return &adminOverviewDTO{
		AdminDashboardOverview: o,
		Quota: adminQuotaDTO{
			Today: quotaToYuan(o.Quota.Today),
			Week:  quotaToYuan(o.Quota.Week),
			Month: quotaToYuan(o.Quota.Month),
			Total: quotaToYuan(o.Quota.Total),
		},
	}
}

// adminTopUserDTO 管理看板排行榜输出：quota（消费）/ balance（余额）以「元」为单位。
// Admin dashboard leaderboard response: quota (consumption) / balance in CNY yuan.
type adminTopUserDTO struct {
	*model.AdminTopUserRow
	Quota   float64 `json:"quota"`
	Balance float64 `json:"balance"`
}

// toAdminTopUserDTO 把排行榜单行转换为「元」口径的输出结构。
// Convert a single leaderboard row to the yuan-denominated response shape.
func toAdminTopUserDTO(r *model.AdminTopUserRow) *adminTopUserDTO {
	return &adminTopUserDTO{
		AdminTopUserRow: r,
		Quota:           quotaToYuan(r.Quota),
		Balance:         quotaToYuan(r.Balance),
	}
}

// toAdminTopUserDTOs 批量转换排行榜；入参为空时返回空切片而非 nil。
// Convert the leaderboard; returns an empty slice (not nil) when the input is empty.
func toAdminTopUserDTOs(rows []model.AdminTopUserRow) []*adminTopUserDTO {
	out := make([]*adminTopUserDTO, 0, len(rows))
	for i := range rows {
		out = append(out, toAdminTopUserDTO(&rows[i]))
	}
	return out
}

// logStatisticDTO 按 day × model 聚合的日志统计输出：Quota 以「元」为单位。
// Log-statistic (grouped by day × model) response: Quota expressed in CNY yuan.
//
// 字段名保持 model.LogStatistic 原始的大写形式（该结构未声明 json tag），
// 与 /api/user/dashboard 的结构保持一致。
// Field names keep model.LogStatistic's original capitalized form (no json tags declared),
// matching the /api/user/dashboard shape.
type logStatisticDTO struct {
	*model.LogStatistic
	Quota float64 `json:"Quota"`
}

// toLogStatisticDTOs 批量转换聚合统计；入参为 nil 时返回空切片而非 nil。
// Convert aggregated statistics; returns an empty slice (not nil) when the input is empty.
func toLogStatisticDTOs(stats []*model.LogStatistic) []*logStatisticDTO {
	out := make([]*logStatisticDTO, 0, len(stats))
	for _, s := range stats {
		if s == nil {
			continue
		}
		out = append(out, &logStatisticDTO{
			LogStatistic: s,
			Quota:        quotaToYuan(int64(s.Quota)),
		})
	}
	return out
}
