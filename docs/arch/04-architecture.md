# Architecture: Layer Structure and Data Flow

## High-Level Layers

```
┌─────────────────────────────────────────────────────────────┐
│  cmd/                    (entrypoints: server, client)      │
├─────────────────────────────────────────────────────────────┤
│  server/rest/            (HTTP routing, middleware stack)   │
├─────────────────────────────────────────────────────────────┤
│  logic/rest/             (handlers: business logic)         │
├─────────────────────────────────────────────────────────────┤
│  internal/generated/                                              │
│  ├── oapi/rest/          (OpenAPI: types, interface, server)│
│  └── rdb/                (sqlc: Queries, models)            │
├─────────────────────────────────────────────────────────────┤
│  store/                  (state: health, future: cache/DB)  │
│  └── memory/             (in-memory implementation)         │
├─────────────────────────────────────────────────────────────┤
│  logging/                (structured slog wrapper)          │
└─────────────────────────────────────────────────────────────┘
```

**Dependency Rule**: Outer layers depend on inner. Inner never imports outer.
- `cmd` → `server` → `logic` → `internal/generated` + `store`
- `store/memory` has no deps (portable)

## Component Responsibilities

### cmd/server.go
- Config loading (`config.yaml` → struct)
- DB connection (pgxpool) or memory store selection
- Logger initialization
- HTTP server lifecycle (graceful shutdown)
- Wires: `store` → `logic` → `server` → `http.Server`

### server/rest/register.go
- `http.ServeMux` composition
- Mounts `logic/rest` handler at `/`
- Adds request logging middleware (via `logging.Logger`)

### logic/rest/rest.go
- Implements `oapi.StrictServerInterface` (generated)
- **No DB calls directly** — delegates to repository (future) or store
- Current: only health/ping → uses `store/memory` + context IP extraction
- Middleware: `strictIPRecordHandler` injects client IP into context

### store/memory/
- `HealthCheckState` with mutex-protected state machine
- `GetHealthCheckState()` / `SetHealthCheckState()`
- **Interface**: none yet (implicit). Will define `HealthStore` interface when PostgreSQL impl added.

### internal/generated/rdb/
- `Queries` struct: all DB operations (CRUD per table)
- Models: `User`, `Device`, `Record`, `Credential` (pgx types)
- **No business logic** — pure data access

### internal/generated/oapi/rest/
- `StrictServerInterface`: all operations as methods
- Request/response types per endpoint
- `Handler()`: `http.Handler` constructor

## Data Flow Examples

### GET /healthcheck/ready
```
HTTP Request
    → server/rest (ServeMux)
    → logic/rest.handler.GetHealthcheckReady()
    → store/memory.GetHealthCheckState()
    → returns state → maps to 200/503/500
    → JSON response
```

### Future: POST /devices
```
HTTP Request
    → server/rest
    → logic/rest.handler.CreateDevice()
    → (validate) → rdb.Queries.CreateDevice(ctx, params)
    → pgxpool → PostgreSQL
    → returns Device → map to API response
    → JSON response
```

## Why This Layering?

| Layer | Responsibility | Testability |
|-------|----------------|-------------|
| `cmd` | Composition, lifecycle | Integration test only |
| `server` | HTTP concerns (routing, middleware) | `httptest` |
| `logic` | Business rules, orchestration | Unit test with mock store/DB |
| `generated` | Transport/data contracts | Compile-time verified |
| `store` | State persistence abstraction | Swap impl (memory ↔ postgres) |

**KISS rationale**: 
- No DI framework — explicit wiring in `cmd/server.go`.
- No interface explosion — define when second impl exists.
- Generated code as "contract" — logic implements, doesn't wrap.

## Current Gaps / Technical Debt

1. **No repository layer**: `logic` calls `rdb.Queries` directly (not yet). When added, place in `logic/repository/` or `store/rdb/`.
2. **Config → DB connection**: `config.yaml` has `database.type: memory|postgres` but `cmd/server.go` not fully wired. `store/memory` only.
3. **Context propagation**: `logging.Logger` passed manually. Consider `context.WithValue` for request-scoped logger.
4. **Error handling**: No unified error type. `logic` returns `error` → `oapi` maps to 500. Need `AppError` with code for proper HTTP status mapping.

## Unresolved / Open Questions

1. **Repository pattern location**: `logic/repository/` (domain) vs `store/rdb/` (infrastructure)? Preference: `store/rdb/` — keeps `logic` pure.
2. **Transaction management**: Who owns `pgx.Tx`? Repository methods accept `*Queries` (uses pool). For multi-table ops, need `WithTransaction(ctx, fn)` helper.
3. **Observability**: Metrics (Prometheus), tracing (OTel) — where to inject? Middleware in `server/rest` or `cmd`?
4. **Graceful shutdown**: `store/memory` has `StateShuttingDown`. PostgreSQL store needs connection drain.
5. **Multi-tenancy**: Current schema has `user_id` on devices. API will need user context (from auth). Where to enforce? Middleware → context → repository.