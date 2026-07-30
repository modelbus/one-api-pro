package model

type SyncEvent struct {
	Id        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	TableName string `json:"table_name" gorm:"size:64;index"`
	RowId     int64  `json:"row_id" gorm:"index"`
	RowKey    string `json:"row_key" gorm:"size:128"`
	Action    string `json:"action" gorm:"size:16"`
	Data      string `json:"data" gorm:"type:text"`
	NodeId    int    `json:"node_id"`
	CreatedAt int64  `json:"created_at" gorm:"bigint;index"`
	Pushed    int    `json:"pushed" gorm:"default:0"`
}