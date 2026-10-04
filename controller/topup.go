// Package controller 提供 HTTP handler 实现。
// 本文件实现"在线支付充值（topup）"业务：
//   - CreateTopupOrder: 用户自助下单（POST /api/topup/order）
//   - GetTopupSettings: 管理员读取充值设置（GET /api/setting/topup）
//   - PutTopupSettings: 管理员保存充值设置（PUT /api/setting/topup）
//
// 版本: v0.0.22
// 日期: 2026-10-03
// 作者: opencode
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common/payment"
	"github.com/modelbus/one-api-pro/model"
)

// CreateTopupOrderRequest 是 POST /api/topup/order 的请求体。
// Amount 与 PresetAmount 至少传一个：PresetAmount > 0 命中预设（按预设的到账额度发放）；
// 否则按 Amount 走自定义金额（恒 1:1，到账 = Amount × QuotaPerUnit，需开启 allow_custom）。
//
// 版本: v0.0.24
// 日期: 2026-10-03
type CreateTopupOrderRequest struct {
	Amount       float64 `json:"amount"`
	PresetAmount float64 `json:"preset_amount"`
	PayMethod    string  `json:"pay_method"`
}

// CreateTopupOrder handles POST /api/topup/order (user self-service).
// 流程：参数校验 → 解析金额与配额 → 创建订单（type=2）→ 调 buildPayInfo。
// 前置条件：任意支付通道已启用；topup.enabled=true。
//
// 版本: v0.0.10
// 日期: 2026-09-06
func CreateTopupOrder(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未登录"})
		return
	}
	var req CreateTopupOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	if req.Amount <= 0 && req.PresetAmount <= 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "amount 或 preset_amount 至少传一个"})
		return
	}
	if req.PayMethod == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "pay_method 不能为空"})
		return
	}
	if !model.IsValidPayMethod(req.PayMethod) {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "不支持的支付方式"})
		return
	}
	switch req.PayMethod {
	case model.OrderPayMethodWechat, model.OrderPayMethodAlipay, model.OrderPayMethodBank:
		// ok
	default:
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "自助充值仅支持 wechat / alipay / bank",
		})
		return
	}

	// Refuse if no payment channel enabled at all.
	if anyEnabled, _, _ := payment.AnyChannelEnabled(); !anyEnabled {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": noPaymentEnabledMsg})
		return
	}

	order, payAmount, bonusQuota, err := model.CreateTopupOrder(model.CreateTopupOrderInput{
		UserId:       userId,
		Amount:       req.Amount,
		PresetAmount: req.PresetAmount,
		PayMethod:    req.PayMethod,
		Source:       model.OrderSourceUserSelf,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	payInfo := buildPayInfo(req.PayMethod, order.OrderNo, order.Amount, "余额充值")

	// bonus_quota 在 API 边界统一换算为「元」（见 quota_dto.go）。
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "",
		"order":       order,
		"amount":      payAmount,
		"bonus_quota": quotaToYuan(bonusQuota),
		"pay":         payInfo,
	})
}

// TopupSettingsResponse 充值设置的 API 响应结构。
// presets[].bonus_quota 以「元」返回（见 quota_dto.go）。
//
// 版本: v0.0.25
// 日期: 2026-10-04
type TopupSettingsResponse struct {
	Enabled     bool             `json:"enabled"`
	AllowCustom bool             `json:"allow_custom"`
	Presets     []topupPresetDTO `json:"presets"`
}

// GetTopupSettings handles GET /api/setting/topup (root only).
//
// 版本: v0.0.10
// 日期: 2026-09-06
func GetTopupSettings(c *gin.Context) {
	enabled, allowCustom, presets := model.GetTopupSettings()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": TopupSettingsResponse{
			Enabled:     enabled,
			AllowCustom: allowCustom,
			Presets:     toTopupPresetDTOs(presets),
		},
	})
}

// topupPresetInput 充值快捷金额入参：bonus_quota 以「元」传入。
// Top-up preset input: bonus_quota is expressed in CNY yuan.
type topupPresetInput struct {
	Amount     float64 `json:"amount"`
	BonusQuota float64 `json:"bonus_quota"`
}

// topupPresetsToModel 把「元」口径的入参换算回 model.TopupPreset 的微元口径。
// Convert yuan-denominated input back to the micro-quota model.TopupPreset.
func topupPresetsToModel(in []topupPresetInput) []model.TopupPreset {
	out := make([]model.TopupPreset, 0, len(in))
	for _, p := range in {
		out = append(out, model.TopupPreset{
			Amount:     p.Amount,
			BonusQuota: yuanToQuota(p.BonusQuota),
		})
	}
	return out
}

// PutTopupSettingsRequest 是 PUT /api/setting/topup 的请求体。
// presets[].bonus_quota 以「元」传入（见 quota_dto.go）。
//
// 版本: v0.0.25
// 日期: 2026-10-04
type PutTopupSettingsRequest struct {
	Enabled     bool               `json:"enabled"`
	AllowCustom bool               `json:"allow_custom"`
	Presets     []topupPresetInput `json:"presets"`
}

// PutTopupSettings handles PUT /api/setting/topup (root only).
//
// 版本: v0.0.10
// 日期: 2026-09-06
func PutTopupSettings(c *gin.Context) {
	var req PutTopupSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	if err := model.SaveTopupSettings(req.Enabled, req.AllowCustom, topupPresetsToModel(req.Presets)); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "已保存"})
}
