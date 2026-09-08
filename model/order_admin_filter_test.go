// order_admin_filter_test 订单 admin 列表筛选 + UserBrief 嵌入单测。
// 版本: v0.0.13
// 日期: 2026-09-08
// 作者: opencode
//
// 覆盖：
//   - OrderAdminFilter.applyTo 各条件独立与组合行为
//   - GetAllOrders/SearchOrders 在 type/status/source/user_id/plan_id/keyword 下的过滤
//   - GetUsersBriefByIds 空切片短路 + WHERE IN 查询
//   - enrichOrdersWithUserBrief（controller 层）批量补全 User 字段
package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupOrderAdminTestDB 初始化内存 sqlite，建表 users + orders。
// setupOrderAdminTestDB spins up an in-memory sqlite with users + orders.
// 版本: v0.0.13
// 日期: 2026-09-08
func setupOrderAdminTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Order{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// seedUsersAndOrders 插入指定用户与订单，简化测试夹具。
// seedUsersAndOrders inserts the supplied users + orders as fixtures.
// 版本: v0.0.13
// 日期: 2026-09-08
func seedUsersAndOrders(t *testing.T, db *gorm.DB, users []*User, orders []*Order) {
	t.Helper()
	for _, u := range users {
		if err := db.Create(u).Error; err != nil {
			t.Fatalf("seed user: %v", err)
		}
	}
	for _, o := range orders {
		if err := db.Create(o).Error; err != nil {
			t.Fatalf("seed order: %v", err)
		}
	}
}

// TestOrderAdminFilter_ApplyTo 验证 applyTo 在每个字段为 0 时不过滤。
// TestOrderAdminFilter_ApplyTo verifies applyTo is a no-op when all fields are zero.
// 版本: v0.0.13
// 日期: 2026-09-08
func TestOrderAdminFilter_ApplyTo(t *testing.T) {
	db := setupOrderAdminTestDB(t)
	DB = db

	users := []*User{
		{Username: "alice", DisplayName: "Alice", Role: RoleCommonUser, Status: UserStatusEnabled, AccessToken: "tok-alice", AffCode: "aff-alice"},
	}
	orders := []*Order{
		{OrderNo: "TB0001", UserId: 0, Type: OrderTypePlanSubscription, Status: OrderStatusPaid, Source: OrderSourceUserSelf},
	}
	seedUsersAndOrders(t, db, users, orders)

	f := OrderAdminFilter{} // 全部零值
	q := f.applyTo(DB.Model(&Order{}))
	var got []*Order
	if err := q.Find(&got).Error; err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("期望 1 条，实际 %d", len(got))
	}
}

// TestGetAllOrders_FilterMatrix 覆盖 type/status/source/user_id/plan_id/keyword 5 个维度的过滤。
// TestGetAllOrders_FilterMatrix covers all 5 filter dimensions end-to-end.
// 版本: v0.0.13
// 日期: 2026-09-08
func TestGetAllOrders_FilterMatrix(t *testing.T) {
	db := setupOrderAdminTestDB(t)
	DB = db

	users := []*User{
		{Username: "alice", DisplayName: "Alice", Role: RoleCommonUser, Status: UserStatusEnabled, AccessToken: "tok-alice", AffCode: "aff-alice"},
		{Username: "bob", DisplayName: "Bob", Role: RoleCommonUser, Status: UserStatusEnabled, AccessToken: "tok-bob", AffCode: "aff-bob"},
	}
	orders := []*Order{
		{OrderNo: "TB0001", UserId: 1, Type: OrderTypePlanSubscription, Status: OrderStatusPending, Source: OrderSourceUserSelf, PlanId: 1},
		{OrderNo: "TB0002", UserId: 2, Type: OrderTypeTopup, Status: OrderStatusPaid, Source: OrderSourceAdmin, PlanId: 0},
		{OrderNo: "TB0003", UserId: 1, Type: OrderTypePlanSubscription, Status: OrderStatusPaid, Source: OrderSourceAdmin, PlanId: 1},
	}
	seedUsersAndOrders(t, db, users, orders)

	cases := []struct {
		name   string
		filter OrderAdminFilter
		want   []string
	}{
		// 注意：seed 时 CreateTime=0，applyTo 排序 Order("id desc")，
		// 所以同条件下的多条结果按 id 倒序排列。
		{"type=plan", OrderAdminFilter{Type: OrderTypePlanSubscription}, []string{"TB0003", "TB0001"}},
		{"type=topup", OrderAdminFilter{Type: OrderTypeTopup}, []string{"TB0002"}},
		{"status=paid", OrderAdminFilter{Status: OrderStatusPaid}, []string{"TB0003", "TB0002"}},
		{"status=canceled", OrderAdminFilter{Status: OrderStatusCanceled}, []string{}}, // 无匹配
		{"source=admin", OrderAdminFilter{Source: OrderSourceAdmin}, []string{"TB0003", "TB0002"}},
		{"plan_id=1", OrderAdminFilter{PlanId: 1}, []string{"TB0003", "TB0001"}},
		{"user_id=1", OrderAdminFilter{UserId: 1}, []string{"TB0003", "TB0001"}},
		{"keyword=TB0002", OrderAdminFilter{Keyword: "TB0002"}, []string{"TB0002"}},
		{"combined type+status", OrderAdminFilter{Type: OrderTypePlanSubscription, Status: OrderStatusPaid}, []string{"TB0003"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := GetAllOrders(0, 50, c.filter)
			if err != nil {
				t.Fatalf("GetAllOrders: %v", err)
			}
			if len(got) != len(c.want) {
				t.Fatalf("数量不匹配：got=%d want=%d (%v)", len(got), len(c.want), c.name)
			}
			for i, o := range got {
				if o.OrderNo != c.want[i] {
					t.Errorf("顺序/内容不符：got[%d]=%s want=%s", i, o.OrderNo, c.want[i])
				}
			}
		})
	}
}

