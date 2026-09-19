---
title: Balance Update
description: Periodically pulling upstream account balance and showing it on the channel list.
category: channel
order: 5
---

# Balance Update

## What it does

Periodically asks upstream providers "how much credit do I have left" and shows the number on the channel list — so you don't have to log in to each provider's dashboard.

## Which providers are supported

Any provider with a balance query API:

- OpenAI (Usage API)
- Azure (Subscription API)
- Most Chinese providers (Zhipu, DeepSeek …)

Channels without a balance API show "N/A" — they still work.

## How to enable

Admin → Channels → open a channel → bottom of the page → "Auto update balance" toggle → pick interval (default 1h) → save.

## Where you see it

Admin → Channels list: each channel has a "Balance" column.

- Green = plenty
- Yellow = low (below threshold)
- Red = almost gone

The threshold is set under System Settings → Reminder Threshold.

## Why it shows nothing

- The provider has no balance query API.
- Credential is wrong / expired (same symptom as [Channel Test] failure).
- Network issue (some providers are flaky from certain regions).
- The auto-update toggle is off.

## Related

- [Channel Overview](./overview)
- [Channel Test](./channel-test)