package geminiv2

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "gemini_openai_compatible",
		Name:           "Gemini OpenAI Compatible",
		DefaultBaseURL: "https://generativelanguage.googleapis.com/v1beta/openai/",
		LegacyType:     51,
	}, func() adaptor.Adaptor {
		return &Adaptor{Adaptor: &openai.Adaptor{}}
	})
}
