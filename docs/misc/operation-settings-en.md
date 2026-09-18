---
title: Operation Settings
description: "Quota policy, monitoring thresholds, error-strategy and invitation / registration bonuses."
category: misc
order: 15
---

# Operation Settings

> Site-wide operational knobs: registration / invitation / new-user quota, monitoring thresholds, channel routing and error strategy. UI: `web/default-pro/src/views/setting/OperationSetting.vue`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/option/` | `GET` | Root | All `OptionMap` entries except those whose key ends with `Token` or `Secret` |
| `/api/option/` | `PUT` | Root | Single `{ key, value }` save with per-key validation |

Implementation: `controller/option.go`. Persisted in the `options` table and mirrored into the in-memory `config.OptionMap`.

> `options` is a parallel mechanism to `system_settings`: `options` is a generic key/value table loaded into `OptionMap` at startup; toggles whose key ends with `Enabled` also flip the corresponding `config.*` runtime variable. `system_settings` is the newer category-grouped store maintained by `model.UpsertSystemSetting`. This page covers `options` only.

## Sections

The UI splits settings by purpose; switches usually save on `@change`, while numeric fields share a section-level **Save** button.

| Section | Fields |
|---|---|
| Quota policy | `QuotaForNewUser` / `PreConsumedQuota` / `QuotaForInviter` / `QuotaForInvitee` |
| Monitor | `ChannelDisableThreshold` / `QuotaRemindThreshold` + `AutomaticDisableChannelEnabled` / `AutomaticEnableChannelEnabled` / `LogConsumeEnabled` |
| Log pruning | Date picker + **Clean logs** button → `DELETE /api/log/?target_timestamp=…` |
| General | `TopUpLink` / `ChatLink` / `QuotaPerUnit` / `RetryTimes` + `DisplayInCurrencyEnabled` / `DisplayTokenStatEnabled` / `ApproximateTokenEnabled` |
| Channel routing | `ChannelDefaultCooldownSeconds` / `ChannelMaxCooldownSeconds` + `ChannelConcurrencyEnabled` / `ChannelStickySessionEnabled` |
| Error strategy | `ErrorNext` JSON: `{ passthrough, retry, disable, cooldown }` |
| Plan | `plan.upgrade_mode` (see [plan-setting](plan-setting)) |

## Field Semantics

| Field | Default | Notes |
|---|---|---|
| `QuotaForNewUser` | env | New-user signup bonus (quota) |
| `PreConsumedQuota` | 0 | Pre-deducted per request; refunded on failure |
| `QuotaForInviter` | 0 | Bonus to the inviter when an invitee signs up |
| `QuotaForInvitee` | 0 | Bonus to the invitee when signing up |
| `ChannelDisableThreshold` | 0 | Auto-disable channel after N consecutive failures (requires `AutomaticDisableChannelEnabled=true`) |
| `QuotaRemindThreshold` | 0 | UI low-balance alert |
| `QuotaPerUnit` | env | Currency conversion: `QuotaPerUnit` quota = 1 CNY |
| `RetryTimes` | env | Per-channel retry attempts |
| `ChannelDefaultCooldownSeconds` / `ChannelMaxCooldownSeconds` | env | Channel cooldown window |
| `ErrorNext` | `{passthrough:true,retry:true,disable:true,cooldown:true}` | Relay error-chain strategy (see below) |

`TopUpLink` / `ChatLink` are the landing-page CTA buttons.

## ErrorNext Strategy

Stored as `{"errived": bool, "retry": bool, "disable": bool, "cooldown": bool}`:
- `passthrough` — surface the error to the client without retries.
- `retry` — retry up to `RetryTimes` within the same channel.
- `disable` — disable the channel after hitting the failure threshold.
- `cooldown` — apply `ChannelDefaultCooldownSeconds` after a failure.

The relay engine reads the JSON; the UI just persists the four booleans.

## Toggle Precondition Checks

`PUT /api/option/` rejects enabling the following keys until their preconditions are met (`controller/option.go`):

| Key | Prerequisite |
|---|---|
| `Theme` | Value must be in `config.ValidThemes` |
| `GitHubOAuthEnabled` | `GitHubClientId` and `GitHubClientSecret` both non-empty |
| `EmailDomainRestrictionEnabled` | `EmailDomainWhitelist` non-empty |
| `WeChatAuthEnabled` | `WeChatServerAddress` non-empty |
| `TurnstileCheckEnabled` | `TurnstileSiteKey` non-empty |

Failures return actionable error messages.

## Frontend Guide

- Welcome bar + loading spinner; the page loads via `GET /api/option/` on mount.
- Numeric fields use `<a-input-number size="large">`; each section has its own **Save** button that PUTs all the section's keys.
- Most toggles call `saveSwitch(key)` immediately on `@change`.
- **Log pruning** sends `DELETE /api/log/?target_timestamp=<unix-seconds>`; `target=0` is rejected by the backend.
- **Error strategy** is saved as a single JSON string with the four booleans.

## Implementation Pointers

| Concern | Location |
|---|---|
| Handler + validation | `controller/option.go` |
| `OptionMap` bootstrap | `model/option.go::InitOptionMap` / `SyncOptions` |
| Runtime sync | `model/option.go::UpdateOption` → `updateOptionMap` |
| Relation with `system_settings` | See header comment in `model/system_setting.go` |
| Routes | `router/api.go` |