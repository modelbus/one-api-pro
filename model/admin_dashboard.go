// admin_dashboard.go 管理员仪表盘聚合查询
// Admin dashboard aggregate queries for the operations overview.
// 版本: v0.0.16
// 日期: 2026-09-11
// 作者: opencode
//
// 提供「全站 KPI 概览」「日级趋势」「活跃用户排行榜」三类聚合。
// All SQL filters use status/status predicates to exclude soft-deleted rows.

package model

import (
	"fmt"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
)

// AdminDashboardUsers 用户层 KPI
// AdminDashboardUsers holds user-level KPI aggregates.
// 版本: v0.0.16
// 日期: 2026-09-11
type AdminDashboardUsers struct {
	Total    int64 `json:"total" gorm:"column:total"`
	Enabled  int64 `json:"enabled" gorm:"column:enabled"`
	Disabled int64 `json:"disabled" gorm:"column:disabled"`
	Deleted  int64 `json:"deleted" gorm:"column:deleted"`
	NewToday int64 `json:"new_today" gorm:"column:new_today"`
	New7d    int64 `json:"new_7d" gorm:"column:new_7d"`
	New30d   int64 `json:"new_30d" gorm:"column:new_30d"`
	Active7d int64 `json:"active_7d" gorm:"column:active_7d"`
}

// AdminDashboardResources 资源计数（令牌 / 渠道 / 套餐 / 兑换码 / 订阅）
// AdminDashboardResources holds resource counts.
// 版本: v0.0.16
// 日期: 2026-09-11
type AdminDashboardResources struct {
	Total   int64 `json:"total" gorm:"column:total"`
	Enabled int64 `json:"enabled,omitempty" gorm:"column:enabled"`
	Used    int64 `json:"used,omitempty" gorm:"column:used"`
	Unused  int64 `json:"unused,omitempty" gorm:"column:unused"`
	Active  int64 `json:"active,omitempty" gorm:"column:active"`
	Expired int64 `json:"expired,omitempty" gorm:"column:expired"`
}

// AdminDashboardQuota 各时间窗的 quota 消耗
// AdminDashboardQuota holds quota consumption for each time window.
// 版本: v0.0.16
// 日期: 2026-09-11
type AdminDashboardQuota struct {
	Today int64 `json:"today"`
	Week  int64 `json:"week"`
	Month int64 `json:"month"`
	Total int64 `json:"total"`
}

// AdminDashboardRevenue 收入 KPI（金额单位：元）
// AdminDashboardRevenue holds revenue KPIs (CNY).
// 版本: v0.0.16
// 日期: 2026-09-11
type AdminDashboardRevenue struct {
	Total        float64 `json:"total"`
	Topup        float64 `json:"topup"`
	Subscription float64 `json:"subscription"`
	Refund       float64 `json:"refund"`
}

// AdminDashboardOverview 单次响应聚合：纯 KPI。
// 图表数据由 /api/admin/dashboard/charts 独立提供（结构对齐 api/user/dashboard）。
// AdminDashboardOverview is the KPI-only aggregate response.
// Chart data lives in /api/admin/dashboard/charts (same shape as api/user/dashboard).
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminDashboardOverview struct {
	Users         AdminDashboardUsers     `json:"users"`
	Tokens        AdminDashboardResources `json:"tokens"`
	Channels      AdminDashboardResources `json:"channels"`
	Plans         AdminDashboardResources `json:"plans"`
	Redemptions   AdminDashboardResources `json:"redemptions"`
	Subscriptions AdminDashboardResources `json:"subscriptions"`
	Quota         AdminDashboardQuota     `json:"quota"`
	Revenue       AdminDashboardRevenue   `json:"revenue"`
	Range         string                  `json:"range"`
	GeneratedAt   int64                   `json:"generated_at"`
}

// AdminTopUserRow 排行榜单行
// AdminTopUserRow is one row in the active-user leaderboard.
// 版本: v0.0.16
// 日期: v0.0.16
type AdminTopUserRow struct {
	Id              int    `json:"id" gorm:"column:id"`
	Username        string `json:"username" gorm:"column:username"`
	Email           string `json:"email" gorm:"column:email"`
	RequestCount    int64  `json:"request_count" gorm:"column:request_count_col"`
	Quota           int64  `json:"quota" gorm:"column:quota_col"`
	Balance         int64  `json:"balance" gorm:"column:quota"`
	CurrentPlanName string `json:"current_plan_name"`
}

// AdminDashboardRange 解析后的时间筛选参数
// AdminDashboardRange is the parsed range query parameter.
// 版本: v0.0.16
// 日期: 2026-09-11
type AdminDashboardRange struct {
	Key       string
	StartTs   int64
	EndTs     int64
	TrendDays int
}

