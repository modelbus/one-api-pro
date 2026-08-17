package common

import (
	"os"
	"path/filepath"
	"strings"
)

// loadVersionFromFile 从 VERSION 文件读取版本号并覆盖编译期默认值。
// 优先读取当前工作目录下的 VERSION 文件，其次读取可执行文件同目录下的 VERSION 文件。
// 兼容有无 "v" 前缀两种写法：无前缀时自动补 "v"。
// 文件不存在或内容为空时静默跳过，保留编译期（或 ldflags 注入）的版本号。
func loadVersionFromFile() {
	candidates := []string{"VERSION"}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "VERSION"))
	}

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		v := strings.TrimSpace(string(data))
		if v == "" {
			continue
		}
		if !strings.HasPrefix(v, "v") {
			v = "v" + v
		}
		Version = v
		return
	}
}
