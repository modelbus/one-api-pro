package alibailian

import (
	"fmt"

	"github.com/Leon-PanPan/one-api-pro/relay/meta"
	"github.com/Leon-PanPan/one-api-pro/relay/relaymode"
	oa "github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
)

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"qwen-turbo",
	"qwen-plus",
	"qwen-long",
	"qwen-max",
	"qwen-coder-plus",
	"qwen-coder-plus-latest",
	"qwen-coder-turbo",
	"qwen-coder-turbo-latest",
	"qwen-mt-plus",
	"qwen-mt-turbo",
	"qwq-32b-preview",
	"deepseek-r1",
	"deepseek-v3",
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	switch meta.Mode {
	case relaymode.ChatCompletions:
		return fmt.Sprintf("%s/compatible-mode/v1/chat/completions", meta.BaseURL), nil
	case relaymode.Embeddings:
		return fmt.Sprintf("%s/compatible-mode/v1/embeddings", meta.BaseURL), nil
	default:
	}
	return "", fmt.Errorf("unsupported relay mode %d for ali bailian", meta.Mode)
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
