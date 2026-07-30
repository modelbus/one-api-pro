package registry

import (
	"sync"

	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
)

type ChannelMeta struct {
	ID             string
	Name           string
	DefaultBaseURL string
	LegacyType     int
}

type AdaptorFactory func() adaptor.Adaptor

type channelEntry struct {
	meta    ChannelMeta
	factory AdaptorFactory
}

var (
	mu               sync.RWMutex
	channelsByID     = make(map[string]*channelEntry)
	channelsByLegacy = make(map[int]*channelEntry)
	channelIDs       []string
)

func Register(meta ChannelMeta, factory AdaptorFactory) {
	mu.Lock()
	defer mu.Unlock()

	entry := &channelEntry{meta: meta, factory: factory}
	channelsByID[meta.ID] = entry
	if meta.LegacyType != 0 {
		channelsByLegacy[meta.LegacyType] = entry
	}
	channelIDs = append(channelIDs, meta.ID)
}

func GetAdaptor(id string) adaptor.Adaptor {
	mu.RLock()
	defer mu.RUnlock()

	if entry, ok := channelsByID[id]; ok {
		return entry.factory()
	}
	return nil
}

func GetAdaptorByLegacyType(typ int) adaptor.Adaptor {
	mu.RLock()
	defer mu.RUnlock()

	if entry, ok := channelsByLegacy[typ]; ok {
		return entry.factory()
	}
	return nil
}

func GetChannelMeta(id string) (ChannelMeta, bool) {
	mu.RLock()
	defer mu.RUnlock()

	if entry, ok := channelsByID[id]; ok {
		return entry.meta, true
	}
	return ChannelMeta{}, false
}

func GetChannelMetaByLegacyType(typ int) (ChannelMeta, bool) {
	mu.RLock()
	defer mu.RUnlock()

	if entry, ok := channelsByLegacy[typ]; ok {
		return entry.meta, true
	}
	return ChannelMeta{}, false
}

func GetDefaultBaseURL(id string) string {
	mu.RLock()
	defer mu.RUnlock()

	if entry, ok := channelsByID[id]; ok {
		return entry.meta.DefaultBaseURL
	}
	return ""
}

func AllChannelIDs() []string {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]string, len(channelIDs))
	copy(result, channelIDs)
	return result
}

func AllChannelMetas() []ChannelMeta {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]ChannelMeta, 0, len(channelIDs))
	for _, id := range channelIDs {
		if entry, ok := channelsByID[id]; ok {
			result = append(result, entry.meta)
		}
	}
	return result
}

func IDByLegacyType(typ int) string {
	mu.RLock()
	defer mu.RUnlock()

	if entry, ok := channelsByLegacy[typ]; ok {
		return entry.meta.ID
	}
	return ""
}
