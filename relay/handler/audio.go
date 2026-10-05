package controller

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/client"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/ctxkey"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/billing"
	billingratio "github.com/modelbus/one-api-pro/relay/billing/ratio"
	"github.com/modelbus/one-api-pro/relay/meta"
	"github.com/modelbus/one-api-pro/relay/relaymode"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

func RelayAudioHelper(c *gin.Context, relayMode int) *relaymodel.ErrorWithStatusCode {
	ctx := c.Request.Context()
	meta := meta.GetByContext(c)
	audioModel := "whisper-1"

	tokenId := c.GetInt(ctxkey.TokenId)
	channelId := c.GetInt(ctxkey.ChannelId)
	userId := c.GetInt(ctxkey.Id)
	group := c.GetString(ctxkey.Group)
	tokenName := c.GetString(ctxkey.TokenName)

	var ttsRequest openai.TextToSpeechRequest
	if relayMode == relaymode.AudioSpeech {
		err := common.UnmarshalBodyReusable(c, &ttsRequest)
		if err != nil {
			return openai.ErrorWrapper(err, "invalid_json", http.StatusBadRequest)
		}
		audioModel = ttsRequest.Model
		if len(ttsRequest.Input) > 4096 {
			return openai.ErrorWrapper(errors.New("input is too long (over 4096 characters)"), "text_too_long", http.StatusBadRequest)
		}
	}

	originAudioModel := audioModel
	modelMapping := c.GetStringMapString(ctxkey.ModelMapping)
	if modelMapping != nil && modelMapping[audioModel] != "" {
		audioModel = modelMapping[audioModel]
	}

	priceResult, err := billingratio.GetModelPrice(audioModel, originAudioModel)
	if err != nil {
		return openai.ErrorWrapper(err, "model_price_not_found", http.StatusUnprocessableEntity)
	}
	groupDiscount := billingratio.GetGroupDiscount(group, audioModel, originAudioModel)

	var quota int64
	var preConsumedQuota int64
	var bizErr *relaymodel.ErrorWithStatusCode
	if priceResult.BillingType == model.BillingTypePerRequest {
		// per_request：按次计费，费用在请求前已完全确定，精确预扣。
		// 旧实现按 token 占比估算预扣、却按次全额结算，单次即可透支余额。
		// per_request: flat price known upfront — pre-consume the exact amount
		// so the settlement delta is 0 and the balance cannot go negative.
		quota = billingratio.CalculatePerRequestQuota(priceResult.PerRequestPrice, 1, 1, groupDiscount)
		preConsumedQuota, bizErr = preConsumeExactQuota(ctx, quota, meta)
	} else {
		switch relayMode {
		case relaymode.AudioSpeech:
			ratio := (priceResult.InputPrice + priceResult.OutputPrice) / 2.0 / billingratio.Million * config.QuotaPerUnit
			if ratio == 0 {
				ratio = 1.0
			}
			// 与 CalculateTokenQuota 一致：用 math.Round 而非截断。
			quota = int64(math.Round(float64(len(ttsRequest.Input)) * ratio))
			// 语音合成费用由输入文本长度在请求前确定，预扣即结算。
			preConsumedQuota, bizErr = preConsumeExactQuota(ctx, quota, meta)
		default:
			// PreConsumedQuota 是「额外预留的额度（微元）」，不是 token 数，
			// 不再乘 ratio（历史实现把它当 token 数会放大 ratio 倍）。
			// PreConsumedQuota is an extra reserved quota in micro-quota, not a token count.
			preConsumedQuota, bizErr = preConsumeAmount(ctx, config.PreConsumedQuota, meta)
		}
	}
	if bizErr != nil {
		return bizErr
	}
	succeed := false
	defer func() {
		if succeed {
			return
		}
		if preConsumedQuota > 0 {
			defer func(ctx context.Context) {
				go func() {
					err := model.PostConsumeTokenQuota(tokenId, -preConsumedQuota)
					if err != nil {
						logger.Error(ctx, fmt.Sprintf("error rollback pre-consumed quota: %s", err.Error()))
					}
				}()
			}(c.Request.Context())
		}
	}()

	baseURL := meta.BaseURL
	requestURL := c.Request.URL.String()
	if c.GetString(ctxkey.BaseURL) != "" {
		baseURL = c.GetString(ctxkey.BaseURL)
	}

	fullRequestURL := openai.GetFullRequestURL(baseURL, requestURL, meta.ChannelID)
	if meta.ChannelID == "azure" {
		apiVersion := meta.Config.APIVersion
		if relayMode == relaymode.AudioTranscription {
			fullRequestURL = fmt.Sprintf("%s/openai/deployments/%s/audio/transcriptions?api-version=%s", baseURL, audioModel, apiVersion)
		} else if relayMode == relaymode.AudioSpeech {
			fullRequestURL = fmt.Sprintf("%s/openai/deployments/%s/audio/speech?api-version=%s", baseURL, audioModel, apiVersion)
		}
	}

	requestBody := &bytes.Buffer{}
	_, err = io.Copy(requestBody, c.Request.Body)
	if err != nil {
		return openai.ErrorWrapper(err, "new_request_body_failed", http.StatusInternalServerError)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody.Bytes()))
	responseFormat := c.DefaultPostForm("response_format", "json")

	req, err := http.NewRequest(c.Request.Method, fullRequestURL, requestBody)
	if err != nil {
		return openai.ErrorWrapper(err, "new_request_failed", http.StatusInternalServerError)
	}

	if (relayMode == relaymode.AudioTranscription || relayMode == relaymode.AudioSpeech) && meta.ChannelID == "azure" {
		apiKey := c.Request.Header.Get("Authorization")
		apiKey = strings.TrimPrefix(apiKey, "Bearer ")
		req.Header.Set("api-key", apiKey)
		req.ContentLength = c.Request.ContentLength
	} else {
		req.Header.Set("Authorization", c.Request.Header.Get("Authorization"))
	}
	req.Header.Set("Content-Type", c.Request.Header.Get("Content-Type"))
	req.Header.Set("Accept", c.Request.Header.Get("Accept"))

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return openai.ErrorWrapper(err, "do_request_failed", http.StatusInternalServerError)
	}

	err = req.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_request_body_failed", http.StatusInternalServerError)
	}
	err = c.Request.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_request_body_failed", http.StatusInternalServerError)
	}

	if relayMode != relaymode.AudioSpeech {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return openai.ErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError)
		}
		err = resp.Body.Close()
		if err != nil {
			return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError)
		}

		var openAIErr openai.SlimTextResponse
		if err = json.Unmarshal(responseBody, &openAIErr); err == nil {
			if openAIErr.Error.Message != "" {
				return openai.ErrorWrapper(fmt.Errorf("type %s, code %v, message %s", openAIErr.Error.Type, openAIErr.Error.Code, openAIErr.Error.Message), "request_error", http.StatusInternalServerError)
			}
		}

		var text string
		switch responseFormat {
		case "json":
			text, err = getTextFromJSON(responseBody)
		case "text":
			text, err = getTextFromText(responseBody)
		case "srt":
			text, err = getTextFromSRT(responseBody)
		case "verbose_json":
			text, err = getTextFromVerboseJSON(responseBody)
		case "vtt":
			text, err = getTextFromVTT(responseBody)
		default:
			return openai.ErrorWrapper(errors.New("unexpected_response_format"), "unexpected_response_format", http.StatusInternalServerError)
		}
		if err != nil {
			return openai.ErrorWrapper(err, "get_text_from_body_err", http.StatusInternalServerError)
		}
		// per_request：费用为每次固定价，与转写文本长度无关，结算时
		// 必须沿用按次价，否则会又变回按 token 扣费。
		// per_request: flat per-call price, independent of transcript length.
		if priceResult.BillingType == model.BillingTypePerRequest {
			quota = billingratio.CalculatePerRequestQuota(priceResult.PerRequestPrice, 1, 1, groupDiscount)
		} else {
			// 与 CalculateTokenQuota 保持一致：math.Round 取整 + 最小 1 微元兜底，
			// 避免极短转写文本（或空文本）导致 0 扣费。
			tokens := openai.CountTokenText(text, audioModel)
			quota = int64(math.Round(float64(tokens) * priceResult.InputPrice / billingratio.Million * config.QuotaPerUnit * groupDiscount))
			if quota <= 0 && tokens > 0 {
				quota = 1
			}
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
	}
	if resp.StatusCode != http.StatusOK {
		return RelayErrorHandler(resp)
	}
	succeed = true
	quotaDelta := quota - preConsumedQuota
	defer func(ctx context.Context) {
		go billing.PostConsumeQuota(ctx, &billing.ConsumeQuotaParams{
			TokenId:         tokenId,
			UserId:          userId,
			ChannelId:       channelId,
			QuotaDelta:      quotaDelta,
			TotalQuota:      quota,
			ModelName:       originAudioModel,
			TokenName:       tokenName,
			InputPrice:      priceResult.InputPrice,
			OutputPrice:     priceResult.OutputPrice,
			CachedPrice:     priceResult.CachedPrice,
			PerRequestPrice: priceResult.PerRequestPrice,
			GroupDiscount:   groupDiscount,
			BillingType:     priceResult.BillingType,
		})
	}(c.Request.Context())

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Writer.WriteHeader(resp.StatusCode)

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		return openai.ErrorWrapper(err, "copy_response_body_failed", http.StatusInternalServerError)
	}
	err = resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError)
	}
	return nil
}

func getTextFromVTT(body []byte) (string, error) {
	return getTextFromSRT(body)
}

func getTextFromVerboseJSON(body []byte) (string, error) {
	var whisperResponse openai.WhisperVerboseJSONResponse
	if err := json.Unmarshal(body, &whisperResponse); err != nil {
		return "", fmt.Errorf("unmarshal_response_body_failed err :%w", err)
	}
	return whisperResponse.Text, nil
}

func getTextFromSRT(body []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	var builder strings.Builder
	var textLine bool
	for scanner.Scan() {
		line := scanner.Text()
		if textLine {
			builder.WriteString(line)
			textLine = false
			continue
		} else if strings.Contains(line, "-->") {
			textLine = true
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func getTextFromText(body []byte) (string, error) {
	return strings.TrimSuffix(string(body), "\n"), nil
}

func getTextFromJSON(body []byte) (string, error) {
	var whisperResponse openai.WhisperJSONResponse
	if err := json.Unmarshal(body, &whisperResponse); err != nil {
		return "", fmt.Errorf("unmarshal_response_body_failed err :%w", err)
	}
	return whisperResponse.Text, nil
}
