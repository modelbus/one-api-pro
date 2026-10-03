// topup_amount_test 充值金额一致性测试：订单金额（元）↔ 到账额度 ↔ 用户余额（元）
// Top-up amount consistency: order amount (CNY) <-> credited quota <-> user balance (CNY)
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
//
// 覆盖：
//   - 快捷金额「充 10 得 15」：订单支付 ¥10、快照 credit_amount=15、到账 15_000_000 quota
//   - 自定义金额恒 1:1
//   - 历史订单多余的 exchange_rate 字段被安全忽略（按快照 bonus_quota 发放）
//   - 快照缺失时按 1:1（YuanToQuota(order.Amount)）兜底

package model

import (
	"encoding/json"
	"testing"

	"github.com/modelbus/one-api-pro/common/config"
)

// quotaToYuanForTest 用系统常量把 quota 折算为元（测试内联，避免依赖 relay 包）。
func quotaToYuanForTest(quota int64) float64 {
	return float64(quota) / config.QuotaPerUnit
}

func readUserQuota(t *testing.T, id int) int64 {
	t.Helper()
	var u User
	if err := DB.First(&u, "id = ?", id).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	return u.Quota
}

func seedTopupUser(t *testing.T, id int, quota int64) {
	t.Helper()
	u := &User{Id: id, Username: "topup-amount", Quota: quota, Status: UserStatusEnabled}
	if err := DB.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
}

// TestTopupOrder_PresetBonus_AmountAndCredit 验证快捷金额的「支付金额 / 到账金额」双列口径。
// 充 10 得 15：订单支付额为 ¥10，到账为 15_000_000 quota（= ¥15）。
// 版本: v0.0.24
func TestTopupOrder_PresetBonus_AmountAndCredit(t *testing.T) {
	DB = setupTopupTestDB(t)

	if err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: YuanToQuota(15)},
	}); err != nil {
		t.Fatalf("save settings: %v", err)
	}
	seedTopupUser(t, 1, 0)

	order, payAmount, bonusQuota, err := CreateTopupOrder(CreateTopupOrderInput{
		UserId:       1,
		PresetAmount: 10,
		PayMethod:    OrderPayMethodWechat,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	// 支付金额仍是 ¥10（不得被赠送改写）
	if payAmount != 10 || order.Amount != 10 {
		t.Fatalf("支付金额 = %v / order.Amount = %v，期望 10", payAmount, order.Amount)
	}
	// 到账额度 = ¥15
	if bonusQuota != YuanToQuota(15) {
		t.Fatalf("到账额度 = %d，期望 %d", bonusQuota, YuanToQuota(15))
	}
	if quotaToYuanForTest(bonusQuota) != 15 {
		t.Fatalf("到账金额 = ¥%v，期望 ¥15", quotaToYuanForTest(bonusQuota))
	}

	// 订单快照 JSON 口径
	var info TopupOrderPlanInfo
	if err := json.Unmarshal([]byte(order.PlanInfo), &info); err != nil {
		t.Fatalf("unmarshal plan_info: %v", err)
	}
	if info.Amount != 10 || info.PresetAmount != 10 {
		t.Fatalf("快照 amount=%v preset_amount=%v，期望 10/10", info.Amount, info.PresetAmount)
	}
	if info.CreditAmount != 15 {
		t.Fatalf("快照 credit_amount = %v，期望 15", info.CreditAmount)
	}
	if info.BonusQuota != YuanToQuota(15) {
		t.Fatalf("快照 bonus_quota = %d，期望 %d", info.BonusQuota, YuanToQuota(15))
	}

	// 激活后用户余额 = ¥15
	if err := ActivateTopupByOrder(order); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if got := readUserQuota(t, 1); got != YuanToQuota(15) {
		t.Fatalf("余额 = %d quota（¥%v），期望 %d（¥15）", got, quotaToYuanForTest(got), YuanToQuota(15))
	}
}

// TestTopupOrder_CustomAmount_OneToOne 自定义金额恒 1:1。
// 版本: v0.0.24
func TestTopupOrder_CustomAmount_OneToOne(t *testing.T) {
	DB = setupTopupTestDB(t)
	if err := SaveTopupSettings(true, true, nil); err != nil {
		t.Fatalf("save settings: %v", err)
	}
	seedTopupUser(t, 1, 0)

	order, payAmount, bonusQuota, err := CreateTopupOrder(CreateTopupOrderInput{
		UserId:    1,
		Amount:    25,
		PayMethod: OrderPayMethodAlipay,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if payAmount != 25 || bonusQuota != YuanToQuota(25) {
		t.Fatalf("pay=%v credit=%d，期望 25 / %d", payAmount, bonusQuota, YuanToQuota(25))
	}
	if err := ActivateTopupByOrder(order); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if got := readUserQuota(t, 1); got != YuanToQuota(25) {
		t.Fatalf("余额 = %d，期望 %d", got, YuanToQuota(25))
	}
}

// TestActivateTopup_LegacyExchangeRateIgnored 历史订单快照带 exchange_rate 时应被忽略，
// 仍以 bonus_quota 为到账依据（避免旧倍率被重新应用导致多发）。
// 版本: v0.0.24
func TestActivateTopup_LegacyExchangeRateIgnored(t *testing.T) {
	DB = setupTopupTestDB(t)
	seedTopupUser(t, 1, 0)

	legacyJSON := `{"amount":10,"credit_amount":15,"bonus_quota":15000000,"exchange_rate":2}`
	order := &Order{
		Type:      OrderTypeTopup,
		Source:    OrderSourceUserSelf,
		OrderNo:   "TP-LEGACY-1",
		UserId:    1,
		PlanInfo:  legacyJSON,
		Amount:    10,
		Status:    OrderStatusPending,
		PayStatus: OrderPayStatusPending,
		PayMethod: OrderPayMethodWechat,
	}
	if err := order.Insert(); err != nil {
		t.Fatalf("insert order: %v", err)
	}
	if err := ActivateTopupByOrder(order); err != nil {
		t.Fatalf("activate: %v", err)
	}
	// 到账应为快照 bonus_quota = ¥15，而非 amount × exchange_rate = ¥20
	if got := readUserQuota(t, 1); got != YuanToQuota(15) {
		t.Fatalf("余额 = %d（¥%v），期望 %d（¥15）", got, quotaToYuanForTest(got), YuanToQuota(15))
	}
}

// TestActivateTopup_MissingSnapshotFallback 快照缺失时按 1:1 兜底。
// 版本: v0.0.24
func TestActivateTopup_MissingSnapshotFallback(t *testing.T) {
	DB = setupTopupTestDB(t)
	seedTopupUser(t, 1, 0)

	order := &Order{
		Type:      OrderTypeTopup,
		Source:    OrderSourceUserSelf,
		OrderNo:   "TP-FALLBACK-1",
		UserId:    1,
		PlanInfo:  "",
		Amount:    10,
		Status:    OrderStatusPending,
		PayStatus: OrderPayStatusPending,
		PayMethod: OrderPayMethodWechat,
	}
	if err := order.Insert(); err != nil {
		t.Fatalf("insert order: %v", err)
	}
	if err := ActivateTopupByOrder(order); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if got := readUserQuota(t, 1); got != YuanToQuota(10) {
		t.Fatalf("余额 = %d，期望 %d（按 ¥10 兜底）", got, YuanToQuota(10))
	}
}
