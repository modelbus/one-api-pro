// plan_usage_sqli_test.go 套餐用量累加的 SQL 注入回归测试
// Regression tests ensuring IncrementPlanUsage binds every value as a parameter.
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
//
// 覆盖：
//   - MySQL 方言下生成的 SQL 必须使用 ? 占位符，不得内联 model 值
//   - SQLite 下恶意 model 名按字面量落库，表结构与累加语义不受影响

package model

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// captureLogger 截获 GORM Trace 中的 SQL，供注入回归断言使用。
type captureLogger struct {
	logger.Interface
	sql string
}

// Trace 记录最近一条 SQL；其余日志行为沿用被嵌入的默认 logger。
// 版本: v0.0.25
// 日期: 2026-10-05
func (l *captureLogger) Trace(_ context.Context, _ time.Time, fc func() (string, int64), _ error) {
	l.sql, _ = fc()
}

// rawSQLDialector 覆盖 Explain，使 Trace 拿到带 ? 占位符的原始 SQL 而非内联参数后的 SQL。
type rawSQLDialector struct {
	gorm.Dialector
}

// Explain 直接返回原始 SQL，避免 GORM 在日志中把参数内联回语句。
// 版本: v0.0.25
// 日期: 2026-10-05
func (d rawSQLDialector) Explain(sql string, _ ...interface{}) string {
	return sql
}

// TestIncrementPlanUsage_MySQLParameterized 断言 MySQL 方言生成的 upsert 语句参数化。
// 版本: v0.0.25
// 日期: 2026-10-05
func TestIncrementPlanUsage_MySQLParameterized(t *testing.T) {
	dialector := mysql.New(mysql.Config{
		DSN:                       "user:pass@tcp(127.0.0.1:3306)/oneapi?charset=utf8mb4&parseTime=True",
		SkipInitializeWithVersion: true,
	})
	cap := &captureLogger{Interface: logger.Default.LogMode(logger.Silent)}
	db, err := gorm.Open(rawSQLDialector{Dialector: dialector}, &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
		Logger:                 cap,
	})
	if err != nil {
		t.Fatalf("open mysql dry-run: %v", err)
	}

	oldDB := DB
	DB = db
	t.Cleanup(func() { DB = oldDB })

	payload := "m' ON DUPLICATE KEY UPDATE requests = 999999 -- "
	if err := IncrementPlanUsage(1, payload, WindowTypePeriod, 0, 1, 2, 3, 4); err != nil {
		t.Fatalf("IncrementPlanUsage: %v", err)
	}

	if !strings.Contains(cap.sql, "ON DUPLICATE KEY UPDATE") {
		t.Fatalf("expected MySQL upsert clause, got: %s", cap.sql)
	}
	if strings.Contains(cap.sql, payload) || strings.Contains(cap.sql, "999999") {
		t.Fatalf("payload was inlined into SQL: %s", cap.sql)
	}
	if !strings.Contains(cap.sql, "?") {
		t.Fatalf("expected placeholders in SQL, got: %s", cap.sql)
	}
}

// TestIncrementPlanUsage_SQLiteLiteralAndAccumulate 断言恶意模型名按字面量落库且累加语义正确。
// 版本: v0.0.25
// 日期: 2026-10-05
func TestIncrementPlanUsage_SQLiteLiteralAndAccumulate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&PlanUsage{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	oldDB := DB
	DB = db
	t.Cleanup(func() { DB = oldDB })

	payload := "m'); DROP TABLE plan_usages; -- "
	if err := IncrementPlanUsage(1, payload, WindowTypePeriod, 0, 1, 10, 20, 5); err != nil {
		t.Fatalf("first increment: %v", err)
	}
	if err := IncrementPlanUsage(1, payload, WindowTypePeriod, 0, 2, 1, 2, 3); err != nil {
		t.Fatalf("second increment: %v", err)
	}

	if !db.Migrator().HasTable(&PlanUsage{}) {
		t.Fatal("plan_usages table was dropped by injected payload")
	}

	var pu PlanUsage
	if err := db.Where("user_plan_id = ? AND model = ? AND window_type = ? AND window_index = ?",
		1, payload, WindowTypePeriod, 0).First(&pu).Error; err != nil {
		t.Fatalf("query plan usage: %v", err)
	}
	if pu.Model != payload {
		t.Fatalf("model stored as %q, want literal %q", pu.Model, payload)
	}
	if pu.Requests != 3 || pu.PromptTokens != 11 || pu.CompletionTokens != 22 || pu.CachedTokens != 8 {
		t.Fatalf("unexpected accumulated usage: requests=%d prompt=%d completion=%d cached=%d",
			pu.Requests, pu.PromptTokens, pu.CompletionTokens, pu.CachedTokens)
	}
}
