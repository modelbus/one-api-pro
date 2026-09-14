// PreConsumeQuota unit tests
// 版本: v0.0.20
// 日期: 2026-09-14
// 作者: opencode
//
// 覆盖：
//   - getPreConsumedQuota 的纯函数公式（PreConsumedQuota + promptTokens + MaxTokens）× ratio
//   - preConsumeQuota 在 Redis 关闭 + SQLite 内存库下的三条分支：
//     (1) 余额不足 → 403 不动 DB
//     (2) 余额充足（> 100×preConsumedQuota）→ 免预扣，不动 DB
//     (3) 正常预扣 → DB user.quota 与 token.remain_quota 都扣减
//
// 注：本次修复涉及 Redis 与 DB 协同，但项目测试体系未引入 miniredis，
// 这里只覆盖 Redis 关闭时的直读 DB 路径。Redis 启用路径靠手测 / 集成测试。

package controller

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	dbmodel "github.com/modelbus/one-api-pro/model"
	"github.com/modelbus/one-api-pro/relay/meta"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

// setupPreConsumeTestDB wires an in-memory SQLite, migrates User & Token,
// and disables Redis so CacheGetUserQuota falls through to the direct DB read.
// Redis 关闭以便 CacheGetUserQuota 直接读 DB；这覆盖了修复后 Redis 路径的
// 「DB 是真相」核心逻辑，Redis 路径需依赖集成测试。
func setupPreConsumeTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&dbmodel.User{}, &dbmodel.Token{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	dbmodel.DB = db
	common.RedisEnabled = false
}

// seedUser creates a User row with the given quota and returns the id.
// seedUser 在 DB 中创建一个指定 quota 的 User，返回 id。
func seedUser(t *testing.T, quota int64) int {
	t.Helper()
	u := &dbmodel.User{
		Username: "tester",
		Password: "x",
		Role:     dbmodel.RoleCommonUser,
		Status:   dbmodel.UserStatusEnabled,
		Group:    "default",
		Quota:    quota,
	}
	if err := dbmodel.DB.Create(u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u.Id
}

// seedToken creates a Token row with the given remain_quota for the user.
// seedToken 为指定 user 创建一个指定 remain_quota 的 Token。
func seedToken(t *testing.T, userId int, remain int64) int {
	t.Helper()
	tk := &dbmodel.Token{
		UserId:      userId,
		Key:         "sk-test-" + randSuffix(),
		Status:      1,
		Name:        "test-token",
		RemainQuota: remain,
	}
	if err := dbmodel.DB.Create(tk).Error; err != nil {
		t.Fatalf("create token: %v", err)
	}
	return tk.Id
}

// randSuffix produces a unique-enough suffix so multiple token inserts in one
// test do not collide on the unique key index.
// randSuffix 为同一测试内多个 token 提供足够唯一的后缀，避免 unique key 冲突。
func randSuffix() string {
	// nanosecond + a counter-less salt is fine for in-memory SQLite
	return "uniq"
}

// TestGetPreConsumedQuota_BasicFormula 验证公式：preConsumed = (PreConsumedQuota + promptTokens + MaxTokens) * ratio。
// 版本: v0.0.20
func TestGetPreConsumedQuota_BasicFormula(t *testing.T) {
	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 100}
	// ratio=1 → preConsumed = 500 + 1000 + 100 = 1600
	got := getPreConsumedQuota(req, 1000, 1.0)
	if got != 1600 {
		t.Fatalf("ratio=1 want 1600, got %d", got)
	}

	// ratio=2.625 → 1600 * 2.625 = 4200
	got = getPreConsumedQuota(req, 1000, 2.625)
	if got != 4200 {
		t.Fatalf("ratio=2.625 want 4200, got %d", got)
	}

	// MaxTokens=0 不加入 → preConsumed = (500 + 1000) * 2.625 = 3937
	req2 := &relaymodel.GeneralOpenAIRequest{MaxTokens: 0}
	got = getPreConsumedQuota(req2, 1000, 2.625)
	if got != 3937 {
		t.Fatalf("ratio=2.625 no MaxTokens want 3937, got %d", got)
	}
}

// TestPreConsumeQuota_InsufficientBalance 余额 < preConsumedQuota 时返回 403，
// 且 DB 中 user.quota 与 token.remain_quota 不应变动。
// 版本: v0.0.20
func TestPreConsumeQuota_InsufficientBalance(t *testing.T) {
	setupPreConsumeTestDB(t)

	uid := seedUser(t, 1000) // 只够 1000 quota
	tid := seedToken(t, uid, 5_000_000)

	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 0}
	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}

	// preConsumed = (500 + 5000) * 2.625 = 14437 > 1000
	preConsumed, errObj := preConsumeQuota(context.Background(), req, 5000, 2.625, m)
	if errObj == nil {
		t.Fatalf("expected insufficient_user_quota error, got nil (preConsumed=%d)", preConsumed)
	}
	if errObj.Error.Code != "insufficient_user_quota" {
		t.Fatalf("expected code=insufficient_user_quota, got %s", errObj.Error.Code)
	}

	// DB 不应被修改
	var user dbmodel.User
	if err := dbmodel.DB.First(&user, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if user.Quota != 1000 {
		t.Fatalf("user.quota mutated: got %d, want 1000", user.Quota)
	}
	var token dbmodel.Token
	if err := dbmodel.DB.First(&token, "id = ?", tid).Error; err != nil {
		t.Fatalf("read token: %v", err)
	}
	if token.RemainQuota != 5_000_000 {
		t.Fatalf("token.remain_quota mutated: got %d, want 5_000_000", token.RemainQuota)
	}
}

