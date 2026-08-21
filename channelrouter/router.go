package channelrouter

import (
	"context"
	"fmt"
	"time"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
)

type ChannelRouter struct {
	filters      []ChannelFilter
	selector     ChannelSelector
	Cooldown     *CooldownManager
	Concurrency  *ConcurrencyTracker
	Sticky       *StickySessionStore
	RPM          *RPMTracker
}

var DefaultRouter *ChannelRouter

func NewChannelRouter() *ChannelRouter {
	r := &ChannelRouter{
		Cooldown:    NewCooldownManager(),
		Concurrency: NewConcurrencyTracker(),
		Sticky:      NewStickySessionStore(),
		RPM:         NewRPMTracker(),
		selector:    &PriorityRandomSelector{},
	}
	r.filters = []ChannelFilter{
		&StatusFilter{},
		&FallbackFilter{},
		&CooldownFilter{cooldown: r.Cooldown},
		&ConcurrencyFilter{concurrency: r.Concurrency},
		&RPMFilter{rpmTracker: r.RPM},
		&StickySessionFilter{sticky: r.Sticky},
		&PriorityFilter{},
	}
	return r
}

func InitRouter() {
	DefaultRouter = NewChannelRouter()
	DefaultRouter.Cooldown.StartCleanupLoop(30 * time.Second)
	if config.ChannelStickySessionEnabled {
		DefaultRouter.Sticky.LoadFromLogDB()
	}
	logger.SysLog("channel router initialized")
}

func (r *ChannelRouter) Route(ctx context.Context, req *RouteRequest, candidates []*model.Channel) (*model.Channel, error) {
	filtered := candidates
	for _, f := range r.filters {
		filtered = f.Filter(ctx, filtered, req)
		if len(filtered) == 0 {
			break
		}
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no available channel for group=%s model=%s", req.Group, req.Model)
	}
	return r.selector.Select(ctx, filtered, req)
}

func (r *ChannelRouter) SetCooldown(channelId int, seconds int, reason string, statusCode int) {
	r.Cooldown.SetCooldown(channelId, seconds, reason, statusCode)
}

func (r *ChannelRouter) IsInCooldown(channelId int) bool {
	return r.Cooldown.IsInCooldown(channelId)
}

func (r *ChannelRouter) TryAcquireConcurrency(channelId int, maxConcurrency int) bool {
	return r.Concurrency.TryAcquire(channelId, maxConcurrency)
}

func (r *ChannelRouter) ReleaseConcurrency(channelId int) {
	r.Concurrency.Release(channelId)
}

func (r *ChannelRouter) SetStickySession(sessionKey string, channelId int) {
	if sessionKey == "" {
		return
	}
	r.Sticky.Set(sessionKey, channelId)
}

func (r *ChannelRouter) GetStickySession(sessionKey string) int {
	if sessionKey == "" {
		return 0
	}
	return r.Sticky.Get(sessionKey)
}

func (r *ChannelRouter) IncrementRPM(channelId int) {
	r.RPM.Increment(channelId)
}