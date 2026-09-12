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