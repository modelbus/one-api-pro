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
	"sort"
	"time"

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

// AdminTrendPoint 单日趋势点
// AdminTrendPoint is one point in a daily trend series.
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminTrendPoint struct {
	Day    string `json:"day"`
	Count  int64  `json:"count,omitempty"`
	Quota  int64  `json:"quota,omitempty"`
	Tokens int64  `json:"tokens,omitempty"`
}

// AdminDashboardTrends 趋势数据
// AdminDashboardTrends holds trend series.
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminDashboardTrends struct {
	Requests []AdminTrendPoint `json:"requests"`
	Quota    []AdminTrendPoint `json:"quota"`
	Tokens   []AdminTrendPoint `json:"tokens"`
}

// AdminDashboardOverview 单次响应聚合：KPI + 7 日趋势
// AdminDashboardOverview is the single-shot aggregate response.
// 版本: v0.0.16
// 日期: 2026-09-11
type AdminDashboardOverview struct {
	Users         AdminDashboardUsers      `json:"users"`
	Tokens        AdminDashboardResources  `json:"tokens"`
	Channels      AdminDashboardResources  `json:"channels"`
	Plans         AdminDashboardResources  `json:"plans"`
	Redemptions   AdminDashboardResources  `json:"redemptions"`
	Subscriptions AdminDashboardResources  `json:"subscriptions"`
	Quota         AdminDashboardQuota      `json:"quota"`
	Revenue       AdminDashboardRevenue    `json:"revenue"`
	Trends        AdminDashboardTrends     `json:"trends"`
	Range         string                   `json:"range"`
	GeneratedAt   int64                    `json:"generated_at"`
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
		startOfDay := now - (now%daySec)
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

// GetAdminDashboardOverview 一次性聚合所有 KPI + 7 日趋势。
// 各子查询均为单条 SQL，避免 N+1。
// GetAdminDashboardOverview aggregates all KPIs in a single call.
// Each subsection is one SQL to avoid N+1.
// 版本: v0.0.16
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

	// --- 趋势（始终 7 日，与 range 解耦；后期按需扩展）---
	trendStart := now - 7*daySec
	requests, err := aggregateTrends(trendStart, now, "count")
	if err != nil {
		return nil, err
	}
	overview.Trends.Requests = requests
	quotaSeries, err := aggregateTrends(trendStart, now, "quota")
	if err != nil {
		return nil, err
	}
	overview.Trends.Quota = quotaSeries
	tokensSeries, err := aggregateTrends(trendStart, now, "tokens")
	if err != nil {
		return nil, err
	}
	overview.Trends.Tokens = tokensSeries

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

// aggregateTrends 按 day 聚合请求量/quota/tokens。
// metric: "count"（request_count）/ "quota" / "tokens"。
// 三种 DB 各自一份 SQL，参考 model/log.go:245 的跨库兼容写法。
// aggregateTrends aggregates daily requests / quota / tokens, with cross-DB support.
// 版本: v0.0.17
// 日期: 2026-09-11
func aggregateTrends(startTs, endTs int64, metric string) ([]AdminTrendPoint, error) {
	var dayExpr string
	switch {
	case common.UsingPostgreSQL:
		dayExpr = "TO_CHAR(date_trunc('day', to_timestamp(created_at)), 'YYYY-MM-DD')"
	case common.UsingSQLite:
		dayExpr = "strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
	default:
		dayExpr = "DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
	}
	col := "count(1)"
	switch metric {
	case "quota":
		col = "sum(quota)"
	case "tokens":
		col = "sum(prompt_tokens + completion_tokens)"
	}
	q := fmt.Sprintf(`
		SELECT %s AS day, %s AS value
		FROM logs
		WHERE type = ? AND created_at BETWEEN ? AND ?
		GROUP BY day
		ORDER BY day
	`, dayExpr, col)

	type row struct {
		Day   string
		Value int64
	}
	var rows []row
	if err := LOG_DB.Raw(q, LogTypeConsume, startTs, endTs).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("aggregate trends(%s): %w", metric, err)
	}

	// 补齐缺失日期为 0，让前端图表 X 轴连续。
	bucket := make(map[string]int64, len(rows))
	for _, r := range rows {
		bucket[r.Day] = r.Value
	}

	daySec := int64(86400)
	startDay := startTs - (startTs % daySec)
	endDay := endTs - (endTs % daySec)
	days := int((endDay-startDay)/daySec) + 1
	if days < 1 {
		days = 1
	}
	out := make([]AdminTrendPoint, 0, days)
	for i := 0; i < days; i++ {
		t := time.Unix(startDay+int64(i)*daySec, 0).UTC()
		key := t.Format("2006-01-02")
		var v int64
		if got, ok := bucket[key]; ok {
			v = got
		}
		switch metric {
		case "quota":
			out = append(out, AdminTrendPoint{Day: key, Quota: v})
		case "tokens":
			out = append(out, AdminTrendPoint{Day: key, Tokens: v})
		default:
			out = append(out, AdminTrendPoint{Day: key, Count: v})
		}
	}
	return out, nil
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

// AdminModelDistributionRow 全站模型分布的单行：按模型聚合 quota/requests，附带天数窗口。
// AdminModelDistributionRow is one model row in the admin distribution chart.
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminModelDistributionRow struct {
	ModelName        string `json:"model_name" gorm:"column:model_name"`
	RequestCount     int64  `json:"request_count" gorm:"column:request_count"`
	Quota            int64  `json:"quota" gorm:"column:quota"`
	PromptTokens     int64  `json:"prompt_tokens" gorm:"column:prompt_tokens"`
	CompletionTokens int64  `json:"completion_tokens" gorm:"column:completion_tokens"`
	// DayQuota: 该模型每天的 quota 序列（按 Days 顺序对齐），用于堆叠柱图。
	// DayQuota: per-day quota series aligned with Days, for stacked bar chart.
	DayQuota []int64 `json:"day_quota,omitempty"`
}

// AdminModelDistribution 全站模型分布聚合：按模型 × 日 双维度，限定窗口内 Top N 模型。
// AdminModelDistribution aggregates model distribution across all users.
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminModelDistribution struct {
	Days  []string                    `json:"days"`
	Items []AdminModelDistributionRow `json:"items"`
	Range string                      `json:"range"`
}

// GetAdminModelDistribution 全站模型用量分布：返回按 quota 降序的 Top N 模型 + 该窗口的 7/30 天日期序列。
// 复用 model/log.go 的三库 day 表达式，但去掉 user_id 过滤。
// GetAdminModelDistribution returns top-N models by quota + day series in the window.
// 版本: v0.0.17
// 日期: 2026-09-11
func GetAdminModelDistribution(rawRange string, topN int) (*AdminModelDistribution, error) {
	if topN <= 0 || topN > 50 {
		topN = 8
	}
	r := ParseAdminDashboardRange(rawRange)

	var dayExpr string
	switch {
	case common.UsingPostgreSQL:
		dayExpr = "TO_CHAR(date_trunc('day', to_timestamp(created_at)), 'YYYY-MM-DD')"
	case common.UsingSQLite:
		dayExpr = "strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
	default:
		dayExpr = "DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
	}

	// 取窗口内 quota 总额 Top N 模型名。
	// First pick top-N models by total quota within the window.
	type modelAgg struct {
		ModelName string
		Total     int64
	}
	var top []modelAgg
	if err := LOG_DB.Raw(`
		SELECT model_name, COALESCE(SUM(quota), 0) AS total
		FROM logs
		WHERE type = ?
		  AND model_name != ''
		  AND (? = 0 OR created_at >= ?)
		  AND (? = 0 OR created_at <= ?)
		GROUP BY model_name
		ORDER BY total DESC
		LIMIT ?
	`, LogTypeConsume, r.StartTs, r.StartTs, r.EndTs, r.EndTs, topN).Scan(&top).Error; err != nil {
		return nil, fmt.Errorf("aggregate top models: %w", err)
	}

	out := &AdminModelDistribution{
		Days:  []string{},
		Items: []AdminModelDistributionRow{},
		Range: r.Key,
	}

	// 构造日期序列：默认固定 7 天，便于前端堆叠柱图横轴对齐。
	// 即使没有数据也要返回 7 天序列，让前端 X 轴稳定。
	// 必须用 UTC，与 SQLite/PostgreSQL 的 strftime/to_timestamp 默认时区一致，
	// 否则非 UTC 时区下 day 对齐会偏差一天。
	// Build day series (fixed 7 days by default, aligned with admin trends).
	// Always return 7 days even when empty so the chart X axis stays stable.
	// Must use UTC to match SQLite strftime / PG to_timestamp default timezone,
	// otherwise day alignment drifts by 1 in non-UTC zones.
	now := helper.GetTimestamp()
	daySec := int64(86400)
	todayStart := now - (now % daySec)
	days := 7
	for i := days - 1; i >= 0; i-- {
		t := time.Unix(todayStart-int64(i)*daySec, 0).UTC()
		out.Days = append(out.Days, t.Format("2006-01-02"))
	}
	dayStart := todayStart - int64(days-1)*daySec

	if len(top) == 0 {
		return out, nil
	}

	// 拉 Top N 模型 × 7 天的 quota + requests。
	// Pull per-model per-day aggregates for the top-N models over the last `days`.
	names := make([]string, 0, len(top))
	for _, m := range top {
		names = append(names, m.ModelName)
	}
	type row struct {
		ModelName        string
		Day              string
		RequestCount     int64
		Quota            int64
		PromptTokens     int64
		CompletionTokens int64
	}
	var rows []row
	q := fmt.Sprintf(`
		SELECT model_name, %s AS day,
		       COUNT(1) AS request_count,
		       COALESCE(SUM(quota), 0) AS quota,
		       COALESCE(SUM(prompt_tokens), 0) AS prompt_tokens,
		       COALESCE(SUM(completion_tokens), 0) AS completion_tokens
		FROM logs
		WHERE type = ?
		  AND model_name IN ?
		  AND created_at BETWEEN ? AND ?
		GROUP BY model_name, day
	`, dayExpr)
	if err := LOG_DB.Raw(q, LogTypeConsume, names, dayStart, now).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("aggregate model-day: %w", err)
	}

	// 汇总到 Items 数组（按 Top N 顺序）。
	bucket := make(map[string]*AdminModelDistributionRow, len(top))
	for _, r := range rows {
		key := r.ModelName
		if _, ok := bucket[key]; !ok {
			bucket[key] = &AdminModelDistributionRow{
				ModelName: key,
				DayQuota:  make([]int64, len(out.Days)),
			}
		}
		bucket[key].RequestCount += r.RequestCount
		bucket[key].Quota += r.Quota
		bucket[key].PromptTokens += r.PromptTokens
		bucket[key].CompletionTokens += r.CompletionTokens
		// 把当日 quota 写入对应索引，保持与 out.Days 顺序一致。
		// Fill day-aligned quota at the matching index to keep Days order.
		for di, day := range out.Days {
			if day == r.Day {
				bucket[key].DayQuota[di] += r.Quota
				break
			}
		}
	}
	for _, m := range top {
		if row, ok := bucket[m.ModelName]; ok {
			out.Items = append(out.Items, *row)
		} else {
			out.Items = append(out.Items, AdminModelDistributionRow{
				ModelName: m.ModelName,
				DayQuota:  make([]int64, len(out.Days)),
			})
		}
	}
	return out, nil
}

