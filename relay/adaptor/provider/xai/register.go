package xai

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "xai",
		Name:           "xAI",
		DefaultBaseURL: "https://api.x.ai",
		LegacyType:     45,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
