package openrouter

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "openrouter",
		Name:           "OpenRouter",
		DefaultBaseURL: "https://openrouter.ai/api",
		LegacyType:     20,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
