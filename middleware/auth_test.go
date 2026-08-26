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
const (
	commonUserKey = "usertokenkeyaaaa"
	adminKey      = "admintokenkeybbbb"
	suffixKey     = "suffixkeycccc"
)

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

// newProxyRouter builds the same route shape as router/relay.go for the
// channel-pinning endpoint, applying only TokenAuth. A terminal handler
// echoes back so the test can confirm TokenAuth let the request through.
func newProxyRouter() *gin.Engine {
	r := gin.New()
	r.Any("/v1/oneapi/proxy/:channelid/*target", TokenAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"channel_id": c.GetString(ctxkey.SpecificChannelId),
		})
	})
	return r
}

// doReq sends a request through the given router and returns the recorder.
func doReq(r *gin.Engine, path, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(""))
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TestTokenAuth_URLChannelPinning verifies that the URL-parameter channel
// pinning path is gated by admin role like the suffix path is.
func TestTokenAuth_URLChannelPinning(t *testing.T) {
	seedUser(t, "alice", model.RoleCommonUser, commonUserKey, true)
	seedUser(t, "root", model.RoleAdminUser, adminKey, true)

	r := newProxyRouter()

	t.Run("common user is forbidden from pinning by URL", func(t *testing.T) {
		w := doReq(r, "/v1/oneapi/proxy/1/v1/chat/completions", "Bearer sk-"+commonUserKey)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403; body=%s", w.Code, w.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		errObj, ok := body["error"].(map[string]any)
		if !ok {
			t.Fatalf("expected error object in body, got %v", body)
		}
		msg, _ := errObj["message"].(string)
		if !strings.HasPrefix(msg, "普通用户不支持指定渠道") {
			t.Fatalf("error message = %q, want prefix %q", msg, "普通用户不支持指定渠道")
		}
	})

	t.Run("admin can pin by URL", func(t *testing.T) {
		w := doReq(r, "/v1/oneapi/proxy/1/v1/chat/completions", "Bearer sk-"+adminKey)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got, _ := body["channel_id"].(string); got != "1" {
			t.Fatalf("SpecificChannelId = %q, want %q", got, "1")
		}
	})

	t.Run("admin suffix pinning still works", func(t *testing.T) {
		seedUser(t, "root2", model.RoleAdminUser, suffixKey, true)
		// Hit a non-:channelid route to confirm the suffix-channel path is
		// unaffected by the new admin gate on the URL :channelid path.
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
	})
}