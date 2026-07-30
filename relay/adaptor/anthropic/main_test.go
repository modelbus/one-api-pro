package anthropic

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

		Convey("message_delta with usage", func() {
			jsonStr := `{
				"type": "message_delta",
				"delta": {"stop_reason": "end_turn"},
				"usage": {
					"input_tokens": 0,
					"output_tokens": 200,
					"cache_read_input_tokens": 0
				}
			}`
			var sr StreamResponse
			err := json.Unmarshal([]byte(jsonStr), &sr)
			So(err, ShouldBeNil)
			So(sr.Usage, ShouldNotBeNil)
			So(sr.Usage.OutputTokens, ShouldEqual, 200)
		})
	})
}