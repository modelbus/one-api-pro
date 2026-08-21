package ratio

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/model"
)

func TestModelPriceIntegration(t *testing.T) {
	Convey("Model Price Integration Tests", t, func() {
		originalQuotaPerUnit := config.QuotaPerUnit
		config.QuotaPerUnit = 500000
		defer func() {
			config.QuotaPerUnit = originalQuotaPerUnit
		}()

		Convey("full token billing pipeline", func() {
			inputPrice := 2.5
			outputPrice := 10.0
			cachedPrice := 1.25
			promptTokens := 1000000
			completionTokens := 500000
			cachedTokens := 200000

			groupDiscount := 0.8
			quota := CalculateTokenQuota(inputPrice, outputPrice, cachedPrice, promptTokens, completionTokens, cachedTokens, groupDiscount)

			inputTokens := promptTokens - cachedTokens
			expected := int64((inputPrice*float64(inputTokens) + outputPrice*float64(completionTokens) + cachedPrice*float64(cachedTokens)) * groupDiscount / 1_000_000 * 500000)
			So(quota, ShouldEqual, expected)
		})

		Convey("per request billing pipeline", func() {
			perRequestPrice := 0.04
			sizeRatio := 2.0
			n := 3
			groupDiscount := 0.9

			quota := CalculatePerRequestQuota(perRequestPrice, sizeRatio, n, groupDiscount)
			expected := int64(perRequestPrice * sizeRatio * float64(n) * groupDiscount * 500000)
			So(quota, ShouldEqual, expected)
		})

		Convey("deepseek-chat pricing", func() {
			inputPrice := 0.14
			outputPrice := 0.28
			cachedPrice := 0.014

			promptTokens := 500000
			completionTokens := 100000
			cachedTokens := 100000
			groupDiscount := 1.0

			quota := CalculateTokenQuota(inputPrice, outputPrice, cachedPrice, promptTokens, completionTokens, cachedTokens, groupDiscount)

			inputTokens := promptTokens - cachedTokens
			expected := int64((inputPrice*float64(inputTokens) + outputPrice*float64(completionTokens) + cachedPrice*float64(cachedTokens)) * groupDiscount / 1_000_000 * 500000)
			So(quota, ShouldEqual, expected)
			So(quota, ShouldBeGreaterThan, 0)
		})
	})
}

func TestModelPriceConstants(t *testing.T) {
	Convey("Billing type constants match between model and ratio packages", t, func() {
		So(model.BillingTypeToken, ShouldEqual, "token")
		So(model.BillingTypePerRequest, ShouldEqual, "per_request")
	})
}

func TestEdgeCases(t *testing.T) {
	Convey("Edge cases", t, func() {
		originalQuotaPerUnit := config.QuotaPerUnit
		config.QuotaPerUnit = 500000
		defer func() {
			config.QuotaPerUnit = originalQuotaPerUnit
		}()

		Convey("all cached tokens", func() {
			quota := CalculateTokenQuota(2.5, 10, 1.25, 1000, 0, 1000, 1.0)
			So(quota, ShouldBeGreaterThan, 0)
		})

		Convey("very small token count rounds up to 1", func() {
			quota := CalculateTokenQuota(0.001, 0.001, 0, 1, 1, 0, 1.0)
			So(quota, ShouldEqual, 1)
		})

		Convey("zero price model with tokens still returns 1 (minimum)", func() {
			quota := CalculateTokenQuota(0, 0, 0, 1000, 500, 0, 1.0)
			So(quota, ShouldEqual, 1)
		})

		Convey("large token count", func() {
			quota := CalculateTokenQuota(2.5, 10, 1.25, 1000000, 500000, 200000, 1.0)
			So(quota, ShouldBeGreaterThan, 0)
		})

		Convey("very small group discount", func() {
			quota := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 0.01)
			So(quota, ShouldBeGreaterThan, 0)
		})

		Convey("per request with zero price", func() {
			quota := CalculatePerRequestQuota(0, 1, 1, 1.0)
			So(quota, ShouldEqual, 1)
		})
	})
}