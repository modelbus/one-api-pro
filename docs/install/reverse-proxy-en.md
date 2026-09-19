---
title: Reverse Proxy
description: Nginx / Caddy / Cloudflare configuration examples.
category: install
order: 8
---

# Reverse Proxy

> Why use one, the three things every reverse-proxy config must do, examples.

## Why use a reverse proxy

Exposing the container port directly is unsafe. The reverse proxy handles:

- **TLS termination**: decrypt HTTPS, proxy plain HTTP to the container.
- **Real client IP**: let One API Pro see the real IP (rate limiting / auditing).
- **HSTS / security headers**.
- **Streaming passthrough**: `/v1/chat/completions` is SSE; do NOT buffer.

## Three things every config must do

All examples below do these:

1. TLS termination + HSTS
2. Streaming passthrough (WebSocket / `Transfer-Encoding: chunked`)
3. Real client IP via `X-Real-IP` / `X-Forwarded-For`

## Nginx

```nginx
server {
  listen 443 ssl http2;
  server_name api.example.com;

  ssl_certificate     /etc/letsencrypt/live/api.example.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

  add_header Strict-Transport-Security "max-age=63072000" always;

  real_ip_header X-Forwarded-For;
  set_real_ip_from 0.0.0.0/0;
  real_ip_recursive on;

  location / {
    proxy_pass http://127.0.0.1:3000;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_buffering off;
    proxy_cache off;
    proxy_read_timeout 300s;
    chunked_transfer_encoding on;
  }
}
```

## Caddy

`Caddyfile`:

```
api.example.com {
  reverse_proxy 127.0.0.1:3000 {
    header_up X-Real-IP {remote_host}
    header_up X-Forwarded-For {remote_host}
    header_up X-Forwarded-Proto {scheme}
    flush_interval -1
  }
}
```

Caddy auto-issues Let's Encrypt and adds HSTS — no manual config needed.

## Cloudflare

If you use Cloudflare for DNS + CDN:

1. SSL/TLS → Full (strict)
2. SSL/TLS → Edge Certificates → Always Use HTTPS → On
3. Rules → Transform Rules → set Host header / X-Forwarded-For
4. Add Cloudflare IP ranges to `set_real_ip_from` in Nginx (https://www.cloudflare.com/ips/)

Note: Cloudflare's default Cache Level (Standard) skips `/v1/*` by design — no extra rule needed.

## Verify

After deploy:

- [ ] `curl -I https://api.example.com/` → 200 + `Strict-Transport-Security` header
- [ ] `curl -H "X-Forwarded-For: 1.2.3.4" https://api.example.com/api/log/self` → log shows `1.2.3.4`
- [ ] Streaming: `curl -N https://api.example.com/v1/chat/completions ...` receives SSE chunks

## Related

- [Docker Single-instance](./docker-deploy)
- [System Requirements](./requirements)