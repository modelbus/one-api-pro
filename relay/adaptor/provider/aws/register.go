package aws

import (
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "aws_claude",
		Name:           "AWS Claude",
		DefaultBaseURL: "",
		LegacyType:     33,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})
}
