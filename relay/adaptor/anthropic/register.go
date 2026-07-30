package anthropic

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "anthropic",
		Name:           "Anthropic",
		DefaultBaseURL: "https://api.anthropic.com",
		LegacyType:     14,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