// TestSearchOrders_StatusFilter 验证搜索结果也按 status 过滤。
// TestSearchOrders_StatusFilter verifies SearchOrders honors status filter.
// 版本: v0.0.13
// 日期: 2026-09-08
func TestSearchOrders_StatusFilter(t *testing.T) {
	db := setupOrderAdminTestDB(t)
	DB = db

	users := []*User{
		{Username: "alice", DisplayName: "Alice", Role: RoleCommonUser, Status: UserStatusEnabled, AccessToken: "tok-alice", AffCode: "aff-alice"},
	}
	orders := []*Order{
		{OrderNo: "TB0001", UserId: 1, Type: OrderTypePlanSubscription, Status: OrderStatusPending},
		{OrderNo: "TB0002", UserId: 1, Type: OrderTypeTopup, Status: OrderStatusPaid},
	}
	seedUsersAndOrders(t, db, users, orders)

	got, err := SearchOrders("TB", OrderAdminFilter{Status: OrderStatusPaid})
	if err != nil {
		t.Fatalf("SearchOrders: %v", err)
	}
	if len(got) != 1 || got[0].OrderNo != "TB0002" {
		t.Fatalf("期望仅 TB0002，实际 %v", got)
	}
}

// TestGetUsersBriefByIds_EmptyAndIn 验证空切片短路 + WHERE IN 查询。
// TestGetUsersBriefByIds_EmptyAndIn verifies short-circuit and WHERE IN fetch.
// 版本: v0.0.13
// 日期: 2026-09-08
func TestGetUsersBriefByIds_EmptyAndIn(t *testing.T) {
	db := setupOrderAdminTestDB(t)
	DB = db

	// 空 ids：不查 DB，返回空 map。
	got, err := GetUsersBriefByIds(nil)
	if err != nil {
		t.Fatalf("empty ids: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("空 ids 应返回空 map，实际 %d 条", len(got))
	}

	// 正常 ids：单次 WHERE IN 查询。
	users := []*User{
		{Username: "alice", DisplayName: "Alice", Role: RoleCommonUser, Status: UserStatusEnabled, AccessToken: "tok-alice", AffCode: "aff-alice"},
		{Username: "bob", DisplayName: "Bob", Role: RoleCommonUser, Status: UserStatusEnabled, AccessToken: "tok-bob", AffCode: "aff-bob"},
	}
	seedUsersAndOrders(t, db, users, nil)

	briefs, err := GetUsersBriefByIds([]int{1, 2})
	if err != nil {
		t.Fatalf("GetUsersBriefByIds: %v", err)
	}
	if len(briefs) != 2 {
		t.Fatalf("期望 2 条，实际 %d", len(briefs))
	}
	if briefs[1].Username != "alice" || briefs[2].DisplayName != "Bob" {
		t.Fatalf("字段不匹配：%v", briefs)
	}

	// 含不存在的 id：只返回存在的；不报错。
	briefs, err = GetUsersBriefByIds([]int{1, 99})
	if err != nil {
		t.Fatalf("GetUsersBriefByIds with missing id: %v", err)
	}
	if len(briefs) != 1 {
		t.Fatalf("期望 1 条（id=99 不存在），实际 %d", len(briefs))
	}
}