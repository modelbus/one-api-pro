package model

type ChannelCounter struct {
	ChannelId   int   `json:"channel_id" gorm:"primaryKey"`
	NodeId      int   `json:"node_id" gorm:"primaryKey"`
	Concurrency int   `json:"concurrency" gorm:"default:0"`
	RpmCount    int   `json:"rpm_count" gorm:"default:0"`
	RpmMinute   int64 `json:"rpm_minute" gorm:"default:0"`
	UpdatedAt   int64 `json:"updated_at" gorm:"bigint;default:0"`
}