package azure

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "azure",
		Name:           "Azure",
		DefaultBaseURL: "",
		LegacyType:     3,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
