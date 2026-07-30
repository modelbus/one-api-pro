package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Option struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

type ModelPrice struct {
	Id              uint    `gorm:"primaryKey;autoIncrement"`
	ModelName       string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	InputPrice      float64 `gorm:"type:decimal(16,6);default:0;not null"`
	OutputPrice     float64 `gorm:"type:decimal(16,6);default:0;not null"`
	CachedPrice     float64 `gorm:"type:decimal(16,6);default:0;not null"`
	PerRequestPrice float64 `gorm:"type:decimal(16,6);default:0;not null"`
	BillingType     string  `gorm:"type:varchar(20);default:'token';not null"`
	Enabled         bool    `gorm:"default:true;not null"`
	CreatedAt       int64   `gorm:"bigint;default:0"`
	UpdatedAt       int64   `gorm:"bigint;default:0"`
}

type GroupPrice struct {
	Id        uint    `gorm:"primaryKey;autoIncrement"`
	GroupName string  `gorm:"type:varchar(32);uniqueIndex:idx_group_model;not null"`
	ModelName string  `gorm:"type:varchar(100);uniqueIndex:idx_group_model;default:'';not null"`
	Discount  float64 `gorm:"type:decimal(10,4);default:1;not null"`
	CreatedAt int64   `gorm:"bigint;default:0"`
	UpdatedAt int64   `gorm:"bigint;default:0"`
}

func main() {
	dsn := os.Getenv("SQL_DSN")
	if dsn == "" {
		fmt.Println("请设置 SQL_DSN 环境变量")
		os.Exit(1)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(&ModelPrice{}, &GroupPrice{}); err != nil {
		fmt.Printf("自动迁移失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("表结构已确保存在")

	nowTs := time.Now().Unix()

	var modelCount int64
	db.Model(&ModelPrice{}).Count(&modelCount)
	if modelCount > 0 {
		fmt.Printf("model_prices 表已有 %d 条数据，跳过 ModelRatio 迁移\n", modelCount)
	} else {
		migrateModelRatio(db, nowTs)
	}

	var groupCount int64
	db.Model(&GroupPrice{}).Count(&groupCount)
	if groupCount > 0 {
		fmt.Printf("group_prices 表已有 %d 条数据，跳过 GroupRatio 迁移\n", groupCount)
	} else {
		migrateGroupRatio(db, nowTs)
	}

	fmt.Println("迁移完成")
}

func migrateModelRatio(db *gorm.DB, nowTs int64) {
	var opt Option
	result := db.Where("`key` = ?", "ModelRatio").First(&opt)
	if result.Error != nil {
		fmt.Println("Option 表中无 ModelRatio 数据，跳过")
		return
	}

	modelRatio := make(map[string]float64)
	if err := json.Unmarshal([]byte(strings.ReplaceAll(opt.Value, "\n", "")), &modelRatio); err != nil {
		fmt.Printf("解析 ModelRatio JSON 失败: %v\n", err)
		return
	}

	var compOpt Option
	completionRatio := make(map[string]float64)
	if err := db.Where("`key` = ?", "CompletionRatio").First(&compOpt).Error; err == nil {
		json.Unmarshal([]byte(strings.ReplaceAll(compOpt.Value, "\n", "")), &completionRatio)
	}

	converted := 0
	skipped := 0
	for name, ratio := range modelRatio {
		if ratio < 0 {
			skipped++
			continue
		}

		inputPrice := ratio * 14.0
		compRatio := 1.0
		if cr, ok := completionRatio[name]; ok && cr > 0 {
			compRatio = cr
		}
		outputPrice := inputPrice * compRatio

		price := ModelPrice{
			ModelName:       name,
			InputPrice:      round6(inputPrice),
			OutputPrice:     round6(outputPrice),
			CachedPrice:     0,
			PerRequestPrice: 0,
			BillingType:     "token",
			Enabled:         true,
			CreatedAt:       nowTs,
			UpdatedAt:       nowTs,
		}
		if err := db.Create(&price).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
				skipped++
			} else {
				fmt.Printf("  插入 %s 失败: %v\n", name, err)
			}
		} else {
			converted++
		}
	}
	fmt.Printf("从 ModelRatio 迁移 %d 条模型定价，跳过 %d 条\n", converted, skipped)
}

func migrateGroupRatio(db *gorm.DB, nowTs int64) {
	var opt Option
	result := db.Where("`key` = ?", "GroupRatio").First(&opt)
	if result.Error != nil {
		fmt.Println("Option 表中无 GroupRatio 数据，插入默认分组")
		defaults := []GroupPrice{
			{GroupName: "default", ModelName: "", Discount: 1.0, CreatedAt: nowTs, UpdatedAt: nowTs},
			{GroupName: "vip", ModelName: "", Discount: 1.0, CreatedAt: nowTs, UpdatedAt: nowTs},
			{GroupName: "svip", ModelName: "", Discount: 1.0, CreatedAt: nowTs, UpdatedAt: nowTs},
		}
		for _, d := range defaults {
			db.Create(&d)
		}
		return
	}

	groupRatio := make(map[string]float64)
	if err := json.Unmarshal([]byte(strings.ReplaceAll(opt.Value, "\n", "")), &groupRatio); err != nil {
		fmt.Printf("解析 GroupRatio JSON 失败: %v\n", err)
		return
	}

	converted := 0
	for name, discount := range groupRatio {
		price := GroupPrice{
			GroupName: name,
			ModelName: "",
			Discount:  discount,
			CreatedAt: nowTs,
			UpdatedAt: nowTs,
		}
		if err := db.Create(&price).Error; err != nil {
			fmt.Printf("  插入分组 %s 失败: %v\n", name, err)
		} else {
			converted++
		}
	}
	fmt.Printf("从 GroupRatio 迁移 %d 条分组折扣\n", converted)
}

func round6(f float64) float64 {
	return float64(int(f*1000000+0.5)) / 1000000
}