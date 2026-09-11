// admin_dashboard_test.go 管理员仪表盘聚合查询单测
// Unit tests for admin_dashboard aggregate queries.
// 版本: v0.0.17
// 日期: 2026-09-11
// 作者: opencode
//
// 覆盖：
//   - ParseAdminDashboardRange 在四种 range 下的边界
//   - GetAdminDashboardOverview 各段聚合（users/tokens/channels/plans/redemptions/subscriptions/revenue/active_7d）
//   - GetAdminTopUsers 排行榜的 request_count>0 过滤与排序
//   - GetAdminModelDistribution 全站 Top N 模型 × 7 日聚合

package model

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
)

// setupAdminDashboardTestDB 同时初始化 DB（业务表）与 LOG_DB（日志表），
// 因为管理员仪表盘需要跨两库聚合。
// setupAdminDashboardTestDB initializes both DB and LOG_DB because the
// admin dashboard aggregates across both.
// 版本: v0.0.16
// 日期: 2026-09-11
func setupAdminDashboardTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&User{}, &Token{}, &Channel{}, &Plan{}, &UserPlan{},
		&Redemption{}, &Order{}, &Log{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	DB = db
	LOG_DB = db
	common.RedisEnabled = false
	common.UsingSQLite = true
	common.UsingPostgreSQL = false
	common.UsingMySQL = false
}

