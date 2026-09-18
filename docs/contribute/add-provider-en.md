---
title: Add a Provider
description: "How to onboard a new LLM provider."
category: contribute
order: 4
---

# Add a Provider

> How to onboard a new LLM provider.
> 如何接入一个新的 LLM Provider。

The core is **implement the `relay/adaptor` interface + register in `init()`**. Zero DB migration; appears in "New Channel" after rebuild.

新增 Provider 的核心是 **实现 `relay/adaptor` 接口 + 通过 `init()` 注册**。整个流程零数据库迁移；build 后立即出现在「新建渠道」中。

## 1. File Layout / 文件结构

```text
relay/adaptor/provider/<provider_id>/
├── register.go        # init() registers ChannelMeta
├── adaptor.go        # implements the Adaptor interface
├── main.go           # Handler / StreamHandler (response handling)
└── constants.go      # ModelList + default tokens / other constants
```

> `relay/adaptor/openai/` is the OpenAI-compatible baseline; most new providers only **embed** it.

## 2. Minimal Example / 最小实现

Reference `relay/adaptor/provider/deepseek/`:

参考 `relay/adaptor/provider/deepseek/`：

```go
// relay/adaptor/provider/deepseek/register.go
package deepseek

import (
    "github.com/modelbus/one-api-pro/relay/adaptor"
    "github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
    registry.Register(registry.ChannelMeta{
        ID:             "deepseek",
        Name:           "DeepSeek",
        DefaultBaseURL: "https://api.deepseek.com",
        LegacyType:     36, // pick an unused positive integer
    }, func() adaptor.Adaptor { return &Adaptor{} })
}
```

```go
// relay/adaptor/provider/deepseek/adaptor.go
package deepseek

import (
    "io"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/modelbus/one-api-pro/relay/adaptor"
    "github.com/modelbus/one-api-pro/relay/adaptor/openai"
    "github.com/modelbus/one-api-pro/relay/meta"
)

type Adaptor struct {
    openaiAdaptor openai.Adaptor
}

func (a *Adaptor) Init(meta *meta.Meta)                                                  {}
func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error)                          { return a.openaiAdaptor.GetRequestURL(meta) }
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, m *meta.Meta) error { return a.openaiAdaptor.SetupRequestHeader(c, req, m) }
func (a *Adaptor) ConvertRequest(c *gin.Context, mode int, req *model.GeneralOpenAIRequest) (any, error) {
    return a.openaiAdaptor.ConvertRequest(c, mode, req)
}
func (a *Adaptor) ConvertImageRequest(req *model.ImageRequest) (any, error) {
    return a.openaiAdaptor.ConvertImageRequest(req)
}
func (a *Adaptor) DoRequest(c *gin.Context, m *meta.Meta, body io.Reader) (*http.Response, error) {
    return adaptor.DoRequestHelper(a, c, m, body)
}
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, m *meta.Meta) (*model.Usage, *model.ErrorWithStatusCode) {
    if m.IsStream {
        err, _, usage := StreamHandler(c, resp)
        return usage, err
    }
    err, usage := Handler(c, resp, m.PromptTokens, m.ActualModelName)
    return usage, err
}
func (a *Adaptor) GetModelList() []string { return ModelList }
```

```go
// relay/adaptor/provider/deepseek/constants.go
package deepseek

var ModelList = []string{
    "deepseek-chat",
    "deepseek-reasoner",
}
```

## 3. Three Key Decisions / 三个关键决定

### 3.1 ID and LegacyType

- `ID` is the **string identifier** for the new channel (`channels.type` column).
- `LegacyType` is the historical integer mapping; pick any unused positive integer. `relay/channeltype/define.go` is the historical allocation table — pick something not auto-reused.
- `Name` is shown in the UI.
- `DefaultBaseURL` is the channel default.

### 3.2 OpenAI Compatibility / 是否完全 OpenAI 兼容

- **Fully compatible**: DeepSeek, Groq, Mistral, etc. — embed `openai.Adaptor` and only override `Usage` parsing if fields differ.
  完全兼容：DeepSeek、Groq、Mistral 等，**嵌入 `openai.Adaptor`** 即可，只在 `main.go` 写自己的 Usage 解析。
