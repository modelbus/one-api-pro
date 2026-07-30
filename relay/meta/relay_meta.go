package meta

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Leon-PanPan/one-api-pro/common/ctxkey"
	"github.com/Leon-PanPan/one-api-pro/model"
	"github.com/Leon-PanPan/one-api-pro/relay/relaymode"
)

type Meta struct {
	Mode         int
	ChannelType  int
	ChannelID    string
	ChannelId    int
	TokenId      int
	TokenName    string
	UserId       int
	Group        string
	ModelMapping map[string]string
	// BaseURL is the proxy url set in the channel config
	BaseURL  string
	APIKey   string
	Config   model.ChannelConfig
	IsStream bool
	// OriginModelName is the model name from the raw user request
	OriginModelName string
	// ActualModelName is the model name after mapping
	ActualModelName    string
	RequestURLPath     string
	PromptTokens       int // only for DoResponse
	ForcedSystemPrompt string
	StartTime          time.Time
	PlanId              int    // UserPlan.Id if using subscription, 0 if using quota
	BillingType         string // "request" or "token"
	DefaultModel        string // plan's default_model for forwarding
	SessionKey          string
}

func GetByContext(c *gin.Context) *Meta {
	meta := Meta{
		Mode:               relaymode.GetByPath(c.Request.URL.Path),
		ChannelType:        c.GetInt(ctxkey.Channel),
		ChannelId:          c.GetInt(ctxkey.ChannelId),
		TokenId:            c.GetInt(ctxkey.TokenId),
		TokenName:          c.GetString(ctxkey.TokenName),
		UserId:             c.GetInt(ctxkey.Id),
		Group:              c.GetString(ctxkey.Group),
		ModelMapping:       c.GetStringMapString(ctxkey.ModelMapping),
		OriginModelName:    c.GetString(ctxkey.RequestModel),
		BaseURL:            c.GetString(ctxkey.BaseURL),
		APIKey:             strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer "),
		RequestURLPath:     c.Request.URL.String(),
		ForcedSystemPrompt: c.GetString(ctxkey.SystemPrompt),
		StartTime:          time.Now(),
		PlanId:            c.GetInt(ctxkey.PlanId),
		BillingType:       c.GetString(ctxkey.BillingType),
		DefaultModel:      c.GetString(ctxkey.DefaultModel),
		SessionKey:        c.GetString(ctxkey.SessionKey),
	}
	cfg, ok := c.Get(ctxkey.Config)
	if ok {
		meta.Config = cfg.(model.ChannelConfig)
	}
	meta.ChannelID = c.GetString(ctxkey.ChannelID)
	return &meta
}
