---
title: Channel Test
description: When and how to test a channel's connectivity.
category: channel
order: 4
---

# Channel Test

> When to use it, which button to click, how to read the result.

## When to use

- Right after creating a channel: test before enabling.
- After a channel has been auto-disabled: running a test re-enables it on success.
- When users report "this model isn't working": isolate whether it's the channel.
- After editing fields on [Add a Channel](./add-channel) (Base URL, model list).

## Two kinds of test

### Single-channel test

Admin → Channels → click a channel → "Test" button.

- **What it does**: sends a minimal request using the channel's actual credentials against one of its allowed models.
- **Time**: typically < 5 s.
- **Outcome**:
  - **Pass**: channel is usable; if it was auto-disabled, it's re-enabled.
  - **Fail**: shows the specific error (401 / 404 / 500 etc.) — diagnose from there.

### Batch test

Admin → Channels → "Batch test" at the top of the list.

- **What it does**: tests every enabled channel.
- **Time**: ~1–2 min for 50 channels.
- **Use**: weekly cron to catch channels that silently went bad.
- **Outcome**: full results on the "Test log" page.

## Which model to test with

The test model must be in the channel's model allow-list. Safe choices:

- OpenAI: `gpt-4o-mini` (cheap, almost always works)
- Anthropic: `claude-haiku-4-5`
- Others: pick the "basic" model from the provider's docs

## Error troubleshooting

| Error | Meaning | Fix |
|---|---|---|
| 401 Unauthorized | Wrong / expired credential | Re-issue the API key in the upstream console and update the channel. |
| 404 Not Found | Wrong Base URL | Cross-check against the provider's docs. |
| 400 Bad Request | Format mismatch | Check the Provider type (OpenAI vs Anthropic). |
| 429 Too Many Requests | Upstream rate limit | Lower the concurrency cap or try later. |
| 500 Server Error | Upstream outage | Retry later; pause the channel if it persists. |
| Timeout | Network | Check if the Base URL is reachable; bump the timeout. |

## Related

- [Channel Overview](./overview)
- [Add a Channel](./add-channel)
- [Channel Routing](./channel-routing)