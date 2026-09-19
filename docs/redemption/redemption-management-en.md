---
title: Redemption Management
description: Bulk generate, switch state, and query redemption codes.
category: redemption
order: 9
---

# Redemption Management

> Bulk generate codes, search by keyword. Admin → Redemption Codes.

## List columns

Each row shows one code (or an aggregated batch row):

- Code value (shown once at creation; afterwards masked `ABCD-****-1234`)
- Type (Quota / Plan)
- Quota amount / Plan name
- Max redemptions / Redeemed
- Expiry
- Status (Enabled / Disabled / Exhausted / Expired)

## Generate a batch

See [Redemption Management (admin)](./admin-guide).

## Search

Top search bar + type / status filters.

## Disable / Delete

Per-row actions. Prefer "Disable" (reversible) over "Delete" (irreversible, loses audit trail).

## Export

Top-right "Export": CSV or TXT. Export immediately after generation; code values don't persist in the list.

## FAQ

- **CSV opens as garbled text in Excel**: import via Data → From Text/CSV with UTF-8.
- **Generated codes never redeemed**: check the exported file actually contains the `key` column.
- **Bulk kill**: disable one by one; or ask the admin to run SQL.

## Related

- [Redemption Overview](./overview)
- [Redemption Management (admin)](./admin-guide)
- [Redemption Use (user)](./user-guide)