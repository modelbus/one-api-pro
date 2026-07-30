package openai

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"
)

func init() {
	registry.Register(registry.ChannelMeta{
		ID:             "openai",
		Name:           "OpenAI",
		DefaultBaseURL: "https://api.openai.com",
		LegacyType:     1,
	}, func() adaptor.Adaptor {
		return &Adaptor{}
	})

	registerLegacy()
}

func registerLegacy() {
	// OpenAI-compatible providers without their own adaptor packages
	legacy := []struct {
		id         string
		name       string
		baseURL    string
		legacyType int
	}{
		{"api2d", "API2D", "https://oa.api2d.net", 2},
		{"closeai", "CloseAI", "https://api.closeai-proxy.xyz", 4},
		{"openai-sb", "OpenAI SB", "https://api.openai-sb.com", 5},
		{"openaimax", "OpenAI Max", "https://api.openaimax.com", 6},
		{"ohmygpt", "OhMyGPT", "https://api.ohmygpt.com", 7},
		{"custom", "Custom", "", 8},
		{"ails", "Ails", "https://api.caipacity.com", 9},
		{"aiproxy", "AIProxy", "https://api.aiproxy.io", 10},
		{"api2gpt", "API2GPT", "https://api.api2gpt.com", 12},
		{"aigc2d", "AIGC2D", "https://api.aigc2d.com", 13},
		{"fastgpt", "FastGPT", "https://fastgpt.run/api/openapi", 22},
		{"openai_compatible", "OpenAI Compatible", "", 50},
	}

	factory := func() adaptor.Adaptor {
		return &Adaptor{}
	}

	for _, l := range legacy {
		registry.Register(registry.ChannelMeta{
			ID:             l.id,
			Name:           l.name,
			DefaultBaseURL: l.baseURL,
			LegacyType:     l.legacyType,
		}, factory)
	}
}
