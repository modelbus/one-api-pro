package tencent

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "tencent",
		Name:           "Tencent",
		DefaultBaseURL: "https://hunyuan.tencentcloudapi.com",
		LegacyType:     23,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
