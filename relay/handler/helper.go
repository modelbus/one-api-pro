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

func preConsumeQuota(ctx context.Context, textRequest *relaymodel.GeneralOpenAIRequest, promptTokens int, ratio float64, meta *meta.Meta) (int64, *relaymodel.ErrorWithStatusCode) {
	if meta.PlanId > 0 {
		return 0, nil
	}

	preConsumedQuota := getPreConsumedQuota(textRequest, promptTokens, ratio)

	userQuota, err := dbmodel.CacheGetUserQuota(ctx, meta.UserId)
	if err != nil {
		return preConsumedQuota, openai.ErrorWrapper(err, "get_user_quota_failed", http.StatusInternalServerError)
	}
	if userQuota-preConsumedQuota < 0 {
		return preConsumedQuota, openai.ErrorWrapper(errors.New("user quota is not enough"), "insufficient_user_quota", http.StatusForbidden)
	}
	err = dbmodel.CacheDecreaseUserQuota(meta.UserId, preConsumedQuota)
	if err != nil {
		return preConsumedQuota, openai.ErrorWrapper(err, "decrease_user_quota_failed", http.StatusInternalServerError)
	}
	if userQuota > 100*preConsumedQuota {
		preConsumedQuota = 0
		logger.Info(ctx, fmt.Sprintf("user %d has enough quota %d, trusted and no need to pre-consume", meta.UserId, userQuota))
	}
	if preConsumedQuota > 0 {
		err := dbmodel.PreConsumeTokenQuota(meta.TokenId, preConsumedQuota)
		if err != nil {
			return preConsumedQuota, openai.ErrorWrapper(err, "pre_consume_token_quota_failed", http.StatusForbidden)
		}
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

	totalTokens := promptTokens + completionTokens
	if totalTokens == 0 {
		quota = 0
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

	logContent := fmt.Sprintf("定价：输入¥%.4f/百万tokens × %d + 输出¥%.4f/百万tokens × %d",
		priceResult.InputPrice, promptTokens, priceResult.OutputPrice, completionTokens)
	if cachedTokens > 0 {
		logContent += fmt.Sprintf(" + 缓存¥%.4f/百万tokens × %d", priceResult.CachedPrice, cachedTokens)
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