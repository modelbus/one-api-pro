// plan_virtual_amount_test.go 套餐配置校验与虚拟余额判定的单元测试
// Unit tests for plan config validation and virtual-amount quota checks
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
//
// 覆盖：
//   - Plan.ValidateConfig 的 billing_type / model_limits / 维度匹配校验
//   - CheckPlanQuota 的虚拟余额耗尽、快照兜底、非法 JSON、空 limits
//   - Increment/GetUserPlanUsedAmount 与 RemainingVirtualAmount
package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
)

// setupPlanQuotaTestDB 初始化 in-memory sqlite + plan/user_plan/plan_usage 表。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func setupPlanQuotaTestDB(t *testing.T) {
	t.Helper()
	common.RedisEnabled = false
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&Plan{}, &UserPlan{}, &PlanUsage{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	DB = db
}

// createTestPlan 插入一个套餐并返回。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func createTestPlan(t *testing.T, plan *Plan) *Plan {
	t.Helper()
	if err := plan.Insert(); err != nil {
		t.Fatalf("insert plan: %v", err)
	}
	return plan
}

// createTestUserPlan 为 plan 插入一条活跃订阅，字段取自套餐（模拟激活时的快照）。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func createTestUserPlan(t *testing.T, plan *Plan, userId int) *UserPlan {
	t.Helper()
	now := helper.GetTimestamp()
	up := &UserPlan{
		UserId:        userId,
		PlanId:        int(plan.Id),
		StartTime:     now,
		EndTime:       now + 86400,
		Status:        UserPlanStatusActive,
		BillingType:   plan.BillingType,
		ModelLimits:   plan.ModelLimits,
		VirtualAmount: plan.VirtualAmount,
	}
	if err := up.Insert(); err != nil {
		t.Fatalf("insert user plan: %v", err)
	}
	return up
}

// TestValidateConfig_BillingType 校验计费维度白名单。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestValidateConfig_BillingType(t *testing.T) {
	valid := `{"gpt-4o": {"period_h": 5, "token_period": 1000}}`
	if err := (&Plan{Name: "a", BillingType: BillingTypeToken, ModelLimits: valid}).ValidateConfig(); err != nil {
		t.Fatalf("token plan should be valid, got %v", err)
	}
	if err := (&Plan{Name: "a", BillingType: "per_request", ModelLimits: valid}).ValidateConfig(); err == nil {
		t.Fatal("per_request should be rejected")
	}
	if err := (&Plan{Name: "a", BillingType: "", ModelLimits: valid}).ValidateConfig(); err == nil {
		t.Fatal("empty billing type should be rejected")
	}
	if err := (&Plan{Name: "a", BillingType: BillingTypeToken, VirtualAmount: -1, ModelLimits: valid}).ValidateConfig(); err == nil {
		t.Fatal("negative virtual amount should be rejected")
	}
}

