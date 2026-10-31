// validate_test.go 校验错误文案映射的单元测试
// Unit tests for human-readable validation error messages.
// 版本: v0.0.27
// 日期: 2026-10-31
// 作者: opencode
package common

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common/i18n"
)

func newTestContext(t *testing.T, lang string) *gin.Context {
	t.Helper()
	if err := i18n.Init(); err != nil {
		t.Fatalf("failed to init i18n: %v", err)
	}
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(i18n.ContextKey, lang)
	return c
}

// 模拟 model.User 的校验规则，验证文案映射覆盖 Password / Username 等字段。
type User struct {
	Username string `validate:"max=12"`
	Password string `validate:"min=8,max=20"`
}

func TestValidationMessage_PasswordTooShort_ZH(t *testing.T) {
	c := newTestContext(t, "zh-CN")
	err := Validate.Struct(&User{Username: "alice", Password: "123"})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	got := ValidationMessage(c, err)
	if got != "密码长度不能少于 8 位" {
		t.Fatalf("unexpected message: %q", got)
	}
	if strings.Contains(got, "validation for") || strings.Contains(got, "Key:") {
		t.Fatalf("raw validator error leaked to user: %q", got)
	}
}

func TestValidationMessage_PasswordTooShort_EN(t *testing.T) {
	c := newTestContext(t, "en")
	err := Validate.Struct(&User{Username: "alice", Password: "123"})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	got := ValidationMessage(c, err)
	if got != "password length must be at least 8 characters" {
		t.Fatalf("unexpected message: %q", got)
	}
}

func TestValidationMessage_MultipleFields(t *testing.T) {
	c := newTestContext(t, "zh-CN")
	err := Validate.Struct(&User{Username: "this-name-is-way-too-long", Password: "1"})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	got := ValidationMessage(c, err)
	if got != "用户名长度不能超过 12 位；密码长度不能少于 8 位" {
		t.Fatalf("unexpected message: %q", got)
	}
}

// 未登记的结构体/字段回落为通用提示，不暴露英文原始错误。
func TestValidationMessage_UnknownStructFallsBack(t *testing.T) {
	type Unknown struct {
		Field string `validate:"min=3"`
	}
	c := newTestContext(t, "zh-CN")
	err := Validate.Struct(&Unknown{Field: "x"})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	got := ValidationMessage(c, err)
	if got != "该字段长度不能少于 3 位" {
		t.Fatalf("unexpected message: %q", got)
	}
}

func TestValidationMessage_NotValidationError(t *testing.T) {
	c := newTestContext(t, "zh-CN")
	got := ValidationMessage(c, nil)
	if got != "无效的输入，请检查您的输入" {
		t.Fatalf("unexpected message: %q", got)
	}
}
