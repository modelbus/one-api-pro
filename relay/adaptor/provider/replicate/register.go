package replicate

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "replicate",
		Name:           "Replicate",
		DefaultBaseURL: "https://api.replicate.com/v1/models/",
		LegacyType:     46,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
