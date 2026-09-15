package openai

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// tiktokenDownloadTimeout 下载 tiktoken 编码文件的超时时间。
// 在无法访问 openaipublic.blob.core.windows.net 的环境下，
// 避免 http.Get 无限阻塞导致服务无法启动。
const tiktokenDownloadTimeout = 30 * time.Second

// httpBpeLoader 实现 tiktoken.BpeLoader 接口，为编码文件下载增加超时与本地缓存。
type httpBpeLoader struct {
	client *http.Client
}

func newHTTPBpeLoader() *httpBpeLoader {
	return &httpBpeLoader{client: &http.Client{Timeout: tiktokenDownloadTimeout}}
}

func (l *httpBpeLoader) LoadTiktokenBpe(tiktokenBpeFile string) (map[string]int, error) {
	contents, err := l.readFileCached(tiktokenBpeFile)
	if err != nil {
		return nil, err
	}

	bpeRanks := make(map[string]int)
	for _, line := range strings.Split(string(contents), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}
		token, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, err
		}
		rank, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		bpeRanks[string(token)] = rank
	}
	return bpeRanks, nil
}

func (l *httpBpeLoader) readFileCached(blobpath string) ([]byte, error) {
	if !strings.HasPrefix(blobpath, "http://") && !strings.HasPrefix(blobpath, "https://") {
		return os.ReadFile(blobpath)
	}

	// 1) 编译期内嵌的 BPE 文件（默认路径，无需任何环境变量配置）。
	//    这是兜底，避免在无法访问 openaipublic.blob.core.windows.net
	//    的网络环境下因为 tiktoken 初始化失败导致 defaultTokenEncoder 为 nil。
	if data, ok := lookupEmbeddedBpe(blobpath); ok {
		return data, nil
	}

	// 2) 用户在 TIKTOKEN_CACHE_DIR 显式提供的 BPE 文件（可选覆盖）。
	cacheDir := os.Getenv("TIKTOKEN_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = os.Getenv("DATA_GYM_CACHE_DIR")
	}
	if cacheDir != "" {
		cacheKey := fmt.Sprintf("%x", sha1.Sum([]byte(blobpath)))
		cachePath := filepath.Join(cacheDir, cacheKey)
		if data, err := os.ReadFile(cachePath); err == nil {
			return data, nil
		}
	}

	// 3) HTTP 下载（最后的回退，依赖网络）。
	resp, err := l.client.Get(blobpath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download %s: status code %d", blobpath, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(cacheDir, os.ModePerm); err == nil {
		_ = os.WriteFile(filepath.Join(cacheDir, fmt.Sprintf("%x", sha1.Sum([]byte(blobpath)))), data, 0644)
	}
	return data, nil
}
