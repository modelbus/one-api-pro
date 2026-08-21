package channelrouter

import (
	"fmt"
	"sync"
	"time"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/logger"
)

type CooldownEntry struct {
	ChannelId    int
	UntilTime    time.Time
	Reason       string
	OriginalCode int
}

type CooldownManager struct {
	store sync.Map
}

func NewCooldownManager() *CooldownManager {
	return &CooldownManager{}
}

func (m *CooldownManager) SetCooldown(channelId int, seconds int, reason string, statusCode int) {
	if seconds <= 0 {
		return
	}
	if seconds > config.ChannelMaxCooldownSeconds {
		seconds = config.ChannelMaxCooldownSeconds
	}
	entry := &CooldownEntry{
		ChannelId:    channelId,
		UntilTime:    time.Now().Add(time.Duration(seconds) * time.Second),
		Reason:       reason,
		OriginalCode: statusCode,
	}
	m.store.Store(channelId, entry)
	logger.SysLog(fmt.Sprintf("channel #%d entered cooldown for %ds (reason: %s, status: %d)", channelId, seconds, reason, statusCode))
}

func (m *CooldownManager) IsInCooldown(channelId int) bool {
	val, ok := m.store.Load(channelId)
	if !ok {
		return false
	}
	entry := val.(*CooldownEntry)
	if time.Now().Before(entry.UntilTime) {
		return true
	}
	m.store.Delete(channelId)
	return false
}

func (m *CooldownManager) GetCooldownRemaining(channelId int) time.Duration {
	val, ok := m.store.Load(channelId)
	if !ok {
		return 0
	}
	entry := val.(*CooldownEntry)
	remaining := time.Until(entry.UntilTime)
	if remaining <= 0 {
		m.store.Delete(channelId)
		return 0
	}
	return remaining
}

func (m *CooldownManager) CleanExpired() {
	now := time.Now()
	m.store.Range(func(key, val interface{}) bool {
		entry := val.(*CooldownEntry)
		if now.After(entry.UntilTime) {
			m.store.Delete(key)
		}
		return true
	})
}

func (m *CooldownManager) StartCleanupLoop(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			m.CleanExpired()
		}
	}()
}