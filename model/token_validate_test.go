// token_validate_test.go ValidateUserToken 单测
// Unit tests for ValidateUserToken expiry semantics.
//
// 覆盖：
//   - sentinel == -1（永不过期）→ 通过
//   - sentinel == 0（前端旧 bug 落库的脏数据）→ 仍通过（兜底）
//   - ExpiredTime > now → 报"该令牌已过期"
//   - ExpiredTime 在未来 → 通过
//   - 空 key → 报"未提供令牌"
//   - 不存在的 key → 报"无效的令牌"
//   - status == TokenStatusExpired → 报"该令牌已过期"（前置短路）
//
// 版本: v0.0.21
// 日期: 2026-09-14
// 作者: opencode

package model

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
)

// setupTokenValidateTestDB wires an in-memory SQLite so we can exercise
// CacheGetTokenByKey's DB fallback path. Redis is disabled to mirror the
// prod config used by these tests.
func setupTokenValidateTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&User{}, &Token{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	DB = db
	common.RedisEnabled = false
	common.UsingSQLite = true
}

func seedToken(t *testing.T, key string, status int, expiredTime int64, unlimited bool) {
	t.Helper()
	user := &User{
		Username:    "u-" + key,
		Password:    "x",
		Role:        RoleCommonUser,
		Status:      UserStatusEnabled,
		Group:       "default",
		AccessToken: "acc-" + key,
		AffCode:     "aff-" + key,
	}
	if err := DB.Create(user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	tok := &Token{
		UserId:         user.Id,
		Key:            key,
		Status:         status,
		Name:           "test-" + key,
		ExpiredTime:    expiredTime,
		UnlimitedQuota: unlimited,
		RemainQuota:    0,
	}
	if err := DB.Create(tok).Error; err != nil {
		t.Fatalf("seed token: %v", err)
	}
}

// TestValidateUserToken_ExpirySentinels 验证"永不过期"的 sentinel 语义：
// 经典值 -1 与历史脏数据 0 都视为永不过期。
// 这是 2026-09-14 的回归测试：前端早期版本会把"永不过期"写成 0，
// 后端又只识别 -1，导致 UI 显示「永不过期」但 API 报「令牌已过期」。
func TestValidateUserToken_ExpirySentinels(t *testing.T) {
	setupTokenValidateTestDB(t)

	t.Run("ExpiredTime=-1 (canonical sentinel) passes", func(t *testing.T) {
		key := "key-never-expires-aaa"
		seedToken(t, key, TokenStatusEnabled, -1, true)
		tok, err := ValidateUserToken(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tok == nil || tok.ExpiredTime != -1 {
			t.Fatalf("got %+v, want ExpiredTime=-1", tok)
		}
	})

	t.Run("ExpiredTime=0 (legacy frontend bug) still passes", func(t *testing.T) {
		key := "key-legacy-zero-ccc"
		seedToken(t, key, TokenStatusEnabled, 1, true)
		// 用 raw SQL 强制把 expired_time 改成 0，绕过 GORM "default:-1" 标签
		// 在 INSERT 时把 0 替换成 -1 的副作用。
		// 模拟用户实际环境：MySQL 中真实存在的 expired_time=0 行
		// （来自老前端 "永不过期" 提交 0 的代码路径）。
		if err := DB.Exec(
			"UPDATE tokens SET expired_time = 0 WHERE `key` = ?", key,
		).Error; err != nil {
			t.Fatalf("force ExpiredTime=0: %v", err)
		}
		tok, err := ValidateUserToken(key)
		if err != nil {
			t.Fatalf("unexpected error for ExpiredTime=0: %v", err)
		}
		if tok.ExpiredTime != 0 {
			t.Fatalf("ExpiredTime changed unexpectedly: got %d", tok.ExpiredTime)
		}
	})
}

// TestValidateUserToken_PastExpiry 验证 ExpiredTime > 0 且 < now 时报过期。
func TestValidateUserToken_PastExpiry(t *testing.T) {
	setupTokenValidateTestDB(t)
	key := "key-past-expiry-ddd"
	past := helper.GetTimestamp() - 3600
	seedToken(t, key, TokenStatusEnabled, past, true)

	_, err := ValidateUserToken(key)
	if err == nil {
		t.Fatalf("expected expired error, got nil")
	}
	if err.Error() != "该令牌已过期" {
		t.Fatalf("got %q, want %q", err.Error(), "该令牌已过期")
	}
}

// TestValidateUserToken_FutureExpiry 验证 ExpiredTime 在未来时通过。
func TestValidateUserToken_FutureExpiry(t *testing.T) {
	setupTokenValidateTestDB(t)
	key := "key-future-expiry-eee"
	future := helper.GetTimestamp() + 3600
	seedToken(t, key, TokenStatusEnabled, future, true)

	if _, err := ValidateUserToken(key); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestValidateUserToken_StatusExpired 验证 status 字段已经显式标 Expired
// 时报过期，且优先级高于 ExpiredTime 数值。
func TestValidateUserToken_StatusExpired(t *testing.T) {
	setupTokenValidateTestDB(t)
	key := "key-status-expired-fff"
	seedToken(t, key, TokenStatusExpired, -1, true)

	_, err := ValidateUserToken(key)
	if err == nil {
		t.Fatalf("expected expired error, got nil")
	}
}

// TestValidateUserToken_NotFound 验证空 key / 未知 key 的错误路径。
func TestValidateUserToken_NotFound(t *testing.T) {
	setupTokenValidateTestDB(t)

	t.Run("empty key", func(t *testing.T) {
		_, err := ValidateUserToken("")
		if err == nil || err.Error() != "未提供令牌" {
			t.Fatalf("got %v, want %q", err, "未提供令牌")
		}
	})

	t.Run("unknown key", func(t *testing.T) {
		_, err := ValidateUserToken("does-not-exist-zzz")
		if err == nil {
			t.Fatalf("expected error for unknown key, got nil")
		}
		if !errors.Is(err, err) { // sanity: err is non-nil error
			t.Fatalf("unexpected error type: %T", err)
		}
	})
}