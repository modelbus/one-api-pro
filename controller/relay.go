package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/channelrouter"
	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/ctxkey"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/middleware"
	dbmodel "github.com/modelbus/one-api-pro/model"
	"github.com/modelbus/one-api-pro/monitor"
	"github.com/modelbus/one-api-pro/relay/handler"
	"github.com/modelbus/one-api-pro/relay/interceptor"
	"github.com/modelbus/one-api-pro/relay/relaymode"
	"github.com/modelbus/one-api-pro/relay/schema"
)

var errorHandlerChain *interceptor.ErrorHandlerChain

func initErrorInterceptorChain() {
	actionHandler := &interceptor.ChannelActionHandler{
		Router: channelrouter.DefaultRouter,
	}
	errorHandlerChain = interceptor.NewErrorHandlerChain(
		&interceptor.ResponseMapperHandler{},
		&interceptor.RetryHandler{},
		actionHandler,
	)
}

func relayHelper(c *gin.Context, relayMode int) *model.ErrorWithStatusCode {
	var err *model.ErrorWithStatusCode
	switch relayMode {
	case relaymode.ImagesGenerations:
		err = controller.RelayImageHelper(c, relayMode)
	case relaymode.AudioSpeech:
		fallthrough
	case relaymode.AudioTranslation:
		fallthrough
	case relaymode.AudioTranscription:
		err = controller.RelayAudioHelper(c, relayMode)
	case relaymode.Proxy:
		err = controller.RelayProxyHelper(c, relayMode)
	default:
		err = controller.RelayTextHelper(c)
	}
	return err
}

func Relay(c *gin.Context) {
	if errorHandlerChain == nil {
		initErrorInterceptorChain()
	}

	ctx := c.Request.Context()
	relayMode := relaymode.GetByPath(c.Request.URL.Path)
	if config.DebugEnabled {
		requestBody, _ := common.GetRequestBody(c)
		logger.Debugf(ctx, "request body: %s", string(requestBody))
	}
	channelId := c.GetInt(ctxkey.ChannelId)
	userId := c.GetInt(ctxkey.Id)
	bizErr := relayHelper(c, relayMode)
	if bizErr == nil {
		monitor.Emit(channelId, true)

		sessionKey := c.GetString(ctxkey.SessionKey)
		if sessionKey != "" && channelrouter.DefaultRouter != nil {
			channelrouter.DefaultRouter.SetStickySession(sessionKey, channelId)
		}

		return
	}

	lastFailedChannelId := channelId
	channelName := c.GetString(ctxkey.ChannelName)
	group := c.GetString(ctxkey.Group)
	originalModel := c.GetString(ctxkey.OriginalModel)

	cooldownSeconds := getChannelCooldownSeconds(channelId)
	errCtx := interceptor.BuildErrorContext(c, bizErr, channelId, channelName, cooldownSeconds)
	processedErr := errorHandlerChain.Process(errCtx)

	monitor.Emit(channelId, false)

	requestId := c.GetString(helper.RequestIdKey)
	retryTimes := config.RetryTimes
	if !errCtx.ShouldRetry {
		logger.Errorf(ctx, "relay error happen, status code is %d, code: %v, type: %s, message: %s, won't retry in this case", bizErr.StatusCode, bizErr.Error.Code, bizErr.Error.Type, bizErr.Error.Message)
		retryTimes = 0
	}
	for i := retryTimes; i > 0; i-- {
		ch, err := routeChannel(c, group, originalModel, userId, i != retryTimes)
		if err != nil {
			logger.Errorf(ctx, "routeChannel failed: %+v", err)
			break
		}
		logger.Infof(ctx, "using channel #%d to retry (remain times %d)", ch.Id, i)
		if ch.Id == lastFailedChannelId {
			continue
		}

		if config.ChannelConcurrencyEnabled && ch.GetMaxConcurrency() > 0 {
			if !channelrouter.DefaultRouter.TryAcquireConcurrency(ch.Id, ch.GetMaxConcurrency()) {
				continue
			}
		}

		middleware.SetupContextForSelectedChannel(c, ch, originalModel)
		requestBody, err := common.GetRequestBody(c)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		bizErr = relayHelper(c, relayMode)
		if bizErr == nil {
			if config.ChannelConcurrencyEnabled && ch.GetMaxConcurrency() > 0 {
				channelrouter.DefaultRouter.ReleaseConcurrency(ch.Id)
			}
			return
		}

		if config.ChannelConcurrencyEnabled && ch.GetMaxConcurrency() > 0 {
			channelrouter.DefaultRouter.ReleaseConcurrency(ch.Id)
		}

		lastFailedChannelId = ch.Id
		channelName = ch.Name

		retryCooldown := getChannelCooldownSeconds(ch.Id)
		errCtx = interceptor.BuildErrorContext(c, bizErr, ch.Id, channelName, retryCooldown)
		errorHandlerChain.Process(errCtx)

		monitor.Emit(ch.Id, false)
	}

	// Fallback path: if every normal channel for the requested model failed
	// (and the failure was retryable), attempt a single call on a configured
	// fallback channel. The fallback channel typically advertises a single
	// model in its `models` field — that model is what we substitute into the
	// request body before re-invoking the relay helper.
	if errCtx != nil && errCtx.ShouldRetry && bizErr != nil {
		if fbErr := tryFallbackChannel(c, group, relayMode, bizErr, lastFailedChannelId); fbErr == nil {
			return
		} else if fbErr != bizErr {
			// Fallback returned a *new* error (different status/message than the
			// upstream one). Surface it to the client so we don't silently
			// mask what happened.
			bizErr = fbErr
			processedErr = fbErr
		}
	}

	if processedErr != nil {
		processedErr.Error.Message = helper.MessageWithRequestId(processedErr.Error.Message, requestId)
		c.JSON(processedErr.StatusCode, gin.H{
			"error": processedErr.Error,
		})
	}
}

