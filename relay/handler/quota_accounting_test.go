// quota_accounting_test 端到端额度会计测试：单价（元/百万 tokens）→ 扣费额度 → 余额（元）
// End-to-end quota accounting: unit price (CNY/M tokens) -> consumed quota -> balance (CNY)
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
//
// 覆盖：
//   - 预扣 + 结算的净扣费恰好等于按真实 usage 计算的实际费用（不漂移、不重复扣）
//   - 扣费额度按 1 元 = QuotaPerUnit 折算回人民币，与单价口径一致
//   - user.quota / token.remain_quota / token.used_quota 三方账目一致

package controller

import (
	"context"
	"testing"

	dbmodel "github.com/modelbus/one-api-pro/model"

	"github.com/modelbus/one-api-pro/common/config"
	billingratio "github.com/modelbus/one-api-pro/relay/billing/ratio"
	"github.com/modelbus/one-api-pro/relay/meta"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

// extendPreConsumeTestDB 在 pre-consume 测试库基础上补齐日志/渠道/订单表，
// 以便 postConsumeQuota 的落库路径可执行。
func extendPreConsumeTestDB(t *testing.T) {
	t.Helper()
	if err := dbmodel.DB.AutoMigrate(&dbmodel.Log{}, &dbmodel.Channel{}, &dbmodel.Order{}); err != nil {
		t.Fatalf("migrate extend: %v", err)
	}
	// 消费日志走独立的 LOG_DB；测试中与主库共用同一内存实例，避免 nil panic。
	dbmodel.LOG_DB = dbmodel.DB
}

func readUserQuota(t *testing.T, uid int) int64 {
	t.Helper()
	var u dbmodel.User
	if err := dbmodel.DB.First(&u, "id = ?", uid).Error; err != nil {
		t.Fatalf("read user: %v", err)
	}
	return u.Quota
}

func readToken(t *testing.T, tid int) dbmodel.Token {
	t.Helper()
	var tk dbmodel.Token
	if err := dbmodel.DB.First(&tk, "id = ?", tid).Error; err != nil {
		t.Fatalf("read token: %v", err)
	}
	return tk
}

// TestConsumePipeline_NetDeductionEqualsActual 验证「预扣 → 结算」后的净扣费
// 恰好等于按实际 usage 计算的费用，且换算回人民币与单价一致。
//
// 场景：输入 ¥2.5/M、输出 ¥10/M，实际用量 1000 prompt + 500 completion。
// 实际费用 = 2.5×1000 + 10×500 = 7500 quota = ¥0.0075。
// 版本: v0.0.24
func TestConsumePipeline_NetDeductionEqualsActual(t *testing.T) {
	setupPreConsumeTestDB(t)
	extendPreConsumeTestDB(t)

	const (
		initialUserQuota  = int64(200_000) // ¥0.20
		initialTokenQuota = int64(5_000_000)
	)
	uid := seedUser(t, initialUserQuota)
	tid := seedToken(t, uid, initialTokenQuota)

	inputPrice, outputPrice, cachedPrice := 2.5, 10.0, 1.25
	promptTokens, completionTokens := 1000, 500

	// 与 RelayTextHelper 相同的预扣比例：(输入价 + 输出价) / 2，单位「quota / token」
	ratio := (inputPrice + outputPrice) / 2.0 / billingratio.Million * config.QuotaPerUnit
	if ratio != 6.25 {
		t.Fatalf("ratio = %v，期望 6.25", ratio)
	}

	m := &meta.Meta{UserId: uid, TokenId: tid, ChannelId: 0, PlanId: 0}
	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 0}

	preConsumed, errObj := preConsumeQuota(context.Background(), req, promptTokens, ratio, m)
	if errObj != nil {
		t.Fatalf("preConsumeQuota 失败：%+v", errObj)
	}
	// (500 + 1000) × 6.25 = 9375
	if preConsumed != 9375 {
		t.Fatalf("preConsumed = %d，期望 9375", preConsumed)
	}
	if got := readUserQuota(t, uid); got != initialUserQuota-preConsumed {
		t.Fatalf("预扣后 user.quota = %d，期望 %d", got, initialUserQuota-preConsumed)
	}

	usage := &relaymodel.Usage{PromptTokens: promptTokens, CompletionTokens: completionTokens}
	priceResult := &billingratio.PriceResult{
		InputPrice:  inputPrice,
		OutputPrice: outputPrice,
		CachedPrice: cachedPrice,
		BillingType: dbmodel.BillingTypeToken,
		Found:       true,
	}
	postConsumeQuota(context.Background(), usage, m, req, preConsumed, priceResult, 1.0, false)

	// 实际费用（quota）：直接由单价 × 真实 tokens 得出
	expectedQuota := int64(inputPrice*float64(promptTokens) + outputPrice*float64(completionTokens))

	// 1) user.quota 净扣费 == 实际费用
	if got := readUserQuota(t, uid); got != initialUserQuota-expectedQuota {
		t.Fatalf("结算后 user.quota = %d，期望 %d（净扣 %d）",
			got, initialUserQuota-expectedQuota, initialUserQuota-got)
	}
	// 2) token 账目：remain 减 actual，used 增 actual
	tk := readToken(t, tid)
	if tk.RemainQuota != initialTokenQuota-expectedQuota {
		t.Fatalf("token.remain_quota = %d，期望 %d", tk.RemainQuota, initialTokenQuota-expectedQuota)
	}
	if tk.UsedQuota != expectedQuota {
		t.Fatalf("token.used_quota = %d，期望 %d", tk.UsedQuota, expectedQuota)
	}
	if tk.RemainQuota+tk.UsedQuota != initialTokenQuota {
		t.Fatalf("token 账目不平：remain(%d) + used(%d) ≠ 初始(%d)",
			tk.RemainQuota, tk.UsedQuota, initialTokenQuota)
	}

	// 3) 换算回人民币：余额差 = ¥0.0075
	spentYuan := billingratio.QuotaToUnit(initialUserQuota - readUserQuota(t, uid))
	if spentYuan != 0.0075 {
		t.Fatalf("本次消费 = ¥%v，期望 ¥0.0075", spentYuan)
	}
	// 4) 消费日志记录的 quota 与实际一致
	var log dbmodel.Log
	if err := dbmodel.DB.Where("user_id = ?", uid).Order("id desc").First(&log).Error; err != nil {
		t.Fatalf("read consume log: %v", err)
	}
	if int64(log.Quota) != expectedQuota {
		t.Fatalf("日志 quota = %d，期望 %d", log.Quota, expectedQuota)
	}
}

