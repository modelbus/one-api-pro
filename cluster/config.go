package cluster

import (
	"os"
	"strconv"
	"strings"

	"github.com/modelbus/one-api-pro/common/logger"
)

var (
	Enabled            bool
	NodeID             int
	NodeName           string
	NodeAddress        string
	Seeds              []string
	PushInterval       int
	DiscoveryInterval  int
	DeadPingInterval   int
	MaxPingFailures    int
	SyncLogs           bool
	BatchSize          int
)

func LoadConfig() {
	Enabled = os.Getenv("CLUSTER_ENABLED") == "true"
	if !Enabled {
		logger.SysLog("集群模式未启用，以单节点模式运行")
		return
	}
	NodeID, _ = strconv.Atoi(os.Getenv("CLUSTER_NODE_ID"))
	if NodeID <= 0 || NodeID >= 50 {
		logger.FatalLog("CLUSTER_NODE_ID 必须在 1-49 之间")
	}
	NodeName = os.Getenv("CLUSTER_NODE_NAME")
	if NodeName == "" {
		NodeName = "node-" + strconv.Itoa(NodeID)
	}
	NodeAddress = os.Getenv("CLUSTER_NODE_ADDRESS")
	if NodeAddress == "" {
		logger.FatalLog("CLUSTER_NODE_ADDRESS 不能为空")
	}
	if os.Getenv("CLUSTER_SECRET") == "" {
		logger.FatalLog("CLUSTER_SECRET 不能为空（用作本机 secret 初始值）")
	}
	seeds := os.Getenv("CLUSTER_SEEDS")
	if seeds != "" {
		Seeds = strings.Split(seeds, ",")
	}
	PushInterval = envInt("CLUSTER_PUSH_INTERVAL", 3)
	DiscoveryInterval = envInt("CLUSTER_DISCOVERY_INTERVAL", 30)
	DeadPingInterval = envInt("CLUSTER_DEAD_PING_INTERVAL", 120)
	MaxPingFailures = envInt("CLUSTER_MAX_PING_FAILURES", 3)
	SyncLogs = os.Getenv("CLUSTER_SYNC_LOGS") != "false"
	BatchSize = envInt("CLUSTER_BATCH_SIZE", 50)
	logger.SysLogf("集群模式已启用，节点编号=%d，名称=%s，地址=%s", NodeID, NodeName, NodeAddress)
}

func envInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}