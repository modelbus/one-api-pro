package ali

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
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
