package cluster

import (
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/Leon-PanPan/one-api-pro/model"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	LoadConfig()
	if !Enabled {
		return
	}

	SetDB(db)

	if err := db.AutoMigrate(&model.ClusterNode{}); err != nil {
		logger.FatalLog("自动迁移 cluster_nodes 表失败: " + err.Error())
	}
	if err := db.AutoMigrate(&model.SyncEvent{}); err != nil {
		logger.FatalLog("自动迁移 sync_events 表失败: " + err.Error())
	}
	if err := db.AutoMigrate(&model.ChannelCounter{}); err != nil {
		logger.FatalLog("自动迁移 channel_counters 表失败: " + err.Error())
	}

	RegisterCallbacks(db)

	SaveLocalNode(db)

	discoverOnce(db)

	if len(Seeds) > 0 {
		pingSeedNodes(db)
	}

	go StartDiscovery(db)
	go StartPusher(db)
	go StartEventCleanup(db)

	// 补推启动前的未推送事件
	var count int64
	WithSkipHook(db).Model(&model.SyncEvent{}).Where("pushed = 0").Count(&count)
	if count > 0 {
		logger.SysLogf("[集群] 发现 %d 条未推送事件，开始补推", count)
		NotifySyncEvent()
	}

	logClusterInfo("集群模块初始化完成")
}