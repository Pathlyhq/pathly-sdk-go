# Changelog

## 0.1.0 — 2026-09-23

- First publish: Go client for Pathly monitoring HTTP scenarios, webhooks,
  maintenance windows and SLA targets.
- `Idempotency-Key` on creates, retries (`Retry-After`, max 4, 90 s cap).
- `Ping` via `GET /v1/usage` (403 accepted).
- Stdlib only.
