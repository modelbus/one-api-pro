package controller

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common/payment"
	"github.com/modelbus/one-api-pro/model"
)

// processNotify reads the raw body, asks the channel to verify the
// signature, marks the order paid and (re-)activates the subscription.
// Returns the channel-specific success payload to write back (WeChat
// needs XML, Alipay needs the literal string "success").
func processNotify(c *gin.Context, payMethod string) (string, error) {
	ch, err := payment.New(payMethod)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}
	notif, err := ch.VerifyNotify(body)
	if err != nil {
		return "", err
	}
	if !notif.Paid {
		return "", errors.New("payment not marked paid by channel")
	}
	order, err := model.GetOrderByOrderNo(notif.OutTradeNo)
	if err != nil {
		return "", err
	}
	if order.Amount > 0 && notif.Amount > 0 && notif.Amount != order.Amount {
		return "", errors.New("amount mismatch")
	}
	if err := model.ActivatePackageByOrder(order, model.OrderUpgradeModeStack); err != nil {
		return "", err
	}
	switch payMethod {
	case model.OrderPayMethodWechat:
		return `<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>`, nil
	default:
		return "success", nil
	}
}

// WechatNotify handles POST /api/payment/wechat/notify (no auth; public).
func WechatNotify(c *gin.Context) {
	body, err := processNotify(c, model.OrderPayMethodWechat)
	if err != nil {
		// Per WeChat spec, return FAIL on signature/business errors so
		// the channel can retry.
		c.Data(http.StatusOK, "application/xml",
			[]byte(`<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[`+err.Error()+`]]></return_msg></xml>`))
		return
	}
	c.Data(http.StatusOK, "application/xml", []byte(body))
}

// AlipayNotify handles POST /api/payment/alipay/notify (no auth; public).
func AlipayNotify(c *gin.Context) {
	body, err := processNotify(c, model.OrderPayMethodAlipay)
	if err != nil {
		// Alipay expects literal "fail" on errors.
		c.String(http.StatusOK, "fail")
		return
	}
	c.String(http.StatusOK, body) // "success"
}

// MockPayRequest is the body of POST /api/payment/mock/notify (root).
type MockPayRequest struct {
	OrderNo string `json:"order_no"`
	Status  int    `json:"status"` // 1=paid, 3=refunded
}

// GetPaymentStatus handles GET /api/payment/status (public).
// Returns the enabled status of every registered payment channel so
// the user-facing pages can decide whether to show the purchase UI or
// short-circuit with an error message.
//
// Response shape:
//
//	{
//	  "success": true,
//	  "data": {
//	    "any_enabled": bool,
//	    "methods": [{ "name": "wechat", "label": "微信支付", "enabled": true }, ...]
//	  }
//	}
func GetPaymentStatus(c *gin.Context) {
	anyEnabled, methods, _ := payment.AnyChannelEnabled()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"any_enabled": anyEnabled,
			"methods":     methods,
		},
	})
}

// noPaymentEnabledMsg is the user-facing error returned by both
// CreatePlanOrder and PayMyOrder when the admin has not enabled any
// payment channel.
const noPaymentEnabledMsg = "系统尚未开通任何支付通道，请设置后开启支付"

// MockPay handles POST /api/payment/mock/notify (root).
// Used by tests / manual operations to mark an order paid without
// going through a real payment channel.
func MockPay(c *gin.Context) {
	var req MockPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	if req.OrderNo == "" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "order_no 不能为空"})
		return
	}
	order, err := model.GetOrderByOrderNo(req.OrderNo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单不存在"})
		return
	}
	switch req.Status {
	case 1:
		if err := model.ActivatePackageByOrder(order, model.OrderUpgradeModeStack); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "激活失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "订单已支付，套餐已激活"})
	case 3:
		if err := model.MarkOrderRefunded(order); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "订单已标记为退款"})
	default:
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "不支持的状态值（仅支持 1 或 3）"})
	}
	_ = bytes.NewBuffer(nil) // keep imports happy if unused
}
