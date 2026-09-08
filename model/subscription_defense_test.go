// subscription_defense_test.go 订阅写入防御测试
// 版本: v0.0.13
// 日期: 2026-09-08
// 作者: opencode
//
// 覆盖：
//   - UserPlan.Insert 拒绝 plan_id<=0
//   - ActivatePackageByOrder 拒绝 order.plan_id<=0
package model

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupDefenseTestDB 初始化内存 sqlite，建 user_plans + orders。
func setupDefenseTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&UserPlan{}, &Order{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// TestUserPlan_InsertRejectsPlanIdZero 验证 UserPlan.Insert 拒绝 plan_id<=0。
// UserPlan.Insert should refuse plan_id<=0 to keep the admin list
// from showing empty套餐 cells.
// 版本: v0.0.13
// 日期: 2026-09-08
func TestUserPlan_InsertRejectsPlanIdZero(t *testing.T) {
	DB = setupDefenseTestDB(t)

	bad := &UserPlan{UserId: 1, PlanId: 0, StartTime: 1, EndTime: 2, Status: 1}
	if err := bad.Insert(); err == nil {
		t.Fatal("期望 plan_id=0 被拒绝，却通过了")
	} else if !strings.Contains(err.Error(), "plan_id") {
		t.Fatalf("错误信息应包含 plan_id: %v", err)
	}

	neg := &UserPlan{UserId: 1, PlanId: -1, StartTime: 1, EndTime: 2, Status: 1}
	if err := neg.Insert(); err == nil {
		t.Fatal("期望 plan_id<0 被拒绝")
	}

	good := &UserPlan{UserId: 1, PlanId: 5, StartTime: 1, EndTime: 2, Status: 1}
	if err := good.Insert(); err != nil {
		t.Fatalf("正常 plan_id 应通过：%v", err)
	}
}

// TestActivatePackageByOrder_RejectsZeroPlanId 验证 ActivatePackageByOrder
// 拒绝 order.plan_id<=0 的订单(防止 user_plan 写入脏数据)。
// 版本: v0.0.13
// 日期: 2026-09-08
func TestActivatePackageByOrder_RejectsZeroPlanId(t *testing.T) {
	DB = setupDefenseTestDB(t)

	o := &Order{OrderNo: "TB0000", Type: OrderTypePlanSubscription, Status: OrderStatusPending, PlanId: 0}
	if err := DB.Create(o).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}
	err := ActivatePackageByOrder(o, OrderUpgradeModeStack)
	if err == nil {
		t.Fatal("期望 plan_id=0 的订单激活被拒绝")
	}
	if !strings.Contains(err.Error(), "plan_id") {
		t.Fatalf("错误信息应包含 plan_id: %v", err)
	}

	var count int64
	DB.Model(&UserPlan{}).Count(&count)
	if count != 0 {
		t.Fatalf("拒绝激活后不应有 user_plan 行，实际 %d", count)
	}
}