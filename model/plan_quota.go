package model

import (
	"math"

	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/common/logger"
)

// QuotaPoolCapacity defines the total capacity of a weighted quota pool.
// Each model's limit contributes weight = QuotaPoolCapacity / limit.
// A plan is exhausted when Σ(consumed × QuotaPoolCapacity / limit) ≥ QuotaPoolCapacity.
// Using 100 makes the values directly interpretable as percentages.
const QuotaPoolCapacity float64 = 100

// QuotaCheckResult holds the result of checking whether a user's plan is usable.
type QuotaCheckResult struct {
	Usable          bool    `json:"usable"`
	PlanId          int     `json:"plan_id"`
	BillingType     string  `json:"billing_type"`
	PeriodWeighted  float64 `json:"period_weighted"`
	WeekWeighted    float64 `json:"week_weighted"`
	MonthWeighted   float64 `json:"month_weighted"`
	ExhaustedWindow string  `json:"exhausted_window,omitempty"`
	// 虚拟余额（微元）：VirtualAmount<=0 表示不限额度。
	VirtualAmount   int64 `json:"virtual_amount"`
	UsedAmount      int64 `json:"used_amount"`
	RemainingAmount int64 `json:"remaining_amount"`
}

// CheckPlanQuota iterates through a user's active plans (sorted by EndTime ASC)
// and returns the first plan that is not exhausted based on weighted pool calculation.
//
// 判定顺序（也是修复后的完整口径）：
//  1. 套餐未过期；
//  2. model_limits 能合法解析（解析失败 / 为空 = 该套餐不覆盖任何模型，跳过）；
//  3. 虚拟余额（virtual_amount）未耗尽 —— 只要耗尽，整个套餐即不可用；
//  4. 请求的模型在该套餐的 model_limits 内；
//  5. 该模型对应的加权池未打满（任一窗口达到 100% 即视为套餐耗尽）。
//
// 任一条不通过就尝试下一个套餐；全部不通过则回落用户全局余额按量计费。
func CheckPlanQuota(userId int, requestModel string) (*QuotaCheckResult, error) {
	ups, err := CacheGetUserActivePlans(userId)
	if err != nil {
		return nil, err
	}
	if len(ups) == 0 {
		return &QuotaCheckResult{Usable: false}, nil
	}

	now := helper.GetTimestamp()

	for _, up := range ups {
		if up.EndTime <= now {
			go func(userId int, planId uint) {
				dbPlan, _ := GetUserPlanById(int(planId))
				if dbPlan != nil && dbPlan.Status == UserPlanStatusActive {
					dbPlan.Status = UserPlanStatusExpired
					_ = dbPlan.Update()
					CacheDeleteUserActivePlans(userId)
				}
			}(userId, up.Id)
			continue
		}

		billingType := up.EffectiveBillingType()

		// model_limits 解析失败说明套餐配置是脏数据。历史实现会把它当成
		// 「不限量」直接放行（配置写错 = 全免费），此处必须判为不可用。
		limits, limitsErr := up.EffectiveModelLimits()
		if limitsErr != nil {
			logger.Debugf(nil, "CheckPlanQuota: UserPlan %d has invalid model_limits: %s, skipping", up.Id, limitsErr.Error())
			continue
		}
		// 未配置任何模型限制 = 不覆盖任何模型，跳过并回落全局余额。
		if limits == nil {
			logger.Debugf(nil, "CheckPlanQuota: UserPlan %d has no model limits, skipping", up.Id)
			continue
		}

		// 虚拟余额是套餐级别的共享池：一旦耗尽，套餐内所有模型都不可用。
		virtualAmount := up.EffectiveVirtualAmount()
		var usedAmount int64
		if virtualAmount > 0 {
			used, err := GetUserPlanUsedAmount(int(up.Id))
			if err != nil {
				logger.Debugf(nil, "CheckPlanQuota: GetUserPlanUsedAmount(%d) error: %s", up.Id, err.Error())
				continue
			}
			usedAmount = used
			if RemainingVirtualAmount(virtualAmount, used) <= 0 {
				continue
			}
		}

		// 模型不在套餐限额内 → 该套餐不覆盖此模型，跳过并尝试下一个套餐。
		if _, found := FindLimit(limits, requestModel); !found {
			continue
		}

		usages, err := GetPlanUsageByUserPlanId(int(up.Id))
		if err != nil {
			logger.Debugf(nil, "CheckPlanQuota: GetPlanUsageByUserPlanId(%d) error: %s", up.Id, err.Error())
			continue
		}

		result := &QuotaCheckResult{
			Usable:          true,
			PlanId:          int(up.Id),
			BillingType:     billingType,
			VirtualAmount:   virtualAmount,
			UsedAmount:      usedAmount,
			RemainingAmount: RemainingVirtualAmount(virtualAmount, usedAmount),
		}

		windowTypes := []string{WindowTypePeriod, WindowTypeWeek, WindowTypeMonth}
		weightedValues := make(map[string]float64, 3)

		for _, windowType := range windowTypes {
			weighted := WeightedUsage(limits, usages, windowType, billingType, now, up.StartTime)
			weightedValues[windowType] = weighted
			if IsExhausted(weighted) {
				result.Usable = false
				result.ExhaustedWindow = windowType
				break
			}
		}

		result.PeriodWeighted = weightedValues[WindowTypePeriod]
		result.WeekWeighted = weightedValues[WindowTypeWeek]
		result.MonthWeighted = weightedValues[WindowTypeMonth]

		if result.Usable {
			return result, nil
		}
	}

	return &QuotaCheckResult{Usable: false}, nil
}

