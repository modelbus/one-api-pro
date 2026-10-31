-- Migration: remove legacy proxy channel type (43)
-- The Proxy channel type forwarded arbitrary HTTP methods, paths and bodies to
-- a channel-provided BASE_URL. That made /v1/oneapi/proxy an unrestricted SSRF
-- gateway into loopback and private networks. Both the adaptor and the route
-- are now removed; existing rows are disabled so they are never routed again.

UPDATE channels SET status = 2 WHERE type = 43 AND status = 1;