// makeUser 工厂：填齐 admin 仪表盘聚合所需的最小字段集。
// makeUser factory: minimum fields required by the aggregates.
// 版本: v0.0.16
// 日期: 2026-09-11
func makeUser(username string, status int, role int, requestCount int, quota int64) *User {
	now := helper.GetTimestamp()
	return &User{
		Username:     username,
		Password:     "x",
		DisplayName:  username,
		Role:         role,
		Status:       status,
		Group:        "default",
		AccessToken:  username + "-token",
		AffCode:      username + "-aff",
		Quota:        quota,
		UsedQuota:    0,
		RequestCount: requestCount,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// makeConsumeLog 工厂：用于构造排行榜的日志。
// makeConsumeLog factory for top-users test data.
// 版本: v0.0.16
// 日期: 2026-09-11
func makeConsumeLog(userId int, quota int, createdAt int64) *Log {
	return &Log{
		UserId:           userId,
		CreatedAt:        createdAt,
		Type:             LogTypeConsume,
		ModelName:        "gpt-4o-mini",
		Quota:            quota,
		PromptTokens:     100,
		CompletionTokens: 200,
		ChannelId:        1,
	}
}

// TestParseAdminDashboardRange 验证 4 个预设的范围计算。
// TestParseAdminDashboardRange covers all four range presets.
// 版本: v0.0.16
// 日期: 2026-09-11
func TestParseAdminDashboardRange(t *testing.T) {
	now := time.Now().Unix()
	daySec := int64(86400)
	startOfToday := now - (now % daySec)

	cases := []struct {
		in        string
		key       string
		startMin  int64
		startMax  int64
		endExpect int64
		trendDays int
	}{
		{"today", "today", startOfToday - 2, startOfToday + 2, now, 1},
		{"7d", "7d", now - 7*daySec - 2, now - 7*daySec + 2, now, 7},
		{"30d", "30d", now - 30*daySec - 2, now - 30*daySec + 2, now, 30},
		{"all", "all", 0, 0, 0, 30},
		{"unknown", "7d", now - 7*daySec - 2, now - 7*daySec + 2, now, 7},
	}
	for _, c := range cases {
		got := ParseAdminDashboardRange(c.in)
		if got.Key != c.key {
			t.Errorf("range=%q key=%q want=%q", c.in, got.Key, c.key)
		}
		if got.StartTs < c.startMin || got.StartTs > c.startMax {
			t.Errorf("range=%q StartTs=%d out of [%d,%d]", c.in, got.StartTs, c.startMin, c.startMax)
		}
		if c.endExpect != 0 {
			if got.EndTs < c.endExpect-2 || got.EndTs > c.endExpect+2 {
				t.Errorf("range=%q EndTs=%d want≈%d", c.in, got.EndTs, c.endExpect)
			}
		} else if got.EndTs != 0 {
			t.Errorf("range=%q expected EndTs=0, got %d", c.in, got.EndTs)
		}
		if got.TrendDays != c.trendDays {
			t.Errorf("range=%q TrendDays=%d want=%d", c.in, got.TrendDays, c.trendDays)
		}
	}
}

// TestGetAdminDashboardOverview_Smoke 验证空库 + 单条样本数据的聚合。
// TestGetAdminDashboardOverview_Smoke exercises the aggregator on a
// representative dataset.
// 版本: v0.0.16
// 日期: 2026-09-11
func TestGetAdminDashboardOverview_Smoke(t *testing.T) {
	setupAdminDashboardTestDB(t)
	now := helper.GetTimestamp()
	daySec := int64(86400)

	// Users: 1 enabled (new today, request_count=10), 1 disabled, 1 deleted
	if err := DB.Create(makeUser("alice", UserStatusEnabled, RoleCommonUser, 10, 100000)).Error; err != nil {
		t.Fatalf("create alice: %v", err)
	}
	if err := DB.Create(makeUser("bob", UserStatusDisabled, RoleCommonUser, 0, 0)).Error; err != nil {
		t.Fatalf("create bob: %v", err)
	}
	if err := DB.Create(makeUser("ghost", UserStatusDeleted, RoleCommonUser, 0, 0)).Error; err != nil {
		t.Fatalf("create ghost: %v", err)
	}

	// Tokens: 1 enabled, 1 disabled
	if err := DB.Create(&Token{
		UserId:      1,
		Key:         "k1",
		Status:      TokenStatusEnabled,
		Name:        "t1",
		CreatedTime: now,
	}).Error; err != nil {
		t.Fatalf("create token1: %v", err)
	}
	if err := DB.Create(&Token{
		UserId:      1,
		Key:         "k2",
		Status:      TokenStatusDisabled,
		Name:        "t2",
		CreatedTime: now,
	}).Error; err != nil {
		t.Fatalf("create token2: %v", err)
	}

	// Channels
	if err := DB.Create(&Channel{
		Type:   1,
		Key:    "ck",
		Name:   "c1",
		Status: ChannelStatusEnabled,
		Models: "gpt-4o-mini",
		Group:  "default",
	}).Error; err != nil {
		t.Fatalf("create channel: %v", err)
	}

	// Plans
	if err := DB.Create(&Plan{
		Name:         "Pro",
		Price:        99.0,
		DurationDays: 30,
		Status:       PlanStatusEnabled,
	}).Error; err != nil {
		t.Fatalf("create plan: %v", err)
	}

	// Redemptions: 2 used + 1 enabled
	for i := 0; i < 2; i++ {
		if err := DB.Create(&Redemption{
			Key:         "r-used-" + string(rune('a'+i)),
			Status:      RedemptionCodeStatusUsed,
			Name:        "r",
			Quota:       1000,
			CreatedTime: now,
		}).Error; err != nil {
			t.Fatalf("create redemption used: %v", err)
		}
	}
	if err := DB.Create(&Redemption{
		Key:         "r-enabled",
		Status:      RedemptionCodeStatusEnabled,
		Name:        "r",
		Quota:       1000,
		CreatedTime: now,
	}).Error; err != nil {
		t.Fatalf("create redemption enabled: %v", err)
	}

	// UserPlans: 1 active, 1 expired (status=0 is overridden by GORM default:1 tag,
	// so we explicitly UPDATE the column afterwards).
	if err := DB.Create(&UserPlan{
		UserId:      1,
		PlanId:      1,
		StartTime:   now - daySec,
		EndTime:     now + daySec,
		Status:      UserPlanStatusActive,
		BillingType: BillingTypeToken,
	}).Error; err != nil {
		t.Fatalf("create user_plan active: %v", err)
	}
	expired := &UserPlan{
		UserId:      1,
		PlanId:      1,
		StartTime:   now - 2*daySec,
		EndTime:     now - daySec,
		Status:      UserPlanStatusActive,
		BillingType: BillingTypeToken,
	}
	if err := DB.Create(expired).Error; err != nil {
		t.Fatalf("create user_plan expired: %v", err)
	}
	if err := DB.Model(expired).Update("status", UserPlanStatusExpired).Error; err != nil {
		t.Fatalf("flip expired status: %v", err)
	}

	// Orders: 1 paid topup + 1 paid subscription + 1 refunded
	if err := DB.Create(&Order{
		Type:       OrderTypeTopup,
		Source:     OrderSourceUserSelf,
		OrderNo:    "TP-test-1",
		UserId:     1,
		Amount:     50.0,
		Status:     OrderStatusPaid,
		PayStatus:  OrderPayStatusPaid,
		CreateTime: now,
		PayTime:    now - 3600,
	}).Error; err != nil {
		t.Fatalf("create order topup: %v", err)
	}
	if err := DB.Create(&Order{
		Type:       OrderTypePlanSubscription,
		Source:     OrderSourceUserSelf,
		OrderNo:    "TB-test-1",
		UserId:     1,
		PlanId:     1,
		Amount:     99.0,
		Status:     OrderStatusPaid,
		PayStatus:  OrderPayStatusPaid,
		CreateTime: now,
		PayTime:    now - 1800,
	}).Error; err != nil {
		t.Fatalf("create order subscription: %v", err)
	}
	if err := DB.Create(&Order{
		Type:       OrderTypeTopup,
		Source:     OrderSourceUserSelf,
		OrderNo:    "TP-test-refund",
		UserId:     1,
		Amount:     10.0,
		Status:     OrderStatusRefunded,
		PayStatus:  OrderPayStatusRefunded,
		CreateTime: now,
		PayTime:    now - 7200,
	}).Error; err != nil {
		t.Fatalf("create order refund: %v", err)
	}

	// Logs: 2 consume rows for alice today
	for i := 0; i < 2; i++ {
		if err := LOG_DB.Create(makeConsumeLog(1, 500, now-100-int64(i))).Error; err != nil {
			t.Fatalf("create log: %v", err)
		}
	}

	ov, err := GetAdminDashboardOverview("7d")
	if err != nil {
		t.Fatalf("overview: %v", err)
	}

	if ov.Users.Total != 3 {
		t.Errorf("users.total=%d want=3", ov.Users.Total)
	}
	if ov.Users.Enabled != 1 || ov.Users.Disabled != 1 || ov.Users.Deleted != 1 {
		t.Errorf("users.status=%+v want enabled=1 disabled=1 deleted=1", ov.Users)
	}
	if ov.Users.NewToday != 3 {
		t.Errorf("users.new_today=%d want=3", ov.Users.NewToday)
	}
	if ov.Users.New7d != 3 {
		t.Errorf("users.new_7d=%d want=3", ov.Users.New7d)
	}
	if ov.Users.New30d != 3 {
		t.Errorf("users.new_30d=%d want=3", ov.Users.New30d)
	}
	if ov.Users.Active7d != 1 {
		t.Errorf("users.active_7d=%d want=1", ov.Users.Active7d)
	}
	if ov.Tokens.Total != 2 || ov.Tokens.Enabled != 1 {
		t.Errorf("tokens=%+v want total=2 enabled=1", ov.Tokens)
	}
	if ov.Channels.Total != 1 || ov.Channels.Enabled != 1 {
		t.Errorf("channels=%+v", ov.Channels)
	}
	if ov.Plans.Total != 1 || ov.Plans.Enabled != 1 {
		t.Errorf("plans=%+v", ov.Plans)
	}
	if ov.Redemptions.Total != 3 || ov.Redemptions.Used != 2 || ov.Redemptions.Unused != 1 {
		t.Errorf("redemptions=%+v want total=3 used=2 unused=1", ov.Redemptions)
	}
	if ov.Subscriptions.Total != 2 || ov.Subscriptions.Active != 1 || ov.Subscriptions.Expired != 1 {
		t.Errorf("subscriptions=%+v", ov.Subscriptions)
	}
	if ov.Quota.Today != 1000 {
		t.Errorf("quota.today=%d want=1000", ov.Quota.Today)
	}
	if ov.Quota.Total != 1000 {
		t.Errorf("quota.total=%d want=1000", ov.Quota.Total)
	}
	if ov.Revenue.Total != 149.0 || ov.Revenue.Topup != 50.0 || ov.Revenue.Subscription != 99.0 {
		t.Errorf("revenue=%+v", ov.Revenue)
	}
	if ov.Revenue.Refund != 10.0 {
		t.Errorf("refund=%v want=10", ov.Revenue.Refund)
	}
	if ov.Range != "7d" {
		t.Errorf("range=%q want=7d", ov.Range)
	}
}

// TestGetAdminTopUsers 验证排行榜：request_count>0 过滤 + 排序 + 当前套餐名。
// TestGetAdminTopUsers verifies the leaderboard filter, ordering, and plan name.
// 版本: v0.0.16
// 日期: 2026-09-11
func TestGetAdminTopUsers(t *testing.T) {
	setupAdminDashboardTestDB(t)
	now := helper.GetTimestamp()
	daySec := int64(86400)

	// alice: 3 logs, request_count=10, active plan
	if err := DB.Create(makeUser("alice", UserStatusEnabled, RoleCommonUser, 10, 5000)).Error; err != nil {
		t.Fatal(err)
	}
	// bob: 1 log, request_count=2, no plan
	if err := DB.Create(makeUser("bob", UserStatusEnabled, RoleCommonUser, 2, 1000)).Error; err != nil {
		t.Fatal(err)
	}
	// carol: request_count=0 (应被过滤)
	if err := DB.Create(makeUser("carol", UserStatusEnabled, RoleCommonUser, 0, 1000)).Error; err != nil {
		t.Fatal(err)
	}
	// ghost: deleted
	if err := DB.Create(makeUser("ghost", UserStatusDeleted, RoleCommonUser, 99, 0)).Error; err != nil {
		t.Fatal(err)
	}

	// Plan
	if err := DB.Create(&Plan{
		Name:         "Pro",
		Price:        99.0,
		DurationDays: 30,
		Status:       PlanStatusEnabled,
	}).Error; err != nil {
		t.Fatal(err)
	}
	// Alice 持有激活套餐
	if err := DB.Create(&UserPlan{
		UserId:      1,
		PlanId:      1,
		StartTime:   now - daySec,
		EndTime:     now + 30*daySec,
		Status:      UserPlanStatusActive,
		BillingType: BillingTypeToken,
	}).Error; err != nil {
		t.Fatal(err)
	}

	// Logs: 3 for alice, 1 for bob, 2 for ghost (该被 status 过滤)
	for i := 0; i < 3; i++ {
		if err := LOG_DB.Create(makeConsumeLog(1, 100, now-int64(i+1)*10)).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := LOG_DB.Create(makeConsumeLog(2, 200, now-5)).Error; err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := LOG_DB.Create(makeConsumeLog(4, 50, now-int64(i+1)*5)).Error; err != nil {
			t.Fatal(err)
		}
	}

	rows, preset, err := GetAdminTopUsers("7d", 20)
	if err != nil {
		t.Fatalf("top users: %v", err)
	}
	if preset != "7d" {
		t.Errorf("preset=%q want=7d", preset)
	}
	if len(rows) != 2 {
		t.Fatalf("len=%d want=2 (carol filtered out, ghost deleted)", len(rows))
	}
	if rows[0].Username != "alice" || rows[0].RequestCount != 3 || rows[0].CurrentPlanName != "Pro" {
		t.Errorf("row0=%+v", rows[0])
	}
	if rows[1].Username != "bob" || rows[1].RequestCount != 1 {
		t.Errorf("row1=%+v", rows[1])
	}
	if rows[1].CurrentPlanName != "" {
		t.Errorf("row1.plan=%q want empty", rows[1].CurrentPlanName)
	}
}

// makeConsumeLogForModel 工厂：可指定 model_name 与 created_at。
// makeConsumeLogForModel factory for model-distribution test data.
// 版本: v0.0.17
// 日期: 2026-09-11
func makeConsumeLogForModel(userId int, modelName string, quota int, createdAt int64) *Log {
	return &Log{
		UserId:           userId,
		CreatedAt:        createdAt,
		Type:             LogTypeConsume,
		ModelName:        modelName,
		Quota:            quota,
		PromptTokens:     100,
		CompletionTokens: 200,
		ChannelId:        1,
	}
}

// TestParseAdminChartsRange 验证图表窗口的 day 对齐与天数。
// TestParseAdminChartsRange verifies day-aligned windows and day counts.
// 版本: v0.0.17
// 日期: 2026-09-11
func TestParseAdminChartsRange(t *testing.T) {
	cases := []struct {
		in   string
		key  string
		days int
	}{
		{"today", "today", 1},
		{"7d", "7d", 7},
		{"30d", "30d", 30},
		{"all", "all", 30},
		{"unknown", "7d", 7},
	}
	daySec := int64(86400)
	for _, c := range cases {
		got := ParseAdminChartsRange(c.in)
		if got.Key != c.key || got.Days != c.days {
			t.Errorf("in=%q key=%q days=%d, want key=%q days=%d", c.in, got.Key, got.Days, c.key, c.days)
		}
		// End - Start 应覆盖 days 天（end 为 now，start 为 todayStart-(days-1)）
		span := got.End - got.Start
		wantMin := int64(c.days-1) * daySec
		if span < wantMin || span > wantMin+daySec {
			t.Errorf("in=%q span=%d want ~[%d,%d]", c.in, span, wantMin, wantMin+daySec)
		}
	}
}

// TestSearchAdminLogsByDayAndModel 验证全站 day×model 聚合结构与字段。
// TestSearchAdminLogsByDayAndModel verifies the day×model aggregate shape.
// 版本: v0.0.17
// 日期: 2026-09-11
func TestSearchAdminLogsByDayAndModel(t *testing.T) {
	setupAdminDashboardTestDB(t)
	now := helper.GetTimestamp()
	daySec := int64(86400)
	todayStart := now - (now % daySec)

	// 两个模型 × 两天
	type rec struct {
		user  int
		model string
		quota int
		day   int64
	}
	recs := []rec{
		{1, "gpt-4o", 1000, 0},
		{2, "gpt-4o", 2000, 0}, // 不同用户，聚合时应合并
		{1, "gpt-4o", 500, 1},
		{1, "claude-3.5", 800, 0},
		{2, "claude-3.5", 600, 1},
	}
	for _, r := range recs {
		ts := todayStart + 100 - r.day*daySec
		if err := LOG_DB.Create(makeConsumeLogForModel(r.user, r.model, r.quota, ts)).Error; err != nil {
			t.Fatalf("create log: %v", err)
		}
	}

	rows, err := SearchAdminLogsByDayAndModel(todayStart-6*daySec, now)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	// 预期 4 行：2 模型 × 2 天
	if len(rows) != 4 {
		t.Fatalf("rows=%d want=4: %+v", len(rows), rows)
	}
	// 字段齐全（无 omitempty 导致的缺字段）
	// 每行 prompt=100×N，completion=200×N（N=合并的日志条数）
	for _, r := range rows {
		if r.Day == "" || r.ModelName == "" {
			t.Errorf("empty day/model: %+v", r)
		}
		if r.PromptTokens == 0 || r.PromptTokens%100 != 0 {
			t.Errorf("prompt_tokens mismatch: %+v", r)
		}
		if r.CompletionTokens != r.PromptTokens*2 {
			t.Errorf("completion/prompt ratio mismatch: %+v", r)
		}
	}
	// 今天 gpt-4o 应合并两个用户：1000+2000=3000, request_count=2
	found := false
	for _, r := range rows {
		if r.Day == rows[0].Day && r.ModelName == "gpt-4o" {
			// rows[0].Day 是今天还是昨天取决于排序（day ASC -> 昨天在前）
			// 直接找 quota=3000 的那行
		}
		if r.ModelName == "gpt-4o" && r.Quota == 3000 {
			found = true
			if r.RequestCount != 2 {
				t.Errorf("gpt-4o request_count=%d want=2", r.RequestCount)
			}
		}
	}
	if !found {
		t.Errorf("gpt-4o merged row (quota=3000) not found: %+v", rows)
	}
}
