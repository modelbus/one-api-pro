package vertexai

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
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