// TestValidateConfig_ModelLimits 校验 model_limits 的格式与维度匹配。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestValidateConfig_ModelLimits(t *testing.T) {
	cases := []struct {
		name        string
		billingType string
		limits      string
		wantErr     bool
	}{
		{"empty limits", BillingTypeToken, "", true},
		{"empty object", BillingTypeToken, "{}", true},
		{"invalid json", BillingTypeToken, "{oops}", true},
		{"legacy number value", BillingTypeToken, `{"gpt-4": 1000}`, true},
		{"token plan without token limit", BillingTypeToken, `{"gpt-4o": {"period_h": 5, "request_period": 10}}`, true},
		{"request plan without request limit", BillingTypeRequest, `{"gpt-4o": {"period_h": 5, "token_period": 10}}`, true},
		{"valid token plan", BillingTypeToken, `{"gpt-4o": {"period_h": 5, "token_period": 1000}}`, false},
		{"valid request plan", BillingTypeRequest, `{"gpt-4o": {"period_h": 5, "request_period": 10}}`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := (&Plan{Name: "p", BillingType: c.billingType, ModelLimits: c.limits}).ValidateConfig()
			if c.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestCheckPlanQuota_VirtualAmountExhausted 虚拟余额耗尽 → 套餐不可用。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestCheckPlanQuota_VirtualAmountExhausted(t *testing.T) {
	setupPlanQuotaTestDB(t)
	plan := createTestPlan(t, &Plan{
		Name:          "budget",
		BillingType:   BillingTypeToken,
		VirtualAmount: 100_000_000, // 100 元
		ModelLimits:   `{"gpt-4o": {"period_h": 5, "token_period": 1000000}}`,
	})
	up := createTestUserPlan(t, plan, 1)
	if err := IncrementUserPlanUsedAmount(int(up.Id), 100_000_000); err != nil {
		t.Fatalf("increment used amount: %v", err)
	}

	res, err := CheckPlanQuota(1, "gpt-4o")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if res.Usable {
		t.Fatalf("plan should be exhausted, got usable=true")
	}
}

// TestCheckPlanQuota_VirtualAmountRemaining 虚拟余额仍有剩余 → 套餐可用。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestCheckPlanQuota_VirtualAmountRemaining(t *testing.T) {
	setupPlanQuotaTestDB(t)
	plan := createTestPlan(t, &Plan{
		Name:          "budget",
		BillingType:   BillingTypeToken,
		VirtualAmount: 100_000_000,
		ModelLimits:   `{"gpt-4o": {"period_h": 5, "token_period": 1000000}}`,
	})
	up := createTestUserPlan(t, plan, 1)
	if err := IncrementUserPlanUsedAmount(int(up.Id), 30_000_000); err != nil {
		t.Fatalf("increment used amount: %v", err)
	}

	res, err := CheckPlanQuota(1, "gpt-4o")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if !res.Usable {
		t.Fatal("plan should be usable")
	}
	if res.RemainingAmount != 70_000_000 {
		t.Fatalf("remaining = %d, want 70000000", res.RemainingAmount)
	}
	if res.UsedAmount != 30_000_000 {
		t.Fatalf("used = %d, want 30000000", res.UsedAmount)
	}
}

// TestCheckPlanQuota_InvalidJSONNotFree 非法 model_limits 不得被当成「不限量免费」。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestCheckPlanQuota_InvalidJSONNotFree(t *testing.T) {
	setupPlanQuotaTestDB(t)
	plan := createTestPlan(t, &Plan{
		Name:          "broken",
		BillingType:   BillingTypeToken,
		VirtualAmount: 100_000_000,
		ModelLimits:   `{"gpt-4o": 1000}`, // 非法：值应为对象
	})
	createTestUserPlan(t, plan, 1)

	res, err := CheckPlanQuota(1, "gpt-4o")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if res.Usable {
		t.Fatal("invalid model_limits must NOT be usable")
	}
}

// TestCheckPlanQuota_EmptyLimitsNotFree 空 model_limits 不得被当成「不限量免费」。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestCheckPlanQuota_EmptyLimitsNotFree(t *testing.T) {
	setupPlanQuotaTestDB(t)
	plan := createTestPlan(t, &Plan{
		Name:          "empty",
		BillingType:   BillingTypeToken,
		VirtualAmount: 100_000_000,
		ModelLimits:   "",
	})
	createTestUserPlan(t, plan, 1)

	res, err := CheckPlanQuota(1, "gpt-4o")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if res.Usable {
		t.Fatal("empty model_limits must NOT be usable")
	}
}

// TestCheckPlanQuota_ModelNotCovered 未在 model_limits 内的模型不参与套餐。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestCheckPlanQuota_ModelNotCovered(t *testing.T) {
	setupPlanQuotaTestDB(t)
	plan := createTestPlan(t, &Plan{
		Name:          "only-gpt4o",
		BillingType:   BillingTypeToken,
		VirtualAmount: 100_000_000,
		ModelLimits:   `{"gpt-4o": {"period_h": 5, "token_period": 1000000}}`,
	})
	createTestUserPlan(t, plan, 1)

	res, err := CheckPlanQuota(1, "claude-3")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if res.Usable {
		t.Fatal("model outside model_limits must not use the plan")
	}
}

