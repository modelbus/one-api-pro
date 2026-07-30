package proxy

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "proxy",
		Name:           "Proxy",
		DefaultBaseURL: "",
		LegacyType:     43,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
