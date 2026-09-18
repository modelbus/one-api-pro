---
title: System Settings
description: "Site appearance, login / register policy, OAuth / Turnstile / SMTP / announcements / home content."
category: admin
order: 16
---

# System Settings

> Site-wide configuration: server address, login & registration switches, GitHub / Lark / WeChat OAuth, Turnstile, SMTP, appearance, announcements. UI: `web/default-pro/src/views/setting/SystemSetting.vue`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/option/` | `GET` | Root | List all options (`*Token` / `*Secret` keys are filtered out) |
| `/api/option/` | `PUT` | Root | Single-key persistence |

Implementation: `controller/option.go`. Backed by the `options` table + `config.OptionMap`.

## Sections

| Section | Fields |
|---|---|
| Server Address | `ServerAddress` (external API base) |
| Login / Register | `PasswordLoginEnabled` / `PasswordRegisterEnabled` / `RegisterEnabled` / `EmailVerificationEnabled` / `GitHubOAuthEnabled` / `LarkOAuthEnabled` / `WeChatAuthEnabled` / `TurnstileCheckEnabled` |
| GitHub OAuth | `GitHubClientId` / `GitHubClientSecret` |
| Lark OAuth | `LarkClientId` / `LarkClientSecret` |
| WeChat login | `WeChatServerAddress` / `WeChatServerToken` / `WeChatAccountQRCodeImageURL` |
| Turnstile | `TurnstileSiteKey` / `TurnstileSecretKey` |
| SMTP | `SMTPServer` / `SMTPPort` / `SMTPAccount` / `SMTPFrom` / `SMTPToken` |
| Appearance | `SystemName` / `Logo` / `Theme` |
| Content | `Notice` / `HomePageContent` |

## Field Semantics

| Field | Notes |
|---|---|
| `ServerAddress` | External API base URL shown on the Token page (`{{ baseUrl }}/v1`) |
| `PasswordLoginEnabled` | When off, login is only via OAuth / WeChat / Lark |
| `PasswordRegisterEnabled` | When off, signup must go through invitations or OAuth |
| `RegisterEnabled` | Master signup switch; `false` rejects `POST /api/user/register` outright |
| `EmailVerificationEnabled` | Require email + verification code on signup |
| `EmailDomainRestrictionEnabled` + `EmailDomainWhitelist` | Email-domain whitelist (configured under `/setting/operation`) |
| `GitHubOAuthEnabled` / `LarkOAuthEnabled` / `WeChatAuthEnabled` | Enable the OAuth login entries; Client ID / Secret must be filled first |
| `TurnstileCheckEnabled` | Force Cloudflare Turnstile on critical flows (signup, password reset) |
| `Notice` | Top banner (exposed via `GET /api/notice`) |
| `HomePageContent` | Landing-page hero copy (exposed via `GET /api/home_page_content`) |
| `SystemName` / `Logo` / `Theme` | Site name / logo / frontend theme (must be in `config.ValidThemes`) |

## Toggle Preconditions

`controller/option.go::UpdateOption` validates prerequisites before enabling certain keys — see [operation-setting](operation-setting).

## OAuth Setup

1. Create an OAuth app on the provider (GitHub / Lark / WeChat). Capture the Client ID and Secret.
2. Set the redirect URL to `{ServerAddress}/api/oauth/{github|lark|wechat}`.
3. Fill the Client ID / Secret on this page.
4. Toggle the matching `*Enabled` switch and save.

## SMTP

Provide server, port (587 by default), account, `From` address and token (an app-specific password is recommended). Once saved:

- Email-verification codes on signup;
- Password-reset emails;
- Invitation notifications;

are dispatched via this SMTP relay.

## Turnstile

1. Create a Turnstile widget in the Cloudflare dashboard; obtain a Site Key (frontend) + Secret Key (backend).
2. Enter both on this page and toggle `TurnstileCheckEnabled`.
3. Signup / password reset / redemption pages automatically embed the Turnstile challenge.

## Frontend Guide

- Top: single-input **Server Address** + Save button.
- **Login / Register** has eight switches; each saves on `@change`.
- Each OAuth / Turnstile / SMTP block has its own Save button.
- **Appearance** saves the three fields together.
- **Content** saves both textareas together.

## Implementation Pointers

| Concern | Location |
|---|---|
| Handler + validation | `controller/option.go` |
| Boot-time init / sync | `model/option.go::InitOptionMap` / `SyncOptions` |
| Public notice / home content | `controller/misc.go` (`GetNotice` / `GetHomePageContent`) |
| GitHub OAuth | `controller/auth/github.go` |
| Lark OAuth | `controller/auth/lark.go` |
| WeChat OAuth | `controller/auth/wechat.go` |
| Routes | `router/api.go` |