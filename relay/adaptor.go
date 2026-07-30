package relay

import (
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/registry"

	// Register all adaptors via init()
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/ai360"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/aiproxy"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/ali"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/alibailian"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/anthropic"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/aws"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/azure"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/baichuan"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/baidu"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/baiduv2"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/cloudflare"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/cohere"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/coze"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/deepl"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/deepseek"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/doubao"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/gemini"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/geminiv2"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/groq"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/lingyiwanwu"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/minimax"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/mistral"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/moonshot"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/novita"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/ollama"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/openrouter"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/palm"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/proxy"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/replicate"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/siliconflow"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/stepfun"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/tencent"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/togetherai"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/vertexai"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/xai"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/xunfei"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/xunfeiv2"
	_ "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/zhipu"
)

func GetAdaptorByChannel(channelType int) adaptor.Adaptor {
	return registry.GetAdaptorByLegacyType(channelType)
}

func GetAdaptorByChannelID(channelID string) adaptor.Adaptor {
	return registry.GetAdaptor(channelID)
}

func GetAdaptor(apiType int) adaptor.Adaptor {
	return registry.GetAdaptorByLegacyType(apiType)
}
