package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/cluster"
	"github.com/modelbus/one-api-pro/common/helper"
	cluster_model "github.com/modelbus/one-api-pro/model"
)

func GetAllClusterNodes(c *gin.Context) {
	if !cluster.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "集群模式未启用",
		})
		return
	}
	nodes := cluster.GetAllNodes(cluster.GetDB())
	type nodeWithSelfFlag struct {
		cluster_model.ClusterNode
		IsSelf bool `json:"is_self"`
	}
	result := make([]nodeWithSelfFlag, 0, len(nodes))
	for _, n := range nodes {
		result = append(result, nodeWithSelfFlag{
			ClusterNode: n,
			IsSelf:      n.NodeId == cluster.NodeID,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    result,
	})
}

func GetClusterNode(c *gin.Context) {
	if !cluster.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "集群模式未启用",
		})
		return
	}
	nodeId, _ := strconv.Atoi(c.Param("id"))
	var node cluster_model.ClusterNode
	result := cluster.WithSkipHook(cluster.GetDB()).Where("node_id = ?", nodeId).First(&node)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    node,
	})
}

type ClusterNodeRequest struct {
	NodeId   int    `json:"node_id"`
	NodeName string `json:"node_name"`
	Address  string `json:"address"`
	Secret   string `json:"secret"`
}

func AddClusterNode(c *gin.Context) {
	if !cluster.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "集群模式未启用",
		})
		return
	}
	var req ClusterNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if req.NodeId <= 0 || req.NodeId >= 50 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点编号必须在 1-49 之间",
		})
		return
	}
	if req.Address == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点地址不能为空",
		})
		return
	}
	if req.Secret == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "secret 不能为空（其他节点连接此节点时的认证密钥）",
		})
		return
	}
	now := helper.GetTimestamp()

	skipDB := cluster.WithSkipHook(cluster.GetDB())
	var existing cluster_model.ClusterNode
	if err := skipDB.Where("node_id = ?", req.NodeId).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点编号已存在",
		})
		return
	}

	node := cluster_model.ClusterNode{
		NodeId:        req.NodeId,
		NodeName:      req.NodeName,
		Address:       req.Address,
		SecretKey:     req.Secret,
		Status:        1,
		Disabled:      false,
		LastHeartbeat: now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := skipDB.Create(&node).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    node,
	})
}

func UpdateClusterNode(c *gin.Context) {
	if !cluster.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "集群模式未启用",
		})
		return
	}
	var req ClusterNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if req.NodeId <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点编号无效",
		})
		return
	}

	now := helper.GetTimestamp()
	skipDB := cluster.WithSkipHook(cluster.GetDB())
	updates := map[string]interface{}{
		"updated_at": now,
	}
	if req.NodeName != "" {
		updates["node_name"] = req.NodeName
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Secret != "" {
		updates["secret_key"] = req.Secret
		updates["status"] = 1
		updates["ping_failures"] = 0
	}

	result := skipDB.Model(&cluster_model.ClusterNode{}).Where("node_id = ?", req.NodeId).Updates(updates)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

// DeleteClusterNode 软删除：设置 disabled = true
func DeleteClusterNode(c *gin.Context) {
	if !cluster.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "集群模式未启用",
		})
		return
	}
	nodeId, _ := strconv.Atoi(c.Param("id"))
	if nodeId == cluster.NodeID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "不能禁用当前节点",
		})
		return
	}
	skipDB := cluster.WithSkipHook(cluster.GetDB())
	now := helper.GetTimestamp()
	result := skipDB.Model(&cluster_model.ClusterNode{}).Where("node_id = ?", nodeId).Updates(map[string]interface{}{
		"disabled":   true,
		"status":     2,
		"updated_at": now,
	})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "节点已禁用（物理删除需要手动 SQL: DELETE FROM cluster_nodes WHERE node_id = ?）",
	})
}

// EnableClusterNode 重新启用已禁用的节点
func EnableClusterNode(c *gin.Context) {
	if !cluster.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "集群模式未启用",
		})
		return
	}
	nodeId, _ := strconv.Atoi(c.Param("id"))
	skipDB := cluster.WithSkipHook(cluster.GetDB())
	now := helper.GetTimestamp()
	result := skipDB.Model(&cluster_model.ClusterNode{}).Where("node_id = ?", nodeId).Updates(map[string]interface{}{
		"disabled":       false,
		"status":         1,
		"ping_failures":  0,
		"last_heartbeat": now,
		"updated_at":     now,
	})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func PingClusterNode(c *gin.Context) {
	if !cluster.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "集群模式未启用",
		})
		return
	}
	nodeId, _ := strconv.Atoi(c.Param("id"))
	var node cluster_model.ClusterNode
	result := cluster.WithSkipHook(cluster.GetDB()).Where("node_id = ?", nodeId).First(&node)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "节点不存在",
		})
		return
	}

	resp, err := cluster.PingNode(cluster.GetDB(), &node)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "ping 失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    resp,
	})
}