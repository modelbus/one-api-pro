package channelrouter

import (
	"fmt"
	"sync"
	"time"

	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
)

type StickySessionStore struct {
	mu    sync.RWMutex
	store map[string]int
}

func NewStickySessionStore() *StickySessionStore {
	return &StickySessionStore{
		store: make(map[string]int),
	}
}

func (s *StickySessionStore) Get(sessionKey string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store[sessionKey]
}

func (s *StickySessionStore) Set(sessionKey string, channelId int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[sessionKey] = channelId
}

func (s *StickySessionStore) Delete(sessionKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, sessionKey)
}

func MakeSessionKey(userId int, modelName string) string {
	return fmt.Sprintf("%d:%s", userId, modelName)
}

func (s *StickySessionStore) LoadFromLogDB() {
	var results []struct {
		SessionKey string
		ChannelId  int
	}
	cutoff := time.Now().Add(-24 * time.Hour).Unix()
	err := model.LOG_DB.Raw(`
		SELECT session_key, channel_id
		FROM logs
		WHERE type = 2
		  AND session_key != ''
		  AND created_at > ?
		ORDER BY created_at DESC
	`, cutoff).Scan(&results).Error
	if err != nil {
		logger.SysError("failed to load sticky sessions from log DB: " + err.Error())
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range results {
		if _, exists := s.store[r.SessionKey]; !exists {
			s.store[r.SessionKey] = r.ChannelId
		}
	}
	logger.SysLog(fmt.Sprintf("loaded %d sticky sessions from log DB", len(s.store)))
}