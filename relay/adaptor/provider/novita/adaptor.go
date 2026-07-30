package novita

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
	"meta-llama/llama-3-8b-instruct",
	"meta-llama/llama-3-70b-instruct",
	"nousresearch/hermes-2-pro-llama-3-8b",
	"nousresearch/nous-hermes-llama2-13b",
	"mistralai/mistral-7b-instruct",
	"cognitivecomputations/dolphin-mixtral-8x22b",
	"sao10k/l3-70b-euryale-v2.1",
	"sophosympatheia/midnight-rose-70b",
	"gryphe/mythomax-l2-13b",
	"Nous-Hermes-2-Mixtral-8x7B-DPO",
	"lzlv_70b",
	"teknium/openhermes-2.5-mistral-7b",
	"microsoft/wizardlm-2-8x22b",
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	if meta.Mode == relaymode.ChatCompletions {
		return fmt.Sprintf("%s/chat/completions", meta.BaseURL), nil
	}
	return "", fmt.Errorf("unsupported relay mode %d for novita", meta.Mode)
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
