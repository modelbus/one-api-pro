// Package model provides data persistence and business logic.
// 本文件实现"用户在线支付充值（topup）"业务：预设金额、自定义金额、订单生成、到账发放。
//
// 命名沿用历史 topup 命名（与 LogTypeTopup/AdminTopUp/TopUp 等保持全局一致），
// 仅 UI 文案对外显示为"充值"。
//
// 口径约定：1 元 = config.QuotaPerUnit 额度（常量 1_000_000），
// 即「充值支付金额」与「到账额度」1:1 对齐；赠送能力通过快捷金额的到账金额体现。
//
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/helper"
)

// TopupPreset 表示一个充值快捷金额配置项。
// amount: 用户需支付的金额（元，浮点，保留两位小数）
// bonus_quota: 用户实际到账的系统额度（quota 整数，>= amount 折算额度）
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
// 用于对账时还原下单时的金额与额度快照。
// CreditAmount 为到账金额（元），与 BonusQuota 等价表述，便于人工对账；
// 历史订单中多余的 exchange_rate 字段会被 JSON 解码器自动忽略（安全兼容）。
//
// 版本: v0.0.24
// 日期: 2026-10-03
type TopupOrderPlanInfo struct {
	Amount       float64 `json:"amount"`
	PresetAmount float64 `json:"preset_amount,omitempty"`
	CreditAmount float64 `json:"credit_amount"`
	BonusQuota   int64   `json:"bonus_quota"`
}

// YuanToQuota 将金额（元）按系统常量基准换算为额度（quota）。
// 1 元 = config.QuotaPerUnit 额度；使用 math.Round 消除浮点截断误差
// （例如 0.07 元 × 1e6 在浮点下为 69999.99…，直接取整会少 1）。
//
// 版本: v0.0.24
// 日期: 2026-10-03
func YuanToQuota(amount float64) int64 {
	if amount <= 0 {
		return 0
	}
	return int64(math.Round(amount * config.QuotaPerUnit))
}

// GetTopupSettings 读取充值设置。返回结构化对象，便于上层直接使用。
// presets 为空时返回空切片（不是 nil）。
//
// 版本: v0.0.24
// 日期: 2026-10-03
func GetTopupSettings() (enabled bool, allowCustom bool, presets []TopupPreset) {
	enabled = GetSystemSettingString(SystemSettingKeyTopupEnabled) == "true"
	allowCustom = GetSystemSettingString(SystemSettingKeyTopupAllowCustom) == "true"

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
// 校验规则：金额必须大于 0；金额不允许重复；
// 到账额度不得低于支付金额按基准折算的额度（允许赠送，不允许缩水）。
//
// 版本: v0.0.24
// 日期: 2026-10-03
func SaveTopupSettings(enabled, allowCustom bool, presets []TopupPreset) error {
	if presets == nil {
		presets = []TopupPreset{}
	}
	seenAmounts := make(map[float64]bool, len(presets))
	for i, p := range presets {
		if p.Amount <= 0 {
			return fmt.Errorf("第 %d 项金额必须大于 0", i+1)
		}
		// 快捷金额金额不允许重复（v0.0.10 2026-09-06 新增校验）
		if seenAmounts[p.Amount] {
			return fmt.Errorf("快捷金额重复：%.2f 元已存在", p.Amount)
		}
		seenAmounts[p.Amount] = true
		// 到账额度不得低于支付金额折算额度（营销赠送允许高于）。
		minQuota := YuanToQuota(p.Amount)
		if p.BonusQuota < minQuota {
			return fmt.Errorf("第 %d 项到账额度不能低于支付金额折算额度（%.2f 元 = %d 额度）",
				i+1, p.Amount, minQuota)
		}
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
	return UpsertSystemSetting(SystemSettingKeyTopupPresets,
		string(presetsJSON), SystemSettingCategoryTopup, "充值快捷金额（支付金额 / 到账金额）")
}

// ResolveTopupAmount 根据入参解析实际支付金额与到账 quota。
// 返回值：payAmount(元), bonusQuota(quota), error
// 自定义金额恒 1:1（不提供兑换比例配置），到账额度 = 金额 × QuotaPerUnit。
//
// 版本: v0.0.24
// 日期: 2026-10-03
func ResolveTopupAmount(in CreateTopupOrderInput, presets []TopupPreset, allowCustom bool) (float64, int64, error) {
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
	return in.Amount, YuanToQuota(in.Amount), nil
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

	enabled, allowCustom, presets := GetTopupSettings()
	if !enabled {
		return nil, 0, 0, errors.New("充值功能未开启")
	}

	payAmount, bonusQuota, err := ResolveTopupAmount(in, presets, allowCustom)
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
		CreditAmount: float64(bonusQuota) / config.QuotaPerUnit,
		BonusQuota:   bonusQuota,
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

	// 解析快照获取到账 quota；若快照缺失则按 1:1（金额 × QuotaPerUnit）兜底
	bonusQuota := int64(0)
	if info := parseTopupOrderPlanInfo(order.PlanInfo); info != nil {
		bonusQuota = info.BonusQuota
	} else {
		bonusQuota = YuanToQuota(order.Amount)
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
