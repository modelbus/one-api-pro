---
title: Redemption Overview
description: What redemption codes are, how they're issued, and how users redeem them.
category: redemption
order: 1
---

# Redemption Overview

> Admin generates in bulk → distribute to users → users redeem at `/redeem` → quota arrives.

## What is a redemption code

A "voucher". Looks like `ABCD-1234-EFGH-5678-...` (24 chars).

Each code is bound to:

- A quota amount (added directly to the user balance)
- Or a plan (redeeming starts a subscription)
- Redemption count (default 1)
- Expiry

## Use cases

- **New user promotion**: send a 50-credit code after sign-up
- **Community**: drop codes in a community
- **Customer support**: small compensation for an issue
- **Partner channel**: bulk codes for the partner to distribute

## Three code types

When creating in the admin:

| Type | What the user gets on redeem |
|---|---|
| **Quota code** | [User balance](../schema/redemption) += N (converted via `QuotaPerUnit`) |
| **Plan code** | A new [Subscription](../schema/subscription) with the plan's `duration_days` |
| **Free-trial code** | Internal "limited free" marker (in development) |

## Full lifecycle

```
Admin creates in bulk
    │
    ▼
Export (CSV / TXT)
    │
    ▼
Distribute (community / email / event)
    │
    ▼
User enters at /redeem
    │
    ▼
Validate: not expired / not exhausted / not disabled
    │
    ▼
Credit: balance += N or activate subscription
```

## vs. other ways

| Method | What the user gets | Issued by |
|---|---|---|
| Redemption code | Balance / Subscription | Admin backend |
| Plan order | Subscription | User purchases |
| Top-up | Balance | User purchases |

Codes suit "bulk distribution"; orders / top-ups suit "user-initiated purchase".

## Related

- [Redemption Schema](../schema/redemption)
- [Redemption Management (admin)](./admin-guide)
- [Redemption Use (user)](./user-guide)