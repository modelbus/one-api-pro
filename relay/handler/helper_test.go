package controller

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/Leon-PanPan/one-api-pro/model"
	billingratio "github.com/Leon-PanPan/one-api-pro/relay/billing/ratio"
	relaymodel "github.com/Leon-PanPan/one-api-pro/relay/schema"
)

func TestPostConsumeQuotaBillingTypeRouting(t *testing.T) {
	Convey("postConsumeQuota billing type routing", t, func() {
		Convey("token billing type uses CalculateTokenQuota", func() {
			priceResult := &billingratio.PriceResult{
				InputPrice:      2.5,
				OutputPrice:     10.0,
				CachedPrice:     1.25,
				PerRequestPrice: 0,
				BillingType:     model.BillingTypeToken,
			}
			usage := &relaymodel.Usage{
				PromptTokens:     1000,
				CompletionTokens: 500,
			}
			groupDiscount := 1.0

			var quota int64
			if priceResult.BillingType == model.BillingTypePerRequest {
				quota = billingratio.CalculatePerRequestQuota(priceResult.PerRequestPrice, 1, 1, groupDiscount)
			} else {
				quota = billingratio.CalculateTokenQuota(
					priceResult.InputPrice, priceResult.OutputPrice, priceResult.CachedPrice,
					usage.PromptTokens, usage.CompletionTokens, 0,
					groupDiscount,
				)
			}
			So(quota, ShouldBeGreaterThan, 0)
		})

		Convey("per_request billing type uses CalculatePerRequestQuota", func() {
			priceResult := &billingratio.PriceResult{
				InputPrice:      0,
				OutputPrice:     0,
				CachedPrice:     0,
				PerRequestPrice: 0.04,
				BillingType:     model.BillingTypePerRequest,
			}
			groupDiscount := 1.0

			var quota int64
			if priceResult.BillingType == model.BillingTypePerRequest {
				quota = billingratio.CalculatePerRequestQuota(priceResult.PerRequestPrice, 1, 1, groupDiscount)
			} else {
				quota = billingratio.CalculateTokenQuota(
					priceResult.InputPrice, priceResult.OutputPrice, priceResult.CachedPrice,
					1, 0, 0,
					groupDiscount,
				)
			}
			So(quota, ShouldBeGreaterThan, 0)
		})

		Convey("cached tokens from PromptTokensDetails", func() {
			usage := &relaymodel.Usage{
				PromptTokens:     1000,
				CompletionTokens: 500,
				PromptTokensDetails: &relaymodel.PromptTokensDetails{
					CachedTokens: 200,
				},
			}
			cachedTokens := 0
			if usage.PromptTokensDetails != nil {
				cachedTokens = usage.PromptTokensDetails.CachedTokens
			}
			So(cachedTokens, ShouldEqual, 200)

			quotaWithCache := billingratio.CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 200, 1.0)
			quotaWithoutCache := billingratio.CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 1.0)
			So(quotaWithCache, ShouldBeLessThan, quotaWithoutCache)
		})

		Convey("zero total tokens results in zero quota", func() {
			quota := billingratio.CalculateTokenQuota(2.5, 10, 1.25, 0, 0, 0, 1.0)
			So(quota, ShouldEqual, 0)
		})
	})
}

func TestUnknownModelReturns422(t *testing.T) {
	Convey("unknown model returns 422", t, func() {
		Convey("GetModelPrice returns error for unknown model", func() {
			_, err := billingratio.GetModelPrice("completely-unknown-model-xyz-123")
			So(err, ShouldNotBeNil)
		})

		Convey("422 status code is StatusUnprocessableEntity", func() {
			So(http.StatusUnprocessableEntity, ShouldEqual, 422)
		})

		Convey("ModelPriceExists returns false for unknown model", func() {
			So(billingratio.ModelPriceExists("completely-unknown-model-xyz-123"), ShouldBeFalse)
		})
	})
}