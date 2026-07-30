package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Leon-PanPan/one-api-pro/common/helper"
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/Leon-PanPan/one-api-pro/model"
	"gorm.io/gorm"
)

var syncNotifyChan = make(chan struct{}, 256)

func NotifySyncEvent() {
	select {
	case syncNotifyChan <- struct{}{}:
	default:
	}
}

func StartPusher(db *gorm.DB) {
	for {
		_, ok := <-syncNotifyChan
		if !ok {
			return
		}
		time.Sleep(100 * time.Millisecond)
		drainNotifyChan()
		pushEvents(db)
	}
}

func drainNotifyChan() {
	for {
		select {
		case <-syncNotifyChan:
		default:
			return
		}
	}
}

func pushEvents(db *gorm.DB) {
	var events []model.SyncEvent
	skipDB := WithSkipHook(db)
	skipDB.Where("pushed = 0").Order("created_at asc").Limit(BatchSize).Find(&events)
	if len(events) == 0 {
		return
	}

	now := helper.GetTimestamp()
	for i, event := range events {
		if event.Action != "delete" {
			// 尝试从数据库读取最新数据（仅用于非复合主键的表）
			// 对于复合主键表（abilities/channel_counters），watcher 端已经写入了完整数据
			data := fetchCurrentRowDataForPush(db, event.TableName, event.RowId, event.RowKey)
			if data != "" {
				events[i].Data = data
			}
		}
		// 注意：之前会在这里更新本机数据库的 updated_at，但：
		// 1. abilities/channel_counters 等表无 id 列
		// 2. watcher 端已经写入了完整数据，接收方会用 data 中的 updated_at
		// 3. 接收方在 applyInsert/applyUpdate 中会处理 updated_at
		// 因此这里不再做额外的 updated_at 更新
		_ = now
	}

	nodes := GetAliveNodesForSync(db)
	if len(nodes) == 0 {
		return
	}

	payload, err := json.Marshal(map[string]interface{}{
		"events": events,
	})
	if err != nil {
		logger.SysError("序列化同步事件失败: " + err.Error())
		return
	}

	allSuccess := true
	for _, node := range nodes {
		if !pushToNode(db, node, payload) {
			allSuccess = false
		}
	}

	if allSuccess {
		ids := make([]int64, len(events))
		for i, e := range events {
			ids[i] = e.Id
		}
		skipDB.Where("id IN ?", ids).Delete(&model.SyncEvent{})
	}
}

func fetchCurrentRowDataForPush(db *gorm.DB, tableName string, rowId int64, rowKey string) string {
	// 复合主键表（abilities, channel_counters）：watcher 端已经写入了完整数据，跳过重新查询
	if tableName == "abilities" || tableName == "channel_counters" {
		return ""
	}
	// options 表：用 key 查
	if tableName == "options" && rowKey != "" {
		var opt struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		db.Table(tableName).Where("`key` = ?", rowKey).Select("`key`, `value`").Scan(&opt)
		if opt.Key == "" {
			return ""
		}
		bytes, _ := json.Marshal(opt)
		return string(bytes)
	}
	// 单主键表：用 id 查
	if rowId > 0 {
		var result map[string]interface{}
		db.Table(tableName).Where("id = ?", rowId).Scan(&result)
		if result == nil {
			return ""
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
	if rowKey != "" {
		var result map[string]interface{}
		db.Table(tableName).Where("`key` = ?", rowKey).Scan(&result)
		if result == nil {
			return ""
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
	return ""
}

func pushToNode(db *gorm.DB, node model.ClusterNode, payload []byte) bool {
	url := node.Address + "/api/cluster/sync"
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		logger.SysError("推送事件到节点 " + node.Address + " 创建请求失败: " + err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	// 用目标节点的 secret（从本地 DB 查）
	secret := GetNodeSecret(db, node.NodeId)
	if secret == "" {
		secret = GetLocalSecret(db)
	}
	req.Header.Set("X-Cluster-Secret", secret)
	req.Header.Set("X-Cluster-Node-Id", fmt.Sprintf("%d", NodeID))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.SysError("推送事件到节点 " + node.Address + " 失败: " + err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		logger.SysError("推送事件到节点 " + node.Address + " 返回非200状态码: " + fmt.Sprintf("%d, body: %s", resp.StatusCode, string(body)))
		return false
	}
	return true
}

func doPost(url string, payload interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	secret := GetLocalSecret(clusterDB)
	req.Header.Set("X-Cluster-Secret", secret)
	req.Header.Set("X-Cluster-Node-Id", fmt.Sprintf("%d", NodeID))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func StartEventCleanup(db *gorm.DB) {
	for {
		time.Sleep(1 * time.Hour)
		skipDB := WithSkipHook(db)
		threshold := helper.GetTimestamp() - 7*24*3600
		skipDB.Where("created_at < ?", threshold).Delete(&model.SyncEvent{})
	}
}