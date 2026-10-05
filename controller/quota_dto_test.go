// quota_dto_test.go API 边界额度换算（微元 <-> 元）的单元测试
// Unit tests for API-boundary quota conversion (micro-quota <-> CNY yuan)
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode
package controller

import (
	"encoding/json"
	"testing"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/model"
)

// TestQuotaToYuan_Basic 验证微元 → 元换算与基准常量一致。
func TestQuotaToYuan_Basic(t *testing.T) {
	if config.QuotaPerUnit != 1_000_000.0 {
		t.Fatalf("QuotaPerUnit = %v，本测试假设其为 1_000_000", config.QuotaPerUnit)
	}
	cases := []struct {
		micro int64
		want  float64
	}{
		{0, 0},
		{10_000_000, 10},
		{2_500_000, 2.5},
		{1, 0.000001},
		{-1, -0.000001}, // 令牌「无限额度」哨兵值
		{10_500_000, 10.5},
	}
	for _, c := range cases {
		if got := quotaToYuan(c.micro); got != c.want {
			t.Fatalf("quotaToYuan(%d) = %v，期望 %v", c.micro, got, c.want)
		}
	}
}

// TestYuanToQuota_RoundsInsteadOfTruncating 验证元 → 微元使用四舍五入。
func TestYuanToQuota_RoundsInsteadOfTruncating(t *testing.T) {
	cases := []struct {
		yuan float64
		want int64
	}{
		{10, 10_000_000},
		{2.5, 2_500_000},
		// 0.07 × 1e6 在浮点下是 69999.99…，直接截断会少 1
		{0.07, 70_000},
		{0, 0},
		{-1, 0},
	}
	for _, c := range cases {
		if got := yuanToQuota(c.yuan); got != c.want {
			t.Fatalf("yuanToQuota(%v) = %d，期望 %d", c.yuan, got, c.want)
		}
	}
}

// TestUserDTO_ShadowsAmountFields 验证 DTO 用 float64 覆盖内嵌 model 的 int64 字段，
// 且其余字段原样透传（这是整套换算方案的基础机制）。
func TestUserDTO_ShadowsAmountFields(t *testing.T) {
	u := &model.User{Id: 7, Username: "alice", Quota: 10_000_000, UsedQuota: 2_500_000}
	b, err := json.Marshal(toUserDTO(u))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["quota"] != 10.0 {
		t.Fatalf("quota = %v，期望 10（元）", got["quota"])
	}
	if got["used_quota"] != 2.5 {
		t.Fatalf("used_quota = %v，期望 2.5（元）", got["used_quota"])
	}
	if got["username"] != "alice" {
		t.Fatalf("username = %v，期望 alice（透传字段被破坏）", got["username"])
	}
	if got["id"] != 7.0 {
		t.Fatalf("id = %v，期望 7（透传字段被破坏）", got["id"])
	}
}

// TestTokenDTO_ShadowsAmountFields 验证令牌 DTO 换算。
func TestTokenDTO_ShadowsAmountFields(t *testing.T) {
	tk := &model.Token{Id: 3, Name: "t", RemainQuota: 5_000_000, UsedQuota: 1_000_000}
	b, err := json.Marshal(toTokenDTO(tk))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]interface{}
	_ = json.Unmarshal(b, &got)
	if got["remain_quota"] != 5.0 || got["used_quota"] != 1.0 {
		t.Fatalf("remain_quota=%v used_quota=%v，期望 5 / 1", got["remain_quota"], got["used_quota"])
	}
	if got["name"] != "t" {
		t.Fatalf("name = %v，期望 t", got["name"])
	}
}

// TestTokenDTO_UnlimitedSentinel 验证无限额度令牌的 -1 哨兵在输出端仍为负数。
func TestTokenDTO_UnlimitedSentinel(t *testing.T) {
	tk := &model.Token{RemainQuota: -1, UnlimitedQuota: true}
	b, _ := json.Marshal(toTokenDTO(tk))
	var got map[string]interface{}
	_ = json.Unmarshal(b, &got)
	if v, ok := got["remain_quota"].(float64); !ok || v >= 0 {
		t.Fatalf("remain_quota = %v，期望负数（前端据 <0 判定为不限额）", got["remain_quota"])
	}
}

// TestChannelDTO_KeepsBalanceUntouched 验证渠道余额（USD）不被换算。
func TestChannelDTO_KeepsBalanceUntouched(t *testing.T) {
	ch := &model.Channel{Id: 1, UsedQuota: 3_000_000, Balance: 12.34}
	b, _ := json.Marshal(toChannelDTO(ch))
	var got map[string]interface{}
	_ = json.Unmarshal(b, &got)
	if got["used_quota"] != 3.0 {
		t.Fatalf("used_quota = %v，期望 3（元）", got["used_quota"])
	}
	if got["balance"] != 12.34 {
		t.Fatalf("balance = %v，期望 12.34（USD 原样）", got["balance"])
	}
}

