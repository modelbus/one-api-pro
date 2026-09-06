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

// payStatusSuccess / payStatusWarning are the user-facing pay.status
// values returned alongside CreatePlanOrder / PayMyOrder. They let the
// frontend decide whether to render the QR/redirect modal or surface a
// warning without guessing from whether pay_url is empty.
const (
	payStatusSuccess = "success"
	payStatusWarning = "warning"
)

// buildPayInfo initializes a payment for the given order and returns
// the "pay" object that CreatePlanOrder / PayMyOrder embed in their
// JSON response.
//
// status:
//   - "success" — pay_url / qr_code is ready (or bank-transfer note is
//     present); the frontend should render the payment modal.
//   - "warning" — payment channel is not registered, disabled, or the
//     SDK call failed; the order is still persisted and the warning
//     text is included so the frontend can show a toast.
func buildPayInfo(payMethod, orderNo string, amount float64, packageName string) gin.H {
	if payMethod == model.OrderPayMethodBank {
		return gin.H{
			"status":  payStatusSuccess,
			"pay_url": "",
			"qr_code": "",
			"note":    "请按订单详情中的账户信息完成转账，等待管理员确认",
		}
	}
	ch, chErr := payment.New(payMethod)
	if chErr != nil {
		return gin.H{
			"status":  payStatusWarning,
			"pay_url": "",
			"qr_code": "",
			"warning": chErr.Error(),
		}
	}
	enabled, _ := ch.IsEnabled()
	if !enabled {
		return gin.H{
			"status":  payStatusWarning,
			"pay_url": "",
			"qr_code": "",
			"warning": "该支付方式尚未启用",
		}
	}
	r, prepErr := ch.PrePay(orderNo, amount, "TBUS-"+packageName)
	if prepErr != nil {
		return gin.H{
			"status":  payStatusWarning,
			"pay_url": "",
			"qr_code": "",
			"warning": prepErr.Error(),
		}
	}
	return gin.H{
		"status":    payStatusSuccess,
		"pay_url":   r.PayURL,
		"qr_code":   r.QRCode,
		"expire_at": r.ExpireAt,
		"trade_no":  r.TradeNo,
	}
}

// CreatePlanOrder handles POST /api/order/plan (user self-service).
// It validates the plan, runs the upgrade/stack logic, and (for
// wechat / alipay) returns a pre-payment URL / QR. The order row
// stays at status=0 (pending) until the payment channel's async
// notification arrives.
//
// If the admin has not enabled any payment channel, the request is
// rejected with no order created — the user must be told to ask the
// admin to enable a channel first.
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

	// Refuse to create the order when no payment channel is configured.
	// The order row is intentionally NOT persisted in this case so the
	// /api/order/self list stays clean.
	anyEnabled, _, _ := payment.AnyChannelEnabled()
	if !anyEnabled {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": noPaymentEnabledMsg})
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
	payInfo := buildPayInfo(req.PayMethod, out.Order.OrderNo, out.Amount, out.PackageName)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "",
		"order":     out.Order,
		"amount":    out.Amount,
		"plan_name": out.PackageName,
		"mode":      out.Mode,
		"pay":       payInfo,
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

// PayMyOrderRequest is the body of POST /api/order/self/:id/pay.
type PayMyOrderRequest struct {
	PayMethod string `json:"pay_method"`
}

// PayMyOrder handles POST /api/order/self/:id/pay (user, ownership
// enforced). It re-issues the pre-payment URL/QR for an existing
// pending order, optionally switching the pay_method first.
//
// Returns the same response shape as CreatePlanOrder so the frontend
// can reuse its payment modal.
func PayMyOrder(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未登录"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var req PayMyOrderRequest
	_ = c.ShouldBindJSON(&req) // body is optional

	o, err := model.GetOrderById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单不存在"})
		return
	}
	if o.UserId != userId {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无权访问此订单"})
		return
	}
	if o.Status != model.OrderStatusPending {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "订单当前不可支付"})
		return
	}

	// Pick the pay_method to use. The request may override the one
	// stored at order-creation time.
	payMethod := req.PayMethod
	if payMethod == "" {
		payMethod = o.PayMethod
	}
	if !model.IsValidPayMethod(payMethod) {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "不支持的支付方式"})
		return
	}
	switch payMethod {
	case model.OrderPayMethodWechat, model.OrderPayMethodAlipay, model.OrderPayMethodBank:
		// ok
	default:
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "自助下单仅支持 wechat / alipay / bank",
		})
		return
	}

	// Refuse if no channel is enabled at all.
	anyEnabled, _, _ := payment.AnyChannelEnabled()
	if !anyEnabled {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": noPaymentEnabledMsg})
		return
	}

	// Persist the chosen pay_method if it changed.
	if payMethod != o.PayMethod {
		o.PayMethod = payMethod
		_ = o.Update()
	}

	packageName := ""
	if p := model.GetPlanByOrderPlanInfo(o.PlanInfo); p != nil {
		packageName = p.Name
	}

	payInfo := buildPayInfo(payMethod, o.OrderNo, o.Amount, packageName)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "",
		"order":     o,
		"amount":    o.Amount,
		"plan_name": packageName,
		"pay":       payInfo,
	})
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
