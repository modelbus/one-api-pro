package anthropic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/common"
	"github.com/Leon-PanPan/one-api-pro/common/helper"
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/Leon-PanPan/one-api-pro/common/render"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	"github.com/Leon-PanPan/one-api-pro/relay/schema"
)

type AnthropicStreamCtx struct {
	ID                 string
	ModelName          string
	CreatedTime        int64
	LastToolCallChoice openai.ChatCompletionsStreamResponseChoice
}

type AnthropicPipeline struct{}

func NewAnthropicPipeline() *AnthropicPipeline {
	return &AnthropicPipeline{}
}

func (p *AnthropicPipeline) Run(c *gin.Context, resp *http.Response) (*model.ErrorWithStatusCode, *model.Usage) {
	ctx := &AnthropicStreamCtx{
		CreatedTime: helper.GetTimestamp(),
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := strings.Index(string(data), "\n"); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})

	common.SetEventStreamHeaders(c)

	var usage model.Usage

	for scanner.Scan() {
		data := scanner.Text()
		if len(data) < 6 || !strings.HasPrefix(data, "data:") {
			continue
		}
		data = strings.TrimPrefix(data, "data:")
		data = strings.TrimSpace(data)

		var claudeResponse StreamResponse
		err := json.Unmarshal([]byte(data), &claudeResponse)
		if err != nil {
			logger.SysError("error unmarshalling stream response: " + err.Error())
			continue
		}

		response, meta := StreamResponseClaude2OpenAI(&claudeResponse)
		if meta != nil {
			usage.PromptTokens += meta.Usage.InputTokens
			usage.CompletionTokens += meta.Usage.OutputTokens
			if meta.Usage.CacheReadTokens > 0 {
				if usage.PromptTokensDetails == nil {
					usage.PromptTokensDetails = &model.PromptTokensDetails{}
				}
				usage.PromptTokensDetails.CachedTokens += meta.Usage.CacheReadTokens
			}
			if len(meta.Id) > 0 {
				ctx.ModelName = meta.Model
				ctx.ID = fmt.Sprintf("chatcmpl-%s", meta.Id)
				continue
			} else {
				if len(ctx.LastToolCallChoice.Delta.ToolCalls) > 0 {
					lastArgs := &ctx.LastToolCallChoice.Delta.ToolCalls[len(ctx.LastToolCallChoice.Delta.ToolCalls)-1].Function
					if len(lastArgs.Arguments.(string)) == 0 {
						lastArgs.Arguments = "{}"
						response.Choices[len(response.Choices)-1].Delta.Content = nil
						response.Choices[len(response.Choices)-1].Delta.ToolCalls = ctx.LastToolCallChoice.Delta.ToolCalls
					}
				}
			}
		}
		if response == nil {
			continue
		}

		response.Id = ctx.ID
		response.Model = ctx.ModelName
		response.Created = ctx.CreatedTime

		for _, choice := range response.Choices {
			if len(choice.Delta.ToolCalls) > 0 {
				ctx.LastToolCallChoice = choice
			}
		}
		err = render.ObjectData(c, response)
		if err != nil {
			logger.SysError(err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		logger.SysError("error reading stream: " + err.Error())
	}

	render.Done(c)

	err := resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}
	return nil, &usage
}
