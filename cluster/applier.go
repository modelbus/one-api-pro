package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
	"gorm.io/gorm"
)

func ApplyEvents(events []model.SyncEvent) {
	if !Enabled {
		return
	}
	// 防御性：捕获 ApplyEvents 中的 panic，防止单个事件处理失败导致整个进程崩溃
	defer func() {
		if r := recover(); r != nil {
			logger.SysError(fmt.Sprintf("[集群] ApplyEvents panic recovered: %v\n%s", r, debug.Stack()))
		}
	}()
	skipDB := WithSkipHook(clusterDB)
	now := helper.GetTimestamp()

	for _, event := range events {
		if event.NodeId == NodeID {
			continue
		}
		applyOne(skipDB, event, now)
	}
}

func applyOne(db *gorm.DB, event model.SyncEvent, now int64) {
	switch event.Action {
	case "insert":
		applyInsert(db, event, now)
	case "update":
		applyUpdate(db, event, now)
	case "delete":
		applyDelete(db, event)
	}
	invalidateCache(event)
}

func applyInsert(db *gorm.DB, event model.SyncEvent, now int64) {
	if event.TableName == "options" {
		applyUpsertOptions(db, event, now)
		return
	}
	if event.TableName == "channel_counters" {
		applyUpsertChannelCounter(db, event)
		return
	}
	if event.TableName == "abilities" {
		applyUpsertAbility(db, event, now)
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Data), &data); err != nil {
		logClusterError("解析插入数据失败: " + err.Error())
		return
	}

	// 标准化 created_at/updated_at 为 int64
	if v, ok := data["created_at"].(float64); ok {
		data["created_at"] = int64(v)
	}
	if v, ok := data["updated_at"].(float64); ok {
		data["updated_at"] = int64(v)
	}

	// 提取主键 id
	idValue, ok := data["id"]
	if !ok {
		logClusterError("插入数据缺少 id 字段")
		return
	}
	var id int64
	switch v := idValue.(type) {
	case float64:
		id = int64(v)
	case int64:
		id = v
	case int:
		id = int64(v)
	default:
		logClusterError("不支持的 id 类型")
		return
	}

	// 检查主键是否已存在
	var count int64
	db.Table(event.TableName).Where("id = ?", id).Count(&count)
	if count > 0 {
		// 已存在，转为 update 逻辑
		applyUpdate(db, event, now)
		return
	}

	// 不存在，插入
	delete(data, "id") // 让数据库生成新的 id？不对——要保留原 id 以保证一致
	// 重新加上 id 用于 create
	data["id"] = id
	if _, ok := data["updated_at"]; ok {
		data["updated_at"] = now
	}
	db.Table(event.TableName).Create(&data)
}

// applyUpsertAbility 处理 abilities 表的复合主键 (group, model, channel_id)
func applyUpsertAbility(db *gorm.DB, event model.SyncEvent, now int64) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Data), &data); err != nil {
		logClusterError("解析 abilities 数据失败: " + err.Error())
		return
	}

	// 提取复合主键字段
	group, _ := data["group"].(string)
	model_, _ := data["model"].(string)
	channelId, _ := toInt64(data["channel_id"])
	if group == "" || model_ == "" || channelId == 0 {
		logClusterError("abilities 数据缺少复合主键字段")
		return
	}

	// 检查主键是否已存在
	var count int64
	db.Table("abilities").Where("`group` = ? AND model = ? AND channel_id = ?", group, model_, channelId).Count(&count)
	if count > 0 {
		// 已存在：更新
		db.Table("abilities").Where("`group` = ? AND model = ? AND channel_id = ?", group, model_, channelId).Updates(data)
		return
	}

	// 不存在，插入
	db.Table("abilities").Create(&data)
}

func toInt64(v interface{}) (int64, bool) {
	switch x := v.(type) {
	case float64:
		return int64(x), true
	case int64:
		return x, true
	case int:
		return int64(x), true
	default:
		return 0, false
	}
}

func applyUpdate(db *gorm.DB, event model.SyncEvent, now int64) {
	if event.TableName == "options" {
		applyUpsertOptions(db, event, now)
		return
	}
	if event.TableName == "channel_counters" {
		applyUpsertChannelCounter(db, event)
		return
	}
	if event.TableName == "abilities" {
		applyUpsertAbility(db, event, now)
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Data), &data); err != nil {
		logClusterError("解析更新数据失败: " + err.Error())
		return
	}

	// 标准化时间字段
	if v, ok := data["updated_at"].(float64); ok {
		data["updated_at"] = int64(v)
	}
	if v, ok := data["created_at"].(float64); ok {
		data["created_at"] = int64(v)
	}

	// updated_at 比较
	incomingUpdatedAt, _ := data["updated_at"].(int64)
	var localUpdatedAt int64
	db.Table(event.TableName).Where("id = ?", event.RowId).Select("updated_at").Scan(&localUpdatedAt)
	if localUpdatedAt > 0 && incomingUpdatedAt <= localUpdatedAt {
		// 本地数据较新或相同，跳过
		return
	}

	if _, ok := data["updated_at"]; ok {
		data["updated_at"] = now
	}
	delete(data, "id")
	db.Table(event.TableName).Where("id = ?", event.RowId).Updates(data)
}

