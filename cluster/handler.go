package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
	"gorm.io/gorm"
)

func clusterAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.GetHeader("X-Cluster-Secret")
		// 优先用 header 验证，header 不匹配则用 body 中的 secret
		if secret == "" {
			secret = c.GetHeader("X-Cluster-Secret-Body")
		}
		if secret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "集群认证失败：缺少 secret"})
			c.Abort()
			return
		}
		// 从数据库读本机 secret
		localSecret := GetLocalSecret(clusterDB)
		if localSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "集群未初始化 secret"})
			c.Abort()
			return
		}
		if secret != localSecret {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "集群认证失败"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RegisterRoutes(r *gin.Engine) {
	if !Enabled {
		return
	}
	g := r.Group("/api/cluster")
	g.Use(clusterAuth())
	{
		g.POST("/ping", handlePing)
		g.POST("/sync", handleSync)
	}
}

func handlePing(c *gin.Context) {
	var req PingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	if req.Secret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "集群认证失败：缺少 secret"})
		return
	}

	// 验证 secret 匹配本机
	localSecret := GetLocalSecret(clusterDB)
	if req.Secret != localSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "集群认证失败"})
		return
	}

	if req.NodeId == NodeID {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "不能 ping 自己"})
		return
	}

	now := helper.GetTimestamp()
	skipDB := WithSkipHook(clusterDB)

	// 查找已存在的节点记录
	var existing model.ClusterNode
	result := skipDB.Where("node_id = ?", req.NodeId).First(&existing)

	if result.Error != nil {
		// 新节点：创建
		newNode := model.ClusterNode{
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
		skipDB.Create(&newNode)
		logger.SysLogf("[集群] 新节点注册: %s (id=%d)", req.Address, req.NodeId)
	} else {
		// 已存在：检查是否被禁用
		if existing.Disabled {
			// 节点已被禁用：仅响应 ping（让对方知道我们在线），不更新对方信息
			c.JSON(http.StatusOK, PingResponse{
				Success: true,
				Node: ClusterNodeInfo{
					NodeId:    NodeID,
					NodeName:  NodeName,
					Address:   NodeAddress,
					Status:    1,
					SecretKey: localSecret,
				},
				Nodes: []ClusterNodeInfo{}, // 不返回任何节点信息
			})
			return
		}
		// 正常更新
		skipDB.Model(&model.ClusterNode{}).Where("node_id = ?", req.NodeId).Updates(map[string]interface{}{
			"node_name":      req.NodeName,
			"address":        req.Address,
			"status":         1,
			"ping_failures":  0,
			"last_heartbeat": now,
			"updated_at":     now,
		})
	}

	// 响应：返回本机 + 所有未禁用的节点
	allNodes := GetAllNodes(clusterDB)
	nodeInfos := make([]ClusterNodeInfo, 0, len(allNodes))
	for _, n := range allNodes {
		nodeInfos = append(nodeInfos, ClusterNodeInfo{
			NodeId:    n.NodeId,
			NodeName:  n.NodeName,
			Address:   n.Address,
			Status:    n.Status,
			SecretKey: n.SecretKey,
		})
	}

	selfNode := ClusterNodeInfo{
		NodeId:    NodeID,
		NodeName:  NodeName,
		Address:   NodeAddress,
		Status:    1,
		SecretKey: localSecret,
	}

	c.JSON(http.StatusOK, PingResponse{
		Success: true,
		Node:    selfNode,
		Nodes:   nodeInfos,
	})
}

func handleSync(c *gin.Context) {
	var req struct {
		Events []model.SyncEvent `json:"events"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}
	// 防御性：捕获 goroutine 中的 panic，防止单个请求的 panic 拖垮整个进程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.SysError(fmt.Sprintf("[集群] handleSync goroutine panic recovered: %v\n%s", r, debug.Stack()))
			}
		}()
		ApplyEvents(req.Events)
	}()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func parseTables(tables string) []string {
	var result []string
	err := json.Unmarshal([]byte("["+tables+"]"), &result)
	if err != nil {
		result = []string{}
		for _, t := range syncTableList {
			result = append(result, t)
		}
	}
	return result
}

func strconvAtoi(s string) int {
	if s == "" {
		return 0
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

var clusterDB *gorm.DB

func SetDB(db *gorm.DB) {
	clusterDB = db
}

func GetDB() *gorm.DB {
	return clusterDB
}

func WithSkipHook(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{SkipHooks: true})
}

func logClusterError(msg string) {
	logger.SysError("[集群] " + msg)
}

func logClusterInfo(msg string) {
	logger.SysLog("[集群] " + msg)
}