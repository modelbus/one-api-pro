package lingyiwanwu

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "lingyiwanwu",
		Name:           "LingYiWanWu",
		DefaultBaseURL: "https://api.lingyiwanwu.com",
		LegacyType:     31,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
