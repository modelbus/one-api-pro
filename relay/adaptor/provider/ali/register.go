package ali

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "ali",
		Name:           "Ali",
		DefaultBaseURL: "https://dashscope.aliyuncs.com",
		LegacyType:     17,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
