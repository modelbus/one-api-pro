package cloudflare

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "cloudflare",
		Name:           "Cloudflare",
		DefaultBaseURL: "https://api.cloudflare.com",
		LegacyType:     37,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
