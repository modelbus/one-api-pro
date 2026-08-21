package xunfei

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
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
