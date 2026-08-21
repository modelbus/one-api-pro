package model

import (
	"time"

	"github.com/modelbus/one-api-pro/common/helper"
	"gorm.io/gorm"
)

type ClusterNode struct {
	Id              int    `json:"id" gorm:"primaryKey;autoIncrement"`
	NodeId          int    `json:"node_id" gorm:"uniqueIndex"`
	NodeName        string `json:"node_name" gorm:"size:64"`
	Address         string `json:"address" gorm:"size:256"`
	SecretKey       string `json:"secret_key" gorm:"size:128"`
	Status          int    `json:"status" gorm:"default:1"`
	Disabled        bool   `json:"disabled" gorm:"default:false"`
	LastHeartbeat   int64  `json:"last_heartbeat" gorm:"bigint;default:0"`
	PingFailures    int    `json:"ping_failures" gorm:"default:0"`
	LastPingAttempt int64  `json:"last_ping_attempt" gorm:"bigint;default:0"`
	CreatedAt       int64  `json:"created_at" gorm:"bigint;default:0"`
	UpdatedAt       int64  `json:"updated_at" gorm:"bigint;default:0"`
}

func (n *ClusterNode) BeforeCreate(tx *gorm.DB) error {
	if n.CreatedAt == 0 {
		n.CreatedAt = helper.GetTimestamp()
	}
	if n.UpdatedAt == 0 {
		n.UpdatedAt = n.CreatedAt
	}
	return nil
}

func (n *ClusterNode) BeforeUpdate(tx *gorm.DB) error {
	n.UpdatedAt = helper.GetTimestamp()
	return nil
}

func (n *ClusterNode) IsAlive() bool {
	return n.Status == 1 && !n.Disabled
}

func (n *ClusterNode) MarkHeartbeat() {
	n.LastHeartbeat = time.Now().Unix()
	n.UpdatedAt = n.LastHeartbeat
}