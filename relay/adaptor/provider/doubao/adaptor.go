package doubao

import (
	"fmt"

	"github.com/modelbus/one-api-pro/relay/meta"
	"github.com/modelbus/one-api-pro/relay/relaymode"
	oa "github.com/modelbus/one-api-pro/relay/adaptor/openai"
)

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"Doubao-pro-128k",
	"Doubao-pro-32k",
	"Doubao-pro-4k",
	"Doubao-lite-128k",
	"Doubao-lite-32k",
	"Doubao-lite-4k",
	"Doubao-embedding",
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	switch meta.Mode {
	case relaymode.ChatCompletions:
		return fmt.Sprintf("%s/api/v3/chat/completions", meta.BaseURL), nil
	case relaymode.Embeddings:
		return fmt.Sprintf("%s/api/v3/embeddings", meta.BaseURL), nil
	default:
	}
	return "", fmt.Errorf("unsupported relay mode %d for doubao", meta.Mode)
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