// TestConsumePipeline_ReturningPreconsume 验证实际用量低于预扣时，
// 差额被正确退回（余额回升而非继续减少）。
// 版本: v0.0.24
func TestConsumePipeline_ReturningPreconsume(t *testing.T) {
	setupPreConsumeTestDB(t)
	extendPreConsumeTestDB(t)

	const (
		initialUserQuota  = int64(200_000)
		initialTokenQuota = int64(5_000_000)
	)
	uid := seedUser(t, initialUserQuota)
	tid := seedToken(t, uid, initialTokenQuota)

	priceResult := &billingratio.PriceResult{
		InputPrice:  2.5,
		OutputPrice: 10.0,
		CachedPrice: 1.25,
		BillingType: dbmodel.BillingTypeToken,
		Found:       true,
	}
	ratio := (priceResult.InputPrice + priceResult.OutputPrice) / 2.0 / billingratio.Million * config.QuotaPerUnit

	m := &meta.Meta{UserId: uid, TokenId: tid, PlanId: 0}
	req := &relaymodel.GeneralOpenAIRequest{MaxTokens: 1000}
	preConsumed, errObj := preConsumeQuota(context.Background(), req, 100, ratio, m)
	if errObj != nil {
		t.Fatalf("preConsumeQuota 失败：%+v", errObj)
	}
	// (500 + 100 + 1000) × 6.25 = 10000
	if preConsumed != 10000 {
		t.Fatalf("preConsumed = %d，期望 10000", preConsumed)
	}

	// 实际只用了 100 prompt + 50 completion → 8250 quota
	usage := &relaymodel.Usage{PromptTokens: 100, CompletionTokens: 50}
	postConsumeQuota(context.Background(), usage, m, req, preConsumed, priceResult, 1.0, false)

	expectedQuota := int64(2.5*100 + 10.0*50) // 750
	if got := readUserQuota(t, uid); got != initialUserQuota-expectedQuota {
		t.Fatalf("user.quota = %d，期望 %d（预扣应被退回差额）", got, initialUserQuota-expectedQuota)
	}
	tk := readToken(t, tid)
	if tk.UsedQuota != expectedQuota {
		t.Fatalf("token.used_quota = %d，期望 %d", tk.UsedQuota, expectedQuota)
	}
}
