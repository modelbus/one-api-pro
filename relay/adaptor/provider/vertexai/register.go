package vertexai

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "vertexai",
		Name:           "Vertex AI",
		DefaultBaseURL: "",
		LegacyType:     42,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
