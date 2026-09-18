package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/model"
)

// TestMain wires the package-level DB to an in-memory SQLite so the handler
// can call into the real model layer (GetAllUsers / SearchUsers /
// GetUserSubscriptionBriefs) without touching the network or any external
// service. Redis is disabled to mirror middleware/auth_test.go.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to open in-memory sqlite: " + err.Error())
	}
	if err := db.AutoMigrate(&model.User{}, &model.Plan{}, &model.UserPlan{}, &model.ModelPrice{}, &model.GroupPrice{}, &model.Log{}); err != nil {
		panic("failed to migrate schema: " + err.Error())
	}
	model.DB = db
	model.LOG_DB = db
	common.RedisEnabled = false
	m.Run()
}

// seedUser creates a user with a known access_token so the tests can later
// assert whether that token was redacted from the response payload.
func seedUser(t *testing.T, username string, role int, accessToken string) *model.User {
	t.Helper()
	u := &model.User{
		Username:    username,
		Password:    "hashed-password",
		Role:        role,
		Status:      model.UserStatusEnabled,
		Group:       "default",
		AccessToken: accessToken,
		AffCode:     username + "-aff-code",
	}
	if err := model.DB.Create(u).Error; err != nil {
		t.Fatalf("seed user %s: %v", username, err)
	}
	return u
}

// runHandler invokes fn with a synthesized gin.Context that carries the given
// path/query string, then decodes the JSON body into out.
func runHandler(t *testing.T, fn gin.HandlerFunc, target string, out interface{}) {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, target, strings.NewReader(""))
	fn(c)
	if w.Code != http.StatusOK {
		t.Fatalf("handler returned status %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if err := json.Unmarshal(w.Body.Bytes(), out); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, w.Body.String())
	}
}

// response mirrors the gin.H envelope used by the real handlers.
type response struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []model.User `json:"data"`
}

// assertNoLeakedAccessToken fails the test if any user in data carries a
// non-empty access_token. Each call also reports how many users were checked.
func assertNoLeakedAccessToken(t *testing.T, users []model.User) {
	t.Helper()
	if len(users) == 0 {
		t.Fatalf("expected at least one user in response, got 0")
	}
	for _, u := range users {
		if u.AccessToken != "" {
			t.Fatalf("user %d (role=%d) leaked access_token=%q", u.Id, u.Role, u.AccessToken)
		}
	}
}

// TestGetAllUsers_OmitsAccessToken verifies that the controller redacts
// access_token from every user in the list response. This guards against the
// admin→root privilege-escalation path where a leaked root token could be
// replayed against RootAuth() endpoints.
func TestGetAllUsers_OmitsAccessToken(t *testing.T) {
	seedUser(t, "common1", model.RoleCommonUser, "common1-token")
	seedUser(t, "admin1", model.RoleAdminUser, "admin1-token")
	seedUser(t, "root1", model.RoleRootUser, "root1-token-secret")

	var resp response
	runHandler(t, GetAllUsers, "/api/user/?p=0", &resp)
	if !resp.Success {
		t.Fatalf("response success=false, message=%q", resp.Message)
	}
	assertNoLeakedAccessToken(t, resp.Data)
	if len(resp.Data) != 3 {
		t.Fatalf("got %d users, want 3", len(resp.Data))
	}
}

// TestSearchUsers_OmitsAccessToken verifies the same redaction for the
// keyword search endpoint.
func TestSearchUsers_OmitsAccessToken(t *testing.T) {
	seedUser(t, "alice", model.RoleCommonUser, "alice-token")
	seedUser(t, "alex", model.RoleAdminUser, "alex-token-secret")
	seedUser(t, "bob", model.RoleRootUser, "bob-token-top-secret")

	var resp response
	runHandler(t, SearchUsers, "/api/user/search?keyword=al", &resp)
	if !resp.Success {
		t.Fatalf("response success=false, message=%q", resp.Message)
	}
	assertNoLeakedAccessToken(t, resp.Data)
	for _, u := range resp.Data {
		if !strings.HasPrefix(u.Username, "al") {
			t.Fatalf("search returned non-matching user %q", u.Username)
		}
	}
}

