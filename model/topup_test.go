// Topup business unit tests
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
//
// 覆盖：
//   - SaveTopupSettings 的去重/正数/「到账不低于支付」校验
//   - GetTopupSettings 的默认值（关闭、空 presets）
//   - ResolveTopupAmount 在 preset 命中、自定义 1:1、关闭自定义等情况的行为
//   - 自定义金额与快捷金额的赠送能力（到账 > 支付）
//   - ActivateTopupByOrder 的幂等

package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/config"
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

// TestYuanToQuota_RoundsFloatError 验证浮点截断误差被四舍五入消除。
// 0.07 元 × 1e6 在 IEEE 754 下约为 69999.999…，直接取整会少 1。
//
// 版本: v0.0.24
// 日期: 2026-10-03
func TestYuanToQuota_RoundsFloatError(t *testing.T) {
	if got := YuanToQuota(0.07); got != 70_000 {
		t.Fatalf("YuanToQuota(0.07) = %d，期望 70000", got)
	}
	if got := YuanToQuota(1); got != int64(config.QuotaPerUnit) {
		t.Fatalf("YuanToQuota(1) = %d，期望 %d", got, int64(config.QuotaPerUnit))
	}
	if got := YuanToQuota(0); got != 0 {
		t.Fatalf("YuanToQuota(0) = %d，期望 0", got)
	}
	if got := YuanToQuota(-1); got != 0 {
		t.Fatalf("YuanToQuota(-1) = %d，期望 0", got)
	}
}

// TestSaveTopupSettings_DuplicateAmount 验证金额重复时返回错误。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestSaveTopupSettings_DuplicateAmount(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: YuanToQuota(10)},
		{Amount: 10, BonusQuota: YuanToQuota(20)},
	})
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestSaveTopupSettings_OK 正常保存与读取（含「充 10 得 15」赠送）。
// 版本: v0.0.24
// 日期: 2026-10-03
func TestSaveTopupSettings_OK(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: YuanToQuota(15)},
		{Amount: 50, BonusQuota: YuanToQuota(50)},
	})
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	enabled, allowCustom, presets := GetTopupSettings()
	if !enabled || !allowCustom {
		t.Fatalf("读取异常：enabled=%v allow_custom=%v", enabled, allowCustom)
	}
	if len(presets) != 2 {
		t.Fatalf("期望 2 条 preset，实际 %d", len(presets))
	}
	if presets[0].BonusQuota != YuanToQuota(15) {
		t.Fatalf("赠送生效失败：期望到账 %d，实际 %d", YuanToQuota(15), presets[0].BonusQuota)
	}
}

// TestSaveTopupSettings_Defaults 验证默认值：关闭、presets 为空切片。
// 版本: v0.0.24
// 日期: 2026-10-03
func TestSaveTopupSettings_Defaults(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	enabled, allowCustom, presets := GetTopupSettings()
	if enabled || allowCustom {
		t.Fatalf("默认应为关闭：enabled=%v allow_custom=%v", enabled, allowCustom)
	}
	if presets == nil {
		t.Fatal("presets 应返回空切片而不是 nil")
	}
	if len(presets) != 0 {
		t.Fatalf("默认 presets 应为空，实际 %d", len(presets))
	}
}

// TestSaveTopupSettings_BonusBelowAmount 到账额度低于支付金额折算额度时拒绝（不允许缩水）。
// 版本: v0.0.24
// 日期: 2026-10-03
func TestSaveTopupSettings_BonusBelowAmount(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: YuanToQuota(10) - 1},
	})
	if err == nil {
		t.Fatal("期望拒绝到账额度低于支付金额的配置")
	}
	// 恰好等于折算额度应通过（1:1）
	if err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: YuanToQuota(10)},
	}); err != nil {
		t.Fatalf("1:1 配置应通过，实际报错：%v", err)
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
	})
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestResolveTopupAmount_PresetHit 预设金额命中（含赠送倍率）。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestResolveTopupAmount_PresetHit(t *testing.T) {
	presets := []TopupPreset{
		{Amount: 10, BonusQuota: YuanToQuota(15)},
		{Amount: 50, BonusQuota: YuanToQuota(50)},
	}
	amt, bonus, err := ResolveTopupAmount(CreateTopupOrderInput{PresetAmount: 50}, presets, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if amt != 50 || bonus != YuanToQuota(50) {
		t.Fatalf("amt=%v bonus=%v", amt, bonus)
	}
}

// TestResolveTopupAmount_PresetMiss 预设金额不存在时返回错误。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestResolveTopupAmount_PresetMiss(t *testing.T) {
	presets := []TopupPreset{{Amount: 10, BonusQuota: YuanToQuota(10)}}
	_, _, err := ResolveTopupAmount(CreateTopupOrderInput{PresetAmount: 99}, presets, true)
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestResolveTopupAmount_CustomDisabled 自定义关闭时拒绝。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestResolveTopupAmount_CustomDisabled(t *testing.T) {
	_, _, err := ResolveTopupAmount(CreateTopupOrderInput{Amount: 25}, nil, false)
	if err == nil {
		t.Fatal("期望返回错误，但通过了")
	}
}

// TestResolveTopupAmount_CustomOK 自定义金额恒 1:1（1 元 = QuotaPerUnit 额度）。
// 版本: v0.0.24
// 日期: 2026-10-03
func TestResolveTopupAmount_CustomOK(t *testing.T) {
	amt, bonus, err := ResolveTopupAmount(CreateTopupOrderInput{Amount: 25}, nil, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if amt != 25 || bonus != YuanToQuota(25) {
		t.Fatalf("amt=%v bonus=%v（期望 25 元 → %d 额度）", amt, bonus, YuanToQuota(25))
	}
}

// TestResolveTopupAmount_CustomZeroRejected 自定义金额必须大于 0。
// 版本: v0.0.24
// 日期: 2026-10-03
func TestResolveTopupAmount_CustomZeroRejected(t *testing.T) {
	if _, _, err := ResolveTopupAmount(CreateTopupOrderInput{Amount: 0}, nil, true); err == nil {
		t.Fatal("期望拒绝 0 元自定义金额")
	}
}

// TestActivateTopupByOrder_Idempotent 幂等：第二次调用不报错、不变更时间。
// 版本: v0.0.10
// 日期: 2026-09-06
func TestActivateTopupByOrder_Idempotent(t *testing.T) {
	db := setupTopupTestDB(t)
	DB = db

	// 保存设置（开启，充 10 得 15）
	if err := SaveTopupSettings(true, true, []TopupPreset{
		{Amount: 10, BonusQuota: YuanToQuota(15)},
	}); err != nil {
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
