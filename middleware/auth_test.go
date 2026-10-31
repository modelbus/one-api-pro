package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/ctxkey"
	"github.com/modelbus/one-api-pro/model"
)

// TestMain wires the package-level DB to an in-memory SQLite before any test runs.
// It also disables Redis so middleware code paths skip the redis client.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to open in-memory sqlite: " + err.Error())
	}
	if err := db.AutoMigrate(&model.User{}, &model.Token{}); err != nil {
		panic("failed to migrate schema: " + err.Error())
	}
	model.DB = db
	common.RedisEnabled = false
	m.Run()
}

// seedUser inserts a user and an enabled token. If unlimited is true the token
// has unlimited quota and does not get blocked by quota checks.
// TokenAuth uses strings.Split(key, "-") to derive the (key, channelId) pair
// from the Authorization header, so token keys must not contain '-' or the
// suffix-channel feature will silently consume the wrong prefix.
const suffixKey = "suffixkeycccc"

func seedUser(t *testing.T, username string, role int, key string, unlimited bool) (int, int) {
	t.Helper()
	user := &model.User{
		Username:    username,
		Password:    "hashed-password",
		Role:        role,
		Status:      model.UserStatusEnabled,
		Group:       "default",
		AccessToken: username + "-access-token",
		AffCode:     username + "-aff-code",
	}
	if err := model.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to seed user %s: %v", username, err)
	}
	token := &model.Token{
		UserId:         user.Id,
		Key:            key,
		Status:         model.TokenStatusEnabled,
		Name:           "test-token",
		ExpiredTime:    -1,
		UnlimitedQuota: unlimited,
		RemainQuota:    0,
	}
	if err := model.DB.Create(token).Error; err != nil {
		t.Fatalf("failed to seed token for %s: %v", username, err)
	}
	return user.Id, token.Id
}

// TestTokenAuth_SuffixChannelPinning verifies that admin token-suffix channel
// pinning still works.
func TestTokenAuth_SuffixChannelPinning(t *testing.T) {
	seedUser(t, "root2", model.RoleAdminUser, suffixKey, true)
	r2 := gin.New()
	r2.POST("/v1/chat/completions", TokenAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"channel_id": c.GetString(ctxkey.SpecificChannelId)})
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(""))
	req.Header.Set("Authorization", "Bearer sk-"+suffixKey+"-42")
	w := httptest.NewRecorder()
	r2.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got, _ := body["channel_id"].(string); got != "42" {
		t.Fatalf("SpecificChannelId = %q, want %q (from suffix)", got, "42")
	}
}
