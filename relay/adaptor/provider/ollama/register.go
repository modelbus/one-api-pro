package ollama

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "ollama",
		Name:           "Ollama",
		DefaultBaseURL: "http://localhost:11434",
		LegacyType:     30,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
