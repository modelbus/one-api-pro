---
title: Reverse Proxy
description: "Nginx / Caddy / Cloudflare example configs."
category: install
order: 8
---

# Reverse Proxy

> Nginx / Caddy / Cloudflare example configs.

The examples below assume One API Pro listens on `127.0.0.1:3000` and is published as `https://api.example.com`.
All three examples do the same three things:

1. TLS termination + HSTS.
2. Pass-through for streaming endpoints (WebSocket / `Transfer-Encoding: chunked`).
3. Real client IP forwarding (`X-Real-IP` / `X-Forwarded-For`) so rate limiting and audit logs work.

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

    # Global rate limit (per IP, 20 r/s burst 40)
    limit_req_zone $binary_remote_addr zone=oneapi:10m rate=20r/s;

    location / {
        limit_req zone=oneapi burst=40 nodelay;

        proxy_pass         http://one_api_pro;
        proxy_http_version 1.1;

        # Real client IP
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Streaming / WebSocket pass-through
        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        "upgrade";
        proxy_buffering off;
        proxy_cache    off;
        proxy_read_timeout  300s;
        proxy_send_timeout  300s;
    }
}
```

Highlights:

- `proxy_buffering off` + `proxy_cache off` prevent buffering of streaming responses.
- `Upgrade` / `Connection: upgrade` are required for WebSocket-style streams on `/v1/chat/completions`.
- `limit_req_zone` and `GLOBAL_*_RATE_LIMIT` are two independent limits; keep Nginx as the DDoS backstop.

## Caddy

`Caddyfile`:

```caddy
api.example.com {
    encode zstd gzip
    reverse_proxy 127.0.0.1:3000 {
        header_up Host {host}
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}
        # Streaming / WebSocket pass-through
        flush_interval -1
        transport http {
            dial_timeout 5s
            response_header_timeout 60s
        }
    }

    # Caddy auto-issues certificates; supply your own with `tls <cert> <key>` if needed.
}
```

`flush_interval -1` flushes each chunk immediately for streaming responses.

## Cloudflare

Point `api.example.com` DNS at your One API Pro origin (the orange-cloud proxy is recommended) and pick **Full (Strict)** under **SSL/TLS → Overview**.

Cloudflare already injects `X-Forwarded-For`; tell Nginx / Caddy to trust the Cloudflare ranges:

```nginx
# Nginx: trust Cloudflare ranges
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

real_ip_header CF-Connecting-IP;   # when the CF-Connecting-IP header is enabled
# or:
# real_ip_header X-Forwarded-For;
```

> Cloudflare IP ranges change over time; check <https://www.cloudflare.com/ips/> before deploying.

## Multi-node load balancing

For a decentralized cluster, prefer `ip_hash` so the same client is pinned to the same node — this keeps session and plan rate-limits consistent:

```nginx
upstream one_api_cluster {
    ip_hash;
    server cn.example.com:3000;
    server us.example.com:3000;
    server eu.example.com:3000;
    keepalive 32;
}
```

> `ip_hash` is critical — it pins a client to a single node so plan rate-limits and Redis cache stay consistent.

See [Multi-node Deployment](/en/decentralization/deployment) for details.

Next: [Configuration](/en/install/config) · [Cluster Overview](/en/decentralization/overview).