// TestCheckPlanQuota_SnapshotFallback 套餐行被删除后，订阅快照仍然生效。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestCheckPlanQuota_SnapshotFallback(t *testing.T) {
	setupPlanQuotaTestDB(t)
	plan := createTestPlan(t, &Plan{
		Name:          "deletable",
		BillingType:   BillingTypeToken,
		VirtualAmount: 100_000_000,
		ModelLimits:   `{"gpt-4o": {"period_h": 5, "token_period": 1000000}}`,
	})
	up := createTestUserPlan(t, plan, 1)

	// 模拟管理员删除套餐行（历史实现会让订阅直接失败）。
	if err := DB.Delete(&Plan{}, plan.Id).Error; err != nil {
		t.Fatalf("delete plan: %v", err)
	}

	res, err := CheckPlanQuota(1, "gpt-4o")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if !res.Usable {
		t.Fatal("subscription should still work via snapshot after plan deletion")
	}
	if res.VirtualAmount != 100_000_000 {
		t.Fatalf("snapshot virtual amount = %d, want 100000000", res.VirtualAmount)
	}

	// 快照余额耗尽后同样不可用。
	if err := IncrementUserPlanUsedAmount(int(up.Id), 100_000_000); err != nil {
		t.Fatalf("increment used amount: %v", err)
	}
	res, err = CheckPlanQuota(1, "gpt-4o")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if res.Usable {
		t.Fatal("snapshot plan should be exhausted")
	}
}

// TestCheckPlanQuota_UnlimitedVirtualAmount VirtualAmount=0 表示不限额度。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestCheckPlanQuota_UnlimitedVirtualAmount(t *testing.T) {
	setupPlanQuotaTestDB(t)
	plan := createTestPlan(t, &Plan{
		Name:          "unlimited",
		BillingType:   BillingTypeToken,
		VirtualAmount: 0,
		ModelLimits:   `{"gpt-4o": {"period_h": 5, "token_period": 1000000}}`,
	})
	up := createTestUserPlan(t, plan, 1)
	if err := IncrementUserPlanUsedAmount(int(up.Id), 999_000_000); err != nil {
		t.Fatalf("increment used amount: %v", err)
	}

	res, err := CheckPlanQuota(1, "gpt-4o")
	if err != nil {
		t.Fatalf("check plan quota: %v", err)
	}
	if !res.Usable {
		t.Fatal("unlimited plan should stay usable")
	}
	if res.RemainingAmount != -1 {
		t.Fatalf("remaining = %d, want -1 for unlimited", res.RemainingAmount)
	}
}

// TestRemainingVirtualAmount 剩余额度计算的边界。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestRemainingVirtualAmount(t *testing.T) {
	if got := RemainingVirtualAmount(0, 5); got != -1 {
		t.Fatalf("unlimited remaining = %d, want -1", got)
	}
	if got := RemainingVirtualAmount(100, 100); got != 0 {
		t.Fatalf("exhausted remaining = %d, want 0", got)
	}
	if got := RemainingVirtualAmount(100, 150); got != 0 {
		t.Fatalf("over-used remaining = %d, want 0", got)
	}
	if got := RemainingVirtualAmount(100, 40); got != 60 {
		t.Fatalf("remaining = %d, want 60", got)
	}
}

// TestEffectiveBillingType_Fallback 计费维度的回退顺序。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestEffectiveBillingType_Fallback(t *testing.T) {
	up := &UserPlan{BillingType: BillingTypeRequest}
	if got := up.EffectiveBillingType(); got != BillingTypeRequest {
		t.Fatalf("got %s, want request", got)
	}
	up.BillingType = ""
	if got := up.EffectiveBillingType(); got != BillingTypeToken {
		t.Fatalf("got %s, want token fallback", got)
	}
	up.BillingType = "garbage"
	up.Plan = &Plan{BillingType: BillingTypeRequest}
	if got := up.EffectiveBillingType(); got != BillingTypeRequest {
		t.Fatalf("got %s, want plan value", got)
	}
}
