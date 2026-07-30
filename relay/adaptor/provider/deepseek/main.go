package deepseek

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/common"
	"github.com/Leon-PanPan/one-api-pro/common/render"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	"github.com/Leon-PanPan/one-api-pro/relay/schema"
)

type DeepSeekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`

	PromptTokensDetails     *model.PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *model.CompletionTokensDetails `json:"completion_tokens_details,omitempty"`

	PromptCacheHitTokens int `json:"prompt_cache_hit_tokens"`
}

type DeepSeekSlimTextResponse struct {
	Choices     []openai.TextResponseChoice `json:"choices"`
	Usage       DeepSeekUsage               `json:"usage"`
	Error       model.Error                 `json:"error"`
}

type DeepSeekStreamResponse struct {
	Id      string                                  `json:"id"`
	Object  string                                  `json:"object"`
	Created int64                                   `json:"created"`
	Model   string                                  `json:"model"`
	Choices []openai.ChatCompletionsStreamResponseChoice `json:"choices"`
	Usage   *DeepSeekUsage                          `json:"usage,omitempty"`
}

func convertDeepSeekUsage(usage *DeepSeekUsage) *model.Usage {
	if usage == nil {
		return nil
	}
	u := &model.Usage{
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
	}
	cachedTokens := usage.PromptCacheHitTokens
	if usage.PromptTokensDetails != nil && usage.PromptTokensDetails.CachedTokens > 0 {
		cachedTokens = usage.PromptTokensDetails.CachedTokens
	}
	if cachedTokens > 0 {
		u.PromptTokensDetails = &model.PromptTokensDetails{
			CachedTokens: cachedTokens,
		}
	}
	if usage.CompletionTokensDetails != nil {
		u.CompletionTokensDetails = usage.CompletionTokensDetails
	}
	return u
}

func StreamHandler(c *gin.Context, resp *http.Response) (*model.ErrorWithStatusCode, string, *model.Usage) {
	responseText := ""
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)
	var usage *model.Usage

	common.SetEventStreamHeaders(c)

	doneRendered := false
	for scanner.Scan() {
		data := scanner.Text()
		if len(data) < 6 {
			continue
		}
		if !strings.HasPrefix(data, "data: ") && !strings.HasPrefix(data, "data:") {
			continue
		}
		dataStr := data
		if strings.HasPrefix(data, "data: ") {
			dataStr = data[6:]
		} else {
			dataStr = data[5:]
		}
		dataStr = strings.TrimSpace(dataStr)
		if dataStr == "[DONE]" {
			render.StringData(c, data)
			doneRendered = true
			continue
		}
		var streamResponse DeepSeekStreamResponse
		err := json.Unmarshal([]byte(dataStr), &streamResponse)
		if err != nil {
			render.StringData(c, data)
			continue
		}
		if streamResponse.Usage != nil {
			usage = convertDeepSeekUsage(streamResponse.Usage)
		}
		render.StringData(c, data)
		for _, choice := range streamResponse.Choices {
			if choice.Delta.Content != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					responseText += content
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		_ = resp.Body.Close()
		return openai.ErrorWrapper(err, "read_stream_failed", http.StatusInternalServerError), "", nil
	}

	if !doneRendered {
		render.Done(c)
	}

	_ = resp.Body.Close()
	return nil, responseText, usage
}

func Handler(c *gin.Context, resp *http.Response, promptTokens int, modelName string) (*model.ErrorWithStatusCode, *model.Usage) {
	var textResponse DeepSeekSlimTextResponse
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return openai.ErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError), nil
	}
	err = resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}
	err = json.Unmarshal(responseBody, &textResponse)
	if err != nil {
		return openai.ErrorWrapper(err, "unmarshal_response_body_failed", http.StatusInternalServerError), nil
	}
	if textResponse.Error.Type != "" {
		return &model.ErrorWithStatusCode{
			Error:      textResponse.Error,
			StatusCode: resp.StatusCode,
		}, nil
	}
	resp.Body = io.NopCloser(strings.NewReader(string(responseBody)))

	for k, v := range resp.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Writer.WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		return openai.ErrorWrapper(err, "copy_response_body_failed", http.StatusInternalServerError), nil
	}
	err = resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}

	usage := convertDeepSeekUsage(&textResponse.Usage)
	if usage.TotalTokens == 0 || (usage.PromptTokens == 0 && usage.CompletionTokens == 0) {
		var contentBuilder strings.Builder
		for _, choice := range textResponse.Choices {
			contentBuilder.WriteString(choice.Message.StringContent())
		}
		completionTokens := openai.CountTokenText(contentBuilder.String(), modelName)
		usage = &model.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		}
	}
	return nil, usage
}