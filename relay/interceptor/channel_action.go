package interceptor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Leon-PanPan/one-api-pro/channelrouter"
	"github.com/Leon-PanPan/one-api-pro/common/config"
	dbmodel "github.com/Leon-PanPan/one-api-pro/model"
	"github.com/Leon-PanPan/one-api-pro/monitor"
	"github.com/Leon-PanPan/one-api-pro/relay/errcode"
)

type ChannelActionHandler struct {
	Router *channelrouter.ChannelRouter
}

func (h *ChannelActionHandler) Name() string {
	return "channel_action"
}

func (h *ChannelActionHandler) Handle(ctx *ErrorContext) bool {
	effectiveCfg := getEffectiveConfig(ctx.StatusCode)

	errMsg := h.buildErrorMessage(ctx)
	dbmodel.UpdateChannelLastError(ctx.ChannelId, errMsg)

	switch effectiveCfg.Action {
	case errcode.ActionDisable:
		if config.AutomaticDisableChannelEnabled {
			monitor.DisableChannel(ctx.ChannelId, ctx.ChannelName, h.getDisableReason(ctx))
		}
	case errcode.ActionCooldown:
		seconds := effectiveCfg.CooldownSeconds
		if seconds <= 0 {
			seconds = ctx.CooldownSeconds
		}
		if seconds <= 0 {
			seconds = config.ChannelDefaultCooldownSeconds
		}
		if maxSec := config.ChannelMaxCooldownSeconds; maxSec > 0 && seconds > maxSec {
			seconds = maxSec
		}
		if h.Router != nil {
			h.Router.SetCooldown(ctx.ChannelId, seconds, fmt.Sprintf("upstream_%d", ctx.StatusCode), ctx.StatusCode)
		}
	}

	if ctx.StatusCode == http.StatusUnauthorized && config.AutomaticDisableChannelEnabled {
		monitor.DisableChannel(ctx.ChannelId, ctx.ChannelName, h.getDisableReason(ctx))
		return true
	}

	if h.shouldDisableBasedOnErrorContent(ctx) {
		monitor.DisableChannel(ctx.ChannelId, ctx.ChannelName, h.getDisableReasonFromContent(ctx))
	}

	return true
}

func (h *ChannelActionHandler) buildErrorMessage(ctx *ErrorContext) string {
	if ctx.OriginalError != nil && ctx.OriginalError.Error.Message != "" {
		msg := ctx.OriginalError.Error.Message
		if len(msg) > 500 {
			msg = msg[:500]
		}
		return fmt.Sprintf("[%d] %s", ctx.StatusCode, msg)
	}
	return fmt.Sprintf("upstream returned %d", ctx.StatusCode)
}

func (h *ChannelActionHandler) shouldDisableBasedOnErrorContent(ctx *ErrorContext) bool {
	if !config.AutomaticDisableChannelEnabled {
		return false
	}
	if ctx.OriginalError == nil {
		return false
	}
	err := &ctx.OriginalError.Error
	switch err.Type {
	case "insufficient_quota", "authentication_error", "permission_error", "forbidden":
		return true
	}
	if err.Code == "invalid_api_key" || err.Code == "account_deactivated" {
		return true
	}
	lowerMessage := strings.ToLower(err.Message)
	disabledKeywords := []string{
		"your access was terminated",
		"violation of our policies",
		"your credit balance is too low",
		"organization has been disabled",
		"permission denied",
		"organization has been restricted",
		"api key not valid",
		"api key expired",
		"已欠费",
	}
	for _, keyword := range disabledKeywords {
		if strings.Contains(lowerMessage, keyword) {
			return true
		}
	}
	if strings.Contains(lowerMessage, "credit") || strings.Contains(lowerMessage, "balance") {
		return true
	}
	return false
}

func (h *ChannelActionHandler) getDisableReason(ctx *ErrorContext) string {
	return fmt.Sprintf("upstream returned %d (auto-disabled by error policy)", ctx.StatusCode)
}

func (h *ChannelActionHandler) getDisableReasonFromContent(ctx *ErrorContext) string {
	if ctx.OriginalError != nil && ctx.OriginalError.Error.Message != "" {
		return ctx.OriginalError.Error.Message
	}
	return fmt.Sprintf("upstream error (status %d)", ctx.StatusCode)
}