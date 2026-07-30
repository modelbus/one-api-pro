package xunfei

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "xunfei",
		Name:           "Xunfei",
		DefaultBaseURL: "",
		LegacyType:     18,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
