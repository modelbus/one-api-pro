-- Migration: Add fallback channel support
-- This migration marks channels that should be reserved exclusively for fallback
-- (i.e. only used when all normal channels for the requested model are exhausted).

ALTER TABLE channels ADD COLUMN IF NOT EXISTS is_fallback BOOLEAN DEFAULT 0;
ALTER TABLE channels ADD COLUMN IF NOT EXISTS fallback_priority BIGINT DEFAULT 0;