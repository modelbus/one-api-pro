package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/channelrouter"
	"github.com/Leon-PanPan/one-api-pro/common"
	"github.com/Leon-PanPan/one-api-pro/common/config"
	"github.com/Leon-PanPan/one-api-pro/common/ctxkey"
	"github.com/Leon-PanPan/one-api-pro/common/helper"
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/Leon-PanPan/one-api-pro/middleware"
	dbmodel "github.com/Leon-PanPan/one-api-pro/model"
	"github.com/Leon-PanPan/one-api-pro/monitor"
	"github.com/Leon-PanPan/one-api-pro/relay/handler"
	"github.com/Leon-PanPan/one-api-pro/relay/interceptor"
	"github.com/Leon-PanPan/one-api-pro/relay/schema"
	"github.com/Leon-PanPan/one-api-pro/relay/relaymode"
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