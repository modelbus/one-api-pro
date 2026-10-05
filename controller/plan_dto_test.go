// plan_dto_test.go 套餐 DTO 单位换算（元 <-> 微元）的单元测试
// Unit tests for plan DTO unit conversion (yuan <-> micro-quota)
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode

package controller

import (
	"encoding/json"
	"testing"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/model"
)

// TestPlanWriteRequest_ToPlan 请求体元 → 存储微元。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestPlanWriteRequest_ToPlan(t *testing.T) {
	req := planWriteRequest{VirtualAmount: 100}
	p := req.toPlan()
	if p.VirtualAmount != int64(100*config.QuotaPerUnit) {
		t.Fatalf("virtual amount = %d, want %d", p.VirtualAmount, int64(100*config.QuotaPerUnit))
	}

	// 小数元也必须精确（避免解码到 int64 报错的老问题）。
	req = planWriteRequest{VirtualAmount: 10.5}
	if got := req.toPlan().VirtualAmount; got != 10_500_000 {
		t.Fatalf("virtual amount = %d, want 10500000", got)
	}

	// 0 表示不限额度。
	zero := planWriteRequest{VirtualAmount: 0}
	if got := zero.toPlan().VirtualAmount; got != 0 {
		t.Fatalf("virtual amount = %d, want 0", got)
	}
}

// TestPlanWriteRequest_Unmarshal 请求体按「元」解析 virtual_amount。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestPlanWriteRequest_Unmarshal(t *testing.T) {
	body := `{"name":"p","billing_type":"token","virtual_amount":99.99}`
	var req planWriteRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.VirtualAmount != 99.99 {
		t.Fatalf("parsed virtual amount = %v, want 99.99", req.VirtualAmount)
	}
	// 内层 int64 字段不应被 JSON 直接赋值（避免单位混淆）。
	if req.Plan.VirtualAmount != 0 {
		t.Fatalf("embedded int64 should stay zero before toPlan, got %d", req.Plan.VirtualAmount)
	}
	if got := req.toPlan().VirtualAmount; got != 99_990_000 {
		t.Fatalf("converted virtual amount = %d, want 99990000", got)
	}
}

// TestToPlanDTO 存储微元 → 输出元。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestToPlanDTO(t *testing.T) {
	dto := toPlanDTO(&model.Plan{Id: 1, Name: "p", VirtualAmount: 100_000_000})
	if dto == nil {
		t.Fatal("dto should not be nil")
	}
	if dto.VirtualAmount != 100 {
		t.Fatalf("virtual amount = %v, want 100", dto.VirtualAmount)
	}
	// 序列化后 virtual_amount 必须是元（100），不是微元。
	b, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded["virtual_amount"] != float64(100) {
		t.Fatalf("json virtual_amount = %v, want 100", decoded["virtual_amount"])
	}
}

// TestToUserPlanDTO 订阅输出中 virtual/used 均为元，且内嵌 plan 同样元化。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestToUserPlanDTO(t *testing.T) {
	up := &model.UserPlan{
		Id:            1,
		VirtualAmount: 100_000_000,
		UsedAmount:    30_000_000,
		Plan:          &model.Plan{Id: 1, VirtualAmount: 100_000_000},
	}
	dto := toUserPlanDTO(up)
	if dto.VirtualAmount != 100 || dto.UsedAmount != 30 {
		t.Fatalf("virtual=%v used=%v, want 100/30", dto.VirtualAmount, dto.UsedAmount)
	}
	if dto.Plan == nil || dto.Plan.VirtualAmount != 100 {
		t.Fatalf("embedded plan should be converted to yuan")
	}
}

// TestPlanInfoQuotaToYuan_VirtualAmount 订单快照中的 virtual_amount 由微元改写为元。
//
// 版本: v0.0.25
// 日期: 2026-10-05
func TestPlanInfoQuotaToYuan_VirtualAmount(t *testing.T) {
	raw := `{"name":"p","virtual_amount":100000000}`
	out := planInfoQuotaToYuan(raw)
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded["virtual_amount"] != float64(100) {
		t.Fatalf("virtual_amount = %v, want 100", decoded["virtual_amount"])
	}
	if decoded["name"] != "p" {
		t.Fatalf("other fields should be preserved, got %v", decoded["name"])
	}
}