// WeightedUsage computes the weighted usage for a given window type.
// This is a pure function: sum of (consumed × QuotaPoolCapacity / limit) across all models
// that have a limit > 0 and match the billingType dimension.
// Usages are filtered to the current window index for the given plan.
func WeightedUsage(
	limits map[string]ModelLimitRule,
	usages []*PlanUsage,
	windowType string,
	billingType string,
	now int64,
	startTime int64,
) float64 {
	total := 0.0
	for _, usage := range usages {
		if usage.WindowType != windowType {
			continue
		}
		rule, found := FindLimit(limits, usage.Model)
		if !found {
			continue
		}

		// Only count usage in the current window
		windowIndex := CalcWindowIndex(now, startTime, windowType, rule.PeriodH)
		if usage.WindowIndex != windowIndex {
			continue
		}

		var limit int64
		var consumed int64
		switch billingType {
		case BillingTypeRequest:
			limit = getRequestLimit(rule, windowType)
			consumed = usage.Requests
		case BillingTypeToken:
			limit = getTokenLimit(rule, windowType)
			consumed = usage.PromptTokens + usage.CompletionTokens
		default:
			limit = getRequestLimit(rule, windowType)
			consumed = usage.Requests
			if limit <= 0 {
				limit = getTokenLimit(rule, windowType)
				consumed = usage.PromptTokens + usage.CompletionTokens
			}
		}

		if limit <= 0 {
			continue
		}
		total += float64(consumed) * QuotaPoolCapacity / float64(limit)
	}
	return math.Round(total*1e8) / 1e8
}

// FindLimit returns the ModelLimitRule for a given model.
// Only an exact model-name match counts — the "other" key is treated as a regular model name.
func FindLimit(limits map[string]ModelLimitRule, model string) (ModelLimitRule, bool) {
	rule, ok := limits[model]
	return rule, ok
}

// IsExhausted returns true if the weighted usage has reached or exceeded QuotaPoolCapacity.
func IsExhausted(weighted float64) bool {
	return weighted >= QuotaPoolCapacity
}

// getRequestLimit returns the request limit for the given window type.
func getRequestLimit(rule ModelLimitRule, windowType string) int64 {
	switch windowType {
	case WindowTypePeriod:
		return rule.RequestPeriod
	case WindowTypeWeek:
		return rule.RequestWeek
	case WindowTypeMonth:
		return rule.RequestMonth
	}
	return 0
}

// getTokenLimit returns the token limit for the given window type.
func getTokenLimit(rule ModelLimitRule, windowType string) int64 {
	switch windowType {
	case WindowTypePeriod:
		return rule.TokenPeriod
	case WindowTypeWeek:
		return rule.TokenWeek
	case WindowTypeMonth:
		return rule.TokenMonth
	}
	return 0
}

