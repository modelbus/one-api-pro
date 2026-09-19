---
title: User Subscription Guide
description: See your subscription progress, days left, per-model usage.
category: subscription
order: 2
---

# User Subscription Guide

> Your view of the subscription: progress, days left, per-model usage.

## Where

After login, visit `/subscription`. The page lists every active subscription you have.

Each card shows:

- Plan name
- Expiry / days left
- Progress bars for each window (5h / 7d / 30d)
- Models covered by the subscription
- "View usage" button → drill-down page

## How to read the progress bars

Each subscription card shows up to 3 bars:

- **5-hour window** (period): used in the last 5 hours
- **7-day window**: this week
- **30-day window**: this month

If any bar nears 100% → it turns red. **If any hits 100%** → plan is exhausted; the next call falls back to pay-as-you-go.

## Renew / Upgrade

1. Click "Upgrade" on the subscription card or go to `/pricing`
2. System creates an order per the [Upgrade & Downgrade](./upgrade-downgrade) rules
3. After payment it takes effect immediately

You can also click "Cancel" on the card to end early (most plans don't refund).

## What happens when a subscription ends

- It auto-expires at the end of the period
- Calls still work, but billing falls back to [pay-as-you-go](./billing-rules)
- To regain the discount, resubscribe

## FAQ

- **Progress not full but calls get rejected**: maybe one window is at 100%. Click "View usage" for per-model detail.
- **Active subscription but still failing**: also check your balance ([Dashboard](../user/dashboard) → Total Quota). Balance might be empty.
- **Upgraded but nothing changed**: upgrade is async; usually effective within seconds.

## Related

- [Subscription Overview](./overview)
- [Billing Rules](./billing-rules)
- [Upgrade & Downgrade](./upgrade-downgrade)
- [Dashboard](../user/dashboard)