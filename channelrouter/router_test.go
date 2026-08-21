package channelrouter

import (
	"context"
	"testing"
	"time"

	"github.com/Leon-PanPan/one-api-pro/model"
)

func makeChannel(id, status, priority int, maxConc int) *model.Channel {
	p := int64(priority)
	mc := maxConc
	ch := &model.Channel{
		Id:             id,
		Status:         status,
		Models:         "gpt-4",
		Group:          "default",
		Priority:       &p,
		MaxConcurrency: &mc,
	}
	return ch
}

func TestStatusFilter(t *testing.T) {
	f := &StatusFilter{}
	candidates := []*model.Channel{
		makeChannel(1, model.ChannelStatusEnabled, 0, 0),
		makeChannel(2, model.ChannelStatusAutoDisabled, 0, 0),
		makeChannel(3, model.ChannelStatusManuallyDisabled, 0, 0),
		makeChannel(4, model.ChannelStatusEnabled, 0, 0),
	}
	result := f.Filter(context.Background(), candidates, &RouteRequest{})
	if len(result) != 2 {
		t.Errorf("StatusFilter: got %d, want 2", len(result))
	}
}

func TestCooldownFilter(t *testing.T) {
	cm := NewCooldownManager()
	f := &CooldownFilter{cooldown: cm}
	candidates := []*model.Channel{
		makeChannel(1, model.ChannelStatusEnabled, 0, 0),
		makeChannel(2, model.ChannelStatusEnabled, 0, 0),
	}
	cm.SetCooldown(1, 60, "test", 429)
	result := f.Filter(context.Background(), candidates, &RouteRequest{})
	if len(result) != 1 || result[0].Id != 2 {
		t.Errorf("CooldownFilter: got %v", result)
	}
}

func TestConcurrencyFilter(t *testing.T) {
	ct := NewConcurrencyTracker()
	f := &ConcurrencyFilter{concurrency: ct}
	candidates := []*model.Channel{
		makeChannel(1, model.ChannelStatusEnabled, 0, 2),
		makeChannel(2, model.ChannelStatusEnabled, 0, 0),
	}
	ct.TryAcquire(1, 2)
	ct.TryAcquire(1, 2)
	result := f.Filter(context.Background(), candidates, &RouteRequest{})
	if len(result) != 1 || result[0].Id != 2 {
		t.Errorf("ConcurrencyFilter: got %v", result)
	}
	ct.Release(1)
	result = f.Filter(context.Background(), candidates, &RouteRequest{})
	if len(result) != 2 {
		t.Errorf("ConcurrencyFilter after release: got %d, want 2", len(result))
	}
}

func TestStickySessionFilter(t *testing.T) {
	ss := NewStickySessionStore()
	f := &StickySessionFilter{sticky: ss}
	candidates := []*model.Channel{
		makeChannel(1, model.ChannelStatusEnabled, 0, 0),
		makeChannel(2, model.ChannelStatusEnabled, 0, 0),
	}
	ss.Set("user1:gpt-4", 1)
	result := f.Filter(context.Background(), candidates, &RouteRequest{SessionKey: "user1:gpt-4"})
	if len(result) != 1 || result[0].Id != 1 {
		t.Errorf("StickySessionFilter: got %v", result)
	}
	result = f.Filter(context.Background(), candidates, &RouteRequest{SessionKey: ""})
	if len(result) != 2 {
		t.Errorf("StickySessionFilter no key: got %d, want 2", len(result))
	}
}

func TestCooldownManager(t *testing.T) {
	cm := NewCooldownManager()
	if cm.IsInCooldown(1) {
		t.Error("should not be in cooldown initially")
	}
	cm.SetCooldown(1, 60, "test", 429)
	if !cm.IsInCooldown(1) {
		t.Error("should be in cooldown after set")
	}
	remaining := cm.GetCooldownRemaining(1)
	if remaining <= 0 {
		t.Error("remaining should be positive")
	}
	val, ok := cm.store.Load(1)
	if !ok {
		t.Fatal("cooldown entry should exist")
	}
	entry := val.(*CooldownEntry)
	entry.UntilTime = entry.UntilTime.Add(-61 * time.Second)
	cm.CleanExpired()
	if cm.IsInCooldown(1) {
		t.Error("should not be in cooldown after manual expiry")
	}
}

func TestConcurrencyTracker(t *testing.T) {
	ct := NewConcurrencyTracker()
	if !ct.TryAcquire(1, 2) {
		t.Error("first acquire should succeed")
	}
	if !ct.TryAcquire(1, 2) {
		t.Error("second acquire should succeed")
	}
	if ct.TryAcquire(1, 2) {
		t.Error("third acquire should fail (at capacity)")
	}
	ct.Release(1)
	if !ct.TryAcquire(1, 2) {
		t.Error("acquire after release should succeed")
	}
	ct.Release(1)
	ct.Release(1)
	if ct.GetActiveCount(1) != 0 {
		t.Errorf("count should be 0, got %d", ct.GetActiveCount(1))
	}
}