// CalculateWeightedUsage returns weighted usage percentages for all window types
// of a single plan, suitable for displaying to the user.
func CalculateWeightedUsage(
	plan *Plan,
	usages []*PlanUsage,
	billingType string,
	now int64,
	startTime int64,
) map[string]float64 {
	limits := plan.GetModelLimits()
	if limits == nil {
		return nil
	}
	result := make(map[string]float64, 3)
	for _, windowType := range []string{WindowTypePeriod, WindowTypeWeek, WindowTypeMonth} {
		result[windowType] = WeightedUsage(limits, usages, windowType, billingType, now, startTime)
	}
	return result
}

// ModelUsageDetail holds per-model usage and percentage for a single window.
type ModelUsageDetail struct {
	Model            string  `json:"model"`
	Requests         int64   `json:"requests"`
	PromptTokens     int64   `json:"prompt_tokens"`
	CompletionTokens int64   `json:"completion_tokens"`
	CachedTokens     int64   `json:"cached_tokens"`
	RequestPercent   float64 `json:"request_percent"`
	TokenPercent     float64 `json:"token_percent"`
}

// CalcModelUsageDetails computes per-model usage details for the current window of each type.
// No DB queries inside — operates on already-fetched usages and limits.
// It iterates over actual usage records (not limit keys) so that model names match.
func CalcModelUsageDetails(
	limits map[string]ModelLimitRule,
	usages []*PlanUsage,
	billingType string,
	now int64,
	startTime int64,
) map[string][]ModelUsageDetail {
	type usageKey struct {
		Model      string
		WindowType string
	}

	// Build current-window usage lookup keyed by model|windowType
	usageMap := make(map[usageKey]*PlanUsage, len(usages))
	for _, u := range usages {
		rule, found := FindLimit(limits, u.Model)
		if !found {
			continue
		}
		windowIndex := CalcWindowIndex(now, startTime, u.WindowType, rule.PeriodH)
		if u.WindowIndex == windowIndex {
			usageMap[usageKey{u.Model, u.WindowType}] = u
		}
	}

	// For each window type, build detail rows
	result := make(map[string][]ModelUsageDetail, 3)
	for _, windowType := range []string{WindowTypePeriod, WindowTypeWeek, WindowTypeMonth} {
		var details []ModelUsageDetail

		// Add rows for models that have an explicit limit rule
		for modelName, rule := range limits {
			key := usageKey{modelName, windowType}

			req := int64(0)
			pt := int64(0)
			ct := int64(0)
			cacheT := int64(0)
			if pu, ok := usageMap[key]; ok {
				req = pu.Requests
				pt = pu.PromptTokens
				ct = pu.CompletionTokens
				cacheT = pu.CachedTokens
			}

			reqLimit := getRequestLimit(rule, windowType)
			tokLimit := getTokenLimit(rule, windowType)
			reqPercent := 0.0
			if reqLimit > 0 {
				reqPercent = math.Round(float64(req)*QuotaPoolCapacity/float64(reqLimit)*1e6) / 1e6
			}
			tokPercent := 0.0
			if tokLimit > 0 {
				tokPercent = math.Round(float64(pt+ct)*QuotaPoolCapacity/float64(tokLimit)*1e6) / 1e6
			}

			details = append(details, ModelUsageDetail{
				Model:            modelName,
				Requests:         req,
				PromptTokens:     pt,
				CompletionTokens: ct,
				CachedTokens:     cacheT,
				RequestPercent:   reqPercent,
				TokenPercent:     tokPercent,
			})
		}

		result[windowType] = details
	}
	return result
}

// CalcNextResetTime computes the next window reset time (Unix seconds) for a given window type.
// nextReset = startTime + (currentIndex + 1) * windowDuration
func CalcNextResetTime(now int64, startTime int64, windowType string, periodH int) int64 {
	windowIndex := CalcWindowIndex(now, startTime, windowType, periodH)
	duration := GetWindowDurationSeconds(windowType, periodH)
	return startTime + int64(windowIndex+1)*duration
}
