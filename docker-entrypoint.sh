#!/bin/sh
# docker-entrypoint.sh - one-api-pro 容器启动入口
#
# 职责：
#   1. 自动加载 $CONFIG_DIR/.env（如果存在）
#   2. 透传用户自定义 CLI 参数
#   3. exec 替换为子进程，保证 SIGTERM 正确传递给 one-api-pro
#
# 支持的调用方式：
#   docker run image                              # 默认行为：自动加载 /app/config/.env
#   docker run image --port 8080 --log-dir /xxx   # 透传 CLI 参数
#   docker run image bash                         # 调试模式：进入 shell，不启动主程序
#   docker run -e CONFIG_DIR=/custom image        # 自定义配置目录

set -e

CONFIG_DIR="${CONFIG_DIR:-/app/config}"
ENV_FILE="$CONFIG_DIR/.env"
BINARY="/app/one-api-pro"

echo "==> CONFIG_DIR=$CONFIG_DIR"
echo "==> USER=$(whoami)"

# 解析 CMD：默认 CMD 是 ["one-api-pro"]，需要跳过这个虚拟参数
SKIP_BINARY=0
case "${1:-}" in
  one-api-pro|/app/one-api-pro)
    SKIP_BINARY=1
    shift
    ;;
esac

# 如果第一个参数不是 one-api-pro 也不是它的路径，说明用户传了其他命令
# 比如 bash / sh / ls -la，用于调试
if [ "$SKIP_BINARY" -eq 0 ] && [ $# -gt 0 ]; then
  exec "$@"
fi

# 收集传递给 one-api-pro 的额外 CLI 参数
EXTRA_ARGS="$*"

# 1) 自动加载 .env
ARGS=""
if [ -f "$ENV_FILE" ]; then
  echo "==> Loading env file: $ENV_FILE"
  ARGS="--env $ENV_FILE"
else
  echo "==> No .env file at $ENV_FILE (only -e environment variables will be used)"
fi

# 2) 透传用户 CLI 参数
if [ -n "$EXTRA_ARGS" ]; then
  ARGS="$ARGS $EXTRA_ARGS"
fi

# 启动
VERSION=$($BINARY --version 2>/dev/null || echo "unknown")
echo "==> one-api-pro version: $VERSION"
echo "==> Starting: $BINARY $ARGS"

exec $BINARY $ARGS
