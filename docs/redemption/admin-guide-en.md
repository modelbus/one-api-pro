---
title: Redemption Management (Admin)
description: Bulk generate, switch state, copy, and query redemption codes.
category: redemption
order: 2
---

# Redemption Management (Admin)

> Admin → Redemption Codes.

## Generate

1. Admin → Redemption Codes → "Bulk generate"
2. Pick type (Quota / Plan)
3. Fill value (quota amount / pick a plan)
4. Max redemptions (default 1)
5. Expiry (default 30 days / custom)
6. Quantity (e.g. 100)
7. "Generate"

Codes appear in the list immediately. **The codes are shown only once** — export CSV / TXT right away.

## Export

Top-right "Export":

- CSV: each line `key,type,value,expires_at`
- TXT: one code per line

Send the file to the community / support team / partner.

## Search

Top search bar:

- Exact match on key
- Filter by type (Quota / Plan)
- Filter by status (Enabled / Disabled / Exhausted / Expired)

## Disable

1. Row → "Disable"
2. Double-confirm → immediately invalid; subsequent redempts are rejected

Common for: leaking code → emergency lock.

## Delete

1. Row → "Delete"

> Delete is irreversible and loses audit trail. Prefer "Disable".

## FAQ

- **Lost the codes after generation**: codes are shown once at generation. Export CSV / TXT immediately.
- **User says code is invalid**: check expired / disabled / exhausted; the admin can look up the code's status.
- **Generated 100 codes, only 30 redeemed**: remaining 70 still work. To kill them all, disable one by one (consider a script / SQL).

## Related

- [Redemption Overview](./overview)
- [Redemption Schema](../schema/redemption)
- [Redemption Use (user)](./user-guide)