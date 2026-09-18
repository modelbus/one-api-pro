---
title: Plan Admin
description: "Plan CRUD, publish toggle, features and model-limits management."
category: admin
order: 5
---

# Plan Admin

> CRUD for plans (subscriptions), publish toggle, and per-model rate limits; served by `web/default-pro/src/views/setting/PlanSetting.vue`.

## Data Model

`model.Plan` (table `plans`):

| Field | Type | Notes |
|---|---|---|
| `name` | `varchar(100)` | Display name |
| `description` | `text` | Subtitle shown on the cards |
| `price` | `decimal(10,2)` | Price in CNY |
| `duration_days` | `int` | Validity length in days |
| `duration_text` | `varchar(50)` | Free-form label ("月卡", "季卡", …) |
| `status` | `int` | `1=published` / `0=hidden` |
| `recommended` | `bool` | Highlight flag |
| `sort` | `int` | Display order (ascending) |
| `features` | `text` (JSON) | Feature bullets (`StringSlice`, multi-row editor in the UI) |
| `model_limits` | `text` (JSON) | `{ "gpt-4": { request_period, token_period, period_h, … }, … }` |
| `default_model` | `varchar(100)` | Default model for the subscription; rejected by `ValidateDefaultModel()` unless it appears in `model_limits` |

`Plan.ValidateDefaultModel()` returns "default_model 'X' is set but model_limits is empty" or "... not found in model_limits" when the default is not covered.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/plan/list` | `GET` | Public | All published plans (user-facing) |
| `/api/plan/detail/:id` | `GET` | Public | Single plan detail |
| `/api/plan/` | `GET` | Admin | Paginated full list (`config.ItemsPerPage`) |
| `/api/plan/search?keyword=` | `GET` | Admin | `name LIKE keyword%` |
| `/api/plan/:id` | `GET` | Admin | Single plan |
| `/api/plan/` | `POST` | **Root** | Create (`name` required) |
| `/api/plan/` | `PUT` | **Root** | Update with `ValidateDefaultModel()` |
| `/api/plan/:id` | `DELETE` | **Root** | Delete |
| `/api/plan/current` | `GET` | User | Current active subscription; flattens the embedded `Plan` fields to the top level |

Implementation: `controller/plan.go`; frontend uses `@/api/plan.js`.

## Publish Toggle

The **Enable / Disable** action on each row issues `PUT /api/plan/` with `status` flipped between `1` and `0`. Plans with `status=0` are filtered out by `GetPublicPlans`, so they no longer appear to self-service buyers.

## Features / Model Limits Editing

- `features`: `StringSlice` exposed to the UI as a JSON array. `PlanSetting.vue` uses a dynamic-row `<a-input>` editor (`featuresFromRecord` / `sanitizeFeaturesList` helpers in `web/default-pro/src/utils/plan.js`).
- `model_limits`: a JSON object stored as raw JSON. Example:

```json
{
  "gpt-4o": {
    "period_h": 5,
    "request_period": 100,
    "request_week": 1000,
    "request_month": 5000,
    "token_period": 50000,
    "token_week": 500000,
    "token_month": 2000000
  },
  "claude-3.5-sonnet": {
    "period_h": 5,
    "request_period": 50,
    "token_period": 20000,
    "request_week": 500,
    "token_month": 2500
  }
}
```

Field semantics:
- `period_h` defines the `WindowTypePeriod` window size (in hours).
- `request_period/week/month` = per-window request caps; `token_*` works the same way for tokens.

## Frontend Guide

- Table columns: ID / name / price / duration_days / default_model / recommended ★ / status tag / sort / actions.
- **Add Plan** opens a 900 px modal (`arco-modal-plan-wide`):
  1. Name, description (textarea);
  2. Price (precision 2), `duration_days`, `duration_text`;
  3. Sort, status (enabled/disabled), recommended switch;
  4. Features editor (dynamic rows, Enter to add, per-row delete);
  5. Default model (must already be in `model_limits`);
  6. Model limits (JSON textarea; placeholder `{"gpt-4": 1000, "gpt-3.5": 5000}`).
- The **Enable / Disable** button commits in place; **Delete** uses a confirmation popover.
- Plans are sorted by `sort` ascending — set `sort` low for the cards you want on top.

## Implementation Pointers

| Concern | Location |
|---|---|
| CRUD handler | `controller/plan.go` |
| Public fetches | `controller/plan.go::GetPublicPlans` / `GetPublicPlanDetail` |
| Current subscription | `controller/plan.go::GetCurrentPlan` → `model.GetActiveUserPlansByUserId` |
| `StringSlice` bridge | `model/plan.go` |
| Limit structure | `model/plan.go::ModelLimitRule` |
| Frontend edit helpers | `web/default-pro/src/utils/plan.js` |
| Routes | `router/api.go` |