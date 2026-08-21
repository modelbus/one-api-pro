package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/ctxkey"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
	"github.com/modelbus/one-api-pro/relay"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	billingratio "github.com/modelbus/one-api-pro/relay/billing/ratio"
	"github.com/modelbus/one-api-pro/relay/meta"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

func getImageRequest(c *gin.Context, _ int) (*relaymodel.ImageRequest, error) {
	imageRequest := &relaymodel.ImageRequest{}
	err := common.UnmarshalBodyReusable(c, imageRequest)
	if err != nil {
		return nil, err
	}
	if imageRequest.N == 0 {
		imageRequest.N = 1
	}
	if imageRequest.Size == "" {
		imageRequest.Size = "1024x1024"
	}
	if imageRequest.Model == "" {
		imageRequest.Model = "dall-e-2"
	}
	return imageRequest, nil
}

func isValidImageSize(model string, size string) bool {
	if model == "cogview-3" || billingratio.ImageSizeRatios[model] == nil {
		return true
	}
	_, ok := billingratio.ImageSizeRatios[model][size]
	return ok
}

func isValidImagePromptLength(model string, promptLength int) bool {
	maxPromptLength, ok := billingratio.ImagePromptLengthLimitations[model]
	return !ok || promptLength <= maxPromptLength
}

func isWithinRange(element string, value int) bool {
	amounts, ok := billingratio.ImageGenerationAmounts[element]
	return !ok || (value >= amounts[0] && value <= amounts[1])
}

func getImageSizeRatio(model string, size string) float64 {
	if ratio, ok := billingratio.ImageSizeRatios[model][size]; ok {
		return ratio
	}
	return 1
}

func validateImageRequest(imageRequest *relaymodel.ImageRequest, _ *meta.Meta) *relaymodel.ErrorWithStatusCode {
	if imageRequest.Prompt == "" {
		return openai.ErrorWrapper(errors.New("prompt is required"), "prompt_missing", http.StatusBadRequest)
	}

	if !isValidImageSize(imageRequest.Model, imageRequest.Size) {
		return openai.ErrorWrapper(errors.New("size not supported for this image model"), "size_not_supported", http.StatusBadRequest)
	}

	if !isValidImagePromptLength(imageRequest.Model, len(imageRequest.Prompt)) {
		return openai.ErrorWrapper(errors.New("prompt is too long"), "prompt_too_long", http.StatusBadRequest)
	}

	if !isWithinRange(imageRequest.Model, imageRequest.N) {
		return openai.ErrorWrapper(errors.New("invalid value of n"), "n_not_within_range", http.StatusBadRequest)
	}
	return nil
}

func getImageCostRatio(imageRequest *relaymodel.ImageRequest) (float64, error) {
	if imageRequest == nil {
		return 0, errors.New("imageRequest is nil")
	}
	imageCostRatio := getImageSizeRatio(imageRequest.Model, imageRequest.Size)
	if imageRequest.Quality == "hd" && imageRequest.Model == "dall-e-3" {
		if imageRequest.Size == "1024x1024" {
			imageCostRatio *= 2
		} else {
			imageCostRatio *= 1.5
		}
	}
	return imageCostRatio, nil
}

