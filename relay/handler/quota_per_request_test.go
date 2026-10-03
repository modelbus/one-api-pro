// quota_per_request_test 按次（per_request）计费的精确预扣回归测试
// Per-request exact pre-consume regression tests
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
//
// 背景 / Background:
//   per_request 计费下费用在请求前已完全确定，但旧实现仍用 token 估算口径预扣
//   （几百 quota），结算时才按按次价格扣数万 quota；两者无比例关系，结算 delta
//   直接落库，单次调用即可把 users.quota 扣成负数。
//
//   修复后：per_request 预扣额 = CalculatePerRequestQuota(...)，与结算额一致，
//   delta 恒为 0；余额不足时在预扣阶段即被 403 拒绝。
//
// 覆盖 / Coverage:
//   - 精确预扣：预扣额 == 按次价格折算的 quota，余额精确减少该金额
//   - 余额不足拒绝：预扣阶段 403，余额不变（旧实现会放行并透支）
//   - 端到端不透支：pre + settle 后余额 = 初始 - 按次价，且 usage=0 不退费
//   - token 计费仍保留 usage=0 免计费守卫

package controller

import (
	"context"
	"strings"
	"testing"

	"github.com/modelbus/one-api-pro/common/config"
	dbmodel "github.com/modelbus/one-api-pro/model"
	billingratio "github.com/modelbus/one-api-pro/relay/billing/ratio"
	"github.com/modelbus/one-api-pro/relay/meta"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

// perRequestPriceResult 构造一个 per_request 定价的 PriceResult。
// 版本: v0.0.24
func perRequestPriceResult(price float64) *billingratio.PriceResult {
	return &billingratio.PriceResult{
		PerRequestPrice: price,
		BillingType:     dbmodel.BillingTypePerRequest,
		Found:           true,
	}
}

// TestPerRequestPreConsume_ExactAmount 验证 per_request 预扣额恰好等于
// 按次价格折算的 quota，且余额精确减少该金额。
// 版本: v0.0.24
func TestPerRequestPreConsume_ExactAmount(t *testing.T) {
	setupPreConsumeTestDB(t)

	const (
		perRequestPrice   = 0.04 // ¥0.04/次 → 40_000 quota
		initialUserQuota  = int64(200_000)
		initialTokenQuota = int64(5_000_000)
	)
	uid := seedUser(t, initialUserQuota)
	tid := seedToken(t, uid, initialTokenQuota)

	exactQuota := billingratio.CalculatePerRequestQuota(perRequestPrice, 1, 1, 1.0)
	if exactQuota != 40_000 {
		t.Fatalf("exactQuota = %d，期望 40000", exactQuota)
	}

	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}
	preConsumed, errObj := preConsumeExactQuota(context.Background(), exactQuota, m)
	if errObj != nil {
		t.Fatalf("preConsumeExactQuota 失败：%+v", errObj)
	}
	// 余额 200_000 < 100×40_000 → 正常预扣，不做 trusted 跳过
	if preConsumed != exactQuota {
		t.Fatalf("preConsumed = %d，期望 %d（精确按次预扣）", preConsumed, exactQuota)
	}
	if got := readUserQuota(t, uid); got != initialUserQuota-exactQuota {
		t.Fatalf("预扣后 user.quota = %d，期望 %d", got, initialUserQuota-exactQuota)
	}
	tk := readToken(t, tid)
	if tk.RemainQuota != initialTokenQuota-exactQuota {
		t.Fatalf("token.remain_quota = %d，期望 %d", tk.RemainQuota, initialTokenQuota-exactQuota)
	}
}

// TestPerRequestPreConsume_InsufficientBalanceRejects 是透支 bug 的核心回归测试：
// 余额小于按次价格时必须 403 且余额不变。
//
// 旧实现下，预扣用的是 token 估算额（此处约 (500+100)×ratio≈几百），远小于余额，
// 会放行；随后结算按 40_000 扣款，直接把余额扣成负数。
// 版本: v0.0.24
func TestPerRequestPreConsume_InsufficientBalanceRejects(t *testing.T) {
	setupPreConsumeTestDB(t)

	const (
		perRequestPrice  = 0.04 // 40_000 quota
		initialUserQuota = int64(20_000)
	)
	uid := seedUser(t, initialUserQuota)
	tid := seedToken(t, uid, 5_000_000)

	exactQuota := billingratio.CalculatePerRequestQuota(perRequestPrice, 1, 1, 1.0)

	// 旧口径的 token 估算额：明显小于余额，旧实现不会拒绝
	legacyEstimate := getPreConsumedQuota(&relaymodel.GeneralOpenAIRequest{MaxTokens: 0}, 100, 1.0)
	if legacyEstimate >= initialUserQuota {
		t.Fatalf("测试前提不成立：旧估算额 %d 应远小于余额 %d", legacyEstimate, initialUserQuota)
	}

	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}
	_, errObj := preConsumeExactQuota(context.Background(), exactQuota, m)
	if errObj == nil {
		t.Fatalf("期望 insufficient_user_quota，实际放行（余额 %d，按次价 %d）", initialUserQuota, exactQuota)
	}
	if errObj.Error.Code != "insufficient_user_quota" {
		t.Fatalf("期望 code=insufficient_user_quota，实际 %s", errObj.Error.Code)
	}
	if got := readUserQuota(t, uid); got != initialUserQuota {
		t.Fatalf("拒绝时余额不应变动：got %d, want %d", got, initialUserQuota)
	}
}

