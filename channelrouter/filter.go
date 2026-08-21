package channelrouter

import (
	"context"

	"github.com/modelbus/one-api-pro/model"
)

type StatusFilter struct{}

func (f *StatusFilter) Name() string {
	return "status"
}

func (f *StatusFilter) Filter(_ context.Context, candidates []*model.Channel, _ *RouteRequest) []*model.Channel {
	result := make([]*model.Channel, 0, len(candidates))
	for _, ch := range candidates {
		if ch.Status == model.ChannelStatusEnabled {
			result = append(result, ch)
		}
	}
	return result
}

type CooldownFilter struct {
	cooldown *CooldownManager
}

func (f *CooldownFilter) Name() string {
	return "cooldown"
}

func (f *CooldownFilter) Filter(_ context.Context, candidates []*model.Channel, _ *RouteRequest) []*model.Channel {
	result := make([]*model.Channel, 0, len(candidates))
	for _, ch := range candidates {
		if !f.cooldown.IsInCooldown(ch.Id) {
			result = append(result, ch)
		}
	}
	return result
}

type ConcurrencyFilter struct {
	concurrency *ConcurrencyTracker
}

func (f *ConcurrencyFilter) Name() string {
	return "concurrency"
}

func (f *ConcurrencyFilter) Filter(_ context.Context, candidates []*model.Channel, _ *RouteRequest) []*model.Channel {
	result := make([]*model.Channel, 0, len(candidates))
	for _, ch := range candidates {
		maxConc := ch.GetMaxConcurrency()
		if maxConc <= 0 {
			result = append(result, ch)
			continue
		}
		if !f.concurrency.IsAtCapacity(ch.Id, maxConc) {
			result = append(result, ch)
		}
	}
	return result
}

type StickySessionFilter struct {
	sticky *StickySessionStore
}

func (f *StickySessionFilter) Name() string {
	return "sticky_session"
}

func (f *StickySessionFilter) Filter(_ context.Context, candidates []*model.Channel, req *RouteRequest) []*model.Channel {
	if req.SessionKey == "" {
		return candidates
	}
	prevChannelId := f.sticky.Get(req.SessionKey)
	if prevChannelId <= 0 {
		return candidates
	}
	for _, ch := range candidates {
		if ch.Id == prevChannelId && ch.Status == model.ChannelStatusEnabled {
			return []*model.Channel{ch}
		}
	}
	return candidates
}

type PriorityFilter struct{}

func (f *PriorityFilter) Name() string {
	return "priority"
}

func (f *PriorityFilter) Filter(_ context.Context, candidates []*model.Channel, req *RouteRequest) []*model.Channel {
	if len(candidates) == 0 {
		return candidates
	}
	if !req.IgnoreFirstPriority {
		return candidates
	}
	highestPriority := candidates[0].GetPriority()
	for i, ch := range candidates {
		if ch.GetPriority() != highestPriority {
			return candidates[i:]
		}
	}
	if len(candidates) > 1 {
		return candidates[1:]
	}
	return candidates
}

type RPMFilter struct {
	rpmTracker *RPMTracker
}

func (f *RPMFilter) Name() string {
	return "rpm"
}

func (f *RPMFilter) Filter(_ context.Context, candidates []*model.Channel, _ *RouteRequest) []*model.Channel {
	result := make([]*model.Channel, 0, len(candidates))
	for _, ch := range candidates {
		maxRPM := ch.GetRPM()
		if maxRPM <= 0 {
			result = append(result, ch)
			continue
		}
		if f.rpmTracker.CurrentRPM(ch.Id) < maxRPM {
			result = append(result, ch)
		}
	}
	return result
}

// FallbackFilter excludes channels that are reserved for fallback only
// (channels with is_fallback=true). These channels are not used for normal
// requests and are only picked by the dedicated fallback path after all
// normal channels for the requested model are exhausted.
type FallbackFilter struct{}

func (f *FallbackFilter) Name() string {
	return "fallback"
}

func (f *FallbackFilter) Filter(_ context.Context, candidates []*model.Channel, _ *RouteRequest) []*model.Channel {
	result := make([]*model.Channel, 0, len(candidates))
	for _, ch := range candidates {
		if ch.GetIsFallback() {
			continue
		}
		result = append(result, ch)
	}
	return result
}