// TestLogDTO_ConvertsIntQuota 验证日志的 `quota`（model.Log 中是 int）也被换算。
func TestLogDTO_ConvertsIntQuota(t *testing.T) {
	l := &model.Log{Id: 1, Quota: 2_500_000, ModelName: "gpt-4o"}
	b, _ := json.Marshal(toLogDTO(l))
	var got map[string]interface{}
	_ = json.Unmarshal(b, &got)
	if got["quota"] != 2.5 {
		t.Fatalf("quota = %v，期望 2.5（元）", got["quota"])
	}
	if got["model_name"] != "gpt-4o" {
		t.Fatalf("model_name = %v，期望 gpt-4o", got["model_name"])
	}
}

// TestRedemptionDTO_ConvertsQuota 验证兑换码 DTO 换算。
func TestRedemptionDTO_ConvertsQuota(t *testing.T) {
	r := &model.Redemption{Id: 1, Quota: 20_000_000}
	b, _ := json.Marshal(toRedemptionDTO(r))
	var got map[string]interface{}
	_ = json.Unmarshal(b, &got)
	if got["quota"] != 20.0 {
		t.Fatalf("quota = %v，期望 20（元）", got["quota"])
	}
}

// TestLogStatisticDTO_KeepsCapitalizedFieldName 验证聚合统计沿用 model 的大写字段名。
func TestLogStatisticDTO_KeepsCapitalizedFieldName(t *testing.T) {
	s := &model.LogStatistic{Day: "2026-10-04", ModelName: "gpt-4o", RequestCount: 3, Quota: 7_500_000}
	b, _ := json.Marshal(toLogStatisticDTO(s))
	var got map[string]interface{}
	_ = json.Unmarshal(b, &got)
	if got["Quota"] != 7.5 {
		t.Fatalf("Quota = %v，期望 7.5（元）", got["Quota"])
	}
	if _, exists := got["quota"]; exists {
		t.Fatalf("不应出现小写 quota 字段（会破坏既有调用方）")
	}
	if got["Day"] != "2026-10-04" || got["RequestCount"] != 3.0 {
		t.Fatalf("透传字段被破坏: %+v", got)
	}
}

