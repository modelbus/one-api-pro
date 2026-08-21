package deepseek

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/modelbus/one-api-pro/relay/schema"
)

func TestConvertDeepSeekUsage(t *testing.T) {
	Convey("convertDeepSeekUsage", t, func() {
		Convey("nil usage returns nil", func() {
			result := convertDeepSeekUsage(nil)
			So(result, ShouldBeNil)
		})

		Convey("basic usage without cache", func() {
			usage := &DeepSeekUsage{
				PromptTokens:     100,
				CompletionTokens: 50,
				TotalTokens:      150,
			}
			result := convertDeepSeekUsage(usage)
			So(result.PromptTokens, ShouldEqual, 100)
			So(result.CompletionTokens, ShouldEqual, 50)
			So(result.TotalTokens, ShouldEqual, 150)
			So(result.PromptTokensDetails, ShouldBeNil)
		})

		Convey("usage with prompt_cache_hit_tokens", func() {
			usage := &DeepSeekUsage{
				PromptTokens:          1000,
				CompletionTokens:      200,
				TotalTokens:           1200,
				PromptCacheHitTokens: 500,
			}
			result := convertDeepSeekUsage(usage)
			So(result.PromptTokens, ShouldEqual, 1000)
			So(result.CompletionTokens, ShouldEqual, 200)
			So(result.PromptTokensDetails, ShouldNotBeNil)
			So(result.PromptTokensDetails.CachedTokens, ShouldEqual, 500)
		})

		Convey("usage with prompt_tokens_details takes priority", func() {
			usage := &DeepSeekUsage{
				PromptTokens:     1000,
				CompletionTokens: 200,
				TotalTokens:      1200,
				PromptTokensDetails: &model.PromptTokensDetails{
					CachedTokens: 300,
				},
				PromptCacheHitTokens: 500,
			}
			result := convertDeepSeekUsage(usage)
			So(result.PromptTokensDetails.CachedTokens, ShouldEqual, 300)
		})

		Convey("zero cache hit tokens", func() {
			usage := &DeepSeekUsage{
				PromptTokens:          100,
				CompletionTokens:     50,
				TotalTokens:          150,
				PromptCacheHitTokens: 0,
			}
			result := convertDeepSeekUsage(usage)
			So(result.PromptTokensDetails, ShouldBeNil)
		})

		Convey("completion tokens details preserved", func() {
			usage := &DeepSeekUsage{
				PromptTokens:     100,
				CompletionTokens: 50,
				TotalTokens:      150,
				CompletionTokensDetails: &model.CompletionTokensDetails{
					ReasoningTokens: 20,
				},
			}
			result := convertDeepSeekUsage(usage)
			So(result.CompletionTokensDetails, ShouldNotBeNil)
			So(result.CompletionTokensDetails.ReasoningTokens, ShouldEqual, 20)
		})
	})
}

func TestDeepSeekUsageJSONParsing(t *testing.T) {
	Convey("DeepSeekUsage JSON parsing", t, func() {
		Convey("parse response with cache tokens", func() {
			jsonStr := `{
				"prompt_tokens": 1000,
				"completion_tokens": 200,
				"total_tokens": 1200,
				"prompt_cache_hit_tokens": 500
			}`
			var usage DeepSeekUsage
			err := json.Unmarshal([]byte(jsonStr), &usage)
			So(err, ShouldBeNil)
			So(usage.PromptTokens, ShouldEqual, 1000)
			So(usage.CompletionTokens, ShouldEqual, 200)
			So(usage.PromptCacheHitTokens, ShouldEqual, 500)
		})

		Convey("parse response without cache tokens", func() {
			jsonStr := `{
				"prompt_tokens": 100,
				"completion_tokens": 50,
				"total_tokens": 150
			}`
			var usage DeepSeekUsage
			err := json.Unmarshal([]byte(jsonStr), &usage)
			So(err, ShouldBeNil)
			So(usage.PromptCacheHitTokens, ShouldEqual, 0)
		})
	})
}

func TestDeepSeekSlimTextResponseJSONParsing(t *testing.T) {
	Convey("DeepSeekSlimTextResponse JSON parsing", t, func() {
		jsonStr := `{
			"choices": [{"index": 0, "message": {"role": "assistant", "content": "hello"}, "finish_reason": "stop"}],
			"usage": {
				"prompt_tokens": 1000,
				"completion_tokens": 200,
				"total_tokens": 1200,
				"prompt_cache_hit_tokens": 500
			}
		}`
		var resp DeepSeekSlimTextResponse
		err := json.Unmarshal([]byte(jsonStr), &resp)
		So(err, ShouldBeNil)
		So(resp.Usage.PromptCacheHitTokens, ShouldEqual, 500)

		result := convertDeepSeekUsage(&resp.Usage)
		So(result.PromptTokensDetails.CachedTokens, ShouldEqual, 500)
	})
}