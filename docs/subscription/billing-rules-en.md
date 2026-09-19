---
title: Billing Rules
description: Subscription window quotas, deduction formula, and billing details.
category: subscription
order: 4
---

# Billing Rules

> How subscription quota is deducted and when the plan is considered exhausted.

## Three windows

Each plan can set three window quotas per model:

| Window | Length | Meaning |
|---|---|---|
| `period` | 5 hours (default) | Short rolling window |
| `week` | 7 days | Weekly quota |
| `month` | 30 days | Monthly quota |

**If any window is exhausted, the plan is considered used up.** E.g. if monthly is gone, calls fall back to pay-as-you-go even if weekly still has room.

## Weighted usage

System computes a "usage" for each window:

```
usage = Σ (consumed_in_window / quota) × 100
```

Any window ≥ 100% → plan exhausted.

## Counting rules

Each plan/model can count by one of two dimensions:

- **By request count**: calls per day / week / month
- **By token total**: tokens per day / week / month

Which one is used depends on the plan setting ([Plan Management](../subscription/plan-management)).

## Per-call deduction order

1. **Check subscription**: does one of your active subscriptions cover this model?
   - Yes → deduct within the subscription window; **no balance charged**
   - No → fall to step 2
2. **Pay-as-you-go**: deduct from [User balance](../schema/user)
   - Enough → call succeeds
   - Not enough → 429 reject
3. **Cache**: subscription state is cached for 5 minutes; updates invalidate.

## Plan expiry

- Subscription auto-expires when the period ends.
- Calls still work, but billing falls back to pay-as-you-go.
- To regain the discount: **resubscribe**.

## FAQ

- **"Quota remaining" jumps**: subscription cache is 5 minutes; restart or invalidate to refresh immediately.
- **Just topped up but no balance deducted**: check whether the call hit a subscription; subscriptions are deducted first.
- **Monthly quota exhausted but weekly still has room**: yes, each window is independent; any one triggers pay-as-you-go.

## Related

- [Subscription Overview](./overview)
- [Subscription Schema](../schema/subscription)
- [Plan Schema](../schema/plan)
- [Plan Settings](./plan-settings)