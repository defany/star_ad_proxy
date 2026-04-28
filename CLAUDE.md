# star_ad_proxy

Webhook proxy in front of Subgram's account-level webhook. Splits incoming batches by `bot_id` and synchronously forwards each per-bot sub-batch to its registered downstream consumer (typically `flamie`).

## Why

Subgram (https://api.subgram.org) only allows one webhook URL per account, but a single account can host multiple bots — each with its own `bot_id` and per-bot `Api-Key`. Downstream consumers want per-bot routing.

## Structure

```
cmd/app/                  — entry point
internal/
  app/                    — orchestration (lc-based) + DI
  app/di/                 — context-based Once dependency injection
  config/                 — koanf JSON+env loader
  controller/http/        — Fiber v3 handlers (webhook, register)
  service/                — business logic (route, forwarder)
  repo/route/redis/       — Dragonfly-backed route table
  model/                  — domain types
pkg/fiber_http/           — generic httpwrap.Post[B] helper around fiber/v3/client
config/                   — JSON config files (local.json gitignored except example)
```

## Stack

- HTTP: `gofiber/fiber/v3` (server) + `gofiber/fiber/v3/client` (outbound)
- Storage: Dragonfly via `redis/go-redis/v9`
- Concurrency: `sourcegraph/conc/pool` (NEVER `golang.org/x/sync/errgroup`)
- Config: `koanf/v2` (JSON file via `CONFIG_PATH` + env overlay; `DRAGONFLY_ADDR` → `dragonfly.addr`)
- Logging: `log/slog` via `goblin/slogx`
- DI/lifecycle: `goblin/inject`, `goblin/lc`
- Error wrapping: `goblin/errfmt.WithSource`

## Commands

```bash
task dev:infra   # start Dragonfly
task dev         # hot-reload via air
task lint        # golangci-lint
task prod        # docker compose prod
```

## Routing model

- Dragonfly key `proxy:routes:<bot_id>` — HASH `{target_url, api_key, created_at, updated_at}`
- Dragonfly key `proxy:routes:index` — SET of bot_ids (for `List()` without SCAN)
- Updates wrapped in `TxPipelined` to keep HASH ↔ index consistent.

## Critical invariant

Forwarding is **synchronous**. On any forward failure the proxy returns 502 to subgram so it redelivers the whole batch. **Downstream consumers MUST be idempotent on `webhook_id`** — a single failing target poisons the batch and triggers redelivery of all events including ones that already succeeded.

## API quick reference

```bash
# Register bot route
curl -X POST http://localhost:8080/register \
  -H 'Authorization: Bearer change-me-locally' \
  -H 'Content-Type: application/json' \
  -d '{"bot_id":5811111111,"api_key":"sk_test","target_url":"http://flamie:8080/webhook/subgram"}'

# Delete
curl -X DELETE http://localhost:8080/register/5811111111 \
  -H 'Authorization: Bearer change-me-locally'

# List
curl http://localhost:8080/register \
  -H 'Authorization: Bearer change-me-locally'

# Subgram-side webhook (no auth on this edge — proxy stamps per-bot Api-Key on forward)
curl -X POST http://localhost:8080/webhook/subgram \
  -H 'Content-Type: application/json' \
  -d '{"webhooks":[{"webhook_id":1,"bot_id":5811111111,"ads_id":123,"link":"x","user_id":1,"status":"subscribed","subscribe_date":"2026-04-28"}]}'
```

## Conventions

- Outbound HTTP: ALWAYS `fiber/v3/client.Client` + `pkg/fiber_http.Post[B]`. Never `net/http`.
- Goroutine pools: ALWAYS `conc/pool.NewWithResults[T]().WithCollectErrored()`. Never `errgroup`.
- Goroutine return values: `(T, error)` — NEVER encode success/failure as `OK bool` or `Error string` fields on result struct.
- Wrap external errors with `errfmt.WithSource(err)` at the boundary, not in every layer.
- No `fmt.Print*`, no `time.Sleep`, no `http.DefaultClient` (forbidigo).
