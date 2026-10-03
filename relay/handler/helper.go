package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/relay/constant/role"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/logger"
	dbmodel "github.com/modelbus/one-api-pro/model"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	billingratio "github.com/modelbus/one-api-pro/relay/billing/ratio"
	"github.com/modelbus/one-api-pro/relay/handler/validator"
	"github.com/modelbus/one-api-pro/relay/meta"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
	"github.com/modelbus/one-api-pro/relay/relaymode"
)

func getAndValidateTextRequest(c *gin.Context, relayMode int) (*relaymodel.GeneralOpenAIRequest, error) {
	textRequest := &relaymodel.GeneralOpenAIRequest{}
	err := common.UnmarshalBodyReusable(c, textRequest)
	if err != nil {
		return nil, err
	}
	if relayMode == relaymode.Moderations && textRequest.Model == "" {
		textRequest.Model = "text-moderation-latest"
	}
	if relayMode == relaymode.Embeddings && textRequest.Model == "" {
		textRequest.Model = c.Param("model")
	}
	err = validator.ValidateTextRequest(textRequest, relayMode)
	if err != nil {
		return nil, err
	}
	return textRequest, nil
}

func getPromptTokens(textRequest *relaymodel.GeneralOpenAIRequest, relayMode int) int {
	switch relayMode {
	case relaymode.ChatCompletions:
		return openai.CountTokenMessages(textRequest.Messages, textRequest.Model)
	case relaymode.Completions:
		return openai.CountTokenInput(textRequest.Prompt, textRequest.Model)
	case relaymode.Moderations:
		return openai.CountTokenInput(textRequest.Input, textRequest.Model)
	}
	return 0
}

func getPreConsumedQuota(textRequest *relaymodel.GeneralOpenAIRequest, promptTokens int, ratio float64) int64 {
	preConsumedTokens := config.PreConsumedQuota + int64(promptTokens)
	if textRequest.MaxTokens != 0 {
		preConsumedTokens += int64(textRequest.MaxTokens)
	}
	return int64(float64(preConsumedTokens) * ratio)
}

// preConsumeQuota pre-consumes an *estimated* quota derived from the token
// pricing ratio. Used for token-billed models, where the final cost is only
// known after the upstream returns its usage.
//
// preConsumeQuota 按 token 定价 ratio 估算并预扣额度，用于 token 计费模型
// （最终费用需等上游返回 usage 才能确定）。
func preConsumeQuota(ctx context.Context, textRequest *relaymodel.GeneralOpenAIRequest, promptTokens int, ratio float64, meta *meta.Meta) (int64, *relaymodel.ErrorWithStatusCode) {
	preConsumedQuota := getPreConsumedQuota(textRequest, promptTokens, ratio)
	return preConsumeAmount(ctx, preConsumedQuota, meta)
}

// preConsumeExactQuota pre-consumes an *exact* quota that is fully
// deterministic before the upstream call. Used for per_request billing, where
// the cost is a flat price per call — pre-consume and settlement use the very
// same amount, so the settlement delta is always 0.
//
// Why this matters: with the old code the per_request path reused the
// token-based estimate as the pre-consume amount (a few hundred quota) while
// settlement charged the full per-request price (tens of thousands of quota).
// The delta was applied straight to users.quota, draining the balance into the
// negative on a single call, and the "is balance enough" guard compared the
// balance against the tiny estimate, so it never rejected.
//
// preConsumeExactQuota 按请求价格精确预扣，用于 per_request 计费：费用在请求前
// 已完全确定，预扣与结算使用同一金额，结算 delta 恒为 0。
//
// 版本: v0.0.24
// 日期: 2026-10-03
func preConsumeExactQuota(ctx context.Context, preConsumedQuota int64, meta *meta.Meta) (int64, *relaymodel.ErrorWithStatusCode) {
	return preConsumeAmount(ctx, preConsumedQuota, meta)
}

