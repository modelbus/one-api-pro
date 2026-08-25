package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/ctxkey"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/model"
)

func GetAllSubscriptions(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	userId, _ := strconv.Atoi(c.Query("user_id"))
	status := -1
	if s := c.Query("status"); s != "" {
		status, _ = strconv.Atoi(s)
	}
	ups, err := model.GetAllUserPlans(p*config.ItemsPerPage, config.ItemsPerPage, userId, status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    ups,
	})
}

func SearchSubscriptions(c *gin.Context) {
	keyword := c.Query("keyword")
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	ups, err := model.SearchUserPlans(keyword, p*config.ItemsPerPage, config.ItemsPerPage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    ups,
	})
}

func GetSubscriptionDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	up, err := model.GetUserPlanById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	plan, _ := model.GetPlanById(up.PlanId)
	up.Plan = plan
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    up,
	})
}

type AddSubscriptionRequest struct {
	UserId       int    `json:"user_id"`
	PlanId       int    `json:"plan_id"`
	BillingType  string `json:"billing_type"`
	DurationDays int    `json:"duration_days"`
	Notes        string `json:"notes"`
	PayMethod    string `json:"pay_method"` // wechat/alipay/bank/offline/free
}

func AddSubscription(c *gin.Context) {
	var req AddSubscriptionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if req.UserId == 0 || req.PlanId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户 ID 和套餐 ID 不能为空",
		})
		return
	}
	if req.PayMethod == "" {
		req.PayMethod = model.OrderPayMethodFree
	}
	if !model.IsValidPayMethod(req.PayMethod) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "不支持的支付方式: " + req.PayMethod,
		})
		return
	}
	plan, err := model.GetPlanById(req.PlanId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "套餐不存在",
		})
		return
	}
	if plan.Status != model.PlanStatusEnabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "套餐已下架",
		})
		return
	}
	user, err := model.GetUserById(req.UserId, false)
	if err != nil || user.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// Build an order row so the admin grant is auditable in /api/order.
	notes := req.Notes
	if notes == "" {
		notes = "管理员开通"
	}
	out, err := model.CreatePlanOrder(model.CreatePlanOrderInput{
		UserId:    req.UserId,
		PlanId:    req.PlanId,
		PayMethod: req.PayMethod,
		Source:    model.OrderSourceAdmin,
		Notes:     notes,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if req.BillingType != model.BillingTypeRequest && req.BillingType != model.BillingTypeToken {
		req.BillingType = model.BillingTypeToken
	}

	// For "free" the order is paid immediately and the subscription is
	// activated. For other admin pay methods (wechat/alipay/bank/offline)
	// the subscription is also activated right now — the admin has
	// confirmed the payment method offline and the order row stays as
	// a paid audit trail. This is intentional per spec: admin grants
	// always take effect immediately.
	if err := model.ActivatePackageByOrder(out.Order, model.OrderUpgradeModeStack); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "激活套餐失败: " + err.Error(),
		})
		return
	}
	up, err := model.GetUserPlanByOrderId(out.Order.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "读取新订阅失败: " + err.Error(),
		})
		return
	}
	if up != nil {
		up.Plan = plan
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    up,
		"order":   out.Order,
	})
}

type UpdateSubscriptionRequest struct {
	Id          uint   `json:"id"`
	EndTime     *int64 `json:"end_time"`
	Status      *int   `json:"status"`
	Notes       string `json:"notes"`
	BillingType string `json:"billing_type"`
}

func UpdateSubscription(c *gin.Context) {
	var req UpdateSubscriptionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if req.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "订阅 ID 不能为空",
		})
		return
	}
	up, err := model.GetUserPlanById(int(req.Id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	if req.EndTime != nil {
		up.EndTime = *req.EndTime
	}
	if req.Status != nil {
		up.Status = *req.Status
	}
	if req.BillingType != "" {
		up.BillingType = req.BillingType
	}
	if req.Notes != "" {
		up.Notes = req.Notes
	}
	up.UpdatedTime = helper.GetTimestamp()
	err = up.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.CacheDeleteUserActivePlans(up.UserId)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func DeleteSubscription(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	up, err := model.GetUserPlanById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	err = model.DeleteUserPlanById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.CacheDeleteUserActivePlans(up.UserId)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func GetUserSubscriptionInfo(c *gin.Context) {
	userId := c.GetInt(ctxkey.Id)
	info, err := model.GetUserSubscriptionInfo(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    info,
	})
}

func GetUserSubscriptions(c *gin.Context) {
	userId := c.GetInt(ctxkey.Id)
	ups, err := model.GetAllUserPlans(0, 1000, userId, -1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    ups,
	})
}

func GetSubscriptionUsage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	up, err := model.GetUserPlanById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	pus, err := model.GetPlanUsageByUserPlanId(int(up.Id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	plan, _ := model.GetPlanById(up.PlanId)
	up.Plan = plan

	now := helper.GetTimestamp()
	limits := plan.GetModelLimits()

	// Compute weighted usage for each window type
	var weighted map[string]float64
	var modelUsage map[string][]model.ModelUsageDetail
	var nextReset map[string]int64
	if limits != nil {
		weighted = model.CalculateWeightedUsage(plan, pus, up.BillingType, now, up.StartTime, "*")
		modelUsage = model.CalcModelUsageDetails(limits, pus, up.BillingType, now, up.StartTime, plan.DefaultModel)
		// Get period_h from default_model or first rule
		var periodH int
		if plan.DefaultModel != "" {
			if rule, ok := limits[plan.DefaultModel]; ok {
				periodH = rule.PeriodH
			}
		}
		if periodH == 0 {
			for _, r := range limits {
				periodH = r.PeriodH
				break
			}
		}
		nextReset = map[string]int64{
			model.WindowTypePeriod: model.CalcNextResetTime(now, up.StartTime, model.WindowTypePeriod, periodH),
			model.WindowTypeWeek:   model.CalcNextResetTime(now, up.StartTime, model.WindowTypeWeek, periodH),
			model.WindowTypeMonth:  model.CalcNextResetTime(now, up.StartTime, model.WindowTypeMonth, periodH),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"subscription": up,
			"usage":         pus,
			"weighted":      weighted,
			"limits":        limits,
			"model_usage":   modelUsage,
			"next_reset":    nextReset,
			"now":           now,
			"start_time":    up.StartTime,
			"billing_type":  up.BillingType,
		},
	})
}