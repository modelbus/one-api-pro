---
title: 新增 Provider
description: "如何接入一个新的 LLM Provider。"
category: contribute
order: 4
---

# 新增 Provider

> 如何接入一个新的 LLM Provider。

新增 Provider 的核心是 **实现 `relay/adaptor` 接口 + 通过 `init()` 注册**。整个流程零数据库迁移；build 后立即出现在「新建渠道」中。

## 1. 文件结构

```text
relay/adaptor/provider/<provider_id>/
├── register.go        # init() 注册 ChannelMeta
├── adaptor.go        # 实现 Adaptor 接口
├── main.go           # Handler / StreamHandler（响应处理）
└── constants.go      # ModelList + 默认 token / 其它常量
```

> `relay/adaptor/openai/` 是 OpenAI 兼容协议基线；多数新 Provider 只需**继承**它（嵌入 `openai.Adaptor`）。

## 2. 最小实现

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
        LegacyType:     36, // 选个未占用的整数，与历史 type 表兼容
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

## 3. 三个关键决定

### 3.1 ID 与 LegacyType

- `ID` 是新渠道的**字符串唯一标识**（`channels.type` 字段）。
- `LegacyType` 是历史整型映射；新 Provider 可填任意未占用的正整数（`relay/channeltype/define.go` 是历史分配表，新 Provider 选一个不会被自动复用的即可）。
- `Name` 显示在 UI。
- `DefaultBaseURL` 渠道默认值。

### 3.2 是否完全 OpenAI 兼容

- **完全兼容**：像 DeepSeek、Groq、Mistral、`GrokYunwu` 等，**嵌入 `openai.Adaptor`** 即可，只在 `main.go` 写自己的 Usage 解析（如果字段不同）。
- **部分兼容**：保留 `Adaptor` 全部方法，自己实现 `ConvertRequest` / `DoResponse`。
- **完全不同**（如 Anthropic 单独的 `/v1/messages` 协议）：参考 `relay/adaptor/anthropic/`，重写全部方法。

### 3.3 流式响应

流式（`IsStream=true`）通常用 `bufio.Scanner` 读 SSE，每行解析 JSON 然后 `render.StringData(c, data)` 透传；`[DONE]` 哨兵行必须原样转发。

非流式：直接 `io.Copy(c.Writer, resp.Body)`，用 `convertDeepSeekUsage` 等提取 Usage。

## 4. 自注册如何被发现

`init()` 在包被导入时执行；只要新目录被以下任意位置引用，编译即可自动 link：

```go
// relay/adaptor/import.go（推荐显式 import，便于 grep）
import (
    _ "github.com/modelbus/one-api-pro/relay/adaptor/provider/deepseek"
    _ "github.com/modelbus/one-api-pro/relay/adaptor/provider/<your_provider>"
)
```

或隐式（被 `controller` / `monitor` 引用的某 provider 间接导入）。

> **不要**在 `relay/adaptor/openai/register.go::registerLegacy` 加你的 Provider：那是历史 OpenAI 兼容 alias 的专属位置，会与下游 ID 冲突。

## 5. 验证

```bash
# 后端构建
go build ./...

# 跑 channeltype 单测，确保 ID
go test ./relay/...

# 重启后到「新建渠道」页面，下拉里应该出现 <provider>
```

## 6. 提交

按 [提交规范 §3](/zh/contribute/commit-convention) 一文件一 commit，body 写明：

```text
feat(channel): 新增 <provider> 适配器

Add <provider> adaptor.

- 实现 Adaptor 接口：嵌入 openai.Adaptor，复用 ConvertRequest
- 自定义 DoResponse：解析 <provider> usage 字段
- 注册到 registry：ID="<id>" / LegacyType=<N>
- 常量 ModelList：["<model1>", "<model2>"]
```

## 7. 进阶

- **图片/Embedding/Audio**：若 Provider 支持，复用 `openai.Adaptor` 的 `ConvertImageRequest` 或自己实现。
- **特殊 Header**：如 Azure 需要 `api-key` 而不是 `Authorization: Bearer`，重写 `SetupRequestHeader`。
- **Balance Update**：若 Provider 有余额查询接口，在 `monitor/channel_balance.go` 注册一个 fetcher（一般不需要，开箱支持是 OpenAI 系）。

下一步 / Next: [新增支付通道](/zh/contribute/add-payment) · [发版流程](/zh/contribute/release-process)。
