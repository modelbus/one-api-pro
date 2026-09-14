package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/common/random"
)

var (
	TokenCacheSeconds         = config.SyncFrequency
	UserId2GroupCacheSeconds  = config.SyncFrequency
	UserId2QuotaCacheSeconds  = config.SyncFrequency
	UserId2StatusCacheSeconds = config.SyncFrequency
	GroupModelsCacheSeconds   = config.SyncFrequency
)

func CacheGetTokenByKey(key string) (*Token, error) {
	keyCol := "`key`"
	if common.UsingPostgreSQL {
		keyCol = `"key"`
	}
	var token Token
	if !common.RedisEnabled {
		err := DB.Where(keyCol+" = ?", key).First(&token).Error
		return &token, err
	}
	tokenObjectString, err := common.RedisGet(fmt.Sprintf("token:%s", key))
	if err != nil {
		err := DB.Where(keyCol+" = ?", key).First(&token).Error
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(token)
		if err != nil {
			return nil, err
		}
		err = common.RedisSet(fmt.Sprintf("token:%s", key), string(jsonBytes), time.Duration(TokenCacheSeconds)*time.Second)
		if err != nil {
			logger.SysError("Redis set token error: " + err.Error())
		}
		return &token, nil
	}
	err = json.Unmarshal([]byte(tokenObjectString), &token)
	return &token, err
}

func CacheGetUserGroup(id int) (group string, err error) {
	if !common.RedisEnabled {
		return GetUserGroup(id)
	}
	group, err = common.RedisGet(fmt.Sprintf("user_group:%d", id))
	if err != nil {
		group, err = GetUserGroup(id)
		if err != nil {
			return "", err
		}
		err = common.RedisSet(fmt.Sprintf("user_group:%d", id), group, time.Duration(UserId2GroupCacheSeconds)*time.Second)
		if err != nil {
			logger.SysError("Redis set user group error: " + err.Error())
		}
	}
	return group, err
}

func fetchAndUpdateUserQuota(ctx context.Context, id int) (quota int64, err error) {
	quota, err = GetUserQuota(id)
	if err != nil {
		return 0, err
	}
	err = common.RedisSet(fmt.Sprintf("user_quota:%d", id), fmt.Sprintf("%d", quota), time.Duration(UserId2QuotaCacheSeconds)*time.Second)
	if err != nil {
		logger.Error(ctx, "Redis set user quota error: "+err.Error())
	}
	return
}

// userQuotaLowWaterMark is the cache freshness threshold for user_quota:<id>.
// When the cached value drops to or below this mark, the next read will refresh
// from DB to avoid serving a stale value that no longer reflects users.quota.
//
// Why 50_000 (instead of config.PreConsumedQuota = 500):
//   - Single pre-consume for high-price models (e.g. minimax-m3 ¥8.40/1M tokens)
//     can be ~80,000 quota, far above 500.
//   - With the old threshold, the cache could drift into the dead zone
//     [501, 80_000) — large enough to look "has quota" but too small to cover
//     the next pre-consume — and stay there indefinitely until it dropped
//     further to ≤500, producing spurious 403 "insufficient_user_quota".
//   - 50_000 keeps the dead zone narrow (~30K wide vs ~80K before) and well
//     below typical single-request pre-consume amounts.
//
// userQuotaLowWaterMark 是 user_quota:<id> 的缓存新鲜度阈值。
// 缓存值跌至或低于该阈值时，下次读取会回源 DB 刷新，避免返回已偏离
// users.quota 的陈旧值。
//
// 设为 50_000（而非 config.PreConsumedQuota = 500）的原因：
//   - 高单价模型（如 minimax-m3 ¥8.40/百万 tokens）单次预扣可达 ~80,000 quota，
//     远超 500。
//   - 旧阈值下，缓存可能漂移到 [501, 80_000) 死区——看似还有钱但已不够
//     下一次预扣——并长期停留，直至进一步跌至 ≤500 才触发刷新，造成
//     误报 403 insufficient_user_quota。
//   - 50_000 让死区收窄到 ~30K（原来 ~80K），且远低于典型单次预扣。
const userQuotaLowWaterMark int64 = 50_000

func CacheGetUserQuota(ctx context.Context, id int) (quota int64, err error) {
	if !common.RedisEnabled {
		return GetUserQuota(id)
	}
	quotaString, err := common.RedisGet(fmt.Sprintf("user_quota:%d", id))
	if err != nil {
		return fetchAndUpdateUserQuota(ctx, id)
	}
	quota, err = strconv.ParseInt(quotaString, 10, 64)
	if err != nil {
		return 0, nil
	}
	if quota <= userQuotaLowWaterMark { // when user's cached quota is at or below the low-water mark, we need to fetch from db
		logger.Infof(ctx, "user %d's cached quota is too low: %d, refreshing from db", quota, id)
		return fetchAndUpdateUserQuota(ctx, id)
	}
	return quota, nil
}