- **Partially compatible**: keep the full `Adaptor` surface and implement `ConvertRequest` / `DoResponse` yourself.
  部分兼容：保留 `Adaptor` 全部方法，自己实现 `ConvertRequest` / `DoResponse`。
- **Fully different protocol** (e.g. Anthropic's `/v1/messages`): see `relay/adaptor/anthropic/` for a full rewrite.
  完全不同（如 Anthropic 单独的 `/v1/messages` 协议）：参考 `relay/adaptor/anthropic/`。

### 3.3 Streaming / 流式响应

Streaming (`IsStream=true`) typically uses `bufio.Scanner` to read SSE; each line is JSON-parsed and forwarded verbatim via `render.StringData(c, data)`; the `[DONE]` sentinel must be forwarded as-is.

流式（`IsStream=true`）通常用 `bufio.Scanner` 读 SSE，每行解析 JSON 然后 `render.StringData(c, data)` 透传；`[DONE]` 哨兵行必须原样转发。

Non-stream: `io.Copy(c.Writer, resp.Body)` and extract Usage via `convertDeepSeekUsage` (or similar).

非流式：直接 `io.Copy(c.Writer, resp.Body)`，用 `convertDeepSeekUsage` 等提取 Usage。

## 4. How the Registration is Discovered / 自注册如何被发现

`init()` runs when the package is imported. As long as the new directory is reachable from any of these, the binary will link it automatically:

`init()` 在包被导入时执行；只要新目录被以下任意位置引用，编译即可自动 link：

```go
// relay/adaptor/import.go (recommended: explicit import for grep-ability)
import (
    _ "github.com/modelbus/one-api-pro/relay/adaptor/provider/deepseek"
    _ "github.com/modelbus/one-api-pro/relay/adaptor/provider/<your_provider>"
)
```

Or indirectly — if any imported provider transitively imports it.

或隐式（被 `controller` / `monitor` 引用的某 provider 间接导入）。

> **Don't** add your provider to `relay/adaptor/openai/register.go::registerLegacy` — that file is for legacy aliases and will collide with the new ID.
> **不要**在 `relay/adaptor/openai/register.go::registerLegacy` 加你的 Provider：那是历史 OpenAI 兼容 alias 的专属位置。

## 5. Verify / 验证

```bash
# Backend build
go build ./...

# Run the channeltype tests to ensure no ID / LegacyType collision
go test ./relay/...

# After restart, "New Channel" dropdown should include <provider>
```

## 6. Commit / 提交

Per [Commit Convention §3](/en/contribute/commit-convention), one commit per file, body includes:

按 [提交规范 §3](/zh/contribute/commit-convention) 一文件一 commit，body 写明：

```text
feat(channel): add <provider> adaptor

Add <provider> adaptor.

- Implement Adaptor: embed openai.Adaptor, reuse ConvertRequest
- Custom DoResponse: parse <provider> usage fields
- Register with registry: ID="<id>" / LegacyType=<N>
- Constants ModelList: ["<model1>", "<model2>"]
```

## 7. Advanced / 进阶

- **Image / Embedding / Audio**: reuse `openai.Adaptor`'s `ConvertImageRequest` if compatible.
  图片/Embedding/Audio：若 Provider 支持，复用 `openai.Adaptor` 的 `ConvertImageRequest`。
- **Special Headers**: e.g. Azure uses `api-key` instead of `Authorization: Bearer`; override `SetupRequestHeader`.
  特殊 Header：如 Azure 需要 `api-key`，重写 `SetupRequestHeader`。
- **Balance Update**: if the provider exposes a balance endpoint, register a fetcher in `monitor/channel_balance.go` (typically not needed; OpenAI-compatible defaults work).
  Balance Update：在 `monitor/channel_balance.go` 注册 fetcher（一般不需要）。

Next: [Add a Payment Channel](/en/contribute/add-payment) · [Release Process](/en/contribute/release-process).