func routeChannel(c *gin.Context, group, modelName string, userId int, ignoreFirstPriority bool) (*dbmodel.Channel, error) {
	if channelrouter.DefaultRouter == nil {
		return dbmodel.CacheGetRandomSatisfiedChannel(group, modelName, ignoreFirstPriority)
	}

	candidates := dbmodel.GetChannelCandidates(group, modelName)
	if len(candidates) == 0 {
		return dbmodel.CacheGetRandomSatisfiedChannel(group, modelName, ignoreFirstPriority)
	}

	sessionKey := ""
	if config.ChannelStickySessionEnabled && userId > 0 && modelName != "" {
		sessionKey = channelrouter.MakeSessionKey(userId, modelName)
	}

	req := &channelrouter.RouteRequest{
		Group:               group,
		Model:               modelName,
		UserId:              userId,
		IgnoreFirstPriority: ignoreFirstPriority,
		SessionKey:          sessionKey,
	}

	ch, err := channelrouter.DefaultRouter.Route(c.Request.Context(), req, candidates)
	if err != nil {
		return dbmodel.CacheGetRandomSatisfiedChannel(group, modelName, ignoreFirstPriority)
	}
	return ch, nil
}

func getChannelCooldownSeconds(channelId int) int {
	if ch, ok := dbmodel.CacheGetChannelById(channelId); ok {
		cooldownSeconds := ch.CooldownSeconds
		if cooldownSeconds <= 0 {
			cooldownSeconds = config.ChannelDefaultCooldownSeconds
		}
		return cooldownSeconds
	}
	ch, err := dbmodel.GetChannelById(channelId, false)
	if err != nil {
		return config.ChannelDefaultCooldownSeconds
	}
	cooldownSeconds := ch.CooldownSeconds
	if cooldownSeconds <= 0 {
		cooldownSeconds = config.ChannelDefaultCooldownSeconds
	}
	return cooldownSeconds
}