// TestUpdateUserRequest_ToModel 验证 PUT /api/user/ 的元 → 微元换算。
func TestUpdateUserRequest_ToModel(t *testing.T) {
	var req updateUserRequest
	if err := json.Unmarshal([]byte(`{"id":9,"username":"bob","quota":10.5}`), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	u := req.toModel()
	if u.Quota != 10_500_000 {
		t.Fatalf("quota = %d，期望 10500000", u.Quota)
	}
	if u.Id != 9 || u.Username != "bob" {
		t.Fatalf("透传字段丢失: %+v", u)
	}
}

// TestUpdateUserRequest_FractionalYuanIsAccepted 验证小数元不会解码失败。
// 若直接解码到 model.User.Quota(int64)，10.5 会报
// "cannot unmarshal number 10.5 into ... of type int64"。
func TestUpdateUserRequest_FractionalYuanIsAccepted(t *testing.T) {
	var req updateUserRequest
	if err := json.Unmarshal([]byte(`{"id":1,"quota":0.000001}`), &req); err != nil {
		t.Fatalf("小数元应可解码，实际报错: %v", err)
	}
	if got := req.toModel().Quota; got != 1 {
		t.Fatalf("quota = %d，期望 1（1 微元）", got)
	}
}

// TestUpdateUserRequest_OmittedQuotaIsNil 验证未传 quota 时指针为 nil，
// 调用方据此跳过「零值补写」，避免把额度误置为 0。
func TestUpdateUserRequest_OmittedQuotaIsNil(t *testing.T) {
	var req updateUserRequest
	if err := json.Unmarshal([]byte(`{"id":9,"username":"bob"}`), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Quota != nil {
		t.Fatalf("未传 quota 时 Quota 应为 nil，实际 %v", *req.Quota)
	}
	if got := req.toModel().Quota; got != 0 {
		t.Fatalf("quota = %d，期望 0", got)
	}
}

// TestUpdateUserRequest_ExplicitZeroQuotaIsNotNil 验证 quota=0 与「未传」可区分，
// 这是管理端把额度清零能真正落库的前提。
func TestUpdateUserRequest_ExplicitZeroQuotaIsNotNil(t *testing.T) {
	var req updateUserRequest
	if err := json.Unmarshal([]byte(`{"id":9,"quota":0}`), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Quota == nil {
		t.Fatal("显式 quota=0 应被识别为已传值")
	}
	if got := req.toModel().Quota; got != 0 {
		t.Fatalf("quota = %d，期望 0", got)
	}
}

// TestTokenWriteRequest_UnlimitedSentinel 验证 -1 哨兵不被换算成 -1000000。
func TestTokenWriteRequest_UnlimitedSentinel(t *testing.T) {
	var req tokenWriteRequest
	if err := json.Unmarshal([]byte(`{"name":"t","remain_quota":-1,"unlimited_quota":true}`), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	tk := req.toToken()
	if tk.RemainQuota != -1 {
		t.Fatalf("remain_quota = %d，期望 -1（哨兵原样透传）", tk.RemainQuota)
	}
	if !tk.UnlimitedQuota {
		t.Fatalf("unlimited_quota 应透传为 true")
	}
}

// TestTokenWriteRequest_NormalQuota 验证常规令牌额度换算。
func TestTokenWriteRequest_NormalQuota(t *testing.T) {
	var req tokenWriteRequest
	if err := json.Unmarshal([]byte(`{"remain_quota":1}`), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := req.toToken().RemainQuota; got != 1_000_000 {
		t.Fatalf("remain_quota = %d，期望 1000000（1 元）", got)
	}
}

// TestRedemptionWriteRequest_ToRedemption 验证兑换码元 → 微元换算。
func TestRedemptionWriteRequest_ToRedemption(t *testing.T) {
	var req redemptionWriteRequest
	if err := json.Unmarshal([]byte(`{"name":"n","count":5,"quota":30}`), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	r := req.toRedemption()
	if r.Quota != 30_000_000 {
		t.Fatalf("quota = %d，期望 30000000", r.Quota)
	}
	if r.Name != "n" || r.Count != 5 {
		t.Fatalf("透传字段丢失: %+v", r)
	}
}

// TestTopupPresetsToModel 验证充值快捷金额的元 → 微元换算。
func TestTopupPresetsToModel(t *testing.T) {
	in := []topupPresetInput{
		{Amount: 10, BonusQuota: 15},   // 充 10 得 15
		{Amount: 20, BonusQuota: 20},   // 1:1
		{Amount: 30, BonusQuota: 30.5}, // 小数
	}
	out := topupPresetsToModel(in)
	if len(out) != 3 {
		t.Fatalf("长度 = %d，期望 3", len(out))
	}
	if out[0].Amount != 10 || out[0].BonusQuota != 15_000_000 {
		t.Fatalf("第 1 项 = %+v，期望 amount=10 bonus_quota=15000000", out[0])
	}
	if out[1].BonusQuota != 20_000_000 {
		t.Fatalf("第 2 项 bonus_quota = %d，期望 20000000", out[1].BonusQuota)
	}
	if out[2].BonusQuota != 30_500_000 {
		t.Fatalf("第 3 项 bonus_quota = %d，期望 30500000", out[2].BonusQuota)
	}
}

// TestPlanInfoQuotaToYuan 验证订单快照内的 bonus_quota 改写，其余字段保留。
func TestPlanInfoQuotaToYuan(t *testing.T) {
	raw := `{"amount":10,"credit_amount":15,"bonus_quota":15000000}`
	got := planInfoQuotaToYuan(raw)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("结果不是合法 JSON: %v", err)
	}
	if m["bonus_quota"] != 15.0 {
		t.Fatalf("bonus_quota = %v，期望 15（元）", m["bonus_quota"])
	}
	if m["amount"] != 10.0 || m["credit_amount"] != 15.0 {
		t.Fatalf("其余字段被破坏: %+v", m)
	}
}

// TestPlanInfoQuotaToYuan_Passthrough 验证无 bonus_quota 或非法输入时原样返回。
func TestPlanInfoQuotaToYuan_Passthrough(t *testing.T) {
	cases := []string{
		"",                      // 空
		`{"plan_name":"pro"}`,   // 套餐订单，无 bonus_quota
		`not-json`,              // 非法 JSON
		`{"bonus_quota":"abc"}`, // 非数值
	}
	for _, c := range cases {
		if got := planInfoQuotaToYuan(c); got != c {
			t.Fatalf("planInfoQuotaToYuan(%q) = %q，期望原样返回", c, got)
		}
	}
}

// TestToUserDTO_NilSafe 验证 nil 输入不会 panic（内嵌指针为 nil 时 json 仍可序列化）。
func TestToUserDTO_NilSafe(t *testing.T) {
	if toUserDTO(nil) != nil {
		t.Fatalf("toUserDTO(nil) 应返回 nil")
	}
	if toTokenDTO(nil) != nil || toChannelDTO(nil) != nil || toRedemptionDTO(nil) != nil || toLogDTO(nil) != nil {
		t.Fatalf("nil 输入应返回 nil")
	}
	if got := toUserDTOs(nil); got == nil || len(got) != 0 {
		t.Fatalf("toUserDTOs(nil) 应返回空切片而非 nil")
	}
}
