// Package model provides data persistence and business logic.
// 本文件实现"用户在线支付充值（topup）"业务：预设金额、自定义金额、订单生成、到账发放。
//
// 命名沿用历史 topup 命名（与 LogTypeTopup/AdminTopUp/TopUp 等保持全局一致），
// 仅 UI 文案对外显示为"充值"。
//
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/modelbus/one-api-pro/common/helper"
)

// TopupPreset 表示一个充值快捷金额配置项。
// amount: 用户需支付的金额（元，浮点，保留两位小数）
// bonus_quota: 用户实际获得的系统额度（quota 整数）
//
// 版本: v0.0.10
// 日期: 2026-09-06
type TopupPreset struct {
	Amount     float64 `json:"amount"`
	BonusQuota int64   `json:"bonus_quota"`
}

// CreateTopupOrderInput 创建充值订单的入参。
// Amount > 0 时表示使用自定义金额；PresetAmount > 0 时表示走预设金额并
// 自动使用对应 preset 的 bonus_quota。两者至少传一个；PresetAmount 命中
// 的 preset 必须存在于系统设置中。
//
// 版本: v0.0.10
// 日期: 2026-09-06
type CreateTopupOrderInput struct {
	UserId       int
	Amount       float64
	PresetAmount float64
	PayMethod    string
	Source       int // OrderSourceUserSelf 或 OrderSourceAdmin
}

// TopupOrderPlanInfo 写入 orders.plan_info 字段的 JSON 结构。
// 用于对账时还原下单时的金额与配额快照。
//
// 版本: v0.0.10
// 日期: 2026-09-06
type TopupOrderPlanInfo struct {
	Amount       float64 `json:"amount"`
	PresetAmount float64 `json:"preset_amount,omitempty"`
	BonusQuota   int64   `json:"bonus_quota"`
	ExchangeRate int64   `json:"exchange_rate"`
}

// defaultTopupExchangeRate 默认自定义金额换算比例 1:1（即 1 元 = 1 quota）。
// 版本: v0.0.10
const defaultTopupExchangeRate int64 = 1

// GetTopupSettings 读取充值设置。返回结构化对象，便于上层直接使用。
// presets 为空时返回空切片（不是 nil）。
//
// 版本: v0.0.10
// 日期: 2026-09-06
func GetTopupSettings() (enabled bool, allowCustom bool, presets []TopupPreset, exchangeRate int64) {
	enabled = GetSystemSettingString(SystemSettingKeyTopupEnabled) == "true"
	allowCustom = GetSystemSettingString(SystemSettingKeyTopupAllowCustom) == "true"
	exchangeRate = parseInt64Setting(SystemSettingKeyTopupExchangeRate, defaultTopupExchangeRate)

	raw := GetSystemSettingString(SystemSettingKeyTopupPresets)
	if raw == "" {
		presets = []TopupPreset{}
		return
	}
	var list []TopupPreset
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		presets = []TopupPreset{}
		return
	}
	presets = list
	return
}

// SaveTopupSettings 整体覆盖保存充值设置。
//
// 版本: v0.0.10
// 日期: 2026-09-06
func SaveTopupSettings(enabled, allowCustom bool, presets []TopupPreset, exchangeRate int64) error {
	if presets == nil {
		presets = []TopupPreset{}
	}
	seenAmounts := make(map[float64]bool, len(presets))
	for i, p := range presets {
		if p.Amount <= 0 {
			return fmt.Errorf("第 %d 项金额必须大于 0", i+1)
		}
		if p.BonusQuota < 0 {
			return fmt.Errorf("第 %d 项额度不能为负数", i+1)
		}
		// 快捷金额金额不允许重复（v0.0.10 2026-09-06 新增校验）
		if seenAmounts[p.Amount] {
			return fmt.Errorf("快捷金额重复：%.2f 元已存在", p.Amount)
		}
		seenAmounts[p.Amount] = true
	}
	if exchangeRate <= 0 {
		return fmt.Errorf("兑换比例必须大于 0")
	}
	presetsJSON, err := json.Marshal(presets)
	if err != nil {
		return err
	}
	if err := UpsertSystemSetting(SystemSettingKeyTopupEnabled,
		fmt.Sprintf("%t", enabled), SystemSettingCategoryTopup, "在线充值总开关"); err != nil {
		return err
	}
	if err := UpsertSystemSetting(SystemSettingKeyTopupAllowCustom,
		fmt.Sprintf("%t", allowCustom), SystemSettingCategoryTopup, "允许自定义金额"); err != nil {
		return err
	}
	if err := UpsertSystemSetting(SystemSettingKeyTopupPresets,
		string(presetsJSON), SystemSettingCategoryTopup, "充值快捷金额"); err != nil {
		return err
	}
	return UpsertSystemSetting(SystemSettingKeyTopupExchangeRate,
		fmt.Sprintf("%d", exchangeRate), SystemSettingCategoryTopup, "自定义金额 1 元 = X quota")
}

