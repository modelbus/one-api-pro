package baidu

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "baidu",
		Name:           "Baidu",
		DefaultBaseURL: "https://aip.baidubce.com",
		LegacyType:     15,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
