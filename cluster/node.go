package cluster

import (
	"encoding/json"
	"os"
	"time"

	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
	"gorm.io/gorm"
)

type ClusterNodeInfo struct {
	NodeId    int    `json:"node_id"`
	NodeName  string `json:"node_name"`
	Address   string `json:"address"`
	Status    int    `json:"status"`
	SecretKey string `json:"secret_key"`
}

type PingRequest struct {
	NodeId   int    `json:"node_id"`
	NodeName string `json:"node_name"`
	Address  string `json:"address"`
	Secret   string `json:"secret"`
}

type PingResponse struct {
	Success bool              `json:"success"`
	Node    ClusterNodeInfo   `json:"node"`
	Nodes   []ClusterNodeInfo `json:"nodes"`
}

func SaveLocalNode(db *gorm.DB) {
	now := helper.GetTimestamp()
	skipDB := WithSkipHook(db)

	// 读取本机节点记录（如果存在）
	var existing model.ClusterNode
	findResult := skipDB.Where("node_id = ?", NodeID).First(&existing)
	found := findResult.Error == nil

	if found {
		// 已存在：更新 address、node_name、状态、心跳（保留 secret_key，由 admin UI 单独管理）
		skipDB.Model(&existing).Updates(map[string]interface{}{
			"address":        NodeAddress,
			"node_name":      NodeName,
			"status":         1,
			"last_heartbeat": now,
			"ping_failures":  0,
			"updated_at":     now,
		})
	} else {
		// 不存在：创建
		initialSecret := os.Getenv("CLUSTER_SECRET")
		if initialSecret == "" {
			logger.FatalLog("CLUSTER_SECRET 环境变量不能为空")
		}
		node := model.ClusterNode{
			NodeId:        NodeID,
			NodeName:      NodeName,
			Address:       NodeAddress,
			SecretKey:     initialSecret,
			Status:        1,
			Disabled:      false,
			LastHeartbeat: now,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := skipDB.Create(&node).Error; err != nil {
			logger.FatalLog("保存本节点信息失败: " + err.Error())
		}
		logger.SysLogf("[集群] 本节点 (id=%d) 信息已初始化到数据库", NodeID)
	}
}

func GetLocalSecret(db *gorm.DB) string {
	var node model.ClusterNode
	result := WithSkipHook(db).Where("node_id = ?", NodeID).First(&node)
	if result.Error != nil {
		return ""
	}
	return node.SecretKey
}

func GetNodeSecret(db *gorm.DB, nodeId int) string {
	var node model.ClusterNode
	result := WithSkipHook(db).Where("node_id = ?", nodeId).First(&node)
	if result.Error != nil {
		return ""
	}
	return node.SecretKey
}

func GetAliveNodesForSync(db *gorm.DB) []model.ClusterNode {
	var nodes []model.ClusterNode
	skipDB := WithSkipHook(db)
	skipDB.Where("status = 1 AND disabled = ? AND node_id != ?", false, NodeID).Find(&nodes)
	return nodes
}

func GetAllRemoteNodes(db *gorm.DB) []model.ClusterNode {
	var nodes []model.ClusterNode
	skipDB := WithSkipHook(db)
	skipDB.Where("node_id != ?", NodeID).Find(&nodes)
	return nodes
}

func GetAllNodes(db *gorm.DB) []model.ClusterNode {
	var nodes []model.ClusterNode
	skipDB := WithSkipHook(db)
	skipDB.Where("disabled = ?", false).Order("node_id asc").Find(&nodes)
	return nodes
}

func StartDiscovery(db *gorm.DB) {
	for {
		time.Sleep(time.Duration(DiscoveryInterval) * time.Second)
		discoverOnce(db)
	}
}

func discoverOnce(db *gorm.DB) {
	nodes := GetAllRemoteNodes(db)
	now := helper.GetTimestamp()
	skipDB := WithSkipHook(db)

	skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", NodeID).Updates(map[string]interface{}{
		"last_heartbeat": now,
		"status":         1,
		"updated_at":     now,
	})

	for i := range nodes {
		node := &nodes[i]
		if node.Disabled {
			continue
		}
		if node.Status == 1 {
			pingAliveNode(db, node, now)
		} else {
			pingDeadNode(db, node, now)
		}
	}
}

