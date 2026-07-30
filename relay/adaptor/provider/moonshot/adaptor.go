package moonshot

import oa "github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"moonshot-v1-8k",
	"moonshot-v1-32k",
	"moonshot-v1-128k",
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