// ParseAdminDashboardRange 把 ?range=today|7d|30d|all 解析为窗口参数。
// all 表示全量，StartTs=0 / EndTs=0（与 model/log.go 的 0=无下界约定一致）。
// ParseAdminDashboardRange maps ?range= to start/end timestamps.
// "all" → StartTs=0 / EndTs=0 (meaning "no bound", see model/log.go).
// 版本: v0.0.16
// 日期: 2026-09-11
func ParseAdminDashboardRange(raw string) AdminDashboardRange {
	now := helper.GetTimestamp()
	daySec := int64(86400)
	switch raw {
	case "today":
		startOfDay := now - (now % daySec)
		return AdminDashboardRange{Key: "today", StartTs: startOfDay, EndTs: now, TrendDays: 1}
	case "30d":
		return AdminDashboardRange{Key: "30d", StartTs: now - 30*daySec, EndTs: now, TrendDays: 30}
	case "all":
		return AdminDashboardRange{Key: "all", StartTs: 0, EndTs: 0, TrendDays: 30}
	default:
		return AdminDashboardRange{Key: "7d", StartTs: now - 7*daySec, EndTs: now, TrendDays: 7}
	}
}

// resolveRangePreset 取小写预设名；与 ParseAdminDashboardRange 兼容。
// resolveRangePreset returns the canonical range key (lowercased, default 7d).
// 版本: v0.0.16
// 日期: 2026-09-11
func resolveRangePreset(raw string) string {
	switch raw {
	case "today", "7d", "30d", "all":
		return raw
	}
	return "7d"
}

// GetAdminDashboardOverview 一次性聚合所有 KPI（不含图表数据）。
// 各子查询均为单条 SQL，避免 N+1。
// GetAdminDashboardOverview aggregates all KPIs (excluding chart data).
// Each subsection is one SQL to avoid N+1.
// 版本: v0.0.17
// 日期: 2026-09-11
func GetAdminDashboardOverview(rawRange string) (*AdminDashboardOverview, error) {
	preset := resolveRangePreset(rawRange)
	r := ParseAdminDashboardRange(rawRange)
	now := helper.GetTimestamp()
	daySec := int64(86400)
	todayStart := now - (now % daySec)
	weekStart := todayStart - 7*daySec
	monthStart := todayStart - 30*daySec

	overview := &AdminDashboardOverview{
		Range:       preset,
		GeneratedAt: now,
	}

	// --- Users ---
	if err := DB.Raw(`
		SELECT
		  COUNT(*) AS total,
		  COALESCE(SUM(status = ?), 0) AS enabled,
		  COALESCE(SUM(status = ?), 0) AS disabled,
		  COALESCE(SUM(status = ?), 0) AS deleted,
		  COALESCE(SUM(created_at >= ?), 0) AS new_today,
		  COALESCE(SUM(created_at >= ?), 0) AS new_7d,
		  COALESCE(SUM(created_at >= ?), 0) AS new_30d
		FROM users
	`, UserStatusEnabled, UserStatusDisabled, UserStatusDeleted,
		todayStart, weekStart, monthStart).Scan(&overview.Users).Error; err != nil {
		return nil, fmt.Errorf("aggregate users: %w", err)
	}

	// --- Tokens ---
	if err := DB.Raw(`
		SELECT COUNT(*) AS total,
		       COALESCE(SUM(status = ?), 0) AS enabled
		FROM tokens
	`, TokenStatusEnabled).Scan(&overview.Tokens).Error; err != nil {
		return nil, fmt.Errorf("aggregate tokens: %w", err)
	}

	// --- Channels ---
	if err := DB.Raw(`
		SELECT COUNT(*) AS total,
		       COALESCE(SUM(status = ?), 0) AS enabled
		FROM channels
	`, ChannelStatusEnabled).Scan(&overview.Channels).Error; err != nil {
		return nil, fmt.Errorf("aggregate channels: %w", err)
	}

	// --- Plans ---
	if err := DB.Raw(`
		SELECT COUNT(*) AS total,
		       COALESCE(SUM(status = ?), 0) AS enabled
		FROM plans
	`, PlanStatusEnabled).Scan(&overview.Plans).Error; err != nil {
		return nil, fmt.Errorf("aggregate plans: %w", err)
	}

	// --- Redemptions ---
	if err := DB.Raw(`
		SELECT COUNT(*) AS total,
		       COALESCE(SUM(status = ?), 0) AS used,
		       COALESCE(SUM(status = ?), 0) AS unused
		FROM redemptions
	`, RedemptionCodeStatusUsed, RedemptionCodeStatusEnabled).Scan(&overview.Redemptions).Error; err != nil {
		return nil, fmt.Errorf("aggregate redemptions: %w", err)
	}

	// --- Subscriptions (user_plans) ---
	if err := DB.Raw(`
		SELECT COUNT(*) AS total,
		       COALESCE(SUM(status = ?), 0) AS active,
		       COALESCE(SUM(status = ?), 0) AS expired
		FROM user_plans
	`, UserPlanStatusActive, UserPlanStatusExpired).Scan(&overview.Subscriptions).Error; err != nil {
		return nil, fmt.Errorf("aggregate subscriptions: %w", err)
	}

	// --- Quota 消耗（按窗口复用 SumUsedQuota） ---
	overview.Quota.Today = SumUsedQuota(LogTypeConsume, todayStart, now, "", "", "", 0)
	overview.Quota.Week = SumUsedQuota(LogTypeConsume, weekStart, now, "", "", "", 0)
	overview.Quota.Month = SumUsedQuota(LogTypeConsume, monthStart, now, "", "", "", 0)
	overview.Quota.Total = SumUsedQuota(LogTypeConsume, 0, 0, "", "", "", 0)

	// --- 收入 ---
	if err := aggregateRevenue(&overview.Revenue, r.StartTs, r.EndTs); err != nil {
		return nil, err
	}

	// --- 活跃用户（近 7 日有 consume 日志 且 request_count > 0）---
	if err := DB.Raw(`
		SELECT COUNT(DISTINCT l.user_id) AS active_7d
		FROM logs l
		INNER JOIN users u ON u.id = l.user_id
		WHERE l.type = ?
		  AND l.created_at >= ?
		  AND u.request_count > 0
	`, LogTypeConsume, weekStart).Scan(&overview.Users.Active7d).Error; err != nil {
		return nil, fmt.Errorf("aggregate active users: %w", err)
	}

	return overview, nil
}

