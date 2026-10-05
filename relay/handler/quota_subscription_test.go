// quota_subscription_test.go 订阅计费的落库行为测试
// Subscription billing persistence tests
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
//
// 覆盖：
//   - 订阅消费不扣 users.quota（只消耗 user_plans.used_amount）
//   - 订阅消费计入 users.used_quota / request_count
//   - plan_usages 窗口用量被累加
package controller

import (
	"context"
	"testing"

	dbmodel "github.com/modelbus/one-api-pro/model"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
	billingratio "github.com/modelbus/one-api-pro/relay/billing/ratio"
	"github.com/modelbus/one-api-pro/relay/meta"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

// extendSubscriptionTestDB 在预扣测试库基础上补齐订阅计费所需的表，
// 并把 DB 方言切到 SQLite（IncrementPlanUsage 会按方言选择 upsert 语法）。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func extendSubscriptionTestDB(t *testing.T) {
	t.Helper()
	common.UsingSQLite = true
	if err := dbmodel.DB.AutoMigrate(&dbmodel.Plan{}, &dbmodel.UserPlan{}, &dbmodel.PlanUsage{}); err != nil {
		t.Fatalf("migrate subscription tables: %v", err)
	}
}

// seedSubscription 创建套餐 + 订阅（快照字段取自套餐），返回 user_plan id。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func seedSubscription(t *testing.T, userId int, modelName string, virtualAmount int64) int {
	t.Helper()
	limits := `{"` + modelName + `": {"period_h": 5, "token_period": 1000000}}`
	plan := &dbmodel.Plan{
		Name:          "test-plan",
		BillingType:   dbmodel.BillingTypeToken,
		VirtualAmount: virtualAmount,
		ModelLimits:   limits,
	}
	if err := plan.Insert(); err != nil {
		t.Fatalf("insert plan: %v", err)
	}
	now := helper.GetTimestamp()
	up := &dbmodel.UserPlan{
		UserId:        userId,
		PlanId:        int(plan.Id),
		StartTime:     now,
		EndTime:       now + 86400,
		Status:        dbmodel.UserPlanStatusActive,
		BillingType:   dbmodel.BillingTypeToken,
		ModelLimits:   limits,
		VirtualAmount: virtualAmount,
	}
	if err := up.Insert(); err != nil {
		t.Fatalf("insert user plan: %v", err)
	}
	return int(up.Id)
}

// TestPostConsumeQuota_SubscriptionCountsStatsNotBalance 订阅消费的账目断言：
// 不扣余额、统计照算、套餐额度累加。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestPostConsumeQuota_SubscriptionCountsStatsNotBalance(t *testing.T) {
	setupPreConsumeTestDB(t)
	extendPreConsumeTestDB(t)
	extendSubscriptionTestDB(t)

	const (
		initialUserQuota  = int64(500_000_000) // ¥500
		initialTokenQuota = int64(500_000_000)
		virtualAmount     = int64(100_000_000) // ¥100
	)
	uid := seedUser(t, initialUserQuota)
	tid := seedToken(t, uid, initialTokenQuota)
	upId := seedSubscription(t, uid, "gpt-4o", virtualAmount)

	m := &meta.Meta{UserId: uid, TokenId: tid, ChannelId: 0, PlanId: upId, OriginModelName: "gpt-4o"}
	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 0}
	usage := &relaymodel.Usage{PromptTokens: 1000, CompletionTokens: 500}
	priceResult := &billingratio.PriceResult{
		InputPrice:  2.5,
		OutputPrice: 10.0,
		CachedPrice: 1.25,
		BillingType: dbmodel.BillingTypeToken,
		Found:       true,
	}

	postConsumeQuota(context.Background(), usage, m, req, 0, priceResult, 1.0, false)

	expectedQuota := int64(2.5*1000 + 10.0*500)

	// 1) 余额不变（订阅消费不扣 users.quota / token.remain_quota）。
	if got := readUserQuota(t, uid); got != initialUserQuota {
		t.Fatalf("user.quota = %d，期望不变 %d", got, initialUserQuota)
	}
	tk := readToken(t, tid)
	if tk.RemainQuota != initialTokenQuota {
		t.Fatalf("token.remain_quota = %d，期望不变 %d", tk.RemainQuota, initialTokenQuota)
	}

	// 2) 统计照常累加：used_quota / request_count。
	var u dbmodel.User
	if err := dbmodel.DB.First(&u, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if u.UsedQuota != expectedQuota {
		t.Fatalf("user.used_quota = %d，期望 %d", u.UsedQuota, expectedQuota)
	}
	if u.RequestCount != 1 {
		t.Fatalf("user.request_count = %d，期望 1", u.RequestCount)
	}

	// 3) 套餐虚拟余额被累加。
	used, err := dbmodel.GetUserPlanUsedAmount(upId)
	if err != nil {
		t.Fatalf("get user plan used amount: %v", err)
	}
	if used != expectedQuota {
		t.Fatalf("user_plan.used_amount = %d，期望 %d", used, expectedQuota)
	}

	// 4) 窗口用量被累加（三个窗口类型各一条）。
	pus, err := dbmodel.GetPlanUsageByUserPlanId(upId)
	if err != nil {
		t.Fatalf("get plan usage: %v", err)
	}
	if len(pus) != 3 {
		t.Fatalf("plan_usages rows = %d，期望 3", len(pus))
	}
	for _, pu := range pus {
		if pu.Requests != 1 {
			t.Fatalf("plan_usage(%s).requests = %d，期望 1", pu.WindowType, pu.Requests)
		}
		if pu.PromptTokens != 1000 || pu.CompletionTokens != 500 {
			t.Fatalf("plan_usage(%s) tokens = %d/%d，期望 1000/500", pu.WindowType, pu.PromptTokens, pu.CompletionTokens)
		}
	}
}
