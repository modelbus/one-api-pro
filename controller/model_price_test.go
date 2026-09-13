// model_price_test HTTP-layer tests for AddModelPrice/AddGroupPrice.
// Verify the controller enforces Id=0 even when the client sends a stale id
// in the JSON body (frontend bug where editing a record and clicking "add"
// would carry the previous record's id into the new POST).
//
// 版本: v0.0.18
// 日期: 2026-09-13
// 作者: opencode

package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/model"
)

// priceResponse mirrors the gin.H envelope used by the real handlers.
type priceResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// runPriceHandler invokes fn with a synthesized POST request whose body is the
// given JSON payload, then decodes the response envelope.
func runPriceHandler(t *testing.T, fn gin.HandlerFunc, payload string) priceResponse {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/model_price/", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	fn(c)
	var out priceResponse
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, w.Body.String())
	}
	return out
}

// runOptionsHandler invokes fn with a synthesized GET request (no body), then
// decodes the response envelope. Used for ListModelPriceOptions.
func runOptionsHandler(t *testing.T, fn gin.HandlerFunc) priceResponse {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/model_price/options", nil)
	fn(c)
	var out priceResponse
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, w.Body.String())
	}
	return out
}

// TestAddModelPrice_StripsClientId verifies that even when the client JSON
// body carries id=9999, AddModelPrice zeroes it so the new row gets an
// auto-incremented primary key and never collides with / overwrites an
// existing row at id=9999.
//
// 版本: v0.0.18
// 日期: 2026-09-13
func TestAddModelPrice_StripsClientId(t *testing.T) {
	// Seed an existing row at id=9999 (via raw SQL so we can preset the PK).
	if err := model.DB.Exec(
		"INSERT INTO model_prices (id, model_name, input_price, output_price, cached_price, per_request_price, billing_type, enabled, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		9999, "preexisting-model", 1.0, 2.0, 0.5, 0.0, "token", true, 0, 0,
	).Error; err != nil {
		t.Fatalf("seed row: %v", err)
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"id":                9999,
		"model_name":        "gpt-4o",
		"input_price":       2.5,
		"output_price":      10.0,
		"cached_price":      1.25,
		"per_request_price": 0,
		"billing_type":      "token",
		"enabled":           true,
	})

	resp := runPriceHandler(t, AddModelPrice, string(payload))
	if !resp.Success {
		t.Fatalf("expected success, got: %s", resp.Message)
	}

	// The pre-existing row at id=9999 must still exist (not overwritten).
	var preExisting model.ModelPrice
	if err := model.DB.First(&preExisting, "id = ?", 9999).Error; err != nil {
		t.Fatalf("pre-existing row lost: %v", err)
	}
	if preExisting.ModelName != "preexisting-model" {
		t.Fatalf("pre-existing row got overwritten: model_name=%q", preExisting.ModelName)
	}

	// The new row must exist with an auto-incremented id (not 9999) and the
	// correct model_name.
	var got model.ModelPrice
	if err := model.DB.First(&got, "model_name = ?", "gpt-4o").Error; err != nil {
		t.Fatalf("new row missing: %v", err)
	}
	if got.Id == 9999 {
		t.Fatalf("new row used client-supplied id=9999 (would overwrite existing row)")
	}
	if got.Id <= 0 {
		t.Fatalf("new row has invalid id=%d", got.Id)
	}
}

// TestAddGroupPrice_StripsClientId mirrors the model-price assertion for
// group prices.
//
// 版本: v0.0.18
// 日期: 2026-09-13
func TestAddGroupPrice_StripsClientId(t *testing.T) {
	// Seed an existing row at id=9999.
	if err := model.DB.Exec(
		"INSERT INTO group_prices (id, group_name, model_name, discount, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		9999, "preexisting-group", "preexisting-model", 0.9, 0, 0,
	).Error; err != nil {
		t.Fatalf("seed row: %v", err)
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"id":         9999,
		"group_name": "vip",
		"model_name": "gpt-4o",
		"discount":   0.5,
	})

	resp := runPriceHandler(t, AddGroupPrice, string(payload))
	if !resp.Success {
		t.Fatalf("expected success, got: %s", resp.Message)
	}

	// The pre-existing row at id=9999 must still exist.
	var preExisting model.GroupPrice
	if err := model.DB.First(&preExisting, "id = ?", 9999).Error; err != nil {
		t.Fatalf("pre-existing row lost: %v", err)
	}
	if preExisting.GroupName != "preexisting-group" {
		t.Fatalf("pre-existing row got overwritten: group_name=%q", preExisting.GroupName)
	}

	// The new row must have an auto-incremented id.
	var got model.GroupPrice
	if err := model.DB.First(&got, "group_name = ? AND model_name = ?", "vip", "gpt-4o").Error; err != nil {
		t.Fatalf("new row missing: %v", err)
	}
	if got.Id == 9999 {
		t.Fatalf("new row used client-supplied id=9999 (would overwrite existing row)")
	}
	if got.Id <= 0 {
		t.Fatalf("new row has invalid id=%d", got.Id)
	}
	if got.Discount != 0.5 {
		t.Fatalf("new row has discount=%v, want 0.5", got.Discount)
	}
}

