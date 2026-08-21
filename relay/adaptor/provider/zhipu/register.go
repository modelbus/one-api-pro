package zhipu

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "zhipu",
		Name:           "Zhipu",
		DefaultBaseURL: "https://open.bigmodel.cn",
		LegacyType:     16,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