// ResolveTopupAmount 根据入参解析实际支付金额与到账 quota。
// 返回值：payAmount(元), bonusQuota(quota), presetMatched(bool)
//
// 版本: v0.0.10
// 日期: 2026-09-06
func ResolveTopupAmount(in CreateTopupOrderInput, presets []TopupPreset, allowCustom bool, exchangeRate int64) (float64, int64, error) {
	if in.PresetAmount > 0 {
		// 命中预设：从列表中找到对应金额的 preset
		for _, p := range presets {
			if p.Amount == in.PresetAmount {
				return p.Amount, p.BonusQuota, nil
			}
		}
		return 0, 0, errors.New("快捷金额未配置")
	}
	// 自定义金额
	if !allowCustom {
		return 0, 0, errors.New("未开启自定义金额")
	}
	if in.Amount <= 0 {
		return 0, 0, errors.New("充值金额必须大于 0")
	}
	if exchangeRate <= 0 {
		exchangeRate = defaultTopupExchangeRate
	}
	bonus := int64(in.Amount * float64(exchangeRate))
	return in.Amount, bonus, nil
}

// CreateTopupOrder 创建充值订单（type=2）。
// 不调用支付预下单（由 controller 调用 buildPayInfo 完成）。
//
// 版本: v0.0.10
// 日期: 2026-09-06
func CreateTopupOrder(in CreateTopupOrderInput) (*Order, float64, int64, error) {
	if in.UserId == 0 {
		return nil, 0, 0, errors.New("user_id 不能为空")
	}
	if in.PayMethod == "" {
		return nil, 0, 0, errors.New("pay_method 不能为空")
	}
	if in.Source == 0 {
		in.Source = OrderSourceUserSelf
	}

	enabled, allowCustom, presets, exchangeRate := GetTopupSettings()
	if !enabled {
		return nil, 0, 0, errors.New("充值功能未开启")
	}

	payAmount, bonusQuota, err := ResolveTopupAmount(in, presets, allowCustom, exchangeRate)
	if err != nil {
		return nil, 0, 0, err
	}

	now := helper.GetTimestamp()
	orderNo, err := GenerateOrderNo("TP")
	if err != nil {
		return nil, 0, 0, err
	}

	info := TopupOrderPlanInfo{
		Amount:       payAmount,
		PresetAmount: in.PresetAmount,
		BonusQuota:   bonusQuota,
		ExchangeRate: exchangeRate,
	}
	infoJSON, _ := json.Marshal(info)

	order := &Order{
		Type:       OrderTypeTopup,
		Source:     in.Source,
		OrderNo:    orderNo,
		UserId:     in.UserId,
		PlanId:     0,
		PlanInfo:   string(infoJSON),
		Amount:     payAmount,
		Status:     OrderStatusPending,
		PayStatus:  OrderPayStatusPending,
		PayMethod:  in.PayMethod,
		CreateTime: now,
		UpdateTime: now,
	}
	if err := order.Insert(); err != nil {
		return nil, 0, 0, err
	}
	return order, payAmount, bonusQuota, nil
}

// ActivateTopupByOrder 标记充值订单已支付并给用户加 quota。
// 幂等：订单已支付直接返回 nil。
// 退款（status=3）暂不回退 quota，留 TODO 等后续统一订单管理功能。
//
// 版本: v0.0.10
// 日期: 2026-09-06
func ActivateTopupByOrder(order *Order) error {
	if order == nil {
		return errors.New("order 不能为空")
	}
	if order.Type != OrderTypeTopup {
		return errors.New("非充值订单")
	}
	if order.Status == OrderStatusPaid {
		return nil // already activated
	}

	// 解析快照获取到账 quota；若快照缺失则按 amount × exchange_rate 兜底
	bonusQuota := int64(0)
	if info := parseTopupOrderPlanInfo(order.PlanInfo); info != nil {
		bonusQuota = info.BonusQuota
	} else {
		_, _, _, exchangeRate := GetTopupSettings()
		if exchangeRate <= 0 {
			exchangeRate = defaultTopupExchangeRate
		}
		bonusQuota = int64(order.Amount * float64(exchangeRate))
	}

	if bonusQuota > 0 {
		if err := IncreaseUserQuota(order.UserId, bonusQuota); err != nil {
			return fmt.Errorf("增加用户额度失败: %w", err)
		}
	}

	if err := order.MarkOrderPaid(order.PayMethod, order.PayTradeNo); err != nil {
		return err
	}

	// TODO: 退款时回退 quota（待后续统一订单管理功能实现）
	return nil
}

func parseTopupOrderPlanInfo(raw string) *TopupOrderPlanInfo {
	if raw == "" {
		return nil
	}
	var info TopupOrderPlanInfo
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return nil
	}
	return &info
}

// parseInt64Setting 读取一个 int64 系统设置；缺失或解析失败时返回 def。
//
// 版本: v0.0.10
// 日期: 2026-09-06
func parseInt64Setting(key string, def int64) int64 {
	raw := GetSystemSettingString(key)
	if raw == "" {
		return def
	}
	var v int64
	if err := json.Unmarshal([]byte(raw), &v); err == nil {
		return v
	}
	// 兼容纯数字字符串
	var n int64
	if _, err := fmt.Sscanf(raw, "%d", &n); err == nil {
		return n
	}
	return def
}
