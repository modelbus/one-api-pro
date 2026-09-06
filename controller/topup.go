// Package controller 提供 HTTP handler 实现。
// 本文件实现"在线支付充值（topup）"业务：
//   - CreateTopupOrder: 用户自助下单（POST /api/topup/order）
//   - GetTopupSettings: 管理员读取充值设置（GET /api/setting/topup）
//   - PutTopupSettings: 管理员保存充值设置（PUT /api/setting/topup）
//
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common/payment"
	"github.com/modelbus/one-api-pro/model"
)

// CreateTopupOrderRequest 是 POST /api/topup/order 的请求体。
// Amount 与 PresetAmount 至少传一个：PresetAmount > 0 命中预设；
// 否则按 Amount × exchange_rate 算 bonus_quota（需开启 allow_custom）。
//
// 版本: v0.0.10
// 日期: 2026-09-06
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
	_ = payAmount
	_ = bonusQuota

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "",
		"order":       order,
		"amount":      payAmount,
		"bonus_quota": bonusQuota,
		"pay":         payInfo,
	})
}

// TopupSettingsResponse 充值设置的 API 响应结构。
//
// 版本: v0.0.10
// 日期: 2026-09-06
type TopupSettingsResponse struct {
	Enabled      bool              `json:"enabled"`
	AllowCustom  bool              `json:"allow_custom"`
	ExchangeRate int64             `json:"exchange_rate"`
	Presets      []model.TopupPreset `json:"presets"`
}

// GetTopupSettings handles GET /api/setting/topup (root only).
//
// 版本: v0.0.10
// 日期: 2026-09-06
func GetTopupSettings(c *gin.Context) {
	enabled, allowCustom, presets, exchangeRate := model.GetTopupSettings()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": TopupSettingsResponse{
			Enabled:      enabled,
			AllowCustom:  allowCustom,
			ExchangeRate: exchangeRate,
			Presets:      presets,
		},
	})
}

// PutTopupSettingsRequest 是 PUT /api/setting/topup 的请求体。
//
// 版本: v0.0.10
// 日期: 2026-09-06
type PutTopupSettingsRequest struct {
	Enabled      bool                `json:"enabled"`
	AllowCustom  bool                `json:"allow_custom"`
	ExchangeRate int64               `json:"exchange_rate"`
	Presets      []model.TopupPreset `json:"presets"`
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
	if err := model.SaveTopupSettings(req.Enabled, req.AllowCustom, req.Presets, req.ExchangeRate); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "已保存"})
}
