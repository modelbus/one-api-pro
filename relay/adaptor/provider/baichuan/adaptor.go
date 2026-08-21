package baichuan

import oa "github.com/modelbus/one-api-pro/relay/adaptor/openai"

type Adaptor struct {
	*oa.Adaptor
}

func (a *Adaptor) NeedsRequestBodyConversion() bool {
	return true
}

var ModelList = []string{
	"Baichuan2-Turbo",
	"Baichuan2-Turbo-192k",
	"Baichuan-Text-Embedding",
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