// tryFallbackChannel is invoked after the normal retry loop has been exhausted
// with a retryable error. It attempts a single call on a configured fallback
// channel (channels.is_fallback = 1) in the user's group, rewriting the
// request body's `model` field to the fallback channel's primary model.
//
// Returns:
//   - nil on success (the response was written to the client)
//   - the same `prevErr` (or a new error) when no fallback applies / fallback
//     itself failed; caller should surface the returned error to the client.
func tryFallbackChannel(c *gin.Context, group string, relayMode int, prevErr *model.ErrorWithStatusCode, lastFailedChannelId int) *model.ErrorWithStatusCode {
	fb, err := dbmodel.GetFallbackChannel(group)
	if err != nil {
		logger.Errorf(c.Request.Context(), "GetFallbackChannel error: %+v", err)
		return prevErr
	}
	if fb == nil {
		return prevErr
	}
	if fb.Id == lastFailedChannelId {
		// Don't re-attempt a channel that just failed as a "fallback" — pick
		// another one if available.
		logger.Debugf(c.Request.Context(), "fallback channel #%d equals lastFailedChannelId, skipping", fb.Id)
		return prevErr
	}

	fbModel := pickFallbackModel(fb)
	if fbModel == "" {
		logger.Debugf(c.Request.Context(), "fallback channel #%d has no models configured, skipping", fb.Id)
		return prevErr
	}

	// Rewrite request body to swap the model name for the fallback channel.
	if relayMode != relaymode.Proxy {
		if rerr := rewriteRequestModel(c, fbModel); rerr != nil {
			logger.Errorf(c.Request.Context(), "rewriteRequestModel failed: %+v", rerr)
			return prevErr
		}
	}

	// Re-bind gin context for the fallback channel (auth headers, base URL,
	// model mapping, etc.). We pass the fallback model name so RequestModel /
	// OriginalModel in the context reflect the actually-sent model.
	middleware.SetupContextForSelectedChannel(c, fb, fbModel)

	// Restore the request body so the helper re-reads it.
	requestBody, _ := common.GetRequestBody(c)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	logger.Infof(c.Request.Context(), "fallback: using channel #%d model=%s after original model failure", fb.Id, fbModel)

	fbErr := relayHelper(c, relayMode)
	if fbErr == nil {
		monitor.Emit(fb.Id, true)
		return nil
	}
	monitor.Emit(fb.Id, false)
	logger.Errorf(c.Request.Context(), "fallback channel #%d also failed: status=%d code=%v type=%s message=%s",
		fb.Id, fbErr.StatusCode, fbErr.Error.Code, fbErr.Error.Type, fbErr.Error.Message)
	return fbErr
}

// pickFallbackModel returns the first non-empty model from a fallback channel's
// `models` field (comma-separated). Returns "" if no model is configured.
func pickFallbackModel(ch *dbmodel.Channel) string {
	parts := strings.Split(ch.Models, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			return p
		}
	}
	return ""
}

// rewriteRequestModel rewrites the JSON `model` field in the cached request
// body. It also updates ctxkey.RequestModel so downstream code reads the new
// value consistently.
func rewriteRequestModel(c *gin.Context, newModel string) error {
	body, err := common.GetRequestBody(c)
	if err != nil {
		return err
	}
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		// Non-JSON payloads: only safe option is to leave the body alone and
		// rely on channel-side model_mapping to translate. Bail out so the
		// caller knows no rewrite happened.
		c.Set(ctxkey.RequestModel, newModel)
		return nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return fmt.Errorf("request body is not a JSON object: %w", err)
	}
	modelJSON, _ := json.Marshal(newModel)
	m["model"] = modelJSON
	newBody, err := json.Marshal(m)
	if err != nil {
		return err
	}
	c.Set(ctxkey.KeyRequestBody, newBody)
	c.Set(ctxkey.RequestModel, newModel)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
	return nil
}

func RelayNotImplemented(c *gin.Context) {
	err := model.Error{
		Message: "API not implemented",
		Type:    "one_api_error",
		Param:   "",
		Code:    "api_not_implemented",
	}
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": err,
	})
}

func RelayNotFound(c *gin.Context) {
	err := model.Error{
		Message: fmt.Sprintf("Invalid URL (%s %s)", c.Request.Method, c.Request.URL.Path),
		Type:    "invalid_request_error",
		Param:   "",
		Code:    "",
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": err,
	})
}