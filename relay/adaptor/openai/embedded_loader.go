// embedded_loader.go 把 tiktoken BPE 编码文件嵌入二进制，
// 避免在无法访问 openaipublic.blob.core.windows.net 的网络环境下，
// 因为 tiktoken 初始化失败导致 defaultTokenEncoder 为 nil，
// 进一步被 SyncOptions 的 ApproximateTokenEnabled=false 覆盖后 panic。
//
// BPE 文件是二进制 vendor 资产，不参与编译产物（不导出 Go 符号），
// 仅通过 lookupEmbeddedBpe() 在运行时按需读取。
//
// 版本: v0.0.21
// 日期: 2026-09-15
// 作者: opencode

package openai

import (
	"crypto/sha1"
	"embed"
	"fmt"
)

//go:embed tiktoken_cache/*
var embeddedBpeFS embed.FS

// lookupEmbeddedBpe 按 tiktoken-go 的 cacheKey 算法（sha1(blobpath)）查找
// 内嵌的 BPE 文件。返回 (nil, false) 表示未嵌入对应编码。
//
// 调用方需要自行负责把 []byte 解析成 map[string]int（见 loadTiktokenBpe）。
func lookupEmbeddedBpe(blobpath string) ([]byte, bool) {
	cacheKey := fmt.Sprintf("%x", sha1.Sum([]byte(blobpath)))
	data, err := embeddedBpeFS.ReadFile("tiktoken_cache/" + cacheKey)
	if err != nil {
		return nil, false
	}
	return data, true
}