package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common/ctxkey"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
)

func PlanQuotaCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		userId := c.GetInt(ctxkey.Id)
		requestModel := c.GetString(ctxkey.RequestModel)

		result, err := model.CheckPlanQuota(userId, requestModel)
		if err != nil {
			logger.SysError("CheckPlanQuota error: " + err.Error())
			c.Next()
			return
		}

		if result.Usable {
			c.Set(ctxkey.PlanId, result.PlanId)
			c.Set(ctxkey.BillingType, result.BillingType)
			logger.Debugf(c.Request.Context(), "user %d using plan %d, billing_type %s for model %s, weighted: period=%.4f week=%.4f month=%.4f",
				userId, result.PlanId, result.BillingType, requestModel,
				result.PeriodWeighted, result.WeekWeighted, result.MonthWeighted)
		}
		// 套餐不可用（无套餐 / 已耗尽 / 不覆盖该模型）时不拦截：PlanId 保持 0，
		// 下游按全局余额（按量计费）处理；余额不足由按量计费链路自行拒绝。

		c.Next()
	}
}