func RelayImageHelper(c *gin.Context, relayMode int) *relaymodel.ErrorWithStatusCode {
	ctx := c.Request.Context()
	meta := meta.GetByContext(c)
	imageRequest, err := getImageRequest(c, meta.Mode)
	if err != nil {
		logger.Errorf(ctx, "getImageRequest failed: %s", err.Error())
		return openai.ErrorWrapper(err, "invalid_image_request", http.StatusBadRequest)
	}

	var isModelMapped bool
	if meta.DefaultModel != "" {
		imageRequest.Model = meta.DefaultModel
	}
	meta.OriginModelName = imageRequest.Model
	imageRequest.Model, isModelMapped = getMappedModelName(imageRequest.Model, meta.ModelMapping)
	meta.ActualModelName = imageRequest.Model

	bizErr := validateImageRequest(imageRequest, meta)
	if bizErr != nil {
		return bizErr
	}

	imageCostRatio, err := getImageCostRatio(imageRequest)
	if err != nil {
		return openai.ErrorWrapper(err, "get_image_cost_ratio_failed", http.StatusInternalServerError)
	}

	imageModel := imageRequest.Model
	imageRequest.Model, _ = getMappedModelName(imageRequest.Model, billingratio.ImageOriginModelName)
	c.Set("response_format", imageRequest.ResponseFormat)

	priceResult, priceErr := billingratio.GetModelPrice(imageModel, meta.OriginModelName)
	if priceErr != nil {
		return openai.ErrorWrapper(priceErr, "model_price_not_found", http.StatusUnprocessableEntity)
	}
	groupDiscount := billingratio.GetGroupDiscount(meta.Group, imageModel, meta.OriginModelName)

	var quota int64
	if priceResult.BillingType == model.BillingTypePerRequest {
		quota = billingratio.CalculatePerRequestQuota(priceResult.PerRequestPrice, imageCostRatio, imageRequest.N, groupDiscount)
	} else {
		quota = billingratio.CalculateTokenQuota(
			priceResult.InputPrice, priceResult.OutputPrice, priceResult.CachedPrice,
			1, 0, 0,
			groupDiscount,
		)
		quota = int64(float64(quota) * imageCostRatio * float64(imageRequest.N))
	}

	if meta.PlanId == 0 {
		userQuota, err := model.CacheGetUserQuota(ctx, meta.UserId)
		if err != nil {
			return openai.ErrorWrapper(err, "get_user_quota_failed", http.StatusInternalServerError)
		}
		if userQuota-quota < 0 {
			return openai.ErrorWrapper(errors.New("user quota is not enough"), "insufficient_user_quota", http.StatusForbidden)
		}
	}

	var requestBody io.Reader
	if isModelMapped || meta.ChannelID == "azure" {
		jsonStr, err := json.Marshal(imageRequest)
		if err != nil {
			return openai.ErrorWrapper(err, "marshal_image_request_failed", http.StatusInternalServerError)
		}
		requestBody = bytes.NewBuffer(jsonStr)
	} else {
		requestBody = c.Request.Body
	}

	adaptor := relay.GetAdaptorByChannelID(meta.ChannelID)
	if adaptor == nil {
		return openai.ErrorWrapper(fmt.Errorf("invalid channel type: %s", meta.ChannelID), "invalid_channel_type", http.StatusBadRequest)
	}
	adaptor.Init(meta)

	switch meta.ChannelID {
	case "zhipu", "ali", "replicate", "baidu":
		finalRequest, err := adaptor.ConvertImageRequest(imageRequest)
		if err != nil {
			return openai.ErrorWrapper(err, "convert_image_request_failed", http.StatusInternalServerError)
		}
		jsonStr, err := json.Marshal(finalRequest)
		if err != nil {
			return openai.ErrorWrapper(err, "marshal_image_request_failed", http.StatusInternalServerError)
		}
		requestBody = bytes.NewBuffer(jsonStr)
	}

	resp, err := adaptor.DoRequest(c, meta, requestBody)
	if err != nil {
		logger.Errorf(ctx, "DoRequest failed: %s", err.Error())
		return openai.ErrorWrapper(err, "do_request_failed", http.StatusInternalServerError)
	}

	defer func(ctx context.Context) {
		if resp != nil &&
			resp.StatusCode != http.StatusCreated &&
			resp.StatusCode != http.StatusOK {
			return
		}

		if meta.PlanId > 0 {
			now := helper.GetTimestamp()
			ups, upsErr := model.CacheGetUserActivePlans(meta.UserId)
			if upsErr == nil {
				for _, up := range ups {
					if int(up.Id) == meta.PlanId && up.Plan != nil {
						limits := up.Plan.GetModelLimits()
						rule, resolvedModel, found := model.FindLimit(limits, meta.OriginModelName, up.Plan.DefaultModel)
						if !found {
							continue
						}
						for _, windowType := range []string{model.WindowTypePeriod, model.WindowTypeWeek, model.WindowTypeMonth} {
							windowIndex := model.CalcWindowIndex(now, up.StartTime, windowType, rule.PeriodH)
							_ = model.IncrementPlanUsage(int(up.Id), resolvedModel, windowType, windowIndex, 1, 0, 0, 0)
						}
						break
					}
				}
			}
		} else {
			err := model.PostConsumeTokenQuota(meta.TokenId, quota)
			if err != nil {
				logger.SysError("error consuming token remain quota: " + err.Error())
			}
			err = model.CacheUpdateUserQuota(ctx, meta.UserId)
			if err != nil {
				logger.SysError("error update user quota cache: " + err.Error())
			}
			model.UpdateUserUsedQuotaAndRequestCount(meta.UserId, quota)
		}
		channelId := c.GetInt(ctxkey.ChannelId)
		model.UpdateChannelUsedQuota(channelId, quota)
		if quota != 0 {
			tokenName := c.GetString(ctxkey.TokenName)
			logContent := fmt.Sprintf("按次计费：¥%.4f/次", priceResult.PerRequestPrice)
			if priceResult.BillingType == model.BillingTypeToken {
				logContent = fmt.Sprintf("定价：输入¥%.4f/百万tokens × 1", priceResult.InputPrice)
			}
			if groupDiscount != 1.0 {
				logContent += fmt.Sprintf(" × 分组折扣%.2f", groupDiscount)
			}
			if meta.PlanId > 0 {
				logContent = fmt.Sprintf("订阅计费 | %s", logContent)
			}
			model.RecordConsumeLog(ctx, &model.Log{
				UserId:            meta.UserId,
				ChannelId:         meta.ChannelId,
				PromptTokens:      0,
				CompletionTokens:  0,
				CachedTokens:      0,
				ModelName:        meta.OriginModelName,
				TokenName:        tokenName,
				Quota:            int(quota),
				Content:          logContent,
				BillingSource:    func() int { if meta.PlanId > 0 { return 1 }; return 0 }(),
				PlanId:           meta.PlanId,
			})
		}
	}(c.Request.Context())

	_, respErr := adaptor.DoResponse(c, resp, meta)
	if respErr != nil {
		logger.Errorf(ctx, "respErr is not nil: %+v", respErr)
		return respErr
	}

	return nil
}