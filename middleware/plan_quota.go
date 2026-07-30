package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/common/ctxkey"
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/Leon-PanPan/one-api-pro/model"
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
			if result.DefaultModel != "" {
				c.Set(ctxkey.DefaultModel, result.DefaultModel)
				c.Set(ctxkey.RequestModel, result.DefaultModel)
				logger.Debugf(c.Request.Context(), "user %d model %s not in plan limits, forwarding to default_model %s for plan %d",
					userId, requestModel, result.DefaultModel, result.PlanId)
			}
			logger.Debugf(c.Request.Context(), "user %d using plan %d, billing_type %s for model %s, weighted: period=%.4f week=%.4f month=%.4f",
				userId, result.PlanId, result.BillingType, requestModel,
				result.PeriodWeighted, result.WeekWeighted, result.MonthWeighted)
		} else {
			// Plan is not usable or user has no plan - check if user has quota-based access
			// If the user has active plans with model_limits and the model is not in limits and no default_model,
			// we need to block the request with 422
			ups, upsErr := model.CacheGetUserActivePlans(userId)
			if upsErr == nil && len(ups) > 0 {
				for _, up := range ups {
					if up.Plan == nil {
						continue
					}
					limits := up.Plan.GetModelLimits()
					if limits == nil {
						continue
					}
					_, _, found := model.FindLimit(limits, requestModel, "")
					if !found && up.Plan.DefaultModel == "" {
						c.JSON(http.StatusUnprocessableEntity, gin.H{
							"success": false,
							"message": fmt.Sprintf("模型 %s 不在套餐限额配置中，且套餐未设置默认模型，请联系管理员", requestModel),
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}