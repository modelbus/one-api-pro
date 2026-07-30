package channelrouter

import (
	"context"

	"github.com/Leon-PanPan/one-api-pro/model"
)

type RouteRequest struct {
	Group               string
	Model               string
	UserId              int
	IgnoreFirstPriority bool
	SessionKey          string
}

type ChannelFilter interface {
	Name() string
	Filter(ctx context.Context, candidates []*model.Channel, req *RouteRequest) []*model.Channel
}

type ChannelSelector interface {
	Select(ctx context.Context, candidates []*model.Channel, req *RouteRequest) (*model.Channel, error)
}