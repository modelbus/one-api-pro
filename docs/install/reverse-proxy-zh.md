---
title: 反向代理
description: Nginx / Caddy / Cloudflare 配置示例。
category: install
order: 8
---

# 反向代理

> 为什么要用反向代理、配置时要做对的三件事、各家示例。

## 为什么要用反向代理

直接暴露容器端口到公网不安全。反向代理负责：

- **TLS 终止**：把 HTTPS 解密，反代到容器内 HTTP
- **真实 IP 透传**：让 One API Pro 看到客户端真实 IP（限流 / 审计）
- **HSTS / 安全头**：HTTP Strict Transport Security 等
- **流式透传**：`/v1/chat/completions` 是 SSE 流式响应，不能被反代缓存

## 配置必做三件事

下面所有示例都做到：

1. TLS 终止与 HSTS
2. 流式接口所需的 WebSocket / `Transfer-Encoding: chunked` 透传
3. 真实客户端 IP 透传（`X-Real-IP` / `X-Forwarded-For`）

## Nginx

```nginx
server {
  listen 443 ssl http2;
  server_name api.example.com;

  ssl_certificate     /etc/letsencrypt/live/api.example.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

  add_header Strict-Transport-Security "max-age=63072000" always;

  # 真实客户端 IP
  real_ip_header X-Forwarded-For;
  set_real_ip_from 0.0.0.0/0;
  real_ip_recursive on;

  location / {
    proxy_pass http://127.0.0.1:3000;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    # 流式必须
    proxy_buffering off;
    proxy_cache off;
    proxy_read_timeout 300s;
    chunked_transfer_encoding on;
  }
}
```

## Caddy

`Caddyfile`：

```
api.example.com {
  reverse_proxy 127.0.0.1:3000 {
    header_up X-Real-IP {remote_host}
    header_up X-Forwarded-For {remote_host}
    header_up X-Forwarded-Proto {scheme}
    # 流式
    flush_interval -1
  }
}
```

Caddy 自动申请 Let's Encrypt 证书，自动 HSTS，无需手写。

## Cloudflare

如果你用 Cloudflare 做 DNS + CDN：

1. 后台 → SSL/TLS → Full (strict)
2. 后台 → SSL/TLS → Edge Certificates → Always Use HTTPS → On
3. 后台 → Rules → Transform Rules → 添加 Host header / X-Forwarded-For
4. Cloudflare 的 IP 段通过 `set_real_ip_from` 加到 Nginx 配置里（[Cloudflare IP 段](https://www.cloudflare.com/ips/)）

注意：Cloudflare 默认会缓存 HTML 页面，但 `/v1/*` 不缓存即可（默认 Cache Level = Standard 已避免）。

## 验证

部署后：

- [ ] `curl -I https://api.example.com/` 看到 200 + `Strict-Transport-Security`
- [ ] `curl -H "X-Forwarded-For: 1.2.3.4" https://api.example.com/api/log/self` 在后台日志里看到 IP 是 `1.2.3.4`
- [ ] 流式调用：`curl -N https://api.example.com/v1/chat/completions ...` 能流式收到 SSE 块

## 相关文档

- [Docker 单实例部署](./docker-deploy)
- [系统要求](./requirements)