// AdminUsageDetailRow 使用明细的单行：每个模型 × 每天的请求/消耗明细。
// AdminUsageDetailRow is one row in the admin usage-details table.
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminUsageDetailRow struct {
	Day              string `json:"day" gorm:"column:day"`
	ModelName        string `json:"model_name" gorm:"column:model_name"`
	RequestCount     int64  `json:"request_count" gorm:"column:request_count"`
	Quota            int64  `json:"quota" gorm:"column:quota"`
	PromptTokens     int64  `json:"prompt_tokens" gorm:"column:prompt_tokens"`
	CompletionTokens int64  `json:"completion_tokens" gorm:"column:completion_tokens"`
}

// AdminUsageDetails 全站使用明细：每个 Top N 模型 × 每个 day 的明细行。
// 与 model-distribution 的差别：这里返回明细行（透视前），方便前端做表格。
// AdminUsageDetails is the full per-model per-day breakdown.
// 版本: v0.0.17
// 日期: 2026-09-11
type AdminUsageDetails struct {
	Days  []string              `json:"days"`
	Items []AdminUsageDetailRow `json:"items"`
	Range string                `json:"range"`
}

// GetAdminUsageDetails 返回 Top N 模型 × 该窗口内每日明细（透视前）。
// 与 GetAdminModelDistribution 共用同样的窗口/Top-N 逻辑，但返回明细分行而非聚合。
// GetAdminUsageDetails returns the breakdown rows (unpivoted) for the same window.
// 版本: v0.0.17
// 日期: 2026-09-11
func GetAdminUsageDetails(rawRange string, topN int) (*AdminUsageDetails, error) {
	if topN <= 0 || topN > 50 {
		topN = 8
	}
	r := ParseAdminDashboardRange(rawRange)

	var dayExpr string
	switch {
	case common.UsingPostgreSQL:
		dayExpr = "TO_CHAR(date_trunc('day', to_timestamp(created_at)), 'YYYY-MM-DD')"
	case common.UsingSQLite:
		dayExpr = "strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
	default:
		dayExpr = "DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
	}

	// 窗口内 quota 总额 Top N 模型名（与 distribution 共用逻辑，但窗口独立）
	// Top-N models by total quota within the window.
	type modelAgg struct {
		ModelName string
		Total     int64
	}
	var top []modelAgg
	if err := LOG_DB.Raw(`
		SELECT model_name, COALESCE(SUM(quota), 0) AS total
		FROM logs
		WHERE type = ?
		  AND model_name != ''
		  AND (? = 0 OR created_at >= ?)
		  AND (? = 0 OR created_at <= ?)
		GROUP BY model_name
		ORDER BY total DESC
		LIMIT ?
	`, LogTypeConsume, r.StartTs, r.StartTs, r.EndTs, r.EndTs, topN).Scan(&top).Error; err != nil {
		return nil, fmt.Errorf("aggregate top models: %w", err)
	}

	out := &AdminUsageDetails{
		Days:  []string{},
		Items: []AdminUsageDetailRow{},
		Range: r.Key,
	}
	if len(top) == 0 {
		return out, nil
	}

	names := make([]string, 0, len(top))
	for _, m := range top {
		names = append(names, m.ModelName)
	}

	q := fmt.Sprintf(`
		SELECT %s AS day,
		       model_name,
		       COUNT(1) AS request_count,
		       COALESCE(SUM(quota), 0) AS quota,
		       COALESCE(SUM(prompt_tokens), 0) AS prompt_tokens,
		       COALESCE(SUM(completion_tokens), 0) AS completion_tokens
		FROM logs
		WHERE type = ?
		  AND model_name IN ?
		  AND (? = 0 OR created_at >= ?)
		  AND (? = 0 OR created_at <= ?)
		GROUP BY day, model_name
		ORDER BY day ASC, quota DESC
	`, dayExpr)
	var rows []AdminUsageDetailRow
	if err := LOG_DB.Raw(q, LogTypeConsume, names, r.StartTs, r.StartTs, r.EndTs, r.EndTs).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("aggregate usage details: %w", err)
	}

	// 构造 day 序列（按 day ASC 去重），便于前端画表格时间轴。
	daySet := make(map[string]bool, len(rows))
	for _, r := range rows {
		daySet[r.Day] = true
	}
	days := make([]string, 0, len(daySet))
	for d := range daySet {
		days = append(days, d)
	}
	sort.Strings(days)
	out.Days = days
	out.Items = rows
	return out, nil
}
