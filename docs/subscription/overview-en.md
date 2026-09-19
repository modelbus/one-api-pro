---
title: Subscription Overview
description: "Three billing modes: subscription, pay-as-you-go, top-up quota."
category: subscription
order: 1
---

# Subscription Overview

> How your quota is deducted, how much subscriptions save you, and the relationship to balance.

## Three billing modes

One API Pro supports three billing modes in parallel. They share one account but work differently:

| Mode | What it is | What gets deducted |
|---|---|---|
| **Subscription (Token Plan)** | A fixed quota pool per week / month / day | Plan-internal pool — **not** account balance |
| **Pay-as-you-go** | Pay per call | Account balance (`users.quota`) |
| **Top-up quota** | Add credit to your balance | Provides the balance for pay-as-you-go |

A single user can hold **all three** at once.

## Deduction order

On every call:

1. **Check subscription first**: does one of your active subscriptions cover this model?
   - Yes → use the plan's quota (**no balance deducted**)
   - No → fall through to pay-as-you-go
2. **Pay-as-you-go**: deduct from your [User quota](../schema/user)
   - Enough → call succeeds
   - Not enough → 429 reject

In short:

> Subscription = quota pool inside the plan (doesn't touch your balance)
> Balance = your fallback money; once gone, you can't call any more

## Subscription vs Top-up

| Dimension | Subscription | Top-up |
|---|---|---|
| Stored in | [Subscription](../schema/subscription) row | [User.quota](../schema/user) field |
| Expires? | Yes, when the period ends | No, never |
| Discount? | Yes (plan multiplier) | No |
| Drawn first? | Yes | Only after subscription is exhausted |
| After expiry? | Falls back to pay-as-you-go | Keeps being valid |

## Upgrading

When you already have a subscription and buy a more expensive one:

- **Default (price-diff)**: pay only the difference; remaining quota is pro-rated.
- **Stack**: pay full price for the new plan; the old one keeps running.

Admin toggles this in [Plan Settings](./plan-settings).

## When a subscription expires

- It auto-expires.
- Calls still work, but billing falls back to pay-as-you-go (deducts balance).
- To regain the discount, **resubscribe**.

## Related

- [Subscription Schema](../schema/subscription)
- [Plan Schema](../schema/plan)
- [User Schema](../schema/user)
- [Upgrade & Downgrade](./upgrade-downgrade)
- [Billing Rules](./billing-rules)