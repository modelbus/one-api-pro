package geminiv2

import (
	"fmt"
	"strings"

	"github.com/modelbus/one-api-pro/relay/meta"
	oa "github.com/modelbus/one-api-pro/relay/adaptor/openai"
)

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"gemini-pro", "gemini-1.0-pro",
	"gemini-1.5-flash", "gemini-1.5-flash-8b",
	"gemini-1.5-pro", "gemini-1.5-pro-experimental",
	"text-embedding-004", "aqa",
	"gemini-2.0-flash", "gemini-2.0-flash-exp",
	"gemini-2.0-flash-lite-preview-02-05",
	"gemini-2.0-flash-thinking-exp-01-21",
	"gemini-2.0-pro-exp-02-05",
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	baseURL := strings.TrimSuffix(meta.BaseURL, "/")
	requestPath := strings.TrimPrefix(meta.RequestURLPath, "/v1")
	return fmt.Sprintf("%s%s", baseURL, requestPath), nil
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
