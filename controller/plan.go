package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/model"
)

func GetAllPlans(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	plans, err := model.GetAllPlans(p*config.ItemsPerPage, config.ItemsPerPage)
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
		"data":    plans,
	})
}

func SearchPlans(c *gin.Context) {
	keyword := c.Query("keyword")
	plans, err := model.SearchPlans(keyword)
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
		"data":    plans,
	})
}

func GetPlan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	plan, err := model.GetPlanById(id)
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
		"data":    plan,
	})
}

func AddPlan(c *gin.Context) {
	plan := model.Plan{}
	err := c.ShouldBindJSON(&plan)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if plan.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "套餐名称不能为空",
		})
		return
	}
	if err := plan.ValidateDefaultModel(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	plan.CreatedTime = helper.GetTimestamp()
	plan.UpdatedTime = helper.GetTimestamp()
	err = plan.Insert()
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
		"data":    plan,
	})
}

func UpdatePlan(c *gin.Context) {
	plan := model.Plan{}
	err := c.ShouldBindJSON(&plan)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if plan.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "套餐 ID 不能为空",
		})
		return
	}
	if err := plan.ValidateDefaultModel(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	plan.UpdatedTime = helper.GetTimestamp()
	err = plan.Update()
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
	})
}

func DeletePlan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := model.DeletePlanById(id)
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
	})
}

// GetPublicPlans returns enabled plans for the user-facing subscription page.
// No auth required.
func GetPublicPlans(c *gin.Context) {
	all, err := model.GetAllPlans(0, 1000)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	enabled := make([]*model.Plan, 0, len(all))
	for _, p := range all {
		if p.Status == model.PlanStatusEnabled {
			enabled = append(enabled, p)
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": enabled})
}

// GetPublicPlanDetail returns a single plan by id (no auth).
func GetPublicPlanDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	plan, err := model.GetPlanById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": plan})
}

// GetCurrentPlan returns the authenticated user's active subscription.
// Auth: UserAuth. The response flattens the embedded Plan fields to
// the top level so the user-facing frontend (which mirrors tbus-web)
// can read fields like sort / duration_days / price directly.
func GetCurrentPlan(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未登录"})
		return
	}
	actives, err := model.GetActiveUserPlansByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(actives) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": nil})
		return
	}
	up := actives[0]
	now := helper.GetTimestamp()
	expireDate := time.Unix(up.EndTime, 0).Format("2006-01-02")
	data := gin.H{
		"id":              up.Id,
		"user_id":         up.UserId,
		"plan_id":         up.PlanId,
		"order_id":        up.OrderId,
		"start_time":      up.StartTime,
		"end_time":        up.EndTime,
		"expire_time":     up.EndTime, // unix seconds (compat with tbus)
		"expire_date":     expireDate, // ISO date (compat with tbus)
		"status":          up.Status,
		"billing_type":    up.BillingType,
		"notes":           up.Notes,
		"created_time":    up.CreatedTime,
		"updated_time":    up.UpdatedTime,
		"is_expired":      up.EndTime <= now,
		"remaining_days":  int64(0),
	}
	if up.EndTime > now {
		data["remaining_days"] = int64((up.EndTime - now + 86399) / 86400)
	}
	if up.Plan != nil {
		data["name"] = up.Plan.Name
		data["price"] = up.Plan.Price
		data["duration_days"] = up.Plan.DurationDays
		data["duration_text"] = up.Plan.DurationText
		data["sort"] = up.Plan.Sort
		data["recommended"] = up.Plan.Recommended
		data["description"] = up.Plan.Description
		data["features"] = up.Plan.Features
		data["model_limits"] = up.Plan.ModelLimits
		data["default_model"] = up.Plan.DefaultModel
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": data})
}