func CacheUpdateUserQuota(ctx context.Context, id int) error {
	if !common.RedisEnabled {
		return nil
	}
	quota, err := CacheGetUserQuota(ctx, id)
	if err != nil {
		return err
	}
	err = common.RedisSet(fmt.Sprintf("user_quota:%d", id), fmt.Sprintf("%d", quota), time.Duration(UserId2QuotaCacheSeconds)*time.Second)
	return err
}

func CacheDecreaseUserQuota(id int, quota int64) error {
	if !common.RedisEnabled {
		return nil
	}
	err := common.RedisDecrease(fmt.Sprintf("user_quota:%d", id), int64(quota))
	return err
}

// CacheIncreaseUserQuota rolls back a previous CacheDecreaseUserQuota when the
// matching DB-side write fails (e.g. PreConsumeTokenQuota rejected the token).
// Without this rollback, the Redis cache would drift further from users.quota
// every time the post-decrease DB write fails, eventually landing in the dead
// zone [lowWaterMark, preConsumedQuota) and causing spurious 403.
//
// CacheIncreaseUserQuota 在对应的 DB 写入失败（如 PreConsumeTokenQuota 拒绝该
// token）时，回滚此前 CacheDecreaseUserQuota 的扣减。若不回滚，每次 DB 写入
// 失败都会让 Redis 缓存进一步偏离 users.quota，最终落入死区
// [lowWaterMark, preConsumedQuota)，导致误报 403。
func CacheIncreaseUserQuota(id int, quota int64) error {
	if !common.RedisEnabled {
		return nil
	}
	err := common.RedisIncrease(fmt.Sprintf("user_quota:%d", id), int64(quota))
	return err
}

func CacheIsUserEnabled(userId int) (bool, error) {
	if !common.RedisEnabled {
		return IsUserEnabled(userId)
	}
	enabled, err := common.RedisGet(fmt.Sprintf("user_enabled:%d", userId))
	if err == nil {
		return enabled == "1", nil
	}

	userEnabled, err := IsUserEnabled(userId)
	if err != nil {
		return false, err
	}
	enabled = "0"
	if userEnabled {
		enabled = "1"
	}
	err = common.RedisSet(fmt.Sprintf("user_enabled:%d", userId), enabled, time.Duration(UserId2StatusCacheSeconds)*time.Second)
	if err != nil {
		logger.SysError("Redis set user enabled error: " + err.Error())
	}
	return userEnabled, err
}

func CacheGetGroupModels(ctx context.Context, group string) ([]string, error) {
	if !common.RedisEnabled {
		return GetGroupModels(ctx, group)
	}
	modelsStr, err := common.RedisGet(fmt.Sprintf("group_models:%s", group))
	if err == nil {
		return strings.Split(modelsStr, ","), nil
	}
	models, err := GetGroupModels(ctx, group)
	if err != nil {
		return nil, err
	}
	err = common.RedisSet(fmt.Sprintf("group_models:%s", group), strings.Join(models, ","), time.Duration(GroupModelsCacheSeconds)*time.Second)
	if err != nil {
		logger.SysError("Redis set group models error: " + err.Error())
	}
	return models, nil
}

var group2model2channels map[string]map[string][]*Channel
var channelSyncLock sync.RWMutex
var channelId2channel map[int]*Channel

func InitChannelCache() {
	newChannelId2channel := make(map[int]*Channel)
	var channels []*Channel
	DB.Where("status = ?", ChannelStatusEnabled).Find(&channels)
	for _, channel := range channels {
		newChannelId2channel[channel.Id] = channel
	}
	var abilities []*Ability
	DB.Find(&abilities)
	groups := make(map[string]bool)
	for _, ability := range abilities {
		groups[ability.Group] = true
	}
	newGroup2model2channels := make(map[string]map[string][]*Channel)
	for group := range groups {
		newGroup2model2channels[group] = make(map[string][]*Channel)
	}
	for _, channel := range channels {
		groups := strings.Split(channel.Group, ",")
		for _, group := range groups {
			models := strings.Split(channel.Models, ",")
			for _, model := range models {
				if _, ok := newGroup2model2channels[group][model]; !ok {
					newGroup2model2channels[group][model] = make([]*Channel, 0)
				}
				newGroup2model2channels[group][model] = append(newGroup2model2channels[group][model], channel)
			}
		}
	}

	// sort by priority
	for group, model2channels := range newGroup2model2channels {
		for model, channels := range model2channels {
			sort.Slice(channels, func(i, j int) bool {
				return channels[i].GetPriority() > channels[j].GetPriority()
			})
			newGroup2model2channels[group][model] = channels
		}
	}

	channelSyncLock.Lock()
	group2model2channels = newGroup2model2channels
	channelId2channel = newChannelId2channel
	channelSyncLock.Unlock()
	logger.SysLog("channels synced from database")
}