func TestPriorityRandomSelector(t *testing.T) {
	s := &PriorityRandomSelector{}
	candidates := []*model.Channel{
		makeChannel(1, model.ChannelStatusEnabled, 10, 0),
		makeChannel(2, model.ChannelStatusEnabled, 10, 0),
		makeChannel(3, model.ChannelStatusEnabled, 5, 0),
	}
	ch, err := s.Select(context.Background(), candidates, &RouteRequest{})
	if err != nil {
		t.Errorf("Select failed: %v", err)
	}
	if ch.Id != 1 && ch.Id != 2 {
		t.Errorf("Select: got channel %d, expected 1 or 2 (highest priority)", ch.Id)
	}
}

func TestMakeSessionKey(t *testing.T) {
	key := MakeSessionKey(123, "gpt-4")
	if key != "123:gpt-4" {
		t.Errorf("MakeSessionKey = %s, want 123:gpt-4", key)
	}
}

func makeFallbackChannel(id, status, priority int, fallback bool, fallbackPriority int64) *model.Channel {
	p := int64(priority)
	isFb := fallback
	fp := fallbackPriority
	ch := &model.Channel{
		Id:               id,
		Status:           status,
		Models:           "gpt-4",
		Group:            "default",
		Priority:         &p,
		IsFallback:       &isFb,
		FallbackPriority: &fp,
	}
	return ch
}

func TestFallbackFilter(t *testing.T) {
	f := &FallbackFilter{}
	candidates := []*model.Channel{
		makeFallbackChannel(1, model.ChannelStatusEnabled, 0, false, 0),
		makeFallbackChannel(2, model.ChannelStatusEnabled, 0, true, 10),
		makeFallbackChannel(3, model.ChannelStatusEnabled, 0, false, 0),
		makeFallbackChannel(4, model.ChannelStatusEnabled, 0, true, 20),
	}
	result := f.Filter(context.Background(), candidates, &RouteRequest{})
	if len(result) != 2 {
		t.Fatalf("FallbackFilter: got %d, want 2 (non-fallback channels)", len(result))
	}
	for _, ch := range result {
		if ch.GetIsFallback() {
			t.Errorf("FallbackFilter leaked fallback channel #%d", ch.Id)
		}
	}
}

func TestFallbackFilter_NoFallback(t *testing.T) {
	f := &FallbackFilter{}
	candidates := []*model.Channel{
		makeFallbackChannel(1, model.ChannelStatusEnabled, 0, false, 0),
		makeFallbackChannel(2, model.ChannelStatusEnabled, 0, false, 0),
	}
	result := f.Filter(context.Background(), candidates, &RouteRequest{})
	if len(result) != 2 {
		t.Errorf("FallbackFilter (no fallback): got %d, want 2", len(result))
	}
}

func TestFallbackFilter_AllFallback(t *testing.T) {
	f := &FallbackFilter{}
	candidates := []*model.Channel{
		makeFallbackChannel(1, model.ChannelStatusEnabled, 0, true, 0),
		makeFallbackChannel(2, model.ChannelStatusEnabled, 0, true, 0),
	}
	result := f.Filter(context.Background(), candidates, &RouteRequest{})
	if len(result) != 0 {
		t.Errorf("FallbackFilter (all fallback): got %d, want 0", len(result))
	}
}

func TestRouteExcludesFallbackChannels(t *testing.T) {
	router := NewChannelRouter()
	candidates := []*model.Channel{
		makeFallbackChannel(10, model.ChannelStatusEnabled, 100, false, 0), // highest priority, normal
		makeFallbackChannel(20, model.ChannelStatusEnabled, 100, true, 50),  // highest priority, fallback
		makeFallbackChannel(30, model.ChannelStatusEnabled, 50, false, 0),   // lower priority, normal
	}
	for i := 0; i < 50; i++ {
		ch, err := router.Route(context.Background(), &RouteRequest{
			Group: "default",
			Model: "gpt-4",
		}, candidates)
		if err != nil {
			t.Fatalf("Route failed: %v", err)
		}
		if ch.GetIsFallback() {
			t.Errorf("Route returned a fallback channel #%d on iteration %d", ch.Id, i)
		}
	}
}

func TestChannelGetters(t *testing.T) {
	t.Run("IsFallback nil pointer", func(t *testing.T) {
		ch := &model.Channel{}
		if ch.GetIsFallback() {
			t.Error("nil IsFallback should return false")
		}
	})
	t.Run("IsFallback true", func(t *testing.T) {
		isFb := true
		ch := &model.Channel{IsFallback: &isFb}
		if !ch.GetIsFallback() {
			t.Error("IsFallback=true should return true")
		}
	})
	t.Run("IsFallback false (explicit)", func(t *testing.T) {
		isFb := false
		ch := &model.Channel{IsFallback: &isFb}
		if ch.GetIsFallback() {
			t.Error("IsFallback=false should return false")
		}
	})
	t.Run("FallbackPriority nil pointer", func(t *testing.T) {
		ch := &model.Channel{}
		if ch.GetFallbackPriority() != 0 {
			t.Errorf("nil FallbackPriority should return 0, got %d", ch.GetFallbackPriority())
		}
	})
	t.Run("FallbackPriority set", func(t *testing.T) {
		fp := int64(42)
		ch := &model.Channel{FallbackPriority: &fp}
		if ch.GetFallbackPriority() != 42 {
			t.Errorf("FallbackPriority=42, got %d", ch.GetFallbackPriority())
		}
	})
}