// aggregateRevenue 按时间窗聚合订单收入与退款。
// aggregateRevenue aggregates revenue by time window.
// 版本: v0.0.16
// 日期: 2026-09-11
func aggregateRevenue(out *AdminDashboardRevenue, startTs, endTs int64) error {
	tx := DB.Table("orders").Where("status = ?", OrderStatusPaid)
	if startTs > 0 {
		tx = tx.Where("pay_time >= ?", startTs)
	}
	if endTs > 0 {
		tx = tx.Where("pay_time <= ?", endTs)
	}
	var row struct {
		Total        float64
		Topup        float64
		Subscription float64
	}
	if err := tx.Select(`
		COALESCE(SUM(amount), 0) AS total,
		COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS topup,
		COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS subscription
	`, OrderTypeTopup, OrderTypePlanSubscription).Scan(&row).Error; err != nil {
		return fmt.Errorf("aggregate paid revenue: %w", err)
	}
	out.Total = row.Total
	out.Topup = row.Topup
	out.Subscription = row.Subscription

	refundTx := DB.Table("orders").Where("status = ?", OrderStatusRefunded)
	if startTs > 0 {
		refundTx = refundTx.Where("pay_time >= ?", startTs)
	}
	if endTs > 0 {
		refundTx = refundTx.Where("pay_time <= ?", endTs)
	}
	if err := refundTx.Select("COALESCE(SUM(amount), 0) AS total").Scan(&out.Refund).Error; err != nil {
		return fmt.Errorf("aggregate refund: %w", err)
	}
	return nil
}

// GetAdminTopUsers 返回活跃用户排行榜（默认按 request_count 降序，quota 降序）。
// 口径：近 range 时间内有 consume 日志、且 request_count > 0 的用户。
// GetAdminTopUsers returns the active-user leaderboard.
// 版本: v0.0.16
// 日期: 2026-09-11
func GetAdminTopUsers(rawRange string, limit int) ([]AdminTopUserRow, string, error) {
	preset := resolveRangePreset(rawRange)
	r := ParseAdminDashboardRange(rawRange)

	if limit <= 0 || limit > 200 {
		limit = 20
	}

	type joined struct {
		Id              int
		Username        string
		Email           string
		Quota           int64
		RequestCountCol int64
		QuotaCol        int64
	}
	var rows []joined
	q := `
		SELECT u.id AS id,
		       u.username AS username,
		       u.email AS email,
		       u.quota AS quota,
		       COUNT(l.id) AS request_count_col,
		       COALESCE(SUM(l.quota), 0) AS quota_col
		FROM users u
		LEFT JOIN logs l
		  ON l.user_id = u.id
		 AND l.type = ?
		 AND l.created_at >= ?
		WHERE u.status != ?
		  AND u.request_count > 0
		GROUP BY u.id
		ORDER BY request_count_col DESC, quota_col DESC
		LIMIT ?
	`
	if err := LOG_DB.Raw(q, LogTypeConsume, r.StartTs, UserStatusDeleted, limit).Scan(&rows).Error; err != nil {
		return nil, preset, fmt.Errorf("aggregate top users: %w", err)
	}

	if len(rows) == 0 {
		return []AdminTopUserRow{}, preset, nil
	}

	// 批量取当前活跃订阅名称（避免 N+1）
	userIds := make([]int, 0, len(rows))
	for _, r := range rows {
		userIds = append(userIds, r.Id)
	}
	planMap, err := batchTopActivePlanNames(userIds)
	if err != nil {
		// 不致命，返回空套餐名
		planMap = map[int]string{}
	}

	out := make([]AdminTopUserRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, AdminTopUserRow{
			Id:              r.Id,
			Username:        r.Username,
			Email:           r.Email,
			RequestCount:    r.RequestCountCol,
			Quota:           r.QuotaCol,
			Balance:         r.Quota,
			CurrentPlanName: planMap[r.Id],
		})
	}
	return out, preset, nil
}

