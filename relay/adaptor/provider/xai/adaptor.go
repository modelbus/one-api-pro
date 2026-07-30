package xai

import oa "github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"grok-2",
	"grok-vision-beta",
	"grok-2-vision-1212",
	"grok-2-vision",
	"grok-2-vision-latest",
	"grok-2-1212",
	"grok-2-latest",
	"grok-beta",
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
