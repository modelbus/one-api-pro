package palm

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "palm",
		Name:           "PaLM",
		DefaultBaseURL: "https://generativelanguage.googleapis.com",
		LegacyType:     11,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
