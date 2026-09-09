// usage.go Claude usage 归一化与累计合并工具
// Usage normalization & cumulative merge helpers for the Anthropic adaptor.
// 版本: v0.0.15
// 日期: 2026-09-09
// 作者: opencode

package anthropic

import (
	"github.com/modelbus/one-api-pro/relay/schema"
)

// ClaudeUsage2OpenAI 将 Anthropic 互斥的 usage 字段归一化为 OpenAI「cached ⊆ prompt」语义。
// Claude 的 input_tokens / cache_read_input_tokens / cache_creation_input_tokens 三者互斥且不含彼此，
// 而共享计费公式 input*(prompt-cached)+cachedPrice*cached 隐含 cached ⊂ prompt；
// 因此把 read/creation 折入 PromptTokens，使公式能还原出 input×input + cachedPrice×read 的正确金额。
// cache_creation 按 input 价计费是已知近似（真实写入价约 1.25×input），见 #13。
// 版本: v0.0.15
// 日期: 2026-09-09
// 作者: opencode
func ClaudeUsage2OpenAI(u Usage) model.Usage {
	prompt := u.InputTokens + u.CacheReadTokens + u.CacheCreationTokens
	usage := model.Usage{
		PromptTokens:     prompt,
		CompletionTokens: u.OutputTokens,
		TotalTokens:      prompt + u.OutputTokens,
	}
	if u.CacheReadTokens > 0 {
		usage.PromptTokensDetails = &model.PromptTokensDetails{
			CachedTokens: u.CacheReadTokens,
		}
	}
	return usage
}

// MergeClaudeUsage 将一个携带累计 usage 的 Claude 流式事件合并进 acc。
// message_start 与 message_delta 都携带 usage，且新版 API 下均为整段请求的累计值
// （message_delta 重复 message_start 的 input/cache 字段并带上最终 output_tokens），
// 因此不能按 += 增量累加（会双计 input/cache），需对各字段取 max（累计值单调递增）。
// 兼容旧形态：delta 仅带 output 或 input 字段为 0 时，max 仍保留 message_start 的 input。
// 版本: v0.0.15
// 日期: 2026-09-09
// 作者: opencode
func MergeClaudeUsage(acc *model.Usage, u Usage) {
	cur := ClaudeUsage2OpenAI(u)
	if cur.PromptTokens > acc.PromptTokens {
		acc.PromptTokens = cur.PromptTokens
	}
	if cur.CompletionTokens > acc.CompletionTokens {
		acc.CompletionTokens = cur.CompletionTokens
	}
	if cur.TotalTokens > acc.TotalTokens {
		acc.TotalTokens = cur.TotalTokens
	}
	if cur.PromptTokensDetails != nil {
		if acc.PromptTokensDetails == nil {
			acc.PromptTokensDetails = &model.PromptTokensDetails{}
		}
		if cur.PromptTokensDetails.CachedTokens > acc.PromptTokensDetails.CachedTokens {
			acc.PromptTokensDetails.CachedTokens = cur.PromptTokensDetails.CachedTokens
		}
	}
}
