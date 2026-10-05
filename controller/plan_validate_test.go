// plan_validate_test.go 套餐 HTTP 层的保存校验测试
// HTTP-layer tests for plan save-time validation
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
//
// 覆盖：
//   - AddPlan / UpdatePlan 拒绝非法 model_limits、非法 billing_type、维度不匹配
//   - 合法请求落库，且 virtual_amount 按「元 → 微元」换算
package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/model"
)

// planEnvelope mirrors the gin.H envelope returned by plan handlers.
type planEnvelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// postPlanJSON 以 POST/PUT 调用套餐 handler 并解析响应。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func postPlanJSON(t *testing.T, fn gin.HandlerFunc, method, payload string) planEnvelope {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, "/api/plan/", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	fn(c)
	var out planEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, w.Body.String())
	}
	return out
}

// TestAddPlan_ValidatesConfig 保存校验：非法配置必须被拒绝，合法配置成功。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestAddPlan_ValidatesConfig(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "valid token plan",
			payload: `{"name":"vt-valid","billing_type":"token","virtual_amount":10,"model_limits":"{\"gpt-4o\":{\"period_h\":5,\"token_period\":1000}}"}`,
			wantOK:  true,
		},
		{
			name:    "legacy number limits",
			payload: `{"name":"vt-legacy","billing_type":"token","model_limits":"{\"gpt-4\":1000}"}`,
			wantOK:  false,
			wantMsg: "格式错误",
		},
		{
			name:    "invalid json limits",
			payload: `{"name":"vt-badjson","billing_type":"token","model_limits":"{oops}"}`,
			wantOK:  false,
			wantMsg: "格式错误",
		},
		{
			name:    "empty limits",
			payload: `{"name":"vt-empty","billing_type":"token","model_limits":""}`,
			wantOK:  false,
			wantMsg: "不能为空",
		},
		{
			name:    "invalid billing type",
			payload: `{"name":"vt-billing","billing_type":"per_request","model_limits":"{\"gpt-4o\":{\"period_h\":5,\"token_period\":1000}}"}`,
			wantOK:  false,
			wantMsg: "计费维度",
		},
		{
			name:    "dimension mismatch",
			payload: `{"name":"vt-mismatch","billing_type":"request","model_limits":"{\"gpt-4o\":{\"period_h\":5,\"token_period\":1000}}"}`,
			wantOK:  false,
			wantMsg: "缺少与计费维度",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := postPlanJSON(t, AddPlan, http.MethodPost, tc.payload)
			if res.Success != tc.wantOK {
				t.Fatalf("success = %v, want %v (message=%s)", res.Success, tc.wantOK, res.Message)
			}
			if !tc.wantOK && !strings.Contains(res.Message, tc.wantMsg) {
				t.Fatalf("message = %q, want to contain %q", res.Message, tc.wantMsg)
			}
		})
	}
}

// TestUpdatePlan_RejectsInvalidAndPersistsYuan 更新时的校验与单位换算。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestUpdatePlan_RejectsInvalidAndPersistsYuan(t *testing.T) {
	seed := &model.Plan{
		Name:        "vt-update-source",
		BillingType: model.BillingTypeToken,
		ModelLimits: `{"gpt-4o":{"period_h":5,"token_period":1000}}`,
	}
	if err := seed.Insert(); err != nil {
		t.Fatalf("seed plan: %v", err)
	}

	// 非法配置：拒绝且不落库。
	bad := `{"id":` + strconv.Itoa(int(seed.Id)) + `,"name":"vt-update-source","billing_type":"token","model_limits":"{\"gpt-4\":1000}"}`
	res := postPlanJSON(t, UpdatePlan, http.MethodPut, bad)
	if res.Success {
		t.Fatalf("invalid update should be rejected, got success (message=%s)", res.Message)
	}
	reloaded, err := model.GetPlanById(int(seed.Id))
	if err != nil {
		t.Fatalf("reload plan: %v", err)
	}
	if reloaded.ModelLimits != seed.ModelLimits {
		t.Fatalf("model_limits should be untouched, got %s", reloaded.ModelLimits)
	}

	// 合法配置：成功，且 virtual_amount 以微元落库。
	good := `{"id":` + strconv.Itoa(int(seed.Id)) + `,"name":"vt-update-source","billing_type":"request","virtual_amount":12.5,"model_limits":"{\"gpt-4o\":{\"period_h\":5,\"request_period\":100}}"}`
	res = postPlanJSON(t, UpdatePlan, http.MethodPut, good)
	if !res.Success {
		t.Fatalf("valid update failed: %s", res.Message)
	}
	reloaded, err = model.GetPlanById(int(seed.Id))
	if err != nil {
		t.Fatalf("reload plan: %v", err)
	}
	if reloaded.BillingType != model.BillingTypeRequest {
		t.Fatalf("billing_type = %s, want request", reloaded.BillingType)
	}
	if reloaded.VirtualAmount != int64(12.5*config.QuotaPerUnit) {
		t.Fatalf("virtual_amount = %d, want %d", reloaded.VirtualAmount, int64(12.5*config.QuotaPerUnit))
	}
}
