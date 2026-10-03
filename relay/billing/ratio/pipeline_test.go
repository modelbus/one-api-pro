// pipeline_test 计费口径一致性测试：单价（元/百万 tokens）与扣费额度、余额（元）三方对齐
// Billing-consistency tests: unit price (CNY/M tokens) vs consumed quota vs balance (CNY)
// 版本: v0.0.24
// 日期: 2026-10-03
// 作者: opencode
//
// 说明：额度基准 QuotaPerUnit 为系统常量（1 元 = 1_000_000 quota），
// 定价以「元 / 百万 tokens」录入，保留 6 位小数——三者精度相互对齐。
// 本测试用真实金额做精确断言，防止后续改动破坏 1:1:1 的换算关系。

package ratio

import (
	"math"
	"testing"

	"github.com/modelbus/one-api-pro/common/config"
)

// TestYuanQuotaRoundTrip 验证「分」粒度的金额经 元→quota→元 往返无偏差。
// 版本: v0.0.24
func TestYuanQuotaRoundTrip(t *testing.T) {
	for cents := int64(1); cents <= 100_000; cents++ {
		yuan := float64(cents) / 100
		quota := int64(math.Round(yuan * config.QuotaPerUnit))
		back := QuotaToUnit(quota)
		if math.Abs(back-yuan) > 1e-9 {
			t.Fatalf("往返偏差：%.2f 元 → %d quota → %.12f 元", yuan, quota, back)
		}
	}
}

// TestTokenBilling_PriceMatchesYuan 验证 token 计费额度折算回人民币后
// 恰好等于「单价 × tokens」。
// 版本: v0.0.24
func TestTokenBilling_PriceMatchesYuan(t *testing.T) {
	cases := []struct {
		name                   string
		in, out, cached        float64
		prompt, completion, cc int
		discount               float64
	}{
		{"gpt-4o 全量", 2.5, 10.0, 1.25, 1_000_000, 500_000, 200_000, 1.0},
		{"deepseek 含缓存", 0.14, 0.28, 0.014, 500_000, 100_000, 100_000, 1.0},
		{"小额度含折扣", 8.40, 8.40, 0, 1_000, 1_000, 0, 0.8},
		{"单次极小额", 2.5, 10.0, 1.25, 1_000, 500, 0, 1.0},
	}
	for _, c := range cases {
		quota := CalculateTokenQuota(c.in, c.out, c.cached, c.prompt, c.completion, c.cc, c.discount)
		got := QuotaToUnit(quota)
		inputTokens := c.prompt - c.cc
		want := (c.in*float64(inputTokens) + c.out*float64(c.completion) + c.cached*float64(c.cc)) * c.discount / Million
		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("%s：实际 ¥%.9f，期望 ¥%.9f", c.name, got, want)
		}
	}
}

// TestPerRequestBilling_PriceMatchesYuan 验证按次计费额度折算回人民币后
// 恰好等于「单次单价 × n × 折扣」。
// 版本: v0.0.24
func TestPerRequestBilling_PriceMatchesYuan(t *testing.T) {
	cases := []struct {
		price          float64
		sizeRatio, n   float64
		discount, want float64
	}{
		{0.04, 1, 1, 1.0, 0.04},
		{0.08, 2, 3, 0.9, 0.08 * 2 * 3 * 0.9},
		{0.12, 1, 2, 1.0, 0.24},
	}
	for _, c := range cases {
		quota := CalculatePerRequestQuota(c.price, c.sizeRatio, int(c.n), c.discount)
		got := QuotaToUnit(quota)
		if math.Abs(got-c.want) > 1e-9 {
			t.Fatalf("按次计费：单价 %.2f → 实际 ¥%.9f，期望 ¥%.9f", c.price, got, c.want)
		}
	}
}

// TestQuotaPerUnitAlignment 验证常量语义：1 元 = QuotaPerUnit 额度，且 6 位小数定价可被精确表示。
// 版本: v0.0.24
func TestQuotaPerUnitAlignment(t *testing.T) {
	if config.QuotaPerUnit != 1_000_000 {
		t.Fatalf("QuotaPerUnit = %v，期望 1e6", config.QuotaPerUnit)
	}
	// 定价最小刻度 0.000001 元 → 1 quota，正好对齐
	if q := int64(math.Round(0.000001 * config.QuotaPerUnit)); q != 1 {
		t.Fatalf("0.000001 元 = %d quota，期望 1（6 位小数定价精度对齐）", q)
	}
	// 1 元 → 1_000_000 quota → 1 元
	if got := QuotaToUnit(int64(config.QuotaPerUnit)); got != 1.0 {
		t.Fatalf("1e6 quota = ¥%v，期望 ¥1", got)
	}
}
