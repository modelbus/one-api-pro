package gemini

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "gemini",
		Name:           "Gemini",
		DefaultBaseURL: "https://generativelanguage.googleapis.com",
		LegacyType:     24,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
