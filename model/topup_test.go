// Topup business unit tests
// 版本: v0.0.22
// 日期: 2026-10-03
// 作者: opencode
//
// 覆盖：
//   - SaveTopupSettings 的去重/非负/正数校验、兑换比例下限校验
//   - GetTopupSettings 的默认值（QuotaPerUnit）与存量脏数据自愈
//   - 赠送倍率（> QuotaPerUnit）保留
//   - ResolveTopupAmount 在 preset 命中、自定义换算、关闭自定义等情况的行为
//   - ActivateTopupByOrder 的幂等

package model

import (
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
)

// setupTopupTestDB 初始化一个 in-memory sqlite + system_settings 表。
// 返回 DB 引用，由调用方赋值给 model.DB。
// 关闭 Redis 缓存，避免本包测试在 RDB 未初始化时 panic（与其他测试一致）。
//
// 版本: v0.0.22
// 日期: 2026-10-03
func setupTopupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	common.RedisEnabled = false
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&SystemSetting{}, &Order{}, &User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// TestSaveTopupSettings_DuplicateAmount 验证金额重复时返回错误。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestSaveTopupSettings_DuplicateAmount(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: 10},
		{Amount: 10, BonusQuota: 20},
	}, 1)
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestSaveTopupSettings_OK 正常保存与读取。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestSaveTopupSettings_OK(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: 10},
		{Amount: 50, BonusQuota: 60},
	}, topupQuotaPerUnit())
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	enabled, allowCustom, presets, rate := GetTopupSettings()
	if !enabled || !allowCustom {
		t.Fatalf("读取异常：enabled=%v allow_custom=%v", enabled, allowCustom)
	}
	if len(presets) != 2 {
		t.Fatalf("期望 2 条 preset，实际 %d", len(presets))
	}
	if rate != topupQuotaPerUnit() {
		t.Fatalf("exchange_rate 期望 %d，实际 %d", topupQuotaPerUnit(), rate)
	}
}

// TestSaveTopupSettings_Defaults 验证默认值。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestSaveTopupSettings_Defaults(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	enabled, allowCustom, presets, rate := GetTopupSettings()
	if enabled || allowCustom {
		t.Fatalf("默认应为关闭：enabled=%v allow_custom=%v", enabled, allowCustom)
	}
	if presets == nil {
		t.Fatal("presets 应返回空切片而不是 nil")
	}
	if rate != topupQuotaPerUnit() {
		t.Fatalf("exchange_rate 默认应为 %d，实际 %d", topupQuotaPerUnit(), rate)
	}
}

