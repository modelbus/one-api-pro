package deepseek

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "deepseek",
		Name:           "DeepSeek",
		DefaultBaseURL: "https://api.deepseek.com",
		LegacyType:     36,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
