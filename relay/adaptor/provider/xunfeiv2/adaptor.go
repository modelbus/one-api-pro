package xunfeiv2

import oa "github.com/modelbus/one-api-pro/relay/adaptor/openai"

type Adaptor struct {
	*oa.Adaptor
}

var ModelList = []string{
	"lite",
	"generalv3",
	"pro-128k",
	"generalv3.5",
	"max-32k",
	"4.0Ultra",
}

func (a *Adaptor) GetModelList() []string  { return ModelList }
