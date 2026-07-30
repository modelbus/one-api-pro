package openrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/meta"
	oa "github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
)

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"01-ai/yi-large",
	"aetherwiing/mn-starcannon-12b",
	"ai21/jamba-1-5-large",
	"ai21/jamba-1-5-mini",
	"anthropic/claude-3-haiku",
	"anthropic/claude-3-opus",
	"anthropic/claude-3-sonnet",
	"anthropic/claude-3.5-haiku",
	"anthropic/claude-3.5-sonnet",
	"cohere/command-r",
	"cohere/command-r-plus",
	"deepseek/deepseek-chat",
	"deepseek/deepseek-r1",
	"google/gemini-2.0-flash-001",
	"google/gemini-flash-1.5",
	"google/gemini-pro",
	"google/gemini-pro-1.5",
	"google/gemma-2-9b-it",
	"meta-llama/llama-3-70b-instruct",
	"meta-llama/llama-3-8b-instruct",
	"meta-llama/llama-3.1-70b-instruct",
	"meta-llama/llama-3.1-8b-instruct",
	"meta-llama/llama-3.2-1b-instruct",
	"meta-llama/llama-3.2-3b-instruct",
	"meta-llama/llama-3.3-70b-instruct",
	"microsoft/phi-3-mini-128k-instruct",
	"microsoft/phi-3-medium-128k-instruct",
	"microsoft/phi-4",
	"mistralai/mistral-7b-instruct",
	"mistralai/mistral-large",
	"mistralai/mistral-medium",
	"mistralai/mistral-small",
	"mistralai/mixtral-8x7b",
	"nousresearch/hermes-3-llama-3.1-70b",
	"openai/gpt-3.5-turbo",
	"openai/gpt-4",
	"openai/gpt-4-turbo",
	"openai/gpt-4o",
	"openai/gpt-4o-mini",
	"openai/o1",
	"openai/o1-mini",
	"openai/o1-preview",
	"openai/o3-mini",
	"openrouter/auto",
	"qwen/qwen-2.5-72b-instruct",
	"qwen/qwen-2.5-7b-instruct",
	"qwen/qwen-2.5-coder-32b-instruct",
	"x-ai/grok-beta",
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("Authorization", "Bearer "+meta.APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/Leon-PanPan/one-api-pro")
	req.Header.Set("X-Title", "One Api Pro")
	return nil
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
