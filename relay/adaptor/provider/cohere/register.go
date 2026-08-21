package cohere

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "cohere",
		Name:           "Cohere",
		DefaultBaseURL: "https://api.cohere.ai",
		LegacyType:     35,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
