# 🛰️ star_ad_proxy

Tiny webhook proxy for [Subgram](https://api.subgram.org/api-docs). Subgram allows **only one webhook URL per account**, but a single account can host multiple bots with their own `bot_id` and `Api-Key`. This proxy takes the firehose, splits each batch by `bot_id`, and forwards every per-bot group to its own target — with that bot's own key stamped on the request.

```
subgram → POST /webhook/subgram → [split by bot_id] → bot A target
                                                    → bot B target
```

## 🚀 Run it

```bash
task dev:infra   # Dragonfly
task dev         # proxy on :8080
```

Config lives in `config/example.jsonc` (copy to `local.json`). Env overrides map `_` → `.` via koanf: `DRAGONFLY_ADDR=…`, `ADMIN_TOKEN=…`, etc.

## 🔌 API

| Method | Path | Auth | What |
|---|---|---|---|
| `POST` | `/webhook/subgram` | — | inbound from subgram, fan-out |
| `POST` | `/register` | `Bearer ADMIN_TOKEN` | upsert a route |
| `GET`  | `/register` | same | list routes (api_key masked) |
| `DELETE` | `/register/:bot_id` | same | drop a route |

Register:
```bash
curl -X POST :8080/register -H 'Authorization: Bearer change-me' -H 'Content-Type: application/json' \
  -d '{"bot_id":7213410106,"api_key":"sk_live_xxx","target_url":"https://svc.example.com/webhook/subgram"}'
```

Responses from `/webhook/subgram`:
- ✅ all good → `200 {"status":"ok","forwarded":N}`
- ⚠️ any group failed → `502 {"status":"partial","error":"…","results":[…]}` → subgram redelivers
- 🤷 unregistered `bot_id` → silently skipped + WARN log (redelivery wouldn't help)

## 🧠 One thing to know

On any forward failure we return `502` so subgram redelivers the **whole** batch. That means **downstream consumers must be idempotent on `webhook_id`** — duplicates will happen.

## 🗄️ Storage

Dragonfly: `proxy:routes:<bot_id>` (HASH) + `proxy:routes:index` (SET). No TTL, registration is boot-time.

## 🛠️ Stack

- 🦊 [**Fiber v3**](https://gofiber.io) — HTTP server *and* outbound client
- 🐉 [**Dragonfly**](https://www.dragonflydb.io) via [`go-redis`](https://github.com/redis/go-redis) — route table
- 🌀 [**`sourcegraph/conc`**](https://github.com/sourcegraph/conc) — typed goroutine pool for the fan-out
- 📜 [**`knadh/koanf/v2`**](https://github.com/knadh/koanf) — JSON file + env overlay config
- 🪵 stdlib `log/slog` — structured logging

No SQL DB, no message broker, no cron.

MIT.
