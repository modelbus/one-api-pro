---
title: Redemption Use (User)
description: How to redeem a code as a user.
category: redemption
order: 3
---

# Redemption Use (User)

> Got a code — here's how to use it.

## How to redeem

1. Log in and visit `/redeem`
2. Paste the 24-char code
3. Click "Redeem"
4. Success → balance credited / subscription activated

You can also find the redemption form in [Personal Center](./user-guide).

## After success

- **Quota code**: your [User balance](../schema/user) += N
- **Plan code**: a new [Subscription](../schema/subscription) is activated; expires at `now + plan.duration_days`

Check [Dashboard](../user/dashboard) for balance; visit `/subscription` to see the new plan.

## On failure

| Message | Meaning |
|---|---|
| Code not found | Typed wrong; re-check |
| Expired | Past `expiry_at`; contact the issuer |
| Exhausted | Past `max_redemptions`; contact the issuer |
| Disabled | Admin disabled it; contact the issuer |
| Already redeemed by your account | Single-use code per account |

## FAQ

- **Code pasted but still fails**: check leading/trailing whitespace; copy-paste sometimes adds invisible characters
- **Quota not arrived**: refresh the page, or log out and back in
- **Re-entering the same code**: a single-use code rejects with "already redeemed by your account"

## Related

- [Redemption Overview](./overview)
- [Redemption Management (admin)](./admin-guide)