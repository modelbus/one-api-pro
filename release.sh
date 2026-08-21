#!/usr/bin/env bash
#
# release.sh - 一键构建所有平台的可执行文件（不压缩，直接可用）
#
# 用法:
#   ./release.sh                      # 使用 VERSION 文件 / git tag 作为版本号
#   ./release.sh v0.1.0               # 指定版本号
#   ./release.sh v0.1.0 --skip-frontend  # 跳过前端构建（复用已有 web/build）
#
# 输出（与 .github/workflows/release.yml 保持一致）:
#   dist/one-api-pro-linux-amd64
#   dist/one-api-pro-linux-arm64
#   dist/one-api-pro-windows-amd64.exe
#   dist/one-api-pro-darwin-amd64
#   dist/one-api-pro-darwin-arm64

set -euo pipefail

# ---------- 参数解析 ----------
VERSION=""
SKIP_FRONTEND=false

for arg in "$@"; do
  case "$arg" in
    --skip-frontend)
      SKIP_FRONTEND=true
      ;;
    --*)
      echo "未知参数: $arg" >&2
      exit 1
      ;;
    *)
      VERSION="$arg"
      ;;
  esac
done

# ---------- 确定版本号 ----------
if [ -z "$VERSION" ]; then
  if [ -f VERSION ]; then
    VERSION="$(tr -d '[:space:]' < VERSION)"
  fi
fi
if [ -z "$VERSION" ]; then
  VERSION="$(git describe --tags --always 2>/dev/null || true)"
fi
if [ -z "$VERSION" ]; then
  VERSION="0.0.1"
fi
case "$VERSION" in
  v*) ;;
  *) VERSION="v$VERSION" ;;
esac

echo "==> 版本号: $VERSION"

# ---------- 依赖检查 ----------
for cmd in go node npm; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "错误: 缺少依赖 $cmd，请先安装" >&2
    exit 1
  fi
done

# ---------- 下载依赖 ----------
echo "==> go mod download"
go mod download

# ---------- 构建前端 ----------
if [ "$SKIP_FRONTEND" = true ]; then
  echo "==> 跳过前端构建（使用已有 web/build）"
else
  echo "==> 构建前端"
  (cd web && sh build.sh)
fi

# ---------- 交叉编译 ----------
LDFLAGS="-s -w -X github.com/modelbus/one-api-pro/common.Version=${VERSION}"
rm -rf dist
mkdir -p dist

echo "==> 构建 one-api-pro-linux-amd64 (linux/amd64)"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/one-api-pro-linux-amd64 .

echo "==> 构建 one-api-pro-linux-arm64 (linux/arm64)"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/one-api-pro-linux-arm64 .

echo "==> 构建 one-api-pro-windows-amd64.exe (windows/amd64)"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/one-api-pro-windows-amd64.exe .

echo "==> 构建 one-api-pro-darwin-amd64 (darwin/amd64)"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/one-api-pro-darwin-amd64 .

echo "==> 构建 one-api-pro-darwin-arm64 (darwin/arm64)"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/one-api-pro-darwin-arm64 .

# ---------- 输出 ----------
echo ""
echo "==> 构建完成:"
ls -lh dist/
