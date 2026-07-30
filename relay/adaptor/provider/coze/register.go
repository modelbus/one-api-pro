package coze

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "coze",
		Name:           "Coze",
		DefaultBaseURL: "https://api.coze.com",
		LegacyType:     34,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
