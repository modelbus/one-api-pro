// migrate_topup_presets 修正 topup.presets 中旧口径的 bonus_quota
// Fix legacy-format bonus_quota values inside topup.presets
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode
//
// 背景：
//
//	v0.0.10 引入充值快捷金额时，bonus_quota 使用「元」口径（与 amount 同量纲）；
//	v0.0.24 将额度基准改为微元（1 元 = 1_000_000），但未迁移已落库的旧值。
//	结果 {amount:10, bonus_quota:10} 被按微元解读成「充 10 元到账 ¥0.00001」，
//	前端校验会直接报「到账金额不能低于支付金额」，设置页无法保存。
//
// 判据：
//
//	model.SaveTopupSettings 保证 bonus_quota >= round(amount × 1e6)（允许赠送、
//	不允许缩水），因此任何 bonus_quota < round(amount × 1e6) 的行必然是旧口径，
//	按 1:1 等值（充多少到账多少）修正。
//
// 用法：
//
//	SQL_DSN=... go run ./cmd/migrate_topup_presets          # 干跑，仅打印差异
//	SQL_DSN=... go run ./cmd/migrate_topup_presets --apply  # 实际写回
//
// SQL_DSN 留空时回退到当前目录的 one-api-pro.db（SQLite）。
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// quotaPerUnit 与 common/config.QuotaPerUnit 保持一致：1 元 = 1_000_000 微元。
// Kept in sync with common/config.QuotaPerUnit: 1 CNY = 1_000_000 micro-quota.
const quotaPerUnit = 1_000_000.0

// settingKeyTopupPresets 与 model.SystemSettingKeyTopupPresets 保持一致。
const settingKeyTopupPresets = "topup.presets"

// topupPreset 与 model.TopupPreset 的持久化结构一致。
type topupPreset struct {
	Amount     float64 `json:"amount"`
	BonusQuota int64   `json:"bonus_quota"`
}

func main() {
	apply := false
	for _, arg := range os.Args[1:] {
		if arg == "--apply" {
			apply = true
		}
	}

	dsn := os.Getenv("SQL_DSN")
	if dsn == "" {
		path := os.Getenv("SQLITE_PATH")
		if path == "" {
			path = "one-api-pro.db"
		}
		fmt.Printf("未设置 SQL_DSN，回退到 SQLite 文件：%s\n", path)
	}

	db, err := openDB(dsn)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}

	var value string
	err = db.Raw("SELECT value FROM system_settings WHERE key = ?", settingKeyTopupPresets).Scan(&value).Error
	if err != nil {
		fmt.Printf("读取 system_settings 失败: %v\n", err)
		os.Exit(1)
	}
	if value == "" {
		fmt.Printf("未找到 %s 配置行，无需迁移\n", settingKeyTopupPresets)
		return
	}

	var presets []topupPreset
	if err := json.Unmarshal([]byte(value), &presets); err != nil {
		fmt.Printf("解析 %s JSON 失败: %v\n原始值: %s\n", settingKeyTopupPresets, err, value)
		os.Exit(1)
	}

	changed := 0
	for i := range presets {
		p := &presets[i]
		// 支付金额折算出的最低到账额度（微元）。
		want := int64(math.Round(p.Amount * quotaPerUnit))
		if p.Amount <= 0 {
			fmt.Printf("  第 %d 行 amount=%.6f 非正数，跳过\n", i+1, p.Amount)
			continue
		}
		if p.BonusQuota >= want {
			continue
		}
		fmt.Printf("  第 %d 行：支付 %.6f 元，到账 %d → %d（= ¥%.6f）\n",
			i+1, p.Amount, p.BonusQuota, want, float64(want)/quotaPerUnit)
		p.BonusQuota = want
		changed++
	}

	if changed == 0 {
		fmt.Printf("%s 中 %d 行均为微元口径，无需迁移\n", settingKeyTopupPresets, len(presets))
		return
	}

	if !apply {
		fmt.Printf("\n干跑完成：%d/%d 行需要修正。确认无误后加 --apply 实际写回。\n", changed, len(presets))
		return
	}

	newValue, err := json.Marshal(presets)
	if err != nil {
		fmt.Printf("序列化失败: %v\n", err)
		os.Exit(1)
	}
	if err := db.Exec("UPDATE system_settings SET value = ? WHERE key = ?", string(newValue), settingKeyTopupPresets).Error; err != nil {
		fmt.Printf("写回失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n已写回 %s：修正 %d/%d 行\n", settingKeyTopupPresets, changed, len(presets))
}

// openDB 按主程序 model.chooseDB 的同一规则选择驱动：
//   - `postgres://` 前缀 → PostgreSQL；
//   - SQL_DSN 非空（且非 postgres）→ MySQL；
//   - SQL_DSN 为空 → SQLite（路径默认 one-api-pro.db，可用 SQLITE_PATH 覆盖）。
//
// 与主程序保持一致可避免「服务连 MySQL、脚本却写 SQLite」这类误操作。
// Pick the driver with the same rule as model.chooseDB so the script cannot write to a
// different database than the running service.
func openDB(dsn string) (*gorm.DB, error) {
	switch {
	case strings.HasPrefix(dsn, "postgres://"), strings.HasPrefix(dsn, "postgresql://"):
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case dsn != "":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		path := os.Getenv("SQLITE_PATH")
		if path == "" {
			path = "one-api-pro.db"
		}
		return gorm.Open(sqlite.Open(path), &gorm.Config{})
	}
}
