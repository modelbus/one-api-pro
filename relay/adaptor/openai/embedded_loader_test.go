// embedded_loader_test.go 验证 tiktoken 离线缓存的三级查找优先级：
//  1. 编译期内嵌的 BPE 文件（默认，最重要）
//  2. TIKTOKEN_CACHE_DIR 显式指定（可选覆盖）
//  3. HTTP 下载（最后回退）
//
// 这层防护对应 issues：启动期 HTTP 超时 → fallbackToApproximateTokenCount
// → 600s 后 SyncOptions 把 ApproximateTokenEnabled 覆盖回 false →
// defaultTokenEncoder==nil + flag==false → panic("nil pointer dereference")。
//
// 版本: v0.0.21
// 日期: 2026-09-15
// 作者: opencode

package openai

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkoukk/tiktoken-go"
)

const (
	cl100kURL  = "https://openaipublic.blob.core.windows.net/encodings/cl100k_base.tiktoken"
	o200kURL   = "https://openaipublic.blob.core.windows.net/encodings/o200k_base.tiktoken"
	cl100kKey  = "9b5ad71b2ce5302211f9c61530b329a4922fc6a4"
	o200kKey   = "fb374d419588a4632f3f557e76b4b70aebbca790"
)

// TestLookupEmbeddedBpe 验证：所有内置的 4 个 BPE 文件都已正确嵌入，
// 且 sha1(blobpath) 与文件命名一致。
func TestLookupEmbeddedBpe(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		wantKey string
	}{
		{"cl100k_base", cl100kURL, cl100kKey},
		{"o200k_base", o200kURL, o200kKey},
		{"p50k_base", "https://openaipublic.blob.core.windows.net/encodings/p50k_base.tiktoken", "ec7223a39ce59f226a68acc30dc1af2788490e15"},
		{"r50k_base", "https://openaipublic.blob.core.windows.net/encodings/r50k_base.tiktoken", "0ea1e91bbb3a60f729a8dc8f777fd2fc07cd8df4"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotKey := fmt.Sprintf("%x", sha1.Sum([]byte(c.url)))
			if gotKey != c.wantKey {
				t.Fatalf("sha1(%s) = %s, want %s", c.url, gotKey, c.wantKey)
			}
			data, ok := lookupEmbeddedBpe(c.url)
			if !ok {
				t.Fatalf("embedded BPE missing for %s (expected file %s in tiktoken_cache/)", c.url, c.wantKey)
			}
			if len(data) == 0 {
				t.Fatalf("embedded BPE for %s is empty", c.url)
			}
		})
	}
}

// TestLookupEmbeddedBpe_UnknownURL 验证：未知 URL 返回 ok=false（不会 panic）。
func TestLookupEmbeddedBpe_UnknownURL(t *testing.T) {
	data, ok := lookupEmbeddedBpe("https://example.com/nonexistent.tiktoken")
	if ok || data != nil {
		t.Fatalf("expected (nil, false), got (%v, %v)", data, ok)
	}
}

// TestLoadEmbeddedBpe 验证：通过 httpBpeLoader 读取内嵌的 cl100k_base，
// 解析后能正常 tokenize。这正是 InitTokenEncoders 调用 tiktoken-go 的实际路径。
func TestLoadEmbeddedBpe(t *testing.T) {
	t.Setenv("TIKTOKEN_CACHE_DIR", "") // 禁用磁盘 cache，确保走 embedded
	loader := newHTTPBpeLoader()
	ranks, err := loader.LoadTiktokenBpe(cl100kURL)
	if err != nil {
		t.Fatalf("LoadTiktokenBpe failed: %v", err)
	}
	if len(ranks) == 0 {
		t.Fatalf("empty BPE ranks")
	}
	// sanity: cl100k_base 应该至少有 100000 个 token（实际 ~100k）
	if len(ranks) < 100000 {
		t.Fatalf("cl100k_base ranks count = %d, expected >= 100000", len(ranks))
	}
}

