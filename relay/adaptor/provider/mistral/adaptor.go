package mistral

import oa "github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"open-mistral-7b",
	"open-mixtral-8x7b",
	"mistral-small-latest",
	"mistral-medium-latest",
	"mistral-large-latest",
	"mistral-embed",
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
