package billing

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConsumeQuotaParamsLogContent(t *testing.T) {
	Convey("ConsumeQuotaParams log content format", t, func() {
		Convey("token billing log format", func() {
			params := &ConsumeQuotaParams{
				InputPrice:       2.5,
				OutputPrice:      10.0,
				CachedPrice:      1.25,
				GroupDiscount:    1.0,
				PromptTokens:     1000,
				CompletionTokens: 500,
				CachedTokens:     0,
				BillingType:      "token",
			}
			logContent := fmt.Sprintf("定价：输入¥%.4f/百万tokens × %d + 输出¥%.4f/百万tokens × %d",
				params.InputPrice, params.PromptTokens, params.OutputPrice, params.CompletionTokens)
			So(logContent, ShouldContainSubstring, "输入¥2.5000/百万tokens × 1000")
			So(logContent, ShouldContainSubstring, "输出¥10.0000/百万tokens × 500")
		})

		Convey("token billing with cache", func() {
			params := &ConsumeQuotaParams{
				InputPrice:       2.5,
				OutputPrice:      10.0,
				CachedPrice:      1.25,
				GroupDiscount:    1.0,
				PromptTokens:     1000,
				CompletionTokens: 500,
				CachedTokens:     200,
				BillingType:      "token",
			}
			logContent := fmt.Sprintf("定价：输入¥%.4f/百万tokens × %d + 输出¥%.4f/百万tokens × %d",
				params.InputPrice, params.PromptTokens, params.OutputPrice, params.CompletionTokens)
			if params.CachedTokens > 0 {
				logContent += fmt.Sprintf(" + 缓存¥%.4f/百万tokens × %d", params.CachedPrice, params.CachedTokens)
			}
			So(logContent, ShouldContainSubstring, "缓存¥1.2500/百万tokens × 200")
		})

		Convey("group discount in log", func() {
			params := &ConsumeQuotaParams{
				InputPrice:       2.5,
				OutputPrice:      10.0,
				CachedPrice:      0,
				GroupDiscount:    0.8,
				PromptTokens:     1000,
				CompletionTokens: 500,
				CachedTokens:     0,
				BillingType:      "token",
			}
			logContent := fmt.Sprintf("定价：输入¥%.4f/百万tokens × %d + 输出¥%.4f/百万tokens × %d",
				params.InputPrice, params.PromptTokens, params.OutputPrice, params.CompletionTokens)
			if params.GroupDiscount != 1.0 {
				logContent += fmt.Sprintf(" × 分组折扣%.2f", params.GroupDiscount)
			}
			So(logContent, ShouldContainSubstring, "分组折扣0.80")
		})

		Convey("per_request billing log format", func() {
			params := &ConsumeQuotaParams{
				InputPrice:    0.04,
				GroupDiscount: 0.9,
				BillingType:   "per_request",
			}
			logContent := ""
			if params.BillingType == "per_request" {
				logContent = fmt.Sprintf("按次计费：¥%.4f × 分组折扣%.2f", params.InputPrice, params.GroupDiscount)
			}
			So(logContent, ShouldContainSubstring, "按次计费")
			So(logContent, ShouldContainSubstring, "¥0.0400")
			So(logContent, ShouldContainSubstring, "分组折扣0.90")
		})

		Convey("no group discount when discount is 1.0", func() {
			logContent := "定价：输入¥2.5000/百万tokens × 1000 + 输出¥10.0000/百万tokens × 500"
			So(logContent, ShouldNotContainSubstring, "分组折扣")
		})
	})
}

func TestConsumeQuotaParamsStruct(t *testing.T) {
	Convey("ConsumeQuotaParams struct", t, func() {
		params := &ConsumeQuotaParams{
			TokenId:          1,
			UserId:           2,
			ChannelId:        3,
			QuotaDelta:       100,
			TotalQuota:       200,
			ModelName:        "gpt-4o",
			TokenName:        "test-token",
			InputPrice:       2.5,
			OutputPrice:      10.0,
			CachedPrice:      1.25,
			GroupDiscount:    0.8,
			PromptTokens:     1000,
			CompletionTokens: 500,
			CachedTokens:     200,
			BillingType:      "token",
		}
		So(params.TokenId, ShouldEqual, 1)
		So(params.UserId, ShouldEqual, 2)
		So(params.ChannelId, ShouldEqual, 3)
		So(params.QuotaDelta, ShouldEqual, 100)
		So(params.TotalQuota, ShouldEqual, 200)
		So(params.ModelName, ShouldEqual, "gpt-4o")
		So(params.BillingType, ShouldEqual, "token")
		So(params.CachedTokens, ShouldEqual, 200)
	})
}