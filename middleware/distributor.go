package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Leon-PanPan/one-api-pro/channelrouter"
	"github.com/Leon-PanPan/one-api-pro/common/config"
	"github.com/Leon-PanPan/one-api-pro/common/ctxkey"
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/Leon-PanPan/one-api-pro/model"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

type ModelRequest struct {
	Model string `json:"model" form:"model"`
}

func Distribute() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userId := c.GetInt(ctxkey.Id)
		userGroup, _ := model.CacheGetUserGroup(userId)
		c.Set(ctxkey.Group, userGroup)
		var requestModel string
		var channel *model.Channel
		channelId, ok := c.Get(ctxkey.SpecificChannelId)
		if ok {
			id, err := strconv.Atoi(channelId.(string))
			if err != nil {
				abortWithMessage(c, http.StatusBadRequest, "无效的渠道 Id")
				return
			}
			channel, err = model.GetChannelById(id, true)
			if err != nil {
				abortWithMessage(c, http.StatusBadRequest, "无效的渠道 Id")
				return
			}
			if channel.Status != model.ChannelStatusEnabled {
				abortWithMessage(c, http.StatusForbidden, "该渠道已被禁用")
				return
			}
		} else {
			requestModel = c.GetString(ctxkey.RequestModel)
			var err error
			channel, err = selectChannelViaRouter(c, userGroup, requestModel, userId, false)
			if err != nil {
				message := fmt.Sprintf("当前分组 %s 下对于模型 %s 无可用渠道", userGroup, requestModel)
				if channel != nil {
					logger.SysError(fmt.Sprintf("渠道不存在：%d", channel.Id))
					message = "数据库一致性已被破坏，请联系管理员"
				}
				abortWithMessage(c, http.StatusServiceUnavailable, message)
				return
			}
		}
		logger.Debugf(ctx, "user id %d, user group: %s, request model: %s, using channel #%d", userId, userGroup, requestModel, channel.Id)

		sessionKey := ""
		if config.ChannelStickySessionEnabled && userId > 0 && requestModel != "" {
			sessionKey = channelrouter.MakeSessionKey(userId, requestModel)
			c.Set(ctxkey.SessionKey, sessionKey)
		}

		if config.ChannelConcurrencyEnabled && channel.GetMaxConcurrency() > 0 {
			if !channelrouter.DefaultRouter.TryAcquireConcurrency(channel.Id, channel.GetMaxConcurrency()) {
				logger.Debugf(ctx, "channel #%d at concurrency limit (%d), trying alternate", channel.Id, channel.GetMaxConcurrency())
				altChannel, altErr := selectChannelViaRouter(c, userGroup, requestModel, userId, false)
				if altErr == nil && altChannel.Id != channel.Id {
					channel = altChannel
				}
			}
			c.Set(ctxkey.ConcurrencyAcquired, true)
		}

		SetupContextForSelectedChannel(c, channel, requestModel)

		if channelrouter.DefaultRouter != nil && channel.GetRPM() > 0 {
			channelrouter.DefaultRouter.IncrementRPM(channel.Id)
		}

		if sessionKey != "" {
			channelrouter.DefaultRouter.SetStickySession(sessionKey, channel.Id)
		}

		c.Next()

		if val, exists := c.Get(ctxkey.ConcurrencyAcquired); exists && val.(bool) {
			channelrouter.DefaultRouter.ReleaseConcurrency(channel.Id)
		}
	}
}

func selectChannelViaRouter(c *gin.Context, group, modelName string, userId int, ignoreFirstPriority bool) (*model.Channel, error) {
	if channelrouter.DefaultRouter == nil {
		return model.CacheGetRandomSatisfiedChannel(group, modelName, ignoreFirstPriority)
	}

	sessionKey := ""
	if config.ChannelStickySessionEnabled && userId > 0 && modelName != "" {
		sessionKey = channelrouter.MakeSessionKey(userId, modelName)
	}

	candidates := model.GetChannelCandidates(group, modelName)
	if len(candidates) == 0 {
		return model.CacheGetRandomSatisfiedChannel(group, modelName, ignoreFirstPriority)
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
		return model.CacheGetRandomSatisfiedChannel(group, modelName, ignoreFirstPriority)
	}
	return ch, nil
}

func SetupContextForSelectedChannel(c *gin.Context, channel *model.Channel, modelName string) {
	c.Set(ctxkey.Channel, channel.Type)
	c.Set(ctxkey.ChannelId, channel.Id)
	c.Set(ctxkey.ChannelName, channel.Name)
	if channel.SystemPrompt != nil && *channel.SystemPrompt != "" {
		c.Set(ctxkey.SystemPrompt, *channel.SystemPrompt)
	}
	c.Set(ctxkey.ModelMapping, channel.GetModelMapping())
	c.Set(ctxkey.OriginalModel, modelName) // for retry
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", channel.Key))
	channelID := registry.IDByLegacyType(channel.Type)
	c.Set(ctxkey.ChannelID, channelID)
	baseURL := channel.GetBaseURL()
	if baseURL == "" {
		baseURL = registry.GetDefaultBaseURL(channelID)
	}
	c.Set(ctxkey.BaseURL, baseURL)
	cfg, _ := channel.LoadConfig()
	// this is for backward compatibility
	if channel.Other != nil {
		switch registry.IDByLegacyType(channel.Type) {
		case "azure":
			if cfg.APIVersion == "" {
				cfg.APIVersion = *channel.Other
			}
		case "xunfei":
			if cfg.APIVersion == "" {
				cfg.APIVersion = *channel.Other
			}
		case "gemini":
			if cfg.APIVersion == "" {
				cfg.APIVersion = *channel.Other
			}
		case "aiproxylibrary":
			if cfg.LibraryID == "" {
				cfg.LibraryID = *channel.Other
			}
		case "ali":
			if cfg.Plugin == "" {
				cfg.Plugin = *channel.Other
			}
		}
	}
	c.Set(ctxkey.Config, cfg)
}