// TestPreConsumeQuota_TrustedDoesNotMutate 是修复 1 的核心回归测试：
// 当 userQuota > 100*preConsumedQuota 时，preConsumeQuota 必须免预扣返回 0，
// 且 DB 的 user.quota 与 token.remain_quota 都不应变动。
//
// 修复前：CacheDecreaseUserQuota 先扣 Redis，再判定免预扣（局部 preConsumedQuota=0），
//         token 表不扣但 Redis 已扣（漂移源）。
// 修复后：守卫 2 提前 return，Redis 与 token 表都不动。
// 版本: v0.0.20
func TestPreConsumeQuota_TrustedDoesNotMutate(t *testing.T) {
	setupPreConsumeTestDB(t)

	uid := seedUser(t, 17_678_941) // admin5 的实际场景
	tid := seedToken(t, uid, 5_000_000)

	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 0}
	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}

	// preConsumed = (500 + 30340) * 2.625 = 80_955
	// 17_678_941 > 100 * 80_955 = 8_095_500 → 守卫 2 触发
	preConsumed, errObj := preConsumeQuota(context.Background(), req, 30340, 2.625, m)
	if errObj != nil {
		t.Fatalf("expected nil error, got %+v", errObj)
	}
	if preConsumed != 0 {
		t.Fatalf("expected preConsumed=0 (trusted), got %d", preConsumed)
	}

	// DB user.quota 维持不变
	var user dbmodel.User
	if err := dbmodel.DB.First(&user, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if user.Quota != 17_678_941 {
		t.Fatalf("user.quota mutated after trusted pre-consume: got %d, want 17_678_941", user.Quota)
	}
	// DB token.remain_quota 维持不变
	var token dbmodel.Token
	if err := dbmodel.DB.First(&token, "id = ?", tid).Error; err != nil {
		t.Fatalf("read token: %v", err)
	}
	if token.RemainQuota != 5_000_000 {
		t.Fatalf("token.remain_quota mutated after trusted pre-consume: got %d, want 5_000_000", token.RemainQuota)
	}
}

// TestPreConsumeQuota_NormalPreConsume 不满足 trusted 条件时正常预扣，
// DB 的 user.quota 与 token.remain_quota 都应被扣减 preConsumedQuota。
//
// 注意：PreConsumeTokenQuota 在 token 与 user 两张表都做扣减（model/token.go:273-279），
// 因此 pre-consume 阶段 user.quota 就已经减少 preConsumedQuota，而不是只在
// post-consume 时改 used_quota。这是与"先预扣 post 再结算"心智模型的差异。
// 版本: v0.0.20
func TestPreConsumeQuota_NormalPreConsume(t *testing.T) {
	setupPreConsumeTestDB(t)

	uid := seedUser(t, 200_000) // 用户余额足够但不满足 100x preConsumed
	tid := seedToken(t, uid, 5_000_000)

	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 0}
	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}

	// preConsumed = (500 + 1000) * 2.625 = 3937
	// 200_000 < 100 * 3937 = 393_700 → 守卫 2 不触发，正常预扣
	preConsumed, errObj := preConsumeQuota(context.Background(), req, 1000, 2.625, m)
	if errObj != nil {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if preConsumed != 3937 {
		t.Fatalf("expected preConsumed=3937, got %d", preConsumed)
	}

	// DB user.quota 在 pre-consume 阶段就被 DecreaseUserQuota 减去 preConsumed
	var user dbmodel.User
	if err := dbmodel.DB.First(&user, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	if user.Quota != 200_000-3937 {
		t.Fatalf("user.quota should be reduced by preConsumed after pre-consume; got %d, want %d", user.Quota, 200_000-3937)
	}

	// DB token.remain_quota 也减少 preConsumed
	var token dbmodel.Token
	if err := dbmodel.DB.First(&token, "id = ?", tid).Error; err != nil {
		t.Fatalf("read token: %v", err)
	}
	if token.RemainQuota != 5_000_000-3937 {
		t.Fatalf("token.remain_quota want %d, got %d", 5_000_000-3937, token.RemainQuota)
	}
}

// TestPreConsumeQuota_PlanIdShortCircuit 当 PlanId>0 时直接返回 0，不读 user/token。
// 版本: v0.0.20
func TestPreConsumeQuota_PlanIdShortCircuit(t *testing.T) {
	setupPreConsumeTestDB(t)

	uid := seedUser(t, 0) // quota=0
	tid := seedToken(t, uid, 0)

	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 100}
	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 99}

	preConsumed, errObj := preConsumeQuota(context.Background(), req, 99999, 99.0, m)
	if errObj != nil {
		t.Fatalf("plan path should not error, got %+v", errObj)
	}
	if preConsumed != 0 {
		t.Fatalf("expected preConsumed=0 for plan path, got %d", preConsumed)
	}
}