func SyncChannelCache(frequency int) {
	for {
		time.Sleep(time.Duration(frequency) * time.Second)
		logger.SysLog("syncing channels from database")
		InitChannelCache()
	}
}

func CacheGetRandomSatisfiedChannel(group string, model string, ignoreFirstPriority bool) (*Channel, error) {
	if !config.MemoryCacheEnabled {
		return GetRandomSatisfiedChannel(group, model, ignoreFirstPriority)
	}
	channelSyncLock.RLock()
	defer channelSyncLock.RUnlock()
	channels := group2model2channels[group][model]
	if len(channels) == 0 {
		return nil, errors.New("channel not found")
	}
	endIdx := len(channels)
	// choose by priority
	firstChannel := channels[0]
	if firstChannel.GetPriority() > 0 {
		for i := range channels {
			if channels[i].GetPriority() != firstChannel.GetPriority() {
				endIdx = i
				break
			}
		}
	}
	idx := rand.Intn(endIdx)
	if ignoreFirstPriority {
		if endIdx < len(channels) { // which means there are more than one priority
			idx = random.RandRange(endIdx, len(channels))
		}
	}
	return channels[idx], nil
}

func GetChannelCandidates(group string, modelName string) []*Channel {
	if !config.MemoryCacheEnabled {
		var channels []*Channel
		err := DB.Where("status = ? AND models LIKE ? AND `group` LIKE ?",
			ChannelStatusEnabled, "%"+modelName+"%", "%"+group+"%").
			Order("priority DESC").
			Find(&channels).Error
		if err != nil {
			logger.SysError("failed to get channel candidates: " + err.Error())
			return nil
		}
		sort.Slice(channels, func(i, j int) bool {
			return channels[i].GetPriority() > channels[j].GetPriority()
		})
		return channels
	}
	channelSyncLock.RLock()
	defer channelSyncLock.RUnlock()
	channels := group2model2channels[group][modelName]
	if len(channels) == 0 {
		return nil
	}
	result := make([]*Channel, len(channels))
	copy(result, channels)
	return result
}

func CacheGetChannelById(channelId int) (*Channel, bool) {
	if !config.MemoryCacheEnabled {
		return nil, false
	}
	channelSyncLock.RLock()
	defer channelSyncLock.RUnlock()
	ch, ok := channelId2channel[channelId]
	return ch, ok
}

// GetFallbackChannel returns an enabled fallback channel for the given user
// group. It returns nil if no fallback channel is configured for that group.
// Among fallback channels with the same priority, the pick is random.
//
// Note: cooldown / concurrency / RPM are intentionally NOT checked here — the
// caller (controller/relay.go) acquires concurrency via channelrouter after
// picking, and a failed fallback call simply returns the error to the client.
func GetFallbackChannel(group string) (*Channel, error) {
	if group == "" {
		return nil, errors.New("group is empty")
	}
	var channels []*Channel
	groupCol := "`group`"
	if common.UsingPostgreSQL {
		groupCol = `"group"`
	}
	err := DB.Where("status = ? AND is_fallback = ? AND "+groupCol+" LIKE ?",
		ChannelStatusEnabled, true, "%"+group+"%").
		Order("fallback_priority DESC, id DESC").
		Find(&channels).Error
	if err != nil {
		logger.SysError("failed to query fallback channels: " + err.Error())
		return nil, err
	}
	if len(channels) == 0 {
		return nil, nil
	}

	// Walk the priority tier; skip channels that have no models configured.
	highest := channels[0].GetFallbackPriority()
	tier := make([]*Channel, 0, len(channels))
	for _, ch := range channels {
		if ch.GetFallbackPriority() != highest {
			break
		}
		if strings.TrimSpace(ch.Models) == "" {
			continue
		}
		tier = append(tier, ch)
	}
	if len(tier) == 0 {
		return nil, nil
	}
	return tier[rand.Intn(len(tier))], nil
}
