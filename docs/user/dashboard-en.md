---
title: Dashboard
description: Your consumption, active subscription and recent orders at a glance.
category: user
order: 2
---

# Dashboard

> The first page you see after logging in — your consumption overview.

## Where

You land on `/dashboard` after login.

## What you see

Top to bottom:

1. **Welcome bar**: your name + role (User / Admin) + system version
2. **Four key metric cards**:
   - Total Tokens consumed
   - Total Requests made
   - Total Quota remaining
   - Current Subscription (name + expiry)
3. **Usage progress bar**: today's usage as a percentage of your daily cap (red over 80%)
4. **Trend chart**: requests in the last 7 days
5. **Top-N models**: which models you used most
6. **Recent activity**: latest calls (time / model / token)
7. **API Key overview**: your most recent API Key (masked)
8. **Quick links**: jump to Tokens / Redemption / Logs

Admin accounts also see an "Operations Dashboard" button.

## How to read it

- **Total Tokens / Total Quota**: how many resources you have left
- **Usage progress**: today's usage vs. your daily cap; if you have no subscription, no bar is shown
- **Current Subscription**: your active plan + expiry; expired plans show "Expired — renew"

If numbers don't match what you expected (consumed more than displayed), see [FAQ](#faq).

## FAQ

- **A few minutes delay in data**: log writes are async, usually < 1 min
- **"Plan expired" but calls still work**: expiry only removes the discount; calls still go through
- **Admin's "Total Quota" looks wrong**: an admin may have adjusted the quota; check Admin → Users → detail

## Related

- [Access Token](./access-token)
- [Profile](./profile)
- [My Orders](./orders)
- [Subscription (Token Plan)](../subscription/overview)