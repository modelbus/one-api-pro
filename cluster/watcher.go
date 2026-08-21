package cluster

import (
	"encoding/json"
	"reflect"

	"github.com/modelbus/one-api-pro/common/helper"
	"gorm.io/gorm"
)

type contextKeyType struct{}

var skipHookKey = contextKeyType{}

func RegisterCallbacks(db *gorm.DB) {
	if !Enabled {
		return
	}
	db.Callback().Create().After("gorm:create").Register("cluster:after_create", afterCreate)
	db.Callback().Update().After("gorm:update").Register("cluster:after_update", afterUpdate)
	db.Callback().Delete().Before("gorm:delete").Register("cluster:before_delete", beforeDelete)
	db.Callback().Delete().After("gorm:delete").Register("cluster:after_delete", afterDelete)
}

func afterCreate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	if isClusterInternalSchema(db.Statement.Schema.Name) {
		return
	}
	tableName := db.Statement.Table
	if !isSyncableTable(tableName) {
		return
	}
	rowId := getPrimaryKeyValue(db)
	data := marshalRowData(db)
	rowKey := ""
	if tableName == "options" {
		rowKey = getOptionKey(db)
	}
	createSyncEvent(db, tableName, rowId, rowKey, "insert", data)
}

func afterUpdate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	if isClusterInternalSchema(db.Statement.Schema.Name) {
		return
	}
	tableName := db.Statement.Table
	if !isSyncableTable(tableName) {
		return
	}
	rowId := getPrimaryKeyValue(db)
	rowKey := ""
	if tableName == "options" {
		rowKey = getOptionKey(db)
	}
	data := fetchCurrentRowData(db, tableName, rowId, rowKey)
	createSyncEvent(db, tableName, rowId, rowKey, "update", data)
}

func beforeDelete(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	if isClusterInternalSchema(db.Statement.Schema.Name) {
		return
	}
	tableName := db.Statement.Table
	if !isSyncableTable(tableName) {
		return
	}
	rowId := getPrimaryKeyValue(db)
	rowKey := ""
	if tableName == "options" {
		rowKey = getOptionKey(db)
	}
	data := fetchCurrentRowData(db, tableName, rowId, rowKey)
	db.InstanceSet("cluster:delete_data", data)
}

func afterDelete(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	if isClusterInternalSchema(db.Statement.Schema.Name) {
		return
	}
	tableName := db.Statement.Table
	if !isSyncableTable(tableName) {
		return
	}
	rowId := getPrimaryKeyValue(db)
	rowKey := ""
	if tableName == "options" {
		rowKey = getOptionKey(db)
	}
	data, _ := db.InstanceGet("cluster:delete_data")
	dataStr, _ := data.(string)
	createSyncEvent(db, tableName, rowId, rowKey, "delete", dataStr)
}

func isClusterInternalSchema(schemaName string) bool {
	return schemaName == "SyncEvent" || schemaName == "ClusterNode" || schemaName == "ClusterSyncProgress"
}

func getPrimaryKeyValue(db *gorm.DB) int64 {
	if db.Statement.Schema == nil {
		return 0
	}
	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			val := db.Statement.ReflectValue
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			if val.Kind() == reflect.Slice {
				if val.Len() == 0 {
					return 0
				}
				f := val.Index(0).FieldByName(field.Name)
				switch f.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					return f.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					return int64(f.Uint())
				}
			} else {
				f := val.FieldByName(field.Name)
				switch f.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					return f.Int()
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					return int64(f.Uint())
				}
			}
		}
	}
	return 0
}

func getOptionKey(db *gorm.DB) string {
	if db.Statement.Schema == nil {
		return ""
	}
	val := reflect.ValueOf(db.Statement.ReflectValue)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Slice {
		return ""
	}
	f := val.FieldByName("Key")
	if f.IsValid() {
		return f.String()
	}
	return ""
}

func marshalRowData(db *gorm.DB) string {
	val := db.Statement.ReflectValue
	if !val.IsValid() {
		return ""
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Slice {
		if val.Len() == 0 {
			return ""
		}
		val = val.Index(0)
	}
	if val.Kind() != reflect.Struct {
		return ""
	}
	data := val.Interface()
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func fetchCurrentRowData(db *gorm.DB, tableName string, rowId int64, rowKey string) string {
	// 使用 NewDB session 隔离 Statement，避免污染原 db 的状态
	skipDB := db.Session(&gorm.Session{NewDB: true})
	if tableName == "options" && rowKey != "" {
		var opt struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		err := skipDB.Table(tableName).Where("`key` = ?", rowKey).Select("`key`, `value`").Scan(&opt).Error
		if err != nil {
			return ""
		}
		bytes, _ := json.Marshal(opt)
		return string(bytes)
	}
	var result map[string]interface{}
	if rowId > 0 {
		skipDB.Table(tableName).Where("id = ?", rowId).Scan(&result)
	} else if rowKey != "" {
		skipDB.Table(tableName).Where("`key` = ?", rowKey).Scan(&result)
	}
	if result == nil {
		return ""
	}
	if _, ok := result["updated_at"]; ok {
		result["updated_at"] = helper.GetTimestamp()
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func createSyncEvent(db *gorm.DB, tableName string, rowId int64, rowKey string, action string, data string) {
	if data == "" && action != "delete" {
		return
	}
	if action == "update" {
		now := helper.GetTimestamp()
		if tableName != "options" && rowId > 0 {
			// 直接 SQL 更新 updated_at（不经过 GORM 回调，避免任何潜在递归）
			WithSkipHook(db).Exec(
				"UPDATE "+tableName+" SET updated_at = ? WHERE id = ?", now, rowId,
			)
		}
		data = fetchCurrentRowData(db, tableName, rowId, rowKey)
		if data == "" {
			return
		}
	}
	// 使用 Exec 直接执行 SQL INSERT，绕开 GORM 回调机制
	// 彻底避免 createSyncEvent -> DB.Create -> afterCreate -> createSyncEvent 的无限递归
	now := helper.GetTimestamp()
	skipDB := WithSkipHook(db)
	sqlStr := "INSERT INTO sync_events (table_name, row_id, row_key, action, data, node_id, created_at, pushed) VALUES (?, ?, ?, ?, ?, ?, ?, 0)"
	// 关键：使用 NewDB session 强制使用新的 statement，避免参数错乱
	tx := skipDB.Session(&gorm.Session{NewDB: true}).Exec(sqlStr, tableName, rowId, rowKey, action, data, NodeID, now)
	if tx.Error != nil {
		logClusterError("创建同步事件失败: " + tx.Error.Error())
		return
	}
	if tx.RowsAffected == 0 {
		logClusterError("创建同步事件失败: RowsAffected=0")
		return
	}
	NotifySyncEvent()
}