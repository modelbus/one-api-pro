package stepfun

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "stepfun",
		Name:           "StepFun",
		DefaultBaseURL: "https://api.stepfun.com",
		LegacyType:     32,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
