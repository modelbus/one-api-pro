package interceptor

import (
	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/common/config"
	"github.com/Leon-PanPan/one-api-pro/relay/errcode"
	relaymodel "github.com/Leon-PanPan/one-api-pro/relay/schema"
)

type ErrorContext struct {
	GinContext      *gin.Context
	StatusCode      int
	OriginalError   *relaymodel.ErrorWithStatusCode
	MappedError     *relaymodel.ErrorWithStatusCode
	ChannelId       int
	ChannelName     string
	CooldownSeconds int
	ShouldRetry     bool
	ChannelAction   string
}

type ErrorHandler interface {
	Name() string
	Handle(ctx *ErrorContext) bool
}

type ErrorHandlerChain struct {
	handlers []ErrorHandler
}

func NewErrorHandlerChain(handlers ...ErrorHandler) *ErrorHandlerChain {
	return &ErrorHandlerChain{handlers: handlers}
}

func (chain *ErrorHandlerChain) Process(ctx *ErrorContext) *relaymodel.ErrorWithStatusCode {
	for _, handler := range chain.handlers {
		continueChain := handler.Handle(ctx)
		if !continueChain {
			break
		}
	}
	if ctx.MappedError != nil {
		return ctx.MappedError
	}
	return ctx.OriginalError
}

func BuildErrorContext(c *gin.Context, originalErr *relaymodel.ErrorWithStatusCode, channelId int, channelName string, cooldownSeconds int) *ErrorContext {
	return &ErrorContext{
		GinContext:      c,
		StatusCode:      originalErr.StatusCode,
		OriginalError:   originalErr,
		MappedError:     nil,
		ChannelId:       channelId,
		ChannelName:     channelName,
		CooldownSeconds: cooldownSeconds,
		ShouldRetry:     false,
		ChannelAction:   "",
	}
}

func getEffectiveConfig(statusCode int) errcode.StatusCodeConfig {
	defaultCfg := errcode.GetDefaultConfig(statusCode)
	en := config.ErrorNext

	switch defaultCfg.Category {
	case errcode.CategoryPassthrough:
		if !en.Passthrough {
			defaultCfg.Action = errcode.ActionRetry
		}
	case errcode.CategoryDisable:
		if !en.Disable {
			defaultCfg.Action = errcode.ActionRetry
		}
	case errcode.CategoryCooldown:
		if !en.Cooldown {
			defaultCfg.Action = errcode.ActionRetry
		}
	case errcode.CategoryRetry:
		if !en.Retry {
			defaultCfg.Action = errcode.ActionReturn
		}
	}
	return defaultCfg
}