func pingAliveNode(db *gorm.DB, node *model.ClusterNode, now int64) {
	resp, err := PingNode(db, node)
	skipDB := WithSkipHook(db)
	if err != nil || (resp != nil && !resp.Success) {
		node.PingFailures++
		if node.PingFailures >= MaxPingFailures {
			skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", node.NodeId).Updates(map[string]interface{}{
				"status":          2,
				"ping_failures":   node.PingFailures,
				"updated_at":      now,
			})
			logger.SysErrorf("[集群] 节点 %s (id=%d) 连续 %d 次 ping 失败，标记为失败状态",
				node.Address, node.NodeId, node.PingFailures)
		} else {
			skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", node.NodeId).Updates(map[string]interface{}{
				"ping_failures":    node.PingFailures,
				"last_ping_attempt": now,
				"updated_at":       now,
			})
		}
		return
	}

	skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", node.NodeId).Updates(map[string]interface{}{
		"status":           1,
		"ping_failures":    0,
		"last_heartbeat":   now,
		"last_ping_attempt": now,
		"updated_at":       now,
	})
	node.PingFailures = 0
	node.Status = 1

	mergeDiscoveredNodes(db, resp.Nodes)
}

func pingDeadNode(db *gorm.DB, node *model.ClusterNode, now int64) {
	if now-node.LastPingAttempt < int64(DeadPingInterval) {
		return
	}

	skipDB := WithSkipHook(db)
	skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", node.NodeId).Updates(map[string]interface{}{
		"last_ping_attempt": now,
		"updated_at":       now,
	})
	node.LastPingAttempt = now

	resp, err := PingNode(db, node)
	if err != nil || (resp != nil && !resp.Success) {
		logger.SysErrorf("[集群] 失败节点 %s (id=%d) ping 仍然失败: %v", node.Address, node.NodeId, err)
		return
	}

	skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", node.NodeId).Updates(map[string]interface{}{
		"status":         1,
		"ping_failures":  0,
		"last_heartbeat": now,
		"updated_at":     now,
	})
	node.Status = 1

	logger.SysLogf("[集群] 失败节点 %s (id=%d) 已恢复", node.Address, node.NodeId)

	mergeDiscoveredNodes(db, resp.Nodes)
}

func PingNode(db *gorm.DB, node *model.ClusterNode) (*PingResponse, error) {
	url := node.Address + "/api/cluster/ping"

	// 用目标节点的 secret（从本地 DB 查），如果没有则用本机的（兼容旧版）
	secret := GetNodeSecret(db, node.NodeId)
	if secret == "" {
		secret = GetLocalSecret(db)
	}

	payload := PingRequest{
		NodeId:   NodeID,
		NodeName: NodeName,
		Address:  NodeAddress,
		Secret:   secret,
	}
	respBody, err := doPost(url, payload)
	if err != nil {
		return nil, err
	}

	var result PingResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		logger.SysError("[集群] 解析 ping 响应失败: " + err.Error())
		return nil, err
	}

	return &result, nil
}

func mergeDiscoveredNodes(db *gorm.DB, nodes []ClusterNodeInfo) {
	skipDB := WithSkipHook(db)
	now := helper.GetTimestamp()
	for _, n := range nodes {
		if n.NodeId == NodeID {
			continue
		}

		var existing model.ClusterNode
		result := skipDB.Where("node_id = ?", n.NodeId).First(&existing)
		if result.Error != nil {
			newNode := model.ClusterNode{
				NodeId:        n.NodeId,
				NodeName:      n.NodeName,
				Address:       n.Address,
				SecretKey:     n.SecretKey,
				Status:        n.Status,
				Disabled:      false,
				LastHeartbeat: now,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			skipDB.Create(&newNode)
			logger.SysLogf("[集群] 发现新节点 %s (id=%d)", n.Address, n.NodeId)
		} else {
			updates := map[string]interface{}{
				"address":    n.Address,
				"node_name":  n.NodeName,
				"secret_key": n.SecretKey,
				"updated_at": now,
			}
			if !existing.Disabled {
				if n.Status == 1 {
					updates["status"] = 1
				}
			}
			skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", n.NodeId).Updates(updates)
		}
	}
}

func pingSeedNodes(db *gorm.DB) {
	if len(Seeds) == 0 {
		return
	}

	for _, seed := range Seeds {
		node := &model.ClusterNode{
			NodeId:   0,
			Address:  seed,
			Status:   1,
		}
		resp, err := PingNode(db, node)
		if err != nil {
			logger.SysError("[集群] ping 种子节点 " + seed + " 失败: " + err.Error())
			continue
		}
		if resp.Success {
			logger.SysLogf("[集群] 成功 ping 种子节点 %s，获取到 %d 个节点", seed, len(resp.Nodes))
			mergeDiscoveredNodes(db, resp.Nodes)
			return
		}
	}
	logger.SysError("[集群] 未能 ping 通任何种子节点")
}