// TestPerRequestPipeline_NoOverdraft 端到端验证 pre + settle 后余额恰好为
// 初始余额减去按次价格，且 usage=0 时也不会把按次费用退回。
// 版本: v0.0.24
func TestPerRequestPipeline_NoOverdraft(t *testing.T) {
	setupPreConsumeTestDB(t)
	extendPreConsumeTestDB(t)

	const (
		perRequestPrice   = 0.04 // 40_000 quota
		initialUserQuota  = int64(200_000)
		initialTokenQuota = int64(5_000_000)
	)
	uid := seedUser(t, initialUserQuota)
	tid := seedToken(t, uid, initialTokenQuota)

	priceResult := perRequestPriceResult(perRequestPrice)
	exactQuota := billingratio.CalculatePerRequestQuota(perRequestPrice, 1, 1, 1.0)

	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}
	req := &relaymodel.GeneralOpenAIRequest{}
	preConsumed, errObj := preConsumeExactQuota(context.Background(), exactQuota, m)
	if errObj != nil {
		t.Fatalf("preConsumeExactQuota 失败：%+v", errObj)
	}

	// 上游返回 usage=0（空响应）也不应把按次费用退回
	usage := &relaymodel.Usage{PromptTokens: 0, CompletionTokens: 0}
	postConsumeQuota(context.Background(), usage, m, req, preConsumed, priceResult, 1.0, false)

	if got := readUserQuota(t, uid); got != initialUserQuota-exactQuota {
		t.Fatalf("结算后 user.quota = %d，期望 %d（净扣应为按次价，不能透支）",
			got, initialUserQuota-exactQuota)
	}
	if got := readUserQuota(t, uid); got < 0 {
		t.Fatalf("余额被扣成负数：%d", got)
	}
	tk := readToken(t, tid)
	if tk.UsedQuota != exactQuota {
		t.Fatalf("token.used_quota = %d，期望 %d", tk.UsedQuota, exactQuota)
	}
	if tk.RemainQuota+tk.UsedQuota != initialTokenQuota {
		t.Fatalf("token 账目不平：remain(%d) + used(%d) ≠ 初始(%d)",
			tk.RemainQuota, tk.UsedQuota, initialTokenQuota)
	}

	var log dbmodel.Log
	if err := dbmodel.DB.Where("user_id = ?", uid).Order("id desc").First(&log).Error; err != nil {
		t.Fatalf("read consume log: %v", err)
	}
	if int64(log.Quota) != exactQuota {
		t.Fatalf("日志 quota = %d，期望 %d", log.Quota, exactQuota)
	}
	if !strings.Contains(log.Content, "按次计费") {
		t.Fatalf("日志内容应标注按次计费，实际：%s", log.Content)
	}
}

// TestTokenBilling_ZeroUsageStillFree 确认 token 计费的 usage=0 免计费守卫
// 在修复 per_request 豁免后仍然保留。
// 版本: v0.0.24
func TestTokenBilling_ZeroUsageStillFree(t *testing.T) {
	setupPreConsumeTestDB(t)
	extendPreConsumeTestDB(t)

	const initialUserQuota = int64(200_000)
	uid := seedUser(t, initialUserQuota)
	tid := seedToken(t, uid, 5_000_000)

	priceResult := &billingratio.PriceResult{
		InputPrice:  2.5,
		OutputPrice: 10.0,
		CachedPrice: 1.25,
		BillingType: dbmodel.BillingTypeToken,
		Found:       true,
	}
	ratio := (priceResult.InputPrice + priceResult.OutputPrice) / 2.0 / billingratio.Million * config.QuotaPerUnit

	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}
	req := &relaymodel.GeneralOpenAIRequest{}
	preConsumed, errObj := preConsumeQuota(context.Background(), req, 0, ratio, m)
	if errObj != nil {
		t.Fatalf("preConsumeQuota 失败：%+v", errObj)
	}

	usage := &relaymodel.Usage{PromptTokens: 0, CompletionTokens: 0}
	postConsumeQuota(context.Background(), usage, m, req, preConsumed, priceResult, 1.0, false)

	if got := readUserQuota(t, uid); got != initialUserQuota {
		t.Fatalf("token 计费 usage=0 应全额退回：user.quota = %d，期望 %d", got, initialUserQuota)
	}
}