// batchTopActivePlanNames 批量取每个用户的「当前激活订阅」的套餐名。
// 只取最早到期的那条，避免展示成多条。
// batchTopActivePlanNames batches current plan name per user.
// 版本: v0.0.16
// 日期: 2026-09-11
func batchTopActivePlanNames(userIds []int) (map[int]string, error) {
	out := make(map[int]string, len(userIds))
	if len(userIds) == 0 {
		return out, nil
	}
	now := helper.GetTimestamp()
	type row struct {
		UserId int
		Name   string
	}
	var rows []row
	if err := DB.Raw(`
		SELECT up.user_id AS user_id, p.name AS name
		FROM user_plans up
		INNER JOIN plans p ON p.id = up.plan_id
		WHERE up.user_id IN ? AND up.status = ? AND up.end_time > ?
		ORDER BY up.user_id, up.end_time ASC
	`, userIds, UserPlanStatusActive, now).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		if _, ok := out[r.UserId]; ok {
			continue
		}
		out[r.UserId] = r.Name
	}
	return out, nil
}

// AdminChartsRange 图表窗口：按 range 解析出 day 对齐的起止时间戳与天数。
// AdminChartsRange holds the day-aligned window for chart queries.
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminChartsRange struct {
	Key   string
	Start int64
	End   int64
	Days  int
}

// ParseAdminChartsRange 把 ?range=today|7d|30d|all 解析为 day 对齐窗口。
// all 收敛为最近 30 天，避免超长历史拖垮响应。
// ParseAdminChartsRange maps ?range= to a day-aligned window.
// "all" is capped to the last 30 days to bound the response size.
// 版本: v0.0.17
// 日期: 2026-09-11
func ParseAdminChartsRange(raw string) AdminChartsRange {
	now := helper.GetTimestamp()
	daySec := int64(86400)
	todayStart := now - (now % daySec)
	key := resolveRangePreset(raw)
	days := 7
	switch key {
	case "today":
		days = 1
	case "30d", "all":
		days = 30
	}
	start := todayStart - int64(days-1)*daySec
	return AdminChartsRange{Key: key, Start: start, End: now, Days: days}
}

// SearchAdminLogsByDayAndModel 全站按 day × model 聚合（无 user_id 过滤）。
// 返回结构与 model.SearchLogsByDayAndModel 完全一致（[]*LogStatistic），
// 前端可与 api/user/dashboard 复用同一套图表构建逻辑。
// SearchAdminLogsByDayAndModel aggregates the whole site by day × model.
// The row shape mirrors model.SearchLogsByDayAndModel ([]*LogStatistic),
// so the frontend can reuse the same chart-building logic as api/user/dashboard.
// 版本: v0.0.17
// 日期: 2026-09-11
func SearchAdminLogsByDayAndModel(start, end int64) ([]*LogStatistic, error) {
	groupSelect := "DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d') as day"

	if common.UsingPostgreSQL {
		groupSelect = "TO_CHAR(date_trunc('day', to_timestamp(created_at)), 'YYYY-MM-DD') as day"
	}

	if common.UsingSQLite {
		groupSelect = "strftime('%Y-%m-%d', datetime(created_at, 'unixepoch')) as day"
	}

	var out []*LogStatistic
	err := LOG_DB.Raw(`
		SELECT `+groupSelect+`,
		model_name, count(1) as request_count,
		sum(quota) as quota,
		sum(prompt_tokens) as prompt_tokens,
		sum(completion_tokens) as completion_tokens
		FROM logs
		WHERE type=2
		AND created_at BETWEEN ? AND ?
		GROUP BY day, model_name
		ORDER BY day, model_name
	`, start, end).Scan(&out).Error

	return out, err
}
