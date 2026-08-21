package minimax

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "minimax",
		Name:           "Minimax",
		DefaultBaseURL: "https://api.minimax.chat",
		LegacyType:     27,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
