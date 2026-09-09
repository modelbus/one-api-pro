package anthropic

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/modelbus/one-api-pro/relay/schema"
)

func TestAnthropicUsageCacheParsing(t *testing.T) {
	Convey("Anthropic Usage with cache fields", t, func() {
		Convey("parse cache_read_input_tokens", func() {
			jsonStr := `{
				"input_tokens": 1000,
				"output_tokens": 200,
				"cache_read_input_tokens": 300,
				"cache_creation_input_tokens": 100
			}`
			var usage Usage
			err := json.Unmarshal([]byte(jsonStr), &usage)
			So(err, ShouldBeNil)
			So(usage.InputTokens, ShouldEqual, 1000)
			So(usage.OutputTokens, ShouldEqual, 200)
			So(usage.CacheReadTokens, ShouldEqual, 300)
			So(usage.CacheCreationTokens, ShouldEqual, 100)
		})

		Convey("parse usage without cache fields", func() {
			jsonStr := `{
				"input_tokens": 1000,
				"output_tokens": 200
			}`
			var usage Usage
			err := json.Unmarshal([]byte(jsonStr), &usage)
			So(err, ShouldBeNil)
			So(usage.CacheReadTokens, ShouldEqual, 0)
			So(usage.CacheCreationTokens, ShouldEqual, 0)
		})
	})
}

func TestAnthropicResponseCacheParsing(t *testing.T) {
	Convey("Anthropic Response with cache usage", t, func() {
		jsonStr := `{
			"id": "msg_test",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Hello"}],
			"model": "claude-3-5-sonnet-20241022",
			"stop_reason": "end_turn",
			"usage": {
				"input_tokens": 1000,
				"output_tokens": 200,
				"cache_read_input_tokens": 500,
				"cache_creation_input_tokens": 100
			}
		}`
		var resp Response
		err := json.Unmarshal([]byte(jsonStr), &resp)
		So(err, ShouldBeNil)
		So(resp.Usage.CacheReadTokens, ShouldEqual, 500)
		So(resp.Usage.CacheCreationTokens, ShouldEqual, 100)
	})
}

func TestStreamResponseCacheParsing(t *testing.T) {
	Convey("Anthropic StreamResponse with cache usage", t, func() {
		Convey("message_start with cache", func() {
			jsonStr := `{
				"type": "message_start",
				"message": {
					"id": "msg_test",
					"type": "message",
					"role": "assistant",
					"content": [],
					"model": "claude-3-5-sonnet-20241022",
					"usage": {
						"input_tokens": 1000,
						"output_tokens": 0,
						"cache_read_input_tokens": 800
					}
				}
			}`
			var sr StreamResponse
			err := json.Unmarshal([]byte(jsonStr), &sr)
			So(err, ShouldBeNil)
			So(sr.Message, ShouldNotBeNil)
			So(sr.Message.Usage.CacheReadTokens, ShouldEqual, 800)
		})

		Convey("message_delta with cumulative usage", func() {
			// 新版 API 下 message_delta 携带整段请求的累计 usage（重复 message_start 的 input/cache），
			// 不再是最早「仅 output、input 为 0」的增量形态（见 #13）。
			jsonStr := `{
				"type": "message_delta",
				"delta": {"stop_reason": "end_turn"},
				"usage": {
					"input_tokens": 1000,
					"output_tokens": 200,
					"cache_read_input_tokens": 400,
					"cache_creation_input_tokens": 0
				}
			}`
			var sr StreamResponse
			err := json.Unmarshal([]byte(jsonStr), &sr)
			So(err, ShouldBeNil)
			So(sr.Usage, ShouldNotBeNil)
			So(sr.Usage.InputTokens, ShouldEqual, 1000)
			So(sr.Usage.OutputTokens, ShouldEqual, 200)
			So(sr.Usage.CacheReadTokens, ShouldEqual, 400)
		})
	})
}

func TestClaudeUsage2OpenAI(t *testing.T) {
	Convey("ClaudeUsage2OpenAI maps disjoint Claude usage to OpenAI subset semantics", t, func() {
		Convey("folds cache read and creation into prompt tokens", func() {
			u := ClaudeUsage2OpenAI(Usage{
				InputTokens:         1000,
				OutputTokens:        200,
				CacheReadTokens:     300,
				CacheCreationTokens: 100,
			})
			So(u.PromptTokens, ShouldEqual, 1400)
			So(u.CompletionTokens, ShouldEqual, 200)
			So(u.TotalTokens, ShouldEqual, 1600)
			So(u.PromptTokensDetails, ShouldNotBeNil)
			So(u.PromptTokensDetails.CachedTokens, ShouldEqual, 300)
		})

		Convey("no cache fields yields plain prompt tokens", func() {
			u := ClaudeUsage2OpenAI(Usage{InputTokens: 1000, OutputTokens: 200})
			So(u.PromptTokens, ShouldEqual, 1000)
			So(u.CompletionTokens, ShouldEqual, 200)
			So(u.PromptTokensDetails, ShouldBeNil)
		})
	})
}

func TestMergeClaudeUsage(t *testing.T) {
	Convey("MergeClaudeUsage takes max of cumulative stream usage", t, func() {
		Convey("message_start + message_delta cumulative sequence does not double count", func() {
			var acc model.Usage
			// message_start 携带的累计 usage
			MergeClaudeUsage(&acc, Usage{InputTokens: 1000, OutputTokens: 3})
			// message_delta 重复 input 并给出最终 output
			MergeClaudeUsage(&acc, Usage{InputTokens: 1000, OutputTokens: 200})
			So(acc.PromptTokens, ShouldEqual, 1000)
			So(acc.CompletionTokens, ShouldEqual, 200)
		})

		Convey("legacy delta shape (input omitted) keeps message_start input", func() {
			var acc model.Usage
			MergeClaudeUsage(&acc, Usage{InputTokens: 1000})
			MergeClaudeUsage(&acc, Usage{OutputTokens: 200})
			So(acc.PromptTokens, ShouldEqual, 1000)
			So(acc.CompletionTokens, ShouldEqual, 200)
		})

		Convey("cache read counted once and stays subset of prompt", func() {
			var acc model.Usage
			MergeClaudeUsage(&acc, Usage{InputTokens: 600, CacheReadTokens: 400, CacheCreationTokens: 100})
			MergeClaudeUsage(&acc, Usage{InputTokens: 600, CacheReadTokens: 400, CacheCreationTokens: 100, OutputTokens: 200})
			So(acc.PromptTokens, ShouldEqual, 1100)
			So(acc.CompletionTokens, ShouldEqual, 200)
			So(acc.PromptTokensDetails.CachedTokens, ShouldEqual, 400)
			So(acc.PromptTokensDetails.CachedTokens, ShouldBeLessThanOrEqualTo, acc.PromptTokens)
		})
	})
}