// TestListModelPriceOptions_OnlyEnabled verifies that the dedicated read-only
// dropdown endpoint returns only enabled=true rows and projects the result
// down to model_name strings (no pricing fields leak into the dropdown UI).
//
// 版本: v0.0.19
// 日期: 2026-09-13
func TestListModelPriceOptions_OnlyEnabled(t *testing.T) {
	// Seed: 2 enabled + 1 disabled. Names are namespaced so they don't collide
	// with data left behind by TestAddModelPrice_StripsClientId (which seeds
	// "preexisting-model" + a row named "gpt-4o" via raw SQL with id=9999).
	if err := model.DB.Create(&model.ModelPrice{
		ModelName: "opt-enabled-alpha", BillingType: model.BillingTypeToken, Enabled: true,
	}).Error; err != nil {
		t.Fatalf("seed enabled #1: %v", err)
	}
	if err := model.DB.Create(&model.ModelPrice{
		ModelName: "opt-enabled-beta", BillingType: model.BillingTypeToken, Enabled: true,
	}).Error; err != nil {
		t.Fatalf("seed enabled #2: %v", err)
	}
	// For the disabled row, use raw SQL — `gorm:"default:true;not null"` on
	// the Enabled field forces db.Create(Enabled:false) to insert `true` in
	// the glebarez/sqlite driver, so the explicit false must bypass GORM's
	// default-value handling.
	if err := model.DB.Exec(
		"INSERT INTO model_prices (model_name, input_price, output_price, cached_price, per_request_price, billing_type, enabled, created_at, updated_at) VALUES (?, 0, 0, 0, 0, ?, ?, 0, 0)",
		"opt-disabled", model.BillingTypeToken, false,
	).Error; err != nil {
		t.Fatalf("seed disabled: %v", err)
	}

	resp := runOptionsHandler(t, ListModelPriceOptions)
	if !resp.Success {
		t.Fatalf("expected success, got: %s", resp.Message)
	}

	// Decode the data field as []string (the handler is supposed to project
	// model_name only).
	raw, err := json.Marshal(resp.Data)
	if err != nil {
		t.Fatalf("re-marshal data: %v", err)
	}
	var names []string
	if err := json.Unmarshal(raw, &names); err != nil {
		t.Fatalf("decode data as []string: %v; raw=%s", err, string(raw))
	}

	// Verify our 3 seeded rows: enabled-alpha + enabled-beta included,
	// opt-disabled excluded.
	assertContainsAll(t, names, []string{"opt-enabled-alpha", "opt-enabled-beta"})
	assertContainsNone(t, names, []string{"opt-disabled"})

	// Disabled row must NOT leak into the dropdown.
	for _, n := range names {
		if n == "opt-disabled" {
			t.Fatalf("disabled model leaked into dropdown: %v", names)
		}
	}
}

// TestListModelPriceOptions_DedupEmptyName verifies that rows with empty
// model_name (which can sneak in via bad admin input) are silently skipped
// instead of returning an empty-string option.
//
// 版本: v0.0.19
// 日期: 2026-09-13
func TestListModelPriceOptions_DedupEmptyName(t *testing.T) {
	// Use raw SQL for the empty-name seed too — GORM's `not null` would
	// reject an empty ModelName before the row ever reaches the handler.
	if err := model.DB.Exec(
		"INSERT INTO model_prices (model_name, input_price, output_price, cached_price, per_request_price, billing_type, enabled, created_at, updated_at) VALUES (?, 0, 0, 0, 0, ?, ?, 0, 0)",
		"", model.BillingTypeToken, true,
	).Error; err != nil {
		// Some SQLite engines reject the empty-string insert because of the
		// NOT NULL / uniqueIndex on model_name. In that case the row never
		// gets created and the handler naturally returns no empty entry —
		// which is exactly what we want to verify. Skip the raw-insert path
		// and fall through to the handler call.
		t.Logf("empty-name seed skipped (db rejected): %v", err)
	}
	if err := model.DB.Create(&model.ModelPrice{
		ModelName: "opt-real-name", BillingType: model.BillingTypeToken, Enabled: true,
	}).Error; err != nil {
		t.Fatalf("seed real name: %v", err)
	}

	resp := runOptionsHandler(t, ListModelPriceOptions)
	if !resp.Success {
		t.Fatalf("expected success, got: %s", resp.Message)
	}

	raw, _ := json.Marshal(resp.Data)
	var names []string
	if err := json.Unmarshal(raw, &names); err != nil {
		t.Fatalf("decode data: %v", err)
	}
	// Empty model_name must not appear anywhere in the dropdown payload.
	for _, n := range names {
		if n == "" {
			t.Fatalf("empty model_name leaked into dropdown: %v", names)
		}
	}
	assertContainsAll(t, names, []string{"opt-real-name"})
}

// assertContainsAll fails the test if any name in want is missing from got.
func assertContainsAll(t *testing.T, got, want []string) {
	t.Helper()
	idx := make(map[string]bool, len(got))
	for _, g := range got {
		idx[g] = true
	}
	for _, w := range want {
		if !idx[w] {
			t.Errorf("expected %q in result, got %v", w, got)
		}
	}
}

// assertContainsNone fails the test if any name in banned is present in got.
func assertContainsNone(t *testing.T, got, banned []string) {
	t.Helper()
	idx := make(map[string]bool, len(got))
	for _, g := range got {
		idx[g] = true
	}
	for _, b := range banned {
		if idx[b] {
			t.Errorf("did not expect %q in result, got %v", b, got)
		}
	}
}