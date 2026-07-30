package deepl

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "deepl",
		Name:           "DeepL",
		DefaultBaseURL: "https://api-free.deepl.com",
		LegacyType:     38,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
