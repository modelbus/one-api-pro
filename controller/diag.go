// diag.go admin 诊断端点
// Diag handlers: 用于排查 admin 端报错的 admin-only diagnostics.
// 版本: v0.0.13
// 日期: 2026-09-08
// 作者: opencode

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/model"
)

// DiagSubscriptionRow 是一行诊断数据：user_plan 行 + 对应 order 行 + user 摘要。
// 用于排查「user_plans.plan_id 显示为 0」类问题。
// DiagSubscriptionRow is one row in the diag table: a user_plan + its
// source order + brief user info. Used to debug issues like
// "user_plans.plan_id always shows 0".
// 版本: v0.0.13
type DiagSubscriptionRow struct {
	UserPlan  *model.UserPlan    `json:"user_plan"`
	OrderNo   string             `json:"order_no"`
	OrderPlan int                `json:"order_plan_id"`
	OrderType int                `json:"order_type"`
	User      *model.UserBrief   `json:"user"`
}

// DiagSubscriptionsResponse 是 /api/diag/subscriptions 的返回体。
type DiagSubscriptionsResponse struct {
	TotalUserPlans int                 `json:"total_user_plans"`
	MissingPlanId  int                 `json:"missing_plan_id"`
	Rows           []DiagSubscriptionRow `json:"rows"`
}

// DiagSubscriptions 处理 GET /api/diag/subscriptions (admin)。
// 返回 user_plans 与 orders 的逐行对比，帮助诊断 plan_id 显示为 0 的问题。
// - total_user_plans: 扫到的 user_plan 总数
// - missing_plan_id: plan_id <= 0 的行数（疑似脏数据）
// - rows: 最多 100 行，每行带 user_plan + order_no + order_plan_id + user
//
// 版本: v0.0.13
// 日期: 2026-09-08
func DiagSubscriptions(c *gin.Context) {
	ups, err := model.GetAllUserPlans(0, 100, 0, -1)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}

	missing := 0
	userIds := make([]int, 0, len(ups))
	for _, u := range ups {
		if u.PlanId <= 0 {
			missing++
		}
		userIds = append(userIds, u.UserId)
	}
	userMap, _ := model.GetUsersBriefByIds(userIds)

	rows := make([]DiagSubscriptionRow, 0, len(ups))
	for _, u := range ups {
		row := DiagSubscriptionRow{
			UserPlan: u,
			User:     userMap[u.UserId],
		}
		if u.OrderId > 0 {
			if o, err := model.GetOrderById(u.OrderId); err == nil && o != nil {
				row.OrderNo = o.OrderNo
				row.OrderPlan = o.PlanId
				row.OrderType = o.Type
			}
		}
		rows = append(rows, row)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": DiagSubscriptionsResponse{
			TotalUserPlans: len(ups),
			MissingPlanId:  missing,
			Rows:           rows,
		},
	})
}