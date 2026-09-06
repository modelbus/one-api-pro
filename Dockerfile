# syntax=docker/dockerfile:1.7

# ============ Stage 1: 前端构建（仅 default-pro） ============
FROM node:20-alpine AS web-builder
WORKDIR /src

COPY web/THEMES      ./web/THEMES
COPY web/default-pro ./web/default-pro

RUN while IFS= read -r theme; do \
      echo "==> npm install ($theme)" && \
      npm install --prefix ./web/"$theme" && \
      echo "==> npm run build ($theme)" && \
      npm run build --prefix ./web/"$theme"; \
    done < web/THEMES


# ============ Stage 2: Go 二进制构建 ============
FROM golang:1.22-alpine AS go-builder
WORKDIR /build

# go.mod 声明 1.25.0，由 GOTOOLCHAIN=auto 自动下载 1.25 toolchain
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=web-builder /src/web/build ./web/build

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
