# syntax=docker/dockerfile:1.7

# ============ Stage 1: 前端构建产物由 CI 预构建 ============
# 关键改动：v0.0.16 之前 Stage 1 在 buildx 多架构构建里跑 npm install。
# 在 QEMU arm64 emulated 下 V8 生成的原子指令触发 SIGILL(exit 132)。
# 改为由 CI workflow 在原生 amd64 runner 预构建 web/build/default-pro,
# Dockerfile 仅 COPY 已构建产物;arm64 构建不再跑 npm。
#
# web/build 由 release-docker.yml 中的 "Pre-build web/default-pro" step
# 产出并通过 GHA cache 跨 workflow 复用。
# 版本: v0.0.16
FROM golang:1.22-alpine AS go-builder
WORKDIR /build

# go.mod 声明 1.25.0，由 GOTOOLCHAIN=auto 自动下载 1.25 toolchain
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

# 注意:web/build 必须在 build context 中存在,由 CI workflow 预构建并 cache
COPY . .

ARG VERSION="dev"
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath \
      -ldflags "-s -w -X 'github.com/modelbus/one-api-pro/common.Version=${VERSION}'" \
      -o /out/one-api-pro


# ============ Stage 3: 运行时镜像 ============
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata wget \
 && addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=go-builder /out/one-api-pro /app/one-api-pro
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

RUN mkdir -p /app/data /app/config \
 && chown -R app:app /app

VOLUME ["/app/config", "/app/data"]

USER app
WORKDIR /app/data

ENV PORT=3000 \
    LOG_DIR=/app/data/logs \
    SQLITE_PATH=/app/data/one-api-pro.db \
    CONFIG_DIR=/app/config \
    TZ=UTC

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- "http://localhost:${PORT}/api/status" || exit 1

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["one-api-pro"]
