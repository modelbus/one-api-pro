package cloudflare

import "github.com/Leon-PanPan/one-api-pro/relay/schema"

type Request struct {
	Messages    []model.Message `json:"messages,omitempty"`
	Lora        string          `json:"lora,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Prompt      string          `json:"prompt,omitempty"`
	Raw         bool            `json:"raw,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
	Temperature *float64        `json:"temperature,omitempty"`
}
