package baichuan

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "baichuan",
		Name:           "Baichuan",
		DefaultBaseURL: "https://api.baichuan-ai.com",
		LegacyType:     26,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
