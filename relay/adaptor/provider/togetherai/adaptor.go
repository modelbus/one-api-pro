package togetherai

import oa "github.com/modelbus/one-api-pro/relay/adaptor/openai"

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"meta-llama/Llama-3-70b-chat-hf",
	"deepseek-ai/deepseek-coder-33b-instruct",
	"mistralai/Mixtral-8x22B-Instruct-v0.1",
	"Qwen/Qwen1.5-72B-Chat",
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