// preConsumeAmount applies the shared pre-consume guards and DB writes for an
// already-computed quota amount:
//
//  1. balance is insufficient → reject with 403, do NOT touch Redis/DB
//  2. balance is large enough (100×) → trust the user and skip pre-consume
//  3. otherwise → decrease Redis + token remain_quota, roll back on failure
//
// preConsumeAmount 对已计算好的预扣额度执行统一守卫与落库：余额不足 → 403；
// 余额足够大（100×）→ 免预扣；否则正常预扣，失败时回滚。
func preConsumeAmount(ctx context.Context, preConsumedQuota int64, meta *meta.Meta) (int64, *relaymodel.ErrorWithStatusCode) {
	if meta.PlanId > 0 {
		return 0, nil
	}

	userQuota, err := dbmodel.CacheGetUserQuota(ctx, meta.UserId)
	if err != nil {
		return preConsumedQuota, openai.ErrorWrapper(err, "get_user_quota_failed", http.StatusInternalServerError)
	}

	// Guard 1: balance is insufficient — reject immediately, do NOT touch Redis.
	if userQuota-preConsumedQuota < 0 {
		return preConsumedQuota, openai.ErrorWrapper(errors.New("user quota is not enough"), "insufficient_user_quota", http.StatusForbidden)
	}

	// Guard 2: balance is large enough to fully trust — skip pre-consume
	// entirely (no Redis change, no token.quota change). The previous
	// implementation set preConsumedQuota = 0 only AFTER Redis was already
	// decreased, causing per-request drift of ~preConsumedQuota on every
	// "trusted" request. For high-price models (e.g. ¥8.40/M tokens) this
	// drift is ~80_000 per request — enough to push the cached value into
	// the dead zone [lowWaterMark, preConsumedQuota) within tens of requests
	// and produce spurious 403s even though users.quota is still healthy.
	if userQuota > 100*preConsumedQuota {
		logger.Info(ctx, fmt.Sprintf("user %d has enough quota %d, trusted and no need to pre-consume", meta.UserId, userQuota))
		return 0, nil
	}

	// Only reach here when we actually need to pre-consume.
	err = dbmodel.CacheDecreaseUserQuota(meta.UserId, preConsumedQuota)
	if err != nil {
		return preConsumedQuota, openai.ErrorWrapper(err, "decrease_user_quota_failed", http.StatusInternalServerError)
	}
	err = dbmodel.PreConsumeTokenQuota(meta.TokenId, preConsumedQuota)
	if err != nil {
		// Roll back the Redis decrement so the cache does not drift further
		// from users.quota when the DB-side pre-consume fails.
		if rbErr := dbmodel.CacheIncreaseUserQuota(meta.UserId, preConsumedQuota); rbErr != nil {
			logger.Error(ctx, fmt.Sprintf("rollback CacheIncreaseUserQuota failed for user %d: %s", meta.UserId, rbErr.Error()))
		}
		return preConsumedQuota, openai.ErrorWrapper(err, "pre_consume_token_quota_failed", http.StatusForbidden)
	}
	return preConsumedQuota, nil
}

