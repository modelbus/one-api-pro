package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/payment"
	"github.com/modelbus/one-api-pro/model"
)

// CreatePlanOrderRequest is the body of POST /api/order/plan.
type CreatePlanOrderRequest struct {
	PlanId    int    `json:"plan_id"`
	PayMethod string `json:"pay_method"`
}

// CreatePlanOrder handles POST /api/order/plan (user self-service).
// It validates the plan, runs the upgrade/stack logic, and (for
// wechat / alipay) returns a pre-payment URL / QR. The order row
// stays at status=0 (pending) until the payment channel's async
// notification arrives.
func CreatePlanOrder(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未登录"})
		return
	}
	var req CreatePlanOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	if req.PlanId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "plan_id 不能为空"})
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
	// Users can only choose wechat / alipay / bank for self-service.
	switch req.PayMethod {
	case model.OrderPayMethodWechat, model.OrderPayMethodAlipay, model.OrderPayMethodBank:
		// ok
	default:
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "自助下单仅支持 wechat / alipay / bank",
		})
		return
	}

	out, err := model.CreatePlanOrder(model.CreatePlanOrderInput{
		UserId:    userId,
		PlanId:    req.PlanId,
		PayMethod: req.PayMethod,
		Source:    model.OrderSourceUserSelf,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	// For online channels, also call the payment SDK to obtain a
	// pre-payment URL. If the channel is disabled or the SDK call
	// fails, the order is still persisted and the caller can show a
	// "configure payment settings" hint.
	var payInfo gin.H
	if req.PayMethod != model.OrderPayMethodBank {
		ch, chErr := payment.New(req.PayMethod)
		if chErr != nil {
			payInfo = gin.H{"pay_url": "", "qr_code": "", "warning": chErr.Error()}
		} else {
			enabled, _ := ch.IsEnabled()
			if !enabled {
				payInfo = gin.H{"pay_url": "", "qr_code": "", "warning": "该支付方式尚未启用"}
			} else {
				r, prepErr := ch.PrePay(out.Order.OrderNo, out.Amount, "TBUS-"+out.PackageName)
				if prepErr != nil {
					payInfo = gin.H{"pay_url": "", "qr_code": "", "warning": prepErr.Error()}
				} else {
					payInfo = gin.H{
						"pay_url":   r.PayURL,
						"qr_code":   r.QRCode,
						"expire_at": r.ExpireAt,
						"trade_no":  r.TradeNo,
					}
				}
			}
		}
	} else {
		payInfo = gin.H{"pay_url": "", "qr_code": "", "note": "请按订单详情中的账户信息完成转账，等待管理员确认"}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "",
		"order":    out.Order,
		"amount":   out.Amount,
		"plan_name": out.PackageName,
		"mode":     out.Mode,
		"pay":      payInfo,
	})
}

// GetMyOrders handles GET /api/order/self?type=1|2 (user).
func GetMyOrders(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未登录"})
		return
	}
	orderType, _ := strconv.Atoi(c.Query("type"))
	orders, err := model.GetUserOrders(userId, orderType, 200)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": orders})
}

// GetMyOrder handles GET /api/order/self/:id (user, ownership enforced).
func GetMyOrder(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未登录"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	o, err := model.GetOrderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单不存在"})
		return
	}
	if o.UserId != userId {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无权访问此订单"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": o})
}

// CancelMyOrder handles POST /api/order/self/:id/cancel (user).
func CancelMyOrder(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未登录"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	o, err := model.GetOrderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单不存在"})
		return
	}
	if o.UserId != userId {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无权访问此订单"})
		return
	}
	if err := model.CancelOrder(o); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "订单已取消"})
}

// ---------- Admin handlers ----------

// GetAllOrders handles GET /api/order (admin, paginated).
func GetAllOrders(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	orderType, _ := strconv.Atoi(c.Query("type"))
	orders, err := model.GetAllOrders(p*config.ItemsPerPage, config.ItemsPerPage, orderType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": orders})
}

// SearchOrders handles GET /api/order/search?keyword=... (admin).
func SearchOrders(c *gin.Context) {
	keyword := c.Query("keyword")
	orderType, _ := strconv.Atoi(c.Query("type"))
	orders, err := model.SearchOrders(keyword, orderType)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": orders})
}

// GetOrder handles GET /api/order/:id (admin, any user).
func GetOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	o, err := model.GetOrderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": o})
}

// MarkOrderPaidRequest is the body of PUT /api/order/:id (admin).
type MarkOrderPaidRequest struct {
	Status     int    `json:"status"`      // 1=paid, 3=refunded
	PayMethod  string `json:"pay_method"`  // optional override
	PayTradeNo string `json:"pay_trade_no"` // optional admin-entered reference
}

// MarkOrderPaid handles PUT /api/order/:id (admin).
// Allows the admin to mark a "wechat / alipay / bank / offline" order
// as paid (status=1) or refunded (status=3). For status=1 with
// PayMethod="offline" or "bank", the subscription is activated
// immediately.
func MarkOrderPaid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	o, err := model.GetOrderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单不存在"})
		return
	}
	var req MarkOrderPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	switch req.Status {
	case 1:
		mode := model.OrderUpgradeModeStack
		if req.PayMethod != "" {
			o.PayMethod = req.PayMethod
		}
		if req.PayTradeNo != "" {
			o.PayTradeNo = req.PayTradeNo
		}
		if err := model.ActivatePackageByOrder(o, mode); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "激活套餐失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "订单已支付，套餐已激活"})
	case 3:
		if err := model.MarkOrderRefunded(o); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "订单已标记为退款"})
	default:
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "不支持的状态值（仅支持 1 或 3）"})
	}
}

// DeleteOrder handles DELETE /api/order/:id (root only).
func DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	o, err := model.GetOrderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单不存在"})
		return
	}
	if err := o.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}