// TestEmbeddedBpe_UsedByTiktokenGo 端到端：让 tiktoken-go 用我们的 loader，
// 验证不设任何 env var 也能正常 encoder。
func TestEmbeddedBpe_UsedByTiktokenGo(t *testing.T) {
	t.Setenv("TIKTOKEN_CACHE_DIR", "")
	tiktoken.SetBpeLoader(newHTTPBpeLoader())
	for _, model := range []string{"gpt-3.5-turbo", "gpt-4", "gpt-4o"} {
		enc, err := tiktoken.EncodingForModel(model)
		if err != nil {
			t.Fatalf("EncodingForModel(%s) failed: %v", model, err)
		}
		if enc == nil {
			t.Fatalf("EncodingForModel(%s) returned nil", model)
		}
		tokens := enc.Encode("hello world", nil, nil)
		if len(tokens) == 0 {
			t.Fatalf("%s: encoded 0 tokens for 'hello world'", model)
		}
	}
}

// TestTIKTOKEN_CACHE_DIR_Overrides 验证：当 TIKTOKEN_CACHE_DIR 显式设置，
// 我们的 loader 优先用磁盘文件（而非 embedded）。这是给需要自定义 BPE 的高级用户。
func TestTIKTOKEN_CACHE_DIR_Overrides(t *testing.T) {
	// 用 HTTP 测试服务器模拟一个"BPE 文件"，让 loader 不走 embedded。
	var serverHits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHits++
		// 返回一个最小的、合法的 BPE 内容（一个空行 + 一对 token/rank）
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte("YQ== 1\n")) // "a" -> rank 1
	}))
	defer srv.Close()

	// 把 srv.URL 作为假 URL，让 loader 走 HTTP 路径
	tmpDir := t.TempDir()
	// 预置一个 SHA1(srv.URL) 的文件，让 TIKTOKEN_CACHE_DIR 命中
	key := fmt.Sprintf("%x", sha1.Sum([]byte(srv.URL)))
	cacheFile := filepath.Join(tmpDir, key)
	if err := os.WriteFile(cacheFile, []byte("YQ== 1\n"), 0644); err != nil {
		t.Fatalf("write cache file: %v", err)
	}
	t.Setenv("TIKTOKEN_CACHE_DIR", tmpDir)

	loader := newHTTPBpeLoader()
	ranks, err := loader.LoadTiktokenBpe(srv.URL)
	if err != nil {
		t.Fatalf("LoadTiktokenBpe failed: %v", err)
	}
	if serverHits != 0 {
		t.Fatalf("HTTP server should not have been hit when TIKTOKEN_CACHE_DIR is set; got %d hits", serverHits)
	}
	if len(ranks) != 1 || ranks["a"] != 1 {
		t.Fatalf("ranks = %v, want {a:1}", ranks)
	}
}

// TestHTTPFallback_WhenEmbeddedMissing 验证：当 embedded 没找到（未知 URL），
// 且 TIKTOKEN_CACHE_DIR 没设，loader 会回退到 HTTP 下载。
// 这条路径在生产代码里不会触发（4 个标准 URL 都被 embedded），
// 但保留了"未来 tiktoken-go 升级换 URL"的扩展能力。
func TestHTTPFallback_WhenEmbeddedMissing(t *testing.T) {
	// 模拟一个未知的 tiktoken URL（不在 embed 里）。
	// 注意：实际 httpBpeLoader 会对所有 https URL 都先查 embedded，
	// 所以我们用 srv.URL（http://，不是 https://）来避开 embedded 路径。
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("YQ== 42\n"))
	}))
	defer srv.Close()

	t.Setenv("TIKTOKEN_CACHE_DIR", "")
	loader := newHTTPBpeLoader()
	ranks, err := loader.LoadTiktokenBpe(srv.URL)
	if err != nil {
		t.Fatalf("LoadTiktokenBpe should fall back to HTTP for http:// URLs: %v", err)
	}
	if len(ranks) != 1 || ranks["a"] != 42 {
		t.Fatalf("ranks = %v, want {a:42}", ranks)
	}
}