package model

import (
	"errors"

	"github.com/modelbus/one-api-pro/common/helper"
)

// SystemSetting category constants
const (
	SystemSettingCategoryPayment = "payment" // 支付配置
	SystemSettingCategoryPlan    = "plan"    // 套餐运营
	SystemSettingCategoryTopup   = "topup"   // 充值设置（在线支付充值）
	SystemSettingCategoryGeneral = "general" // 通用
)

// SystemSetting key constants — payment.*
const (
	SystemSettingKeyWechatEnabled = "payment.wechat.enabled"
	SystemSettingKeyWechatConfig  = "payment.wechat.config"
	SystemSettingKeyAlipayEnabled = "payment.alipay.enabled"
	SystemSettingKeyAlipayConfig  = "payment.alipay.config"
	SystemSettingKeyBankEnabled   = "payment.bank.enabled"
	SystemSettingKeyBankConfig    = "payment.bank.config"
)

// SystemSetting key constants — plan.*
const (
	SystemSettingKeyPlanUpgradeMode = "plan.upgrade_mode"
	// SystemSettingKeyPlanAllowTopup 已迁移至 topup.enabled，DB 行保留不动
)

// SystemSetting key constants — topup.*（用户在线支付充值）
// 版本 v0.0.10  2026-09-06  新增
const (
	SystemSettingKeyTopupEnabled      = "topup.enabled"       // 充值总开关（bool: true/false）
	SystemSettingKeyTopupAllowCustom  = "topup.allow_custom"  // 是否允许用户输入自定义金额（bool）
	SystemSettingKeyTopupPresets      = "topup.presets"       // 快捷金额列表（JSON: [{amount, bonus_quota}, ...]）
	SystemSettingKeyTopupExchangeRate = "topup.exchange_rate" // 自定义金额 1 元 = X quota（int，默认 1 即 1:1）
)

// SystemSetting is a generic key-value configuration row.
// category is used by the API to group settings (payment / plan / general).
type SystemSetting struct {
	Key         string `gorm:"primaryKey;type:varchar(64)" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Category    string `gorm:"type:varchar(32);not null;default:'general';index:idx_category" json:"category"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	UpdatedAt   int64  `gorm:"bigint;default:0" json:"updated_at"`
}

// UpsertSystemSetting inserts or updates a setting row.
func UpsertSystemSetting(key, value, category, description string) error {
	if key == "" {
		return errors.New("key 不能为空")
	}
	now := helper.GetTimestamp()
	setting := SystemSetting{
		Key:         key,
		Value:       value,
		Category:    category,
		Description: description,
		UpdatedAt:   now,
	}
	// GORM upsert: use ON CONFLICT (SQLite / PG) / ON DUPLICATE KEY UPDATE (MySQL).
	return DB.Save(&setting).Error
}

func GetSystemSetting(key string) (*SystemSetting, error) {
	if key == "" {
		return nil, errors.New("key 不能为空")
	}
	s := SystemSetting{Key: key}
	if err := DB.First(&s, "`key` = ?", key).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSystemSettingsByCategory(category string) ([]*SystemSetting, error) {
	var ss []*SystemSetting
	if err := DB.Where("category = ?", category).Order("`key` asc").Find(&ss).Error; err != nil {
		return nil, err
	}
	return ss, nil
}

func GetAllSystemSettings() ([]*SystemSetting, error) {
	var ss []*SystemSetting
	if err := DB.Order("category asc, `key` asc").Find(&ss).Error; err != nil {
		return nil, err
	}
	return ss, nil
}

func DeleteSystemSetting(key string) error {
	return DB.Where("`key` = ?", key).Delete(&SystemSetting{}).Error
}

// GetSystemSettingString returns the Value for key, or "" if missing.
// Used by callers that just need a plain string regardless of existence.
func GetSystemSettingString(key string) string {
	s, err := GetSystemSetting(key)
	if err != nil {
		return ""
	}
	return s.Value
}