func applyDelete(db *gorm.DB, event model.SyncEvent) {
	if event.TableName == "channel_counters" {
		db.Table(event.TableName).Where("channel_id = ? AND node_id != ?", event.RowId, NodeID).Delete(nil)
		return
	}
	// abilities 表用复合主键 (group, model, channel_id)
	if event.TableName == "abilities" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(event.Data), &data); err == nil {
			group, _ := data["group"].(string)
			model_, _ := data["model"].(string)
			channelId, _ := toInt64(data["channel_id"])
			if group != "" && model_ != "" && channelId > 0 {
				db.Table("abilities").Where("`group` = ? AND model = ? AND channel_id = ?", group, model_, channelId).Delete(nil)
			}
		}
		return
	}
	if event.RowId > 0 {
		db.Table(event.TableName).Where("id = ?", event.RowId).Delete(nil)
	} else if event.RowKey != "" {
		db.Table(event.TableName).Where("`key` = ?", event.RowKey).Delete(nil)
	} else {
		// 防御性：row_id 和 row_key 都为空，说明发送方的 Delete 函数有 bug
		// （如 DB.Where("id = ?", id).Delete(&Model{}) 模式，Model 零值导致 GORM 回调取不到 Id）
		// 尝试从 data JSON 中提取 id
		if len(event.Data) > 0 {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(event.Data), &data); err == nil {
				if idValue, ok := data["id"]; ok {
					var id int64
					switch v := idValue.(type) {
					case float64:
						id = int64(v)
					case int64:
						id = v
					case int:
						id = int64(v)
					}
					if id > 0 {
						logger.SysLogf("[集群] 收到表 %s 的删除事件，从 data 中提取 id=%d 进行删除", event.TableName, id)
						db.Table(event.TableName).Where("id = ?", id).Delete(nil)
						return
					}
				}
			}
		}
		logClusterError(fmt.Sprintf("收到表 %s 的删除事件但 row_id、row_key、data.id 都为空，跳过（发送方 Delete 函数可能有 bug）", event.TableName))
	}
}

func applyUpsertOptions(db *gorm.DB, event model.SyncEvent, now int64) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Data), &data); err != nil {
		logClusterError("解析选项数据失败: " + err.Error())
		return
	}
	key, _ := data["key"].(string)
	value, _ := data["value"].(string)
	if key == "" {
		return
	}
	var count int64
	db.Table("options").Where("`key` = ?", key).Count(&count)
	if count > 0 {
		db.Table("options").Where("`key` = ?", key).Updates(map[string]interface{}{
			"value":      value,
			"updated_at": now,
		})
	} else {
		db.Table("options").Create(map[string]interface{}{
			"key":        key,
			"value":      value,
			"created_at": now,
			"updated_at": now,
		})
	}
}

func applyUpsertChannelCounter(db *gorm.DB, event model.SyncEvent) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(event.Data), &data); err != nil {
		logClusterError("解析计数器数据失败: " + err.Error())
		return
	}
	channelId, _ := data["channel_id"].(float64)
	nodeId, _ := data["node_id"].(float64)
	concurrency, _ := data["concurrency"].(float64)
	rpmCount, _ := data["rpm_count"].(float64)
	rpmMinute, _ := data["rpm_minute"].(float64)
	updatedAt, _ := data["updated_at"].(float64)

	if int(nodeId) == NodeID {
		return
	}
	var count int64
	db.Table("channel_counters").Where("channel_id = ? AND node_id = ?", int(channelId), int(nodeId)).Count(&count)
	if count > 0 {
		db.Table("channel_counters").Where("channel_id = ? AND node_id = ?", int(channelId), int(nodeId)).Updates(map[string]interface{}{
			"concurrency": int(concurrency),
			"rpm_count":   int(rpmCount),
			"rpm_minute":  int64(rpmMinute),
			"updated_at":  int64(updatedAt),
		})
	} else {
		db.Table("channel_counters").Create(map[string]interface{}{
			"channel_id":  int(channelId),
			"node_id":     int(nodeId),
			"concurrency": int(concurrency),
			"rpm_count":   int(rpmCount),
			"rpm_minute":  int64(rpmMinute),
			"updated_at":  int64(updatedAt),
		})
	}
}

func invalidateCache(event model.SyncEvent) {
	// 检查 Redis 是否可用，未连接则直接跳过（不修改主程序的 common/redis.go）
	if !common.RedisEnabled || common.RDB == nil {
		return
	}
	switch event.TableName {
	case "users":
		common.RedisDel(fmt.Sprintf("user_group:%d", event.RowId))
		common.RedisDel(fmt.Sprintf("user_quota:%d", event.RowId))
		common.RedisDel(fmt.Sprintf("user_enabled:%d", event.RowId))
		common.RedisDel(fmt.Sprintf("user_plans:%d", event.RowId))
	case "tokens":
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(event.Data), &data); err == nil {
			if key, ok := data["key"].(string); ok {
				common.RedisDel("token:" + key)
			}
		}
	case "channels":
		model.InitChannelCache()
		ctx := context.Background()
		keys, err := common.RDB.Keys(ctx, "group_models:*").Result()
		if err == nil {
			for _, key := range keys {
				common.RedisDel(key)
			}
		}
	case "options":
		model.InitOptionMap()
	case "user_plans":
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(event.Data), &data); err == nil {
			if userId, ok := data["user_id"].(float64); ok {
				common.RedisDel(fmt.Sprintf("user_plans:%d", int(userId)))
			}
		}
	}
}