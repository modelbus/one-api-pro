package stepfun

import oa "github.com/modelbus/one-api-pro/relay/adaptor/openai"

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"step-1-8k",
	"step-1-32k",
	"step-1-128k",
	"step-1-256k",
	"step-1-flash",
	"step-2-16k",
	"step-1v-8k",
	"step-1v-32k",
	"step-1x-medium",
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
