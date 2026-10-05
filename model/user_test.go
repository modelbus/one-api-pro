package model

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
)

// setupUserTestDB wires an in-memory SQLite so we can exercise the User
// BeforeCreate / BeforeUpdate hooks without touching the network or any
// external service. Redis is disabled to mirror controller/user_test.go.
func setupUserTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open in-memory sqlite: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("migrate User: %v", err)
	}
	DB = db
	common.RedisEnabled = false
}

// TestUser_BeforeCreate_AutoFillsTimestamps 验证直接 DB.Create(user) 会自动写入
// CreatedAt / UpdatedAt，避免历史 bug 导致新用户的两个时间戳都是 0。
// 版本: v0.0.12
// 日期: 2026-09-07
func TestUser_BeforeCreate_AutoFillsTimestamps(t *testing.T) {
	setupUserTestDB(t)

	u := &User{
		Username:    "auto-ts-user",
		Password:    "hashed-pwd",
		Role:        RoleCommonUser,
		Status:      UserStatusEnabled,
		Group:       "default",
		AccessToken: "auto-ts-token",
		AffCode:     "auto-aff",
	}
	if err := DB.Create(u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if u.CreatedAt == 0 {
		t.Fatalf("CreatedAt was not populated by BeforeCreate; got 0")
	}
	if u.UpdatedAt == 0 {
		t.Fatalf("UpdatedAt was not populated by BeforeCreate; got 0")
	}
	if u.CreatedAt != u.UpdatedAt {
		t.Fatalf("expected CreatedAt == UpdatedAt on insert, got %d vs %d", u.CreatedAt, u.UpdatedAt)
	}
	// Sanity: re-read from DB to confirm the column was actually persisted.
	var stored User
	if err := DB.First(&stored, "id = ?", u.Id).Error; err != nil {
		t.Fatalf("readback: %v", err)
	}
	if stored.CreatedAt == 0 || stored.UpdatedAt == 0 {
		t.Fatalf("DB row has zero timestamps: created_at=%d updated_at=%d", stored.CreatedAt, stored.UpdatedAt)
	}
}

// TestUser_BeforeUpdate_RefreshesUpdatedAt 验证任何 Updates(...) 路径都会刷新
// UpdatedAt，覆盖 UpdateUser / ManageUser / EmailBind / GenerateAccessToken 等。
// 版本: v0.0.12
// 日期: 2026-09-07
func TestUser_BeforeUpdate_RefreshesUpdatedAt(t *testing.T) {
	setupUserTestDB(t)
	u := &User{
		Username:    "updatable",
		Password:    "hashed-pwd",
		Role:        RoleCommonUser,
		Status:      UserStatusEnabled,
		Group:       "default",
		AccessToken: "updatable-token",
		AffCode:     "updatable-aff",
	}
	if err := DB.Create(u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	originalCreatedAt := u.CreatedAt
	originalUpdatedAt := u.UpdatedAt

	// Sleep long enough so the new Unix-second timestamp is strictly greater.
	time.Sleep(1100 * time.Millisecond)

	u.DisplayName = "new-display"
	if err := DB.Model(u).Updates(u).Error; err != nil {
		t.Fatalf("update user: %v", err)
	}
	if u.UpdatedAt <= originalUpdatedAt {
		t.Fatalf("UpdatedAt was not refreshed; before=%d after=%d", originalUpdatedAt, u.UpdatedAt)
	}
	if u.CreatedAt != originalCreatedAt {
		t.Fatalf("CreatedAt should be immutable on update; before=%d after=%d", originalCreatedAt, u.CreatedAt)
	}

	// Readback sanity.
	var stored User
	if err := DB.First(&stored, "id = ?", u.Id).Error; err != nil {
		t.Fatalf("readback: %v", err)
	}
	if stored.UpdatedAt <= originalUpdatedAt {
		t.Fatalf("DB row UpdatedAt not refreshed; got %d, want > %d", stored.UpdatedAt, originalUpdatedAt)
	}
}

// TestUpdateUserQuotaColumn_AllowsZero 验证按列更新能把额度写成 0。
//
// 背景：User.Update 使用结构体形式的 Updates(user)，GORM 会跳过零值字段，
// 管理端「把额度清零」会接口成功但额度不变；本函数是修复路径。
// 版本: v0.0.25
// 日期: 2026-10-05
func TestUpdateUserQuotaColumn_AllowsZero(t *testing.T) {
	setupUserTestDB(t)
	u := &User{
		Username:    "quota-zero-user",
		Password:    "hashed-pwd",
		Role:        RoleCommonUser,
		Status:      UserStatusEnabled,
		Group:       "default",
		AccessToken: "quota-zero-token",
		AffCode:     "quota-zero-aff",
		Quota:       1_000_000,
	}
	if err := DB.Create(u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := UpdateUserQuotaColumn(u.Id, 0); err != nil {
		t.Fatalf("UpdateUserQuotaColumn: %v", err)
	}
	var stored User
	if err := DB.First(&stored, "id = ?", u.Id).Error; err != nil {
		t.Fatalf("readback: %v", err)
	}
	if stored.Quota != 0 {
		t.Fatalf("quota = %d，期望 0", stored.Quota)
	}

	// 非零值同样可写入。
	if err := UpdateUserQuotaColumn(u.Id, 2_500_000); err != nil {
		t.Fatalf("UpdateUserQuotaColumn(non-zero): %v", err)
	}
	if err := DB.First(&stored, "id = ?", u.Id).Error; err != nil {
		t.Fatalf("readback#2: %v", err)
	}
	if stored.Quota != 2_500_000 {
		t.Fatalf("quota = %d，期望 2500000", stored.Quota)
	}

	// id 为空时拒绝，避免误更新整表。
	if err := UpdateUserQuotaColumn(0, 10); err == nil {
		t.Fatal("userId=0 应返回错误")
	}
}

// _ ensures the helper import is referenced even if future tests don't use it.
var _ = helper.GetTimestamp