func postConsumeQuota(ctx context.Context, usage *relaymodel.Usage, meta *meta.Meta, textRequest *relaymodel.GeneralOpenAIRequest, preConsumedQuota int64, priceResult *billingratio.PriceResult, groupDiscount float64, systemPromptReset bool) {
	if usage == nil {
		logger.Error(ctx, "usage is nil, which is unexpected")
		return
	}
	promptTokens := usage.PromptTokens
	completionTokens := usage.CompletionTokens
	cachedTokens := 0
	if usage.PromptTokensDetails != nil {
		cachedTokens = usage.PromptTokensDetails.CachedTokens
	}

	var quota int64
	if priceResult.BillingType == dbmodel.BillingTypePerRequest {
		quota = billingratio.CalculatePerRequestQuota(priceResult.PerRequestPrice, 1, 1, groupDiscount)
	} else {
		quota = billingratio.CalculateTokenQuota(
			priceResult.InputPrice, priceResult.OutputPrice, priceResult.CachedPrice,
			promptTokens, completionTokens, cachedTokens,
			groupDiscount,
		)
	}

	// A zero-token usage means the upstream reported no billable token usage
	// (e.g. an empty/aborted stream). For token billing there is nothing to
	// charge, so zero the quota out. per_request billing is a flat per-call
	// price that does not depend on tokens — zeroing it here would silently
	// refund the exact amount pre-consumed in RelayTextHelper and make the
	// call free, so it must be exempted.
	//
	// usage 为 0 token 表示上游未上报可计费 token（空响应/中断流）。token 计费
	// 下无费用可收，置 0 合理；per_request 是「每次固定价」，与 token 无关，
	// 置 0 会把精确预扣的按次费用全额退回，因此必须豁免。
	if priceResult.BillingType != dbmodel.BillingTypePerRequest {
		if promptTokens+completionTokens == 0 {
			quota = 0
		}
	}

	planId := meta.PlanId
	billingSource := 0
	if planId > 0 {
		billingSource = 1
		now := helper.GetTimestamp()
		ups, err := dbmodel.CacheGetUserActivePlans(meta.UserId)
		if err != nil || len(ups) == 0 {
			logger.Error(ctx, "failed to get active plans for subscription billing: "+err.Error())
		} else {
			for _, up := range ups {
				if int(up.Id) == planId {
					limits := up.Plan.GetModelLimits()
					modelName := meta.OriginModelName
					rule, resolvedModel, found := dbmodel.FindLimit(limits, modelName, up.Plan.DefaultModel)
					if !found {
						continue
					}
					resolvedName := resolvedModel
				for _, windowType := range []string{dbmodel.WindowTypePeriod, dbmodel.WindowTypeWeek, dbmodel.WindowTypeMonth} {
					windowIndex := dbmodel.CalcWindowIndex(now, up.StartTime, windowType, rule.PeriodH)
					err := dbmodel.IncrementPlanUsage(int(up.Id), resolvedName, windowType, windowIndex, 1, int64(promptTokens), int64(completionTokens), int64(cachedTokens))
						if err != nil {
							logger.Error(ctx, fmt.Sprintf("failed to increment plan usage: %s", err.Error()))
						}
					}
					break
				}
			}
		}
		if preConsumedQuota != 0 {
			go func(ctx context.Context) {
				err := dbmodel.PostConsumeTokenQuota(meta.TokenId, -preConsumedQuota)
				if err != nil {
					logger.Error(ctx, "error returning pre-consumed quota for subscription: "+err.Error())
				}
			}(ctx)
		}
	} else {
		quotaDelta := quota - preConsumedQuota
		err := dbmodel.PostConsumeTokenQuota(meta.TokenId, quotaDelta)
		if err != nil {
			logger.Error(ctx, "error consuming token remain quota: "+err.Error())
		}
		err = dbmodel.CacheUpdateUserQuota(ctx, meta.UserId)
		if err != nil {
			logger.Error(ctx, "error update user quota cache: "+err.Error())
		}
	}

	logContent := fmt.Sprintf("定价：输入¥%.6f/百万tokens × %d + 输出¥%.6f/百万tokens × %d",
		priceResult.InputPrice, promptTokens, priceResult.OutputPrice, completionTokens)
	if priceResult.BillingType == dbmodel.BillingTypePerRequest {
		// per_request prices are per call, not per token — reporting the token
		// formula here would be misleading, and the token counts are irrelevant.
		// per_request 为「每次固定价」，日志不应展示 token 公式。
		logContent = fmt.Sprintf("按次计费：¥%.6f/次", priceResult.PerRequestPrice)
	} else if cachedTokens > 0 {
		logContent += fmt.Sprintf(" + 缓存¥%.6f/百万tokens × %d", priceResult.CachedPrice, cachedTokens)
	}
	if groupDiscount != 1.0 {
		logContent += fmt.Sprintf(" × 分组折扣%.2f", groupDiscount)
	}
	if billingSource == 1 {
		logContent = fmt.Sprintf("订阅计费 | %s", logContent)
	}
	dbmodel.RecordConsumeLog(ctx, &dbmodel.Log{
		UserId:            meta.UserId,
		ChannelId:         meta.ChannelId,
		PromptTokens:      promptTokens,
		CompletionTokens:  completionTokens,
		CachedTokens:      cachedTokens,
		ModelName:         meta.OriginModelName,
		TokenName:         meta.TokenName,
		Quota:             int(quota),
		Content:           logContent,
		IsStream:          meta.IsStream,
		ElapsedTime:       helper.CalcElapsedTime(meta.StartTime),
		SystemPromptReset: systemPromptReset,
		BillingSource:     billingSource,
		PlanId:            planId,
		SessionKey:        meta.SessionKey,
	})
	if billingSource == 0 {
		dbmodel.UpdateUserUsedQuotaAndRequestCount(meta.UserId, quota)
		dbmodel.UpdateChannelUsedQuota(meta.ChannelId, quota)
	} else {
		dbmodel.UpdateChannelUsedQuota(meta.ChannelId, quota)
	}
}

func getMappedModelName(modelName string, mapping map[string]string) (string, bool) {
	if mapping == nil {
		return modelName, false
	}
	mappedModelName := mapping[modelName]
	if mappedModelName != "" {
		return mappedModelName, true
	}
	return modelName, false
}

func isErrorHappened(meta *meta.Meta, resp *http.Response) bool {
	if resp == nil {
		if meta.ChannelID == "aws_claude" {
			return false
		}
		return true
	}
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated {
		return true
	}
	if meta.ChannelID == "deepl" {
		return false
	}

	if meta.IsStream && strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") &&
		meta.ChannelID != "replicate" {
		return true
	}
	return false
}

func setSystemPrompt(ctx context.Context, request *relaymodel.GeneralOpenAIRequest, prompt string) (reset bool) {
	if prompt == "" {
		return false
	}
	if len(request.Messages) == 0 {
		return false
	}
	if request.Messages[0].Role == role.System {
		request.Messages[0].Content = prompt
		logger.Infof(ctx, "rewrite system prompt")
		return true
	}
	request.Messages = append([]relaymodel.Message{{
		Role:    role.System,
		Content: prompt,
	}}, request.Messages...)
	logger.Infof(ctx, "add system prompt")
	return true
}