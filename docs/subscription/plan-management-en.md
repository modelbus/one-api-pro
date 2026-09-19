---
title: Plan Management (Admin)
description: How admins create, publish, unpublish, and delete plans.
category: subscription
order: 6
---

# Plan Management (Admin)

> Admin → Plans.

## What a Plan is

A Plan is the "product" template users see on the subscribe page. Each plan defines:

| Field | Meaning |
|---|---|
| Name | What users see |
| Price | ¥ |
| Quota | Calls or tokens (unit price set in [Model Price](../schema/model-price)) |
| Validity | Number of days |
| Discount multiplier | Discount applied during the subscription (e.g. 0.8) |
| Description | Rich text shown publicly |
| Status | Published / Unpublished |
| Recommended | Toggle: shows a "Recommended" badge |

## Create a plan

Admin → Plans → Add:

1. Name, price, validity
2. Total quota (e.g. 5M tokens / month)
3. Per-model window quotas (optional; empty = unlimited for that model)
4. Discount multiplier (default 1.0 = no discount)
5. Description (rich text)
6. Save

## Publish / Unpublish

- **Publish**: users see this plan on the subscribe page
- **Unpublish**: users can't start new subscriptions; existing ones are unaffected

Typical flow: develop → test → publish. Trouble? Unpublish immediately.

## Delete

Don't delete directly (breaks existing subscriptions). Recommended flow:

1. Unpublish (stops new subscriptions)
2. Wait for existing subscriptions to expire (naturally)
3. Then delete

## Impact of editing

| Changed | Effect |
|---|---|
| Price | Existing subs use the old price |
| Validity | Existing subs use the old validity |
| Discount multiplier | Existing subs use the old multiplier |
| Status (unpublish) | Existing subs continue |

**Plans never retroactively update active subscriptions.**

## FAQ

- **New plan but users don't see it**: check it's published; verify [Plan Settings](./plan-settings) didn't disable it.
- **Plan price change — existing users charged at the new rate?** No, they keep the old rate until they re-subscribe.
- **Delete blocked by "still has active subscriptions"**: wait for them to expire, then delete.

## Related

- [Plan Schema](../schema/plan)
- [Plan Pricing (concept)](../pricing/plan)
- [Plan Settings](./plan-settings)
- [Subscription Management (admin)](./admin-guide)