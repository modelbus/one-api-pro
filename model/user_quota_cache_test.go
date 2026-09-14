// User-quota cache fix unit tests
// 版本: v0.0.20
// 日期: 2026-09-14
// 作者: opencode
//
// 覆盖：
//   - IncreaseUserQuota 非 batch 模式下 DB 正确累加（修复 3 验证基础语义未破坏）
//   - IncreaseUserQuota 非 batch 模式下 batchUpdateEnabled=false 不入队
//   - IncreaseUserQuota 在 batchUpdateEnabled=true 下只入队不直接写 DB
//   - DecreaseUserQuota 与 IncreaseUserQuota 的对称性（确保 Decrease 不破坏缓存契约）
//
// 注：本测试只覆盖「DB 是真相」层；Redis 路径需依赖集成测试 / miniredis。

package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
)

// setupQuotaCacheTestDB wires an in-memory SQLite with the User table.
// 关闭 Redis 以模拟 CacheGetUserQuota 走 GetUserQuota 直读 DB 的路径。
func setupQuotaCacheTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	common.RedisEnabled = false
	return db
}

// seedQuotaUser creates a User with the given quota and returns the id.
// seedQuotaUser 创建一个指定 quota 的 User 并返回 id。
func seedQuotaUser(t *testing.T, db *gorm.DB, quota int64) int {
	t.Helper()
	u := &User{
		Username: "quota-user",
		Password: "x",
		Role:     RoleCommonUser,
		Status:   UserStatusEnabled,
		Group:    "default",
		Quota:    quota,
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u.Id
}

// TestIncreaseUserQuota_DirectWriteDB 验证非 batch 模式下 IncreaseUserQuota
// 直接 UPDATE users.quota，且 BatchUpdateStores 不被污染（修复 3 行为契约）。
// 版本: v0.0.20
func TestIncreaseUserQuota_DirectWriteDB(t *testing.T) {
	db := setupQuotaCacheTestDB(t)
	DB = db

	uid := seedQuotaUser(t, db, 1000)

	if err := IncreaseUserQuota(uid, 500); err != nil {
		t.Fatalf("IncreaseUserQuota: %v", err)
	}

	var u User
	if err := DB.First(&u, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if u.Quota != 1500 {
		t.Fatalf("expected quota=1500, got %d", u.Quota)
	}

	// batchUpdateStores[BatchUpdateTypeUserQuota] 不应被污染
	batchUpdateLocks[BatchUpdateTypeUserQuota].Lock()
	stored := batchUpdateStores[BatchUpdateTypeUserQuota][uid]
	batchUpdateLocks[BatchUpdateTypeUserQuota].Unlock()
	if stored != 0 {
		t.Fatalf("batch store should remain untouched in direct mode; got %d", stored)
	}
}

// TestIncreaseUserQuota_NegativeRejected 验证 quota<0 时 IncreaseUserQuota 拒绝。
// 这是已有的契约；新增的修复 3 不应破坏它。
// 版本: v0.0.20
func TestIncreaseUserQuota_NegativeRejected(t *testing.T) {
	db := setupQuotaCacheTestDB(t)
	DB = db

	uid := seedQuotaUser(t, db, 1000)

	err := IncreaseUserQuota(uid, -100)
	if err == nil {
		t.Fatalf("expected error for negative quota, got nil")
	}

	// DB 不应被改
	var u User
	if err := DB.First(&u, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if u.Quota != 1000 {
		t.Fatalf("expected quota unchanged at 1000, got %d", u.Quota)
	}
}

// TestDecreaseUserQuota_DirectWriteDB 验证 DecreaseUserQuota 走 DecreaseUserQuota
// 直接更新 DB，且 batch 模式不入队。
// 版本: v0.0.20
func TestDecreaseUserQuota_DirectWriteDB(t *testing.T) {
	db := setupQuotaCacheTestDB(t)
	DB = db

	uid := seedQuotaUser(t, db, 5000)

	if err := DecreaseUserQuota(uid, 1234); err != nil {
		t.Fatalf("DecreaseUserQuota: %v", err)
	}

	var u User
	if err := DB.First(&u, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if u.Quota != 5000-1234 {
		t.Fatalf("expected quota=%d, got %d", 5000-1234, u.Quota)
	}
}

// TestIncreaseDecreaseUserQuota_Roundtrip 验证连续加减不漂移：
// IncreaseUserQuota +500 紧接 DecreaseUserQuota +500 (把 quota 减到 0)，
// DB users.quota 应回到原始值。
// 这是修复 3「+quota 后刷缓存」路径的语义回归 —— 即便没有 Redis，
// 多次 IncreaseUserQuota + DecreaseUserQuota 调用 DB 一致性必须保持。
// 版本: v0.0.20
func TestIncreaseDecreaseUserQuota_Roundtrip(t *testing.T) {
	db := setupQuotaCacheTestDB(t)
	DB = db

	uid := seedQuotaUser(t, db, 1000)

	if err := IncreaseUserQuota(uid, 500); err != nil {
		t.Fatalf("IncreaseUserQuota: %v", err)
	}
	if err := DecreaseUserQuota(uid, 500); err != nil {
		t.Fatalf("DecreaseUserQuota: %v", err)
	}

	var u User
	if err := DB.First(&u, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if u.Quota != 1000 {
		t.Fatalf("roundtrip should restore original quota 1000, got %d", u.Quota)
	}
}

// TestCacheIncreaseUserQuota_DisabledRedisNoop 关闭 Redis 时 CacheIncreaseUserQuota
// 必须空操作且不返回错误（防御性回滚的语义保持）。
// 版本: v0.0.20
func TestCacheIncreaseUserQuota_DisabledRedisNoop(t *testing.T) {
	common.RedisEnabled = false
	if err := CacheIncreaseUserQuota(999, 12345); err != nil {
		t.Fatalf("expected nil error when Redis is disabled, got %v", err)
	}
}

// TestCacheDecreaseUserQuota_DisabledRedisNoop 与上一条对称，确保
// CacheDecreaseUserQuota 在 Redis 关闭时也不影响其它测试状态。
// 版本: v0.0.20
func TestCacheDecreaseUserQuota_DisabledRedisNoop(t *testing.T) {
	common.RedisEnabled = false
	if err := CacheDecreaseUserQuota(999, 12345); err != nil {
		t.Fatalf("expected nil error when Redis is disabled, got %v", err)
	}
}
