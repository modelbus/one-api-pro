---
title: Redemption Quota Calculation
description: How redemption codes credit quota and how plans are activated.
category: redemption
order: 4
---

# Redemption Quota Calculation

> How much balance does a code add? Depends on its type.

## Two types, different math

### Quota codes

- The code stores an integer
- Balance credited = integer × `QuotaPerUnit` (configured in System Settings)
- Default: 1 quota unit = ¥1

Example:
- Code quota value = 100,000
- `QuotaPerUnit` = 500,000 (default)
- 100,000 quota units credited (per the system's current rate)

Changing `QuotaPerUnit` in [System Settings](../misc/system-settings) affects how all future codes are credited.

### Plan codes

- The code stores a plan id
- Redeeming creates a Subscription with `expire_at = now + plan.duration_days`
- Discount multiplier uses the plan's current setting

Example:
- Code points to Plan A (¥30/month, 30 days, multiplier 0.8)
- Redeeming → 30 days of subscription with calls charged at 0.8×

## Expiry

Each code has `expire_at` (Unix seconds). Past that time the code is invalid regardless of redemption status.

Admin sets expiry when generating (default 30 days).

## FAQ

- **Want to give a user ¥5 balance**: create 5 quota codes, each 100,000 (= ¥5 worth at default rate)
- **Want to give a monthly plan instead**: pick the plan when generating the code; a quota code can't activate a plan
- **Code credit doesn't match the plan price**: the code credit is integer × `QuotaPerUnit`; it doesn't depend on plan prices

## Related

- [Redemption Overview](./overview)
- [Redemption Schema](../schema/redemption)
- [System Settings (admin)](../misc/system-settings)