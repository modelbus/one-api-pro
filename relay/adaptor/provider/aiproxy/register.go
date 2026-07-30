package aiproxy

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "aiproxylibrary",
		Name:           "AIProxy Library",
		DefaultBaseURL: "https://api.aiproxy.io",
		LegacyType:     21,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
