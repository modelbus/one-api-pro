package moonshot

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "moonshot",
		Name:           "Moonshot",
		DefaultBaseURL: "https://api.moonshot.cn",
		LegacyType:     25,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
