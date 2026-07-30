package cluster

var syncTableList = []string{
	"users", "tokens", "channels", "abilities",
	"options", "redemptions", "plans",
	"user_plans", "plan_usages",
	"channel_counters",
}

func isSyncableTable(tableName string) bool {
	for _, t := range syncTableList {
		if t == tableName {
			return true
		}
	}
	if tableName == "logs" {
		return SyncLogs
	}
	return false
}