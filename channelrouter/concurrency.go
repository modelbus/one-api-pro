package channelrouter

import (
	"sync"
	"sync/atomic"

	"github.com/modelbus/one-api-pro/cluster"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/model"
)

type ConcurrencyTracker struct {
	activeRequests sync.Map
}

func NewConcurrencyTracker() *ConcurrencyTracker {
	return &ConcurrencyTracker{}
}

func (t *ConcurrencyTracker) TryAcquire(channelId int, maxConcurrency int) bool {
	if maxConcurrency <= 0 {
		return true
	}
	if cluster.Enabled {
		var total int64
		model.DB.Model(&model.ChannelCounter{}).
			Where("channel_id = ?", channelId).
			Select("SUM(concurrency)").Scan(&total)
		if int(total) >= maxConcurrency {
			return false
		}
		now := helper.GetTimestamp()
		model.DB.Exec(`
			INSERT INTO channel_counters (channel_id, node_id, concurrency, rpm_count, rpm_minute, updated_at)
			VALUES (?, ?, 1, 0, 0, ?)
			ON DUPLICATE KEY UPDATE concurrency = concurrency + 1, updated_at = ?
		`, channelId, cluster.NodeID, now, now)
		return true
	}
	counter := t.getCounter(channelId)
	current := counter.Load()
	if int(current) >= maxConcurrency {
		return false
	}
	if !counter.CompareAndSwap(current, current+1) {
		current = counter.Load()
		if int(current) >= maxConcurrency {
			return false
		}
		if !counter.CompareAndSwap(current, current+1) {
			return false
		}
	}
	return true
}

func (t *ConcurrencyTracker) Release(channelId int) {
	if cluster.Enabled {
		now := helper.GetTimestamp()
		model.DB.Exec(`
			UPDATE channel_counters SET concurrency = GREATEST(concurrency - 1, 0), updated_at = ?
			WHERE channel_id = ? AND node_id = ?
		`, now, channelId, cluster.NodeID)
		return
	}
	counter := t.getCounter(channelId)
	newVal := counter.Add(-1)
	if newVal < 0 {
		counter.Store(0)
	}
}

func (t *ConcurrencyTracker) GetActiveCount(channelId int) int64 {
	if cluster.Enabled {
		var total int64
		model.DB.Model(&model.ChannelCounter{}).
			Where("channel_id = ?", channelId).
			Select("SUM(concurrency)").Scan(&total)
		return total
	}
	counter := t.getCounter(channelId)
	return counter.Load()
}

func (t *ConcurrencyTracker) IsAtCapacity(channelId int, maxConcurrency int) bool {
	if maxConcurrency <= 0 {
		return false
	}
	return int(t.GetActiveCount(channelId)) >= maxConcurrency
}

func (t *ConcurrencyTracker) getCounter(channelId int) *atomic.Int64 {
	val, ok := t.activeRequests.Load(channelId)
	if !ok {
		newCounter := &atomic.Int64{}
		actual, loaded := t.activeRequests.LoadOrStore(channelId, newCounter)
		if loaded {
			return actual.(*atomic.Int64)
		}
		return newCounter
	}
	return val.(*atomic.Int64)
}