// TestSaveTopupSettings_NegativeQuota 验证 bonus_quota 不能为负。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestSaveTopupSettings_NegativeQuota(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: -1},
	}, topupQuotaPerUnit())
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestSaveTopupSettings_ZeroAmount 验证金额必须大于 0。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestSaveTopupSettings_ZeroAmount(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 0, BonusQuota: 10},
	}, topupQuotaPerUnit())
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestResolveTopupAmount_PresetHit 预设金额命中。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestResolveTopupAmount_PresetHit(t *testing.T) {
	presets := []TopupPreset{
		{Amount: 10, BonusQuota: 10},
		{Amount: 50, BonusQuota: 60},
	}
	amt, bonus, err := ResolveTopupAmount(CreateTopupOrderInput{PresetAmount: 50}, presets, true, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if amt != 50 || bonus != 60 {
		t.Fatalf("amt=%v bonus=%v", amt, bonus)
	}
}

// TestResolveTopupAmount_PresetMiss 预设金额不存在时返回错误。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestResolveTopupAmount_PresetMiss(t *testing.T) {
	presets := []TopupPreset{{Amount: 10, BonusQuota: 10}}
	_, _, err := ResolveTopupAmount(CreateTopupOrderInput{PresetAmount: 99}, presets, true, 1)
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestResolveTopupAmount_CustomDisabled 自定义关闭时拒绝。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestResolveTopupAmount_CustomDisabled(t *testing.T) {
	_, _, err := ResolveTopupAmount(CreateTopupOrderInput{Amount: 25}, nil, false, 1)
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestResolveTopupAmount_CustomOK 自定义金额按基准单位换算（1 元 = QuotaPerUnit 额度）。
// 版本: v0.0.22
// 日期: 2026-10-03
func TestResolveTopupAmount_CustomOK(t *testing.T) {
	base := topupQuotaPerUnit()
	amt, bonus, err := ResolveTopupAmount(CreateTopupOrderInput{Amount: 25}, nil, true, base)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if amt != 25 || bonus != 25*base {
		t.Fatalf("amt=%v bonus=%v（期望 25 元 → %d 额度）", amt, bonus, 25*base)
	}
}

// TestResolveTopupAmount_CustomRate 自定义金额按 exchange_rate 换算（10 元 × 500000）。
// 版本: v0.0.22
// 日期: 2026-10-03
func TestResolveTopupAmount_CustomRate(t *testing.T) {
	amt, bonus, err := ResolveTopupAmount(CreateTopupOrderInput{Amount: 10}, nil, true, 500000)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if amt != 10 || bonus != 5_000_000 {
		t.Fatalf("amt=%v bonus=%v（期望 5000000）", amt, bonus)
	}
}

// TestResolveTopupAmount_BelowBaseHealed 低于基准单位的比例被自愈为基准值。
// 版本: v0.0.22
// 日期: 2026-10-03
func TestResolveTopupAmount_BelowBaseHealed(t *testing.T) {
	base := topupQuotaPerUnit()
	amt, bonus, err := ResolveTopupAmount(CreateTopupOrderInput{Amount: 10}, nil, true, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if amt != 10 || bonus != 10*base {
		t.Fatalf("amt=%v bonus=%v（期望 10 元 → %d 额度）", amt, bonus, 10*base)
	}
}

// TestGetTopupSettings_LegacyRateSelfHeal 存量脏数据（1 / 100000）读取时自愈为基准值。
// 版本: v0.0.22
// 日期: 2026-10-03
func TestGetTopupSettings_LegacyRateSelfHeal(t *testing.T) {
	for _, legacy := range []string{"1", "100000"} {
		db := setupTopupTestDB(t)
		DB = db
		if err := UpsertSystemSetting(SystemSettingKeyTopupExchangeRate, legacy,
			SystemSettingCategoryTopup, "legacy"); err != nil {
			t.Fatalf("seed legacy %s: %v", legacy, err)
		}
		_, _, _, rate := GetTopupSettings()
		if rate != topupQuotaPerUnit() {
			t.Fatalf("legacy %s 应自愈为 %d，实际 %d", legacy, topupQuotaPerUnit(), rate)
		}
	}
}

// TestGetTopupSettings_BonusMultiplierKept 高于基准单位的比例视为赠送倍率并保留。
// 版本: v0.0.22
// 日期: 2026-10-03
func TestGetTopupSettings_BonusMultiplierKept(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db
	multiplier := topupQuotaPerUnit() * 2
	if err := UpsertSystemSetting(SystemSettingKeyTopupExchangeRate,
		fmt.Sprintf("%d", multiplier), SystemSettingCategoryTopup, "bonus"); err != nil {
		t.Fatalf("seed bonus: %v", err)
	}
	_, _, _, rate := GetTopupSettings()
	if rate != multiplier {
		t.Fatalf("赠送倍率应保留 %d，实际 %d", multiplier, rate)
	}
}

// TestSaveTopupSettings_BelowBaseRejected 低于基准单位直接拒绝。
// 版本: v0.0.22
// 日期: 2026-10-03
func TestSaveTopupSettings_BelowBaseRejected(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db
	if err := SaveTopupSettings(true, true, nil, 1); err == nil {
		t.Fatal("期望拒绝低于基准单位的兑换比例")
	}
}

// TestActivateTopupByOrder_Idempotent 幂等：第二次调用不报错、不变更时间。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestActivateTopupByOrder_Idempotent(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	// 保存设置（开启）
	if err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: 100},
	}, topupQuotaPerUnit()); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	// 准备一个用户（ActivateTopupByOrder 内部会调用 IncreaseUserQuota）
	u := &User{Id: 1, Username: "topup-test", Quota: 0, Status: UserStatusEnabled}
	if err := DB.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	order, _, _, err := CreateTopupOrder(CreateTopupOrderInput{
		UserId:       1,
		PresetAmount: 10,
		PayMethod:    OrderPayMethodWechat,
		Source:       OrderSourceUserSelf,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if order.Type != OrderTypeTopup {
		t.Fatalf("订单类型应为 %d，实际 %d", OrderTypeTopup, order.Type)
	}

	if err := ActivateTopupByOrder(order); err != nil {
		t.Fatalf("first activate: %v", err)
	}
	firstPaidTime := order.PayTime

	// 二次调用应直接返回 nil
	if err := ActivateTopupByOrder(order); err != nil {
		t.Fatalf("second activate should be no-op: %v", err)
	}
	if order.PayTime != firstPaidTime {
		t.Fatal("幂等失败：第二次调用修改了 PayTime")
	}
}

// TestActivateTopupByOrder_WrongType 非充值订单拒绝。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestActivateTopupByOrder_WrongType(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	order := &Order{Type: OrderTypePlanSubscription, Status: OrderStatusPending}
	err := ActivateTopupByOrder(order)
	if err == nil {
		t.Fatal("期望对非充值订单返回错误")
	}
}
