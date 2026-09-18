---
title: 反向代理
description: "Nginx / Caddy / Cloudflare 配置示例。"
category: install
order: 8
---

# 反向代理

> Nginx / Caddy / Cloudflare 配置示例。

下面示例假设 One API Pro 监听 `127.0.0.1:3000`，对外暴露 `https://api.example.com`。
所有示例都做了三件事：

1. TLS 终止与 HSTS；
2. 流式接口所需的 WebSocket / `Transfer-Encoding: chunked` 透传；
3. 真实客户端 IP 透传（`X-Real-IP` / `X-Forwarded-For`），便于限流与审计生效。

## Nginx

```nginx
upstream one_api_pro {
    server 127.0.0.1:3000;
    keepalive 32;
}

server {
    listen 80;
    server_name api.example.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate     /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    client_max_body_size 50m;

    # 全局限速（每 IP 每秒 20 次突发）
    limit_req_zone $binary_remote_addr zone=oneapi:10m rate=20r/s;

    location / {
        limit_req zone=oneapi burst=40 nodelay;

        proxy_pass         http://one_api_pro;
        proxy_http_version 1.1;

        # 真实 IP
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 流式 / WebSocket 兼容
        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        "upgrade";
        proxy_buffering off;
        proxy_cache    off;
        proxy_read_timeout  300s;
        proxy_send_timeout  300s;
    }
}
```

要点 / Highlights:

- `proxy_buffering off` 与 `proxy_cache off` 避免流式响应被缓冲；
- `Upgrade` / `Connection: upgrade` 头对 `/v1/chat/completions` 流式响应可选启用 WebSocket；
- `limit_req_zone` 与 `GLOBAL_*_RATE_LIMIT` 是两层独立限流，建议至少保留 Nginx 一层做 DDoS 兜底；

## Caddy

`Caddyfile`：

```caddy
api.example.com {
    encode zstd gzip
    reverse_proxy 127.0.0.1:3000 {
        header_up Host {host}
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}
        # 流式 / WebSocket
        flush_interval -1
        transport http {
            dial_timeout 5s
            response_header_timeout 60s
        }
    }

    # Caddy 自动签发证书；如使用自有证书可加 tls <cert> <key>
}
```

`flush_interval -1` 让 Caddy 立即转发每个 chunk，配合流式响应。

## Cloudflare

将 `api.example.com` 的 DNS 指向 One API Pro 的公网入口（推荐开启 Cloudflare 代理小黄云），并在 **SSL/TLS → Overview** 选择 **Full (Strict)** 模式。

Cloudflare 代理已自动注入 `X-Forwarded-For`；Nginx / Caddy 侧需信任 Cloudflare 网段：

```nginx
# Nginx: set_real_ip_from 信任 Cloudflare 网段
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 131.0.72.0/22;
# IPv6
set_real_ip_from 2400:cb00::/32;
set_real_ip_from 2606:4700::/32;
set_real_ip_from 2803:f800::/32;
set_real_ip_from 2405:b500::/32;
set_real_ip_from 2405:8100::/32;
set_real_ip_from 2a06:98c0::/29;
set_real_ip_from 2c0f:f248::/32;

real_ip_header CF-Connecting-IP;   # 当开启 CF-Connecting-IP 头时
# 或者
# real_ip_header X-Forwarded-For;
```

> Cloudflare 网段会变动，部署前请以 <https://www.cloudflare.com/ips/> 为准。

## 多节点负载均衡

去中心化集群下推荐 `ip_hash`，将同一客户端锁定到同一节点，避免会话与计划配额被切碎：

```nginx
upstream one_api_cluster {
    ip_hash;
    server cn.example.com:3000;
    server us.example.com:3000;
    server eu.example.com:3000;
    keepalive 32;
}
```

> `ip_hash` 是关键 — 它把同一客户端绑定到同一节点，让会话与 Redis 缓存命中一致。

详见 [多节点部署](/zh/decentralization/deployment)。
See [Multi-node Deployment](/en/decentralization/deployment) for details.

下一步 / Next: [配置项与环境变量](/zh/install/config) · [Cluster 概览](/zh/decentralization/overview)。