// TestGetUserById_StillOmitsAccessToken is a regression guard for the
// single-user read path that was already protected by GetUserById(_, false).
// We assert the existing DAO-level omission so future refactors cannot
// silently regress that behavior.
func TestGetUserById_StillOmitsAccessToken(t *testing.T) {
	u := seedUser(t, "solo", model.RoleAdminUser, "solo-token")
	got, err := model.GetUserById(u.Id, false)
	if err != nil {
		t.Fatalf("GetUserById returned error: %v", err)
	}
	if got.AccessToken != "" {
		t.Fatalf("GetUserById(_, false) leaked access_token=%q", got.AccessToken)
	}
	// sanity: selectAll=true should still expose the token (used by
	// GenerateAccessToken). This is intentional and not under test here.
	full, err := model.GetUserById(u.Id, true)
	if err != nil {
		t.Fatalf("GetUserById(_, true) returned error: %v", err)
	}
	if full.AccessToken != "solo-token" {
		t.Fatalf("GetUserById(_, true) expected access_token=%q, got %q", "solo-token", full.AccessToken)
	}
	// And the wire-level user object that the controller serializes must
	// also strip the field. We re-use runHandler to drive the real handler
	// (GetUser reads ":id" from the path).
	// Note: GetUser enforces a role hierarchy check; seed caller as root
	// by injecting the role into the context via a tiny wrapper.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%d", u.Id), strings.NewReader(""))
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", u.Id)}}
	c.Set("role", model.RoleRootUser)
	GetUser(c)
	if w.Code != http.StatusOK {
		t.Fatalf("GetUser status=%d body=%s", w.Code, w.Body.String())
	}
	var single struct {
		Success bool       `json:"success"`
		Data    model.User `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &single); err != nil {
		t.Fatalf("decode single response: %v", err)
	}
	if single.Data.AccessToken != "" {
		t.Fatalf("GetUser wire response leaked access_token=%q", single.Data.AccessToken)
	}

	// runManage is a tiny harness for the POST /api/user/manage handler.
	// It accepts a JSON body, injects `role` / `username` into the context
	// (mimicking what middleware/auth.go would set after AdminAuth), and
	// decodes the standard {success, message, data} envelope.
	type batchResp struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	runManage := func(t *testing.T, body string, role int, caller string) batchResp {
		t.Helper()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/user/manage", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("role", role)
		c.Set("username", caller)
		ManageUser(c)
		if w.Code != http.StatusOK {
			t.Fatalf("ManageUser status=%d body=%s", w.Code, w.Body.String())
		}
		var r batchResp
		if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
			t.Fatalf("decode ManageUser response: %v; body=%s", err, w.Body.String())
		}
		return r
	}

	t.Run("batch_disable_non_root_forbidden", func(t *testing.T) {
		seedUser(t, "alice-rbac", model.RoleCommonUser, "alice-rbac-tk")
		seedUser(t, "bob-rbac", model.RoleCommonUser, "bob-rbac-tk")
		body := `{"action":"batch-disable","usernames":["alice-rbac","bob-rbac"]}`
		// admin (non-root) must be rejected
		r := runManage(t, body, model.RoleAdminUser, "admin-caller")
		if r.Success {
			t.Fatalf("non-root admin should be rejected; got success=true, data=%v", r.Data)
		}
		// both users must still be enabled
		for _, name := range []string{"alice-rbac", "bob-rbac"} {
			var u model.User
			if err := model.DB.Where("username = ?", name).First(&u).Error; err != nil {
				t.Fatalf("re-fetch %s: %v", name, err)
			}
			if u.Status != model.UserStatusEnabled {
				t.Fatalf("%s should still be enabled, got status=%d", name, u.Status)
			}
		}
	})

	t.Run("batch_delete_root_success_partial", func(t *testing.T) {
		// Seed fresh users for this subtest so state is independent.
		carol := seedUser(t, "carol-batch", model.RoleCommonUser, "carol-batch-tk")
		dave := seedUser(t, "dave-batch", model.RoleCommonUser, "dave-batch-tk")
		// existing root from earlier subtest must be skipped, not deleted.
		body := `{"action":"batch-delete","usernames":["carol-batch","dave-batch","root1"]}`
		r := runManage(t, body, model.RoleRootUser, "root1")
		if !r.Success {
			t.Fatalf("partial-success batch should be success=true; got message=%q data=%v", r.Message, r.Data)
		}
		if r.Data["succeeded_count"].(float64) != 2 {
			t.Fatalf("expected 2 succeeded, got %v", r.Data["succeeded_count"])
		}
		if r.Data["failed_count"].(float64) != 1 {
			t.Fatalf("expected 1 failed (root), got %v", r.Data["failed_count"])
		}
		failed, _ := r.Data["failed"].(map[string]interface{})
		if failed["root1"] == nil {
			t.Fatalf("expected root1 to appear in failed map, got %v", failed)
		}
		// carol-batch/dave-batch should be soft-deleted (status = UserStatusDeleted).
		// user.Delete() renames Username → deleted_<uuid>; query by Id instead.
		for _, u0 := range []*model.User{carol, dave} {
			var u model.User
			if err := model.DB.First(&u, u0.Id).Error; err != nil {
				t.Fatalf("re-fetch id=%d: %v", u0.Id, err)
			}
			if u.Status != model.UserStatusDeleted {
				t.Fatalf("id=%d should be soft-deleted, got status=%d", u0.Id, u.Status)
			}
		}
	})

	t.Run("batch_all_failed_returns_success_false", func(t *testing.T) {
		seedUser(t, "eve-batch", model.RoleCommonUser, "eve-batch-tk")
		// caller is root1, list contains root1 only — root self-skip fails the whole batch.
		bodySelf := `{"action":"batch-delete","usernames":["root1"]}`
		r := runManage(t, bodySelf, model.RoleRootUser, "root1")
		if r.Success {
			t.Fatalf("all-failed batch must return success=false; got message=%q", r.Message)
		}
		if r.Data["failed_count"].(float64) != 1 {
			t.Fatalf("expected 1 failed, got %v", r.Data["failed_count"])
		}
		// eve-batch untouched
		var u model.User
		if err := model.DB.Where("username = ?", "eve-batch").First(&u).Error; err != nil {
			t.Fatalf("re-fetch eve-batch: %v", err)
		}
		if u.Status != model.UserStatusEnabled {
			t.Fatalf("eve-batch should still be enabled, got status=%d", u.Status)
		}
	})

	t.Run("batch_empty_usernames_rejected", func(t *testing.T) {
		body := `{"action":"batch-disable","usernames":[]}`
		r := runManage(t, body, model.RoleRootUser, "root1")
		if r.Success {
			t.Fatalf("empty usernames must return success=false; got message=%q", r.Message)
		}
	})
}